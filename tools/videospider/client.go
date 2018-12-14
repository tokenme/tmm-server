package videospider

import (
	"bytes"
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/levigross/grequests"
	"github.com/mkideal/log"
	"github.com/panjf2000/ants"
	"github.com/tokenme/tmm/common"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TaskVideo struct {
	Id        uint64
	Link      string
	VideoLink string
}

type Client struct {
	service             *common.Service
	config              common.Config
	proxy               *Proxy
	httpClient          *grequests.Session
	resolvers           map[string]Resolver
	TLSHandshakeTimeout time.Duration
	DialTimeout         time.Duration
	exitCh              chan struct{}
	canExitCh           chan struct{}
}

func NewClient(service *common.Service, config common.Config) *Client {
	ro := &grequests.RequestOptions{
		UserAgent:    "Mozilla/5.0 (Linux; U; Android 4.3; en-us; SM-N900T Build/JSS15J) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
		UseCookieJar: false,
	}
	c := &Client{
		service:    service,
		config:     config,
		proxy:      NewProxy(service.Redis.Master, config.ProxyApiKey),
		httpClient: grequests.NewSession(ro),
		resolvers:  make(map[string]Resolver),
		exitCh:     make(chan struct{}, 1),
		canExitCh:  make(chan struct{}, 1),
	}
	c.RegisterAll()
	return c
}

func (this *Client) Register(resolver Resolver) {
	this.resolvers[resolver.Name()] = resolver
}

func (this *Client) Get(link string) (info Video, err error) {
	for _, resolver := range this.resolvers {
		if resolver.MatchUrl(link) {
			info, err = resolver.Get(link)
			if err != nil {
				return info, err
			}
			var files []VideoLink
			for _, f := range info.Files {
				if f.Link == "" {
					continue
				}
				files = append(files, f)
			}
			info.Files = files
			if len(files) == 0 {
				return info, errors.New("no video file found")
			}
			return info, nil
		}
	}
	return info, errors.New("resolver not found")
}

func (this *Client) Save(v Video) error {
	if len(v.Files) == 0 {
		return errors.New("invalid video")
	}
	sorter := NewVideoSorter(v.Files)
	sort.Sort(sort.Reverse(sorter))
	db := this.service.Db
	_, _, err := db.Query(`INSERT INTO tmm.share_tasks (creator, title, summary, link, image, video_link, is_video, points, points_left, bonus, max_viewers) VALUES (0, '%s', '%s', '%s', '%s', '%s', 1, 5000, 5000, 5, 20)`, db.Escape(v.Title), db.Escape(v.Desc), db.Escape(v.Link), db.Escape(v.Cover), db.Escape(sorter[0].Link))
	return err
}

func (this *Client) StartUpdateVideosService() {
	updateVideoCh := make(chan struct{}, 1)
	this.UpdateVideos(updateVideoCh)
	for {
		select {
		case <-updateVideoCh:
			go this.UpdateVideos(updateVideoCh)
		case <-this.exitCh:
			log.Warn("ExitCh")
			close(updateVideoCh)
			this.canExitCh <- struct{}{}
			return
		}
	}
}

func (this *Client) Stop() {
	this.exitCh <- struct{}{}
	log.Info("Can Exit")
	<-this.canExitCh
}

func (this *Client) UpdateVideos(updateCh chan<- struct{}) error {
	log.Info("Start Update Videos")
	defer func() {
		time.Sleep(10 * time.Second)
		updateCh <- struct{}{}
	}()
	db := this.service.Db
	var wg sync.WaitGroup
	videoFetchPool, _ := ants.NewPoolWithFunc(10, func(req interface{}) error {
		defer wg.Done()
		task := req.(*TaskVideo)
		//log.Info("Updating:%s", task.Link)
		video, err := this.Get(task.Link)
		if err != nil {
			log.Error("Update:%s, Failed:%s", task.Link, err.Error())
			return err
		}
		if len(video.Files) == 0 {
			log.Error("Update:%s, Failed: no video found", task.Link)
			return errors.New("invalid video")
		}
		sorter := NewVideoSorter(video.Files)
		sort.Sort(sort.Reverse(sorter))
		task.VideoLink = sorter[0].Link
		//log.Info("Updated:%s, Video:%s", task.Link, task.VideoLink)
		return nil
	})
	var (
		startId uint64
		endId   uint64
		where   string
	)
	for {
		if startId > 0 {
			where = fmt.Sprintf(" AND id<%d", startId)
		}
		endId = startId
		rows, _, err := db.Query(`SELECT id, link FROM tmm.share_tasks WHERE online_status=1 AND is_video=1 AND video_updated_at<DATE_SUB(NOW(), INTERVAL 30 MINUTE)%s ORDER BY id DESC LIMIT 1000`, where)
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			break
		}
		var tasks []*TaskVideo
		for _, row := range rows {
			task := &TaskVideo{
				Id:   row.Uint64(0),
				Link: row.Str(1),
			}
			endId = task.Id
			wg.Add(1)
			videoFetchPool.Serve(task)
			tasks = append(tasks, task)
		}
		wg.Wait()
		var (
			val        []string
			offlineIds []string
		)
		for _, task := range tasks {
			if task.VideoLink != "" {
				val = append(val, fmt.Sprintf("(%d, '%s', NOW())", task.Id, db.Escape(task.VideoLink)))
			} else {
				offlineIds = append(offlineIds, strconv.FormatUint(task.Id, 10))
			}
		}
		if len(val) > 0 {
			log.Warn("Saving:%d Videos", len(val))
			_, _, err := db.Query(`INSERT INTO tmm.share_tasks (id, video_link, video_updated_at) VALUES %s ON DUPLICATE KEY UPDATE video_link=VALUES(video_link), video_updated_at=VALUES(video_updated_at)`, strings.Join(val, ","))
			if err != nil {
				log.Error(err.Error())
			}
		}
		if len(offlineIds) > 0 {
			log.Warn("Offline:%d Videos", len(val))
			_, _, err := db.Query(`UPDATE tmm.share_tasks SET online_status=-1 WHERE id IN (%s)`, strings.Join(offlineIds, ","))
			if err != nil {
				log.Error(err.Error())
			}
		}
		if endId == startId {
			break
		}
		startId = endId
	}
	return nil
}

func (this *Client) RegisterAll() {
	toutiao := NewToutiao(this)
	this.Register(toutiao)
	pearVideo := NewPearVideo(this)
	this.Register(pearVideo)
	weibo := NewWeibo(this)
	this.Register(weibo)
	miaopai := NewMiaoPai(this)
	this.Register(miaopai)
	krcom := NewKrcom(this)
	this.Register(krcom)
}

func (this *Client) GetHtml(link string, ro *grequests.RequestOptions) (string, error) {
	proxyUrl, _ := this.proxy.Get()
	if proxyUrl != nil {
		if ro == nil {
			ro = &grequests.RequestOptions{
				Proxies:             map[string]*url.URL{"https": proxyUrl},
				TLSHandshakeTimeout: this.TLSHandshakeTimeout,
				DialTimeout:         this.DialTimeout,
			}
		} else {
			ro.Proxies = map[string]*url.URL{"https": proxyUrl}
			ro.TLSHandshakeTimeout = this.TLSHandshakeTimeout
			ro.DialTimeout = this.DialTimeout
		}
	}
	resp, err := this.httpClient.Get(link, ro)
	if err != nil {
		this.proxy.Update()
		return "", err
	}
	return resp.String(), nil
}

func (this *Client) GetBytes(link string, ro *grequests.RequestOptions) ([]byte, error) {
	proxyUrl, _ := this.proxy.Get()
	if proxyUrl != nil {
		if ro == nil {
			ro = &grequests.RequestOptions{
				Proxies:             map[string]*url.URL{"https": proxyUrl},
				TLSHandshakeTimeout: this.TLSHandshakeTimeout,
				DialTimeout:         this.DialTimeout,
			}
		} else {
			ro.Proxies = map[string]*url.URL{"https": proxyUrl}
			ro.TLSHandshakeTimeout = this.TLSHandshakeTimeout
			ro.DialTimeout = this.DialTimeout
		}
	}
	resp, err := this.httpClient.Get(link, ro)
	if err != nil {
		this.proxy.Update()
		return nil, err
	}
	return resp.Bytes(), nil
}

func (this *Client) GetReader(link string, ro *grequests.RequestOptions) (*bytes.Reader, error) {
	buf, err := this.GetBytes(link, ro)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}
