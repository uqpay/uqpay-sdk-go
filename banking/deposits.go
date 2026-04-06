package banking

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// DepositsClient handles deposit operations
type DepositsClient struct {
	client *common.APIClient
}

// Deposit represents a deposit transaction
type Deposit struct {
	DepositID              string         `json:"deposit_id"`
	ShortReferenceID       string         `json:"short_reference_id"`
	Currency               string         `json:"currency"`
	Amount                 string         `json:"amount"`
	DepositFee             string         `json:"deposit_fee"`
	DepositStatus          string         `json:"deposit_status"`
	ReceiverAccountNumber  string         `json:"receiver_account_number"`
	DepositReference       string         `json:"deposit_reference"`
	Sender                 *DepositSender `json:"sender,omitempty"`
	CreateTime             string         `json:"create_time"`
	CompleteTime           string         `json:"complete_time"`
}

// DepositSender represents the sender information for a deposit
type DepositSender struct {
	SenderName          string `json:"sender_name"`
	SenderCountry       string `json:"sender_country"`
	SenderAccountNumber string `json:"sender_account_number"`
	SenderSwiftCode     string `json:"sender_swift_code"`
	SenderAddress       string `json:"sender_address"`
}

// ListDepositsRequest represents a deposit list request
type ListDepositsRequest struct {
	PageSize      int    `json:"page_size"`      // required, 10-100
	PageNumber    int    `json:"page_number"`    // required, >=1
	StartTime     string `json:"start_time"`     // optional, ISO8601
	EndTime       string `json:"end_time"`       // optional, ISO8601
	DepositStatus string `json:"deposit_status"` // optional: PENDING, COMPLETED, FAILED
	Currency      string `json:"currency"`       // optional
}

// ListDepositsResponse represents a deposit list response
type ListDepositsResponse struct {
	TotalPages int       `json:"total_pages"`
	TotalItems int       `json:"total_items"`
	Data       []Deposit `json:"data"`
}

// List lists deposits
func (c *DepositsClient) List(ctx context.Context, req *ListDepositsRequest) (*ListDepositsResponse, error) {
	var resp ListDepositsResponse
	params := url.Values{}
	params.Set("page_size", strconv.Itoa(req.PageSize))
	params.Set("page_number", strconv.Itoa(req.PageNumber))
	if req.StartTime != "" {
		params.Set("start_time", req.StartTime)
	}
	if req.EndTime != "" {
		params.Set("end_time", req.EndTime)
	}
	if req.DepositStatus != "" {
		params.Set("deposit_status", req.DepositStatus)
	}
	if req.Currency != "" {
		params.Set("currency", req.Currency)
	}
	path := "/v1/deposit?" + params.Encode()

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list deposits: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific deposit
func (c *DepositsClient) Get(ctx context.Context, depositID string) (*Deposit, error) {
	var resp Deposit
	path := fmt.Sprintf("/v1/deposit/%s", depositID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get deposit: %w", err)
	}
	return &resp, nil
}
