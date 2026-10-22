package issuing_test

import (
	"context"
	"encoding/json"
	"fmt"
	uqpay "github.com/uqpay/uqpay-sdk-go/v3"
	"github.com/uqpay/uqpay-sdk-go/v3/configuration"
	"github.com/uqpay/uqpay-sdk-go/v3/issuing"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPublicAssignMetadata(t *testing.T) {
	var body map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/connect/token" {
			fmt.Fprint(w, `{"auth_token":"offline-token","expired_at":4102444800}`)
			return
		}
		if r.Method != "POST" || r.URL.Path != "/v1/issuing/cards/assign" {
			t.Errorf("route %s %s", r.Method, r.URL.Path)
		}
		body = nil
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, `{"card_id":"card-1"}`)
	}))
	defer server.Close()
	client, err := uqpay.NewClient("offline-client", "offline-key", &configuration.Environment{BaseURL: server.URL, FilesBaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	for _, metadata := range []*issuing.AssignCardMetadata{nil, {CardSKUCode: "SKU123456789"}} {
		result, err := client.Issuing.Cards.Assign(context.Background(), &issuing.AssignCardRequest{CardholderID: "holder-1", CardNumber: "example-pan", CardCurrency: "USD", CardMode: "SINGLE", Metadata: metadata})
		if err != nil {
			t.Fatal(err)
		}
		if result.CardID != "card-1" {
			t.Fatal(result)
		}
		expected := map[string]interface{}{"cardholder_id": "holder-1", "card_number": "example-pan", "card_currency": "USD", "card_mode": "SINGLE"}
		if metadata != nil {
			expected["metadata"] = map[string]interface{}{"card_sku_code": "SKU123456789"}
		}
		if !reflect.DeepEqual(body, expected) {
			t.Fatalf("wire %#v", body)
		}
	}
}
