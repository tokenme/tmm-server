package good

import (
	//"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	//"github.com/mkideal/log"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/tools/ykt"
	"net/http"
)

const DEFAULT_PAGE_SIZE uint = 30

type ListRequest struct {
	Page     uint `json:"page" form:"page"`
	PageSize uint `json:"page_size" form:"page_size"`
}

func ListHandler(c *gin.Context) {
	var req ListRequest
	if CheckErr(c.Bind(&req), c) {
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 || req.PageSize > DEFAULT_PAGE_SIZE {
		req.PageSize = DEFAULT_PAGE_SIZE
	}
	yktReq := ykt.GoodListRequest{
		Source:   1,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	res, err := yktReq.Run()
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, res.Data.Items)
}
