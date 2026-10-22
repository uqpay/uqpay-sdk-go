package payment

import (
	"context"
	"encoding/json"
	"github.com/uqpay/uqpay-sdk-go/v3/common"
	"github.com/uqpay/uqpay-sdk-go/v3/configuration"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPaymentAPISync(t *testing.T) {
	var bodies []map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPaths := []string{"/v2/payment_intents/create", "/v2/payment_intents/pi-1/confirm", "/v2/payment_intents/pi-1/capture"}
		if len(bodies) >= len(expectedPaths) || r.Method != http.MethodPost || r.URL.Path != expectedPaths[len(bodies)] {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var body map[string]interface{}
		if r.Method == "POST" {
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
			bodies = append(bodies, body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"payment_intent_id":"pi-1","amount":"10.00","captured_amount":"6.00","intent_status":"PARTIALLY_CAPTURED","authorization_type":"pre_authorization","authorization_expires_at":null,"available_to_capture":"0.00","available_to_void":"0.00","capture_id":"cap-1","merchant_capture_order_id":"order-cap-1","capture_records":[{"capture_id":"cap-1","merchant_capture_order_id":null,"capture_sequence":1,"capture_amount":"6.00","capture_mode":"final","channel_capture_status":0,"channel_capture_time":null,"complete_time":null}],"current_funds_action":{"type":"release"}}`))
	}))
	defer server.Close()
	cfg := &configuration.Configuration{ClientID: "offline-client", Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}
	c := NewClient(common.NewAPIClient(cfg, &staticTokenProvider{token: "offline-token"}))
	no := false
	yes := true
	_, err := c.PaymentIntents.Create(context.Background(), &CreatePaymentIntentRequest{AuthorizationType: "pre_authorization", PaymentMethod: &PaymentMethod{Type: "card_present", CardPresent: &CardPresent{AuthorizationType: "pre_authorization", AutoCapture: &no}}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.PaymentIntents.Confirm(context.Background(), "pi-1", &ConfirmPaymentIntentRequest{AuthorizationType: "pre_authorization"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := c.PaymentIntents.Capture(context.Background(), "pi-1", &CapturePaymentIntentRequest{AmountToCapture: "6.00", MerchantCaptureOrderID: "order-cap-1", FinalCapture: &yes})
	if err != nil {
		t.Fatal(err)
	}
	if bodies[0]["authorization_type"] != "pre_authorization" || bodies[1]["authorization_type"] != "pre_authorization" {
		t.Fatal(bodies)
	}
	card := bodies[0]["payment_method"].(map[string]interface{})["card_present"].(map[string]interface{})
	if card["auto_capture"] != false || card["authorization_type"] != "pre_authorization" {
		t.Fatal(card)
	}
	if bodies[2]["amount_to_capture"] != "6.00" || bodies[2]["final_capture"] != true || bodies[2]["merchant_capture_order_id"] != "order-cap-1" {
		t.Fatal(bodies[2])
	}
	if result.CaptureID != "cap-1" || result.CapturedAmount != "6.00" || result.AuthorizationExpiresAt != nil || result.AvailableToCapture != "0.00" || result.CaptureRecords[0].ChannelCaptureStatus == nil || *result.CaptureRecords[0].ChannelCaptureStatus != 0 || result.CaptureRecords[0].CompleteTime != nil || result.CaptureRecords[0].MerchantCaptureOrderID != nil {
		t.Fatalf("lost capture response: %+v", result)
	}

	record := result.CaptureRecords[0]
	if record.CaptureSequence != 1 || record.CaptureAmount != "6.00" || record.CaptureMode != "final" || record.ChannelCaptureTime != nil || record.CompleteTime != nil || record.MerchantCaptureOrderID != nil {
		t.Fatalf("lost capture record fields: %+v", record)
	}
	b, _ := json.Marshal(&CapturePaymentIntentRequest{})
	if string(b) != "{}" {
		t.Fatalf("omitted capture must use defaults: %s", b)
	}
	var attempt PaymentAttempt
	if err := json.Unmarshal([]byte(`{"attempt_id":"a-1","authorization_type":"pre_authorization"}`), &attempt); err != nil || attempt.AuthorizationType != "pre_authorization" {
		t.Fatal(err)
	}
}

func TestPaymentConfirmOptionsAndLegacyCall(t *testing.T) {
	var headers []http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v2/payment_intents/pi-1/confirm" {
			t.Errorf("unexpected route %s %s", r.Method, r.URL.Path)
		}
		headers = append(headers, r.Header.Clone())
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"payment_intent_id":"pi-1"}`))
	}))
	defer server.Close()
	cfg := &configuration.Configuration{ClientID: "configured-client", Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}
	c := NewClient(common.NewAPIClient(cfg, &staticTokenProvider{token: "offline-token"}))
	req := &ConfirmPaymentIntentRequest{AuthorizationType: "pre_authorization"}
	// Existing three-argument callers retain the configured client ID.
	if _, err := c.PaymentIntents.Confirm(context.Background(), "pi-1", req); err != nil {
		t.Fatal(err)
	}
	opts := &common.RequestOptions{OnBehalfOf: "sub-account", IdempotencyKey: "caller-key"}
	if _, err := c.PaymentIntents.Confirm(context.Background(), "pi-1", req, opts); err != nil {
		t.Fatal(err)
	}
	override := &common.RequestOptions{ClientID: "override-client", OnBehalfOf: "other-account", IdempotencyKey: "override-key"}
	if _, err := c.PaymentIntents.Confirm(context.Background(), "pi-1", req, override); err != nil {
		t.Fatal(err)
	}
	if headers[0].Get("x-client-id") != "configured-client" || headers[0].Get("x-idempotency-key") == "" {
		t.Fatal(headers[0])
	}
	if headers[1].Get("x-client-id") != "configured-client" || headers[1].Get("x-on-behalf-of") != "sub-account" || headers[1].Get("x-idempotency-key") != "caller-key" {
		t.Fatal(headers[1])
	}
	if headers[2].Get("x-client-id") != "override-client" || headers[2].Get("x-on-behalf-of") != "other-account" || headers[2].Get("x-idempotency-key") != "override-key" {
		t.Fatal(headers[2])
	}
	if opts.ClientID != "" {
		t.Fatal("caller options mutated")
	}
}

func TestCaptureRecordStatusPresence(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want *int
	}{
		{raw: `{}`},
		{raw: `{"channel_capture_status":0}`, want: paymentStatusPointer(0)},
		{raw: `{"channel_capture_status":1}`, want: paymentStatusPointer(1)},
		{raw: `{"channel_capture_status":2}`, want: paymentStatusPointer(2)},
		{raw: `{"channel_capture_status":3}`, want: paymentStatusPointer(3)},
	} {
		var record CaptureRecord
		if err := json.Unmarshal([]byte(tc.raw), &record); err != nil {
			t.Fatal(err)
		}
		if tc.want == nil {
			if record.ChannelCaptureStatus != nil {
				t.Fatal("missing status must remain absent")
			}
		} else if record.ChannelCaptureStatus == nil || *record.ChannelCaptureStatus != *tc.want {
			t.Fatalf("status mismatch: %+v", record)
		}
		raw, err := json.Marshal(record)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]interface{}
		if err := json.Unmarshal(raw, &fields); err != nil {
			t.Fatal(err)
		}
		status, present := fields["channel_capture_status"]
		if tc.want == nil {
			if present {
				t.Fatal("marshal invented a status")
			}
		} else if !present || status != float64(*tc.want) {
			t.Fatalf("status lost: %s", raw)
		}
	}
}
func paymentStatusPointer(status int) *int { return &status }
