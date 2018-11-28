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

type ShareTask struct {
	Id            uint64 `json:"id"`
	Title         string `json:"title"`
	Summary       string `json:"summary"`
	Link          string `json:"link"`
	ShareLink     string `json:"share_link"`
	Image         string `json:"image,omitempty"`
	Points        string `json:"points,omitempty"`
	PointsLeft    string `json:"points_left,omitempty"`
	Bonus         string `json:"bonus,omitempty"`
	MaxViewers    uint   `json:"max_viewers,omitempty"`
	Viewers       uint   `json:"viewers,omitempty"`
	InsertedAt    string `json:"inserted_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
	OnlineStatus  int8   `json:"online_status,omitempty"`
	InIframe      bool   `json:"-"`
	ShowBonusHint bool   `json:"show_bonus_hint,omitempty"`
	Cid           []int  `json:"cid"`
}
