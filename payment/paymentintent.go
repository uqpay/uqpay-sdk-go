package payment

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// PaymentIntentsClient handles payment intent operations
type PaymentIntentsClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// CreatePaymentIntentRequest represents a payment intent creation request
type CreatePaymentIntentRequest struct {
	Amount          string            `json:"amount"`
	Currency        string            `json:"currency"`
	MerchantOrderID string            `json:"merchant_order_id,omitempty"`
	Description     string            `json:"description,omitempty"`
	ReturnURL       string            `json:"return_url,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	PaymentMethod   *PaymentMethod    `json:"payment_method,omitempty"`
}

// PaymentMethod represents the payment method details
type PaymentMethod struct {
	Type string `json:"type"` // e.g., "card"
	Card *Card  `json:"card,omitempty"`
}

// Card represents card payment details
type Card struct {
	CardName          string   `json:"card_name,omitempty"`
	CardNumber        string   `json:"card_number,omitempty"`
	ExpiryMonth       string   `json:"expiry_month,omitempty"`
	ExpiryYear        string   `json:"expiry_year,omitempty"`
	CVC               string   `json:"cvc,omitempty"`
	Network           string   `json:"network,omitempty"` // e.g., "mastercard", "visa"
	Billing           *Billing `json:"billing,omitempty"`
	AutoCapture       *bool    `json:"auto_capture,omitempty"`
	AuthorizationType string   `json:"authorization_type,omitempty"`
	ThreeDSAction     string   `json:"three_ds_action,omitempty"`
}

// Billing represents billing information for a card payment
type Billing struct {
	FirstName   string   `json:"first_name,omitempty"`
	LastName    string   `json:"last_name,omitempty"`
	Email       string   `json:"email,omitempty"`
	PhoneNumber string   `json:"phone_number,omitempty"`
	Address     *Address `json:"address,omitempty"`
}

// Address represents a billing address
type Address struct {
	CountryCode string `json:"country_code,omitempty"`
	State       string `json:"state,omitempty"`
	City        string `json:"city,omitempty"`
	Street      string `json:"street,omitempty"`
	Postcode    string `json:"postcode,omitempty"`
}

// UpdatePaymentIntentRequest represents a payment intent update request
type UpdatePaymentIntentRequest struct {
	Description   string            `json:"description,omitempty"`
	ReturnURL     string            `json:"return_url,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	PaymentMethod *PaymentMethod    `json:"payment_method,omitempty"`
}

// ConfirmPaymentIntentRequest represents a payment intent confirmation request
type ConfirmPaymentIntentRequest struct {
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty"`
	ReturnURL     string         `json:"return_url,omitempty"`
}

// CapturePaymentIntentRequest represents a payment intent capture request
type CapturePaymentIntentRequest struct {
	Amount string `json:"amount,omitempty"` // Optional: for partial capture, specify amount less than authorized
}

// CancelPaymentIntentRequest represents a payment intent cancellation request
type CancelPaymentIntentRequest struct {
	CancellationReason string `json:"cancellation_reason,omitempty"` // e.g., "requested_by_customer", "duplicate", "fraudulent", "abandoned"
}

// ListPaymentIntentsRequest represents a payment intents list request
type ListPaymentIntentsRequest struct {
	PageSize   int    `json:"page_size"`   // Number of items per page (default: 10)
	PageNumber int    `json:"page_number"` // Page number (1-based)
	Status     string `json:"status"`      // Filter by status: requires_payment_method, requires_confirmation, requires_action, processing, requires_capture, succeeded, canceled
	StartTime  string `json:"start_time"`  // Filter by creation time (ISO8601)
	EndTime    string `json:"end_time"`    // Filter by creation time (ISO8601)
	Currency   string `json:"currency"`    // Filter by currency
}

// ============================================================================
// Response Structures
// ============================================================================

// PaymentIntent represents a payment intent response
type PaymentIntent struct {
	PaymentIntentID             string                 `json:"payment_intent_id"`
	Amount                      string                 `json:"amount"`
	Currency                    string                 `json:"currency"`
	IntentStatus                string                 `json:"intent_status"`
	MerchantOrderID             string                 `json:"merchant_order_id,omitempty"`
	Description                 string                 `json:"description,omitempty"`
	ReturnURL                   string                 `json:"return_url,omitempty"`
	Metadata                    map[string]string      `json:"metadata,omitempty"`
	AvailablePaymentMethodTypes []string               `json:"available_payment_method_types,omitempty"`
	CapturedAmount              string                 `json:"captured_amount,omitempty"`
	ClientSecret                string                 `json:"client_secret,omitempty"`
	CancellationReason          string                 `json:"cancellation_reason,omitempty"`
	LatestPaymentAttempt        map[string]interface{} `json:"latest_payment_attempt,omitempty"`
	NextAction                  map[string]interface{} `json:"next_action,omitempty"`
	CreateTime                  string                 `json:"create_time,omitempty"`
	UpdateTime                  string                 `json:"update_time,omitempty"`
	CancelTime                  string                 `json:"cancel_time,omitempty"`
	CompleteTime                string                 `json:"complete_time,omitempty"`
}

// ListPaymentIntentsResponse represents a paginated list of payment intents
type ListPaymentIntentsResponse struct {
	TotalPages int             `json:"total_pages"`
	TotalItems int             `json:"total_items"`
	Data       []PaymentIntent `json:"data"`
}

// ============================================================================
// API Methods
// ============================================================================

// Create creates a new payment intent
func (c *PaymentIntentsClient) Create(ctx context.Context, req *CreatePaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	if err := c.client.Post(ctx, "/v2/payment_intents/create", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific payment intent by ID
func (c *PaymentIntentsClient) Get(ctx context.Context, paymentIntentID string) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s", paymentIntentID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}
	return &resp, nil
}

// Update updates properties on a payment intent without confirming
// Note: Updating payment_method requires subsequent confirmation
func (c *PaymentIntentsClient) Update(ctx context.Context, paymentIntentID string, req *UpdatePaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s", paymentIntentID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update payment intent: %w", err)
	}
	return &resp, nil
}

// Confirm confirms a payment intent for payment authorization
func (c *PaymentIntentsClient) Confirm(ctx context.Context, paymentIntentID string, req *ConfirmPaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s/confirm", paymentIntentID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}
	return &resp, nil
}

// Capture captures the funds of an uncaptured payment intent
// The payment intent must have status "requires_capture"
func (c *PaymentIntentsClient) Capture(ctx context.Context, paymentIntentID string, req *CapturePaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s/capture", paymentIntentID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to capture payment intent: %w", err)
	}
	return &resp, nil
}

// Cancel cancels a payment intent and prevents further payment attempts
func (c *PaymentIntentsClient) Cancel(ctx context.Context, paymentIntentID string, req *CancelPaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s/cancel", paymentIntentID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to cancel payment intent: %w", err)
	}
	return &resp, nil
}

// List returns a paginated list of payment intents with optional filters
func (c *PaymentIntentsClient) List(ctx context.Context, req *ListPaymentIntentsRequest) (*ListPaymentIntentsResponse, error) {
	var resp ListPaymentIntentsResponse

	path := "/v2/payment_intents"
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
		separator = "&"
	}
	if req.Currency != "" {
		path += fmt.Sprintf("%scurrency=%s", separator, req.Currency)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list payment intents: %w", err)
	}
	return &resp, nil
}
