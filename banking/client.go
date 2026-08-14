package banking

import "github.com/uqpay/uqpay-sdk-go/v2/common"

// Client represents the Banking API client
type Client struct {
	Transfers                  *TransfersClient
	Balances                   *BalancesClient
	VirtualAccounts            *VirtualAccountsClient
	VirtualAccountApplications *VirtualAccountApplicationsClient
	Deposits                   *DepositsClient
	Beneficiaries              *BeneficiariesClient
	Payouts                    *PayoutsClient
	Conversions                *ConversionClient
	ExchangeRates              *ExchangeRatesClient
}

// NewClient creates a new Banking API client
func NewClient(apiClient *common.APIClient) *Client {
	return &Client{
		Transfers:                  &TransfersClient{client: apiClient},
		Balances:                   &BalancesClient{client: apiClient},
		VirtualAccounts:            &VirtualAccountsClient{client: apiClient},
		VirtualAccountApplications: &VirtualAccountApplicationsClient{client: apiClient},
		Deposits:                   &DepositsClient{client: apiClient},
		Beneficiaries:              &BeneficiariesClient{client: apiClient},
		Payouts:                    &PayoutsClient{client: apiClient},
		Conversions:                &ConversionClient{client: apiClient},
		ExchangeRates:              &ExchangeRatesClient{client: apiClient},
	}
}

func firstRequestOptions(opts []*common.RequestOptions) *common.RequestOptions {
	if len(opts) == 0 {
		return nil
	}
	return opts[0]
}
