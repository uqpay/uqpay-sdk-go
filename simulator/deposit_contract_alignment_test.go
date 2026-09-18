package simulator

import (
	"encoding/json"
	"testing"
)

func TestDepositAccountIDContract(t *testing.T) {
	b, err := json.Marshal(CreateDepositRequest{AccountID: "account-1", Amount: 10, Currency: "SGD", SenderSwiftCode: "WELGBE22"})
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]interface{}
	_ = json.Unmarshal(b, &wire)
	if wire["account_id"] != "account-1" || wire["amount"] != float64(10) {
		t.Fatalf("bad request: %s", b)
	}
}
