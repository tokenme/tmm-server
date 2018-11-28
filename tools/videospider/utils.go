package videospider

import (
	"regexp"
	"strings"
	"time"
)

func RxOf(pattern string, content string, index int) (rcontent string) {
	re, _ := regexp.Compile(pattern)
	submatch := re.FindStringSubmatch(content)
	for i, v := range submatch {
		if i == index {
			rcontent = v
			break
		}
	}
	return
}

func R1(pattern string, content string) (rcontent string) {
	return RxOf(pattern, content, 1)
}

func R1Of(patterns []string, s string) (rcontent string) {
	for _, pattern := range patterns {
		if rcontent = R1(pattern, s); len(rcontent) > 0 {
			break
		}
	}
	return
}

func StringToMilliseconds(fmtstring string, datestring string) int64 {
	tm2, _ := time.Parse(fmtstring, datestring)
	return tm2.Unix()
}

func SafeUrl(link string) string {
	if strings.HasPrefix(link, "http://") {
		return strings.Replace(link, "http://", "https://", -1)
	}
	return link
}
