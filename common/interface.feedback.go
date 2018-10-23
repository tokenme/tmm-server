package common

type Feedback struct {
	Ts      string     `json:"ts"`
	Msg     string     `json:"msg"`
	Image   string     `json:"image,omitempty"`
	Channel string     `json:"-"`
	Bot     bool       `json:"bot,omitempty"`
	Replies []Feedback `json:"replies,omitempty"`
}
