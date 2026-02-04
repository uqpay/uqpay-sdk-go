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
	CardLimit        *float64          `json:"card_limit,omitempty"`
	CardCurrency     string            `json:"card_currency"`
	CardholderID     string            `json:"cardholder_id"`
	CardProductID    string            `json:"card_product_id"`
	SpendingControls []SpendingControl `json:"spending_controls,omitempty"`
	RiskControls     *RiskControls     `json:"risk_controls,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// SpendingControl represents spending control rules for a card
type SpendingControl struct {
	Amount   float64 `json:"amount"`
	Interval string  `json:"interval"` // PER_TRANSACTION
}

// RiskControls represents user-customized risk control settings
type RiskControls struct {
	Allow3DSTransactions *string  `json:"allow_3ds_transactions,omitempty"` // Y or N
	AllowedMCC           []string `json:"allowed_mcc,omitempty"`
	BlockedMCC           []string `json:"blocked_mcc,omitempty"`
}

// CardUpdateRequest represents a card update request
type CardUpdateRequest struct {
	CardLimit           *float64          `json:"card_limit,omitempty"`
	NoPINPaymentAmount  *float64          `json:"no_pin_payment_amount,omitempty"`
	SpendingControls    []SpendingControl `json:"spending_controls,omitempty"`
	RiskControls        *RiskControls     `json:"risk_controls,omitempty"`
	Metadata            map[string]string `json:"metadata,omitempty"`
}

// UpdateCardStatusRequest represents a card status update request
type UpdateCardStatusRequest struct {
	CardStatus   string  `json:"card_status"` // ACTIVE, FROZEN, CANCELLED
	UpdateReason *string `json:"update_reason,omitempty"`
}

// CardOrderRequest represents a card recharge/withdraw request
type CardOrderRequest struct {
	Amount float64 `json:"amount"`
}

// ActivateCardRequest represents a card activation request
type ActivateCardRequest struct {
	CardID             string   `json:"card_id"`
	ActivationCode     string   `json:"activation_code"`
	PIN                string   `json:"pin"`
	NoPINPaymentAmount *float64 `json:"no_pin_payment_amount,omitempty"`
}

// SetPINRequest represents a card PIN reset request
type SetPINRequest struct {
	CardID string `json:"card_id"`
	PIN    string `json:"pin"`
}

// AssignCardRequest represents a card assignment request
type AssignCardRequest struct {
	CardholderID string `json:"cardholder_id"`
	CardNumber   string `json:"card_number"`
	CardCurrency string `json:"card_currency"`
	CardMode     string `json:"card_mode"` // SINGLE or SHARE
}

// BulkCardCreationRequest represents a bulk card creation request
type BulkCardCreationRequest struct {
	CardBIN string `json:"card_bin"`
	Numbers int    `json:"numbers"` // 1-5000
}

// ListCardsRequest represents a card list request
type ListCardsRequest struct {
	PageSize     int     `json:"page_size"`
	PageNumber   int     `json:"page_number"`
	CardNumber   *string `json:"card_number,omitempty"`
	CardStatus   *string `json:"card_status,omitempty"`
	CardholderID *string `json:"cardholder_id,omitempty"`
}

// ============================================================================
// Response Structures
// ============================================================================

// CardCreationResponse represents the response after creating a card
type CardCreationResponse struct {
	CardID      string `json:"card_id"`
	CardOrderID string `json:"card_order_id"`
	CreateTime  string `json:"create_time"`
	CardStatus  string `json:"card_status"`
	OrderStatus string `json:"order_status"`
}

// CardUpdatedResponse represents the response after updating a card
type CardUpdatedResponse struct {
	CardID      string `json:"card_id"`
	CardOrderID string `json:"card_order_id"`
	CardStatus  string `json:"card_status"`
	OrderStatus string `json:"order_status"`
}

// CardStatusResponse represents the response after updating card status
type CardStatusResponse struct {
	CardID       string  `json:"card_id"`
	CardOrderID  string  `json:"card_order_id"`
	OrderStatus  string  `json:"order_status"`
	UpdateReason *string `json:"update_reason,omitempty"`
}

// RetrieveCardResponse represents detailed card information
type RetrieveCardResponse struct {
	CardID              string                  `json:"card_id"`
	CardBIN             string                  `json:"card_bin"`
	CardScheme          string                  `json:"card_scheme"`
	CardCurrency        string                  `json:"card_currency"`
	CardNumber          string                  `json:"card_number"`
	FormFactor          string                  `json:"form_factor"`
	ModeType            string                  `json:"mode_type"`
	CardProductID       string                  `json:"card_product_id"`
	CardLimit           float64                 `json:"card_limit"`
	AvailableBalance    string                  `json:"available_balance"`
	Cardholder          CardholderInfo          `json:"cardholder"`
	SpendingControls    []SpendingControl       `json:"spending_controls,omitempty"`
	NoPINPaymentAmount  string                  `json:"no_pin_payment_amount"`
	RiskControls        *RiskControls           `json:"risk_controls,omitempty"`
	Metadata            map[string]string       `json:"metadata,omitempty"`
	CardStatus          string                  `json:"card_status"`
	UpdateReason        *string                 `json:"update_reason,omitempty"`
	ConsumedAmount      *string                 `json:"consumed_amount,omitempty"`
}

// CardholderInfo represents cardholder information in card response
type CardholderInfo struct {
	CardholderID     string  `json:"cardholder_id"`
	Email            string  `json:"email"`
	NumberOfCards    int     `json:"number_of_cards"`
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	CreateTime       string  `json:"create_time"`
	CardholderStatus string  `json:"cardholder_status"`
	DateOfBirth      *string `json:"date_of_birth,omitempty"`
	CountryCode      *string `json:"country_code,omitempty"`
	PhoneNumber      *string `json:"phone_number,omitempty"`
}

// SecureCardInfo represents secure card information
type SecureCardInfo struct {
	CVV        string `json:"cvv"`
	ExpireDate string `json:"expire_date"`
	CardNumber string `json:"card_number"`
}

// CardOrder represents a card order
type CardOrder struct {
	CardID       string  `json:"card_id"`
	CardOrderID  string  `json:"card_order_id"`
	OrderType    string  `json:"order_type"`
	Amount       float64 `json:"amount"`
	CardCurrency string  `json:"card_currency"`
	CreateTime   string  `json:"create_time"`
	UpdateTime   string  `json:"update_time"`
	CompleteTime string  `json:"complete_time"`
	OrderStatus  string  `json:"order_status"`
}

// ActivateCardResponse represents the response after activating a card
type ActivateCardResponse struct {
	RequestStatus string `json:"request_status"`
}

// SetPINResponse represents the response after resetting PIN
type SetPINResponse struct {
	RequestStatus string `json:"request_status"`
}

// AssignCardResponse represents the response after assigning a card
type AssignCardResponse struct {
	CardID      string `json:"card_id"`
	CardOrderID string `json:"card_order_id"`
	CreateTime  string `json:"create_time"`
	CardStatus  string `json:"card_status"`
	OrderStatus string `json:"order_status"`
}

// BulkCardCreationResponse represents the response after bulk card creation
type BulkCardCreationResponse struct {
	ReportID   string  `json:"report_id"`
	ExpireDate *string `json:"expire_date,omitempty"`
}

// ListCardsResponse represents a card list response
type ListCardsResponse struct {
	TotalPages int                    `json:"total_pages"`
	TotalItems int                    `json:"total_items"`
	Data       []RetrieveCardResponse `json:"data"`
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
