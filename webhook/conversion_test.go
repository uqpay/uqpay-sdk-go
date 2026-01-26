package webhook

import (
	"encoding/json"
	"testing"
)

// ============================================================================
// Test Fixtures - Conversion Webhooks
// ============================================================================

const conversionTradeSettledWebhookJSON = `{
	"version": "V1.6.0",
	"event_name": "CONVERSION",
	"event_type": "conversion.trade.settled",
	"event_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
	"source_id": "CONV-2024-001",
	"data": {
		"account_id": "acc-123456",
		"account_name": "Test Company Ltd",
		"buy_amount": "38.51",
		"buy_currency": "SGD",
		"client_rate": "1.3456",
		"conversion_id": "CONV-2024-001",
		"conversion_status": "TRADE_SETTLED",
		"conversion_way": "API",
		"create_time": "2026-01-21T10:30:00+08:00",
		"creator": "api_user",
		"direct_id": "direct-789",
		"sell_amount": "100",
		"sell_currency": "USD",
		"settle_time": "2026-01-21T10:35:00+08:00",
		"short_reference_id": "CV260121-ABC123"
	}
}`

const conversionFundsAwaitingWebhookJSON = `{
	"version": "V1.6.0",
	"event_name": "CONVERSION",
	"event_type": "conversion.funds.awaiting",
	"event_id": "b2c3d4e5-f6g7-8901-bcde-fg2345678901",
	"source_id": "CONV-2024-002",
	"data": {
		"account_id": "acc-123456",
		"account_name": "Test Company Ltd",
		"buy_amount": "75.00",
		"buy_currency": "EUR",
		"client_rate": "0.9200",
		"conversion_id": "CONV-2024-002",
		"conversion_status": "AWAITING_FUNDS",
		"conversion_way": "WEB",
		"create_time": "2026-01-21T11:00:00+08:00",
		"creator": "web_user",
		"direct_id": "direct-789",
		"sell_amount": "81.52",
		"sell_currency": "USD",
		"short_reference_id": "CV260121-DEF456"
	}
}`

const conversionFundsArrivedWebhookJSON = `{
	"version": "V1.6.0",
	"event_name": "CONVERSION",
	"event_type": "conversion.funds.arrived",
	"event_id": "c3d4e5f6-g7h8-9012-cdef-gh3456789012",
	"source_id": "CONV-2024-003",
	"data": {
		"account_id": "acc-123456",
		"account_name": "Test Company Ltd",
		"buy_amount": "5000.00",
		"buy_currency": "GBP",
		"client_rate": "0.7850",
		"conversion_id": "CONV-2024-003",
		"conversion_status": "FUNDS_ARRIVED",
		"conversion_way": "API",
		"create_time": "2026-01-21T09:00:00+08:00",
		"creator": "api_user",
		"direct_id": "direct-789",
		"sell_amount": "6369.43",
		"sell_currency": "USD",
		"settle_time": "2026-01-21T12:00:00+08:00",
		"short_reference_id": "CV260121-GHI789"
	}
}`

// ============================================================================
// Conversion Trade Settled Tests
// ============================================================================

func TestParseConversionTradeSettledWebhook(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(conversionTradeSettledWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse conversion trade settled webhook: %v", err)
	}

	// Verify envelope fields
	if event.Version != "V1.6.0" {
		t.Errorf("Version mismatch: got %s, want V1.6.0", event.Version)
	}
	if event.EventName != EventNameConversion {
		t.Errorf("EventName mismatch: got %s, want %s", event.EventName, EventNameConversion)
	}
	if event.EventType != EventTypeConversionTradeSettled {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, EventTypeConversionTradeSettled)
	}

	// Verify helper methods
	if !event.IsConversionEvent() {
		t.Error("IsConversionEvent should return true")
	}
	if !event.IsConversionTradeSettledEvent() {
		t.Error("IsConversionTradeSettledEvent should return true")
	}
	if event.IsAcquiringEvent() {
		t.Error("IsAcquiringEvent should return false")
	}
}

func TestParseConversionData_TradeSettled(t *testing.T) {
	var event Event
	json.Unmarshal([]byte(conversionTradeSettledWebhookJSON), &event)

	conversion, err := event.ParseConversionData()
	if err != nil {
		t.Fatalf("Failed to parse conversion data: %v", err)
	}

	// Verify core fields
	if conversion.ConversionID != "CONV-2024-001" {
		t.Errorf("ConversionID mismatch: got %s", conversion.ConversionID)
	}
	if conversion.AccountID != "acc-123456" {
		t.Errorf("AccountID mismatch: got %s", conversion.AccountID)
	}
	if conversion.AccountName != "Test Company Ltd" {
		t.Errorf("AccountName mismatch: got %s", conversion.AccountName)
	}

	// Verify buy side
	if conversion.BuyAmount != "38.51" {
		t.Errorf("BuyAmount mismatch: got %s", conversion.BuyAmount)
	}
	if conversion.BuyCurrency != "SGD" {
		t.Errorf("BuyCurrency mismatch: got %s", conversion.BuyCurrency)
	}

	// Verify sell side
	if conversion.SellAmount != "100" {
		t.Errorf("SellAmount mismatch: got %s", conversion.SellAmount)
	}
	if conversion.SellCurrency != "USD" {
		t.Errorf("SellCurrency mismatch: got %s", conversion.SellCurrency)
	}

	// Verify rate and status
	if conversion.ClientRate != "1.3456" {
		t.Errorf("ClientRate mismatch: got %s", conversion.ClientRate)
	}
	if conversion.ConversionStatus != ConversionStatusTradeSettled {
		t.Errorf("ConversionStatus mismatch: got %s", conversion.ConversionStatus)
	}
	if conversion.ConversionWay != ConversionWayAPI {
		t.Errorf("ConversionWay mismatch: got %s", conversion.ConversionWay)
	}

	// Verify settle time is set
	if conversion.SettleTime != "2026-01-21T10:35:00+08:00" {
		t.Errorf("SettleTime mismatch: got %s", conversion.SettleTime)
	}

	// Verify reference fields
	if conversion.ShortReferenceID != "CV260121-ABC123" {
		t.Errorf("ShortReferenceID mismatch: got %s", conversion.ShortReferenceID)
	}
	if conversion.DirectID != "direct-789" {
		t.Errorf("DirectID mismatch: got %s", conversion.DirectID)
	}
}

// ============================================================================
// Conversion Funds Awaiting Tests
// ============================================================================

func TestParseConversionFundsAwaitingWebhook(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(conversionFundsAwaitingWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse conversion funds awaiting webhook: %v", err)
	}

	if event.EventType != EventTypeConversionFundsAwaiting {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, EventTypeConversionFundsAwaiting)
	}

	if !event.IsConversionEvent() {
		t.Error("IsConversionEvent should return true")
	}

	conversion, err := event.ParseConversionData()
	if err != nil {
		t.Fatalf("Failed to parse conversion data: %v", err)
	}

	if conversion.ConversionStatus != ConversionStatusAwaitingFunds {
		t.Errorf("ConversionStatus mismatch: got %s, want %s", conversion.ConversionStatus, ConversionStatusAwaitingFunds)
	}
	if conversion.ConversionWay != ConversionWayWeb {
		t.Errorf("ConversionWay mismatch: got %s, want %s", conversion.ConversionWay, ConversionWayWeb)
	}

	// Settle time should be empty for awaiting status
	if conversion.SettleTime != "" {
		t.Errorf("SettleTime should be empty for awaiting status, got %s", conversion.SettleTime)
	}
}

// ============================================================================
// Conversion Funds Arrived Tests
// ============================================================================

func TestParseConversionFundsArrivedWebhook(t *testing.T) {
	var event Event
	err := json.Unmarshal([]byte(conversionFundsArrivedWebhookJSON), &event)
	if err != nil {
		t.Fatalf("Failed to parse conversion funds arrived webhook: %v", err)
	}

	if event.EventType != EventTypeConversionFundsArrived {
		t.Errorf("EventType mismatch: got %s, want %s", event.EventType, EventTypeConversionFundsArrived)
	}

	conversion, err := event.ParseConversionData()
	if err != nil {
		t.Fatalf("Failed to parse conversion data: %v", err)
	}

	if conversion.ConversionStatus != ConversionStatusFundsArrived {
		t.Errorf("ConversionStatus mismatch: got %s, want %s", conversion.ConversionStatus, ConversionStatusFundsArrived)
	}

	// Settle time should be set
	if conversion.SettleTime == "" {
		t.Error("SettleTime should not be empty for funds arrived status")
	}
}

// ============================================================================
// All Conversion Event Types Test
// ============================================================================

func TestConversionAllEventTypes(t *testing.T) {
	testCases := []struct {
		name           string
		json           string
		expectedType   string
		expectedStatus string
		hasSettleTime  bool
	}{
		{
			name:           "TradeSettled",
			json:           conversionTradeSettledWebhookJSON,
			expectedType:   EventTypeConversionTradeSettled,
			expectedStatus: ConversionStatusTradeSettled,
			hasSettleTime:  true,
		},
		{
			name:           "FundsAwaiting",
			json:           conversionFundsAwaitingWebhookJSON,
			expectedType:   EventTypeConversionFundsAwaiting,
			expectedStatus: ConversionStatusAwaitingFunds,
			hasSettleTime:  false,
		},
		{
			name:           "FundsArrived",
			json:           conversionFundsArrivedWebhookJSON,
			expectedType:   EventTypeConversionFundsArrived,
			expectedStatus: ConversionStatusFundsArrived,
			hasSettleTime:  true,
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

			if !event.IsConversionEvent() {
				t.Error("IsConversionEvent should return true")
			}

			conversion, err := event.ParseConversionData()
			if err != nil {
				t.Fatalf("Failed to parse: %v", err)
			}

			if conversion.ConversionStatus != tc.expectedStatus {
				t.Errorf("ConversionStatus mismatch: got %s, want %s", conversion.ConversionStatus, tc.expectedStatus)
			}

			if tc.hasSettleTime && conversion.SettleTime == "" {
				t.Error("SettleTime should not be empty")
			}
			if !tc.hasSettleTime && conversion.SettleTime != "" {
				t.Error("SettleTime should be empty")
			}
		})
	}
}

// ============================================================================
// Conversion Constants Tests
// ============================================================================

func TestConversionStatusConstants(t *testing.T) {
	if ConversionStatusTradeSettled != "TRADE_SETTLED" {
		t.Errorf("ConversionStatusTradeSettled mismatch")
	}
	if ConversionStatusAwaitingFunds != "AWAITING_FUNDS" {
		t.Errorf("ConversionStatusAwaitingFunds mismatch")
	}
	if ConversionStatusFundsArrived != "FUNDS_ARRIVED" {
		t.Errorf("ConversionStatusFundsArrived mismatch")
	}
	if ConversionStatusPending != "PENDING" {
		t.Errorf("ConversionStatusPending mismatch")
	}
	if ConversionStatusCompleted != "COMPLETED" {
		t.Errorf("ConversionStatusCompleted mismatch")
	}
	if ConversionStatusCanceled != "CANCELED" {
		t.Errorf("ConversionStatusCanceled mismatch")
	}
	if ConversionStatusFailed != "FAILED" {
		t.Errorf("ConversionStatusFailed mismatch")
	}
}

func TestConversionWayConstants(t *testing.T) {
	if ConversionWayAPI != "API" {
		t.Errorf("ConversionWayAPI mismatch")
	}
	if ConversionWayWeb != "WEB" {
		t.Errorf("ConversionWayWeb mismatch")
	}
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestParseConversionData_WrongEventType(t *testing.T) {
	wrongTypeJSON := `{
		"version": "V1.6.0",
		"event_name": "ACQUIRING",
		"event_type": "acquiring.payment_intent.created",
		"event_id": "test-id",
		"source_id": "test-source",
		"data": {}
	}`

	var event Event
	json.Unmarshal([]byte(wrongTypeJSON), &event)

	_, err := event.ParseConversionData()
	if err == nil {
		t.Error("ParseConversionData should fail for non-conversion event type")
	}
}
