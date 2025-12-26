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

// ============================================================================
// Response Structures
// ============================================================================

// PaymentIntent represents a payment intent response
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
