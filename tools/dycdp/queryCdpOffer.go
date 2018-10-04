package dycdp

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/shopspring/decimal"
)

// QueryCdpOffer invokes the dycdp.QueryCdpOffer API synchronously
// api document: https://help.aliyun.com/api/dycdp/queryCdpOffer.html
func (client *Client) QueryCdpOffer(request *QueryCdpOfferRequest) (response *QueryCdpOfferResponse, err error) {
	response = CreateQueryCdpOfferResponse()
	err = client.DoAction(request, response)
	return
}

// QueryCdpOfferWithChan invokes the dycdp.QueryCdpOffer API asynchronously
// api document: https://help.aliyun.com/api/cdn/queryCdpOffer.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryCdpOfferWithChan(request *QueryCdpOfferRequest) (<-chan *QueryCdpOfferResponse, <-chan error) {
	responseChan := make(chan *QueryCdpOfferResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryCdpOffer(request)
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

// QueryCdpOfferWithCallback invokes the cdn.QueryCdpOffer API asynchronously
// api document: https://help.aliyun.com/api/dycdp/queryCdpOffer.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryCdpOfferWithCallback(request *QueryCdpOfferRequest, callback func(response *QueryCdpOfferResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryCdpOfferResponse
		var err error
		defer close(result)
		response, err = client.QueryCdpOffer(request)
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

// QueryCdpOfferRequest is the request struct for api QueryCdpOffer
type QueryCdpOfferRequest struct {
	*requests.RpcRequest
	Vendor      string `position:"Query" name:"Vendor"`
	ChannelType string `position:"Query" name:"ChannelType"`
	Province    string `position:"Query" name:"Province"`
}

// QueryCdpOfferResponse is the response struct for api QueryCdpOffer
type QueryCdpOfferResponse struct {
	*responses.BaseResponse
	RequestId  string     `json:"RequestId" xml:"RequestId"`
	Code       string     `json:"Code" xml:"Code"`
	Message    string     `json:"Message" xml:"Message"`
	FlowOffers FlowOffers `json:"FlowOffers" xml:"FlowOffers"`
}

type FlowOffers struct {
	FlowOffer []FlowOffer `json:"FlowOffer" xml:"FlowOffer"`
}

type FlowOffer struct {
	Vendor      string          `json:"Vendor" xml:"Vendor"`
	ChannelType string          `json:"ChannelType" xml:"ChannelType"`
	Province    string          `json:"Province" xml:"Province"`
	Grade       string          `json:"Grade" xml:"Grade"`
	FlowType    string          `json:"FlowType" xml:"FlowType"`
	UseEff      string          `json:"UseEff" xml:"UseEff"`
	UseLimit    string          `json:"UseLimit" xml:"UseLimit"`
	UseScene    string          `json:"UseScene" xml:"UseScene"`
	Right       string          `json:"Right" xml:"Right"`
	OfferId     uint64          `json:"OfferId" xml:"OfferId"`
	Price       uint64          `json:"Price" xml:"Price"`
	Discount    decimal.Decimal `json:"Discount" xml:"Discount"`
}

// CreateQueryCdpOfferRequest creates a request to invoke QueryCdpOffer API
func CreateQueryCdpOfferRequest() (request *QueryCdpOfferRequest) {
	request = &QueryCdpOfferRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dycdpapi", "2018-06-10", "queryCdpOffer", "", "")
	request.Domain = "dycdpapi.aliyuncs.com"
	return
}

// CreateQueryCdpOfferResponse creates a response to parse from QueryCdpOffer response
func CreateQueryCdpOfferResponse() (response *QueryCdpOfferResponse) {
	response = &QueryCdpOfferResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
