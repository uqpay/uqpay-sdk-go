package issuing

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v4/banking"
	"github.com/uqpay/uqpay-sdk-go/v4/common"
	"github.com/uqpay/uqpay-sdk-go/v4/configuration"
	"github.com/uqpay/uqpay-sdk-go/v4/payment"
	"github.com/uqpay/uqpay-sdk-go/v4/simulator"
)

func TestRemainingFrozenResponses(t *testing.T) {
	raw, err := os.ReadFile("testdata/remaining-responses.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Operation, Name, Path, Method string
		Request, Body                 json.RawMessage
	}
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range cases {
		t.Run(fixture.Operation+"/"+fixture.Name, func(t *testing.T) {
			requests := make(chan *http.Request, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests <- r
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(fixture.Body)
			}))
			defer server.Close()
			api := common.NewAPIClient(&configuration.Configuration{Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}, &staticTokenProvider{token: "offline"})
			ctx := context.Background()
			pay := payment.NewClient(api)
			var actual interface{}
			decode := func(v interface{}) {
				t.Helper()
				if err := json.Unmarshal(fixture.Request, v); err != nil {
					t.Fatal(err)
				}
			}
			switch fixture.Operation {
			case "payout":
				actual, err = banking.NewClient(api).Payouts.Get(ctx, "po-1")
			case "transaction":
				actual, err = NewClient(api).Transactions.Get(ctx, "tx-1")
			case "authorization":
				var req simulator.AuthorizationRequest
				decode(&req)
				actual, err = simulator.NewClient(api).Issuing.Authorize(ctx, &req)
			case "bank.get":
				actual, err = pay.BankAccounts.Get(ctx, "ba-1")
			case "bank.list":
				actual, err = pay.BankAccounts.List(ctx, &payment.ListBankAccountsRequest{PageSize: 10, PageNumber: 1})
			case "bank.create":
				var req payment.CreateBankAccountRequest
				decode(&req)
				actual, err = pay.BankAccounts.Create(ctx, &req)
			case "intent.get":
				actual, err = pay.PaymentIntents.Get(ctx, "pi-1")
			case "intent.create":
				var req payment.CreatePaymentIntentRequest
				decode(&req)
				actual, err = pay.PaymentIntents.Create(ctx, &req)
			case "intent.confirm":
				var req payment.ConfirmPaymentIntentRequest
				decode(&req)
				actual, err = pay.PaymentIntents.Confirm(ctx, "pi-1", &req)
			case "attempt":
				actual, err = pay.PaymentAttempts.Get(ctx, "pa-1")
			default:
				t.Fatal(fixture.Operation)
			}
			if err != nil {
				t.Fatal(err)
			}
			req := <-requests
			if req.Method != fixture.Method || req.URL.Path != fixture.Path {
				t.Fatal(req.Method, req.URL, fixture.Path)
			}
			var want map[string]interface{}
			if err := json.Unmarshal(fixture.Body, &want); err != nil {
				t.Fatal(err)
			}
			// Compare exposed Go fields before omitempty removes their empty representations.
			switch v := actual.(type) {
			case *banking.PayoutDetailResponse:
				expected := want["payer"].(map[string]interface{})
				if v.Payer.PayerID != expected["payer_id"] || v.Payer.IdentificationType != expected["identification_type"] {
					t.Fatal("payer fields changed")
				}
				if expected["identification_type"] == "" {
					delete(expected, "identification_type")
				}
			case *Transaction:
				if fixture.Name == "optional-absent" && (v.WalletType != nil || v.MerchantData != nil) {
					t.Fatal("invented optional fields")
				}
			case *payment.BankAccount:
				if want["bank_code_type"] == "" {
					if v.BankCodeType != "" {
						t.Fatal(v)
					}
					delete(want, "bank_code_type")
				}
			case *payment.ListBankAccountsResponse:
				if len(v.Data) != 1 {
					t.Fatal(v)
				}
				expected := want["data"].([]interface{})[0].(map[string]interface{})
				if expected["bank_code_type"] == "" {
					if v.Data[0].BankCodeType != "" {
						t.Fatal(v)
					}
					delete(expected, "bank_code_type")
				}
			case *payment.PaymentIntent:
				fields := map[string]interface{}{"metadata": v.Metadata, "next_action": v.NextAction, "latest_payment_attempt": v.LatestPaymentAttempt}
				for key, value := range fields {
					encoded, _ := json.Marshal(value)
					var got interface{}
					if err := json.Unmarshal(encoded, &got); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, want[key]) {
						t.Fatalf("%s: %#v want %#v", key, got, want[key])
					}
					delete(want, key)
				}
				for key, value := range map[string]string{"cancel_time": v.CancelTime, "complete_time": v.CompleteTime} {
					expected, _ := want[key].(string)
					if value != expected {
						t.Fatal(key, value, expected)
					}
					delete(want, key)
				}
			case *payment.PaymentAttempt:
				if method, ok := want["payment_method"].(map[string]interface{}); ok {
					kind := method["type"].(string)
					expected := method[kind].(map[string]interface{})
					var osType string
					if v.PaymentMethod == nil {
						t.Fatal("missing method")
					}
					switch kind {
					case "grabpay":
						if v.PaymentMethod.GrabPay == nil {
							t.Fatal("missing grabpay")
						}
						osType = v.PaymentMethod.GrabPay.OSType
					case "tng":
						if v.PaymentMethod.TNG == nil {
							t.Fatal("missing tng")
						}
						osType = v.PaymentMethod.TNG.OSType
					case "wechatpay":
						if v.PaymentMethod.WeChatPay == nil {
							t.Fatal("missing wechatpay")
						}
						osType = v.PaymentMethod.WeChatPay.OSType
					}
					if osType != expected["os_type"] {
						t.Fatal("OS type lost")
					}
					if osType == "" {
						delete(expected, "os_type")
					}
				}
			}
			encoded, err := json.Marshal(actual)
			if err != nil {
				t.Fatal(err)
			}
			var got interface{}
			if err := json.Unmarshal(encoded, &got); err != nil {
				t.Fatal(err)
			}
			assertProvidedFields(t, fixture.Name, got, want)
		})
	}
}
