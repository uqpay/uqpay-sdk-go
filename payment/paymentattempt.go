package payment

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// PaymentAttemptsClient handles payment attempt operations
type PaymentAttemptsClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// ListPaymentAttemptsRequest represents a payment attempts list request
type ListPaymentAttemptsRequest struct {
	PageSize        int    `json:"page_size"`         // Required. Number of items per page (1-100)
	PageNumber      int    `json:"page_number"`       // Required. Page number, must be >= 1
	PaymentIntentID string `json:"payment_intent_id"` // Optional. Filter by payment intent ID
	AttemptStatus   string `json:"attempt_status"`    // Optional. Filter by status: INITIATED, AUTHENTICATION_REDIRECTED, PENDING_AUTHORIZATION, AUTHORIZED, CAPTURE_REQUESTED, SETTLED, SUCCEEDED, CANCELLED, EXPIRED, FAILED
}

// ============================================================================
// Response Structures
// ============================================================================

// PaymentAttempt represents a payment attempt response
type PaymentAttempt struct {
	AttemptID          string            `json:"attempt_id"`                    // Unique identifier for the attempt (UUID)
	Amount             string            `json:"amount,omitempty"`              // Transaction amount
	Currency           string            `json:"currency,omitempty"`            // ISO 4217 three-letter currency code
	CapturedAmount     string            `json:"captured_amount,omitempty"`     // Funds successfully captured
	RefundedAmount     string            `json:"refunded_amount,omitempty"`     // Funds returned to customer
	AttemptStatus      string            `json:"attempt_status,omitempty"`      // INITIATED, AUTHENTICATION_REDIRECTED, PENDING_AUTHORIZATION, AUTHORIZED, CAPTURE_REQUESTED, SETTLED, SUCCEEDED, CANCELLED, EXPIRED, or FAILED
	CancellationReason string            `json:"cancellation_reason,omitempty"` // Reason for cancelling the payment attempt
	FailureCode        string            `json:"failure_code,omitempty"`        // Error code if payment failed; see Error Code Reference
	PaymentMethod      string            `json:"payment_method,omitempty"`      // Payment method used, e.g. "card"
	Metadata           map[string]string `json:"metadata,omitempty"`            // Key-value metadata pairs
	CreateTime         string            `json:"create_time,omitempty"`         // ISO 8601 creation timestamp
	UpdateTime         string            `json:"update_time,omitempty"`         // ISO 8601 last modification timestamp
	CompleteTime       string            `json:"complete_time,omitempty"`       // ISO 8601 completion timestamp
}

// ListPaymentAttemptsResponse represents a paginated list of payment attempts
type ListPaymentAttemptsResponse struct {
	TotalPages int              `json:"total_pages"` // Total number of available pages
	TotalItems int              `json:"total_items"` // Total count of available items
	Data       []PaymentAttempt `json:"data"`        // List of payment attempt records
}

// ============================================================================
// API Methods
// ============================================================================

// Get retrieves a specific payment attempt by ID
func (c *PaymentAttemptsClient) Get(ctx context.Context, paymentAttemptID string) (*PaymentAttempt, error) {
	var resp PaymentAttempt
	path := fmt.Sprintf("/v2/payment/payment_attempts/%s", paymentAttemptID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get payment attempt: %w", err)
	}
	return &resp, nil
}

// List returns a paginated list of payment attempts with optional filters
func (c *PaymentAttemptsClient) List(ctx context.Context, req *ListPaymentAttemptsRequest) (*ListPaymentAttemptsResponse, error) {
	var resp ListPaymentAttemptsResponse

	path := "/v2/payment/payment_attempts"
	separator := "?"

	if req.PageSize > 0 {
		path += fmt.Sprintf("%spage_size=%d", separator, req.PageSize)
		separator = "&"
	}
	if req.PageNumber > 0 {
		path += fmt.Sprintf("%spage_number=%d", separator, req.PageNumber)
		separator = "&"
	}
	if req.PaymentIntentID != "" {
		path += fmt.Sprintf("%spayment_intent_id=%s", separator, req.PaymentIntentID)
		separator = "&"
	}
	if req.AttemptStatus != "" {
		path += fmt.Sprintf("%sattempt_status=%s", separator, req.AttemptStatus)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list payment attempts: %w", err)
	}
	return &resp, nil
}
