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
	Amount          string            `json:"amount"`                   // Required. Decimal amount to collect, e.g. "10.12"
	Currency        string            `json:"currency"`                 // Required. ISO 4217 three-letter currency code
	MerchantOrderID string            `json:"merchant_order_id"`        // Required. Merchant reference ID, max 128 chars
	Description     string            `json:"description"`              // Required. Descriptor shown to customer, max 32 chars
	ReturnURL       string            `json:"return_url"`               // Required. Redirect URL after payment authentication, max 1024 chars
	PaymentMethod   *PaymentMethod    `json:"payment_method,omitempty"` // Optional. Payment method details
	IPAddress       string            `json:"ip_address,omitempty"`     // Optional. IPv4/IPv6, max 45 chars; required when three_ds_action=enforce_3ds
	PaymentOrders   *PaymentOrders    `json:"payment_orders,omitempty"` // Optional. Purchase order details
	BrowserInfo     *BrowserInfo      `json:"browser_info,omitempty"`   // Optional. Browser data for risk/fraud; required when three_ds_action=enforce_3ds
	Metadata        map[string]string `json:"metadata,omitempty"`       // Optional. Custom key-value pairs, max 512 bytes JSON
}

// ============================================================================
// Payment Method Types
// ============================================================================

// PaymentMethod represents the payment method details
// Only one payment method type should be set based on the Type field
type PaymentMethod struct {
	Type string `json:"type"` // Required. Discriminator: card, card_present, applepay, googlepay, crypto, alipaycn, alipayhk, unionpay, wechatpay, grabpay, paynow, truemoney, tng, gcash, dana, kakaopay, toss, naverpay

	// Card payments (online)
	Card *Card `json:"card,omitempty"` // Optional. Online card payment details

	// Card present payments (POS terminal)
	CardPresent *CardPresent `json:"card_present,omitempty"` // Optional. In-person POS terminal payment

	// Digital wallet token payments
	ApplePay  *ApplePay  `json:"applepay,omitempty"`  // Optional. Apple Pay with network token
	GooglePay *GooglePay `json:"googlepay,omitempty"` // Optional. Google Pay with network token

	// Cryptocurrency payments
	Crypto *Crypto `json:"crypto,omitempty"` // Optional. Crypto payment; currency must be USD

	// E-wallet and QR code payments
	AlipayCN  *WalletPayment `json:"alipaycn,omitempty"`  // Optional. Alipay China
	AlipayHK  *WalletPayment `json:"alipayhk,omitempty"`  // Optional. Alipay Hong Kong
	UnionPay  *WalletPayment `json:"unionpay,omitempty"`  // Optional. UnionPay
	WeChatPay *WeChatPay     `json:"wechatpay,omitempty"` // Optional. WeChat Pay
	GrabPay   *GrabPay       `json:"grabpay,omitempty"`   // Optional. GrabPay
	PayNow    *WalletPayment `json:"paynow,omitempty"`    // Optional. PayNow (Singapore)
	TrueMoney *WalletPayment `json:"truemoney,omitempty"` // Optional. TrueMoney (Thailand)
	TNG       *WalletPayment `json:"tng,omitempty"`       // Optional. Touch 'n Go (Malaysia)
	GCash     *WalletPayment `json:"gcash,omitempty"`     // Optional. GCash (Philippines)
	Dana      *WalletPayment `json:"dana,omitempty"`      // Optional. Dana (Indonesia)
	KakaoPay  *WalletPayment `json:"kakaopay,omitempty"`  // Optional. KakaoPay (South Korea)
	Toss      *WalletPayment `json:"tosspay,omitempty"`   // Optional. Toss (South Korea)
	NaverPay  *WalletPayment `json:"naverpay,omitempty"`  // Optional. NaverPay (South Korea)
}

// ============================================================================
// Card Payment (Online)
// ============================================================================

// Card represents online card payment details
type Card struct {
	CardName          string   `json:"card_name,omitempty"`          // Optional. Cardholder name
	CardNumber        string   `json:"card_number,omitempty"`        // Optional. Full card number
	ExpiryMonth       string   `json:"expiry_month,omitempty"`       // Optional. 2-digit expiry month
	ExpiryYear        string   `json:"expiry_year,omitempty"`        // Optional. 4-digit expiry year
	CVC               string   `json:"cvc,omitempty"`                // Optional. Card verification code
	Network           string   `json:"network,omitempty"`            // Optional. Card network: visa, mastercard, amex, etc.
	Billing           *Billing `json:"billing,omitempty"`            // Optional. Billing information
	AutoCapture       *bool    `json:"auto_capture,omitempty"`       // Optional. Auto-capture after authorization; defaults to true
	AuthorizationType string   `json:"authorization_type,omitempty"` // Optional. Set to "authorization" for manual capture
	ThreeDSAction     string   `json:"three_ds_action,omitempty"`    // Optional. Set to "enforce_3ds" to enforce 3D Secure
}

// Billing represents billing information for a card payment
type Billing struct {
	FirstName   string   `json:"first_name,omitempty"`   // Optional. Cardholder first name
	LastName    string   `json:"last_name,omitempty"`    // Optional. Cardholder last name
	Email       string   `json:"email,omitempty"`        // Optional. Cardholder email
	PhoneNumber string   `json:"phone_number,omitempty"` // Optional. Cardholder phone number
	Address     *Address `json:"address,omitempty"`      // Optional. Billing address
}

// Address represents a billing or shipping address
type Address struct {
	CountryCode string `json:"country_code,omitempty"` // Required. ISO 3166-1 alpha-2 country code
	State       string `json:"state,omitempty"`        // Optional. Required for US/CA, max 100 chars
	City        string `json:"city,omitempty"`         // Required. Max 100 chars
	Street      string `json:"street,omitempty"`       // Required. Max 100 chars
	Postcode    string `json:"postcode,omitempty"`     // Required. Max 10 chars
}

// ============================================================================
// Customer
// ============================================================================

// CustomerRequest represents customer details for a payment
type CustomerRequest struct {
	FirstName   string            `json:"first_name"`             // Required. Customer first name
	LastName    string            `json:"last_name"`              // Required. Customer last name
	Email       string            `json:"email"`                  // Required. Customer email address
	PhoneNumber string            `json:"phone_number,omitempty"` // Optional. Customer phone number
	Description string            `json:"description,omitempty"`  // Optional. Max 255 chars
	Address     *Address          `json:"address,omitempty"`      // Optional. Customer address
	Metadata    map[string]string `json:"metadata,omitempty"`     // Optional. Custom key-value pairs
}

// ============================================================================
// Payment Orders
// ============================================================================

// PaymentOrders represents purchase order information
type PaymentOrders struct {
	Type     string           `json:"type,omitempty"`     // Optional. Industry category, e.g. "physical_goods"
	Products []PaymentProduct `json:"products,omitempty"` // Optional. List of products in the order
}

// PaymentProduct represents a product in a payment order
type PaymentProduct struct {
	Name     string `json:"name"`                // Required. Product name, max 255 chars
	Price    string `json:"price"`               // Required. Price per unit, decimal format
	Quantity int    `json:"quantity"`            // Required. Quantity of the product
	ImageURL string `json:"image_url,omitempty"` // Optional. Product thumbnail URL
}

// ============================================================================
// Browser Info (3DS / Risk)
// ============================================================================

// BrowserInfo represents browser information for risk and fraud prevention
// Required when three_ds_action=enforce_3ds
type BrowserInfo struct {
	AcceptHeader     string         `json:"accept_header,omitempty"`      // Optional. HTTP Accept header value
	Browser          *BrowserDetail `json:"browser,omitempty"`            // Optional. Browser-specific details
	DeviceID         string         `json:"device_id,omitempty"`          // Optional. Unique device identifier
	Language         string         `json:"language,omitempty"`           // Optional. ISO language code, e.g. "en-US"
	Location         *GeoLocation   `json:"location,omitempty"`           // Optional. Device geolocation
	Mobile           *MobileInfo    `json:"mobile,omitempty"`             // Optional. Required for mobile transactions
	ScreenColorDepth int            `json:"screen_color_depth,omitempty"` // Optional. Color depth in bits, 1-48
	ScreenHeight     int            `json:"screen_height,omitempty"`      // Optional. Screen height in pixels, 1-9999
	ScreenWidth      int            `json:"screen_width,omitempty"`       // Optional. Screen width in pixels, 1-9999
	Timezone         string         `json:"timezone,omitempty"`           // Optional. UTC offset, e.g. "-2" or "8"
}

// BrowserDetail represents browser-specific information
type BrowserDetail struct {
	JavaEnabled       bool   `json:"java_enabled,omitempty"`       // Optional. Whether Java is enabled in the browser
	JavascriptEnabled bool   `json:"javascript_enabled,omitempty"` // Optional. Whether JavaScript is enabled
	UserAgent         string `json:"user_agent,omitempty"`         // Optional. Browser user agent string
}

// GeoLocation represents device geolocation
type GeoLocation struct {
	Lat string `json:"lat,omitempty"` // Latitude coordinate
	Lon string `json:"lon,omitempty"` // Longitude coordinate
}

// MobileInfo represents mobile device information
type MobileInfo struct {
	DeviceModel string `json:"device_model,omitempty"` // e.g., "Apple IPHONE 7"
	OSType      string `json:"os_type,omitempty"`      // IOS, ANDROID
	OSVersion   string `json:"os_version,omitempty"`   // e.g., "IOS 14.5"
}

// ============================================================================
// Card Present Payment (POS Terminal)
// ============================================================================

// CardPresent represents card present (POS terminal) payment details
type CardPresent struct {
	CardNumber                   string        `json:"card_number"`                              // Required. Full card number
	ExpiryMonth                  string        `json:"expiry_month"`                             // Required. 2-digit expiry month
	ExpiryYear                   string        `json:"expiry_year"`                              // Required. 4-digit expiry year
	CardholderVerificationMethod string        `json:"cardholder_verification_method,omitempty"` // Optional. online_pin, manual_signature, or skipped
	EncryptedPIN                 string        `json:"encrypted_pin,omitempty"`                  // Optional. Encrypted PIN block
	PANEntryMode                 string        `json:"pan_entry_mode"`                           // Required. manual_entry, chip, magstripe, contactless_chip, or contactless_magstripe
	Fallback                     bool          `json:"fallback,omitempty"`                       // Optional. Whether this is a fallback transaction
	FallbackReason               string        `json:"fallback_reason,omitempty"`                // Optional. e.g. "chip_read_failure"
	EMVTags                      string        `json:"emv_tags,omitempty"`                       // Optional. EMV tag data
	Track1                       string        `json:"track1,omitempty"`                         // Optional. Magnetic stripe track 1 data
	Track2                       string        `json:"track2,omitempty"`                         // Optional. Magnetic stripe track 2 data
	TerminalInfo                 *TerminalInfo `json:"terminal_info,omitempty"`                  // Optional. POS terminal details
	SystemTraceAuditNumber       string        `json:"system_trace_audit_number,omitempty"`      // Optional. 6-digit trace number
}

// TerminalInfo represents POS terminal information
type TerminalInfo struct {
	TerminalID        string `json:"terminal_id,omitempty"`         // Optional. Up to 8 alphanumeric chars
	MobileDevice      bool   `json:"mobile_device,omitempty"`       // Optional. Whether this is a mobile POS device
	UseEmbeddedReader bool   `json:"use_embedded_reader,omitempty"` // Optional. Whether reader is embedded in mobile POS
}

// ============================================================================
// Wallet/E-Wallet Payments (Common Structure)
// ============================================================================

// WalletPayment represents common fields for wallet/e-wallet payments
// Used for: alipaycn, alipayhk, unionpay, paynow, truemoney, tng, gcash, dana, kakaopay, toss, naverpay
type WalletPayment struct {
	Flow        string `json:"flow"`                   // Required. Payment flow: qrcode, securepay, mobile_app, mobile_web, etc.
	OSType      string `json:"os_type,omitempty"`      // Optional. ios or android; required for mobile_web/mobile_app flows
	IsPresent   bool   `json:"is_present,omitempty"`   // Optional. Whether customer is physically present
	PaymentCode string `json:"payment_code,omitempty"` // Optional. Customer's payment code from wallet app
}

// ============================================================================
// WeChat Pay (Extended Wallet)
// ============================================================================

// WeChatPay represents WeChat Pay payment details
// Has additional fields compared to standard wallet payments
type WeChatPay struct {
	Flow        string `json:"flow"`                   // Required. qrcode, mini_program, mobile_app, mobile_web, or official_account
	OSType      string `json:"os_type,omitempty"`      // Optional. ios or android; required for mobile_web/mobile_app flows
	IsPresent   bool   `json:"is_present,omitempty"`   // Optional. Whether customer is physically present
	PaymentCode string `json:"payment_code,omitempty"` // Optional. Customer's payment code from wallet app
	OpenID      string `json:"open_id,omitempty"`      // Optional. Required for mini_program, mobile_app, official_account flows
}

// ============================================================================
// GrabPay (Extended Wallet)
// ============================================================================

// GrabPay represents GrabPay payment details
// Has additional fields compared to standard wallet payments
type GrabPay struct {
	Flow        string `json:"flow"`                   // Required. Payment flow: qrcode
	OSType      string `json:"os_type,omitempty"`      // Optional. ios or android; required for mobile_web
	IsPresent   bool   `json:"is_present,omitempty"`   // Optional. Whether customer is physically present
	PaymentCode string `json:"payment_code,omitempty"` // Optional. Customer's payment code from wallet app
	ShopperName string `json:"shopper_name,omitempty"` // Optional. Name of the shopper
}

// ============================================================================
// Apple Pay
// ============================================================================

// ApplePay represents Apple Pay payment details
type ApplePay struct {
	Flow           string          `json:"flow"`                      // Required. redirect, mobile_web, mobile_app, or contactless
	OSType         string          `json:"os_type,omitempty"`         // Optional. ios; required for mobile_web/mobile_app flows
	IsPresent      bool            `json:"is_present,omitempty"`      // Optional. Whether customer is physically present
	Network        string          `json:"network"`                   // Required. Card network: visa, mastercard, amex, discover, jcb
	CardType       string          `json:"card_type,omitempty"`       // Optional. debit or credit
	TokenType      string          `json:"token_type"`                // Required. decrypted or encrypted
	AuthMethod     string          `json:"auth_method"`               // Required. cryptogram_3ds or pan_only
	NetworkToken   *NetworkToken   `json:"network_token"`             // Required. DPAN token details
	BillingContact *BillingContact `json:"billing_contact,omitempty"` // Optional. Billing contact information
}

// ============================================================================
// Google Pay
// ============================================================================

// GooglePay represents Google Pay payment details
type GooglePay struct {
	Flow           string          `json:"flow"`                      // Required. redirect, mobile_web, mobile_app, or contactless
	OSType         string          `json:"os_type,omitempty"`         // Optional. ios or android; required for mobile_web/mobile_app flows
	IsPresent      bool            `json:"is_present,omitempty"`      // Optional. Whether customer is physically present
	Network        string          `json:"network"`                   // Required. Card network: visa, mastercard, amex, discover, jcb
	CardType       string          `json:"card_type,omitempty"`       // Optional. debit or credit
	TokenType      string          `json:"token_type"`                // Required. decrypted or encrypted
	AuthMethod     string          `json:"auth_method"`               // Required. cryptogram_3ds or pan_only
	NetworkToken   *NetworkToken   `json:"network_token"`             // Required. DPAN token details
	BillingAddress *BillingContact `json:"billing_address,omitempty"` // Optional. Billing address information
}

// NetworkToken represents the DPAN token details for Apple Pay and Google Pay
type NetworkToken struct {
	Number      string `json:"number"`               // Required. DPAN number, 12-52 chars
	ExpiryMonth string `json:"expiry_month"`         // Required. 2-digit expiry month
	ExpiryYear  string `json:"expiry_year"`          // Required. 4-digit expiry year
	Cryptogram  string `json:"cryptogram,omitempty"` // Optional. Base64 encoded; required for cryptogram_3ds auth method
	ECI         string `json:"eci,omitempty"`        // Optional. ECI value; required for cryptogram_3ds auth method
}

// BillingContact represents billing contact information for Apple Pay and Google Pay
type BillingContact struct {
	FirstName string   `json:"first_name,omitempty"` // Optional. Contact first name
	LastName  string   `json:"last_name,omitempty"`  // Optional. Contact last name
	Email     string   `json:"email,omitempty"`      // Optional. Contact email
	Phone     string   `json:"phone,omitempty"`      // Optional. Contact phone number
	Address   *Address `json:"address,omitempty"`    // Optional. Contact address
}

// ============================================================================
// Cryptocurrency
// ============================================================================

// Crypto represents cryptocurrency payment details
// Note: Currency must be USD for crypto payments
type Crypto struct {
	Flow      string `json:"flow"`              // Required. redirect or qrcode
	Network   string `json:"network,omitempty"` // Optional. ETH or TRON; required when flow is qrcode
	IsPresent bool   `json:"is_present"`        // Required. Must be false for crypto payments
}

// UpdatePaymentIntentRequest represents a payment intent update request
type UpdatePaymentIntentRequest struct {
	Amount          string            `json:"amount,omitempty"`            // Optional. Updated payment amount, decimal format
	Currency        string            `json:"currency,omitempty"`          // Optional. ISO 4217 three-letter currency code
	Customer        *CustomerRequest  `json:"customer,omitempty"`          // Optional. Customer details; omit when customer_id is specified
	CustomerID      string            `json:"customer_id,omitempty"`       // Optional. Unique customer ID for recurring payments
	PaymentOrders   *PaymentOrders    `json:"payment_orders,omitempty"`    // Optional. Purchase order details
	MerchantOrderID string            `json:"merchant_order_id,omitempty"` // Optional. Merchant reference ID
	Description     string            `json:"description,omitempty"`       // Optional. Descriptor shown to customer, max 32 chars
	Metadata        map[string]string `json:"metadata,omitempty"`          // Optional. Custom key-value pairs, max 512 bytes JSON
	ReturnURL       string            `json:"return_url,omitempty"`        // Optional. Redirect URL after authentication, max 1024 chars
}

// ConfirmPaymentIntentRequest represents a payment intent confirmation request
type ConfirmPaymentIntentRequest struct {
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty"` // Optional. Payment method to confirm with
	IPAddress     string         `json:"ip_address,omitempty"`     // Optional. IPv4/IPv6; required when three_ds_action=enforce_3ds
	BrowserInfo   *BrowserInfo   `json:"browser_info,omitempty"`   // Optional. Browser data for fraud prevention; required for 3DS
	ReturnURL     string         `json:"return_url,omitempty"`     // Optional. Redirect URL after payment authentication
}

// CapturePaymentIntentRequest represents a payment intent capture request
type CapturePaymentIntentRequest struct {
	AmountToCapture float64 `json:"amount_to_capture,omitempty"` // Optional. Amount to capture; defaults to original amount if omitted
}

// CancelPaymentIntentRequest represents a payment intent cancellation request
type CancelPaymentIntentRequest struct {
	CancellationReason string `json:"cancellation_reason,omitempty"` // Optional. duplicate, fraudulent, requested_by_customer, or abandoned
}

// ListPaymentIntentsRequest represents a payment intents list request
type ListPaymentIntentsRequest struct {
	PageSize            int    `json:"page_size"`             // Required. Items per page, min 1, max 100
	PageNumber          int    `json:"page_number"`           // Required. Page number, 1-based
	PaymentIntentStatus string `json:"payment_intent_status"` // Optional. Filter: REQUIRES_PAYMENT_METHOD, REQUIRES_CUSTOMER_ACTION, REQUIRES_CAPTURE, PENDING, SUCCEEDED, CANCELLED, FAILED
	StartTime           string `json:"start_time"`            // Optional. Exclusive start time filter, ISO 8601 format
	EndTime             string `json:"end_time"`              // Optional. Exclusive end time filter, ISO 8601 format
}

// ============================================================================
// Response Structures
// ============================================================================

// PaymentIntent represents a payment intent response
type PaymentIntent struct {
	PaymentIntentID             string                 `json:"payment_intent_id"`                        // Unique payment intent identifier
	Amount                      string                 `json:"amount"`                                   // Order amount charged, decimal format
	Currency                    string                 `json:"currency"`                                 // ISO 4217 three-letter currency code
	IntentStatus                string                 `json:"intent_status"`                            // REQUIRES_PAYMENT_METHOD, REQUIRES_CUSTOMER_ACTION, REQUIRES_CAPTURE, PENDING, SUCCEEDED, CANCELLED, or FAILED
	MerchantOrderID             string                 `json:"merchant_order_id,omitempty"`              // Merchant reference ID
	Description                 string                 `json:"description,omitempty"`                    // Payment descriptor shown to customer
	ReturnURL                   string                 `json:"return_url,omitempty"`                     // Redirect URL after payment authentication
	Metadata                    map[string]string      `json:"metadata,omitempty"`                       // Custom key-value pairs, max 512 bytes JSON
	AvailablePaymentMethodTypes []string               `json:"available_payment_method_types,omitempty"` // List of supported payment method types
	CapturedAmount              string                 `json:"captured_amount,omitempty"`                // Amount already captured
	Customer                    *CustomerRequest       `json:"customer,omitempty"`                       // Customer details
	ClientSecret                string                 `json:"client_secret,omitempty"`                  // Client secret, valid for 60 minutes
	CancellationReason          string                 `json:"cancellation_reason,omitempty"`            // Reason for cancellation if cancelled
	LatestPaymentAttempt        map[string]interface{} `json:"latest_payment_attempt,omitempty"`         // Most recent payment attempt details
	NextAction                  map[string]interface{} `json:"next_action,omitempty"`                    // Required customer action: redirect_to_url, display_qr_code, display_bank_details, or redirect_iframe
	CreateTime                  string                 `json:"create_time,omitempty"`                    // Creation timestamp, ISO 8601
	UpdateTime                  string                 `json:"update_time,omitempty"`                    // Last update timestamp, ISO 8601
	CancelTime                  string                 `json:"cancel_time,omitempty"`                    // Cancellation timestamp, ISO 8601
	CompleteTime                string                 `json:"complete_time,omitempty"`                  // Completion timestamp, ISO 8601
}

// ListPaymentIntentsResponse represents a paginated list of payment intents
type ListPaymentIntentsResponse struct {
	TotalPages int             `json:"total_pages"` // Total number of pages
	TotalItems int             `json:"total_items"` // Total number of items
	Data       []PaymentIntent `json:"data"`        // Array of payment intent objects
}

// ============================================================================
// API Methods
// ============================================================================

// Create creates a new payment intent
func (c *PaymentIntentsClient) Create(ctx context.Context, req *CreatePaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	opts := &common.RequestOptions{
		ClientID: c.client.Config.ClientID,
	}
	if err := c.client.PostWithOptions(ctx, "/v2/payment_intents/create", req, &resp, opts); err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}
	return &resp, nil
}

// Get retrieves a specific payment intent by ID
func (c *PaymentIntentsClient) Get(ctx context.Context, paymentIntentID string) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s", paymentIntentID)
	opts := &common.RequestOptions{
		ClientID: c.client.Config.ClientID,
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opts); err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}
	return &resp, nil
}

// Update updates properties on a payment intent without confirming
// Note: Updating payment_method requires subsequent confirmation
func (c *PaymentIntentsClient) Update(ctx context.Context, paymentIntentID string, req *UpdatePaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s", paymentIntentID)
	opts := &common.RequestOptions{
		ClientID: c.client.Config.ClientID,
	}
	if err := c.client.PostWithOptions(ctx, path, req, &resp, opts); err != nil {
		return nil, fmt.Errorf("failed to update payment intent: %w", err)
	}
	return &resp, nil
}

// Confirm confirms a payment intent for payment authorization
func (c *PaymentIntentsClient) Confirm(ctx context.Context, paymentIntentID string, req *ConfirmPaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s/confirm", paymentIntentID)
	opts := &common.RequestOptions{
		ClientID: c.client.Config.ClientID,
	}
	if err := c.client.PostWithOptions(ctx, path, req, &resp, opts); err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}
	return &resp, nil
}

// Capture captures the funds of an uncaptured payment intent
// The payment intent must have status "requires_capture"
func (c *PaymentIntentsClient) Capture(ctx context.Context, paymentIntentID string, req *CapturePaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s/capture", paymentIntentID)
	opts := &common.RequestOptions{
		ClientID: c.client.Config.ClientID,
	}
	if err := c.client.PostWithOptions(ctx, path, req, &resp, opts); err != nil {
		return nil, fmt.Errorf("failed to capture payment intent: %w", err)
	}
	return &resp, nil
}

// Cancel cancels a payment intent and prevents further payment attempts
func (c *PaymentIntentsClient) Cancel(ctx context.Context, paymentIntentID string, req *CancelPaymentIntentRequest) (*PaymentIntent, error) {
	var resp PaymentIntent
	path := fmt.Sprintf("/v2/payment_intents/%s/cancel", paymentIntentID)
	opts := &common.RequestOptions{
		ClientID: c.client.Config.ClientID,
	}
	if err := c.client.PostWithOptions(ctx, path, req, &resp, opts); err != nil {
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
	if req.PaymentIntentStatus != "" {
		path += fmt.Sprintf("%spayment_intent_status=%s", separator, req.PaymentIntentStatus)
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

	opts := &common.RequestOptions{
		ClientID: c.client.Config.ClientID,
	}
	if err := c.client.GetWithOptions(ctx, path, &resp, opts); err != nil {
		return nil, fmt.Errorf("failed to list payment intents: %w", err)
	}
	return &resp, nil
}
