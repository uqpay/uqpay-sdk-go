package issuing

import (
	"context"
	"encoding/json"
	"github.com/uqpay/uqpay-sdk-go/v3/common"
	"github.com/uqpay/uqpay-sdk-go/v3/configuration"
	"github.com/uqpay/uqpay-sdk-go/v3/payment"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestKYCBranchesPaginationAndProxyHeaders(t *testing.T) {
	var body map[string]interface{}
	var path, size string
	var headers http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		size = r.URL.Query().Get("page_size")
		headers = r.Header.Clone()
		body = nil
		if r.Body != nil && r.Method == "POST" {
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Error(err)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()
	api := common.NewAPIClient(&configuration.Configuration{ClientID: "client", Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}, &staticTokenProvider{token: "offline"})
	client := NewClient(api)
	ctx := context.Background()
	for _, provider := range []string{"SUMSUB", "MYINFO", "JUMIO", "DIDIT", "SHUFTI", "REGTANK"} {
		for _, n := range []int{9, 10, 64, 65} {
			for _, dob := range []string{"2009-09-17", "2008-09-17", "1947-09-17", "1946-09-17"} {
				ref := strings.Repeat("r", n)
				proof := &KycVerification{Method: "THIRD_PARTY", KycProof: &KycProof{Provider: provider, ReferenceID: ref}}
				_, err := client.Cardholders.Create(ctx, &CreateCardholderRequest{Email: "test@example.test", FirstName: "Test", LastName: "User", CountryCode: "SG", PhoneNumber: "81234567", DateOfBirth: &dob, KycVerification: proof})
				if err != nil {
					t.Fatal(err)
				}
				checkKYCBody(t, body, ref, dob)
				if path != "/v1/issuing/cardholders" {
					t.Fatal(path)
				}
				_, err = client.Cardholders.Update(ctx, "holder-1", &UpdateCardholderRequest{DateOfBirth: &dob, KycVerification: proof})
				if err != nil {
					t.Fatal(err)
				}
				checkKYCBody(t, body, ref, dob)
				if path != "/v1/issuing/cardholders/holder-1" {
					t.Fatal(path)
				}
				email, first, last, country := "test@example.test", "Test", "User", "SG"
				_, err = client.Cards.Create(ctx, &CreateCardRequest{CardholderID: "holder-1", CardCurrency: "SGD", CardProductID: "product-1", CardholderRequiredFields: &CardholderRequiredFields{Email: &email, FirstName: &first, LastName: &last, CountryCode: &country, DateOfBirth: &dob, KycVerification: proof}})
				if err != nil {
					t.Fatal(err)
				}
				inline := body["cardholder_required_fields"].(map[string]interface{})
				checkKYCBody(t, inline, ref, dob)
				if inline["email"] != email || inline["country_code"] != country || inline["first_name"] != first || inline["last_name"] != last {
					t.Fatal(inline)
				}
			}
		}
	}
	for _, n := range []int{1, 10, 100} {
		_, err := client.Cards.List(ctx, &ListCardsRequest{PageSize: n, PageNumber: 1})
		if err != nil {
			t.Fatal(err)
		}
		if size != strconv.Itoa(n) {
			t.Fatal(size)
		}
	}
	payments := payment.NewClient(api)
	// D189-D196: route-specific delegation without caller idempotency keys.
	for _, tc := range []struct {
		path string
		call func(*common.RequestOptions) error
	}{
		{"/v2/payment/balances", func(o *common.RequestOptions) error {
			_, e := payments.Balances.List(ctx, &payment.ListBalancesRequest{}, o)
			return e
		}},
		{"/v2/payment/balances/USD", func(o *common.RequestOptions) error { _, e := payments.Balances.Get(ctx, "USD", o); return e }},
		{"/v2/payment/bankaccount", func(o *common.RequestOptions) error {
			_, e := payments.BankAccounts.List(ctx, &payment.ListBankAccountsRequest{}, o)
			return e
		}},
		{"/v2/payment/bankaccount/ba-1", func(o *common.RequestOptions) error { _, e := payments.BankAccounts.Get(ctx, "ba-1", o); return e }},
		{"/v2/payment/payout", func(o *common.RequestOptions) error {
			_, e := payments.Payouts.List(ctx, &payment.ListPayoutsRequest{}, o)
			return e
		}},
		{"/v2/payment/payout/po-1", func(o *common.RequestOptions) error { _, e := payments.Payouts.Get(ctx, "po-1", o); return e }},
		{"/v2/payment/settlements", func(o *common.RequestOptions) error {
			_, e := payments.Reports.ListSettlements(ctx, &payment.ListSettlementsRequest{}, o)
			return e
		}},
		{"/v2/payment_intents/pi-1", func(o *common.RequestOptions) error { _, e := payments.PaymentIntents.Get(ctx, "pi-1", o); return e }},
	} {
		for _, account := range []string{"sub-account", ""} {
			var opts *common.RequestOptions
			if account != "" {
				opts = &common.RequestOptions{OnBehalfOf: account}
			}
			if err := tc.call(opts); err != nil {
				t.Fatal(err)
			}
			if path != tc.path || headers.Get("x-on-behalf-of") != account || headers.Get("x-client-id") != "client" {
				t.Fatalf("%s: %s %v", tc.path, path, headers)
			}
		}
	}
	_, err := payments.PaymentIntents.Create(ctx, &payment.CreatePaymentIntentRequest{Amount: "1.00", Currency: "USD"}, &common.RequestOptions{IdempotencyKey: "fixed-key"})
	if err != nil {
		t.Fatal(err)
	}
	if headers.Get("x-idempotency-key") != "fixed-key" {
		t.Fatal(headers)
	}
}
func checkKYCBody(t *testing.T, body map[string]interface{}, ref, dob string) {
	t.Helper()
	proof := body["kyc_verification"].(map[string]interface{})["kyc_proof"].(map[string]interface{})
	if proof["reference_id"] != ref || body["date_of_birth"] != dob {
		t.Fatal(body)
	}
}
