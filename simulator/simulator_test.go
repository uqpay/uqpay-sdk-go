package simulator

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

type staticTokenProvider struct{}

func (*staticTokenProvider) GetToken() (string, error) { return "token_123", nil }

func TestSimulatorRoutesMatchPublishedContract(t *testing.T) {
	type captured struct{ method, path, body string }
	requests := make(chan captured, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		requests <- captured{r.Method, r.URL.Path, string(body)}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	config := &configuration.Configuration{
		Environment: &configuration.Environment{BaseURL: server.URL},
		HTTPClient:  server.Client(),
	}
	client := NewClient(common.NewAPIClient(config, &staticTokenProvider{}))
	ctx := context.Background()

	tests := []struct {
		path string
		body string
		call func() error
	}{
		{
			"/v1/simulation/issuing/authorization", `"card_id":"card_123"`,
			func() error {
				_, err := client.Issuing.Authorize(ctx, &AuthorizationRequest{CardID: "card_123"})
				return err
			},
		},
		{
			"/v1/simulation/issuing/reversal", `"transaction_id":"txn_123"`,
			func() error {
				_, err := client.Issuing.Reverse(ctx, &ReversalRequest{TransactionID: "txn_123"})
				return err
			},
		},
		{
			"/v1/simulation/deposit", `"sender_swift_code":"TESTSGSG"`,
			func() error {
				_, err := client.Deposits.Create(ctx, &CreateDepositRequest{SenderSwiftCode: "TESTSGSG"})
				return err
			},
		},
	}

	for _, tt := range tests {
		if err := tt.call(); err != nil {
			t.Fatalf("%s returned an error: %v", tt.path, err)
		}
		got := <-requests
		if got.method != http.MethodPost || got.path != tt.path {
			t.Errorf("request = %s %s, want POST %s", got.method, got.path, tt.path)
		}
		if !strings.Contains(got.body, tt.body) {
			t.Errorf("body = %q, want it to contain %q", got.body, tt.body)
		}
	}
}
