package issuing

import (
	"encoding/json"
	"testing"
)

func TestCardOrderUnmarshalAmount(t *testing.T) {
	tests := []struct {
		name       string
		amountJSON string
		want       float64
		wantErr    bool
	}{
		{name: "string", amountJSON: `"1.25"`, want: 1.25},
		{name: "number", amountJSON: `1.25`, want: 1.25},
		{name: "null", amountJSON: `null`, want: 0},
		{name: "invalid string", amountJSON: `"not-a-number"`, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload := []byte(`{
				"card_id":"card-1",
				"card_order_id":"order-1",
				"order_type":"CREATE_CARD",
				"amount":` + test.amountJSON + `,
				"card_currency":"USD",
				"order_status":"PROCESSING"
			}`)

			var order CardOrder
			err := json.Unmarshal(payload, &order)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unmarshal CardOrder: %v", err)
			}
			if order.Amount != test.want {
				t.Fatalf("Amount = %v, want %v", order.Amount, test.want)
			}
			if order.CardOrderID != "order-1" {
				t.Fatalf("CardOrderID = %q, want order-1", order.CardOrderID)
			}
		})
	}
}
