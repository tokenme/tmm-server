package qutoutiaospider

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	//"github.com/davecgh/go-spew/spew"
	"github.com/levigross/grequests"
	"github.com/mkideal/log"
	"github.com/panjf2000/ants"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/tools/ljson"
	"github.com/tokenme/tmm/tools/qiniu"
	"github.com/tokenme/tmm/utils"
	"github.com/yizenghui/wxarticle2md"
	"gopkg.in/russross/blackfriday.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	SUGGEST_URL = "http://api.1sapp.com/content/outList?cid=255&tn=1&page=%d&limit=50&user=temporary%d&show_time=%s&min_time=%s&content_type=1&dtu=200"
)

type SuggestResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    SuggestData `json:"data"`
}

type SuggestData struct {
	MinTime  int64     `json:"min_time"`
	MaxTime  int64     `json:"max_time"`
	Articles []Article `json:"data"`
}

type Article struct {
	Id       string   `json:"id"`
	Author   string   `json:"source_name"`
	Title    string   `json:"title"`
	Url      string   `json:"url"`
	Digest   string   `json:"introduction"`
	Covers   []string `json:"cover"`
	Markdown string   `json:"-"`
	DateTime string   `json:"publish_time"`
}

type Crawler struct {
	proxy   *Proxy
	service *common.Service
	config  common.Config
}

func NewCrawler(service *common.Service, config common.Config) *Crawler {
	return &Crawler{
		proxy:   NewProxy(service.Redis.Master, config.ProxyApiKey),
		service: service,
		config:  config,
	}
}

func (this *Crawler) Run() (int, error) {
	log.Info("QuToutiao Article Crawler start")
	articles, err := this.Suggest(0, 1, 0, 0)
	if err != nil {
		log.Error(err.Error())
	}
	log.Warn("Finished %d articles in qutoutiao.net", len(articles))
	return len(articles), nil
}

func (this *Crawler) Suggest(userTmp int64, page int, showTime int64, minTime int64) (articles []Article, err error) {
	ro := &grequests.RequestOptions{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.2 Safari/605.1.15",
	}
	proxyUrl, _ := this.proxy.Get()
	if proxyUrl != nil {
		ro.Proxies = map[string]*url.URL{"https": proxyUrl, "http": proxyUrl}
	}
	if page == 0 {
		page = 1
	}
	if userTmp == 0 {
		userTmp = time.Now().UnixNano()
	}
	var (
		sTime string
		mTime string
	)
	if showTime > 0 {
		sTime = strconv.FormatInt(showTime, 10)
	}
	if minTime > 0 {
		mTime = strconv.FormatInt(minTime, 10)
	}
	resp, err := grequests.Get(fmt.Sprintf(SUGGEST_URL, page, userTmp, sTime, mTime), ro)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var searchResp SuggestResponse
	err = ljson.Unmarshal(resp.Bytes(), &searchResp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var ids []string
	for _, a := range searchResp.Data.Articles {
		if a.Id == "" {
			continue
		}
		ids = append(ids, a.Id)
	}
	if len(ids) == 0 {
		log.Warn("NOT FOUND ANY qutoutiao ARTICLES")
		return nil, nil
	}
	log.Info("Got %d qutoutiao search result", len(ids))
	db := this.service.Db
	rows, _, err := db.Query(`SELECT fileid FROM tmm.articles WHERE fileid IN (%s) AND platform=2`, strings.Join(ids, ","))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	idMap := make(map[string]struct{})
	for _, row := range rows {
		idMap[row.Str(0)] = struct{}{}
	}

	var val []string
	for _, a := range searchResp.Data.Articles {
		if a.Id == "" {
			continue
		}
		if _, found := idMap[a.Id]; found {
			log.Warn("Found Article:%s ", a.Id)
			continue
		}
		article, err := this.getArticle(a)
		if err != nil || article.Markdown == "" {
			if err != nil {
				log.Error(err.Error())
			}
			continue
		}
		articles = append(articles, article)
		newA, err := this.updateArticleImages(article)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		publishTime := time.Now()
		if article.DateTime != "" {
			ts, err := strconv.ParseInt(article.DateTime, 10, 64)
			if err == nil {
				publishTime = time.Unix(ts/1000, 0)
			}
		}
		sortId := utils.RangeRandUint64(1, 1000000)
		var cover string
		if len(newA.Covers) > 0 {
			cover = newA.Covers[0]
		}
		val = append(val, fmt.Sprintf("(%s, '%s', '%s', '%s', '%s', '%s', '%s', '%s', 2, %d)", newA.Id, db.Escape(newA.Author), db.Escape(newA.Title), db.Escape(newA.Url), db.Escape(cover), publishTime.Format("2006-01-02 15:04:05"), db.Escape(newA.Markdown), db.Escape(newA.Digest), sortId))
	}
	if len(val) > 0 {
		_, _, err := db.Query(`INSERT IGNORE INTO tmm.articles (fileid, author, title, link, cover, published_at, content, digest, platform, sortid) VALUES %s`, strings.Join(val, ","))
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
	}
	return articles, nil
}

func (this *Crawler) getArticle(article Article) (Article, error) {
	log.Info("Fetching toutiao article: %s", article.Url)
	ro := &grequests.RequestOptions{
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.116 Safari/537.36",
	}
	proxyUrl, _ := this.proxy.Get()
	if proxyUrl != nil {
		ro.Proxies = map[string]*url.URL{"https": proxyUrl, "http": proxyUrl}
	}
	resp, err := grequests.Get(article.Url, ro)
	if err != nil {
		log.Error(err.Error())
		return article, err
	}
	reader := bytes.NewReader(resp.Bytes())
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return article, err
	}
	contentWrapper := doc.Find("div.article div.content")
	contentWrapper.Find("img").Each(func(i int, s *goquery.Selection) {
		if src, found := s.Attr("data-src"); found {
			s.SetAttr("src", src).RemoveAttr("data-src")
		}
	})
	h, err := contentWrapper.Html()
	if err != nil {
		return article, err
	}
	a, err := wxarticle2md.ToAticle(h)
	if err != nil {
		return article, err
	}
	article.Markdown = wxarticle2md.Convert(a)
	log.Info("Fetched qutoutiao article: %s", article.Url)
	return article, nil
}

func (this *Crawler) updateArticleImages(a Article) (Article, error) {
	reader := bytes.NewBuffer(blackfriday.Run([]byte(a.Markdown)))
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return a, err
	}
	var imageMap sync.Map
	var wg sync.WaitGroup
	uploadImagePool, _ := ants.NewPoolWithFunc(10, func(src interface{}) {
		defer wg.Done()
		ori := src.(string)
		link, err := this.uploadImage(ori)
		if err != nil {
			log.Error(err.Error())
			return
		}
		imageMap.Store(ori, link)
		return
	})
	for _, cover := range a.Covers {
		wg.Add(1)
		uploadImagePool.Serve(cover)
	}
	doc.Find("img").Each(func(idx int, s *goquery.Selection) {
		s.SetAttr("class", "image")
		if src, found := s.Attr("src"); found {
			wg.Add(1)
			uploadImagePool.Serve(src)
		} else {
			s.Remove()
		}
	})
	wg.Wait()
	var covers []string
	for _, cover := range a.Covers {
		if src, found := imageMap.Load(cover); found {
			covers = append(covers, src.(string))
		}
	}
	a.Covers = covers
	doc.Find("img").Each(func(idx int, s *goquery.Selection) {
		s.SetAttr("class", "image")
		if src, found := s.Attr("src"); found {
			if link, found := imageMap.Load(src); found {
				s.SetAttr("src", link.(string))
			} else {
				s.Remove()
			}
		} else {
			s.Remove()
		}
	})
	h, err := doc.Find("body").Html()
	if err != nil {
		return a, err
	}
	a.Markdown = h
	return a, nil
}

func (this *Crawler) uploadImage(src string) (string, error) {
	log.Info("Uploading image: %s", src)
	resp, err := http.DefaultClient.Get(src)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fn := base64.URLEncoding.EncodeToString([]byte(src))
	link, _, err := qiniu.Upload(context.Background(), this.config.Qiniu, this.config.Qiniu.ImagePath, fn, body)
	return link, err
}

func (this *Crawler) Publish() error {
	db := this.service.Db
	rows, _, err := db.Query(`SELECT id, title, digest, cover FROM tmm.articles WHERE published=0 AND platform=2 ORDER BY sortid LIMIT 1000`)
	if err != nil {
		return err
	}
	var ids []string
	var val []string
	for _, row := range rows {
		id := row.Uint64(0)
		title := row.Str(1)
		digest := row.Str(2)
		link := fmt.Sprintf("https://tmm.tokenmama.io/article/show/%d", id)
		cover := strings.Replace(row.Str(3), "http://", "https://", -1)
		ids = append(ids, fmt.Sprintf("%d", id))
		val = append(val, fmt.Sprintf("(0, '%s', '%s', '%s', '%s', 500, 500, 5, 20, 1)", db.Escape(title), db.Escape(digest), db.Escape(link), db.Escape(cover)))
	}
	if len(val) > 0 {
		_, _, err := db.Query(`INSERT INTO tmm.share_tasks (creator, title, summary, link, image, points, points_left, bonus, max_viewers, is_crawled) VALUES %s`, strings.Join(val, ","))
		if err != nil {
			return err
		}
		_, _, err = db.Query(`UPDATE tmm.articles SET published=1 WHERE id IN (%s)`, strings.Join(ids, ","))
		if err != nil {
			return err
		}
		log.Info("Published %d articles", len(ids))
	}
	return nil
}
