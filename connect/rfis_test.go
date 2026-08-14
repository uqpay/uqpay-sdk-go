package connect

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v2/common"
	"github.com/uqpay/uqpay-sdk-go/v2/configuration"
)

type capturedRFIRequest struct {
	method string
	path   string
	query  string
	body   string
}

func TestRFIRoutesMatchPublishedContract(t *testing.T) {
	requests := make(chan capturedRFIRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		requests <- capturedRFIRequest{r.Method, r.URL.Path, r.URL.RawQuery, string(body)}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	config := &configuration.Configuration{
		Environment: &configuration.Environment{BaseURL: server.URL},
		HTTPClient:  server.Client(),
	}
	client := NewClient(common.NewAPIClient(config, &requestOptionsTokenProvider{}))
	ctx := context.Background()

	if _, err := client.RFIs.List(ctx, &ListRFIsRequest{
		PageSize: 10, PageNumber: 1, Status: RFIStatusActionRequired,
	}); err != nil {
		t.Fatalf("RFIs.List returned an error: %v", err)
	}
	got := <-requests
	if got.method != http.MethodGet || got.path != "/v1/rfis" {
		t.Errorf("request = %s %s, want GET /v1/rfis", got.method, got.path)
	}
	for _, want := range []string{"page_size=10", "page_number=1", "status=ACTION_REQUIRED"} {
		if !strings.Contains(got.query, want) {
			t.Errorf("query = %q, want it to contain %q", got.query, want)
		}
	}

	if _, err := client.RFIs.Get(ctx, "rfi_123"); err != nil {
		t.Fatalf("RFIs.Get returned an error: %v", err)
	}
	got = <-requests
	if got.method != http.MethodGet || got.path != "/v1/rfis/rfi_123" {
		t.Errorf("request = %s %s, want GET /v1/rfis/rfi_123", got.method, got.path)
	}

	if _, err := client.RFIs.Answer(ctx, &AnswerRFIRequest{RFIID: "rfi_123"}); err != nil {
		t.Fatalf("RFIs.Answer returned an error: %v", err)
	}
	got = <-requests
	if got.method != http.MethodPost || got.path != "/v1/rfis/answer" {
		t.Errorf("request = %s %s, want POST /v1/rfis/answer", got.method, got.path)
	}
	if !strings.Contains(got.body, `"rfi_id":"rfi_123"`) {
		t.Errorf("body = %q, want rfi_id", got.body)
	}
}
