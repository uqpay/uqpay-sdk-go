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
	EventNameOnboarding = "ONBOARDING"
	EventNameAcquiring  = "ACQUIRING"
	EventNameConversion = "CONVERSION"
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

// Event types for conversion
const (
	EventTypeConversionTradeSettled  = "conversion.trade.settled"
	EventTypeConversionFundsAwaiting = "conversion.funds.awaiting"
	EventTypeConversionFundsArrived  = "conversion.funds.arrived"
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
