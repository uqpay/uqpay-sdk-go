package webhook

// PaymentIntentData represents the payment intent information in acquiring webhook events.
// This is returned in the data field for acquiring.payment_intent.* events.
type PaymentIntentData struct {
	// PaymentIntentID is the unique identifier for the payment intent
	PaymentIntentID string `json:"payment_intent_id"`

	// Amount is the payment amount as a string (e.g., "101")
	Amount string `json:"amount"`

	// Currency is the ISO 4217 currency code (e.g., "USD")
	Currency string `json:"currency"`

	// Description is the payment description
	Description string `json:"description,omitempty"`

	// IntentStatus is the current status of the payment intent
	IntentStatus string `json:"intent_status"`

	// MerchantOrderID is the merchant's order reference
	MerchantOrderID string `json:"merchant_order_id,omitempty"`

	// Metadata contains custom key-value pairs
	Metadata map[string]string `json:"metadata,omitempty"`

	// PaymentMethod contains payment method details (if available)
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty"`

	// CreateTime is the creation timestamp
	CreateTime string `json:"create_time,omitempty"`

	// CompleteTime is the completion timestamp (if completed)
	CompleteTime *string `json:"complete_time,omitempty"`

	// CancelTime is the cancellation timestamp (if canceled)
	CancelTime *string `json:"cancel_time,omitempty"`

	// CancellationReason is the reason for cancellation (if canceled)
	CancellationReason string `json:"cancellation_reason,omitempty"`
}

// PaymentMethod represents the payment method used for a payment
type PaymentMethod struct {
	// Type is the payment method type (e.g., "card", "alipaycn", "grabpay")
	Type string `json:"type,omitempty"`

	// Card contains card details (if payment method is card)
	Card *CardDetails `json:"card,omitempty"`

	// AlipayCN contains Alipay China details (if payment method is alipaycn)
	AlipayCN *AlipayDetails `json:"alipaycn,omitempty"`

	// AlipayHK contains Alipay Hong Kong details (if payment method is alipayhk)
	AlipayHK *AlipayDetails `json:"alipayhk,omitempty"`

	// GrabPay contains GrabPay details (if payment method is grabpay)
	GrabPay *AlipayDetails `json:"grabpay,omitempty"`

	// WeChatPay contains WeChat Pay details (if payment method is wechatpay)
	WeChatPay *AlipayDetails `json:"wechatpay,omitempty"`
}

// AlipayDetails represents alternative payment method details (Alipay, GrabPay, etc.)
type AlipayDetails struct {
	// Flow is the payment flow type (e.g., "qrcode", "app", "wap")
	Flow string `json:"flow,omitempty"`

	// OSType is the operating system type (e.g., "ios", "android")
	OSType string `json:"os_type,omitempty"`
}

// CardDetails represents card payment method details
type CardDetails struct {
	// Brand is the card brand (e.g., "visa", "mastercard")
	Brand string `json:"brand,omitempty"`

	// Last4 is the last 4 digits of the card number
	Last4 string `json:"last4,omitempty"`

	// ExpMonth is the expiration month
	ExpMonth int `json:"exp_month,omitempty"`

	// ExpYear is the expiration year
	ExpYear int `json:"exp_year,omitempty"`

	// Funding is the card funding type (e.g., "credit", "debit")
	Funding string `json:"funding,omitempty"`

	// Country is the card issuing country
	Country string `json:"country,omitempty"`
}

// Payment intent status constants
const (
	// IntentStatusRequiresPaymentMethod indicates payment method is needed
	IntentStatusRequiresPaymentMethod = "REQUIRES_PAYMENT_METHOD"

	// IntentStatusRequiresConfirmation indicates confirmation is needed
	IntentStatusRequiresConfirmation = "REQUIRES_CONFIRMATION"

	// IntentStatusRequiresAction indicates additional action is needed (e.g., 3DS)
	IntentStatusRequiresAction = "REQUIRES_ACTION"

	// IntentStatusProcessing indicates payment is being processed
	IntentStatusProcessing = "PROCESSING"

	// IntentStatusSucceeded indicates payment was successful
	IntentStatusSucceeded = "SUCCEEDED"

	// IntentStatusCanceled indicates payment was canceled
	IntentStatusCanceled = "CANCELED"

	// IntentStatusFailed indicates payment failed
	IntentStatusFailed = "FAILED"
)

// PaymentAttemptData represents the payment attempt information in acquiring webhook events.
// This is returned in the data field for acquiring.payment_attempt.* events.
type PaymentAttemptData struct {
	// PaymentAttemptID is the unique identifier for the payment attempt
	PaymentAttemptID string `json:"payment_attempt_id"`

	// PaymentIntentID is the ID of the parent payment intent
	PaymentIntentID string `json:"payment_intent_id"`

	// Amount is the payment amount as a string (e.g., "0.01")
	Amount string `json:"amount"`

	// Currency is the ISO 4217 currency code (e.g., "USD")
	Currency string `json:"currency"`

	// AttemptStatus is the current status of the payment attempt
	AttemptStatus string `json:"attempt_status"`

	// MerchantOrderID is the merchant's order reference
	MerchantOrderID string `json:"merchant_order_id,omitempty"`

	// PaymentMethod contains payment method details
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty"`

	// CapturedAmount is the amount that has been captured
	CapturedAmount string `json:"captured_amount,omitempty"`

	// RefundedAmount is the amount that has been refunded
	RefundedAmount string `json:"refunded_amount,omitempty"`

	// FailureCode is the failure code if the attempt failed
	FailureCode string `json:"failure_code,omitempty"`

	// CreateTime is the creation timestamp
	CreateTime string `json:"create_time,omitempty"`

	// CompleteTime is the completion timestamp (if completed)
	CompleteTime *string `json:"complete_time,omitempty"`

	// CancelTime is the cancellation timestamp (if canceled)
	CancelTime *string `json:"cancel_time,omitempty"`

	// CancellationReason is the reason for cancellation (if canceled)
	CancellationReason string `json:"cancellation_reason,omitempty"`
}

// Payment attempt status constants
const (
	// AttemptStatusInitiated indicates the attempt has been initiated
	AttemptStatusInitiated = "INITIATED"

	// AttemptStatusPending indicates the attempt is pending
	AttemptStatusPending = "PENDING"

	// AttemptStatusCaptureRequested indicates capture has been requested
	AttemptStatusCaptureRequested = "CAPTURE_REQUESTED"

	// AttemptStatusSucceeded indicates the attempt was successful
	AttemptStatusSucceeded = "SUCCEEDED"

	// AttemptStatusFailed indicates the attempt failed
	AttemptStatusFailed = "FAILED"

	// AttemptStatusCanceled indicates the attempt was canceled
	AttemptStatusCanceled = "CANCELED"
)
