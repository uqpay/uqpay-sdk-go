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
	Currency string `json:"currency"`
}

// ListBalancesRequest represents a request to list issuing balances
type ListBalancesRequest struct {
	PageSize   int `json:"page_size"`   // required, 10-100
	PageNumber int `json:"page_number"` // required, >=1
}

// ListBalanceTransactionsRequest represents a request to list issuing balance transactions
type ListBalanceTransactionsRequest struct {
	PageSize   int    `json:"page_size"`   // required, 10-100
	PageNumber int    `json:"page_number"` // required, >=1
	StartTime  string `json:"start_time"`  // optional, max 90 days interval
	EndTime    string `json:"end_time"`    // optional, max 90 days interval
}

// ============================================================================
// Response Structures
// ============================================================================

// IssuingBalance represents an issuing account balance
type IssuingBalance struct {
	BalanceID        string  `json:"balance_id"`
	Currency         string  `json:"currency"`
	AvailableBalance float64 `json:"available_balance"`
	MarginBalance    float64 `json:"margin_balance"`
	FrozenBalance    float64 `json:"frozen_balance"`
	CreateTime       string  `json:"create_time"`
	LastTradeTime    string  `json:"last_trade_time"`
	BalanceStatus    string  `json:"balance_status"` // ACTIVE, PENDING, PROCESSING, CLOSED
}

// ListBalancesResponse represents a paginated list of issuing balances
type ListBalancesResponse struct {
	TotalPages int              `json:"total_pages"`
	TotalItems int              `json:"total_items"`
	Data       []IssuingBalance `json:"data"`
}

// IssuingBalanceTransaction represents an issuing balance transaction
type IssuingBalanceTransaction struct {
	TransactionID      string  `json:"transaction_id"`
	ShortTransactionID string  `json:"short_transaction_id"`
	AccountID          string  `json:"account_id"`
	BalanceID          string  `json:"balance_id"`
	TransactionType    string  `json:"transaction_type"` // DEPOSIT, TRANSFER_IN, TRANSFER_OUT, ISSUING_AUTHORIZATION, ISSUING_REVERSAL, ISSUING_REFUND, CARD_RECHARGE, CARD_WITHDRAW, SETTLEMENT_DEBIT, SETTLEMENT_CREDIT, SETTLEMENT_REVERSAL, FEE, REFUND, ADJUSTMENT, FUNDS_TRANSFER_IN, FUNDS_TRANSFER_OUT, FEE_REFUND, FEE_DEDUCTION, MARGIN_PAYMENT, MARGIN_REFUND, OTHER
	Currency           string  `json:"currency"`
	Amount             float64 `json:"amount"`
	CreateTime         string  `json:"create_time"`
	CompleteTime       string  `json:"complete_time"`
	TransactionStatus  string  `json:"transaction_status"` // FAILED, PENDING, COMPLETED, CANCELLED
	EndingBalance      float64 `json:"ending_balance"`
	Description        string  `json:"description"`
}

// ListBalanceTransactionsResponse represents a paginated list of issuing balance transactions
type ListBalanceTransactionsResponse struct {
	TotalPages int                         `json:"total_pages"`
	TotalItems int                         `json:"total_items"`
	Data       []IssuingBalanceTransaction `json:"data"`
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

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list issuing balance transactions: %w", err)
	}
	return &resp, nil
}
