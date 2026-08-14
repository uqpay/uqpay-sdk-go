package payment

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v2/common"
	"github.com/uqpay/uqpay-sdk-go/v2/configuration"
)

type staticTokenProvider struct {
	token string
}

func (p *staticTokenProvider) GetToken() (string, error) {
	return p.token, nil
}

func TestPaymentIntentOptionsMergeDefaultClientID(t *testing.T) {
	requests := make(chan http.Header, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	config := &configuration.Configuration{
		ClientID:    "configured-client",
		Environment: &configuration.Environment{BaseURL: server.URL},
		HTTPClient:  server.Client(),
	}
	apiClient := common.NewAPIClient(config, &staticTokenProvider{token: "default-token"})
	client := NewClient(apiClient)
	opts := &common.RequestOptions{
		OnBehalfOf:     "account_sub_456",
		IdempotencyKey: "idempotency_456",
		AuthToken:      "request-token",
	}

	if _, err := client.PaymentIntents.Get(context.Background(), "intent_123", opts); err != nil {
		t.Fatalf("PaymentIntents.Get returned an error: %v", err)
	}

	got := <-requests
	wantHeaders := map[string]string{
		"x-on-behalf-of":    "account_sub_456",
		"x-client-id":       "configured-client",
		"x-idempotency-key": "idempotency_456",
		"x-auth-token":      "Bearer request-token",
	}
	for name, want := range wantHeaders {
		if gotValue := got.Get(name); gotValue != want {
			t.Errorf("%s = %q, want %q", name, gotValue, want)
		}
	}
	if opts.ClientID != "" {
		t.Errorf("caller RequestOptions.ClientID mutated to %q, want empty", opts.ClientID)
	}
}

func TestPaymentIntentWithoutOptionsPreservesConfiguredClientID(t *testing.T) {
	requests := make(chan http.Header, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	config := &configuration.Configuration{
		ClientID:    "configured-client",
		Environment: &configuration.Environment{BaseURL: server.URL},
		HTTPClient:  server.Client(),
	}
	apiClient := common.NewAPIClient(config, &staticTokenProvider{token: "default-token"})
	client := NewClient(apiClient)

	if _, err := client.PaymentIntents.Get(context.Background(), "intent_123"); err != nil {
		t.Fatalf("PaymentIntents.Get returned an error: %v", err)
	}

	got := <-requests
	if gotValue := got.Get("x-client-id"); gotValue != "configured-client" {
		t.Errorf("x-client-id = %q, want %q", gotValue, "configured-client")
	}
	if gotValue := got.Get("x-auth-token"); gotValue != "Bearer default-token" {
		t.Errorf("x-auth-token = %q, want %q", gotValue, "Bearer default-token")
	}
	if gotValue := got.Get("x-on-behalf-of"); gotValue != "" {
		t.Errorf("x-on-behalf-of = %q, want empty", gotValue)
	}
}
