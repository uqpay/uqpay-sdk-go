package banking

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v4/common"
	"github.com/uqpay/uqpay-sdk-go/v4/configuration"
)

func TestFrozenDepositResponses(t *testing.T) {
	raw, err := os.ReadFile("testdata/deposit-contract.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name string
		Body json.RawMessage
	}
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range cases {
		t.Run(fixture.Name, func(t *testing.T) {
			requests := make(chan *http.Request, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests <- r
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(fixture.Body)
			}))
			defer server.Close()
			api := common.NewAPIClient(&configuration.Configuration{Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}, &vaStaticTokenProvider{})
			result, err := NewClient(api).Deposits.Get(context.Background(), "deposit-1")
			if err != nil {
				t.Fatal(err)
			}
			request := <-requests
			if request.Method != "GET" || request.URL.Path != "/v1/deposit/deposit-1" {
				t.Fatal(request.Method, request.URL)
			}
			var expected map[string]interface{}
			if err := json.Unmarshal(fixture.Body, &expected); err != nil {
				t.Fatal(err)
			}
			sender := expected["sender"].(map[string]interface{})
			if result.Sender == nil || result.DepositMethod != expected["deposit_method"] || result.Sender.SenderType != sender["sender_type"] || result.Sender.NameType != sender["name_type"] {
				t.Fatalf("classification lost: %+v sender=%+v", result, result.Sender)
			}
			// Public scalar fields keep existing empty/null semantics and omitempty serialization.
			if expected["deposit_method"] == "" {
				delete(expected, "deposit_method")
			}
			for _, key := range []string{"sender_type", "name_type"} {
				if sender[key] == "" {
					delete(sender, key)
				}
			}
			if expected["complete_time"] == nil {
				if result.CompleteTime != "" {
					t.Fatal("invented completion time")
				}
				expected["complete_time"] = ""
			}
			encoded, err := json.Marshal(result)
			if err != nil {
				t.Fatal(err)
			}
			var got interface{}
			if err := json.Unmarshal(encoded, &got); err != nil {
				t.Fatal(err)
			}
			assertBeneficiaryFields(t, fixture.Name, got, expected)
		})
	}
}
