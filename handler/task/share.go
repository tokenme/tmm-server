package task

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	tokenUtils "github.com/tokenme/tmm/utils/token"
    "github.com/tokenme/tmm/tools/wechatmp"
	"github.com/ua-parser/uap-go/uaparser"
	"net/http"
	"net/url"
	"strings"
)

const (
	WX_AUTH_GATEWAY = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect"
	WX_AUTH_URL     = "https://jkgj-isv.isvjcloud.com/rest/m/u/weauth"
)

type ShareData struct {
    AppId      string
    JSConfig   wechatmp.JSConfig
	Task       common.ShareTask
	TrackId    string
	IsIOS      bool
	InviteLink string
	ImpLink    string
    ShareLink  string
	ValidTrack bool
}

func ShareHandler(c *gin.Context) {
	trackId := c.DefaultQuery("track_id", "null")
	taskId, deviceId, err := common.DecryptShareTaskLink(c.Param("encryptedTaskId"), c.Param("encryptedDeviceId"), Config)
	if CheckErr(err, c) {
		return
	}
	if Check(taskId == 0 || deviceId == "", "not found", c) {
		return
	}
	isWx := strings.Contains(strings.ToLower(c.Request.UserAgent()), "micromessenger")
	if isWx && trackId == "null" {
		wxAuthUrl := url.QueryEscape(WX_AUTH_URL)
		shareUri := c.Request.URL.String()
		symbol := "?"
		if strings.Contains(shareUri, "?") {
			symbol = "&"
		}
		wxRedirectUrl := url.QueryEscape(fmt.Sprintf("%s%s%strack_id=___OPENID___", Config.ShareBaseUrl, shareUri, symbol))
		redirectUrl := fmt.Sprintf(WX_AUTH_GATEWAY, Config.Wechat.AppId, wxAuthUrl, wxRedirectUrl)
		c.Redirect(http.StatusFound, redirectUrl)
		return
	}

	db := Service.Db
	query := `SELECT
    st.id,
    st.title,
    st.summary,
    st.link,
    st.image,
    st.max_viewers,
    st.bonus,
    st.points,
    st.points_left,
    dst.viewers,
    ic.id
FROM tmm.share_tasks AS st
LEFT JOIN tmm.device_share_tasks AS dst ON (dst.task_id=st.id AND dst.device_id='%s')
LEFT JOIN tmm.devices AS d ON (d.id=dst.device_id)
LEFT JOIN tmm.invite_codes AS ic ON (ic.user_id=d.user_id)
WHERE st.id=%d
LIMIT 1`
	rows, _, err := db.Query(query, db.Escape(deviceId), taskId)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, "not found", c) {
		log.Error("Not found")
		return
	}
	row := rows[0]
	bonus, _ := decimal.NewFromString(row.Str(6))
	points, _ := decimal.NewFromString(row.Str(7))
	pointsLeft, _ := decimal.NewFromString(row.Str(8))
	task := common.ShareTask{
		Id:         row.Uint64(0),
		Title:      row.Str(1),
		Summary:    row.Str(2),
		Link:       row.Str(3),
		Image:      row.Str(4),
		MaxViewers: row.Uint(5),
		Bonus:      bonus,
		Points:     points,
		PointsLeft: pointsLeft,
	}
	if strings.HasPrefix(task.Link, "https://tmm.tokenmama.io/article/show") {
		task.Link = strings.Replace(task.Link, "https://tmm.tokenmama.io/article/show", "https://static.tianxi100.com/article/show", -1)
		task.TimelineOnly = true
	}
	validTrack := true
	wxFrom := c.Query("from")
	if task.TimelineOnly && (!isWx || wxFrom != "timeline" && wxFrom != "singlemessage" && wxFrom != "groupmessage") {
		validTrack = false
	}
	task.InIframe = task.ShouldUseIframe()
	inviteCode := tokenUtils.Token(row.Uint64(10))
	impLink, _ := task.GetShareImpLink(deviceId, Config)

	parser, err := uaparser.New(Config.UAParserPath)
	var isIOS bool
	if err != nil {
		log.Error(err.Error())
	} else {
		client := parser.Parse(c.Request.UserAgent())
		if strings.Contains(strings.ToLower(client.Os.Family), "ios") {
			isIOS = true
		}
	}

    mpClient := wechatmp.NewClient(Config.Wechat.AppId, Config.Wechat.AppSecret, Service)
    currentUrl := fmt.Sprintf("%s%s", Config.ShareBaseUrl, c.Request.URL.String())
    jsConfig, err := mpClient.GetJSConfig(currentUrl)
    if err != nil {
        log.Error(err.Error())
    }

	c.HTML(http.StatusOK, "share.tmpl", ShareData{
        AppId:      Config.Wechat.AppId,
        JSConfig:   jsConfig,
		Task:       task,
		TrackId:    trackId,
		IsIOS:      isIOS,
		ValidTrack: validTrack,
		InviteLink: fmt.Sprintf("https://tmm.tokenmama.io/invite/%s", inviteCode.Encode()),
        ShareLink:  currentUrl,
		ImpLink:    impLink})
}
