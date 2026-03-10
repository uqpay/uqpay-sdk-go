package payment

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// PaymentRefundsClient handles payment refund operations
type PaymentRefundsClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// CreateRefundRequest represents a refund creation request
type CreateRefundRequest struct {
	PaymentIntentID  string            `json:"payment_intent_id"`            // Required: The ID of the payment intent to refund
	PaymentAttemptID string            `json:"payment_attempt_id,omitempty"` // Optional: The ID of the payment attempt to refund
	Amount           string            `json:"amount"`                       // Required: The amount to refund
	Reason           string            `json:"reason"`                       // Required: The reason for the refund (max 100 chars)
	Metadata         map[string]string `json:"metadata,omitempty"`           // Optional: Additional metadata for the refund
}

// ListRefundsRequest represents a refunds list request
type ListRefundsRequest struct {
	PageSize        int    `json:"page_size"`         // Required. Number of items per page (1-100)
	PageNumber      int    `json:"page_number"`       // Required. Page number to retrieve (1-based)
	StartTime       string `json:"start_time"`        // Optional. Filter start time, ISO 8601 format. Default range is 1 month
	EndTime         string `json:"end_time"`          // Optional. Filter end time, ISO 8601 format. Max range is 3 months
	PaymentIntentID string `json:"payment_intent_id"` // Optional. Filter by payment intent ID
	MerchantOrderID string `json:"merchant_order_id"` // Optional. Filter by merchant reference ID from merchant's system
}

// ============================================================================
// Response Structures
// ============================================================================

// Refund represents a refund response
type Refund struct {
	PaymentRefundID  string            `json:"payment_refund_id"`            // Unique identifier for the refund
	PaymentAttemptID string            `json:"payment_attempt_id,omitempty"` // ID of the payment attempt that was refunded
	Amount           string            `json:"amount,omitempty"`             // Refund amount as decimal string, e.g. "10.01"
	Currency         string            `json:"currency,omitempty"`           // ISO 4217 three-letter currency code, e.g. "USD"
	RefundStatus     string            `json:"refund_status,omitempty"`      // INITIATED, PROCESSING, SUCCEEDED, FAILED, REVERSAL_INITIATED, REVERSAL_PROCESSING, or REVERSAL_SUCCEEDED
	Reason           string            `json:"reason,omitempty"`             // Reason for the refund
	Metadata         map[string]string `json:"metadata,omitempty"`           // User-defined key-value pairs
	CreateTime       string            `json:"create_time,omitempty"`        // Refund creation time, ISO 8601 format
	UpdateTime       string            `json:"update_time,omitempty"`        // Last update time, ISO 8601 format
}

// ListRefundsResponse represents a paginated list of refunds
type ListRefundsResponse struct {
	TotalPages int      `json:"total_pages"` // Total number of available pages
	TotalItems int      `json:"total_items"` // Total count of available items
	Data       []Refund `json:"data"`        // Array of refund objects
}

// ============================================================================
// API Methods
// ============================================================================

// Create creates a new refund for a completed payment
func (c *PaymentRefundsClient) Create(ctx context.Context, req *CreateRefundRequest) (*Refund, error) {
	var resp Refund
	opts := &common.RequestOptions{
		ClientID: c.client.Config.ClientID,
	}
	if err := c.client.PostWithOptions(ctx, "/v2/payment/refunds", req, &resp, opts); err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific refund by ID
func (c *PaymentRefundsClient) Get(ctx context.Context, refundID string) (*Refund, error) {
	var resp Refund
	path := fmt.Sprintf("/v2/payment/refunds/%s", refundID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get refund: %w", err)
	}
	return &resp, nil
}

// List returns a paginated list of refunds with optional filters
func (c *PaymentRefundsClient) List(ctx context.Context, req *ListRefundsRequest) (*ListRefundsResponse, error) {
	var resp ListRefundsResponse

	path := "/v2/payment/refunds"
	separator := "?"

	if req.PageSize > 0 {
		path += fmt.Sprintf("%spage_size=%d", separator, req.PageSize)
		separator = "&"
	}
	if req.PageNumber > 0 {
		path += fmt.Sprintf("%spage_number=%d", separator, req.PageNumber)
		separator = "&"
	}
	if req.StartTime != "" {
		path += fmt.Sprintf("%sstart_time=%s", separator, req.StartTime)
		separator = "&"
	}
	if req.EndTime != "" {
		path += fmt.Sprintf("%send_time=%s", separator, req.EndTime)
		separator = "&"
	}
	if req.PaymentIntentID != "" {
		path += fmt.Sprintf("%spayment_intent_id=%s", separator, req.PaymentIntentID)
		separator = "&"
	}
	if req.MerchantOrderID != "" {
		path += fmt.Sprintf("%smerchant_order_id=%s", separator, req.MerchantOrderID)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list refunds: %w", err)
	}
	return &resp, nil
}
