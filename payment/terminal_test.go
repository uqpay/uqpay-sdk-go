package payment

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

type capturedTerminalRequest struct {
	method string
	path   string
	header http.Header
	body   string
}

func TestTerminalRoutesIncludeConfiguredClientID(t *testing.T) {
	requests := make(chan capturedTerminalRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		requests <- capturedTerminalRequest{r.Method, r.URL.Path, r.Header.Clone(), string(body)}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	config := &configuration.Configuration{
		ClientID:    "client_123",
		Environment: &configuration.Environment{BaseURL: server.URL},
		HTTPClient:  server.Client(),
	}
	client := NewClient(common.NewAPIClient(config, &staticTokenProvider{token: "token_123"}))
	ctx := context.Background()

	if _, err := client.Terminals.Register(ctx, &RegisterTerminalRequest{
		FirmCode: "01", FirmSN: "SN123", TerminalModel: "PAX A920",
	}); err != nil {
		t.Fatalf("Terminals.Register returned an error: %v", err)
	}
	got := <-requests
	assertTerminalRequest(t, got, "/v2/terminal/register", `"firm_code":"01"`)

	if _, err := client.Terminals.GetPINKey(ctx, &GetPINKeyRequest{
		TerminalID: "terminal_123", PrivateKey: "key",
	}); err != nil {
		t.Fatalf("Terminals.GetPINKey returned an error: %v", err)
	}
	got = <-requests
	assertTerminalRequest(t, got, "/v2/terminal/getPinKey", `"prv_key":"key"`)
}

func assertTerminalRequest(t *testing.T, got capturedTerminalRequest, path, bodyFragment string) {
	t.Helper()
	if got.method != http.MethodPost || got.path != path {
		t.Errorf("request = %s %s, want POST %s", got.method, got.path, path)
	}
	if got.header.Get("x-client-id") != "client_123" {
		t.Errorf("x-client-id = %q, want client_123", got.header.Get("x-client-id"))
	}
	if !strings.Contains(got.body, bodyFragment) {
		t.Errorf("body = %q, want it to contain %q", got.body, bodyFragment)
	}
}
