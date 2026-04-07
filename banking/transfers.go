package banking

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// TransfersClient handles transfer operations
type TransfersClient struct {
	client *common.APIClient
}

// Transfer represents a transfer between accounts
type Transfer struct {
	TransferID             string `json:"transfer_id"`
	ReferenceID            string `json:"reference_id"`
	ShortReferenceID       string `json:"short_reference_id"`
	SourceAccountName      string `json:"source_account_name"`
	DestinationAccountName string `json:"destination_account_name"`
	TransferCurrency       string `json:"transfer_currency"`
	TransferAmount         string `json:"transfer_amount"`
	TransferReason         string `json:"transfer_reason"`
	TransferStatus         string `json:"transfer_status"`
	CreatedBy              string `json:"created_by"`
	CreateTime             string `json:"create_time"`
	CompleteTime           string `json:"complete_time"`
}

// ListTransfersRequest represents a transfer list request
type ListTransfersRequest struct {
	PageSize       int    `json:"page_size"`       // 10-100
	PageNumber     int    `json:"page_number"`     // >=1
	StartTime      string `json:"start_time"`      // optional, ISO8601
	EndTime        string `json:"end_time"`        // optional, ISO8601
	TransferStatus string `json:"transfer_status"` // optional: completed, failed
	Currency       string `json:"currency"`        // optional
}

// ListTransfersResponse represents a transfer list response
type ListTransfersResponse struct {
	TotalPages int        `json:"total_pages"`
	TotalItems int        `json:"total_items"`
	Data       []Transfer `json:"data"`
}

// CreateTransferRequest represents a transfer creation request
type CreateTransferRequest struct {
	SourceAccountID string `json:"source_account_id"` // required, UUID
	TargetAccountID string `json:"target_account_id"` // required, UUID
	Currency        string `json:"currency"`          // required
	Amount          string `json:"amount"`            // required
	Reason          string `json:"reason"`            // required
}

// CreateTransferResponse represents a transfer creation response
type CreateTransferResponse struct {
	TransferID       string `json:"transfer_id"`
	ShortReferenceID string `json:"short_reference_id"`
}

// List lists transfers
func (c *TransfersClient) List(ctx context.Context, req *ListTransfersRequest) (*ListTransfersResponse, error) {
	var resp ListTransfersResponse
	params := url.Values{}
	params.Set("page_size", strconv.Itoa(req.PageSize))
	params.Set("page_number", strconv.Itoa(req.PageNumber))
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
func (c *TransfersClient) Create(ctx context.Context, req *CreateTransferRequest) (*CreateTransferResponse, error) {
	var resp CreateTransferResponse
	if err := c.client.Post(ctx, "/v1/transfer", req, &resp); err != nil {
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
