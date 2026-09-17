package issuing

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v3/common"
	"github.com/uqpay/uqpay-sdk-go/v3/configuration"
	"github.com/uqpay/uqpay-sdk-go/v3/connect"
)

func TestRFIAndPINOrderResponses(t *testing.T) {
	raw, err := os.ReadFile("testdata/rfi-orders.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct {
		RFIs   []json.RawMessage `json:"rfis"`
		Orders []json.RawMessage `json:"orders"`
	}
	if err := json.Unmarshal(raw, &fixtures); err != nil {
		t.Fatal(err)
	}
	var response []byte
	requests := make(chan *http.Request, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- r
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(response)
	}))
	defer server.Close()
	api := common.NewAPIClient(&configuration.Configuration{Environment: &configuration.Environment{BaseURL: server.URL}, HTTPClient: server.Client()}, &staticTokenProvider{token: "offline"})
	rfis := connect.NewClient(api).RFIs
	cards := NewClient(api).Cards
	ctx := context.Background()
	checkPath := func(path string) *http.Request {
		t.Helper()
		r := <-requests
		if r.Method != "GET" || r.URL.Path != path {
			t.Fatalf("request %s %s, want GET %s", r.Method, r.URL.Path, path)
		}
		return r
	}
	for _, fixture := range fixtures.RFIs {
		var expected map[string]interface{}
		if err := json.Unmarshal(fixture, &expected); err != nil {
			t.Fatal(err)
		}
		response = fixture
		result, err := rfis.Get(ctx, expected["rfi_id"].(string))
		if err != nil {
			t.Fatal(err)
		}
		checkPath("/v1/rfis/" + expected["rfi_id"].(string))
		assertRFIWireFields(t, *result, expected)
		response = []byte(`{"data":[` + string(fixture) + `],"total_pages":3,"total_items":21}`)
		page, err := rfis.List(ctx, &connect.ListRFIsRequest{PageSize: 10, PageNumber: 2, Status: connect.RFIStatusActionRequired})
		if err != nil {
			t.Fatal(err)
		}
		r := checkPath("/v1/rfis")
		if r.URL.Query().Get("page_size") != "10" || r.URL.Query().Get("page_number") != "2" || r.URL.Query().Get("status") != "ACTION_REQUIRED" {
			t.Fatal(r.URL)
		}
		if page.TotalPages != 3 || page.TotalItems != 21 || len(page.Data) != 1 {
			t.Fatalf("wrong page: %+v", page)
		}
		assertRFIWireFields(t, page.Data[0], expected)
	}
	for _, fixture := range fixtures.Orders {
		var expected map[string]interface{}
		if err := json.Unmarshal(fixture, &expected); err != nil {
			t.Fatal(err)
		}
		response = fixture
		order, err := cards.GetOrder(ctx, expected["card_order_id"].(string))
		if err != nil {
			t.Fatal(err)
		}
		checkPath("/v1/issuing/cards/" + expected["card_order_id"].(string) + "/order")
		fields := map[string]string{"card_id": order.CardID, "card_order_id": order.CardOrderID, "order_type": order.OrderType, "order_status": order.OrderStatus, "create_time": order.CreateTime, "update_time": order.UpdateTime, "complete_time": order.CompleteTime}
		for key, value := range fields {
			if expected[key] != value {
				t.Fatalf("%s = %q, want %v", key, value, expected[key])
			}
		}
		failure, _ := expected["failure_code"].(string)
		if order.FailureCode != failure {
			t.Fatalf("lost failure_code: %+v", order)
		}
		if order.OrderType == "PIN_MANAGEMENT" {
			// Public scalar fields preserve the existing zero-value behavior for absent fields.
			if order.Amount != 0 || order.CardCurrency != "" {
				t.Fatalf("PIN financial fields: %+v", order)
			}
		} else if order.Amount != 1.25 || order.CardCurrency != "USD" {
			t.Fatalf("recharge financial fields: %+v", order)
		}
	}
}

func assertRFIWireFields(t *testing.T, got connect.RFI, expected map[string]interface{}) {
	t.Helper()
	raw, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	var actual map[string]interface{}
	if err := json.Unmarshal(raw, &actual); err != nil {
		t.Fatal(err)
	}
	// An empty attachment slice is omitted on re-serialization; validate it directly.
	items := expected["request"].([]interface{})
	for i, item := range items {
		wantAnswer, ok := item.(map[string]interface{})["answer"].(map[string]interface{})
		if !ok {
			if got.Request[i].Answer != nil {
				t.Fatal("unexpected answer")
			}
			continue
		}
		if attachments, ok := wantAnswer["attachments"].([]interface{}); ok && len(attachments) == 0 {
			if got.Request[i].Answer == nil || len(got.Request[i].Answer.Attachments) != 0 {
				t.Fatal("expected empty attachments")
			}
			actual["request"].([]interface{})[i].(map[string]interface{})["answer"].(map[string]interface{})["attachments"] = []interface{}{}
		}
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("RFI fields mismatch: got %s want %+v", raw, expected)
	}
}
