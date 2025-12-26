package paymentintent

import "github.com/uqpay/uqpay-sdk-go/common"

// Client represents the Payment API client
type Client struct {
	Intents *PaymentIntentsClient
	// Future clients will be added here:
	// Attempts  *PaymentAttemptsClient
	// Refunds   *PaymentRefundsClient
	// Reports   *PaymentReportsClient
	// Balances  *PaymentBalancesClient
	// Payouts   *PaymentPayoutsClient
}

// NewClient creates a new Payment API client
func NewClient(apiClient *common.APIClient) *Client {
	return &Client{
		Intents: &PaymentIntentsClient{client: apiClient},
	}
}
