package webhook

// CardData represents the card information in issuing webhook events.
// This is returned in the data field for card.create.* and card.update.* events.
type CardData struct {
	// CardID is the unique identifier for the card
	CardID string `json:"card_id"`

	// CardProductID is the ID of the card product used
	CardProductID string `json:"card_product_id,omitempty"`

	// CardOrderID is the ID of the card order (used in update events)
	CardOrderID string `json:"card_order_id,omitempty"`

	// CardNumber is the masked card number (e.g., "49372418****4306")
	CardNumber string `json:"card_number"`

	// CardBin is the Bank Identification Number (first 6-8 digits)
	CardBin string `json:"card_bin,omitempty"`

	// CardScheme is the card network (e.g., "VISA", "MASTERCARD")
	CardScheme string `json:"card_scheme"`

	// CardStatus is the current status of the card
	CardStatus string `json:"card_status"`

	// CardCurrency is the ISO 4217 currency code for the card
	CardCurrency string `json:"card_currency,omitempty"`

	// CardLimit is the card limit as a string
	CardLimit string `json:"card_limit,omitempty"`

	// CardAvailableBalance is the available balance (used in card.create events)
	CardAvailableBalance string `json:"card_available_balance,omitempty"`

	// AvailableBalance is the available balance (used in card.update events)
	AvailableBalance string `json:"available_balance,omitempty"`

	// FormFactor indicates if the card is VIRTUAL or PHYSICAL
	FormFactor string `json:"form_factor"`

	// ModeType is the card mode type (e.g., "SINGLE", "MULTI")
	ModeType string `json:"mode_type"`

	// OrderStatus is the status of the card order (e.g., "success")
	OrderStatus string `json:"order_status,omitempty"`

	// NoPinPaymentAmount contains no-pin payment amount info
	NoPinPaymentAmount string `json:"no_pin_payment_amount,omitempty"`

	// Cardholder contains cardholder information
	Cardholder *Cardholder `json:"cardholder,omitempty"`

	// SpendingLimits contains the spending limits (used in card.create events)
	SpendingLimits []SpendingLimit `json:"spending_limits,omitempty"`

	// SpendingControls contains the spending controls (used in card.update events)
	SpendingControls []SpendingLimit `json:"spending_controls,omitempty"`

	// RiskControl contains risk control settings
	RiskControl *RiskControl `json:"risk_control,omitempty"`

	// Metadata contains custom key-value pairs
	Metadata map[string]string `json:"metadata,omitempty"`
}

// GetAvailableBalance returns the available balance from either field
func (c *CardData) GetAvailableBalance() string {
	if c.AvailableBalance != "" {
		return c.AvailableBalance
	}
	return c.CardAvailableBalance
}

// GetSpendingLimits returns spending limits from either field
func (c *CardData) GetSpendingLimits() []SpendingLimit {
	if len(c.SpendingControls) > 0 {
		return c.SpendingControls
	}
	return c.SpendingLimits
}

// RiskControl represents risk control settings for a card
type RiskControl struct {
	// Allow3DSTransactions indicates if 3DS transactions are allowed ("Y" or "N")
	Allow3DSTransactions string `json:"allow_3ds_transactions,omitempty"`

	// AllowOnlineTransactions indicates if online transactions are allowed
	AllowOnlineTransactions string `json:"allow_online_transactions,omitempty"`

	// AllowATMTransactions indicates if ATM transactions are allowed
	AllowATMTransactions string `json:"allow_atm_transactions,omitempty"`

	// AllowContactlessTransactions indicates if contactless transactions are allowed
	AllowContactlessTransactions string `json:"allow_contactless_transactions,omitempty"`

	// AllowInternationalTransactions indicates if international transactions are allowed
	AllowInternationalTransactions string `json:"allow_international_transactions,omitempty"`
}

// Cardholder represents the cardholder information
type Cardholder struct {
	// CardholderID is the unique identifier for the cardholder
	CardholderID string `json:"cardholder_id"`

	// CardholderStatus is the status of the cardholder (e.g., "SUCCESS", "PENDING")
	CardholderStatus string `json:"cardholder_status"`

	// FirstName is the cardholder's first name
	FirstName string `json:"first_name"`

	// LastName is the cardholder's last name
	LastName string `json:"last_name"`

	// Email is the cardholder's email address
	Email string `json:"email,omitempty"`

	// PhoneNumber is the cardholder's phone number
	PhoneNumber string `json:"phone_number,omitempty"`

	// CountryCode is the cardholder's country code (ISO 3166-1 alpha-2)
	CountryCode string `json:"country_code,omitempty"`

	// DateOfBirth is the cardholder's date of birth (YYYY-MM-DD format)
	DateOfBirth string `json:"date_of_birth,omitempty"`

	// NumberOfCards is the number of cards associated with this cardholder
	NumberOfCards int `json:"number_of_cards,omitempty"`

	// CreateTime is the creation timestamp
	CreateTime string `json:"create_time,omitempty"`
}

// SpendingLimit represents a spending limit configuration
type SpendingLimit struct {
	// Amount is the limit amount as a string
	Amount string `json:"amount"`

	// Interval is the time interval for the limit (e.g., "PER_TRANSACTION", "DAILY", "MONTHLY")
	Interval string `json:"interval"`
}

// CardStatusUpdateData represents the card status update information in issuing webhook events.
// This is returned in the data field for card.status.update.* events.
type CardStatusUpdateData struct {
	// CardID is the unique identifier for the card
	CardID string `json:"card_id"`

	// CardNumber is the masked card number (e.g., "49372418****4306")
	CardNumber string `json:"card_number"`

	// CardStatus is the new status of the card after the update
	CardStatus string `json:"card_status"`

	// UpdateReason is the reason for the status update
	UpdateReason string `json:"update_reason,omitempty"`

	// UpdateTime is the timestamp when the status was updated
	UpdateTime string `json:"update_time,omitempty"`
}

// CardActivationCodeData represents the card activation code information in issuing webhook events.
// This is returned in the data field for card.activation.code events.
type CardActivationCodeData struct {
	// CardID is the unique identifier for the card
	CardID string `json:"card_id"`

	// CardNumber is the masked card number (e.g., "12345678****0000")
	CardNumber string `json:"card_number"`

	// ActivationCode is the activation code for the card
	ActivationCode string `json:"activation_code"`
}

// CardRechargeData represents the card recharge/top-up information in issuing webhook events.
// This is returned in the data field for card.recharge.* events.
type CardRechargeData struct {
	// CardID is the unique identifier for the card
	CardID string `json:"card_id"`

	// Amount is the recharge amount as a string
	Amount string `json:"amount"`

	// CardCurrency is the ISO 4217 currency code for the card
	CardCurrency string `json:"card_currency"`

	// CardAvailableBalance is the available balance after the recharge
	CardAvailableBalance string `json:"card_available_balance"`

	// CardStatus is the current status of the card
	CardStatus string `json:"card_status"`

	// OrderStatus is the status of the recharge order (e.g., "SUCCESS", "FAILED")
	OrderStatus string `json:"order_status"`

	// CompleteTime is the timestamp when the recharge was completed
	CompleteTime string `json:"complete_time,omitempty"`

	// UpdateTime is the timestamp when the record was last updated
	UpdateTime string `json:"update_time,omitempty"`
}

// CardTransactionData represents card transaction information in issuing webhook events.
// This is returned in the data field for issuing.fee.card and other transaction events.
type CardTransactionData struct {
	// CardID is the unique identifier for the card
	CardID string `json:"card_id"`

	// CardNumber is the masked card number (e.g., "49372418****4306")
	CardNumber string `json:"card_number"`

	// CardholderID is the unique identifier for the cardholder
	CardholderID string `json:"cardholder_id"`

	// CardAvailableBalance is the available balance on the card after the transaction
	CardAvailableBalance string `json:"card_available_balance,omitempty"`

	// TransactionAmount is the transaction amount as a string
	TransactionAmount string `json:"transaction_amount"`

	// TransactionCurrency is the ISO 4217 currency code for the transaction
	TransactionCurrency string `json:"transaction_currency"`

	// BillingAmount is the billing amount as a string
	BillingAmount string `json:"billing_amount"`

	// BillingCurrency is the ISO 4217 currency code for billing
	BillingCurrency string `json:"billing_currency"`

	// TransactionStatus is the status of the transaction (e.g., "APPROVED", "DECLINED")
	TransactionStatus string `json:"transaction_status"`

	// TransactionType is the type of transaction (e.g., "FEE", "PURCHASE", "REFUND")
	TransactionType string `json:"transaction_type"`

	// TransactionTime is the timestamp when the transaction occurred
	TransactionTime string `json:"transaction_time,omitempty"`

	// PostedTime is the timestamp when the transaction was posted
	PostedTime string `json:"posted_time,omitempty"`

	// ReferenceID is the unique reference identifier for the transaction
	ReferenceID string `json:"reference_id,omitempty"`

	// ShortReferenceID is a short reference identifier for the transaction
	ShortReferenceID string `json:"short_reference_id,omitempty"`

	// Remark contains additional notes or description about the transaction
	Remark string `json:"remark,omitempty"`
}

// Card status constants
const (
	// CardStatusActive indicates the card is active and can be used
	CardStatusActive = "ACTIVE"

	// CardStatusInactive indicates the card is inactive
	CardStatusInactive = "INACTIVE"

	// CardStatusSuspended indicates the card is suspended
	CardStatusSuspended = "SUSPENDED"

	// CardStatusBlocked indicates the card is blocked
	CardStatusBlocked = "BLOCKED"

	// CardStatusFrozen indicates the card is frozen
	CardStatusFrozen = "FROZEN"

	// CardStatusPreCancel indicates the card is pending cancellation
	CardStatusPreCancel = "PRE_CANCEL"

	// CardStatusClosed indicates the card is closed/terminated
	CardStatusClosed = "CLOSED"

	// CardStatusPending indicates the card is pending activation
	CardStatusPending = "PENDING"
)

// Card scheme constants
const (
	// CardSchemeVisa indicates a Visa card
	CardSchemeVisa = "VISA"

	// CardSchemeMastercard indicates a Mastercard
	CardSchemeMastercard = "MASTERCARD"

	// CardSchemeUnionPay indicates a UnionPay card
	CardSchemeUnionPay = "UNIONPAY"
)

// Form factor constants
const (
	// FormFactorVirtual indicates a virtual card
	FormFactorVirtual = "VIRTUAL"

	// FormFactorPhysical indicates a physical card
	FormFactorPhysical = "PHYSICAL"
)

// Mode type constants
const (
	// ModeTypeSingle indicates a single-use card
	ModeTypeSingle = "SINGLE"

	// ModeTypeMulti indicates a multi-use card
	ModeTypeMulti = "MULTI"
)

// Spending limit interval constants
const (
	// SpendingIntervalPerTransaction limits spending per transaction
	SpendingIntervalPerTransaction = "PER_TRANSACTION"

	// SpendingIntervalDaily limits spending per day
	SpendingIntervalDaily = "DAILY"

	// SpendingIntervalWeekly limits spending per week
	SpendingIntervalWeekly = "WEEKLY"

	// SpendingIntervalMonthly limits spending per month
	SpendingIntervalMonthly = "MONTHLY"

	// SpendingIntervalYearly limits spending per year
	SpendingIntervalYearly = "YEARLY"

	// SpendingIntervalAllTime limits total spending for the card's lifetime
	SpendingIntervalAllTime = "ALL_TIME"
)

// Cardholder status constants
const (
	// CardholderStatusSuccess indicates the cardholder was successfully created/verified
	CardholderStatusSuccess = "SUCCESS"

	// CardholderStatusPending indicates the cardholder is pending verification
	CardholderStatusPending = "PENDING"

	// CardholderStatusFailed indicates the cardholder creation/verification failed
	CardholderStatusFailed = "FAILED"
)

// Transaction status constants
const (
	// TransactionStatusApproved indicates the transaction was approved
	TransactionStatusApproved = "APPROVED"

	// TransactionStatusDeclined indicates the transaction was declined
	TransactionStatusDeclined = "DECLINED"

	// TransactionStatusPending indicates the transaction is pending
	TransactionStatusPending = "PENDING"

	// TransactionStatusReversed indicates the transaction was reversed
	TransactionStatusReversed = "REVERSED"
)

// Transaction type constants
const (
	// TransactionTypeFee indicates a fee transaction
	TransactionTypeFee = "FEE"

	// TransactionTypePurchase indicates a purchase transaction
	TransactionTypePurchase = "PURCHASE"

	// TransactionTypeRefund indicates a refund transaction
	TransactionTypeRefund = "REFUND"

	// TransactionTypeWithdrawal indicates a withdrawal/ATM transaction
	TransactionTypeWithdrawal = "WITHDRAWAL"

	// TransactionTypeTopup indicates a card top-up transaction
	TransactionTypeTopup = "TOPUP"
)
