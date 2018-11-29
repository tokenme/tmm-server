package task


const (
	API_OK = "OK"
)

type Response struct {
	code int `json:"code"`
	data interface{} `json:"data,omitempty"`
	message string `json:"message"`
}