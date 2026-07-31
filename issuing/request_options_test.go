package issuing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/uqpay/uqpay-sdk-go/common"
	"github.com/uqpay/uqpay-sdk-go/configuration"
)

type staticTokenProvider struct {
	token string
}

func (p *staticTokenProvider) GetToken() (string, error) {
	return p.token, nil
}

type capturedRequest struct {
	path   string
	header http.Header
}

func newRequestOptionsTestClient(t *testing.T) (*Client, <-chan capturedRequest) {
	t.Helper()

	requests := make(chan capturedRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- capturedRequest{
			path:   r.URL.Path,
			header: r.Header.Clone(),
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	config := &configuration.Configuration{
		Environment: &configuration.Environment{BaseURL: server.URL},
		HTTPClient:  server.Client(),
	}
	apiClient := common.NewAPIClient(config, &staticTokenProvider{token: "default-token"})
	return NewClient(apiClient), requests
}

func TestCardsGetForwardsRequestOptions(t *testing.T) {
	client, requests := newRequestOptionsTestClient(t)

	_, err := client.Cards.Get(context.Background(), "card_123", &common.RequestOptions{
		OnBehalfOf:     "account_sub_123",
		ClientID:       "client_override",
		IdempotencyKey: "idempotency_123",
		AuthToken:      "request-token",
	})
	if err != nil {
		t.Fatalf("Cards.Get returned an error: %v", err)
	}

	got := <-requests
	if got.path != "/v1/issuing/cards/card_123" {
		t.Fatalf("request path = %q, want %q", got.path, "/v1/issuing/cards/card_123")
	}

	wantHeaders := map[string]string{
		"x-on-behalf-of":    "account_sub_123",
		"x-client-id":       "client_override",
		"x-idempotency-key": "idempotency_123",
		"x-auth-token":      "Bearer request-token",
	}
	for name, want := range wantHeaders {
		if gotValue := got.header.Get(name); gotValue != want {
			t.Errorf("%s = %q, want %q", name, gotValue, want)
		}
	}
}

func TestCardsGetWithoutOptionsPreservesLegacyHeaders(t *testing.T) {
	client, requests := newRequestOptionsTestClient(t)

	if _, err := client.Cards.Get(context.Background(), "card_123"); err != nil {
		t.Fatalf("Cards.Get returned an error: %v", err)
	}

	got := <-requests
	if gotValue := got.header.Get("x-auth-token"); gotValue != "Bearer default-token" {
		t.Errorf("x-auth-token = %q, want %q", gotValue, "Bearer default-token")
	}
	if gotValue := got.header.Get("x-on-behalf-of"); gotValue != "" {
		t.Errorf("x-on-behalf-of = %q, want empty", gotValue)
	}
	if gotValue := got.header.Get("x-client-id"); gotValue != "" {
		t.Errorf("x-client-id = %q, want empty", gotValue)
	}
	if gotValue := got.header.Get("x-idempotency-key"); gotValue == "" {
		t.Error("x-idempotency-key is empty, want an automatically generated UUID")
	} else if _, err := uuid.Parse(gotValue); err != nil {
		t.Errorf("x-idempotency-key = %q, want a UUID: %v", gotValue, err)
	}
}

func TestDownloadCenterForwardsRequestOptions(t *testing.T) {
	client, requests := newRequestOptionsTestClient(t)

	resp, err := client.DownloadCenter.Download(context.Background(), "report_123", &common.RequestOptions{
		OnBehalfOf: "account_sub_123",
		AuthToken:  "request-token",
	})
	if err != nil {
		t.Fatalf("DownloadCenter.Download returned an error: %v", err)
	}
	if string(resp.Data) != "{}" {
		t.Errorf("downloaded data = %q, want %q", string(resp.Data), "{}")
	}

	got := <-requests
	if got.path != "/v1/issuing/reports/report_123" {
		t.Fatalf("request path = %q, want %q", got.path, "/v1/issuing/reports/report_123")
	}
	if gotValue := got.header.Get("x-on-behalf-of"); gotValue != "account_sub_123" {
		t.Errorf("x-on-behalf-of = %q, want %q", gotValue, "account_sub_123")
	}
}
