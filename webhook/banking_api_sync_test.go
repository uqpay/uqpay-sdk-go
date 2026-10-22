package webhook

import "testing"

func TestBankingAPISyncCardExpiry(t *testing.T) {
	data := BeneficiaryData{PaymentType: "CARD", BeneficiaryBankDetailsRaw: `{"card_expiry":"2099-12","clearing_system":"Mastercard Send"}`}
	bank, err := data.GetBeneficiaryBankDetails()
	if err != nil || bank.CardExpiry != "2099-12" || bank.ClearingSystem != "Mastercard Send" {
		t.Fatalf("webhook expiry %+v %v", bank, err)
	}
}
