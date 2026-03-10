package payment

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// BankAccountsClient handles bank account operations for settlement purposes
type BankAccountsClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// CreateBankAccountRequest represents a bank account creation request
type CreateBankAccountRequest struct {
	AccountNumber   string `json:"account_number"`             // Required. Bank account number (e.g., IBAN)
	BankName        string `json:"bank_name"`                  // Required. Name of the bank
	SwiftCode       string `json:"swift_code"`                 // Required. SWIFT/BIC code, 8-11 characters
	BankCountryCode string `json:"bank_country_code"`          // Required. ISO 3166-1 alpha-2 country code
	BankAddress     string `json:"bank_address"`               // Required. Physical address of the bank
	Currency        string `json:"currency"`                   // Required. ISO 4217 currency code (e.g., "USD", "GBP")
	BankCodeType    string `json:"bank_code_type,omitempty"`   // Conditional. Allowed: aba, bank_code, sort_code, bsb_code, ifsc, cnaps_number
	BankCodeValue   string `json:"bank_code_value,omitempty"`  // Conditional. Routing code value; format varies by bank_code_type
	BankBranchCode  string `json:"bank_branch_code,omitempty"` // Conditional. Required when currency = "CAD"
}

// UpdateBankAccountRequest represents a bank account update request
type UpdateBankAccountRequest struct {
	AccountNumber   string `json:"account_number"`             // Required. Bank account number (e.g., IBAN)
	BankName        string `json:"bank_name"`                  // Required. Name of the bank
	SwiftCode       string `json:"swift_code"`                 // Required. SWIFT/BIC code, 8-11 characters
	BankCountryCode string `json:"bank_country_code"`          // Required. ISO 3166-1 alpha-2 country code
	BankAddress     string `json:"bank_address"`               // Required. Physical address of the bank
	BankCodeType    string `json:"bank_code_type,omitempty"`   // Conditional. Allowed: aba, bank_code, sort_code, bsb_code, ifsc, cnaps_number
	BankCodeValue   string `json:"bank_code_value,omitempty"`  // Conditional. Routing code value; format varies by bank_code_type
	BankBranchCode  string `json:"bank_branch_code,omitempty"` // Conditional. Required when currency = "CAD"
}

// ListBankAccountsRequest represents a bank accounts list request
type ListBankAccountsRequest struct {
	PageNumber int `json:"page_number"` // Required. Page number to retrieve, must be >= 1
	PageSize   int `json:"page_size"`   // Required. Items per page, range 1-100
}

// ============================================================================
// Response Structures
// ============================================================================

// BankAccount represents a bank account response
type BankAccount struct {
	ID              string `json:"id"`                         // UUID. Unique identifier of the bank account
	AccountNumber   string `json:"account_number"`             // Bank account number (e.g., IBAN)
	AccountName     string `json:"account_name"`               // Name of the account holder
	BankName        string `json:"bank_name"`                  // Name of the bank
	SwiftCode       string `json:"swift_code"`                 // SWIFT/BIC code, 8-11 characters
	BankCountryCode string `json:"bank_country_code"`          // ISO 3166-1 alpha-2 country code
	BankAddress     string `json:"bank_address"`               // Physical address of the bank
	Currency        string `json:"currency"`                   // ISO 4217 currency code (e.g., "USD", "GBP")
	BankCodeType    string `json:"bank_code_type,omitempty"`   // Allowed: aba, bank_code, sort_code, bsb_code, ifsc, cnaps_number
	BankCodeValue   string `json:"bank_code_value,omitempty"`  // Routing code value; format varies by bank_code_type
	BankBranchCode  string `json:"bank_branch_code,omitempty"` // Bank branch code
	AccountStatus   string `json:"account_status"`             // Valid or Invalid
}

// ListBankAccountsResponse represents a paginated list of bank accounts
type ListBankAccountsResponse struct {
	TotalPages int           `json:"total_pages"` // Total number of pages available
	TotalItems int           `json:"total_items"` // Total count of available items
	Data       []BankAccount `json:"data"`        // List of bank account records
}

// ============================================================================
// API Methods
// ============================================================================

// Create creates a new bank account for settlement purposes
// Optional RequestOptions can be provided to set custom headers like x-idempotency-key or x-auth-token
func (c *BankAccountsClient) Create(ctx context.Context, req *CreateBankAccountRequest, opts ...*common.RequestOptions) (*BankAccount, error) {
	var resp BankAccount
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.PostWithOptions(ctx, "/v2/payment/bankaccount/create", req, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to create bank account: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific bank account by ID
// Optional RequestOptions can be provided to set custom headers like x-idempotency-key or x-auth-token
func (c *BankAccountsClient) Get(ctx context.Context, id string, opts ...*common.RequestOptions) (*BankAccount, error) {
	var resp BankAccount
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	path := fmt.Sprintf("/v2/payment/bankaccount/%s", id)
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to get bank account: %w", err)
	}
	return &resp, nil
}

// Update updates an existing bank account
// Optional RequestOptions can be provided to set custom headers like x-idempotency-key or x-auth-token
func (c *BankAccountsClient) Update(ctx context.Context, id string, req *UpdateBankAccountRequest, opts ...*common.RequestOptions) (*BankAccount, error) {
	var resp BankAccount
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	path := fmt.Sprintf("/v2/payment/bankaccount/%s", id)
	if err := c.client.PostWithOptions(ctx, path, req, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to update bank account: %w", err)
	}
	return &resp, nil
}

// List retrieves a paginated list of all bank accounts
// Optional RequestOptions can be provided to set custom headers like x-idempotency-key or x-auth-token
func (c *BankAccountsClient) List(ctx context.Context, req *ListBankAccountsRequest, opts ...*common.RequestOptions) (*ListBankAccountsResponse, error) {
	var resp ListBankAccountsResponse
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	path := fmt.Sprintf("/v2/payment/bankaccount?page_number=%d&page_size=%d", req.PageNumber, req.PageSize)
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to list bank accounts: %w", err)
	}
	return &resp, nil
}
