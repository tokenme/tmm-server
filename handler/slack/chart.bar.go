package slack

import (
	//"github.com/davecgh/go-spew/spew"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/wcharczuk/go-chart"
	"net/http"
)

type BarChartValue struct {
	Value float64
	Label string
}

func ChartBarHandler(c *gin.Context) {
	data := c.Param("data")
	buf, err := base64.URLEncoding.DecodeString(data)
	if CheckErr(err, c) {
		return
	}
	var bars []BarChartValue
	err = json.Unmarshal(buf, &bars)
	if CheckErr(err, c) {
		return
	}
	var values []chart.Value
	for _, bar := range bars {
		values = append(values, chart.Value{
			Value: bar.Value,
			Label: bar.Label,
		})
	}
	sbc := chart.BarChart{
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		XAxis:    chart.StyleShow(),
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
		},
		Bars: values,
	}
	var b bytes.Buffer
	err = sbc.Render(chart.PNG, &b)
	if CheckErr(err, c) {
		return
	}
	c.Data(http.StatusOK, "image/png", b.Bytes())
}
