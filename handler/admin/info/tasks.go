package info

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	. "github.com/tokenme/tmm/handler"
	"encoding/json"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)

const TasksKey = `info-tasks`

func TasksInfoHandler(c *gin.Context) {
	db := Service.Db
	redisConn := Service.Redis.Master.Get()
	context, err := redisConn.Do(`GET`, TasksKey)
	if CheckErr(err, c) {
		return
	}
	if context != nil {
		info := TaskInfo{}
		if json.Unmarshal(context.([]byte), &info) == nil {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    info,
			})
			return
		}
	}
	query := `SELECT SUM(points_left),COUNT(*) FROM tmm.share_tasks`
	rows, _, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	row := rows[0]
	point, err := decimal.NewFromString(row.Str(0))
	if CheckErr(err, c) {
		return
	}

	info := TaskInfo{}
	info.TotalPoint = point
	info.TotalTask = row.Int(1)

	query = `
	SELECT SUM(task.points) AS total_point,
	COUNT(*) AS total_share_task,
	COUNT(distinct devices.user_id) AS total_user
	FROM tmm.device_share_tasks  AS task
	INNER JOIN tmm.devices AS devices ON devices.id = task.device_id  `
	rows, res, err := db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	row = rows[0]
	shareTotalPoint, err := decimal.NewFromString(row.Str(res.Map(`total_point`)))
	if CheckErr(err, c) {
		return
	}
	totalShareTask := row.Int(res.Map(`total_share_task`))
	totalUser := row.Float(res.Map(`total_user`))

	info.ShareTask.TotalTask = totalShareTask
	info.ShareTask.Point = shareTotalPoint
	info.AvgUserShare = float64(totalShareTask) / totalUser

	query = `SELECT COUNT(*) FROM tmm.reading_logs`
	rows, _, err = db.Query(query)
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `not found`, c) {
		return
	}
	totalreadTask := rows[0].Float(0)
	info.ReadTask.TotalTask = int(totalreadTask)
	info.AvgUserRead = totalreadTask / totalUser
	data, err := json.Marshal(&info)
	if CheckErr(err, c) {
		return
	}
	//签到率没有数据等腾讯云数据统计的数据
	redisConn.Do("SET", TasksKey, data, `EX`, KeyAlive)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    info,
	})
}
