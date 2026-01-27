package webhook

import (
	"encoding/json"
	"testing"
)

// ============================================================================
// Test Fixtures - Payment Intent Webhooks
// ============================================================================

const paymentIntentCreatedWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_intent.created",
	"event_id":"5008b4da-a5e0-4a07-86f1-42cf6931afed",
	"source_id":"PI2013833849980588032",
	"data":{
		"amount":"101",
		"cancel_time":null,
		"cancellation_reason":"",
		"complete_time":null,
		"create_time":"2026-01-21T12:39:39.716846826+08:00",
		"currency":"USD",
		"description":"Test payment intent",
		"intent_status":"REQUIRES_PAYMENT_METHOD",
		"merchant_order_id":"test-order-002",
		"metadata":{"test":"true"},
		"payment_intent_id":"PI2013833849980588032",
		"payment_method":null
	}
}`

const paymentIntentSucceededWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_intent.succeeded",
	"event_id":"2a88cf2e-5c1d-4c6c-a6aa-1cccdf259e5e",
	"source_id":"PI2013833849980588032",
	"data":{
		"amount":"101",
		"cancel_time":null,
		"cancellation_reason":"",
		"complete_time":"2026-01-21T13:15:22.456+08:00",
		"create_time":"2026-01-21T12:39:39.717+08:00",
		"currency":"USD",
		"description":"Test payment intent",
		"intent_status":"SUCCEEDED",
		"merchant_order_id":"test-order-002",
		"metadata":{"test":"true"},
		"payment_intent_id":"PI2013833849980588032",
		"payment_method":{
			"type":"card",
			"card":{
				"brand":"visa",
				"last4":"4242",
				"exp_month":12,
				"exp_year":2027,
				"funding":"credit",
				"country":"US"
			}
		}
	}
}`

const paymentIntentFailedWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_intent.failed",
	"event_id":"1e77dbee-4b0b-4b5b-9599-9bbbcf148fd4",
	"source_id":"PI2013833849980588032",
	"data":{
		"amount":"101",
		"cancel_time":null,
		"cancellation_reason":"",
		"complete_time":"2026-01-21T13:09:47.727+08:00",
		"create_time":"2026-01-21T12:39:39.717+08:00",
		"currency":"USD",
		"description":"Test payment intent",
		"intent_status":"FAILED",
		"merchant_order_id":"test-order-002",
		"metadata":{"test":"true"},
		"payment_intent_id":"PI2013833849980588032",
		"payment_method":null
	}
}`

const paymentIntentCanceledWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_intent.canceled",
	"event_id":"3b99df3f-6d2e-5d7d-b7bb-2dddef36af6f",
	"source_id":"PI2013833849980588032",
	"data":{
		"amount":"101",
		"cancel_time":"2026-01-21T13:20:15.123+08:00",
		"cancellation_reason":"requested_by_customer",
		"complete_time":null,
		"create_time":"2026-01-21T12:39:39.717+08:00",
		"currency":"USD",
		"description":"Test payment intent",
		"intent_status":"CANCELED",
		"merchant_order_id":"test-order-002",
		"metadata":{"test":"true"},
		"payment_intent_id":"PI2013833849980588032",
		"payment_method":null
	}
}`

// ============================================================================
// Test Fixtures - Payment Attempt Webhooks
// ============================================================================

const paymentAttemptCreatedWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_attempt.created",
	"event_id":"3b1e591c-c9cd-4e26-be19-d83d7ec907e1",
	"source_id":"PA2013848038472159232",
	"data":{
		"amount":"0.01",
		"attempt_status":"INITIATED",
		"cancel_time":null,
		"cancellation_reason":"",
		"captured_amount":"0.01",
		"complete_time":null,
		"create_time":"2026-01-21T13:36:02.516375218+08:00",
		"currency":"USD",
		"failure_code":"",
		"merchant_order_id":"test-AlipayCN_QRCode-001",
		"payment_attempt_id":"PA2013848038472159232",
		"payment_intent_id":"PI2013848035972354048",
		"payment_method":{
			"alipaycn":{
				"flow":"qrcode",
				"os_type":""
			},
			"type":"alipaycn"
		},
		"refunded_amount":"0"
	}
}`

const paymentAttemptSucceededWebhookJSON = `{
	"version":"V1.6.0",
	"event_name":"ACQUIRING",
	"event_type":"acquiring.payment_attempt.succeeded",
	"event_id":"4c2f602d-d0de-5f37-cf2a-e94e8fd018f2",
	"source_id":"PA2013848038472159232",
	"data":{
		"amount":"0.01",
		"attempt_status":"SUCCEEDED",
		"cancel_time":null,
		"cancellation_reason":"",
		"captured_amount":"0.01",
		"complete_time":"2026-01-21T13:40:15.789+08:00",
		"create_time":"2026-01-21T13:36:02.516+08:00",
		"currency":"USD",
		"failure_code":"",
		"merchant_order_id":"test-AlipayCN_QRCode-001",
		"payment_attempt_id":"PA2013848038472159232",
		"payment_intent_id":"PI2013848035972354048",
		"payment_method":{
			"alipaycn":{
				"flow":"qrcode",
				"os_type":""
			},
			"type":"alipaycn"
		},
		"refunded_amount":"0"
	}
}`

// ============================================================================
// Payment Intent Tests
// ============================================================================

func TestParsePaymentIntentCreatedWebhook(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(paymentIntentCreatedWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment intent webhook event: %v", err)
	}

	// Verify envelope fields
	if event.Version != "V1.6.0" {
		t.Errorf("Version mismatch: got %s, want V1.6.0", event.Version)
	}
	if event.EventName != "ACQUIRING" {
		t.Errorf("EventName mismatch: got %s, want ACQUIRING", event.EventName)
	}
	if event.EventType != "acquiring.payment_intent.created" {
		t.Errorf("EventType mismatch: got %s", event.EventType)
	}

	// Verify helper methods
	if !event.IsAcquiringEvent() {
		t.Error("IsAcquiringEvent should return true")
	}
	if !event.IsPaymentIntentEvent() {
		t.Error("IsPaymentIntentEvent should return true")
	}
	if event.IsOnboardingEvent() {
		t.Error("IsOnboardingEvent should return false")
	}
}

func TestParsePaymentIntentData(t *testing.T) {
	var event Event
	json.Unmarshal([]byte(paymentIntentCreatedWebhookJSON), &event)

	paymentIntent, err := event.ParsePaymentIntentData()
	if err != nil {
		t.Fatalf("Failed to parse payment intent data: %v", err)
	}

	// Verify core fields
	if paymentIntent.PaymentIntentID != "PI2013833849980588032" {
		t.Errorf("PaymentIntentID mismatch: got %s", paymentIntent.PaymentIntentID)
	}
	if paymentIntent.Amount != "101" {
		t.Errorf("Amount mismatch: got %s", paymentIntent.Amount)
	}
	if paymentIntent.Currency != "USD" {
		t.Errorf("Currency mismatch: got %s", paymentIntent.Currency)
	}
	if paymentIntent.Description != "Test payment intent" {
		t.Errorf("Description mismatch: got %s", paymentIntent.Description)
	}
	if paymentIntent.IntentStatus != "REQUIRES_PAYMENT_METHOD" {
		t.Errorf("IntentStatus mismatch: got %s", paymentIntent.IntentStatus)
	}
	if paymentIntent.MerchantOrderID != "test-order-002" {
		t.Errorf("MerchantOrderID mismatch: got %s", paymentIntent.MerchantOrderID)
	}

	// Verify metadata
	if paymentIntent.Metadata == nil {
		t.Fatal("Metadata should not be nil")
	}
	if paymentIntent.Metadata["test"] != "true" {
		t.Errorf("Metadata[test] mismatch: got %s", paymentIntent.Metadata["test"])
	}

	// Verify nullable fields
	if paymentIntent.PaymentMethod != nil {
		t.Error("PaymentMethod should be nil")
	}
}

func TestParsePaymentIntentSucceededWebhook(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(paymentIntentSucceededWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment intent succeeded webhook: %v", err)
	}

	if event.EventType != EventTypePaymentIntentSucceeded {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, EventTypePaymentIntentSucceeded)
	}

	paymentIntent, err := event.ParsePaymentIntentData()
	if err != nil {
		t.Fatalf("Failed to parse payment intent data: %v", err)
	}

	// Verify status is SUCCEEDED
	if paymentIntent.IntentStatus != IntentStatusSucceeded {
		t.Errorf("IntentStatus mismatch: got %s, want %s", paymentIntent.IntentStatus, IntentStatusSucceeded)
	}

	// Verify complete_time is set
	if paymentIntent.CompleteTime == nil {
		t.Error("CompleteTime should not be nil for succeeded event")
	}

	// Verify payment method is populated
	if paymentIntent.PaymentMethod == nil {
		t.Fatal("PaymentMethod should not be nil for succeeded event")
	}
	if paymentIntent.PaymentMethod.Type != "card" {
		t.Errorf("PaymentMethod.Type mismatch: got %s", paymentIntent.PaymentMethod.Type)
	}
	if paymentIntent.PaymentMethod.Card == nil {
		t.Fatal("PaymentMethod.Card should not be nil")
	}
	if paymentIntent.PaymentMethod.Card.Brand != "visa" {
		t.Errorf("Card.Brand mismatch: got %s", paymentIntent.PaymentMethod.Card.Brand)
	}
	if paymentIntent.PaymentMethod.Card.Last4 != "4242" {
		t.Errorf("Card.Last4 mismatch: got %s", paymentIntent.PaymentMethod.Card.Last4)
	}
}

func TestParsePaymentIntentCanceledWebhook(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(paymentIntentCanceledWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment intent canceled webhook: %v", err)
	}

	paymentIntent, err := event.ParsePaymentIntentData()
	if err != nil {
		t.Fatalf("Failed to parse payment intent data: %v", err)
	}

	// Verify status is CANCELED
	if paymentIntent.IntentStatus != IntentStatusCanceled {
		t.Errorf("IntentStatus mismatch: got %s, want %s", paymentIntent.IntentStatus, IntentStatusCanceled)
	}

	// Verify cancel_time is set
	if paymentIntent.CancelTime == nil {
		t.Error("CancelTime should not be nil for canceled event")
	}

	// Verify cancellation_reason is set
	if paymentIntent.CancellationReason != "requested_by_customer" {
		t.Errorf("CancellationReason mismatch: got %s", paymentIntent.CancellationReason)
	}
}

func TestPaymentIntentAllEventTypes(t *testing.T) {
	testCases := []struct {
		name             string
		json             string
		expectedType     string
		expectedStatus   string
		hasCompleteTime  bool
		hasCancelTime    bool
		hasPaymentMethod bool
	}{
		{
			name:             "Created",
			json:             paymentIntentCreatedWebhookJSON,
			expectedType:     EventTypePaymentIntentCreated,
			expectedStatus:   IntentStatusRequiresPaymentMethod,
			hasCompleteTime:  false,
			hasCancelTime:    false,
			hasPaymentMethod: false,
		},
		{
			name:             "Succeeded",
			json:             paymentIntentSucceededWebhookJSON,
			expectedType:     EventTypePaymentIntentSucceeded,
			expectedStatus:   IntentStatusSucceeded,
			hasCompleteTime:  true,
			hasCancelTime:    false,
			hasPaymentMethod: true,
		},
		{
			name:             "Failed",
			json:             paymentIntentFailedWebhookJSON,
			expectedType:     EventTypePaymentIntentFailed,
			expectedStatus:   IntentStatusFailed,
			hasCompleteTime:  true,
			hasCancelTime:    false,
			hasPaymentMethod: false,
		},
		{
			name:             "Canceled",
			json:             paymentIntentCanceledWebhookJSON,
			expectedType:     EventTypePaymentIntentCanceled,
			expectedStatus:   IntentStatusCanceled,
			hasCompleteTime:  false,
			hasCancelTime:    true,
			hasPaymentMethod: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var event Event
			err := json.Unmarshal([]byte(tc.json), &event)
			if err != nil {
				t.Fatalf("Failed to parse event: %v", err)
			}

			if event.EventType != tc.expectedType {
				t.Errorf("EventType mismatch: got %s, want %s", event.EventType, tc.expectedType)
			}

			if !event.IsPaymentIntentEvent() {
				t.Error("IsPaymentIntentEvent should return true")
			}

			paymentIntent, err := event.ParsePaymentIntentData()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			if paymentIntent.IntentStatus != tc.expectedStatus {
				t.Errorf("IntentStatus mismatch: got %s, want %s", paymentIntent.IntentStatus, tc.expectedStatus)
			}

			if tc.hasCompleteTime && paymentIntent.CompleteTime == nil {
				t.Error("CompleteTime should not be nil")
			}
			if tc.hasCancelTime && paymentIntent.CancelTime == nil {
				t.Error("CancelTime should not be nil")
			}
			if tc.hasPaymentMethod && paymentIntent.PaymentMethod == nil {
				t.Error("PaymentMethod should not be nil")
			}
		})
	}
}

// ============================================================================
// Payment Attempt Tests
// ============================================================================

func TestParsePaymentAttemptCreatedWebhook(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(paymentAttemptCreatedWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment attempt created webhook: %v", err)
	}

	if event.EventType != EventTypePaymentAttemptCreated {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, EventTypePaymentAttemptCreated)
	}

	if !event.IsPaymentAttemptEvent() {
		t.Error("IsPaymentAttemptEvent should return true")
	}

	if !event.IsAcquiringEvent() {
		t.Error("IsAcquiringEvent should return true")
	}

	// Should NOT be a payment intent event
	if event.IsPaymentIntentEvent() {
		t.Error("IsPaymentIntentEvent should return false for attempt event")
	}

	attempt, err := event.ParsePaymentAttemptData()
	if err != nil {
		t.Fatalf("Failed to parse payment attempt data: %v", err)
	}

	// Verify core fields
	if attempt.PaymentAttemptID != "PA2013848038472159232" {
		t.Errorf("PaymentAttemptID mismatch: got %s", attempt.PaymentAttemptID)
	}
	if attempt.PaymentIntentID != "PI2013848035972354048" {
		t.Errorf("PaymentIntentID mismatch: got %s", attempt.PaymentIntentID)
	}
	if attempt.Amount != "0.01" {
		t.Errorf("Amount mismatch: got %s", attempt.Amount)
	}
	if attempt.Currency != "USD" {
		t.Errorf("Currency mismatch: got %s", attempt.Currency)
	}
	if attempt.AttemptStatus != AttemptStatusInitiated {
		t.Errorf("AttemptStatus mismatch: got %s", attempt.AttemptStatus)
	}

	// Verify payment method (Alipay)
	if attempt.PaymentMethod == nil {
		t.Fatal("PaymentMethod should not be nil")
	}
	if attempt.PaymentMethod.Type != "alipaycn" {
		t.Errorf("PaymentMethod.Type mismatch: got %s", attempt.PaymentMethod.Type)
	}
	if attempt.PaymentMethod.AlipayCN == nil {
		t.Fatal("PaymentMethod.AlipayCN should not be nil")
	}
	if attempt.PaymentMethod.AlipayCN.Flow != "qrcode" {
		t.Errorf("AlipayCN.Flow mismatch: got %s", attempt.PaymentMethod.AlipayCN.Flow)
	}
}

func TestParsePaymentAttemptSucceededWebhook(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(paymentAttemptSucceededWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse payment attempt succeeded webhook: %v", err)
	}

	if event.EventType != EventTypePaymentAttemptSucceeded {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, EventTypePaymentAttemptSucceeded)
	}

	attempt, err := event.ParsePaymentAttemptData()
	if err != nil {
		t.Fatalf("Failed to parse payment attempt data: %v", err)
	}

	if attempt.AttemptStatus != AttemptStatusSucceeded {
		t.Errorf("AttemptStatus mismatch: got %s, want %s", attempt.AttemptStatus, AttemptStatusSucceeded)
	}

	// Verify complete_time is set
	if attempt.CompleteTime == nil {
		t.Error("CompleteTime should not be nil for succeeded event")
	}

	// Verify amounts
	if attempt.CapturedAmount != "0.01" {
		t.Errorf("CapturedAmount mismatch: got %s", attempt.CapturedAmount)
	}
	if attempt.RefundedAmount != "0" {
		t.Errorf("RefundedAmount mismatch: got %s", attempt.RefundedAmount)
	}
}

// ============================================================================
// Acquiring Status Constants Tests
// ============================================================================

func TestAcquiringStatusConstants(t *testing.T) {
	// Intent status constants
	if IntentStatusRequiresPaymentMethod != "REQUIRES_PAYMENT_METHOD" {
		t.Errorf("IntentStatusRequiresPaymentMethod mismatch")
	}
	if IntentStatusSucceeded != "SUCCEEDED" {
		t.Errorf("IntentStatusSucceeded mismatch")
	}
	if IntentStatusCanceled != "CANCELED" {
		t.Errorf("IntentStatusCanceled mismatch")
	}
	if IntentStatusFailed != "FAILED" {
		t.Errorf("IntentStatusFailed mismatch")
	}

	// Attempt status constants
	if AttemptStatusInitiated != "INITIATED" {
		t.Errorf("AttemptStatusInitiated mismatch")
	}
	if AttemptStatusSucceeded != "SUCCEEDED" {
		t.Errorf("AttemptStatusSucceeded mismatch")
	}
	if AttemptStatusFailed != "FAILED" {
		t.Errorf("AttemptStatusFailed mismatch")
	}
	if AttemptStatusCanceled != "CANCELED" {
		t.Errorf("AttemptStatusCanceled mismatch")
	}

	// Refund status constants
	if RefundStatusInitiated != "INITIATED" {
		t.Errorf("RefundStatusInitiated mismatch")
	}
	if RefundStatusSucceeded != "SUCCEEDED" {
		t.Errorf("RefundStatusSucceeded mismatch")
	}
	if RefundStatusFailed != "FAILED" {
		t.Errorf("RefundStatusFailed mismatch")
	}
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestParsePaymentIntentData_WrongEventType(t *testing.T) {
	wrongTypeJSON := `{
		"version": "V1.6.0",
		"event_name": "ONBOARDING",
		"event_type": "onboarding.account.create",
		"event_id": "test-id",
		"source_id": "test-source",
		"data": {}
	}`

	var event Event
	json.Unmarshal([]byte(wrongTypeJSON), &event)

	_, err := event.ParsePaymentIntentData()
	if err == nil {
		t.Error("ParsePaymentIntentData should fail for non-payment intent event type")
	}
}

func TestParsePaymentAttemptData_WrongEventType(t *testing.T) {
	var event Event
	json.Unmarshal([]byte(paymentIntentCreatedWebhookJSON), &event)

	_, err := event.ParsePaymentAttemptData()
	if err == nil {
		t.Error("ParsePaymentAttemptData should fail for payment intent event type")
	}
}
