package payment

import "github.com/uqpay/uqpay-sdk-go/v2/common"

// Client represents the Payment API client
type Client struct {
	PaymentIntents  *PaymentIntentsClient
	PaymentAttempts *PaymentAttemptsClient
	Refunds         *PaymentRefundsClient
	Reports         *PaymentReportsClient
	Balances        *PaymentBalancesClient
	Payouts         *PaymentPayoutsClient
	BankAccounts    *BankAccountsClient
	Terminals       *TerminalsClient
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
		BankAccounts:    &BankAccountsClient{client: apiClient},
		Terminals:       &TerminalsClient{client: apiClient},
	}
}

func requestOptionsWithClientID(clientID string, opts ...*common.RequestOptions) *common.RequestOptions {
	resolved := common.RequestOptions{
		ClientID: clientID,
	}
	if len(opts) == 0 || opts[0] == nil {
		return &resolved
	}

	resolved = *opts[0]
	if resolved.ClientID == "" {
		resolved.ClientID = clientID
	}
	return &resolved
}
