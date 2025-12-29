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
	Amount          string            `json:"amount"`
	Currency        string            `json:"currency"`
	BeneficiaryID   string            `json:"beneficiary_id,omitempty"`
	MerchantOrderID string            `json:"merchant_order_id,omitempty"`
	Description     string            `json:"description,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// ListPayoutsRequest represents a payouts list request
type ListPayoutsRequest struct {
	PageSize   int    `json:"page_size"`   // Number of items per page
	PageNumber int    `json:"page_number"` // Page number (1-based)
	Status     string `json:"status"`      // Filter by status
	StartTime  string `json:"start_time"`  // Filter by creation time (ISO8601)
	EndTime    string `json:"end_time"`    // Filter by creation time (ISO8601)
}

// ============================================================================
// Response Structures
// ============================================================================

// Payout represents a payout response
type Payout struct {
	PayoutID            string            `json:"payout_id"`
	PayoutAmount        string            `json:"payout_amount,omitempty"`
	PayoutCurrency      string            `json:"payout_currency,omitempty"`
	PayoutStatus        string            `json:"payout_status,omitempty"`
	InternalNote        string            `json:"internal_note,omitempty"`
	StatementDescriptor string            `json:"statement_descriptor,omitempty"`
	BeneficiaryID       string            `json:"beneficiary_id,omitempty"`
	MerchantOrderID     string            `json:"merchant_order_id,omitempty"`
	Metadata            map[string]string `json:"metadata,omitempty"`
	CreateTime          string            `json:"create_time,omitempty"`
	CompletedTime       string            `json:"completed_time,omitempty"`
}

// ListPayoutsResponse represents a paginated list of payouts
type ListPayoutsResponse struct {
	TotalPages int      `json:"total_pages"`
	TotalItems int      `json:"total_items"`
	Data       []Payout `json:"data"`
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
	if req.Status != "" {
		path += fmt.Sprintf("%sstatus=%s", separator, req.Status)
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
