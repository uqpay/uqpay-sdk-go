package banking

import (
	"context"
	"fmt"
	"net/url"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// BalancesClient handles balance operations
type BalancesClient struct {
	client *common.APIClient
}

// Balance represents account balance information
type Balance struct {
	BalanceID        string `json:"balance_id"`        // Required. Unique balance identifier (UUID)
	Currency         string `json:"currency"`          // Required. ISO 4217 currency code (e.g., USD)
	AvailableBalance string `json:"available_balance"` // Required. Funds accessible for transactions
	PrepaidBalance   string `json:"prepaid_balance"`   // Required. Prepaid account balance
	MarginBalance    string `json:"margin_balance"`    // Required. Margin balance available
	FrozenBalance    string `json:"frozen_balance"`    // Required. Funds restricted from use
	CreateTime       string `json:"create_time"`       // Required. Record creation timestamp (ISO 8601)
	LastTradeTime    string `json:"last_trade_time"`   // Required. Most recent transaction timestamp (ISO 8601)
	BalanceStatus    string `json:"balance_status"`    // Required. ACTIVE, PENDING, PROCESSING, or CLOSED
}

// ListBalancesRequest represents a balance list request
type ListBalancesRequest struct {
	PageSize   int `json:"page_size"`   // Required. Items per page (10-100)
	PageNumber int `json:"page_number"` // Required. Page to retrieve (>=1)
}

// ListBalancesResponse represents a balance list response
type ListBalancesResponse struct {
	TotalPages int       `json:"total_pages"` // Total number of pages available
	TotalItems int       `json:"total_items"` // Total count of available items
	Data       []Balance `json:"data"`        // Array of balance objects
}

// BalanceTransaction represents a balance transaction
type BalanceTransaction struct {
	TransactionID     string `json:"transaction_id"`            // Required. Unique transaction identifier (UUID)
	AccountID         string `json:"account_id,omitempty"`      // Optional. Associated account ID
	BalanceID         string `json:"balance_id"`                // Required. Associated balance ID
	TransactionType   string `json:"transaction_type"`          // Required. DEPOSIT, PAYOUT, TRANSFER, CONVERSION, FEE, REFUND, ADJUSTMENT, or INVOICE
	Currency          string `json:"currency"`                  // Required. ISO 4217 currency code (e.g., USD)
	Amount            string `json:"amount"`                    // Required. Transaction amount
	CreditDebitType   string `json:"credit_debit_type"`         // Required. C (Credit) or D (Debit)
	CreateTime        string `json:"create_time"`               // Required. Record creation timestamp (ISO 8601)
	CompleteTime      string `json:"complete_time"`             // Required. Completion timestamp (ISO 8601)
	ReferenceID       string `json:"reference_id"`              // Required. Reference identifier
	TransactionStatus string `json:"transaction_status"`        // Required. FAILED, PENDING, COMPLETED, or CANCELLED
	TransactionWay    string `json:"transaction_way,omitempty"` // Optional. How transaction was initiated (e.g., API)
}

// ListBalanceTransactionsRequest represents a balance transaction list request
type ListBalanceTransactionsRequest struct {
	PageSize          int    `json:"page_size"`          // Required. Items per page (10-100)
	PageNumber        int    `json:"page_number"`        // Required. Page to retrieve (>=1)
	StartTime         string `json:"start_time"`         // Optional. Start of time range, inclusive (ISO 8601)
	EndTime           string `json:"end_time"`           // Optional. End of time range, inclusive (ISO 8601)
	Currency          string `json:"currency"`           // Optional. Filter by ISO 4217 currency code
	TransactionType   string `json:"transaction_type"`   // Optional. ALL, PAYIN, DEPOSIT, PAYOUT, TRANSFER, CONVERSION, FEE, REFUND, ADJUSTMENT, or INVOICE
	TransactionStatus string `json:"transaction_status"` // Optional. ALL, COMPLETED, PENDING, or FAILED
}

// ListBalanceTransactionsResponse represents a balance transaction list response
type ListBalanceTransactionsResponse struct {
	TotalPages int                  `json:"total_pages"` // Total number of pages available
	TotalItems int                  `json:"total_items"` // Total count of available items
	Data       []BalanceTransaction `json:"data"`        // Array of balance transaction objects
}

// Get retrieves balance for a specific currency
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of
func (c *BalancesClient) Get(ctx context.Context, currency string, opts ...*common.RequestOptions) (*Balance, error) {
	var resp Balance
	path := fmt.Sprintf("/v1/balances/%s", currency)
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return &resp, nil
}

// List lists all balances
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of
func (c *BalancesClient) List(ctx context.Context, req *ListBalancesRequest, opts ...*common.RequestOptions) (*ListBalancesResponse, error) {
	if req.PageSize < 10 || req.PageSize > 100 {
		return nil, fmt.Errorf("page_size must be between 10 and 100, got %d", req.PageSize)
	}
	if req.PageNumber < 1 {
		return nil, fmt.Errorf("page_number must be >= 1, got %d", req.PageNumber)
	}

	var resp ListBalancesResponse
	params := url.Values{}
	params.Set("page_size", fmt.Sprintf("%d", req.PageSize))
	params.Set("page_number", fmt.Sprintf("%d", req.PageNumber))

	path := "/v1/balances?" + params.Encode()

	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to list balances: %w", err)
	}
	return &resp, nil
}

// ListTransactions lists balance transactions
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of
func (c *BalancesClient) ListTransactions(ctx context.Context, req *ListBalanceTransactionsRequest, opts ...*common.RequestOptions) (*ListBalanceTransactionsResponse, error) {
	if req.PageSize < 10 || req.PageSize > 100 {
		return nil, fmt.Errorf("page_size must be between 10 and 100, got %d", req.PageSize)
	}
	if req.PageNumber < 1 {
		return nil, fmt.Errorf("page_number must be >= 1, got %d", req.PageNumber)
	}

	var resp ListBalanceTransactionsResponse
	params := url.Values{}
	params.Set("page_size", fmt.Sprintf("%d", req.PageSize))
	params.Set("page_number", fmt.Sprintf("%d", req.PageNumber))

	if req.StartTime != "" {
		params.Set("start_time", req.StartTime)
	}
	if req.EndTime != "" {
		params.Set("end_time", req.EndTime)
	}
	if req.Currency != "" {
		params.Set("currency", req.Currency)
	}
	if req.TransactionType != "" {
		params.Set("transaction_type", req.TransactionType)
	}
	if req.TransactionStatus != "" {
		params.Set("transaction_status", req.TransactionStatus)
	}

	path := "/v1/balances/transactions?" + params.Encode()

	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to list balance transactions: %w", err)
	}
	return &resp, nil
}
