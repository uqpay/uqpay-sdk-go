package payment

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// PaymentPayoutsClient handles payment payout operations
type PaymentPayoutsClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// CreatePayoutRequest represents a payout creation request
type CreatePayoutRequest struct {
	PayoutCurrency      string `json:"payout_currency"`         // Required. ISO 4217 three-letter currency code (e.g., "SGD")
	PayoutAmount        string `json:"payout_amount"`           // Required. Withdrawal amount (e.g., "100.00")
	StatementDescriptor string `json:"statement_descriptor"`    // Required. Reference displayed to recipient's bank, max 15 characters
	InternalNote        string `json:"internal_note,omitempty"` // Optional. Internal remark for the payout
}

// ListPayoutsRequest represents a payouts list request
type ListPayoutsRequest struct {
	PageSize     int    `json:"page_size"`     // Required. Number of items per page (1-100)
	PageNumber   int    `json:"page_number"`   // Required. Page number, 1-based
	PayoutStatus string `json:"payout_status"` // Optional. Filter by status: INITIATED, PROCESSING, COMPLETED, FAILED, or FAILED_REFUNDED
	StartTime    string `json:"start_time"`    // Optional. Payout creation start date, format: YYYY-MM-DD (inclusive, 00:00:00)
	EndTime      string `json:"end_time"`      // Optional. Payout creation end date, format: YYYY-MM-DD (inclusive, 23:59:59)
}

// ============================================================================
// Response Structures
// ============================================================================

// Payout represents a payout response
type Payout struct {
	PayoutID            string `json:"payout_id"`                      // Unique payout identifier (e.g., "PO1968582687224500224")
	PayoutAmount        string `json:"payout_amount,omitempty"`        // Withdrawal amount (e.g., "100.00")
	PayoutCurrency      string `json:"payout_currency,omitempty"`      // ISO 4217 three-letter currency code (e.g., "SGD")
	PayoutStatus        string `json:"payout_status,omitempty"`        // INITIATED, PROCESSING, COMPLETED, FAILED, or FAILED_REFUNDED
	InternalNote        string `json:"internal_note,omitempty"`        // Internal remark for the payout
	StatementDescriptor string `json:"statement_descriptor,omitempty"` // Reference displayed to recipient's bank, max 15 characters
	CreateTime          string `json:"create_time,omitempty"`          // Payout creation timestamp, ISO 8601 format
	CompletedTime       string `json:"completed_time,omitempty"`       // Payout completion timestamp, ISO 8601 format; empty if not yet completed
}

// ListPayoutsResponse represents a paginated list of payouts
type ListPayoutsResponse struct {
	TotalPages int      `json:"total_pages"` // Total number of available result pages
	TotalItems int      `json:"total_items"` // Total count of matching payouts
	Data       []Payout `json:"data"`        // List of payout objects for the current page
}

// ============================================================================
// API Methods
// ============================================================================

// Create creates a new payout order
func (c *PaymentPayoutsClient) Create(ctx context.Context, req *CreatePayoutRequest) (*Payout, error) {
	var resp Payout
	if err := c.client.Post(ctx, "/v2/payment/payout/create", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create payout: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific payout by ID
func (c *PaymentPayoutsClient) Get(ctx context.Context, payoutID string) (*Payout, error) {
	var resp Payout
	path := fmt.Sprintf("/v2/payment/payout/%s", payoutID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get payout: %w", err)
	}
	return &resp, nil
}

// List returns a paginated list of payouts with optional filters
// Note: When filtering by date range, max interval is one month
func (c *PaymentPayoutsClient) List(ctx context.Context, req *ListPayoutsRequest) (*ListPayoutsResponse, error) {
	var resp ListPayoutsResponse

	path := "/v2/payment/payout"
	separator := "?"

	if req.PageSize > 0 {
		path += fmt.Sprintf("%spage_size=%d", separator, req.PageSize)
		separator = "&"
	}
	if req.PageNumber > 0 {
		path += fmt.Sprintf("%spage_number=%d", separator, req.PageNumber)
		separator = "&"
	}
	if req.PayoutStatus != "" {
		path += fmt.Sprintf("%spayout_status=%s", separator, req.PayoutStatus)
		separator = "&"
	}
	if req.StartTime != "" {
		path += fmt.Sprintf("%sstart_time=%s", separator, req.StartTime)
		separator = "&"
	}
	if req.EndTime != "" {
		path += fmt.Sprintf("%send_time=%s", separator, req.EndTime)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list payouts: %w", err)
	}
	return &resp, nil
}
