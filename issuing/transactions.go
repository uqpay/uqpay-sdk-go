package issuing

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// TransactionsClient handles transaction operations
type TransactionsClient struct {
	client *common.APIClient
}

// MerchantData represents merchant details associated with a transaction
type MerchantData struct {
	CategoryCode string `json:"category_code"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Name         string `json:"name"`
}

// Transaction represents a card transaction
type Transaction struct {
	TransactionID         string        `json:"transaction_id"`
	CardID                string        `json:"card_id"`
	CardNumber            string        `json:"card_number"`
	CardholderID          string        `json:"cardholder_id"`
	TransactionType       string        `json:"transaction_type"`
	TransactionAmount     string        `json:"transaction_amount"`
	TransactionCurrency   string        `json:"transaction_currency"`
	BillingAmount         string        `json:"billing_amount"`
	BillingCurrency       string        `json:"billing_currency"`
	TransactionFee        string        `json:"transaction_fee"`
	TransactionFeeCurrency string       `json:"transaction_fee_currency"`
	FeePassThrough        string        `json:"fee_pass_through"` // Y or N
	CardAvailableBalance  string        `json:"card_available_balance"`
	AuthorizationCode     string        `json:"authorization_code"`
	ShortTransactionID    string        `json:"short_transaction_id"`
	OriginalTransactionID string        `json:"original_transaction_id"`
	TransactionStatus     string        `json:"transaction_status"` // APPROVED, DECLINED, PENDING
	TransactionTime       string        `json:"transaction_time"`
	PostedTime            *string       `json:"posted_time,omitempty"`
	MerchantData          *MerchantData `json:"merchant_data,omitempty"`
	Description           string        `json:"description"`
	WalletType            *string       `json:"wallet_type,omitempty"`
}

// ListTransactionsRequest represents a transaction list request
type ListTransactionsRequest struct {
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
	CardID     string `json:"card_id,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
}

// ListTransactionsResponse represents a transaction list response
type ListTransactionsResponse struct {
	TotalPages int           `json:"total_pages"`
	TotalItems int           `json:"total_items"`
	Data       []Transaction `json:"data"`
}

// Get retrieves a transaction by ID
func (c *TransactionsClient) Get(ctx context.Context, transactionID string) (*Transaction, error) {
	var transaction Transaction
	path := fmt.Sprintf("/v1/issuing/transactions/%s", transactionID)
	if err := c.client.Get(ctx, path, &transaction); err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return &transaction, nil
}

// List lists transactions
func (c *TransactionsClient) List(ctx context.Context, req *ListTransactionsRequest) (*ListTransactionsResponse, error) {
	var resp ListTransactionsResponse
	params := url.Values{}
	params.Set("page_size", strconv.Itoa(req.PageSize))
	params.Set("page_number", strconv.Itoa(req.PageNumber))
	if req.CardID != "" {
		params.Set("card_id", req.CardID)
	}
	if req.StartTime != "" {
		params.Set("start_time", req.StartTime)
	}
	if req.EndTime != "" {
		params.Set("end_time", req.EndTime)
	}
	path := "/v1/issuing/transactions?" + params.Encode()
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	return &resp, nil
}
