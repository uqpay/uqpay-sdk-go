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
	AccountBankID  string                    `json:"account_bank_id"`
	AccountHolder  string                    `json:"account_holder"`
	AccountNumber  string                    `json:"account_number"`
	Currency       string                    `json:"currency"`
	CountryCode    string                    `json:"country_code"`
	BankName       string                    `json:"bank_name"`
	BankAddress    string                    `json:"bank_address"`
	Capability     *VirtualAccountCapability `json:"capability,omitempty"`
	ClearingSystem *VirtualAccountClearing   `json:"clearing_system,omitempty"`
	Status         string                    `json:"status"`
	CloseReason    string                    `json:"close_reason"`
}

// VirtualAccountCapability represents the payment capability
type VirtualAccountCapability struct {
	PaymentMethod string `json:"payment_method"`
}

// VirtualAccountClearing represents clearing system details
type VirtualAccountClearing struct {
	Type  string `json:"type"`  // e.g., "bic_swift"
	Value string `json:"value"` // e.g., "SGBDBHB2XXX"
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
	Country       string `json:"country"`                  // required, ISO 3166-1 alpha-2
	Currency      string `json:"currency"`                 // required, one ISO 4217 currency
	PaymentMethod string `json:"payment_method,omitempty"` // optional: "LOCAL" or "SWIFT"
	Nickname      string `json:"nickname,omitempty"`       // optional, max 255 characters
}

// List lists virtual accounts
func (c *VirtualAccountsClient) List(ctx context.Context, req *ListVirtualAccountsRequest, opts ...*common.RequestOptions) (*ListVirtualAccountsResponse, error) {
	var resp ListVirtualAccountsResponse
	path := fmt.Sprintf("/v1/virtual/accounts?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)
	if err := c.client.GetWithOptions(ctx, path, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to list virtual accounts: %w", err)
	}
	return &resp, nil
}

// Create creates a new virtual account
func (c *VirtualAccountsClient) Create(ctx context.Context, req *CreateVirtualAccountRequest, opts ...*common.RequestOptions) (*VirtualAccountApplicationResponse, error) {
	var resp VirtualAccountApplicationResponse
	if err := c.client.PostWithOptions(ctx, "/v1/virtual/accounts", req, &resp, firstRequestOptions(opts)); err != nil {
		return nil, fmt.Errorf("failed to create virtual account: %w", err)
	}
	return &resp, nil
}
