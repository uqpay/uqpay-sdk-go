package banking

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/common"
	"github.com/uqpay/uqpay-sdk-go/configuration"
)

type vaStaticTokenProvider struct{}

func (*vaStaticTokenProvider) GetToken() (string, error) { return "token", nil }

func newVATestClient(t *testing.T, handler http.HandlerFunc) (*Client, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	config := &configuration.Configuration{Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}
	return NewClient(common.NewAPIClient(config, &vaStaticTokenProvider{})), server.Close
}

func TestCreateVirtualAccountApplicationContract(t *testing.T) {
	requestCount := 0
	client, closeServer := newVATestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if r.Method != http.MethodPost || r.URL.Path != "/v1/virtual/accounts" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("x-idempotency-key"); got != "va-replay-key" {
			t.Errorf("idempotency header = %q", got)
		}
		if got := r.Header.Get("x-on-behalf-of"); got != "account-id" {
			t.Errorf("on-behalf header = %q", got)
		}
		var body CreateVirtualAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Country != "BH" || body.Currency != "USD" || body.PaymentMethod != "SWIFT" {
			t.Fatalf("unexpected body: %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"application_id":"app-id","public_version":1,"country":"BH","currency":"USD","status":"SUBMITTED","results":[{"payment_method":"SWIFT","status":"SUBMITTED","virtual_accounts":[],"error":null}]}}`))
	})
	defer closeServer()

	req := &CreateVirtualAccountRequest{Country: "BH", Currency: "USD", PaymentMethod: "SWIFT"}
	opts := &common.RequestOptions{IdempotencyKey: "va-replay-key", OnBehalfOf: "account-id"}
	resp, err := client.VirtualAccounts.Create(context.Background(), req, opts)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data.ApplicationID != "app-id" || resp.Data.PublicVersion != 1 || len(resp.Data.Results) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	replay, err := client.VirtualAccounts.Create(context.Background(), req, opts)
	if err != nil || replay.Data.ApplicationID != resp.Data.ApplicationID || requestCount != 2 {
		t.Fatalf("legal replay must reuse the application: %+v, %v, count=%d", replay, err, requestCount)
	}
}

func TestApplicationListRetrieveAndStrictNotFoundError(t *testing.T) {
	client, closeServer := newVATestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/virtual/applications" {
			if r.URL.Query().Get("page_number") != "1" || r.URL.Query().Get("page_size") != "50" || r.URL.Query().Get("status") != "SUBMITTED" {
				t.Fatalf("unexpected query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"total_pages":1,"total_items":1,"data":[{"application_id":"app-id","public_version":1,"country":"BH","currency":"USD","status":"SUBMITTED","created_at":"2026-08-12T00:00:00Z"}]}`))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"type":"not_found","code":"virtual_account_application_not_found","message":"Virtual account application not found"}`))
	})
	defer closeServer()
	list, err := client.VirtualAccountApplications.List(context.Background(), &ListVirtualAccountApplicationsRequest{PageNumber: 1, PageSize: 50, Status: "SUBMITTED"})
	if err != nil || list.TotalItems != 1 || list.Data[0].ApplicationID != "app-id" {
		t.Fatalf("list: %+v, %v", list, err)
	}
	_, err = client.VirtualAccountApplications.Retrieve(context.Background(), "missing")
	apiErr, ok := err.(*common.APIError)
	if !ok { // resource wraps with %w
		var wrapped *common.APIError
		if !errors.As(err, &wrapped) {
			t.Fatalf("expected APIError, got %T: %v", err, err)
		}
		apiErr = wrapped
	}
	if apiErr.StatusCode != 400 || apiErr.Type != "not_found" || string(apiErr.Code) != "virtual_account_application_not_found" || apiErr.Message != "Virtual account application not found" {
		t.Fatalf("unexpected error: %+v", apiErr)
	}
}
