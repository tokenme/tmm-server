package videospider

import (
	"bytes"
	"errors"
	"github.com/levigross/grequests"
	"github.com/tokenme/tmm/common"
	"net/url"
	"sort"
)

var Resolvers = make(map[string]Resolver)

func Register(resolver Resolver) {
	Resolvers[resolver.Name()] = resolver
}

type Client struct {
	service    *common.Service
	config     common.Config
	proxy      *Proxy
	httpClient *grequests.Session
}

func NewClient(service *common.Service, config common.Config) *Client {
	ro := &grequests.RequestOptions{
		UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.116 Safari/537.36",
		UseCookieJar: false,
	}
	return &Client{
		service:    service,
		config:     config,
		proxy:      NewProxy(service.Redis.Master, config.ProxyApiKey),
		httpClient: grequests.NewSession(ro),
	}
}

func (this *Client) Get(link string) (info Video, err error) {
	if len(Resolvers) == 0 {
		this.RegisterAll()
	}
	for _, resolver := range Resolvers {
		if resolver.MatchUrl(link) {
			return resolver.Get(link)
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
	_, _, err := db.Query(`INSERT INTO tmm.share_tasks (creator, title, summary, link, image, video_link, is_video, points, points_left, bonus, max_viewers) VALUES (0, '%s', '%s', '%s', '%s', '%s', 1, 5000, 5000, 10, 10)`, db.Escape(v.Title), db.Escape(v.Desc), db.Escape(v.Link), db.Escape(v.Cover), db.Escape(sorter[0].Link))
	return err
}

func (this *Client) RegisterAll() {
	toutiao := NewToutiao(this)
	Register(toutiao)
	pearVideo := NewPearVideo(this)
	Register(pearVideo)
	weibo := NewWeibo(this)
	Register(weibo)
	miaopai := NewMiaoPai(this)
	Register(miaopai)
	krcom := NewKrcom(this)
	Register(krcom)
}

func (this *Client) GetHtml(link string, ro *grequests.RequestOptions) (string, error) {
	proxyUrl, _ := this.proxy.Get()
	if proxyUrl != nil {
		if ro == nil {
			ro = &grequests.RequestOptions{
				Proxies: map[string]*url.URL{"https": proxyUrl},
			}
		} else {
			ro.Proxies = map[string]*url.URL{"https": proxyUrl}
		}
	}
	resp, err := this.httpClient.Get(link, ro)
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}

func (this *Client) GetBytes(link string, ro *grequests.RequestOptions) ([]byte, error) {
	proxyUrl, _ := this.proxy.Get()
	if proxyUrl != nil {
		if ro == nil {
			ro = &grequests.RequestOptions{
				Proxies: map[string]*url.URL{"https": proxyUrl},
			}
		} else {
			ro.Proxies = map[string]*url.URL{"https": proxyUrl}
		}

	}
	resp, err := this.httpClient.Get(link, ro)
	if err != nil {
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
