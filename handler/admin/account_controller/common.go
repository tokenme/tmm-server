package account_controller

const (
	ExchangeUc       = "兑换UC"
	ExchangePoint    = "兑换积分"
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
}

type Exchange struct {
	Type   string `json:"type"`
	Pay    string `json:"pay"`
	Get    string `json:"get"`
	When   string `json:"when"`
	Status string `json:"status"`
}

type Task struct {
	Point  string `json:"point"`
	Type   string `json:"type"`
	When   string `json:"when"`
	Status string `json:"status"`
}

type DrawMoney struct {
	Exchange
}
