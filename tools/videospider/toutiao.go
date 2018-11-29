package videospider

import (
	"encoding/base64"
	"fmt"
	//"github.com/mkideal/log"
)

type Toutiao struct {
	Request
}

func NewToutiao(client *Client) *Toutiao {
	return &Toutiao{
		Request{
			client:   client,
			name:     "Toutiao",
			patterns: []string{`toutiao\.com\/item\/(\d+)`, `toutiao\.com\/[ia](\d+)`, `365yg.com\/([\w|\d]+)`, `toutiaoimg.cn/group/\d+/\?iid=(\d+)`},
		},
	}
}

func (this *Toutiao) Get(link string) (info Video, err error) {
	doc, err := this.BuildDoc(link)
	if err != nil {
		return info, err
	}
	html, err := doc.Html()
	if err != nil {
		return info, err
	}
	info.Link = SafeUrl(link)
	info.Title = R1(`title: \'(.+)\'`, html)
	vid := R1(`video(?:id|Id)\:.?\'(\w+)\'`, html)
	createTimeStr := R1(`time: \'(\d+\/\d+\/\d+)\'`, html)
	info.CreateTime = StringToMilliseconds("2006/01/02", createTimeStr)
	keyUrl := fmt.Sprintf("http://i.snssdk.com/video/urls/1/toutiao/mp4/%s", vid)
	vInfo, err := this.BuildJson(keyUrl)
	if err != nil {
		return info, err
	}
	base64Url, _ := vInfo.Get("data").Get("video_list").Get("video_1").Get("main_url").String()
	cover, _ := vInfo.Get("data").Get("poster_url").String()
	info.Cover = SafeUrl(cover)
	duration, _ := vInfo.Get("data").Get("video_duration").Float64()
	info.Duration = int64(duration * 1000)
	ts, err := base64.StdEncoding.DecodeString(base64Url)
	if err != nil {
		return info, err
	}
	width, _ := vInfo.Get("data").Get("video_list").Get("video_1").Get("vwidth").Uint64()
	height, _ := vInfo.Get("data").Get("video_list").Get("video_1").Get("vheight").Uint64()
	size, _ := vInfo.Get("data").Get("video_list").Get("video_1").Get("size").Uint64()
	info.Size = VideoSize{
		Width:  uint(width),
		Height: uint(height),
		Size:   size,
	}
	info.Files = []VideoLink{
		{
			Resolution: Normal,
			Link:       SafeUrl(string(ts)),
		},
	}
	return
}
