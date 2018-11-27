package task

type ShareAddRequest struct {
	Title      string `json:"title" form:"title" binding:"required"`
	Summary    string `json:"summary" form:"summary" binding:"required"`
	Link       string `json:"link" form:"link" binding:"required"`
	Image      string `json:"image" form:"image"`
	Points     string `json:"points" form:"points" binding:"required"`
	Bonus      string `json:"bonus" form:"bonus" binding:"required"`
	MaxViewers uint   `json:"max_viewers" form:"max_viewers" binding:"required"`
	Cid        []int  `json:"cid"  form:"cid"`
}

type OnlineStatusRequest struct {
	TaskId int `json:"id" form:"id"`
	Status int `json:"status" form:"status"`
}
