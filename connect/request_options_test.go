package connect

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v2/common"
	"github.com/uqpay/uqpay-sdk-go/v2/configuration"
)

type requestOptionsTokenProvider struct{}

func (*requestOptionsTokenProvider) GetToken() (string, error) {
	return "default-token", nil
}

func TestAccountsGetWithOptionsPreservesBusinessCode(t *testing.T) {
	requests := make(chan *http.Request, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	config := &configuration.Configuration{
		Environment: &configuration.Environment{BaseURL: server.URL},
		HTTPClient:  server.Client(),
	}
	apiClient := common.NewAPIClient(config, &requestOptionsTokenProvider{})
	client := NewClient(apiClient)

	if _, err := client.Accounts.GetWithOptions(
		context.Background(),
		"account_123",
		&common.RequestOptions{OnBehalfOf: "account_sub_123"},
		"BANKING",
	); err != nil {
		t.Fatalf("Accounts.GetWithOptions returned an error: %v", err)
	}

	got := <-requests
	if got.URL.Path != "/v1/accounts/account_123" {
		t.Errorf("request path = %q, want %q", got.URL.Path, "/v1/accounts/account_123")
	}
	if got.URL.Query().Get("business_code") != "BANKING" {
		t.Errorf("business_code = %q, want %q", got.URL.Query().Get("business_code"), "BANKING")
	}
	if gotValue := got.Header.Get("x-on-behalf-of"); gotValue != "account_sub_123" {
		t.Errorf("x-on-behalf-of = %q, want %q", gotValue, "account_sub_123")
	}
}
