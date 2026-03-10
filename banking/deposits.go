package banking

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// DepositsClient handles deposit operations
type DepositsClient struct {
	client *common.APIClient
}

// Deposit represents a deposit transaction
type Deposit struct {
	DepositID             string         `json:"deposit_id"`                  // Unique deposit identifier
	ShortReferenceID      string         `json:"short_reference_id"`          // System-generated reference code
	Currency              string         `json:"currency"`                    // ISO 4217 three-letter currency code
	Amount                string         `json:"amount"`                      // Gross amount before fees are deducted
	DepositFee            string         `json:"deposit_fee"`                 // Fee amount for the deposit
	DepositStatus         string         `json:"deposit_status"`              // PENDING, COMPLETED, or FAILED
	DepositReference      string         `json:"deposit_reference,omitempty"` // Optional. External reference identifier
	ReceiverAccountNumber string         `json:"receiver_account_number"`     // Destination account that received funds
	Sender                *DepositSender `json:"sender"`                      // Originating party details
	CreateTime            string         `json:"create_time"`                 // ISO 8601 creation timestamp
	CompleteTime          string         `json:"complete_time"`               // ISO 8601 completion timestamp
}

// DepositSender represents the sender information for a deposit
type DepositSender struct {
	SenderName          string `json:"sender_name"`           // Optional. Customer-facing business name of originator
	SenderCountry       string `json:"sender_country"`        // Optional. ISO 3166-1 alpha-2 country code
	SenderAccountNumber string `json:"sender_account_number"` // Optional. Source account number
	SenderSwiftCode     string `json:"sender_swift_code"`     // Optional. SWIFT/BIC code of sender's bank
}

// ListDepositsRequest represents a deposit list request
type ListDepositsRequest struct {
	PageSize   int    `json:"page_size"`   // Required. Max items per page, range: 10-100
	PageNumber int    `json:"page_number"` // Required. Page number to retrieve, min: 1
	StartTime  string `json:"start_time"`  // Optional. ISO 8601 start of creation time range (inclusive)
	EndTime    string `json:"end_time"`    // Optional. ISO 8601 end of creation time range (inclusive)
	Status     string `json:"status"`      // Optional. PENDING, COMPLETED, or FAILED
}

// ListDepositsResponse represents a deposit list response
type ListDepositsResponse struct {
	TotalPages int       `json:"total_pages"` // Total number of pages available
	TotalItems int       `json:"total_items"` // Total count of available items
	Data       []Deposit `json:"data"`        // Collection of deposit records
}

// List lists deposits
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of
func (c *DepositsClient) List(ctx context.Context, req *ListDepositsRequest, opts ...*common.RequestOptions) (*ListDepositsResponse, error) {
	var resp ListDepositsResponse
	path := fmt.Sprintf("/v1/deposit?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

	if req.StartTime != "" {
		path += fmt.Sprintf("&start_time=%s", req.StartTime)
	}
	if req.EndTime != "" {
		path += fmt.Sprintf("&end_time=%s", req.EndTime)
	}
	if req.Status != "" {
		path += fmt.Sprintf("&status=%s", req.Status)
	}

	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to list deposits: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific deposit
// Optional RequestOptions can be provided to set custom headers like x-on-behalf-of
func (c *DepositsClient) Get(ctx context.Context, depositID string, opts ...*common.RequestOptions) (*Deposit, error) {
	var resp Deposit
	path := fmt.Sprintf("/v1/deposit/%s", depositID)
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to get deposit: %w", err)
	}
	return &resp, nil
}
