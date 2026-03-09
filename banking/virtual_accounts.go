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
	AccountBankID  string                    `json:"account_bank_id,omitempty"`
	AccountHolder  string                    `json:"account_holder,omitempty"`
	AccountNumber  string                    `json:"account_number"`
	Currency       string                    `json:"currency"`
	CountryCode    string                    `json:"country_code,omitempty"`
	BankName       string                    `json:"bank_name,omitempty"`
	BankAddress    string                    `json:"bank_address,omitempty"`
	Capability     *VirtualAccountCapability `json:"capability,omitempty"`
	ClearingSystem *VirtualAccountClearing   `json:"clearing_system,omitempty"`
	Status         string                    `json:"status"`
	CloseReason    string                    `json:"close_reason,omitempty"`
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
	PageSize   int    `json:"page_size"`          // required, 10-100
	PageNumber int    `json:"page_number"`        // required, >=1
	Currency   string `json:"currency,omitempty"` // optional, ISO 4217 comma-separated
}

// ListVirtualAccountsResponse represents a virtual account list response
type ListVirtualAccountsResponse struct {
	TotalPages int              `json:"total_pages"`
	TotalItems int              `json:"total_items"`
	Data       []VirtualAccount `json:"data"`
}

// CreateVirtualAccountRequest represents a virtual account creation request
type CreateVirtualAccountRequest struct {
	Currency      string `json:"currency"`                 // required, ISO 4217 currency code(s), e.g., "USD" or "USD,SGD" for multiple
	PaymentMethod string `json:"payment_method,omitempty"` // optional, "LOCAL" or "SWIFT"
}

// CreateVirtualAccountResponse represents the response from creating a virtual account
type CreateVirtualAccountResponse struct {
	Message string `json:"message"` // "SUCCESS"
}

// List lists virtual accounts
func (c *VirtualAccountsClient) List(ctx context.Context, req *ListVirtualAccountsRequest) (*ListVirtualAccountsResponse, error) {
	var resp ListVirtualAccountsResponse
	path := fmt.Sprintf("/v1/virtual/accounts?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

	if req.Currency != "" {
		path += fmt.Sprintf("&currency=%s", req.Currency)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list virtual accounts: %w", err)
	}
	return &resp, nil
}

// Create creates a new virtual account
// Note: The API returns {"message":"SUCCESS"} immediately. Actual account creation
// is confirmed asynchronously via webhooks (virtual.account.create / virtual.account.update).
func (c *VirtualAccountsClient) Create(ctx context.Context, req *CreateVirtualAccountRequest) (*CreateVirtualAccountResponse, error) {
	var resp CreateVirtualAccountResponse
	if err := c.client.Post(ctx, "/v1/virtual/accounts", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create virtual account: %w", err)
	}
	return &resp, nil
}
