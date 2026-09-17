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
	_, err := payments.PaymentIntents.Get(ctx, "pi-1", &common.RequestOptions{OnBehalfOf: "sub-account"})
	if err != nil {
		t.Fatal(err)
	}
	if headers.Get("x-on-behalf-of") != "sub-account" || path != "/v2/payment_intents/pi-1" {
		t.Fatal(headers, path)
	}
	_, err = payments.PaymentIntents.Create(ctx, &payment.CreatePaymentIntentRequest{Amount: "1.00", Currency: "USD"}, &common.RequestOptions{IdempotencyKey: "fixed-key"})
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
