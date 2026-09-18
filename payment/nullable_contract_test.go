package payment

import (
	"context"
	"encoding/json"
	"github.com/uqpay/uqpay-sdk-go/v4/common"
	"github.com/uqpay/uqpay-sdk-go/v4/configuration"
	"net/http"
	"net/http/httptest"
	"testing"
)

// AQ-RESPONSE: empty strings need a populated control to detect dropped fields.
func TestAcquiringOptionalResponsesThroughTransport(t *testing.T) {
	payload := `{}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(payload))
	}))
	defer server.Close()
	api := common.NewAPIClient(&configuration.Configuration{Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}, &staticTokenProvider{token: "offline"})
	client := NewClient(api)
	for _, value := range []string{"", "01"} {
		cvv, stamp := "", ""
		if value != "" {
			cvv = "M"
			stamp = "2026-09-17T00:00:00Z"
		}
		raw, _ := json.Marshal(map[string]interface{}{"advice_code": value, "authentication_data": map[string]interface{}{"cvv_result": cvv}, "complete_time": stamp})
		payload = string(raw)
		got, err := client.PaymentAttempts.Get(context.Background(), "pa-1")
		if err != nil {
			t.Fatal(err)
		}
		if got.AdviceCode != value || got.AuthenticationData["cvv_result"] != cvv || got.CompleteTime != stamp {
			t.Fatal(got)
		}
	}
	for _, raw := range []string{`{}`, `{"metadata":null}`, `{"metadata":{}}`, `{"metadata":{"ref":"0001"}}`} {
		payload = raw
		got, err := client.Refunds.Get(context.Background(), "re-1")
		if err != nil {
			t.Fatal(err)
		}
		var expected struct{ Metadata map[string]string }
		_ = json.Unmarshal([]byte(raw), &expected)
		a, _ := json.Marshal(got.Metadata)
		b, _ := json.Marshal(expected.Metadata)
		if string(a) != string(b) {
			t.Fatalf("%s != %s", a, b)
		}
	}
	for _, raw := range []string{`{}`, `{"completed_time":""}`, `{"completed_time":"2026-09-17T00:00:00Z"}`} {
		payload = raw
		got, err := client.Payouts.Get(context.Background(), "po-1")
		if err != nil {
			t.Fatal(err)
		}
		var expected struct {
			CompletedTime string `json:"completed_time"`
		}
		_ = json.Unmarshal([]byte(raw), &expected)
		if got.CompletedTime != expected.CompletedTime {
			t.Fatal(got)
		}
	}
}
