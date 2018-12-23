package account_controller


import (
	"github.com/gin-gonic/gin"
	."github.com/tokenme/tmm/handler"
	"strconv"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
)
func DrawMoneyByUcHandler(c *gin.Context){
	db := Service.Db
	id, err := strconv.Atoi(c.DefaultQuery(`id`, `-1`))
	if CheckErr(err, c) {
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery(`page`, `-1`))
	if CheckErr(err, c) {
		return
	}
	var offset int
	limit := 10
	if page > 0 {
		offset = limit * (page - 1)
	}
	query := `
SELECT 
	tmm,
	cny,
	inserted_at,
	tx_status
FROM 
	tmm.withdraw_txs 
WHERE 
	user_id = %d
ORDER BY inserted_at DESC
LIMIT %d OFFSET %d`
	rows, _, err := db.Query(query, id, limit, offset)
	if CheckErr(err, c) {
		return
	}

	var DrawMoneyList []*DrawMoney

	for _, row := range rows {
		drawMoney := &DrawMoney{}
		drawMoney.Type = DrawMoneyByUc
		drawMoney.Pay = fmt.Sprintf("%.2f",row.Float(0))
		drawMoney.Get = fmt.Sprintf("%.2f", row.Float(1))
		drawMoney.When = row.Str(2)
		drawMoney.Status = MsgMap[row.Int(3)]
		DrawMoneyList = append(DrawMoneyList, drawMoney)
	}

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    DrawMoneyList,
	})
}
