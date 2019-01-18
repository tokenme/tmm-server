package account_controller

const (
	ExchangeUc       = "UC兑换"
	ExchangePoint    = "积分兑换"
	DrawMoneyByPoint = "积分提现"
	DrawMoneyByUc    = "UC提现"
	TaskSuccessful   = "已结算"
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
	Failed:  "失败",
	Success: "成功",
	Pending: "等待打包中",
}

var typeMap = map[int]string{
	Reading: "阅读奖励",
	Invite:  "拉新好友",
	Share:   "分享文章",
	BfBouns: "好友贡献",
	AppTask: "安装app",
}

var ForMatMap = map[int]string{
	Reading: "已阅读%d秒",
	Share:   "分享%d次数",
}

const (
	Waiting   = 0
	Succeeded = 1
	Refused   = -1
)

var WithDrawMap = map[int]string{
	Waiting:   "等待审核",
	Succeeded: "审核成功",
	Refused:   "审核拒绝",
}

var StatusMap = map[int]string{
	Waiting:   "等待",
	Succeeded: "成功",
	Refused:   "拒绝",
}

const (
	ShareBonus      = 1
	AppBonus        = 2
	ActiveUserBonus = 3
)

const (
	DirectFirend = iota
	InDirectFirend
)

var InviteMap = map[int]string{
	ShareBonus:      "好友贡献(文章)",
	AppBonus:        "好友贡献(app)",
	ActiveUserBonus: "下线好友活跃贡献",
}
var FirendMap = map[int]string{
	DirectFirend:   "直接好友",
	InDirectFirend: "间接好友",
}

type Task struct {
	Type   string `json:"type"`
	Pay    string `json:"pay"`
	Get    string `json:"get"`
	When   string `json:"when"`
	Status string `json:"status"`
	Extra  string `json:"extra"`
}

type PageOptions struct {
	Id      int    `form:"id"`
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
	Types   int    `form:"type"`
	Devices string `form:"devices"`
}
