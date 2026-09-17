package banking

import (
	"encoding/json"
	"testing"
)

func TestBeneficiaryIBANAndDepositFields(t *testing.T) {
	req := BeneficiaryCheckRequest{EntityType: "COMPANY", PaymentMethod: "LOCAL", Currency: "EUR", IBAN: "DE89370400440532013000", BankCountryCode: "DE"}
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]interface{}
	_ = json.Unmarshal(raw, &wire)
	if _, ok := wire["account_number"]; ok {
		t.Fatalf("unexpected empty account: %s", raw)
	}
	if wire["iban"] != req.IBAN || wire["bank_country_code"] != "DE" {
		t.Fatal(string(raw))
	}
	var deposit Deposit
	if err = json.Unmarshal([]byte(`{"deposit_method":"UQPAY_TRANSFER","amount":"-12345678901234567890.12345678","complete_time":null,"sender":{"sender_type":"COMPANY","name_type":"NAMED"}}`), &deposit); err != nil {
		t.Fatal(err)
	}
	if deposit.DepositMethod != "UQPAY_TRANSFER" || deposit.Sender.NameType != "NAMED" || deposit.Sender.SenderType != "COMPANY" || deposit.Amount != "-12345678901234567890.12345678" {
		t.Fatalf("lost response fields: %+v", deposit)
	}
}
