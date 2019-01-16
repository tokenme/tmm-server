package task

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/wechatmp"
	tokenUtils "github.com/tokenme/tmm/utils/token"
	"github.com/ua-parser/uap-go/uaparser"
	"net/http"
	"strings"
	"time"
)

const (
	WX_AUTH_GATEWAY                            = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect"
	WX_AUTH_URL                                = "https://jkgj-isv.isvjcloud.com/rest/m/u/weauth"
	MaxUserRateLimitSecondCounter        int64 = 3
	MaxUserRateLimitSecondBlockCounter   int64 = 5
	MaxUserRateLimitSecondDuration             = 1
	MaxUserRateLimitSecondBlockDurateion       = 600
	MaxUserRateLimitMinuteCounter        int64 = 300
	MaxUserRateLimitMinuteDuration             = 600
	MaxUserRateLimitBlockDurateion             = 3600
	MaxUserRateLimitMinuteBlockCounter   int64 = 500
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
	var trackId string
	taskId, deviceId, err := common.DecryptShareTaskLink(c.Param("encryptedTaskId"), c.Param("encryptedDeviceId"), Config)
	if CheckErr(err, c) {
		return
	}
	if CheckWithCode(taskId == 0 || deviceId == "", NOTFOUND_ERROR, "not found", c) {
		return
	}
	db := Service.Db
	creator, err := getTaskCreator(Service, taskId)
	if CheckErr(err, c) {
		return
	}
	code := c.DefaultQuery("code", "null")
	var (
		isBlocked     bool
		isRateLimited bool
	)
	if creator == 0 {
		isBlocked, isRateLimited, err = checkBlockUser(Service, deviceId, code != "null")
		if err != nil {
			log.Error(err.Error())
		}
	}

	mpClient := wechatmp.NewClient(Config.Wechat.AppId, Config.Wechat.AppSecret, Service, c)
	isWx := strings.Contains(strings.ToLower(c.Request.UserAgent()), "micromessenger")

	if isWx && !isBlocked && !isRateLimited {
		if len(code) > 0 && code != "null" {
			trackId, err = getTrackId(Service, Config, mpClient, code)
			if err != nil {
				log.Error(err.Error())
			}
			if trackId == "" {
				isBlocked = true
			}
		} else if code == "null" {
			shareUri := c.Request.URL.String()
			redirectUrl := fmt.Sprintf("%s%s", Config.ShareBaseUrl, shareUri)
			mpClient.AuthRedirect(redirectUrl, "snsapi_base", "tmm-share")
			return
		}
	}

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

	if isBlocked || isRateLimited {
		c.Redirect(http.StatusFound, task.Link)
		return
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

	currentUrl := fmt.Sprintf("%s%s", Config.ShareBaseUrl, c.Request.URL.String())
	jsConfig, err := mpClient.GetJSConfig(currentUrl)
	if err != nil {
		log.Error(err.Error())
	}
	trackSource := common.TrackFromUnknown
	if isWx {
		trackSource = common.TrackFromWechat
	}
	_, _, err = db.Query(`INSERT INTO tmm.share_task_stats (task_id, record_on, source) VALUES (%d, NOW(), %d) ON DUPLICATE KEY UPDATE pv=pv+1`, task.Id, trackSource)
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

func getTaskCreator(service *common.Service, taskId uint64) (uint64, error) {
	db := service.Db
	rows, _, err := db.Query(`SELECT creator FROM tmm.share_tasks WHERE id=%d LIMIT 1`, taskId)
	if err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, errors.New("not found")
	}
	return rows[0].Uint64(0), nil
}

func checkBlockUser(service *common.Service, deviceId string, haveCode bool) (isBlocked bool, isRateLimited bool, err error) {
	db := service.Db
	var (
		userId uint64
		task   common.ShareTask
	)
	{
		query := `
        SELECT d.user_id, IFNULL(us.blocked, 0) AS blocked, IFNULL(us.block_whitelist, 0) AS block_whitelist
        FROM tmm.devices AS d
        LEFT JOIN tmm.user_settings AS us ON d.user_id = us.user_id
        WHERE d.id = "%s" LIMIT 1
    `
		rows, _, err := db.Query(query, deviceId)
		if err != nil {
			return false, false, err
		}
		if len(rows) > 0 {
			row := rows[0]
			userId = row.Uint64(0)
			isBlocked = row.Int(1) == 1 && row.Int(2) == 0
		}
	}
	log.Info("Share device id:%s, user id:%d, is block:%v", deviceId, userId, isBlocked)
	if isBlocked {
		return true, false, nil
	}
	if userId == 0 {
		return true, false, errors.New("user not found")
	}
	redisConn := service.Redis.Master.Get()
	defer redisConn.Close()

	blockKey := task.UserRateLimitBlockKey(userId)
	secBlockKey := task.UserRateLimitSecondBlockKey(userId)
	blockKeyExists, err := redis.Int(redisConn.Do("EXISTS", blockKey, secBlockKey))
	if err != nil {
		log.Error(err.Error())
	} else if blockKeyExists > 0 {
		log.Warn("RateLimit blocked: %d, duration: %d", userId, MaxUserRateLimitBlockDurateion)
		isRateLimited = true
	}

	if !haveCode {
		secondKey := task.UserRateLimitSecondKey(userId)
		secondCounter, err := redis.Int64(redisConn.Do("INCR", secondKey))
		if err != nil {
			log.Error(err.Error())
		}
		if secondCounter <= 1 {
			_, err := redisConn.Do("EXPIRE", secondKey, MaxUserRateLimitSecondDuration)
			if err != nil {
				log.Error(err.Error())
			}
		} else if secondCounter >= MaxUserRateLimitSecondBlockCounter {
			_, _, err := db.Query(`INSERT INTO tmm.user_settings(user_id, blocked, block_whitelist) VALUES (%d, 1, 0) ON DUPLICATE KEY UPDATE blocked=VALUES(blocked), block_whitelist=VALUES(block_whitelist)`, userId)
			if err != nil {
				log.Error(err.Error())
			} else {
				log.Warn("Block RateLimit user: %d, sec counter: %d", userId, secondCounter)
			}
			isBlocked = true
			isRateLimited = true
		} else if secondCounter >= MaxUserRateLimitSecondCounter {
			log.Warn("RateLimit for user:%d, second counter:%d", userId, secondCounter)
			_, err := redisConn.Do("SET", secBlockKey, "1", "EX", MaxUserRateLimitSecondBlockDurateion)
			if err != nil {
				log.Error(err.Error())
			}
			if _, _, err := db.Query(`INSERT INTO tmm.share_blocked_logs
			(user_id,inserted_at,second_count,minute_count) VALUES(%d,'%s',%d,%d) 
			ON DUPLICATE KEY UPDATE second_count=second_count+VALUES(second_count), minute_count=minute_count+VALUES(minute_count)`,
				userId, time.Now().Format(`2006-01-02`), 1, 0); err != nil {
				return isBlocked, isRateLimited, err
			}
			isRateLimited = true
		}

		if isBlocked {
			return true, isRateLimited, nil
		}

		minuteKey := task.UserRateLimitMinuteKey(userId)
		minuteCounter, err := redis.Int64(redisConn.Do("INCR", minuteKey))
		if err != nil {
			log.Error(err.Error())
		}
		if minuteCounter <= 1 {
			_, err := redisConn.Do("EXPIRE", minuteKey, MaxUserRateLimitMinuteDuration)
			if err != nil {
				log.Error(err.Error())
			}
		} else if minuteCounter >= MaxUserRateLimitMinuteBlockCounter {
			_, _, err := db.Query(`INSERT INTO tmm.user_settings(user_id, blocked, block_whitelist) VALUES (%d, 1, 0) ON DUPLICATE KEY UPDATE blocked=VALUES(blocked), block_whitelist=VALUES(block_whitelist)`, userId)
			if err != nil {
				log.Error(err.Error())
			} else {
				log.Warn("Block RateLimit user: %d, min counter: %d", userId, minuteCounter)
			}
			isBlocked = true
			isRateLimited = true
		} else if minuteCounter >= MaxUserRateLimitMinuteCounter {
			log.Warn("RateLimit for user:%d, minute counter:%d", userId, minuteCounter)
			_, err := redisConn.Do("SET", blockKey, "1", "EX", MaxUserRateLimitBlockDurateion)
			if err != nil {
				log.Error(err.Error())
			}
			if _, _, err := db.Query(`INSERT INTO tmm.share_blocked_logs
			(user_id,inserted_at,second_count,minute_count) VALUES(%d,'%s',%d,%d) 
			ON DUPLICATE KEY UPDATE second_count=second_count+VALUES(second_count), minute_count=minute_count+VALUES(minute_count)`,
				userId, time.Now().Format(`2006-01-02`), 0, 1); err != nil {
				return isBlocked, isRateLimited, err
			}
			isRateLimited = true
		}
	}
	return isBlocked, isRateLimited, nil
}

func getTrackId(service *common.Service, config common.Config, mpClient *wechatmp.Client, code string) (string, error) {
	redisConn := service.Redis.Master.Get()
	defer redisConn.Close()
	codeKey := common.WxCodeKey(code)
	openId, _ := redis.String(redisConn.Do("GET", codeKey))
	if openId != "" {
		return "", errors.New(fmt.Sprintf("Wechat code used, openId:%s", openId))
	}
	oauthAccessToken, err := mpClient.GetOAuthAccessToken(code)
	if err != nil {
		return "", err
	}
	if oauthAccessToken.Openid == "" {
		return "", nil
	}
	_, err = redisConn.Do("SETEX", codeKey, 60*2, oauthAccessToken.Openid)
	if err != nil {
		log.Error(err.Error())
	}
	now := time.Now()
	cryptOpenid := common.CryptOpenid{
		Openid: oauthAccessToken.Openid,
		Ts:     now.Unix(),
	}
	return cryptOpenid.Encode([]byte(config.YktApiSecret))
}
