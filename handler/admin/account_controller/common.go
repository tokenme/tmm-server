package account_controller

const (
	ExchangeUc       = "UC兑换"
	ExchangePoint    = "积分兑换"
	DrawMoneyByPoint = "提现_By_积分"
	DrawMoneyByUc    = "提现_By_UC"
)

const (
	TaskSuccessful = "已结算"
)
const (
	Failed = iota
	Success
	Pending
)
const (
	Reading = iota
	Invite
	Share
	BfBouns
	AppTask
)

var MsgMap = map[int]string{
	Success: "成功",
	Failed:  "失败",
	Pending: "等待打包中",
}

var typeMap = map[int]string{
	Reading: "阅读",
	Invite:  "拉新好友",
	Share:   "分享任务",
	BfBouns: "好友贡献",
	AppTask: "安装app",
}

type Task struct {
	Type   string `json:"type"`
	Pay    string `json:"pay"`
	Get    string `json:"get"`
	When   string `json:"when"`
	Status string `json:"status"`
}

type PageOptions struct {
	Id    int `form:"id"`
	Page  int `form:"page"`
	Limit int `form:"limit"`
	Types int `form:"type"`
}
