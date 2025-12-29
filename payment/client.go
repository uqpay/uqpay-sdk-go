package payment

import "github.com/uqpay/uqpay-sdk-go/common"

// Client represents the Payment API client
type Client struct {
	PaymentIntents  *PaymentIntentsClient
	PaymentAttempts *PaymentAttemptsClient
	Refunds         *PaymentRefundsClient
	Reports         *PaymentReportsClient
	Balances        *PaymentBalancesClient
	Payouts         *PaymentPayoutsClient
}

// NewClient creates a new Payment API client
func NewClient(apiClient *common.APIClient) *Client {
	return &Client{
		PaymentIntents:  &PaymentIntentsClient{client: apiClient},
		PaymentAttempts: &PaymentAttemptsClient{client: apiClient},
		Refunds:         &PaymentRefundsClient{client: apiClient},
		Reports:         &PaymentReportsClient{client: apiClient},
		Balances:        &PaymentBalancesClient{client: apiClient},
		Payouts:         &PaymentPayoutsClient{client: apiClient},
	}
}
