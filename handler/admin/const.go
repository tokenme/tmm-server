package admin



const (
	API_OK = "OK"
)

type Response struct {
	Code int `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Message string `json:"message"`
}