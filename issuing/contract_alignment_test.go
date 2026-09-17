package issuing

import (
	"context"
	"encoding/json"
	"github.com/uqpay/uqpay-sdk-go/v3/common"
	"github.com/uqpay/uqpay-sdk-go/v3/configuration"
	"github.com/uqpay/uqpay-sdk-go/v3/connect"
	"github.com/uqpay/uqpay-sdk-go/v3/simulator"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPINContractAlignment(t *testing.T) {
	req := SetPINRequest{CardID: "card-1", PIN: "135790", Type: "UPDATE", OldPIN: "024680"}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]interface{}
	_ = json.Unmarshal(b, &wire)
	if wire["type"] != "UPDATE" || wire["old_pin"] != "024680" {
		t.Fatalf("lost PIN fields: %s", b)
	}
	var accepted SetPINResponse
	if err := json.Unmarshal([]byte(`{"request_status":"SUCCESS","card_id":"card-1","card_order_id":"order-1","order_status":"PROCESSING","create_time":"2026-09-17T00:00:00Z"}`), &accepted); err != nil {
		t.Fatal(err)
	}
	if accepted.CardOrderID != "order-1" || accepted.OrderStatus != "PROCESSING" {
		t.Fatalf("lost asynchronous order: %+v", accepted)
	}
	var order CardOrder
	if err := json.Unmarshal([]byte(`{"order_type":"PIN_MANAGEMENT","order_status":"FAILED","failure_code":"pin_operation_failed"}`), &order); err != nil {
		t.Fatal(err)
	}
	if order.FailureCode != "pin_operation_failed" {
		t.Fatalf("lost failure: %+v", order)
	}
}
func TestTransactionSettlementContract(t *testing.T) {
	for _, status := range []string{"NOT_APPLICABLE", "UNSETTLED", "UNKNOWN", "SETTLED"} {
		var tx Transaction
		if err := json.Unmarshal([]byte(`{"settlement_status":"`+status+`"}`), &tx); err != nil {
			t.Fatal(err)
		}
		if tx.SettlementStatus == nil || *tx.SettlementStatus != status {
			t.Fatalf("lost settlement: %+v", tx)
		}
	}
	var tx Transaction
	_ = json.Unmarshal([]byte(`{}`), &tx)
	if tx.SettlementStatus != nil {
		t.Fatal("list item must allow absent settlement status")
	}
}

func TestAlignedContractsThroughTransport(t *testing.T) {
	var captured map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		captured = nil
		if r.Body != nil && r.Method == http.MethodPost {
			if err := json.NewDecoder(r.Body).Decode(&captured); err != nil {
				t.Error(err)
			}
		}
		switch r.URL.Path {
		case "/v1/issuing/cards/pin":
			_, _ = w.Write([]byte(`{"request_status":"SUCCESS","card_id":"card-1","card_order_id":"order-1","order_status":"PROCESSING"}`))
		case "/v1/rfis/answer":
			_, _ = w.Write([]byte(`{"rfi_id":"ACTREQ-test","request":[{"answer":{"type":"ATTACHMENT","attachments":[{"file_name":"proof.pdf","size":42}]}}]}`))
		case "/v1/simulation/deposit":
			_, _ = w.Write([]byte(`{"deposit_id":"deposit-1","amount":"10"}`))
		case "/v1/issuing/transactions/tx-1":
			_, _ = w.Write([]byte(`{"settlement_status":"SETTLED","transaction_amount":"123456789.01"}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	api := common.NewAPIClient(&configuration.Configuration{Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}, &staticTokenProvider{token: "offline-token"})
	issuing := NewClient(api)
	ctx := context.Background()
	for _, action := range []string{"", "SET", "RESET", "UPDATE"} {
		req := &SetPINRequest{CardID: "card-1", PIN: "135790", Type: action}
		if action == "UPDATE" {
			req.OldPIN = "024680"
		}
		res, err := issuing.Cards.ResetPIN(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		if res.CardOrderID != "order-1" || res.OrderStatus != "PROCESSING" {
			t.Fatalf("lost accepted response: %+v", res)
		}
		if action == "UPDATE" {
			if captured["old_pin"] != "024680" {
				t.Fatal("lost leading zero")
			}
		} else {
			if _, ok := captured["old_pin"]; ok {
				t.Fatal("unexpected old_pin")
			}
		}
	}
	rfi, err := connect.NewClient(api).RFIs.Answer(ctx, &connect.AnswerRFIRequest{RFIID: "ACTREQ-test", Answer: []connect.RFIAnswerItem{{Key: "note", Type: "TEXT", Text: "source of funds"}}})
	if err != nil {
		t.Fatal(err)
	}
	if rfi.Request[0].Answer.Attachments[0].FileName != "proof.pdf" || captured["rfi_id"] != "ACTREQ-test" {
		t.Fatal("lost RFI contract")
	}
	_, err = simulator.NewClient(api).Deposits.Create(ctx, &simulator.CreateDepositRequest{AccountID: "account-1", Amount: 10, Currency: "SGD", SenderSwiftCode: "WELGBE22"})
	if err != nil {
		t.Fatal(err)
	}
	if captured["account_id"] != "account-1" || captured["amount"] != float64(10) {
		t.Fatal("lost deposit fields")
	}
	tx, err := issuing.Transactions.Get(ctx, "tx-1")
	if err != nil {
		t.Fatal(err)
	}
	if tx.SettlementStatus == nil || *tx.SettlementStatus != "SETTLED" || tx.TransactionAmount != "123456789.01" {
		t.Fatalf("lost transaction fields: %+v", tx)
	}
}
