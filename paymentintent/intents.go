package paymentintent

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// PaymentIntentsClient handles paymentintent intent operations
type PaymentIntentsClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// CreatePaymentIntentRequest represents a paymentintent intent creation request
type CreatePaymentIntentRequest struct {
	Amount          string            `json:"amount"`
	Currency        string            `json:"currency"`
	MerchantOrderID string            `json:"merchant_order_id,omitempty"`
	Description     string            `json:"description,omitempty"`
	ReturnURL       string            `json:"return_url,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	PaymentMethod   *PaymentMethod    `json:"payment_method,omitempty"`
}

// PaymentMethod represents the paymentintent method details
type PaymentMethod struct {
	Type string `json:"type"` // e.g., "card"
	Card *Card  `json:"card,omitempty"`
}

// Card represents card paymentintent details
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

// Billing represents billing information for a card paymentintent
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

// UpdatePaymentIntentRequest represents a paymentintent intent update request
type UpdatePaymentIntentRequest struct {
	Description   string            `json:"description,omitempty"`
	ReturnURL     string            `json:"return_url,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	PaymentMethod *PaymentMethod    `json:"payment_method,omitempty"`
}

// ConfirmPaymentIntentRequest represents a paymentintent intent confirmation request
type ConfirmPaymentIntentRequest struct {
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty"`
	ReturnURL     string         `json:"return_url,omitempty"`
}

// CapturePaymentIntentRequest represents a paymentintent intent capture request
type CapturePaymentIntentRequest struct {
	Amount string `json:"amount,omitempty"` // Optional: for partial capture, specify amount less than authorized
}

// CancelPaymentIntentRequest represents a paymentintent intent cancellation request
type CancelPaymentIntentRequest struct {
	CancellationReason string `json:"cancellation_reason,omitempty"`
}

// ListPaymentIntentsRequest represents a paymentintent intents list request
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

// PaymentIntent represents a paymentintent intent response
type PaymentIntent struct {
	ID              string            `json:"id"`
	Amount          string            `json:"amount"`
	Currency        string            `json:"currency"`
	Status          string            `json:"status"`
	MerchantOrderID string            `json:"merchant_order_id,omitempty"`
	Description     string            `json:"description,omitempty"`
	ReturnURL       string            `json:"return_url,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	PaymentMethod   *PaymentMethod    `json:"payment_method,omitempty"`
	CreatedAt       string            `json:"created_at,omitempty"`
	UpdatedAt       string            `json:"updated_at,omitempty"`
}

// ListPaymentIntentsResponse represents a paginated list of paymentintent intents
type ListPaymentIntentsResponse struct {
	TotalPages int             `json:"total_pages"`
	TotalItems int             `json:"total_items"`
	Data       []PaymentIntent `json:"data"`
}

// ============================================================================
// API Methods
// ============================================================================

// Create creates a new paymentintent intent
func (c *PaymentIntentsClient) Create(ctx context.Context, req *CreatePaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	if err := c.client.Post(ctx, "/v2/payment_intents/create", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create paymentintent intent: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific paymentintent intent by ID
func (c *PaymentIntentsClient) Get(ctx context.Context, paymentIntentID string) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s", paymentIntentID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to get paymentintent intent: %w", err)
	}
	return &resp, nil
}

// Update updates properties on a paymentintent intent without confirming
// Note: Updating payment_method requires subsequent confirmation
func (c *PaymentIntentsClient) Update(ctx context.Context, paymentIntentID string, req *UpdatePaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s", paymentIntentID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update paymentintent intent: %w", err)
	}
	return &resp, nil
}

// Confirm confirms a paymentintent intent for paymentintent authorization
func (c *PaymentIntentsClient) Confirm(ctx context.Context, paymentIntentID string, req *ConfirmPaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s/confirm", paymentIntentID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to confirm paymentintent intent: %w", err)
	}
	return &resp, nil
}

// Capture captures the funds of an uncaptured paymentintent intent
// The paymentintent intent must have status "requires_capture"
func (c *PaymentIntentsClient) Capture(ctx context.Context, paymentIntentID string, req *CapturePaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s/capture", paymentIntentID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to capture paymentintent intent: %w", err)
	}
	return &resp, nil
}

// Cancel cancels a paymentintent intent and prevents further paymentintent attempts
func (c *PaymentIntentsClient) Cancel(ctx context.Context, paymentIntentID string, req *CancelPaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s/cancel", paymentIntentID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to cancel paymentintent intent: %w", err)
	}
	return &resp, nil
}

// List returns a paginated list of paymentintent intents with optional filters
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
		return nil, fmt.Errorf("failed to list paymentintent intents: %w", err)
	}
	return &resp, nil
}
