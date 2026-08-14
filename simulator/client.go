package simulator

import "github.com/uqpay/uqpay-sdk-go/v2/common"

type Client struct {
	Issuing  *IssuingClient
	Deposits *DepositsClient
}

func NewClient(apiClient *common.APIClient) *Client {
	return &Client{
		Issuing:  &IssuingClient{client: apiClient},
		Deposits: &DepositsClient{client: apiClient},
	}
}

func firstRequestOptions(opts []*common.RequestOptions) *common.RequestOptions {
	if len(opts) == 0 {
		return nil
	}
	return opts[0]
}
