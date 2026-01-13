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
	// AccountNumber is the bank account number (e.g., IBAN). Required.
	AccountNumber string `json:"account_number"`
	// BankName is the name of the bank. Required.
	BankName string `json:"bank_name"`
	// SwiftCode is the SWIFT/BIC code of the bank. Required.
	SwiftCode string `json:"swift_code"`
	// BankCountryCode is the two-letter country code (ISO 3166-1 alpha-2). Required.
	BankCountryCode string `json:"bank_country_code"`
	// BankAddress is the physical address of the bank. Required.
	BankAddress string `json:"bank_address"`
	// Currency is the currency of the settlement account. Required.
	Currency string `json:"currency"`
	// BankCodeType is the type of bank routing code. Required based on currency and country:
	// - aba: Required when currency = "USD" and bank_country_code = "US"
	// - bank_code: Required when currency = "CAD" or currency = "HKD"
	// - sort_code: Required when currency = "GBP"
	// - bsb_code: Required when currency = "AUD"
	// - ifsc: Required when currency = "INR"
	// - cnaps_number: Required when currency = "CNH" and bank_country_code = "CN"
	BankCodeType string `json:"bank_code_type,omitempty"`
	// BankCodeValue is the bank identifier value. Format depends on bank_code_type:
	// - aba: Exactly 9 digits
	// - bank_code: Exactly 3 digits
	// - sort_code: Exactly 6 digits
	// - bsb_code: Exactly 6 digits
	// - ifsc: 11 characters (4 letters + 0 + 6 alphanumeric)
	// - cnaps_number: Exactly 12 digits
	BankCodeValue string `json:"bank_code_value,omitempty"`
	// BankBranchCode is the bank branch code. Required when currency = "CAD".
	BankBranchCode string `json:"bank_branch_code,omitempty"`
}

// UpdateBankAccountRequest represents a bank account update request
type UpdateBankAccountRequest struct {
	// AccountNumber is the bank account number (e.g., IBAN). Required.
	AccountNumber string `json:"account_number"`
	// BankName is the name of the bank. Required.
	BankName string `json:"bank_name"`
	// SwiftCode is the SWIFT/BIC code of the bank. Required.
	SwiftCode string `json:"swift_code"`
	// BankCountryCode is the two-letter country code (ISO 3166-1 alpha-2). Required.
	BankCountryCode string `json:"bank_country_code"`
	// BankAddress is the physical address of the bank. Required.
	BankAddress string `json:"bank_address"`
	// BankCodeType is the type of bank routing code. Required based on currency and country:
	// - aba: Required when currency = "USD" and bank_country_code = "US"
	// - bank_code: Required when currency = "CAD" or currency = "HKD"
	// - sort_code: Required when currency = "GBP"
	// - bsb_code: Required when currency = "AUD"
	// - ifsc: Required when currency = "INR"
	// - cnaps_number: Required when currency = "CNH" and bank_country_code = "CN"
	BankCodeType string `json:"bank_code_type,omitempty"`
	// BankCodeValue is the bank identifier value. Format depends on bank_code_type:
	// - aba: Exactly 9 digits
	// - bank_code: Exactly 3 digits
	// - sort_code: Exactly 6 digits
	// - bsb_code: Exactly 6 digits
	// - ifsc: 11 characters (4 letters + 0 + 6 alphanumeric)
	// - cnaps_number: Exactly 12 digits
	BankCodeValue string `json:"bank_code_value,omitempty"`
	// BankBranchCode is the bank branch code. Required when currency = "CAD".
	BankBranchCode string `json:"bank_branch_code,omitempty"`
}

// ListBankAccountsRequest represents a bank accounts list request
type ListBankAccountsRequest struct {
	// PageNumber is the page number to retrieve (must be >= 1). Required.
	PageNumber int `json:"page_number"`
	// PageSize is the maximum number of items per page (1-100). Required.
	PageSize int `json:"page_size"`
}

// ============================================================================
// Response Structures
// ============================================================================

// BankAccount represents a bank account response
type BankAccount struct {
	// ID is the unique identifier (UUID) of the bank account.
	ID string `json:"id"`
	// AccountNumber is the bank account number (e.g., IBAN).
	AccountNumber string `json:"account_number"`
	// AccountName is the name of the account holder.
	AccountName string `json:"account_name"`
	// BankName is the name of the bank.
	BankName string `json:"bank_name"`
	// SwiftCode is the SWIFT/BIC code of the bank.
	SwiftCode string `json:"swift_code"`
	// BankCountryCode is the two-letter country code (ISO 3166-1 alpha-2).
	BankCountryCode string `json:"bank_country_code"`
	// BankAddress is the physical address of the bank.
	BankAddress string `json:"bank_address"`
	// Currency is the currency of the settlement account.
	Currency string `json:"currency"`
	// BankCodeType is the type of bank routing code.
	BankCodeType string `json:"bank_code_type,omitempty"`
	// BankCodeValue is the bank identifier value.
	BankCodeValue string `json:"bank_code_value,omitempty"`
	// BankBranchCode is the bank branch code.
	BankBranchCode string `json:"bank_branch_code,omitempty"`
	// AccountStatus is the status of the bank account (Valid or Invalid).
	AccountStatus string `json:"account_status"`
}

// ListBankAccountsResponse represents a paginated list of bank accounts
type ListBankAccountsResponse struct {
	// TotalPages is the total number of pages available.
	TotalPages int `json:"total_pages"`
	// TotalItems is the total count of available items.
	TotalItems int `json:"total_items"`
	// Data is the list of bank accounts.
	Data []BankAccount `json:"data"`
}

// ============================================================================
// API Methods
// ============================================================================

// Create creates a new bank account for settlement purposes
func (c *BankAccountsClient) Create(ctx context.Context, req *CreateBankAccountRequest) (*BankAccount, error) {
	var resp BankAccount
	if err := c.client.Post(ctx, "/v2/payment/bankaccount/create", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create bank account: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific bank account by ID
func (c *BankAccountsClient) Get(ctx context.Context, id string) (*BankAccount, error) {
	var resp BankAccount
	path := fmt.Sprintf("/v2/payment/bankaccount/%s", id)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get bank account: %w", err)
	}
	return &resp, nil
}

// Update updates an existing bank account
func (c *BankAccountsClient) Update(ctx context.Context, id string, req *UpdateBankAccountRequest) (*BankAccount, error) {
	var resp BankAccount
	path := fmt.Sprintf("/v2/payment/bankaccount/%s", id)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update bank account: %w", err)
	}
	return &resp, nil
}

// List retrieves a paginated list of all bank accounts
func (c *BankAccountsClient) List(ctx context.Context, req *ListBankAccountsRequest) (*ListBankAccountsResponse, error) {
	var resp ListBankAccountsResponse
	path := fmt.Sprintf("/v2/payment/bankaccount?page_number=%d&page_size=%d", req.PageNumber, req.PageSize)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list bank accounts: %w", err)
	}
	return &resp, nil
}
