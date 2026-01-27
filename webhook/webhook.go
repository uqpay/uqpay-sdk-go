// Package webhook provides types for handling UQPAY webhook events.
//
// Webhooks are HTTP callbacks that notify your application when events occur
// in your UQPAY account. This package provides type-safe structures for parsing
// and handling these webhook payloads.
//
// Example usage:
//
//	func handleWebhook(w http.ResponseWriter, r *http.Request) {
//	    var event webhook.Event
//	    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
//	        http.Error(w, "Invalid payload", http.StatusBadRequest)
//	        return
//	    }
//
//	    switch event.EventType {
//	    case webhook.EventTypeAccountCreate:
//	        account, err := event.ParseAccountData()
//	        if err != nil {
//	            http.Error(w, "Failed to parse account", http.StatusBadRequest)
//	            return
//	        }
//	        // Handle account creation
//	    }
//	}
package webhook

import (
	"encoding/json"
	"fmt"
)

// Event names
const (
	EventNameOnboarding  = "ONBOARDING"
	EventNameAcquiring   = "ACQUIRING"
	EventNameConversion  = "CONVERSION"
	EventNameIssuing     = "ISSUING"
	EventNameBeneficiary = "BENEFICIARY"
)

// Event types for onboarding
const (
	EventTypeAccountCreate = "onboarding.account.create"
	EventTypeAccountUpdate = "onboarding.account.update"
)

// Event types for acquiring (payment intents)
const (
	EventTypePaymentIntentCreated   = "acquiring.payment_intent.created"
	EventTypePaymentIntentSucceeded = "acquiring.payment_intent.succeeded"
	EventTypePaymentIntentFailed    = "acquiring.payment_intent.failed"
	EventTypePaymentIntentCanceled  = "acquiring.payment_intent.canceled"
)

// Event types for acquiring (payment attempts)
const (
	EventTypePaymentAttemptCreated          = "acquiring.payment_attempt.created"
	EventTypePaymentAttemptCaptureRequested = "acquiring.payment_attempt.capture_requested"
	EventTypePaymentAttemptSucceeded        = "acquiring.payment_attempt.succeeded"
	EventTypePaymentAttemptFailed           = "acquiring.payment_attempt.failed"
	EventTypePaymentAttemptCanceled         = "acquiring.payment_attempt.canceled"
)

// Event types for acquiring (refunds)
const (
	EventTypeRefundCreated   = "acquiring.refund.created"
	EventTypeRefundSucceeded = "acquiring.refund.succeeded"
	EventTypeRefundFailed    = "acquiring.refund.failed"
)

// Event types for conversion
const (
	EventTypeConversionTradeSettled  = "conversion.trade.settled"
	EventTypeConversionFundsAwaiting = "conversion.funds.awaiting"
	EventTypeConversionFundsArrived  = "conversion.funds.arrived"
)

// Event types for issuing (card events)
const (
	EventTypeCardCreateSucceeded       = "card.create.succeeded"
	EventTypeCardCreateFailed          = "card.create.failed"
	EventTypeCardUpdateSucceeded       = "card.update.succeeded"
	EventTypeCardUpdateFailed          = "card.update.failed"
	EventTypeCardRechargeSucceeded     = "card.recharge.succeeded"
	EventTypeCardRechargeFailed        = "card.recharge.failed"
	EventTypeCardActivationCode        = "card.activation.code"
	EventTypeCardActivated             = "card.activated"
	EventTypeCardSuspended             = "card.suspended"
	EventTypeCardClosed                = "card.closed"
	EventTypeCardStatusUpdateSucceeded = "card.status.update.succeeded"
	EventTypeCardStatusUpdateFailed    = "card.status.update.failed"
)

// Event types for issuing (card transaction events)
const (
	EventTypeIssuingFeeCard = "issuing.fee.card"
)

// Event types for beneficiary
const (
	EventTypeBeneficiarySuccessful = "beneficiary.successful"
	EventTypeBeneficiaryFailed     = "beneficiary.failed"
	EventTypeBeneficiaryPending    = "beneficiary.pending"
)

// Event represents the base webhook event envelope.
// All webhook notifications share this common structure.
type Event struct {
	// Version is the API version number (e.g., "V1.6.0")
	Version string `json:"version"`

	// EventName is the category of the event (e.g., "ONBOARDING")
	EventName string `json:"event_name"`

	// EventType is the specific event type (e.g., "onboarding.account.create")
	EventType string `json:"event_type"`

	// EventID is a unique identifier for this event
	EventID string `json:"event_id"`

	// SourceID is the ID of the resource that triggered the event
	SourceID string `json:"source_id"`

	// Data contains the event-specific payload.
	// Use the Parse* methods to extract typed data.
	Data json.RawMessage `json:"data"`
}

// IsOnboardingEvent returns true if this is an onboarding-related event
func (e *Event) IsOnboardingEvent() bool {
	return e.EventName == EventNameOnboarding
}

// IsAccountCreateEvent returns true if this is an account creation event
func (e *Event) IsAccountCreateEvent() bool {
	return e.EventType == EventTypeAccountCreate
}

// IsAccountUpdateEvent returns true if this is an account update event
func (e *Event) IsAccountUpdateEvent() bool {
	return e.EventType == EventTypeAccountUpdate
}

// ParseAccountData parses the event data as an AccountData struct.
// Returns an error if the event type is not an account event or if parsing fails.
func (e *Event) ParseAccountData() (*AccountData, error) {
	if e.EventType != EventTypeAccountCreate && e.EventType != EventTypeAccountUpdate {
		return nil, fmt.Errorf("event type %s is not an account event", e.EventType)
	}

	var data AccountData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse account data: %w", err)
	}
	return &data, nil
}

// MustParseAccountData is like ParseAccountData but panics on error.
// Use this only when you are certain the event type is correct.
func (e *Event) MustParseAccountData() *AccountData {
	data, err := e.ParseAccountData()
	if err != nil {
		panic(err)
	}
	return data
}

// IsAcquiringEvent returns true if this is an acquiring-related event
func (e *Event) IsAcquiringEvent() bool {
	return e.EventName == EventNameAcquiring
}

// IsPaymentIntentEvent returns true if this is a payment intent event
func (e *Event) IsPaymentIntentEvent() bool {
	switch e.EventType {
	case EventTypePaymentIntentCreated,
		EventTypePaymentIntentSucceeded,
		EventTypePaymentIntentFailed,
		EventTypePaymentIntentCanceled:
		return true
	}
	return false
}

// ParsePaymentIntentData parses the event data as a PaymentIntentData struct.
// Returns an error if the event type is not a payment intent event or if parsing fails.
func (e *Event) ParsePaymentIntentData() (*PaymentIntentData, error) {
	if !e.IsPaymentIntentEvent() {
		return nil, fmt.Errorf("event type %s is not a payment intent event", e.EventType)
	}

	var data PaymentIntentData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse payment intent data: %w", err)
	}
	return &data, nil
}

// MustParsePaymentIntentData is like ParsePaymentIntentData but panics on error.
// Use this only when you are certain the event type is correct.
func (e *Event) MustParsePaymentIntentData() *PaymentIntentData {
	data, err := e.ParsePaymentIntentData()
	if err != nil {
		panic(err)
	}
	return data
}

// IsPaymentAttemptEvent returns true if this is a payment attempt event
func (e *Event) IsPaymentAttemptEvent() bool {
	switch e.EventType {
	case EventTypePaymentAttemptCreated,
		EventTypePaymentAttemptCaptureRequested,
		EventTypePaymentAttemptSucceeded,
		EventTypePaymentAttemptFailed,
		EventTypePaymentAttemptCanceled:
		return true
	}
	return false
}

// ParsePaymentAttemptData parses the event data as a PaymentAttemptData struct.
// Returns an error if the event type is not a payment attempt event or if parsing fails.
func (e *Event) ParsePaymentAttemptData() (*PaymentAttemptData, error) {
	if !e.IsPaymentAttemptEvent() {
		return nil, fmt.Errorf("event type %s is not a payment attempt event", e.EventType)
	}

	var data PaymentAttemptData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse payment attempt data: %w", err)
	}
	return &data, nil
}

// MustParsePaymentAttemptData is like ParsePaymentAttemptData but panics on error.
// Use this only when you are certain the event type is correct.
func (e *Event) MustParsePaymentAttemptData() *PaymentAttemptData {
	data, err := e.ParsePaymentAttemptData()
	if err != nil {
		panic(err)
	}
	return data
}

// IsRefundEvent returns true if this is a refund event
func (e *Event) IsRefundEvent() bool {
	switch e.EventType {
	case EventTypeRefundCreated,
		EventTypeRefundSucceeded,
		EventTypeRefundFailed:
		return true
	}
	return false
}

// ParseRefundData parses the event data as a RefundData struct.
// Returns an error if the event type is not a refund event or if parsing fails.
func (e *Event) ParseRefundData() (*RefundData, error) {
	if !e.IsRefundEvent() {
		return nil, fmt.Errorf("event type %s is not a refund event", e.EventType)
	}

	var data RefundData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse refund data: %w", err)
	}
	return &data, nil
}

// MustParseRefundData is like ParseRefundData but panics on error.
// Use this only when you are certain the event type is correct.
func (e *Event) MustParseRefundData() *RefundData {
	data, err := e.ParseRefundData()
	if err != nil {
		panic(err)
	}
	return data
}

// IsConversionEvent returns true if this is a conversion-related event
func (e *Event) IsConversionEvent() bool {
	return e.EventName == EventNameConversion
}

// IsConversionTradeSettledEvent returns true if this is a conversion trade settled event
func (e *Event) IsConversionTradeSettledEvent() bool {
	return e.EventType == EventTypeConversionTradeSettled
}

// ParseConversionData parses the event data as a ConversionData struct.
// Returns an error if the event type is not a conversion event or if parsing fails.
func (e *Event) ParseConversionData() (*ConversionData, error) {
	if !e.IsConversionEvent() {
		return nil, fmt.Errorf("event type %s is not a conversion event", e.EventType)
	}

	var data ConversionData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse conversion data: %w", err)
	}
	return &data, nil
}

// MustParseConversionData is like ParseConversionData but panics on error.
// Use this only when you are certain the event type is correct.
func (e *Event) MustParseConversionData() *ConversionData {
	data, err := e.ParseConversionData()
	if err != nil {
		panic(err)
	}
	return data
}

// IsIssuingEvent returns true if this is an issuing-related event
func (e *Event) IsIssuingEvent() bool {
	return e.EventName == EventNameIssuing
}

// IsCardEvent returns true if this is a card event
func (e *Event) IsCardEvent() bool {
	switch e.EventType {
	case EventTypeCardCreateSucceeded,
		EventTypeCardCreateFailed,
		EventTypeCardUpdateSucceeded,
		EventTypeCardUpdateFailed,
		EventTypeCardRechargeSucceeded,
		EventTypeCardRechargeFailed,
		EventTypeCardActivationCode,
		EventTypeCardActivated,
		EventTypeCardSuspended,
		EventTypeCardClosed,
		EventTypeCardStatusUpdateSucceeded,
		EventTypeCardStatusUpdateFailed:
		return true
	}
	return false
}

// IsCardStatusUpdateEvent returns true if this is a card status update event
func (e *Event) IsCardStatusUpdateEvent() bool {
	switch e.EventType {
	case EventTypeCardStatusUpdateSucceeded,
		EventTypeCardStatusUpdateFailed:
		return true
	}
	return false
}

// IsCardCreateOrUpdateEvent returns true if this is a card create or update event
func (e *Event) IsCardCreateOrUpdateEvent() bool {
	switch e.EventType {
	case EventTypeCardCreateSucceeded,
		EventTypeCardCreateFailed,
		EventTypeCardUpdateSucceeded,
		EventTypeCardUpdateFailed:
		return true
	}
	return false
}

// ParseCardData parses the event data as a CardData struct.
// Returns an error if the event type is not a card create/update event or if parsing fails.
func (e *Event) ParseCardData() (*CardData, error) {
	if !e.IsCardCreateOrUpdateEvent() {
		return nil, fmt.Errorf("event type %s is not a card create/update event", e.EventType)
	}

	var data CardData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse card data: %w", err)
	}
	return &data, nil
}

// MustParseCardData is like ParseCardData but panics on error.
// Use this only when you are certain the event type is correct.
func (e *Event) MustParseCardData() *CardData {
	data, err := e.ParseCardData()
	if err != nil {
		panic(err)
	}
	return data
}

// IsCardRechargeEvent returns true if this is a card recharge event
func (e *Event) IsCardRechargeEvent() bool {
	switch e.EventType {
	case EventTypeCardRechargeSucceeded,
		EventTypeCardRechargeFailed:
		return true
	}
	return false
}

// IsCardActivationCodeEvent returns true if this is a card activation code event
func (e *Event) IsCardActivationCodeEvent() bool {
	return e.EventType == EventTypeCardActivationCode
}

// ParseCardActivationCodeData parses the event data as a CardActivationCodeData struct.
// Returns an error if the event type is not a card activation code event or if parsing fails.
func (e *Event) ParseCardActivationCodeData() (*CardActivationCodeData, error) {
	if !e.IsCardActivationCodeEvent() {
		return nil, fmt.Errorf("event type %s is not a card activation code event", e.EventType)
	}

	var data CardActivationCodeData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse card activation code data: %w", err)
	}
	return &data, nil
}

// MustParseCardActivationCodeData is like ParseCardActivationCodeData but panics on error.
// Use this only when you are certain the event type is correct.
func (e *Event) MustParseCardActivationCodeData() *CardActivationCodeData {
	data, err := e.ParseCardActivationCodeData()
	if err != nil {
		panic(err)
	}
	return data
}

// ParseCardRechargeData parses the event data as a CardRechargeData struct.
// Returns an error if the event type is not a card recharge event or if parsing fails.
func (e *Event) ParseCardRechargeData() (*CardRechargeData, error) {
	if !e.IsCardRechargeEvent() {
		return nil, fmt.Errorf("event type %s is not a card recharge event", e.EventType)
	}

	var data CardRechargeData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse card recharge data: %w", err)
	}
	return &data, nil
}

// MustParseCardRechargeData is like ParseCardRechargeData but panics on error.
// Use this only when you are certain the event type is correct.
func (e *Event) MustParseCardRechargeData() *CardRechargeData {
	data, err := e.ParseCardRechargeData()
	if err != nil {
		panic(err)
	}
	return data
}

// ParseCardStatusUpdateData parses the event data as a CardStatusUpdateData struct.
// Returns an error if the event type is not a card status update event or if parsing fails.
func (e *Event) ParseCardStatusUpdateData() (*CardStatusUpdateData, error) {
	if !e.IsCardStatusUpdateEvent() {
		return nil, fmt.Errorf("event type %s is not a card status update event", e.EventType)
	}

	var data CardStatusUpdateData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse card status update data: %w", err)
	}
	return &data, nil
}

// MustParseCardStatusUpdateData is like ParseCardStatusUpdateData but panics on error.
// Use this only when you are certain the event type is correct.
func (e *Event) MustParseCardStatusUpdateData() *CardStatusUpdateData {
	data, err := e.ParseCardStatusUpdateData()
	if err != nil {
		panic(err)
	}
	return data
}

// IsCardTransactionEvent returns true if this is a card transaction event
func (e *Event) IsCardTransactionEvent() bool {
	switch e.EventType {
	case EventTypeIssuingFeeCard:
		return true
	}
	return false
}

// ParseCardTransactionData parses the event data as a CardTransactionData struct.
// Returns an error if the event type is not a card transaction event or if parsing fails.
func (e *Event) ParseCardTransactionData() (*CardTransactionData, error) {
	if !e.IsCardTransactionEvent() {
		return nil, fmt.Errorf("event type %s is not a card transaction event", e.EventType)
	}

	var data CardTransactionData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse card transaction data: %w", err)
	}
	return &data, nil
}

// MustParseCardTransactionData is like ParseCardTransactionData but panics on error.
// Use this only when you are certain the event type is correct.
func (e *Event) MustParseCardTransactionData() *CardTransactionData {
	data, err := e.ParseCardTransactionData()
	if err != nil {
		panic(err)
	}
	return data
}

// IsBeneficiaryEvent returns true if this is a beneficiary-related event
func (e *Event) IsBeneficiaryEvent() bool {
	return e.EventName == EventNameBeneficiary
}

// IsBeneficiarySuccessfulEvent returns true if this is a beneficiary successful event
func (e *Event) IsBeneficiarySuccessfulEvent() bool {
	return e.EventType == EventTypeBeneficiarySuccessful
}

// IsBeneficiaryFailedEvent returns true if this is a beneficiary failed event
func (e *Event) IsBeneficiaryFailedEvent() bool {
	return e.EventType == EventTypeBeneficiaryFailed
}

// ParseBeneficiaryData parses the event data as a BeneficiaryData struct.
// Returns an error if the event type is not a beneficiary event or if parsing fails.
func (e *Event) ParseBeneficiaryData() (*BeneficiaryData, error) {
	if !e.IsBeneficiaryEvent() {
		return nil, fmt.Errorf("event type %s is not a beneficiary event", e.EventType)
	}

	var data BeneficiaryData
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return nil, fmt.Errorf("failed to parse beneficiary data: %w", err)
	}
	return &data, nil
}

// MustParseBeneficiaryData is like ParseBeneficiaryData but panics on error.
// Use this only when you are certain the event type is correct.
func (e *Event) MustParseBeneficiaryData() *BeneficiaryData {
	data, err := e.ParseBeneficiaryData()
	if err != nil {
		panic(err)
	}
	return data
}
