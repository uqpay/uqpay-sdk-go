package issuing

import (
	"context"
	"encoding/json"
	"github.com/uqpay/uqpay-sdk-go/v3/common"
	"github.com/uqpay/uqpay-sdk-go/v3/configuration"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestIssuingAPISync(t *testing.T) {
	var body map[string]interface{}
	expectedPath, expectedMethod := "/v1/issuing/cardholders/holder-1", http.MethodPost
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expectedMethod || r.URL.Path != expectedPath {
			t.Errorf("route: %s %s; want %s %s", r.Method, r.URL.Path, expectedMethod, expectedPath)
		}
		if r.Method == "POST" {
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"cardholder_id":"holder-1","residential_address":{"line2":null}}`))
	}))
	defer server.Close()
	api := common.NewAPIClient(&configuration.Configuration{ClientID: "client", Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}, &staticTokenProvider{token: "offline"})
	client := NewClient(api)
	poa, date, country, verification := "cGRm", "2026-10-01", "SG", "a3lj"
	identity := &Identity{ProofOfAddress: &poa, ProofOfAddressIssueDate: &date, ProofOfAddressIssuingCountry: &country, IdentityVerification: &verification}
	if _, err := client.Cardholders.Update(context.Background(), "holder-1", &UpdateCardholderRequest{Identity: identity}); err != nil {
		t.Fatal(err)
	}
	expected := map[string]interface{}{"identity": map[string]interface{}{"proof_of_address": poa, "proof_of_address_issue_date": date, "proof_of_address_issuing_country": country, "identity_verification": verification}}
	if !reflect.DeepEqual(body, expected) {
		t.Fatalf("wire: %#v", body)
	}
	expectedPath = "/v1/issuing/cards"
	if _, err := client.Cards.Create(context.Background(), &CreateCardRequest{CardCurrency: "USD", CardProductID: "product-1"}); err != nil {
		t.Fatal(err)
	}
	if _, ok := body["cardholder_id"]; ok {
		t.Fatalf("unexpected holder: %#v", body)
	}
	expectedPath = "/v1/issuing/cards/assign"
	if _, err := client.Cards.Assign(context.Background(), &AssignCardRequest{CardNumber: "example-pan", CardCurrency: "USD", CardMode: "SINGLE"}); err != nil {
		t.Fatal(err)
	}
	if _, ok := body["cardholder_id"]; ok {
		t.Fatalf("unexpected holder: %#v", body)
	}
	expectedPath, expectedMethod = "/v1/issuing/cardholders/holder-1", http.MethodGet
	holder, err := client.Cardholders.Get(context.Background(), "holder-1")
	if err != nil {
		t.Fatal(err)
	}
	if holder.ResidentialAddress == nil || holder.ResidentialAddress.Line2 != nil {
		t.Fatalf("sparse/null address: %#v", holder.ResidentialAddress)
	}
}
