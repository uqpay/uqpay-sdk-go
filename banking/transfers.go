package banking

import (
	"context"
	"fmt"
	"net/url"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// TransfersClient handles transfer operations
type TransfersClient struct {
	client *common.APIClient
}

// Transfer represents a transfer between accounts
type Transfer struct {
	TransferID             string `json:"transfer_id"`              // Unique identifier for the transfer, UUID format
	ReferenceID            string `json:"reference_id"`             // System-generated reference, e.g. P220406-LLCVLRM
	ShortReferenceID       string `json:"short_reference_id"`       // Abbreviated reference identifier
	SourceAccountName      string `json:"source_account_name"`      // Name of the account sending funds
	DestinationAccountName string `json:"destination_account_name"` // Name of the account receiving funds
	TransferCurrency       string `json:"transfer_currency"`        // ISO 4217 currency code, e.g. USD, SGD
	TransferAmount         string `json:"transfer_amount"`          // Amount of money transferred
	TransferStatus         string `json:"transfer_status"`          // completed or failed
	CreatedBy              string `json:"created_by"`               // User or system that initiated the transfer
	CreateTime             string `json:"create_time"`              // Timestamp when request was initiated, ISO 8601 format
	CompleteTime           string `json:"complete_time"`            // Timestamp when request was completed, ISO 8601 format
}

// ListTransfersRequest represents a transfer list request
type ListTransfersRequest struct {
	PageSize       int    `json:"page_size"`       // Required. Min: 10, Max: 100
	PageNumber     int    `json:"page_number"`     // Required. Min: 1
	StartTime      string `json:"start_time"`      // Optional. Filter by created_time start, ISO 8601 format
	EndTime        string `json:"end_time"`        // Optional. Filter by created_time end, ISO 8601 format
	TransferStatus string `json:"transfer_status"` // Optional. completed or failed
	Currency       string `json:"currency"`        // Optional. ISO 4217 currency code
}

// ListTransfersResponse represents a transfer list response
type ListTransfersResponse struct {
	TotalPages int        `json:"total_pages"` // Total number of pages available
	TotalItems int        `json:"total_items"` // Total count of available items
	Data       []Transfer `json:"data"`        // Array of transfer records
}

// CreateTransferRequest represents a transfer creation request
type CreateTransferRequest struct {
	SourceAccountID string `json:"source_account_id"` // Required. UUID of the source account
	TargetAccountID string `json:"target_account_id"` // Required. UUID of the target account
	Currency        string `json:"currency"`          // Required. ISO 4217 currency code, e.g. USD
	Amount          string `json:"amount"`            // Required. Transfer amount
	Reason          string `json:"reason"`            // Required. Reason for the transfer
}

// CreateTransferResponse represents a transfer creation response
type CreateTransferResponse struct {
	TransferID       string `json:"transfer_id"`        // UUID of the created transfer
	ShortReferenceID string `json:"short_reference_id"` // Human-readable reference, e.g. P220406-LLCVLRM
}

// List lists transfers
func (c *TransfersClient) List(ctx context.Context, req *ListTransfersRequest) (*ListTransfersResponse, error) {
	var resp ListTransfersResponse

	params := url.Values{}
	params.Set("page_size", fmt.Sprintf("%d", req.PageSize))
	params.Set("page_number", fmt.Sprintf("%d", req.PageNumber))

	if req.StartTime != "" {
		params.Set("start_time", req.StartTime)
	}
	if req.EndTime != "" {
		params.Set("end_time", req.EndTime)
	}
	if req.TransferStatus != "" {
		params.Set("transfer_status", req.TransferStatus)
	}
	if req.Currency != "" {
		params.Set("currency", req.Currency)
	}

	path := "/v1/transfer?" + params.Encode()

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list transfers: %w", err)
	}
	return &resp, nil
}

// Create creates a new transfer
// Optional RequestOptions can be provided to set custom headers like x-idempotency-key or x-on-behalf-of
func (c *TransfersClient) Create(ctx context.Context, req *CreateTransferRequest, opts ...*common.RequestOptions) (*CreateTransferResponse, error) {
	var resp CreateTransferResponse
	var opt *common.RequestOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if err := c.client.PostWithOptions(ctx, "/v1/transfer", req, &resp, opt); err != nil {
		return nil, fmt.Errorf("failed to create transfer: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific transfer
func (c *TransfersClient) Get(ctx context.Context, transferID string) (*Transfer, error) {
	var resp Transfer
	path := fmt.Sprintf("/v1/transfer/%s", transferID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get transfer: %w", err)
	}
	return &resp, nil
}
