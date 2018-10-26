package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/tokenme/tmm/tools/xinge"
	"github.com/tokenme/tmm/tools/xinge/push"
	"github.com/tokenme/tmm/tools/xinge/query"
	"github.com/tokenme/tmm/tools/xinge/tag"
	"log"
	"time"
)

func main() {
	client := xinge.NewClient(2200311250, "491e482934603d9d8124469ea74ea72b")
	getAppDeviceNum(client)
}

func getMsgStatus(client *xinge.Client) {
	req := query.GetMsgStatusRequest{
		BaseRequest: client.DefaultBaseRequest(),
		PushIds: []query.PushIdMap{
			{PushId: "3206825"},
		},
	}
	var ret = &xinge.BaseResponse{}
	err := client.Run(req, ret)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(ret)
}

func singleDevice(client *xinge.Client) {
	msg := xinge.IosMessage{
		Aps: xinge.ApsAttr{
			Alert: xinge.ApsAlert{
				Body: "test alert",
			},
		},
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalln(err)
	}
	sendTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	req := push.SingleDeviceRequest{
		BaseRequest: client.DefaultBaseRequest(),
		DeviceToken: "c43e4def735b9787948abfe526d3d11e2a10dcd332c9cf8c403b3a564ed5c073",
		Message:     msg.String(),
		SendTime:    sendTime,
		Environment: 1,
	}
	var ret = &xinge.BaseResponse{}
	err = client.Run(req, ret)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(ret)
}

func tagDevice(client *xinge.Client) {
	msg := xinge.IosMessage{
		Aps: xinge.ApsAttr{
			Alert: xinge.ApsAlert{
				Body: "test alert",
			},
		},
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalln(err)
	}
	sendTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	req := push.TagDeviceRequest{
		BaseRequest: client.DefaultBaseRequest(),
		TagList:     []string{"COUNTRY:86"},
		TagsOp:      xinge.OR,
		Message:     msg.String(),
		SendTime:    sendTime,
		Environment: 1,
	}
	var ret = &xinge.PushIdResponse{}
	err = client.Run(req, ret)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(ret)
}

func createMultiPush(client *xinge.Client) {
	msg := xinge.IosMessage{
		Aps: xinge.ApsAttr{
			Alert: xinge.ApsAlert{
				Body: "test alert",
			},
		},
	}
	req := push.CreateMultiPushRequest{
		BaseRequest: client.DefaultBaseRequest(),
		Message:     msg.String(),
		Environment: 1,
	}
	var ret = &xinge.PushIdResponse{}
	err := client.Run(req, ret)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(ret)
}

func batchSet(client *xinge.Client) {
	req := tag.BatchSetRequest{
		BaseRequest: client.DefaultBaseRequest(),
		TagTokenList: [][]string{
			{"Admin", "c43e4def735b9787948abfe526d3d11e2a10dcd332c9cf8c403b3a564ed5c073"},
		},
	}
	var ret = &xinge.BaseResponse{}
	err := client.Run(req, ret)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(ret)
}

func batchDel(client *xinge.Client) {
	req := tag.BatchDelRequest{
		BaseRequest: client.DefaultBaseRequest(),
		TagTokenList: [][]string{
			{"Admin", "c43e4def735b9787948abfe526d3d11e2a10dcd332c9cf8c403b3a564ed5c073"},
		},
	}
	var ret = &xinge.BaseResponse{}
	err := client.Run(req, ret)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(ret)
}

func getAppDeviceNum(client *xinge.Client) {
	req := query.GetAppDeviceNumRequest{
		BaseRequest: client.DefaultBaseRequest(),
	}
	var ret = &xinge.BaseResponse{}
	err := client.Run(req, ret)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(ret)
}

func getDeviceToken(client *xinge.Client) {
	req := query.GetAppTokenInfoRequest{
		BaseRequest: client.DefaultBaseRequest(),
		DeviceToken: "c43e4def735b9787948abfe526d3d11e2a10dcd332c9cf8c403b3a564ed5c073",
	}
	var ret = &xinge.DeviceTokenResponse{}
	err := client.Run(req, ret)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(ret)
}

func queryAppTags(client *xinge.Client) {
	req := query.QueryAppTagsRequest{
		BaseRequest: client.DefaultBaseRequest(),
	}
	var ret = &xinge.TagsResponse{}
	err := client.Run(req, ret)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(ret)
}

func queryTokenTags(client *xinge.Client) {
	req := query.QueryTokenTagsRequest{
		BaseRequest: client.DefaultBaseRequest(),
		DeviceToken: "c43e4def735b9787948abfe526d3d11e2a10dcd332c9cf8c403b3a564ed5c073",
	}
	var ret = &xinge.TagsResponse{}
	err := client.Run(req, ret)
	if err != nil {
		log.Fatalln(err)
	}
	spew.Dump(ret)
}
