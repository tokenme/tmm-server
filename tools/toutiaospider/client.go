package toutiaospider

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
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
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	REALTIME_NEWS_URL = "https://www.toutiao.com/api/pc/realtime_news/"
)

type RealtimeNewsResponse struct {
	Message string    `json:"message"`
	Data    []Article `json:"data"`
}

type Article struct {
	GroupId  string `json:"group_id"`
	ImageUrl string `json:"image_url"`
	OpenUrl  string `json:"open_url"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Url      string `json:"url"`
	Digest   string `json:"digest"`
	Markdown string `json:"-"`
	DateTime string `json:"-"`
}

type ArticleBaseInfo struct {
	Article ArticleInfo `json:"articleInfo"`
}

type ArticleInfo struct {
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Cover   string         `json:"coverImg"`
	ItemId  string         `json:"itemId"`
	SubInfo ArticleSubInfo `json:"subInfo"`
}

type ArticleSubInfo struct {
	Source string `json:"source"`
	Time   string `json:"time"`
}

func (this ArticleInfo) ToArticle() Article {
	imageUrl := this.Cover
	if strings.HasPrefix(imageUrl, "http://") {
		imageUrl = strings.Replace(imageUrl, "http://", "https://", -1)
	}
	return Article{
		GroupId:  this.ItemId,
		ImageUrl: imageUrl,
		Title:    this.Title,
		Author:   this.SubInfo.Source,
		DateTime: this.SubInfo.Time,
		Markdown: this.Content,
	}
}

type Crawler struct {
	proxy      *Proxy
	service    *common.Service
	config     common.Config
	httpClient *grequests.Session
}

func NewCrawler(service *common.Service, config common.Config) *Crawler {
	ro := &grequests.RequestOptions{
		UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.116 Safari/537.36",
		UseCookieJar: true,
	}
	return &Crawler{
		proxy:      NewProxy(service.Redis.Master, config.ProxyApiKey),
		service:    service,
		config:     config,
		httpClient: grequests.NewSession(ro),
	}
}

func (this *Crawler) ReloadHttpClient() {
	ro := &grequests.RequestOptions{
		UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.116 Safari/537.36",
		UseCookieJar: true,
	}
	this.httpClient = grequests.NewSession(ro)
}

func (this *Crawler) Run() error {
	log.Info("Toutiao Article Crawler start")
	articles, err := this.SearchRealtimeNews()
	if err != nil {
		log.Error(err.Error())
	}
	log.Warn("Finished %d articles in toutiao.com", len(articles))
	return nil
}

func (this *Crawler) SearchRealtimeNews() (articles []Article, err error) {
	log.Warn("SearchRealtimeNews")
	ro := &grequests.RequestOptions{
		UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.116 Safari/537.36",
		UseCookieJar: true,
	}
	proxyUrl, _ := this.proxy.Get()
	if proxyUrl != nil {
		ro = &grequests.RequestOptions{
			Proxies: map[string]*url.URL{"https": proxyUrl},
		}
	}
	resp, err := this.httpClient.Get(REALTIME_NEWS_URL, ro)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var searchResp RealtimeNewsResponse
	err = ljson.Unmarshal(resp.Bytes(), &searchResp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var ids []string
	for _, a := range searchResp.Data {
		if a.GroupId == "" {
			continue
		}
		link := fmt.Sprintf("https://m.toutiao.com%s", a.OpenUrl)
		article, err := this.getArticle(link)
		if err != nil || article.Markdown == "" {
			continue
		}
		articles = append(articles, article)
		ids = append(ids, a.GroupId)
	}
	if len(ids) == 0 {
		return nil, nil
	}
	db := this.service.Db
	rows, _, err := db.Query(`SELECT fileid FROM tmm.articles WHERE fileid IN (%s) AND platform=1`, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	idMap := make(map[string]struct{})
	for _, row := range rows {
		idMap[row.Str(0)] = struct{}{}
	}
	var val []string
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for _, a := range articles {
		if _, found := idMap[a.GroupId]; found {
			continue
		}
		newA, err := this.updateArticleImages(a)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		publishTime := time.Now()
		if a.DateTime != "" {
			t, err := time.ParseInLocation("2006-01-02 15:04:05", a.DateTime, loc)
			if err == nil {
				publishTime = t
			}
		}
		sortId := utils.RangeRandUint64(1, 1000000)
		val = append(val, fmt.Sprintf("(%s, '%s', '%s', '%s', '%s', '%s', '%s', '%s', 1, %d)", newA.GroupId, db.Escape(newA.Author), db.Escape(newA.Title), db.Escape(newA.Url), db.Escape(newA.ImageUrl), publishTime.Format("2006-01-02 15:04:05"), db.Escape(newA.Digest), db.Escape(newA.Markdown), sortId))
	}
	if len(val) > 0 {
		_, _, err := db.Query(`INSERT IGNORE INTO tmm.articles (fileid, author, title, link, cover, published_at, digest, content, platform, sortid) VALUES %s`, strings.Join(val, ","))
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
	}
	return articles, nil
}

func (this *Crawler) getArticle(link string) (article Article, err error) {
	ro := &grequests.RequestOptions{
		UserAgent:    "Mozilla/5.0 (Linux; U; Android 4.3; en-us; SM-N900T Build/JSS15J) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
		UseCookieJar: false,
	}
	proxyUrl, _ := this.proxy.Get()
	if proxyUrl != nil {
		ro = &grequests.RequestOptions{
			Proxies: map[string]*url.URL{"https": proxyUrl},
		}
	}
	resp, err := this.httpClient.Get(link, ro)
	if err != nil {
		return article, err
	}
	reader := bytes.NewReader(resp.Bytes())
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return article, err
	}
	var articleInfo ArticleBaseInfo
	doc.Find("script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.HasPrefix(s.Text(), "var BASE_DATA = {") {
			tmp := strings.TrimPrefix(s.Text(), "var BASE_DATA = ")
			tmp = strings.TrimSuffix(tmp, ";")
			re := regexp.MustCompile("(\n\\s+)")
			tmp = re.ReplaceAllString(tmp, "")
			re = regexp.MustCompile("(,shareInfo: {.*?})")
			tmp = re.ReplaceAllString(tmp, "")
			re = regexp.MustCompile("'([^']*)'")
			tmp = re.ReplaceAllString(tmp, "\"${1}\"")
			err := ljson.Unmarshal([]byte(tmp), &articleInfo)
			if err != nil {
				log.Error(err.Error())
			}
			return false
		}
		return true
	})
	articleInfo.Article.Content = html.UnescapeString(articleInfo.Article.Content)
	if articleInfo.Article.SubInfo.Source == "" || articleInfo.Article.ItemId == "" || articleInfo.Article.Title == "" || articleInfo.Article.Content == "" {
		return article, errors.New("no content")
	}
	a, err := wxarticle2md.ToAticle(articleInfo.Article.Content)
	if err != nil {
		return article, err
	}
	articleInfo.Article.Content = wxarticle2md.Convert(a)
	article = articleInfo.Article.ToArticle()
	article.Url = link
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
	uploadImagePool, _ := ants.NewPoolWithFunc(10, func(src interface{}) error {
		defer wg.Done()
		ori := src.(string)
		link, err := this.uploadImage(ori)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		imageMap.Store(ori, link)
		return nil
	})
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
	rows, _, err := db.Query(`SELECT id, title, digest, cover FROM tmm.articles WHERE published=0 AND platform=1 ORDER BY sortid LIMIT 1000`)
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
