package banking

import (
	"context"
	"encoding/json"
	"github.com/uqpay/uqpay-sdk-go/v3/common"
	"net/http"
	"testing"
)

func TestBankingAPISyncValidationAndMethodRules(t *testing.T) {
	client, closeServer := newVATestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.Header.Get("x-on-behalf-of") != "sub-account" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/payouts/validation":
			q := r.URL.Query()
			if q.Get("currency") != "USD" || q.Get("amount") != "1.00" || q.Get("beneficiary_id") != "saved + id" || q.Get("fee_paid_by") != "OURS" {
				t.Errorf("query: %v", q)
			}
			_, _ = w.Write([]byte(`{"currency":"USD","amount":"1.00","fee_amount":"5","total_amount":"6","channel_code":"MASTERCARD","sufficient_funds":false,"beneficiary_valid":true}`))
		case "/v1/beneficiaries/paymentmethods":
			if r.URL.Query().Get("entity_type") != "INDIVIDUAL" {
				t.Errorf("missing entity type")
			}
			_, _ = w.Write([]byte(`{"data":[{"payment_method":"CARD","clearing_systems":"MASTERCARD SEND","validation_field":["account_number",{"value":"card_expiry","rule":{"pattern":"month"},"children":[{"type":"MONTH"}]}]}]}`))
		case "/v1/beneficiaries":
			if r.URL.Query().Get("payment_method") != "CARD" {
				t.Errorf("missing filter")
			}
			_, _ = w.Write([]byte(`{"data":[]}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	})
	defer closeServer()
	opts := &common.RequestOptions{OnBehalfOf: "sub-account"}
	result, err := client.Payouts.Validate(context.Background(), &ValidatePayoutRequest{BeneficiaryID: "saved + id", Currency: "USD", Amount: "1.00", FeePaidBy: "OURS"}, opts)
	if err != nil || result.ChannelCode != "MASTERCARD" || result.FeeAmount != "5" || result.TotalAmount != "6" || result.SufficientFunds || !result.BeneficiaryValid {
		t.Fatalf("result %+v, %v", result, err)
	}
	methods, err := client.Beneficiaries.ListPaymentMethodsForEntity(context.Background(), "USD", "US", "INDIVIDUAL", opts)
	if err != nil || len(methods) != 1 || len(methods[0].ValidationField) != 2 {
		t.Fatalf("methods %+v, %v", methods, err)
	}
	rule := methods[0].ValidationField[1].(map[string]any)
	if rule["value"] != "card_expiry" {
		t.Fatalf("rule lost: %+v", rule)
	}
	_, err = client.Beneficiaries.List(context.Background(), &ListBeneficiariesRequest{PageSize: 10, PageNumber: 1, PaymentMethod: "CARD"}, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBankingAPISyncCardFieldsAndOptionalPayout(t *testing.T) {
	raw, err := json.Marshal(BeneficiaryCreationRequest{EntityType: "INDIVIDUAL", PaymentMethod: "CARD", BankDetails: &BankDetails{CardExpiry: "2099-12", AccountHolder: "Test Recipient", AccountNumber: "1234567890123", AccountCurrencyCode: "USD", BankCountryCode: "US", ClearingSystem: "MASTERCARD SEND"}})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	_ = json.Unmarshal(raw, &body)
	if _, ok := body["address"]; ok {
		t.Fatalf("individual CARD address should be omitted: %s", raw)
	}
	if body["bank_details"].(map[string]any)["card_expiry"] != "2099-12" {
		t.Fatal("expiry lost")
	}
	raw, _ = json.Marshal(CreatePayoutRequest{Currency: "USD", Amount: "1", PurposeCode: "BILL_PAYMENT"})
	_ = json.Unmarshal(raw, &body)
	var payout map[string]any
	_ = json.Unmarshal(raw, &payout)
	for _, key := range []string{"payout_reference", "fee_paid_by", "payout_date"} {
		if _, ok := payout[key]; ok {
			t.Errorf("optional field %s emitted: %s", key, raw)
		}
	}
}
