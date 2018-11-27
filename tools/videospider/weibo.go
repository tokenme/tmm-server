package videospider

import (
	"github.com/levigross/grequests"
	//"github.com/mkideal/log"
	"strconv"
)

type Weibo struct {
	Request
}

func NewWeibo(client *Client) *Weibo {
	return &Weibo{
		Request{
			client:   client,
			name:     "Weibo",
			patterns: []string{`weibo\.com\/tv\/v\/(\w+)`, `weibo.com\/\d+\/(\w+)`},
		},
	}
}

func (this *Weibo) Get(link string) (info Video, err error) {
	ro := &grequests.RequestOptions{
		UserAgent:    "Mozilla/5.0 (Linux; U; Android 4.3; en-us; SM-N900T Build/JSS15J) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
		UseCookieJar: false,
	}
	html, err := this.GetContent(link, ro)
	if err != nil {
		return info, err
	}
	info.Link = SafeUrl(link)
	info.Title = R1(`"title":\ "(.+)"`, html)
	info.Cover = R1(`"url":\ "(.+)"`, html)
	durationStr := R1(`"duration":\ "(\d+)"`, html)
	duration, _ := strconv.ParseInt(durationStr, 10, 64)
	info.Duration = duration * 1000

	sizeStr := R1(`"size":\ "(\d+)"`, html)
	size, _ := strconv.ParseUint(sizeStr, 10, 64)
	widthStr := R1(`"width":\ "(\d+)"`, html)
	width, _ := strconv.ParseUint(widthStr, 10, 64)
	heightStr := R1(`"height":\ "(\d+)"`, html)
	height, _ := strconv.ParseUint(heightStr, 10, 64)
	info.Size = VideoSize{
		Width:  uint(width),
		Height: uint(height),
		Size:   size,
	}
	streamUrl := R1(`"stream_url":\ "(.+)"`, html)
	streamUrlHd := R1(`"stream_url_hd":\ "(.+)"`, html)
	if streamUrl != "" {
		info.Files = append(info.Files, VideoLink{
			Resolution: Normal,
			Link:       SafeUrl(streamUrl),
		})
	}

	if streamUrlHd != "" {
		info.Files = append(info.Files, VideoLink{
			Resolution: HD,
			Link:       SafeUrl(streamUrlHd),
		})
	}

	return
}
