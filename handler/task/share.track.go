package task

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/utils"
	"net/http"
	"strings"
	"time"
)

func ShareTrackHandler(c *gin.Context) {
	targetUrl := c.Query("url")
	taskId, err := utils.DecryptUint64(c.Param("encryptedTaskId"), []byte(Config.LinkSalt))
	if err != nil {
		log.Error(err.Error())
		c.Redirect(http.StatusFound, targetUrl)
		return
	}
	uid, _ := utils.DecryptUint64(c.Query("uid"), []byte(Config.LinkSalt))
	if targetUrl != "" {
		db := Service.Db
		trackSource := common.TrackFromUnknown
		isWx := strings.Contains(strings.ToLower(c.Request.UserAgent()), "micromessenger")
		isUCoin := strings.Contains(strings.ToLower(c.Request.UserAgent()), "ucoin")
		var uv string
		if isWx {
			trackSource = common.TrackFromWechat
		} else if isUCoin || uid > 0 {
			trackSource = common.TrackFromUCoin
			if checkUV(taskId, uid) {
				uv = ", uv=uv+1"
			}
		}
		_, _, err := db.Query(`INSERT INTO tmm.share_task_stats (task_id, record_on, source) VALUES (%d, NOW(), %d) ON DUPLICATE KEY UPDATE pv=pv+1%s`, taskId, trackSource, uv)
		if err != nil {
			log.Error(err.Error())
		}
	}
	c.Redirect(http.StatusFound, targetUrl)
}

func checkUV(taskId uint64, uid uint64) bool {
	redisConn := Service.Redis.Master.Get()
	defer redisConn.Close()
	task := common.ShareTask{Id: taskId}
	uidKey := task.UidKey(uid)
	_, err := redisConn.Do("GET", uidKey)
	var found bool
	if err == nil {
		found = true
	} else {
		_, err = redisConn.Do("SETEX", uidKey, 60*60*24, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Error(err.Error())
		}
	}
	return found
}
