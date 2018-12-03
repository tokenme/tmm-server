package videospider

import (
	"errors"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	//"github.com/mkideal/log"
	"net/url"
	"regexp"
)

type Krcom struct {
	Request
}

func NewKrcom(client *Client) *Krcom {
	return &Krcom{
		Request{
			client:   client,
			name:     "Krcom",
			patterns: []string{`krcom.cn\/(\d+)\/episodes\/(\d+):(\d+)`},
		},
	}
}

func (this *Krcom) Get(link string) (info Video, err error) {
	re, _ := regexp.Compile(this.patterns[0])
	submatch := re.FindStringSubmatch(link)
	if len(submatch) != 4 {
		return info, errors.New("wrong url")
	}
	channelId := submatch[1]
	videoId := fmt.Sprintf("%s:%s", submatch[2], submatch[3])
	keyUrl := fmt.Sprintf("https://krcom.cn/h5/videodata?channel_id=%s&video_id=%s", channelId, videoId)
	vInfo, err := this.BuildJson(keyUrl, nil)
	if err != nil {
		return info, err
	}
	data, _ := vInfo.Get("data").String()
	dataStr, err := url.QueryUnescape(data)
	if err != nil {
		return info, err
	}
	vInfo, err = simplejson.NewJson([]byte(dataStr))
	if err != nil {
		return info, err
	}
	info.Link = link
	info.Title, _ = vInfo.Get("video_info").Get("video_title").String()
	cover, _ := vInfo.Get("video_info").Get("cover_img").String()
	info.Cover, _ = url.QueryUnescape(cover)
	low, _ := vInfo.Get("video_info").Get("urls").Get("mp4_ld_mp4").String()
	mid, _ := vInfo.Get("video_info").Get("urls").Get("mp4_720p_mp4").String()
	high, _ := vInfo.Get("video_info").Get("urls").Get("mp4_hd_mp4").String()
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
	return
}
