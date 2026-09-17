package simulator

import (
	"encoding/json"
	"testing"
)

func TestAuthorizationResponseDecimalStrings(t *testing.T) {
	var response AuthorizationResponse
	if err := json.Unmarshal([]byte(`{"transaction_amount":"12345678901234567890.12345678","billing_amount":"0.00","card_available_balance":"-0.01"}`), &response); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(response)
	var got map[string]interface{}
	_ = json.Unmarshal(raw, &got)
	if got["transaction_amount"] != "12345678901234567890.12345678" || got["billing_amount"] != "0.00" || got["card_available_balance"] != "-0.01" {
		t.Fatal(string(raw))
	}
}
