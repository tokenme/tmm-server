package common

type Adgroup struct {
	Id           uint64      `json:"id"`
	OnlineStatus uint8       `json:"online_status"`
	Adzone       *Adzone     `json:"-"`
	Creatives    []*Creative `json:"creatives"`
}

type Creative struct {
	Id        uint64 `json:"id"`
	AdgroupId uint64 `json:"adgroup_id,omitempty"`
	Image     string `json:"image,omitempty"`
	Link      string `json:"link,omitempty"`
	Width     uint   `json:"width,omitempty"`
	Height    uint   `json:"height,omitempty"`
}

type Adzone struct {
	Id   uint64 `json:"id"`
	Cid  uint   `json:"cid"`
	Page uint   `json:"page"`
	Idx  int    `json:"idx"`
}
