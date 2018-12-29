package info

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/forex"
	"github.com/tokenme/tmm/handler/admin"
	"net/http"
	"fmt"
)

type Ert struct {
	PointToCashRate decimal.Decimal `json:"point_to_cash_rate"`
	PointToUcRate   string          `json:"point_to_uc_rate"`
	UcToCashRate    decimal.Decimal `json:"uc_to_cash_rate"`
}

func ExchangeRateHandler(c *gin.Context) {
	rt := &Ert{}
	//积分提现汇率
	currency := "CNY"
	PTCPrice := common.GetPointPrice(Service, Config)
	rate := forex.Rate(Service, "USD", currency)
	rt.PointToCashRate = PTCPrice.Mul(rate)

	//积分兑换uc汇率
	//tmmPerTs, err := common.GetTMMPerTs(Config, Service)
	//if CheckErr(err, c) {
	//	return
	//}
	//pointsPerTs, err := common.GetPointsPerTs(Service)
	//if CheckErr(err, c) {
	//	return
	//}
	//rt.PointToUcRate = tmmPerTs.Div(pointsPerTs)

	//uc提现汇率
	UTCPrice := common.GetTMMPrice(Service, Config, common.RecyclePrice)
	rate = forex.Rate(Service, "USD", currency)
	rt.UcToCashRate = UTCPrice.Mul(rate)
	rates,_:=rt.PointToCashRate.Div(rt.UcToCashRate).Float64()
	rt.PointToUcRate = fmt.Sprintf("%.4f",rates)

	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    rt,
	})
}
