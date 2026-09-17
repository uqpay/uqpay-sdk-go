package banking

import (
	"context"
	"encoding/json"
	"github.com/uqpay/uqpay-sdk-go/v3/common"
	"github.com/uqpay/uqpay-sdk-go/v3/configuration"
	"net/http"
	"net/http/httptest"
	"testing"
)

// D122-D129: test actual transport, typed decoding and reserialization.
func TestBalanceDecimalStringsThroughListAndDetail(t *testing.T) {
	fields := []string{"available_balance", "frozen_balance", "margin_balance", "prepaid_balance"}
	amounts := []string{"0.00", "1.23", "-0.01", "12345678901234567890.12", "-12345678901234567890.12", "0.12345678901234567890"}
	var balance map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != "GET" {
			t.Error(r.Method)
		}
		switch r.URL.Path {
		case "/v1/balances/USD":
			_ = json.NewEncoder(w).Encode(balance)
		case "/v1/balances":
			if r.URL.Query().Get("page_size") != "10" || r.URL.Query().Get("page_number") != "1" {
				t.Error(r.URL)
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": []interface{}{balance}, "total_pages": 1, "total_items": 1})
		default:
			t.Error(r.URL)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	api := common.NewAPIClient(&configuration.Configuration{Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}, &vaStaticTokenProvider{})
	client := NewClient(api)
	check := func(got *Balance) {
		raw, err := json.Marshal(got)
		if err != nil {
			t.Fatal(err)
		}
		var wire map[string]interface{}
		if err = json.Unmarshal(raw, &wire); err != nil {
			t.Fatal(err)
		}
		for _, field := range fields {
			if wire[field] != balance[field] {
				t.Fatalf("%s: got %v want %s", field, wire[field], balance[field])
			}
		}
	}
	for i := range amounts {
		balance = map[string]string{"currency": "USD"}
		for j, field := range fields {
			balance[field] = amounts[(i+j)%len(amounts)]
		}
		detail, err := client.Balances.Get(context.Background(), "USD")
		if err != nil {
			t.Fatal(err)
		}
		check(detail)
		list, err := client.Balances.List(context.Background(), &ListBalancesRequest{PageSize: 10, PageNumber: 1})
		if err != nil {
			t.Fatal(err)
		}
		if len(list.Data) != 1 || list.TotalItems != 1 {
			t.Fatal(list)
		}
		check(&list.Data[0])
	}
}
