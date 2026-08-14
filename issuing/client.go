package issuing

import (
	"github.com/uqpay/uqpay-sdk-go/v2/authdecision"
	"github.com/uqpay/uqpay-sdk-go/v2/common"
)

// Client provides access to Issuing APIs
type Client struct {
	Cards          *CardsClient
	Cardholders    *CardholdersClient
	Transactions   *TransactionsClient
	Products       *ProductsClient
	Balances       *BalancesClient
	Transfers      *TransfersClient
	Reports        *ReportsClient
	DownloadCenter *DownloadCenterClient
	MerchantBrands *MerchantBrandsClient
	AuthDecision   *authdecision.Client
}

// NewClient creates a new Issuing client
func NewClient(apiClient *common.APIClient) *Client {
	return &Client{
		Cards:          &CardsClient{client: apiClient},
		Cardholders:    &CardholdersClient{client: apiClient},
		Transactions:   &TransactionsClient{client: apiClient},
		Products:       &ProductsClient{client: apiClient},
		Balances:       &BalancesClient{client: apiClient},
		Transfers:      &TransfersClient{client: apiClient},
		Reports:        &ReportsClient{client: apiClient},
		DownloadCenter: &DownloadCenterClient{client: apiClient},
		MerchantBrands: &MerchantBrandsClient{client: apiClient},
		AuthDecision:   authdecision.NewClient(),
	}
}

func firstRequestOptions(opts []*common.RequestOptions) *common.RequestOptions {
	if len(opts) == 0 {
		return nil
	}
	return opts[0]
}
