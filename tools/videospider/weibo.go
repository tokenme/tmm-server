package videospider

import (
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
			patterns: []string{`weibo\.(cn|com)\/tv\/v\/(\w+)`, `weibo\.(cn|com)\/\d+\/(\w+)`, `weibo\.(cn|com)\/status\/([\w|\d]+)`, `weibo\.(cn|com)\/detail\/([\w|\d]+)`, `weibo\.(cn|com)/\d+/\w+`},
		},
	}
}

func (this *Weibo) Get(link string) (info Video, err error) {
	html, err := this.GetContent(link, nil)
	if err != nil {
		return info, err
	}
	info.Link = SafeUrl(link)
	info.Title = R1(`"title":\ "(.+)"`, html)
	if info.Title == "" {
		info.Title = R1(`"page_title":\ "(.+)"`, html)
	}
	info.Desc = R1(`"content2":\ "(.+)"`, html)
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
