package banking

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// VirtualAccountsClient handles virtual account operations
type VirtualAccountsClient struct {
	client *common.APIClient
}

// VirtualAccount represents a virtual account
type VirtualAccount struct {
	VirtualAccountID   string               `json:"virtual_account_id"`
	VirtualAccountName string               `json:"virtual_account_name"`
	Status             string               `json:"status"`
	CreateTime         string               `json:"create_time"`
	CurrencyBankDetail []CurrencyBankDetail `json:"currency_bank_detail"`
}

// CurrencyBankDetail represents bank details for a specific currency
type CurrencyBankDetail struct {
	Currency        string `json:"currency"`
	BankName        string `json:"bank_name"`
	BankAddress     string `json:"bank_address"`
	BankCountryCode string `json:"bank_country_code"`
	AccountName     string `json:"account_name"`
	AccountNumber   string `json:"account_number"`
	SwiftCode       string `json:"swift_code"`
	RoutingNumber   string `json:"routing_number"`
	IBAN            string `json:"iban"`
	SortCode        string `json:"sort_code"`
	IFSCCode        string `json:"ifsc_code"`
}

// ListVirtualAccountsRequest represents a virtual account list request
type ListVirtualAccountsRequest struct {
	PageSize   int `json:"page_size"`   // required, 10-100
	PageNumber int `json:"page_number"` // required, >=1
}

// ListVirtualAccountsResponse represents a virtual account list response
type ListVirtualAccountsResponse struct {
	TotalPages int              `json:"total_pages"`
	TotalItems int              `json:"total_items"`
	Data       []VirtualAccount `json:"data"`
}

// CreateVirtualAccountRequest represents a virtual account creation request
type CreateVirtualAccountRequest struct {
	Currency      string `json:"currency"`       // required, ISO 4217 currency code(s), e.g., "USD" or "USD,SGD" for multiple
	PaymentMethod string `json:"payment_method"` // required, "LOCAL" or "SWIFT"
}

// List lists virtual accounts
func (c *VirtualAccountsClient) List(ctx context.Context, req *ListVirtualAccountsRequest) (*ListVirtualAccountsResponse, error) {
	var resp ListVirtualAccountsResponse
	path := fmt.Sprintf("/v1/virtual/accounts?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list virtual accounts: %w", err)
	}
	return &resp, nil
}

// Create creates a new virtual account
func (c *VirtualAccountsClient) Create(ctx context.Context, req *CreateVirtualAccountRequest) (*VirtualAccount, error) {
	var resp VirtualAccount
	if err := c.client.Post(ctx, "/v1/virtual/accounts", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create virtual account: %w", err)
	}
	return &resp, nil
}
