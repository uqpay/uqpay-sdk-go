package issuing

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// BalancesClient handles issuing balance operations
type BalancesClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// RetrieveBalanceRequest represents a request to retrieve issuing balance
type RetrieveBalanceRequest struct {
	Currency string `json:"currency"` // Required. ISO 4217 currency code, e.g. "USD"
}

// ListBalancesRequest represents a request to list issuing balances
type ListBalancesRequest struct {
	PageSize   int `json:"page_size"`   // Required. Items per page, min: 10, max: 100, default: 10
	PageNumber int `json:"page_number"` // Required. Page number to retrieve, min: 1, default: 1
}

// ListBalanceTransactionsRequest represents a request to list issuing balance transactions
type ListBalanceTransactionsRequest struct {
	PageSize          int    `json:"page_size"`                    // Required. Items per page, min: 10, max: 100, default: 10
	PageNumber        int    `json:"page_number"`                  // Required. Page number to retrieve, min: 1, default: 1
	StartTime         string `json:"start_time,omitempty"`         // Optional. ISO 8601 format, max 90-day interval with EndTime
	EndTime           string `json:"end_time,omitempty"`           // Optional. ISO 8601 format, max 90-day interval with StartTime
	Currency          string `json:"currency,omitempty"`           // Optional. ISO 4217 currency code filter, e.g. "USD"
	TransactionType   string `json:"transaction_type,omitempty"`   // Optional. DEPOSIT, TRANSFER_IN, TRANSFER_OUT, ISSUING_AUTHORIZATION, ISSUING_REVERSAL, ISSUING_REFUND, CARD_RECHARGE, CARD_WITHDRAW, SETTLEMENT_DEBIT, SETTLEMENT_CREDIT, SETTLEMENT_REVERSAL, FEE, REFUND, ADJUSTMENT, FUNDS_TRANSFER_IN, FUNDS_TRANSFER_OUT, FEE_REFUND, FEE_DEDUCTION, MARGIN_PAYMENT, MARGIN_REFUND, OTHER
	TransactionStatus string `json:"transaction_status,omitempty"` // Optional. COMPLETED, PENDING, or FAILED
	TransactionID     string `json:"transaction_id,omitempty"`     // Optional. UUID, filter by specific transaction
}

// ============================================================================
// Response Structures
// ============================================================================

// IssuingBalance represents an issuing account balance
type IssuingBalance struct {
	BalanceID        string `json:"balance_id"`        // Unique identifier for the account balance
	Currency         string `json:"currency"`          // ISO 4217 currency code
	AvailableBalance string `json:"available_balance"` // Currently accessible funds
	MarginBalance    string `json:"margin_balance"`    // Additional borrowing capacity
	FrozenBalance    string `json:"frozen_balance"`    // Temporarily locked funds
	CreateTime       string `json:"create_time"`       // ISO 8601 creation timestamp
	LastTradeTime    string `json:"last_trade_time"`   // ISO 8601 most recent update timestamp
	BalanceStatus    string `json:"balance_status"`    // ACTIVE, PENDING, PROCESSING, or CLOSED
}

// ListBalancesResponse represents a paginated list of issuing balances
type ListBalancesResponse struct {
	TotalPages int              `json:"total_pages"` // Total pages of available items
	TotalItems int              `json:"total_items"` // Total count of available items
	Data       []IssuingBalance `json:"data"`        // Collection of balance objects
}

// IssuingBalanceTransaction represents an issuing balance transaction
type IssuingBalanceTransaction struct {
	TransactionID      string `json:"transaction_id"`       // UUID, unique identifier for the transaction
	ShortTransactionID string `json:"short_transaction_id"` // Short-form transaction identifier
	AccountID          string `json:"account_id"`           // Account identifier
	AccountName        string `json:"account_name"`         // Account name
	BalanceID          string `json:"balance_id"`           // Associated account balance identifier
	TransactionType    string `json:"transaction_type"`     // DEPOSIT, TRANSFER_IN, TRANSFER_OUT, ISSUING_AUTHORIZATION, ISSUING_REVERSAL, ISSUING_REFUND, CARD_RECHARGE, CARD_WITHDRAW, SETTLEMENT_DEBIT, SETTLEMENT_CREDIT, SETTLEMENT_REVERSAL, FEE, REFUND, ADJUSTMENT, FUNDS_TRANSFER_IN, FUNDS_TRANSFER_OUT, FEE_REFUND, FEE_DEDUCTION, MARGIN_PAYMENT, MARGIN_REFUND, OTHER
	Currency           string `json:"currency"`             // ISO 4217 three-letter currency code
	Amount             string `json:"amount"`               // Transaction amount
	CreateTime         string `json:"create_time"`          // ISO 8601 creation timestamp
	CompleteTime       string `json:"complete_time"`        // ISO 8601 completion timestamp
	TransactionStatus  string `json:"transaction_status"`   // COMPLETED, PENDING, or FAILED
	EndingBalance      string `json:"ending_balance"`       // Balance after the transaction
	Description        string `json:"description"`          // Transaction description
}

// ListBalanceTransactionsResponse represents a paginated list of issuing balance transactions
type ListBalanceTransactionsResponse struct {
	TotalPages int                         `json:"total_pages"` // Total pages of available items
	TotalItems int                         `json:"total_items"` // Total count of available items
	Data       []IssuingBalanceTransaction `json:"data"`        // Collection of transaction objects
}

// ============================================================================
// API Methods
// ============================================================================

// Retrieve retrieves the issuing balance for a specific currency
func (c *BalancesClient) Retrieve(ctx context.Context, req *RetrieveBalanceRequest) (*IssuingBalance, error) {
	var resp IssuingBalance
	if err := c.client.Post(ctx, "/v1/issuing/balances", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to retrieve issuing balance: %w", err)
	}
	return &resp, nil
}

// List lists all issuing balances with pagination
func (c *BalancesClient) List(ctx context.Context, req *ListBalancesRequest) (*ListBalancesResponse, error) {
	var resp ListBalancesResponse
	path := fmt.Sprintf("/v1/issuing/balances?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list issuing balances: %w", err)
	}
	return &resp, nil
}

// ListTransactions lists issuing balance transactions with pagination and optional time filters
func (c *BalancesClient) ListTransactions(ctx context.Context, req *ListBalanceTransactionsRequest) (*ListBalanceTransactionsResponse, error) {
	var resp ListBalanceTransactionsResponse
	path := fmt.Sprintf("/v1/issuing/balances/transactions?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

	if req.StartTime != "" {
		path += fmt.Sprintf("&start_time=%s", req.StartTime)
	}
	if req.EndTime != "" {
		path += fmt.Sprintf("&end_time=%s", req.EndTime)
	}
	if req.Currency != "" {
		path += fmt.Sprintf("&currency=%s", req.Currency)
	}
	if req.TransactionType != "" {
		path += fmt.Sprintf("&transaction_type=%s", req.TransactionType)
	}
	if req.TransactionStatus != "" {
		path += fmt.Sprintf("&transaction_status=%s", req.TransactionStatus)
	}
	if req.TransactionID != "" {
		path += fmt.Sprintf("&transaction_id=%s", req.TransactionID)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list issuing balance transactions: %w", err)
	}
	return &resp, nil
}
