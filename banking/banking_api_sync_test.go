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
	if err != nil || result.ChannelCode != "MASTERCARD" || result.FeeAmount != "5" || result.TotalAmount != "6" || result.SufficientFunds == nil || *result.SufficientFunds || result.BeneficiaryValid == nil || !*result.BeneficiaryValid {
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

func TestBankingAPISyncSparseCardAddressPublicRequests(t *testing.T) {
	calls := 0
	client, closeServer := newVATestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != http.MethodPost {
			t.Errorf("method = %s", r.Method)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if r.URL.Path == "/v1/payouts" {
			body = body["beneficiary"].(map[string]any)
		}
		address := body["address"].(map[string]any)
		if len(address) != 1 || address["nationality"] != "US" {
			t.Errorf("sparse address changed: %+v", address)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"beneficiary_id":"ben-id","payout_id":"payout-id"}`))
	})
	defer closeServer()
	address := &Address{Nationality: "US"}
	bank := &BankDetails{CardExpiry: "2099-12", AccountHolder: "Test Recipient", AccountNumber: "1234567890123", AccountCurrencyCode: "USD", BankCountryCode: "US", ClearingSystem: "MASTERCARD SEND"}
	req := &BeneficiaryCreationRequest{EntityType: "INDIVIDUAL", PaymentMethod: "CARD", BankDetails: bank, Address: address}
	if _, err := client.Beneficiaries.Create(context.Background(), req); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Beneficiaries.Update(context.Background(), "ben-id", req); err != nil {
		t.Fatal(err)
	}
	payout := &CreatePayoutRequest{Currency: "USD", Amount: "1", PurposeCode: "BILL_PAYMENT", Beneficiary: &PayoutInlineBeneficiary{EntityType: "INDIVIDUAL", PaymentMethod: "CARD", BankDetails: bank, Address: address}}
	if _, err := client.Payouts.Create(context.Background(), payout); err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Errorf("calls = %d", calls)
	}
}

func TestBankingAPISyncValidationMissingVersusFalse(t *testing.T) {
	response := `{}`
	client, closeServer := newVATestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	})
	defer closeServer()
	req := &ValidatePayoutRequest{BeneficiaryID: "ben-id", Currency: "USD", Amount: "1", FeePaidBy: "OURS"}
	absent, err := client.Payouts.Validate(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if absent.SufficientFunds != nil || absent.BeneficiaryValid != nil {
		t.Fatalf("missing decision must remain unknown: %+v", absent)
	}
	raw, err := json.Marshal(absent)
	if err != nil {
		t.Fatal(err)
	}
	var encoded map[string]any
	if err := json.Unmarshal(raw, &encoded); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"sufficient_funds", "beneficiary_valid"} {
		if _, ok := encoded[key]; ok {
			t.Errorf("absent %s emitted", key)
		}
	}
	response = `{"sufficient_funds":false,"beneficiary_valid":false}`
	negative, err := client.Payouts.Validate(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if negative.SufficientFunds == nil || *negative.SufficientFunds || negative.BeneficiaryValid == nil || *negative.BeneficiaryValid {
		t.Fatalf("false decision lost: %+v", negative)
	}
	raw, err = json.Marshal(negative)
	if err != nil {
		t.Fatal(err)
	}
	encoded = nil
	if err := json.Unmarshal(raw, &encoded); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"sufficient_funds", "beneficiary_valid"} {
		if value, ok := encoded[key]; !ok || value != false {
			t.Errorf("explicit false %s lost: %s", key, raw)
		}
	}
}

func TestBankingAPISyncAddressResponseAndLegacyRequestSemantics(t *testing.T) {
	// Responses must retain empty strings, including the five frozen address fields.
	raw, err := json.Marshal(Beneficiary{Address: &Address{Nationality: "US"}})
	if err != nil {
		t.Fatal(err)
	}
	var response map[string]any
	if err := json.Unmarshal(raw, &response); err != nil {
		t.Fatal(err)
	}
	address := response["address"].(map[string]any)
	for _, key := range []string{"country", "city", "street_address", "state", "postal_code"} {
		if value, ok := address[key]; !ok || value != "" {
			t.Errorf("response empty %s lost: %s", key, raw)
		}
	}
	for _, paymentMethod := range []string{"LOCAL", "SWIFT", "FUTURE"} {
		raw, err = json.Marshal(BeneficiaryCreationRequest{PaymentMethod: paymentMethod, Address: &Address{Nationality: "US"}})
		if err != nil {
			t.Fatal(err)
		}
		var request map[string]any
		if err := json.Unmarshal(raw, &request); err != nil {
			t.Fatal(err)
		}
		legacy := request["address"].(map[string]any)
		if len(legacy) != 6 || legacy["country"] != "" {
			t.Errorf("%s request changed: %s", paymentMethod, raw)
		}
	}
	full := &Address{StreetAddress: "1 Main Street", City: "Boston", State: "MA", PostalCode: "02101", Country: "US", Nationality: "US"}
	for _, request := range []any{
		BeneficiaryCreationRequest{PaymentMethod: "CARD", Address: full},
		PayoutInlineBeneficiary{PaymentMethod: "CARD", Address: full},
	} {
		raw, err = json.Marshal(request)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]any
		if err := json.Unmarshal(raw, &fields); err != nil {
			t.Fatal(err)
		}
		address := fields["address"].(map[string]any)
		if len(address) != 6 || address["country"] != "US" || address["street_address"] != "1 Main Street" {
			t.Errorf("complete CARD address changed: %s", raw)
		}
	}
}
