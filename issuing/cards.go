package issuing

import (
	"context"
	"fmt"

	"github.com/uqpay/uqpay-sdk-go/common"
)

// CardsClient handles card operations
type CardsClient struct {
	client *common.APIClient
}

// ============================================================================
// Request Structures
// ============================================================================

// CreateCardRequest represents a card creation request
type CreateCardRequest struct {
	CardLimit        *string           `json:"card_limit,omitempty"`        // Optional. Total credit limit for the card, min 0.01 for some BINs
	CardCurrency     string            `json:"card_currency"`               // Required. Card currency, e.g. SGD, USD
	CardholderID     string            `json:"cardholder_id"`               // Required. UUID of the cardholder
	CardProductID    string            `json:"card_product_id"`             // Required. UUID of the card product
	SpendingControls []SpendingControl `json:"spending_controls,omitempty"` // Optional. Rules that control spending for this card
	RiskControls     *RiskControls     `json:"risk_controls,omitempty"`     // Optional. User-customized risk control settings
	Metadata         map[string]string `json:"metadata,omitempty"`          // Optional. Key-value pairs, max 3200 bytes
}

// SpendingControl represents spending control rules for a card
type SpendingControl struct {
	Amount   string `json:"amount"`   // Required. Maximum amount allowed per interval, min 0
	Interval string `json:"interval"` // Required. Interval type: PER_TRANSACTION
}

// RiskControls represents user-customized risk control settings
type RiskControls struct {
	Allow3DSTransactions *string  `json:"allow_3ds_transactions,omitempty"` // Optional. Y or N, defaults to Y
	AllowedMCC           []string `json:"allowed_mcc,omitempty"`            // Optional. Whitelist of Merchant Category Codes
	BlockedMCC           []string `json:"blocked_mcc,omitempty"`            // Optional. Blacklist of Merchant Category Codes
}

// CardUpdateRequest represents a card update request
type CardUpdateRequest struct {
	CardLimit          *string           `json:"card_limit,omitempty"`            // Optional. Credit limit for the card, min 0, up to 2 decimals
	NoPINPaymentAmount *string           `json:"no_pin_payment_amount,omitempty"` // Optional. Max amount for transactions without PIN, default 200 SGD
	SpendingControls   []SpendingControl `json:"spending_controls,omitempty"`     // Optional. Rules controlling card spending
	RiskControls       *RiskControls     `json:"risk_controls,omitempty"`         // Optional. User-customized risk control settings
	Metadata           map[string]string `json:"metadata,omitempty"`              // Optional. Key-value pairs, max 3200 bytes
}

// UpdateCardStatusRequest represents a card status update request
type UpdateCardStatusRequest struct {
	CardStatus   string  `json:"card_status"`             // Required. Target status: ACTIVE, FROZEN, or CANCELLED
	UpdateReason *string `json:"update_reason,omitempty"` // Optional. Reason for the status change, max 100 chars
}

// CardOrderRequest represents a card recharge/withdraw request
type CardOrderRequest struct {
	Amount string `json:"amount"` // Required. Recharge or withdraw amount, must be > 0
}

// ActivateCardRequest represents a card activation request
type ActivateCardRequest struct {
	CardID             string  `json:"card_id"`                         // Required. UUID of the card to activate
	ActivationCode     string  `json:"activation_code"`                 // Required. Activation code for the card
	PIN                string  `json:"pin"`                             // Required. 6-digit numeric PIN for the card
	NoPINPaymentAmount *string `json:"no_pin_payment_amount,omitempty"` // Optional. Max amount without PIN verification, default 200 SGD
}

// SetPINRequest represents a card PIN reset request
type SetPINRequest struct {
	CardID string `json:"card_id"` // Required. UUID of the card
	PIN    string `json:"pin"`     // Required. New 6-digit numeric PIN
}

// AssignCardRequest represents a card assignment request
type AssignCardRequest struct {
	CardholderID string `json:"cardholder_id"` // Required. UUID of the cardholder to assign to
	CardNumber   string `json:"card_number"`   // Required. Card number to assign
	CardCurrency string `json:"card_currency"` // Required. Currency for the card, e.g. SGD, USD
	CardMode     string `json:"card_mode"`     // Required. SINGLE (prepaid) or SHARE (debit/prepaid)
}

// BulkCardCreationRequest represents a bulk card creation request
type BulkCardCreationRequest struct {
	CardBIN string `json:"card_bin"` // Required. Card BIN (Bank Identification Number) prefix
	Numbers int    `json:"numbers"`  // Required. Number of cards to create, range 1-5000
}

// ListCardsRequest represents a card list request
type ListCardsRequest struct {
	PageSize     int     `json:"page_size"`               // Required. Items per page, min 10, max 100, default 10
	PageNumber   int     `json:"page_number"`             // Required. Page to retrieve, min 1, default 1
	CardNumber   *string `json:"card_number,omitempty"`   // Optional. Filter by full card number
	CardStatus   *string `json:"card_status,omitempty"`   // Optional. PENDING|ACTIVE|FROZEN|BLOCKED|CANCELLED|LOST|STOLEN|FAILED
	CardholderID *string `json:"cardholder_id,omitempty"` // Optional. Filter by cardholder UUID
}

// ============================================================================
// Response Structures
// ============================================================================

// CardCreationResponse represents the response after creating a card
type CardCreationResponse struct {
	CardID      string `json:"card_id"`       // UUID of the created card
	CardOrderID string `json:"card_order_id"` // UUID of the card order
	CreateTime  string `json:"create_time"`   // ISO 8601 creation timestamp
	CardStatus  string `json:"card_status"`   // PENDING|ACTIVE|FROZEN|BLOCKED|CANCELLED|LOST|STOLEN|FAILED
	OrderStatus string `json:"order_status"`  // PENDING|PROCESSING|SUCCESS|FAILED
}

// CardUpdatedResponse represents the response after updating a card
type CardUpdatedResponse struct {
	CardID      string `json:"card_id"`       // UUID of the updated card
	CardOrderID string `json:"card_order_id"` // UUID of the card order
	CardStatus  string `json:"card_status"`   // PENDING|ACTIVE|FROZEN|BLOCKED|CANCELLED|LOST|STOLEN|FAILED
	OrderStatus string `json:"order_status"`  // PENDING|PROCESSING|SUCCESS|FAILED
}

// CardStatusResponse represents the response after updating card status
type CardStatusResponse struct {
	CardID       string  `json:"card_id"`                 // UUID of the card
	CardOrderID  string  `json:"card_order_id"`           // UUID of the card order
	OrderStatus  string  `json:"order_status"`            // PENDING|PROCESSING|SUCCESS|FAILED
	UpdateReason *string `json:"update_reason,omitempty"` // Reason provided for the status change
}

// RetrieveCardResponse represents detailed card information
type RetrieveCardResponse struct {
	CardID             string            `json:"card_id"`                     // UUID of the card
	CardBIN            string            `json:"card_bin"`                    // Card number prefix (BIN)
	CardScheme         string            `json:"card_scheme"`                 // Payment scheme, e.g. VISA
	CardCurrency       string            `json:"card_currency"`               // Card currency, e.g. SGD, USD
	CardNumber         string            `json:"card_number"`                 // Masked card number
	FormFactor         string            `json:"form_factor"`                 // VIRTUAL or PHYSICAL
	ModeType           string            `json:"mode_type"`                   // SINGLE (prepaid) or SHARE (debit/prepaid)
	CardProductID      string            `json:"card_product_id"`             // UUID of the associated card product
	CardLimit          string            `json:"card_limit"`                  // Credit limit assigned to the card
	AvailableBalance   string            `json:"available_balance"`           // Current available balance
	Cardholder         CardholderInfo    `json:"cardholder"`                  // Associated cardholder details
	SpendingControls   []SpendingControl `json:"spending_controls,omitempty"` // Transaction spending limit rules
	NoPINPaymentAmount string            `json:"no_pin_payment_amount"`       // Max amount allowed without PIN verification
	RiskControls       *RiskControls     `json:"risk_controls,omitempty"`     // User-customized risk control settings
	Metadata           interface{}       `json:"metadata,omitempty"`          // Key-value pairs, max 3200 bytes. May be object or empty string.
	CardStatus         string            `json:"card_status"`                 // PENDING|ACTIVE|FROZEN|BLOCKED|CANCELLED|LOST|STOLEN|FAILED
	UpdateReason       *string           `json:"update_reason,omitempty"`     // Reason for card status change
	ConsumedAmount     *string           `json:"consumed_amount,omitempty"`   // Cumulative amount already used against card limit
}

// CardholderInfo represents cardholder information in card response
type CardholderInfo struct {
	CardholderID     string  `json:"cardholder_id"`           // UUID of the cardholder
	Email            string  `json:"email"`                   // Cardholder email address
	NumberOfCards    int     `json:"number_of_cards"`         // Total cards associated with this cardholder
	FirstName        string  `json:"first_name"`              // Cardholder first name
	LastName         string  `json:"last_name"`               // Cardholder last name
	CreateTime       string  `json:"create_time"`             // Creation timestamp
	CardholderStatus string  `json:"cardholder_status"`       // FAILED|PENDING|SUCCESS|INCOMPLETE
	DateOfBirth      *string `json:"date_of_birth,omitempty"` // Format: yyyy-mm-dd
	CountryCode      *string `json:"country_code,omitempty"`  // ISO 3166-1 alpha-2 country code
	PhoneNumber      *string `json:"phone_number,omitempty"`  // Contact phone number
}

// SecureCardInfo represents secure card information
type SecureCardInfo struct {
	CVV        string `json:"cvv"`         // 3-digit Card Verification Value
	ExpireDate string `json:"expire_date"` // Card expiry date in MM/YY format
	CardNumber string `json:"card_number"` // Full 16-digit Primary Account Number
}

// CardOrder represents a card order
type CardOrder struct {
	CardID       string `json:"card_id"`       // UUID of the card
	CardOrderID  string `json:"card_order_id"` // UUID of the card order
	OrderType    string `json:"order_type"`    // CARD_CREATE|CARD_RECHARGE|CARD_WITHDRAW|CARD_UPDATE
	Amount       string `json:"amount"`        // Transaction amount
	CardCurrency string `json:"card_currency"` // Currency code, e.g. SGD, USD
	CreateTime   string `json:"create_time"`   // ISO 8601 creation timestamp
	UpdateTime   string `json:"update_time"`   // ISO 8601 last update timestamp
	CompleteTime string `json:"complete_time"` // ISO 8601 completion timestamp
	OrderStatus  string `json:"order_status"`  // PENDING|PROCESSING|SUCCESS|FAILED
}

// ActivateCardResponse represents the response after activating a card
type ActivateCardResponse struct {
	RequestStatus string `json:"request_status"` // Request completion status, e.g. SUCCESS
}

// SetPINResponse represents the response after resetting PIN
type SetPINResponse struct {
	RequestStatus string `json:"request_status"` // Request completion status, e.g. SUCCESS
}

// AssignCardResponse represents the response after assigning a card
type AssignCardResponse struct {
	CardID      string `json:"card_id"`       // UUID of the assigned card
	CardOrderID string `json:"card_order_id"` // UUID of the card order
	CreateTime  string `json:"create_time"`   // ISO 8601 creation timestamp
	CardStatus  string `json:"card_status"`   // PENDING|ACTIVE|FROZEN|BLOCKED|CANCELLED|LOST|STOLEN|FAILED
	OrderStatus string `json:"order_status"`  // PENDING|PROCESSING|SUCCESS|FAILED
}

// BulkCardCreationResponse represents the response after bulk card creation
type BulkCardCreationResponse struct {
	ReportID   string  `json:"report_id"`             // UUID of the bulk creation report
	ExpireDate *string `json:"expire_date,omitempty"` // Report expiration date
}

// CreatePANTokenResponse represents the response after creating a PAN token
type CreatePANTokenResponse struct {
	Token string `json:"token"` // One-time PAN token, expires after 60 seconds
}

// ListCardsResponse represents a card list response
type ListCardsResponse struct {
	TotalPages int                    `json:"total_pages"` // Total number of pages available
	TotalItems int                    `json:"total_items"` // Total number of items available
	Data       []RetrieveCardResponse `json:"data"`        // Array of card objects
}

// ============================================================================
// API Methods
// ============================================================================

// Create creates a new card
func (c *CardsClient) Create(ctx context.Context, req *CreateCardRequest) (*CardCreationResponse, error) {
	var resp CardCreationResponse
	if err := c.client.Post(ctx, "/v1/issuing/cards", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}
	return &resp, nil
}

// Update updates the specified issuing card
func (c *CardsClient) Update(ctx context.Context, cardID string, req *CardUpdateRequest) (*CardUpdatedResponse, error) {
	var resp CardUpdatedResponse
	path := fmt.Sprintf("/v1/issuing/cards/%s", cardID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}
	return &resp, nil
}

// Get retrieves a card by ID
func (c *CardsClient) Get(ctx context.Context, cardID string) (*RetrieveCardResponse, error) {
	var card RetrieveCardResponse
	path := fmt.Sprintf("/v1/issuing/cards/%s", cardID)
	if err := c.client.Get(ctx, path, &card); err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}
	return &card, nil
}

// GetSecure retrieves secure card information
func (c *CardsClient) GetSecure(ctx context.Context, cardID string) (*SecureCardInfo, error) {
	var info SecureCardInfo
	path := fmt.Sprintf("/v1/issuing/cards/%s/secure", cardID)
	if err := c.client.Get(ctx, path, &info); err != nil {
		return nil, fmt.Errorf("failed to get secure card info: %w", err)
	}
	return &info, nil
}

// List lists cards with pagination and filters
func (c *CardsClient) List(ctx context.Context, req *ListCardsRequest) (*ListCardsResponse, error) {
	var resp ListCardsResponse
	path := fmt.Sprintf("/v1/issuing/cards?page_size=%d&page_number=%d", req.PageSize, req.PageNumber)

	if req.CardNumber != nil {
		path += fmt.Sprintf("&card_number=%s", *req.CardNumber)
	}
	if req.CardStatus != nil {
		path += fmt.Sprintf("&card_status=%s", *req.CardStatus)
	}
	if req.CardholderID != nil {
		path += fmt.Sprintf("&cardholder_id=%s", *req.CardholderID)
	}

	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("failed to list cards: %w", err)
	}
	return &resp, nil
}

// UpdateStatus updates card status
func (c *CardsClient) UpdateStatus(ctx context.Context, cardID string, req *UpdateCardStatusRequest) (*CardStatusResponse, error) {
	var resp CardStatusResponse
	path := fmt.Sprintf("/v1/issuing/cards/%s/status", cardID)
	if err := c.client.Post(ctx, path, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to update card status: %w", err)
	}
	return &resp, nil
}

// Recharge recharges a card
func (c *CardsClient) Recharge(ctx context.Context, cardID string, req *CardOrderRequest) (*CardOrder, error) {
	var order CardOrder
	path := fmt.Sprintf("/v1/issuing/cards/%s/recharge", cardID)
	if err := c.client.Post(ctx, path, req, &order); err != nil {
		return nil, fmt.Errorf("failed to recharge card: %w", err)
	}
	return &order, nil
}

// Withdraw withdraws funds from a card
func (c *CardsClient) Withdraw(ctx context.Context, cardID string, req *CardOrderRequest) (*CardOrder, error) {
	var order CardOrder
	path := fmt.Sprintf("/v1/issuing/cards/%s/withdraw", cardID)
	if err := c.client.Post(ctx, path, req, &order); err != nil {
		return nil, fmt.Errorf("failed to withdraw from card: %w", err)
	}
	return &order, nil
}

// GetOrder retrieves a card order by order ID
func (c *CardsClient) GetOrder(ctx context.Context, orderID string) (*CardOrder, error) {
	var order CardOrder
	path := fmt.Sprintf("/v1/issuing/cards/%s/order", orderID)
	if err := c.client.Get(ctx, path, &order); err != nil {
		return nil, fmt.Errorf("failed to get card order: %w", err)
	}
	return &order, nil
}

// Activate activates a physical card
func (c *CardsClient) Activate(ctx context.Context, req *ActivateCardRequest) (*ActivateCardResponse, error) {
	var resp ActivateCardResponse
	if err := c.client.Post(ctx, "/v1/issuing/cards/activate", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to activate card: %w", err)
	}
	return &resp, nil
}

// ResetPIN resets the PIN for a physical card
func (c *CardsClient) ResetPIN(ctx context.Context, req *SetPINRequest) (*SetPINResponse, error) {
	var resp SetPINResponse
	if err := c.client.Post(ctx, "/v1/issuing/cards/pin", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to reset card PIN: %w", err)
	}
	return &resp, nil
}

// Assign assigns a physical card or bulk created virtual card to a cardholder
func (c *CardsClient) Assign(ctx context.Context, req *AssignCardRequest) (*AssignCardResponse, error) {
	var resp AssignCardResponse
	if err := c.client.Post(ctx, "/v1/issuing/cards/assign", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to assign card: %w", err)
	}
	return &resp, nil
}

// BulkCreate creates virtual cards in bulk
func (c *CardsClient) BulkCreate(ctx context.Context, req *BulkCardCreationRequest) (*BulkCardCreationResponse, error) {
	var resp BulkCardCreationResponse
	if err := c.client.Post(ctx, "/v1/issuing/cards/bulk", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to bulk create cards: %w", err)
	}
	return &resp, nil
}

// CreatePANToken creates a one-time PAN token for accessing sensitive card details
// through a secure iframe. The token expires after 60 seconds and can only be used once.
func (c *CardsClient) CreatePANToken(ctx context.Context, cardID string) (*CreatePANTokenResponse, error) {
	var resp CreatePANTokenResponse
	path := fmt.Sprintf("/v1/issuing/cards/%s/token", cardID)
	if err := c.client.Post(ctx, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to create PAN token: %w", err)
	}
	return &resp, nil
}
