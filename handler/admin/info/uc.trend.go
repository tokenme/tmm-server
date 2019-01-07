package info

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/etherscan-api"
	"time"
	"github.com/shopspring/decimal"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/handler/admin"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

const (
	UcTransferAmount = iota
	UcTransferUsers
)
const (
	UcTransferAmountKey = "uc-amount-%s-%s"
	UcTransferUsersKey  = "uc-users-%s-%s"
)

func UcTrendHandler(c *gin.Context) {
	client := etherscan.New(etherscan.Mainnet, Config.EtherscanAPIKey)
	conn := Service.Redis.Master.Get()
	defer conn.Close()
	var req StatsRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	startTime := time.Now().AddDate(0, 0, -7).Format(`2006-01-02`)
	endTime := time.Now().Format(`2006-01-02`)
	geth := Service.Geth
	newBlock, _ := geth.BlockByNumber(c, nil)
	if req.StartTime != "" {
		startTime = req.StartTime
	}
	if req.EndTime != "" {
		endTime = req.EndTime
	}
	var key string
	if req.Type == UcTransferAmount {
		key = fmt.Sprintf(UcTransferAmountKey, startTime, endTime)
	} else {
		key = fmt.Sprintf(UcTransferUsersKey, startTime, endTime)
	}
	bytes, err := redis.Bytes(conn.Do(`GET`, key))
	if err == nil && bytes != nil {
		var data TrendData
		if !CheckErr(json.Unmarshal(bytes, &data), c) {
			c.JSON(http.StatusOK, admin.Response{
				Code:    0,
				Message: admin.API_OK,
				Data:    data,
			})
		}
		return
	}
	tm, _ := time.Parse(`2006-01-02`, startTime)
	end, _ := time.Parse(`2006-01-02`, endTime)
	blockTime := time.Unix(newBlock.Time().Int64(), 0)
	startBlock := int(newBlock.Number().Int64()) - int(blockTime.Sub(tm).Seconds()/12)
	address := Config.TMMTokenAddress
	txs, err := client.ERC20Transfers(&address, nil, &startBlock, nil, 0, 100000, false)
	if CheckErr(err, c) {
		return
	}
	usersCount := make(map[string]map[string]struct{})
	dataMap := make(map[string]decimal.Decimal)
	for _, tx := range txs {
		txtime := tx.TimeStamp.Time()
		if txtime.After(tm) && txtime.Before(end.AddDate(0, 0, 1)) && tx.From != "0x12c9b5a2084decd1f73af37885fc1e0ced5d5ee8" && tx.From != "0x251e3c2de440c185912ea701a421d80bbe5ee8c9" {
			format := txtime.Format(`2006-01-02`)
			if req.Type == UcTransferUsers {
				if value, ok := usersCount[format]; ok {
					value[tx.From] = struct{}{}
				} else {
					_map := make(map[string]struct{})
					_map[tx.From] = struct{}{}
					usersCount[format] = _map
				}
			} else {
				if value, ok := dataMap[format]; ok {
					dataMap[format] = value.Add(decimal.NewFromBigInt(tx.Value.Int(), -1*9))
				} else {
					dataMap[format] = decimal.NewFromBigInt(tx.Value.Int(), -1*9)
				}
			}
		}
		if txtime.After(end.AddDate(0, 0, 1)) {
			break
		}
	}
	var indexName, valueList []string
	var yaxisName, title string
	if req.Type == UcTransferUsers {
		yaxisName = "人数"
		title = "Uc转账人数"
		for {
			if tm.Equal(end) {
				if accountMap, ok := usersCount[tm.Format(`2006-01-02`)]; ok {
					indexName = append(indexName, tm.Format(`2006-01-02`))
					valueList = append(valueList, fmt.Sprintf("%d", len(accountMap)))
				} else {
					indexName = append(indexName, tm.Format(`2006-01-02`))
					valueList = append(valueList, fmt.Sprintf("%d", 0))
				}
				break
			}
			if accountMap, ok := usersCount[tm.Format(`2006-01-02`)]; ok {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf("%d", len(accountMap)))
				tm = tm.AddDate(0, 0, 1)
			} else {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf("%d", 0))
				tm = tm.AddDate(0, 0, 1)
			}
		}
	} else {
		yaxisName = "Uc"
		title = "Uc转账金额"
		for {
			if tm.Equal(end) {
				if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
					indexName = append(indexName, tm.Format(`2006-01-02`))
					valueList = append(valueList, value.Ceil().String())
				} else {
					indexName = append(indexName, tm.Format(`2006-01-02`))
					valueList = append(valueList, fmt.Sprintf("%d", 0))
				}
				break
			}
			if value, ok := dataMap[tm.Format(`2006-01-02`)]; ok {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, value.Ceil().String())
				tm = tm.AddDate(0, 0, 1)
			} else {
				indexName = append(indexName, tm.Format(`2006-01-02`))
				valueList = append(valueList, fmt.Sprintf("%d", 0))
				tm = tm.AddDate(0, 0, 1)
			}
		}
	}
	var data TrendData
	data.Title = title
	data.Xaxis.Data = indexName
	data.Xaxis.Name = "日期"
	data.Yaxis.Name = yaxisName
	data.Series = append(data.Series, Series{Data: valueList, Name: data.Title, Type: "line"})
	bytes, err = json.Marshal(&data)
	if CheckErr(err, c) {
		return
	}
	_,err=conn.Do(`SET`, key, bytes, `EX`, 60*60)
	fmt.Println(err)
	c.JSON(http.StatusOK, admin.Response{
		Code:    0,
		Message: admin.API_OK,
		Data:    data,
	})

}
