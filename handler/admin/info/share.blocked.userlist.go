package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"strings"
	"regexp"
	"strconv"
	"github.com/tokenme/tmm/common"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

func BlockedUserListHandler(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery(`limit`, `10`))
	if CheckErr(err, c) {
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery(`page`, `1`))
	if CheckErr(err, c) {
		return
	}

	var offset int
	if limit < 1 {
		limit = 10
	}
	if page > 0 {
		offset = (page - 1) * limit
	}

	conn := Service.Redis.Master.Get()
	defer conn.Close()
	bytes, _ := conn.Do(`scan`, offset, `match`, `st-*-rate-block`, `COUNT`, limit)
	bytesString := fmt.Sprintf("%s", bytes)
	regexpUserId, _ := regexp.Compile(`-(\d*)-`)
	UserIdList := regexpUserId.FindAllStringSubmatch(bytesString, -1)

	var idList []string
	for _, UserId := range UserIdList {
		idList = append(idList, UserId[1])
	}

	db := Service.Db
	query := `
SELECT
	u.id,
	u.mobile,
	IFNULL(wx.nick,u.nickname) 
FROM 
	ucoin.users AS u
LEFT JOIN tmm.wx AS wx ON wx.user_id = u.id 
WHERE u.id IN (%s)

`
	var userList []*common.User
	var total int
	if len(idList) > 0 {
		rows, _, err := db.Query(query, strings.Join(idList, `,`))
		if CheckErr(err, c) {
			return
		}

		for _, row := range rows {
			userList = append(userList, &common.User{
				Id:     row.Uint64(0),
				Mobile: row.Str(1),
				Nick:   row.Str(2),
			})
		}

		bytes, _ = conn.Do(`scan`, 0, `match`, `st-*-rate-block`, `COUNT`, 10000)
		bytesString = fmt.Sprintf("%s", bytes)
		UserIdList = regexpUserId.FindAllStringSubmatch(bytesString, -1)
		total = len(UserIdList)
	}

	c.JSON(http.StatusOK, admin.Response{
		Message: admin.API_OK,
		Code:    0,
		Data: gin.H{
			`total`: total,
			`date`:  userList,
		},
	})

}
