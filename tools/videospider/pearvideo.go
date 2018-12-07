package videospider

import (
//"github.com/mkideal/log"
)

type PearVideo struct {
	Request
}

func NewPearVideo(client *Client) *PearVideo {
	return &PearVideo{
		Request{
			client:   client,
			name:     "PearVideo",
			patterns: []string{`www\.pearvideo\.com/video_(\d+)`},
		},
	}
}

func (this *PearVideo) Get(link string) (info Video, err error) {
	doc, err := this.BuildDoc(link, nil)
	if err != nil {
		return info, err
	}
	html, _ := doc.Html()
	pubTime := doc.Find("div.date").Text()
	info.Link = SafeUrl(link)
	info.CreateTime = StringToMilliseconds("2006-01-02 15:04", pubTime)
	info.Title = doc.Find("h1.video-tt").Text()
	info.Desc = doc.Find(".details-content").Find("summary").Text()
	cover, _ := doc.Find("#poster").Find("img.img").Attr("src")
	info.Cover = SafeUrl(cover)
	var files []VideoLink

	hdUrl := R1(`hdUrl=\"(.+?)\"`, html)
	sdUrl := R1(`sdUrl=\"(.+?)\"`, html)
	ldUrl := R1(`ldUrl=\"(.+?)\"`, html)
	srcUrl := R1(`srcUrl=\"(.+?)\"`, html)
	if len(hdUrl) > 10 {
		hdUrl = SafeUrl(hdUrl)
		files = append(files, VideoLink{
			Resolution: HD,
			Link:       hdUrl,
		})
	}
	if len(sdUrl) > 10 {
		sdUrl = SafeUrl(sdUrl)
		files = append(files, VideoLink{
			Resolution: Normal,
			Link:       sdUrl,
		})
	}
	if len(ldUrl) > 10 {
		ldUrl = SafeUrl(ldUrl)
		files = append(files, VideoLink{
			Resolution: Low,
			Link:       ldUrl,
		})
	}
	if len(srcUrl) > 10 {
		srcUrl = SafeUrl(srcUrl)
		files = append(files, VideoLink{
			Resolution: Normal,
			Link:       srcUrl,
		})
	}
	info.Files = files
	return
}
