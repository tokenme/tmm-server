package afs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// CreateAfsCheck invokes the jaq.CreateAfsAppCheck API synchronously
// api document: https://help.aliyun.com/api/afs/createAfsCheck.html
func (client *Client) CreateAfsAppCheck(request *CreateAfsAppCheckRequest) (response *CreateAfsAppCheckResponse, err error) {
	response = CreateCreateAfsAppCheckResponse()
	err = client.DoAction(request, response)
	return
}

// CreateAfsCheckWithChan invokes the jaq.CreateAfsCheck API asynchronously
// api document: https://help.aliyun.com/api/cdn/createCdpOrder.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateAfsAppCheckWithChan(request *CreateAfsAppCheckRequest) (<-chan *CreateAfsAppCheckResponse, <-chan error) {
	responseChan := make(chan *CreateAfsAppCheckResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateAfsAppCheck(request)
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

// CreateAfsCheckWithCallback invokes the jaq.createAfsCheck API asynchronously
// api document: https://help.aliyun.com/api/afs/createAfsCheck.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateAfsAppCheckWithCallback(request *CreateAfsAppCheckRequest, callback func(response *CreateAfsAppCheckResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateAfsAppCheckResponse
		var err error
		defer close(result)
		response, err = client.CreateAfsAppCheck(request)
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

// CreateAfsAppCheckRequest is the request struct for api CreateAfsAppCheck
type CreateAfsAppCheckRequest struct {
	*requests.RpcRequest
	Session string `position:"Query" name:"Session"`
}

// CreateAfsAppCheckResponse is the response struct for api CreateAfsAppCheck
type CreateAfsAppCheckResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	Data      Data   `json:"Data" xml:"Data"`
}

type Data struct {
	SecondCheckResult int `json:"SecondCheckResult" xml:"SecondCheckResult"`
}

// CreateCreateAfsAppCheckRequest creates a request to invoke CreateAfsAppCheck API
func CreateCreateAfsAppCheckRequest() (request *CreateAfsAppCheckRequest) {
	request = &CreateAfsAppCheckRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("jaq", "2016-11-23", "AfsAppCheck", "", "")
	request.Domain = "jaq.aliyuncs.com"
	return
}

// CreateCreateAfsCheckResponse creates a response to parse from CreateAsfCheck response
func CreateCreateAfsAppCheckResponse() (response *CreateAfsAppCheckResponse) {
	response = &CreateAfsAppCheckResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
