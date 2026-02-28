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
	PageSize        int    `json:"page_size"`         // Number of items per page
	PageNumber      int    `json:"page_number"`       // Page number (1-based)
	PaymentIntentID string `json:"payment_intent_id"` // Filter by payment intent ID
	AttemptStatus   string `json:"attempt_status"`    // Filter by status: INITIATED, AUTHENTICATION_REDIRECTED, PENDING_AUTHORIZATION, AUTHORIZED, CAPTURE_REQUESTED, SETTLED, SUCCEEDED, CANCELLED, EXPIRED, FAILED
}

// ============================================================================
// Response Structures
// ============================================================================

// PaymentAttempt represents a payment attempt response
type PaymentAttempt struct {
	AttemptID          string            `json:"attempt_id"`
	Amount             string            `json:"amount,omitempty"`
	Currency           string            `json:"currency,omitempty"`
	CapturedAmount     string            `json:"captured_amount,omitempty"`
	RefundedAmount     string            `json:"refunded_amount,omitempty"`
	AttemptStatus      string            `json:"attempt_status,omitempty"`
	CancellationReason string            `json:"cancellation_reason,omitempty"`
	FailureCode        string            `json:"failure_code,omitempty"`
	PaymentMethod      *PaymentMethod    `json:"payment_method,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	CreateTime         string            `json:"create_time,omitempty"`
	UpdateTime         string            `json:"update_time,omitempty"`
	CompleteTime       string            `json:"complete_time,omitempty"`
}

// ListPaymentAttemptsResponse represents a paginated list of payment attempts
type ListPaymentAttemptsResponse struct {
	TotalPages int              `json:"total_pages"`
	TotalItems int              `json:"total_items"`
	Data       []PaymentAttempt `json:"data"`
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
