package issuing

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// TransactionsClient handles transaction operations
type TransactionsClient struct {
	client *common.APIClient
}

// Transaction represents a card transaction
type Transaction struct {
	TransactionID          string        `json:"transaction_id"`           // Required. UUID. Unique identifier for the transaction
	ShortTransactionID     string        `json:"short_transaction_id"`     // Required. Short unique identifier, e.g. "CT2024-03-01"
	OriginalTransactionID  string        `json:"original_transaction_id"`  // Required. UUID. Identifier for the original transaction
	CardID                 string        `json:"card_id"`                  // Required. UUID. Unique card identifier
	CardNumber             string        `json:"card_number"`              // Required. Masked card number, e.g. "************5668"
	CardholderID           string        `json:"cardholder_id"`            // Required. UUID. The cardholder's unique identifier
	TransactionType        string        `json:"transaction_type"`         // Required. AUTHORIZATION, REFUND, FUND COLLECTION, ATM DEPOSIT, REVERSAL, VALIDATION, SETTLEMENT DEBIT, SETTLEMENT CREDIT, SETTLEMENT REVERSAL, CHARGEBACK DEBIT, or CHARGEBACK CREDIT
	TransactionAmount      string        `json:"transaction_amount"`       // Required. Transaction amount in original currency
	TransactionCurrency    string        `json:"transaction_currency"`     // Required. ISO 4217 currency code, e.g. "USD"
	BillingAmount          string        `json:"billing_amount"`           // Required. Billing amount
	BillingCurrency        string        `json:"billing_currency"`         // Required. ISO 4217 billing currency code, e.g. "SGD"
	TransactionFee         string        `json:"transaction_fee"`          // Required. Fee amount charged for the transaction
	TransactionFeeCurrency string        `json:"transaction_fee_currency"` // Required. ISO 4217 fee currency code, e.g. "SGD"
	FeePassThrough         string        `json:"fee_pass_through"`         // Required. Whether fee was deducted from card balance. Y or N
	CardAvailableBalance   string        `json:"card_available_balance"`   // Required. Available card balance after transaction
	AuthorizationCode      string        `json:"authorization_code"`       // Required. Authorization approval code
	MerchantData           *MerchantData `json:"merchant_data"`            // Required. Merchant information object
	MerchantName           string        `json:"merchant_name"`            // Merchant business name (legacy field)
	Description            string        `json:"description"`              // Required. Additional context based on transaction status (failure reasons or remarks)
	TransactionStatus      string        `json:"transaction_status"`       // Required. APPROVED, DECLINED, or PENDING
	TransactionTime        string        `json:"transaction_time"`         // Required. Transaction occurrence time. ISO 8601 format
	PostedTime             string        `json:"posted_time,omitempty"`    // Optional. Transaction posted time. ISO 8601 format
	WalletType             string        `json:"wallet_type,omitempty"`    // Optional. Digital wallet used. ApplePay, GooglePay, GOOGLE ECOMMERCE, GOOGLE, or GOOGLE PAY
}

// MerchantData represents merchant information in a card transaction
type MerchantData struct {
	CategoryCode string `json:"category_code"` // Required. Merchant category code, e.g. "6011"
	Name         string `json:"name"`          // Required. Merchant business name
	City         string `json:"city"`          // Optional. City where the merchant is located
	Country      string `json:"country"`       // Optional. Country where the merchant is located, e.g. "CN"
}

// ListTransactionsRequest represents a transaction list request
type ListTransactionsRequest struct {
	PageSize   int    `json:"page_size"`            // Required. Number of items per page. Min: 10, Max: 100, Default: 10
	PageNumber int    `json:"page_number"`          // Required. Page number to retrieve. Min: 1, Default: 1
	CardID     string `json:"card_id,omitempty"`    // Optional. UUID. Filter transactions by card identifier
	StartTime  string `json:"start_time,omitempty"` // Optional. Earliest transaction time. ISO 8601 format. Max interval with end_time: 90 days
	EndTime    string `json:"end_time,omitempty"`   // Optional. Latest transaction time. ISO 8601 format. Max interval with start_time: 90 days
}

// ListTransactionsResponse represents a transaction list response
type ListTransactionsResponse struct {
	TotalPages int           `json:"total_pages"` // Total number of pages available
	TotalItems int           `json:"total_items"` // Total count of available items
	Data       []Transaction `json:"data"`        // Array of transaction objects
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
	path := fmt.Sprintf("/v1/issuing/transactions?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)
	if req.CardID != "" {
		path = fmt.Sprintf("%s&card_id=%s", path, req.CardID)
	}
	if req.StartTime != "" {
		path = fmt.Sprintf("%s&start_time=%s", path, req.StartTime)
	}
	if req.EndTime != "" {
		path = fmt.Sprintf("%s&end_time=%s", path, req.EndTime)
	}
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	return &resp, nil
}
