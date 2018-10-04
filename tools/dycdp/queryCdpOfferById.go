package dycdp

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// QueryCdpOfferById invokes the dycdp.QueryCdpOfferById API synchronously
// api document: https://help.aliyun.com/api/dycdp/queryCdpOfferById.html
func (client *Client) QueryCdpOfferById(request *QueryCdpOfferByIdRequest) (response *QueryCdpOfferByIdResponse, err error) {
	response = CreateQueryCdpOfferByIdResponse()
	err = client.DoAction(request, response)
	return
}

// QueryCdpOfferByIdWithChan invokes the dycdp.QueryCdpOfferById API asynchronously
// api document: https://help.aliyun.com/api/cdn/queryCdpOffer.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryCdpOfferByIdWithChan(request *QueryCdpOfferByIdRequest) (<-chan *QueryCdpOfferByIdResponse, <-chan error) {
	responseChan := make(chan *QueryCdpOfferByIdResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryCdpOfferById(request)
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

// QueryCdpOfferByIdWithCallback invokes the cdn.QueryCdpOfferById API asynchronously
// api document: https://help.aliyun.com/api/dycdp/queryCdpOffer.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryCdpOfferByIdWithCallback(request *QueryCdpOfferByIdRequest, callback func(response *QueryCdpOfferByIdResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryCdpOfferByIdResponse
		var err error
		defer close(result)
		response, err = client.QueryCdpOfferById(request)
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

// QueryCdpOfferByIdRequest is the request struct for api QueryCdpOfferById
type QueryCdpOfferByIdRequest struct {
	*requests.RpcRequest
	OfferId string `position:"Query" name:"OfferId"`
}

// QueryCdpOfferByIdResponse is the response struct for api QueryCdpOfferById
type QueryCdpOfferByIdResponse struct {
	*responses.BaseResponse
	RequestId  string     `json:"RequestId" xml:"RequestId"`
	Code       string     `json:"Code" xml:"Code"`
	Message    string     `json:"Message" xml:"Message"`
	FlowOffers FlowOffers `json:"FlowOffers" xml:"FlowOffers"`
}

// CreateQueryCdpOfferByIdRequest creates a request to invoke QueryCdpOfferById API
func CreateQueryCdpOfferByIdRequest() (request *QueryCdpOfferByIdRequest) {
	request = &QueryCdpOfferByIdRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dycdpapi", "2018-06-10", "queryCdpOfferById", "", "")
	request.Domain = "dycdpapi.aliyuncs.com"
	return
}

// CreateQueryCdpOfferByIdResponse creates a response to parse from QueryCdpOfferById response
func CreateQueryCdpOfferByIdResponse() (response *QueryCdpOfferByIdResponse) {
	response = &QueryCdpOfferByIdResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
