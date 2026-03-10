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
	AccountBankID  string                    `json:"account_bank_id,omitempty"` // Unique bank identifier
	AccountHolder  string                    `json:"account_holder,omitempty"`  // Holder of the global account
	AccountNumber  string                    `json:"account_number"`            // Account or IBAN number
	Currency       string                    `json:"currency"`                  // ISO 4217 currency code, e.g. "usd"
	CountryCode    string                    `json:"country_code,omitempty"`    // ISO 3166-1 alpha-2 country code
	BankName       string                    `json:"bank_name,omitempty"`       // Name of the account bank
	BankAddress    string                    `json:"bank_address,omitempty"`    // Address of the account bank
	Capability     *VirtualAccountCapability `json:"capability,omitempty"`      // Payment method configuration
	ClearingSystem *VirtualAccountClearing   `json:"clearing_system,omitempty"` // Optional; clearing system details
	Status         string                    `json:"status"`                    // ACTIVE, INACTIVE, or CLOSED
	CloseReason    string                    `json:"close_reason,omitempty"`    // Optional; reason for account closure
}

// VirtualAccountCapability represents the payment capability
type VirtualAccountCapability struct {
	PaymentMethod string `json:"payment_method"` // SWIFT or LOCAL
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
	TotalPages int              `json:"total_pages"` // Total number of pages available
	TotalItems int              `json:"total_items"` // Total count of available items
	Data       []VirtualAccount `json:"data"`        // List of virtual account objects
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
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of
func (c *VirtualAccountsClient) List(ctx context.Context, req *ListVirtualAccountsRequest, opts ...*common.RequestOptions) (*ListVirtualAccountsResponse, error) {
	var resp ListVirtualAccountsResponse
	path := fmt.Sprintf("/v1/virtual/accounts?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

	if req.Currency != "" {
		path += fmt.Sprintf("&currency=%s", req.Currency)
	}

	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to list virtual accounts: %w", err)
	}
	return &resp, nil
}

// Create creates a new virtual account
// Note: The API returns {"message":"SUCCESS"} immediately. Actual account creation
// is confirmed asynchronously via webhooks (virtual.account.create / virtual.account.update).
// Optional RequestOptions can be provided to set custom headers like x-idempotency-key or x-on-behalf-of
func (c *VirtualAccountsClient) Create(ctx context.Context, req *CreateVirtualAccountRequest, opts ...*common.RequestOptions) (*CreateVirtualAccountResponse, error) {
	var resp CreateVirtualAccountResponse
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.PostWithOptions(ctx, "/v1/virtual/accounts", req, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to create virtual account: %w", err)
	}
	return &resp, nil
}
