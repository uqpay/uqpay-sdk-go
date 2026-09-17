package webhook

import (
	"encoding/json"
	"testing"
)

func TestAllPaymentMethodDetails(t *testing.T) {
	for _, method := range []string{"card", "card_present", "wechatpay", "alipay", "alipaycn", "alipayhk", "paynow", "grabpay", "applepay", "googlepay", "unionpay", "crypto", "tng", "truemoney", "gcash", "dana", "kakaopay", "tosspay", "naverpay", "mpay", "kplus", "boost", "rabbitlinepay", "kaspi", "hipay", "shopeepay"} {
		t.Run(method, func(t *testing.T) {
			details := map[string]interface{}{"flow": "qrcode", "os_type": "web", "static_qrcode": "qr-data", "static_qrcode_extension": "png", "static_qrcode_number_plate": "plate"}
			if method == "card" || method == "card_present" {
				details = map[string]interface{}{"card_name": "Test", "card_number": "411111******1111", "network": "VISA"}
			}
			raw, _ := json.Marshal(map[string]interface{}{"type": method, method: details})
			var payment PaymentMethod
			if err := json.Unmarshal(raw, &payment); err != nil {
				t.Fatal(err)
			}
			out, _ := json.Marshal(payment)
			var got map[string]interface{}
			_ = json.Unmarshal(out, &got)
			gotDetails, ok := got[method].(map[string]interface{})
			if !ok {
				t.Fatalf("details lost: %s", out)
			}
			for k, v := range details {
				if gotDetails[k] != v {
					t.Errorf("%s lost: %s", k, out)
				}
			}
		})
	}
}

func TestWebhookNullAndUnknownValues(t *testing.T) {
	raw := []byte(`{"event_type":"acquiring.payment_intent.succeeded","data":{"amount":"12345678901234567890.12345678","complete_time":null,"cancel_time":"","metadata":null,"payment_method":null,"next_action":{"redirect_to_url":{"return_url":""}}}}`)
	var event Event
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatal(err)
	}
	parsed, err := event.ParsePaymentIntentData()
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Amount != "12345678901234567890.12345678" || parsed.CompleteTime != nil || parsed.CancelTime == nil || *parsed.CancelTime != "" || parsed.PaymentMethod != nil || parsed.NextAction == nil {
		t.Fatalf("lost fields: %+v", parsed)
	}
	var tx CardTransactionData
	if err = json.Unmarshal([]byte(`{"wallet_type":"FUTURE_WALLET","transaction_amount":"0.00000001"}`), &tx); err != nil {
		t.Fatal(err)
	}
	if tx.WalletType != "FUTURE_WALLET" || tx.TransactionAmount != "0.00000001" {
		t.Fatal(tx)
	}
	var holder Cardholder
	if err = json.Unmarshal([]byte(`{"reason":"more evidence required"}`), &holder); err != nil {
		t.Fatal(err)
	}
	if holder.Reason != "more evidence required" {
		t.Fatal(holder)
	}
}
