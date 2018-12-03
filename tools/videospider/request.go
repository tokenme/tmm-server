package videospider

import (
	"github.com/PuerkitoBio/goquery"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/levigross/grequests"
)

type Resolution = uint8

const (
	Normal Resolution = 1
	HD     Resolution = 2
	Low    Resolution = 3
)

type Video struct {
	Title      string      `json:"title"`
	Desc       string      `json:"desc"`
	Link       string      `json:"link"`
	Cover      string      `json:"cover"`
	Duration   int64       `json:"duration"`
	Files      []VideoLink `json:"files"`
	CreateTime int64       `json:"createTime"`
	Size       VideoSize   `json:"size"`
}

type VideoLink struct {
	Resolution Resolution `json:"resolution"`
	Link       string     `json:"link"`
}

type VideoSize struct {
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
	Size   uint64 `json:"size"`
}

type Resolver interface {
	MatchUrl(url string) bool
	Get(link string) (info Video, err error)
	Name() string
}

type Request struct {
	client   *Client
	name     string
	patterns []string
}

func (this *Request) Name() string {
	return this.name
}

func (this *Request) MatchUrl(url string) bool {
	if len(R1Of(this.patterns, url)) > 1 {
		return true
	}
	return false
}

func (this *Request) GetContent(url string, ro *grequests.RequestOptions) (string, error) {
	return this.client.GetHtml(url, ro)
}

func (this *Request) BuildDoc(url string, ro *grequests.RequestOptions) (*goquery.Document, error) {
	reader, err := this.client.GetReader(url, ro)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(reader)
}

func (this *Request) BuildJson(url string, ro *grequests.RequestOptions) (*simplejson.Json, error) {
	bjson, err := this.client.GetBytes(url, ro)
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(bjson)
}
