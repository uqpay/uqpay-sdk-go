package banking

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// BalancesClient handles balance operations
type BalancesClient struct {
	client *common.APIClient
}

// Balance represents account balance information
type Balance struct {
	BalanceID        string `json:"balance_id"`
	Currency         string `json:"currency"`
	AvailableBalance string `json:"available_balance"`
	PrepaidBalance   string `json:"prepaid_balance"`
	MarginBalance    string `json:"margin_balance"`
	FrozenBalance    string `json:"frozen_balance"`
	BalanceStatus    string `json:"balance_status"`
	CreateTime       string `json:"create_time"`
	LastTradeTime    string `json:"last_trade_time"`
}

// ListBalancesRequest represents a balance list request
type ListBalancesRequest struct {
	PageSize   int `json:"page_size"`   // required, 10-100
	PageNumber int `json:"page_number"` // required, >=1
}

// ListBalancesResponse represents a balance list response
type ListBalancesResponse struct {
	TotalPages int       `json:"total_pages"`
	TotalItems int       `json:"total_items"`
	Data       []Balance `json:"data"`
}

// BalanceTransaction represents a balance transaction
type BalanceTransaction struct {
	TransactionID     string `json:"transaction_id"`
	AccountID         string `json:"account_id"`
	BalanceID         string `json:"balance_id"`
	Currency          string `json:"currency"`
	Amount            string `json:"amount"`
	CreditDebitType   string `json:"credit_debit_type"`   // C (credit) or D (debit)
	TransactionType   string `json:"transaction_type"`    // CONVERSION, DEPOSIT, PAYOUT, TRANSFER, FEE, etc.
	TransactionStatus string `json:"transaction_status"`  // COMPLETED, PENDING, FAILED
	TransactionWay    string `json:"transaction_way"`     // API, WEB, etc.
	PayoutWay         string `json:"payout_way,omitempty"`
	ReferenceID       string `json:"reference_id"`
	CreateTime        string `json:"create_time"`
	CompleteTime      string `json:"complete_time"`
}

// ListBalanceTransactionsRequest represents a balance transaction list request
type ListBalanceTransactionsRequest struct {
	PageSize          int    `json:"page_size"`          // required, 10-100
	PageNumber        int    `json:"page_number"`        // required, >=1
	StartTime         string `json:"start_time"`         // optional, ISO8601
	EndTime           string `json:"end_time"`           // optional, ISO8601
	Currency          string `json:"currency"`           // optional
	TransactionType   string `json:"transaction_type"`   // optional: ALL, PAYIN, DEPOSIT, etc.
	TransactionStatus string `json:"transaction_status"` // optional: ALL, COMPLETED, PENDING, FAILED
}

// ListBalanceTransactionsResponse represents a balance transaction list response
type ListBalanceTransactionsResponse struct {
	TotalPages int                  `json:"total_pages"`
	TotalItems int                  `json:"total_items"`
	Data       []BalanceTransaction `json:"data"`
}

// Get retrieves balance for a specific currency
func (c *BalancesClient) Get(ctx context.Context, currency string) (*Balance, error) {
	var resp Balance
	path := fmt.Sprintf("/v1/balances/%s", currency)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return &resp, nil
}

// List lists all balances
func (c *BalancesClient) List(ctx context.Context, req *ListBalancesRequest) (*ListBalancesResponse, error) {
	var resp ListBalancesResponse
	path := fmt.Sprintf("/v1/balances?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list balances: %w", err)
	}
	return &resp, nil
}

// ListTransactions lists balance transactions
func (c *BalancesClient) ListTransactions(ctx context.Context, req *ListBalanceTransactionsRequest) (*ListBalanceTransactionsResponse, error) {
	var resp ListBalanceTransactionsResponse
	path := fmt.Sprintf("/v1/balances/transactions?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

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

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list balance transactions: %w", err)
	}
	return &resp, nil
}
