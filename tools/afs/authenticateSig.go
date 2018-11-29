package afs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// CreateAfsCheck invokes the jaq.CreateAfsAppCheck API synchronously
// api document: https://help.aliyun.com/api/afs/createAfsCheck.html
func (client *Client) AuthenticateSig(request *AuthenticateSigRequest) (response *AuthenticateSigResponse, err error) {
	response = CreateAuthenticateSigResponse()
	err = client.DoAction(request, response)
	return
}

// CreateAfsCheckWithChan invokes the jaq.CreateAfsCheck API asynchronously
// api document: https://help.aliyun.com/api/cdn/createCdpOrder.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AuthenticateSigWithChan(request *AuthenticateSigRequest) (<-chan *AuthenticateSigResponse, <-chan error) {
	responseChan := make(chan *AuthenticateSigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AuthenticateSig(request)
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
func (client *Client) AuthenticateSigWithCallback(request *AuthenticateSigRequest, callback func(response *AuthenticateSigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AuthenticateSigResponse
		var err error
		defer close(result)
		response, err = client.AuthenticateSig(request)
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

// CreateAfsAppCheckRequest is the request struct for api AuthenticateSig
type AuthenticateSigRequest struct {
	*requests.RpcRequest
	SessionId string `position:"Query" name:"SessionId"`
	Token     string `position:"Query" name:"Token"`
	Sig       string `position:"Query" name:"Sig"`
	Scene     string `position:"Query" name:"Scene"`
	AppKey    string `position:"Query" name:"AppKey"`
	RemoteIp  string `position:"Query" name:"RemoteIp"`
}

// CreateAfsAppCheckResponse is the response struct for api CreateAfsAppCheck
type AuthenticateSigResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Code      int    `json:"Code" xml:"Code"`
	Msg       string `json:"Msg" xml:"Msg"`
	RiskLevel string `json:"RiskLevel" xml:"RiskLevel"`
	Detail    string `json:"Detail" xml:"Detail"`
}

// CreateCreateAfsAppCheckRequest creates a request to invoke CreateAfsAppCheck API
func CreateAuthenticateSigRequest() (request *AuthenticateSigRequest) {
	request = &AuthenticateSigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("afs", "2018-01-12", "AuthenticateSig", "", "")
	request.Domain = "afs.aliyuncs.com"
	return
}

// CreateAuthenticateSigResponse creates a response to parse from CreateAsfCheck response
func CreateAuthenticateSigResponse() (response *AuthenticateSigResponse) {
	response = &AuthenticateSigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
