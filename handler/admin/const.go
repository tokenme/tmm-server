package admin

import "github.com/shopspring/decimal"

const (
	API_OK = "OK"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

type AddRequest struct {
	Title         string          `json:"title" form:"title" `
	Summary       string          `json:"summary" form:"summary" `
	Link          string          `json:"link" form:"link" `
	Image         string          `json:"image" form:"image"`
	FileExtension string          `json:"image_extension" from:"image_extension"`
	Points        decimal.Decimal `json:"points" form:"points" `
	Bonus         decimal.Decimal `json:"bonus" form:"bonus" `
	MaxViewers    uint            `json:"max_viewers" form:"max_viewers" `
	Cid           []int           `json:"cid" form:"cid"`
}
