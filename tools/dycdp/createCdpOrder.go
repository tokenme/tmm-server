package dycdp

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// CreateCdpOrder invokes the dycdp.CreateCdpOrder API synchronously
// api document: https://help.aliyun.com/api/dycdp/createCdpOrder.html
func (client *Client) CreateCdpOrder(request *CreateCdpOrderRequest) (response *CreateCdpOrderResponse, err error) {
	response = CreateCreateCdpOrderResponse()
	err = client.DoAction(request, response)
	return
}

// CreateCdpOrderWithChan invokes the dycdp.CreateCdpOrder API asynchronously
// api document: https://help.aliyun.com/api/cdn/createCdpOrder.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateCdpOrderWithChan(request *CreateCdpOrderRequest) (<-chan *CreateCdpOrderResponse, <-chan error) {
	responseChan := make(chan *CreateCdpOrderResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateCdpOrder(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// CreateCdpOrderWithCallback invokes the cdn.CreateCdpOrder API asynchronously
// api document: https://help.aliyun.com/api/dycdp/queryCdpOffer.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateCdpOrderWithCallback(request *CreateCdpOrderRequest, callback func(response *CreateCdpOrderResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateCdpOrderResponse
		var err error
		defer close(result)
		response, err = client.CreateCdpOrder(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// CreateCdpOrderRequest is the request struct for api CreateCdpOrder
type CreateCdpOrderRequest struct {
	*requests.RpcRequest
	PhoneNumber string           `position:"Query" name:"PhoneNumber"`
	OfferId     requests.Integer `position:"Query" name:"OfferId"`
	OutOrderId  string           `position:"Query" name:"OutOrderId"`
	Reason      string           `position:"Query" name:"Reason"`
	MaxPrice    requests.Integer `position:"Query" name:"MaxPrice"`
	ParamList   string           `position:"Query" name:"ParamList"`
}

// CreateCdpOrderResponse is the response struct for api CreateCdpOrder
type CreateCdpOrderResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	Data      Data   `json:"Data" xml:"Data"`
}

type Data struct {
	ResultCode  string `json:"ResultCode" xml:"ResultCode"`
	ResultMsg   string `json:"ResultMsg" xml:"ResultMsg"`
	OrderId     string `json:"OrderId" xml:"OrderId"`
	ExtendParam string `json:"ExtendParam" xml:"ExtendParam"`
}

// CreateCreateCdpOrderRequest creates a request to invoke CreateCdpOrder API
func CreateCreateCdpOrderRequest() (request *CreateCdpOrderRequest) {
	request = &CreateCdpOrderRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dycdpapi", "2018-06-10", "createCdpOrder", "", "")
	request.Domain = "dycdpapi.aliyuncs.com"
	return
}

// CreateCreateCdpOrderResponse creates a response to parse from CreateCdpOrder response
func CreateCreateCdpOrderResponse() (response *CreateCdpOrderResponse) {
	response = &CreateCdpOrderResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
