package webhook

import (
	"encoding/json"
	"testing"
)

// ============================================================================
// Test Fixtures
// ============================================================================

const testEventJSON = `{
	"version": "V1.6.0",
	"event_name": "ONBOARDING",
	"event_type": "onboarding.account.create",
	"event_id": "8a78af1e-de83-43a5-b177-ecbc6a8a9fc6",
	"source_id": "f5bb6498-552e-40a5-b14b-616aa04ac1c1",
	"data": {
		"account_id": "f5bb6498-552e-40a5-b14b-616aa04ac1c1",
		"status": "PROCESSING"
	}
}`

// ============================================================================
// Event Struct Tests
// ============================================================================

func TestParseEvent(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(testEventJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse event: %v", err)
	}

	// Verify envelope fields
	if event.Version != "V1.6.0" {
		t.Errorf("Version mismatch: got %s, want V1.6.0", event.Version)
	}
	if event.EventName != "ONBOARDING" {
		t.Errorf("EventName mismatch: got %s, want ONBOARDING", event.EventName)
	}
	if event.EventType != "onboarding.account.create" {
		t.Errorf("EventType mismatch: got %s, want onboarding.account.create", event.EventType)
	}
	if event.EventID != "8a78af1e-de83-43a5-b177-ecbc6a8a9fc6" {
		t.Errorf("EventID mismatch: got %s", event.EventID)
	}
	if event.SourceID != "f5bb6498-552e-40a5-b14b-616aa04ac1c1" {
		t.Errorf("SourceID mismatch: got %s", event.SourceID)
	}
	if event.Data == nil {
		t.Error("Data should not be nil")
	}
}

func TestEventHelperMethods(t *testing.T) {
	testCases := []struct {
		name          string
		eventName     string
		eventType     string
		isOnboard     bool
		isAcquiring   bool
		isConvert     bool
		isIssuing     bool
		isBeneficiary bool
	}{
		{
			name:          "Onboarding account create",
			eventName:     EventNameOnboarding,
			eventType:     EventTypeAccountCreate,
			isOnboard:     true,
			isAcquiring:   false,
			isConvert:     false,
			isIssuing:     false,
			isBeneficiary: false,
		},
		{
			name:          "Acquiring payment intent",
			eventName:     EventNameAcquiring,
			eventType:     EventTypePaymentIntentCreated,
			isOnboard:     false,
			isAcquiring:   true,
			isConvert:     false,
			isIssuing:     false,
			isBeneficiary: false,
		},
		{
			name:          "Conversion trade settled",
			eventName:     EventNameConversion,
			eventType:     EventTypeConversionTradeSettled,
			isOnboard:     false,
			isAcquiring:   false,
			isConvert:     true,
			isIssuing:     false,
			isBeneficiary: false,
		},
		{
			name:          "Issuing card create",
			eventName:     EventNameIssuing,
			eventType:     EventTypeCardCreateSucceeded,
			isOnboard:     false,
			isAcquiring:   false,
			isConvert:     false,
			isIssuing:     true,
			isBeneficiary: false,
		},
		{
			name:          "Beneficiary successful",
			eventName:     EventNameBeneficiary,
			eventType:     EventTypeBeneficiarySuccessful,
			isOnboard:     false,
			isAcquiring:   false,
			isConvert:     false,
			isIssuing:     false,
			isBeneficiary: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := Event{
				EventName: tc.eventName,
				EventType: tc.eventType,
			}

			if event.IsOnboardingEvent() != tc.isOnboard {
				t.Errorf("IsOnboardingEvent: got %v, want %v", event.IsOnboardingEvent(), tc.isOnboard)
			}
			if event.IsAcquiringEvent() != tc.isAcquiring {
				t.Errorf("IsAcquiringEvent: got %v, want %v", event.IsAcquiringEvent(), tc.isAcquiring)
			}
			if event.IsConversionEvent() != tc.isConvert {
				t.Errorf("IsConversionEvent: got %v, want %v", event.IsConversionEvent(), tc.isConvert)
			}
			if event.IsIssuingEvent() != tc.isIssuing {
				t.Errorf("IsIssuingEvent: got %v, want %v", event.IsIssuingEvent(), tc.isIssuing)
			}
			if event.IsBeneficiaryEvent() != tc.isBeneficiary {
				t.Errorf("IsBeneficiaryEvent: got %v, want %v", event.IsBeneficiaryEvent(), tc.isBeneficiary)
			}
		})
	}
}

// ============================================================================
// Event Name Constants Tests
// ============================================================================

func TestEventNameConstants(t *testing.T) {
	testCases := []struct {
		constant string
		expected string
	}{
		{EventNameOnboarding, "ONBOARDING"},
		{EventNameAcquiring, "ACQUIRING"},
		{EventNameConversion, "CONVERSION"},
		{EventNameIssuing, "ISSUING"},
		{EventNameBeneficiary, "BENEFICIARY"},
	}

	for _, tc := range testCases {
		if tc.constant != tc.expected {
			t.Errorf("Constant mismatch: got %s, want %s", tc.constant, tc.expected)
		}
	}
}

// ============================================================================
// Event Type Constants Tests
// ============================================================================

func TestEventTypeConstants(t *testing.T) {
	// Onboarding event types
	onboardingTypes := map[string]string{
		EventTypeAccountCreate: "onboarding.account.create",
		EventTypeAccountUpdate: "onboarding.account.update",
	}

	for constant, expected := range onboardingTypes {
		if constant != expected {
			t.Errorf("Onboarding constant mismatch: got %s, want %s", constant, expected)
		}
	}

	// Acquiring event types - payment intents
	paymentIntentTypes := map[string]string{
		EventTypePaymentIntentCreated:   "acquiring.payment_intent.created",
		EventTypePaymentIntentSucceeded: "acquiring.payment_intent.succeeded",
		EventTypePaymentIntentFailed:    "acquiring.payment_intent.failed",
		EventTypePaymentIntentCanceled:  "acquiring.payment_intent.canceled",
	}

	for constant, expected := range paymentIntentTypes {
		if constant != expected {
			t.Errorf("Payment intent constant mismatch: got %s, want %s", constant, expected)
		}
	}

	// Acquiring event types - payment attempts
	paymentAttemptTypes := map[string]string{
		EventTypePaymentAttemptCreated:          "acquiring.payment_attempt.created",
		EventTypePaymentAttemptCaptureRequested: "acquiring.payment_attempt.capture_requested",
		EventTypePaymentAttemptSucceeded:        "acquiring.payment_attempt.succeeded",
		EventTypePaymentAttemptFailed:           "acquiring.payment_attempt.failed",
		EventTypePaymentAttemptCanceled:         "acquiring.payment_attempt.canceled",
	}

	for constant, expected := range paymentAttemptTypes {
		if constant != expected {
			t.Errorf("Payment attempt constant mismatch: got %s, want %s", constant, expected)
		}
	}

	// Acquiring event types - refunds
	refundTypes := map[string]string{
		EventTypeRefundCreated:   "acquiring.refund.created",
		EventTypeRefundSucceeded: "acquiring.refund.succeeded",
		EventTypeRefundFailed:    "acquiring.refund.failed",
	}

	for constant, expected := range refundTypes {
		if constant != expected {
			t.Errorf("Refund constant mismatch: got %s, want %s", constant, expected)
		}
	}

	// Conversion event types
	conversionTypes := map[string]string{
		EventTypeConversionTradeSettled:  "conversion.trade.settled",
		EventTypeConversionFundsAwaiting: "conversion.funds.awaiting",
		EventTypeConversionFundsArrived:  "conversion.funds.arrived",
	}

	for constant, expected := range conversionTypes {
		if constant != expected {
			t.Errorf("Conversion constant mismatch: got %s, want %s", constant, expected)
		}
	}

	// Issuing event types
	issuingTypes := map[string]string{
		EventTypeCardCreateSucceeded:       "card.create.succeeded",
		EventTypeCardCreateFailed:          "card.create.failed",
		EventTypeCardUpdateSucceeded:       "card.update.succeeded",
		EventTypeCardUpdateFailed:          "card.update.failed",
		EventTypeCardRechargeSucceeded:     "card.recharge.succeeded",
		EventTypeCardRechargeFailed:        "card.recharge.failed",
		EventTypeCardActivationCode:        "card.activation.code",
		EventTypeCardActivated:             "card.activated",
		EventTypeCardSuspended:             "card.suspended",
		EventTypeCardClosed:                "card.closed",
		EventTypeCardStatusUpdateSucceeded: "card.status.update.succeeded",
		EventTypeCardStatusUpdateFailed:    "card.status.update.failed",
		EventTypeIssuingFeeCard:            "issuing.fee.card",
	}

	for constant, expected := range issuingTypes {
		if constant != expected {
			t.Errorf("Issuing constant mismatch: got %s, want %s", constant, expected)
		}
	}

	// Beneficiary event types
	beneficiaryTypes := map[string]string{
		EventTypeBeneficiarySuccessful: "beneficiary.successful",
		EventTypeBeneficiaryFailed:     "beneficiary.failed",
		EventTypeBeneficiaryPending:    "beneficiary.pending",
	}

	for constant, expected := range beneficiaryTypes {
		if constant != expected {
			t.Errorf("Beneficiary constant mismatch: got %s, want %s", constant, expected)
		}
	}
}
