package videospider

import (
	"fmt"
	"github.com/levigross/grequests"
	//"github.com/mkideal/log"
)

type MiaoPai struct {
	Request
}

func NewMiaoPai(client *Client) *MiaoPai {
	return &MiaoPai{
		Request{
			client:   client,
			name:     "MiaoPai",
			patterns: []string{`n.miaopai.com\/media\/(.+)`, `static.xiaokaxiu.com\/xkx\/h5share\/index.html?scid=(.+)`, `v.xiaokaxiu.com\/v\/(.+).html`},
		},
	}
}

func (this *MiaoPai) Get(link string) (info Video, err error) {
	vid := R1(this.patterns[0], link)
	keyUrl := fmt.Sprintf("https://n.miaopai.com/api/aj_media/info.json?smid=%s&appid=530", vid)
	ro := &grequests.RequestOptions{
		UserAgent:    "Mozilla/5.0 (Linux; U; Android 4.3; en-us; SM-N900T Build/JSS15J) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
		UseCookieJar: false,
	}
	vInfo, err := this.BuildJson(keyUrl, ro)
	if err != nil {
		return info, err
	}
	info.Link = SafeUrl(link)
	info.Title, _ = vInfo.Get("data").Get("description").String()
	info.CreateTime, _ = vInfo.Get("data").Get("created_at").Int64()
	metadata := vInfo.Get("data").Get("meta_data").GetIndex(0)
	cover, _ := metadata.Get("pics").Get("l").String()
	info.Cover = SafeUrl(cover)
	low, _ := metadata.Get("play_urls").Get("l").String()
	mid, _ := metadata.Get("play_urls").Get("m").String()
	high, _ := metadata.Get("play_urls").Get("n").String()
	if low != "" {
		low = SafeUrl(low)
		info.Files = append(info.Files, VideoLink{
			Resolution: Low,
			Link:       low,
		})
	}
	if mid != "" {
		mid = SafeUrl(mid)
		info.Files = append(info.Files, VideoLink{
			Resolution: Normal,
			Link:       mid,
		})
	}
	if high != "" {
		high = SafeUrl(high)
		info.Files = append(info.Files, VideoLink{
			Resolution: HD,
			Link:       high,
		})
	}
	width, _ := metadata.Get("upload").Get("width").Uint64()
	height, _ := metadata.Get("upload").Get("height").Uint64()
	duration, _ := metadata.Get("upload").Get("length").Int64()
	info.Duration = duration * 1000
	info.Size = VideoSize{
		Width:  uint(width),
		Height: uint(height),
	}
	return
}
