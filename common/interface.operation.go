package common

import (
)

type Activity struct {
	Id              uint64      `json:"id"`
    RowId           uint        `json:"row_id"`
	OnlineStatus    int8        `json:"online_status"`
    Title           string      `json:"title"`
    Image           string      `json:"image"`
    ShareImage      string      `json:"share_image,omitempty"`
    Link            string      `json:"link,omitempty"`
	Width           uint        `json:"width,omitempty"`
	Height          uint        `json:"height,omitempty"`
    Action          string      `json:"action,omitempty"`
}

type Headline struct {
	Id              uint64      `json:"id"`
	OnlineStatus    int8        `json:"online_status"`
    Title           string      `json:"title"`
    Link            string      `json:"link,omitempty"`
}
