package webhook

import (
	"encoding/json"
	"testing"
)

func TestPaymentAPISyncAlertRawAndIssuerCountry(t *testing.T) {
	// Chargeback alerts use Event.Data; no dedicated helper narrows their fields or enums.
	var event Event
	if err := json.Unmarshal([]byte(`{"event_type":"acquiring.chargeback.alert.created","data":{"order_status":"SUCCEEDED","alert_status":"EXPIRED","card_scheme":"VISA","card_masked":"411111******1111","auth_code":"test","arn":"reference"}}`), &event); err != nil {
		t.Fatal(err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(event.Data, &data); err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]interface{}{"order_status": "SUCCEEDED", "alert_status": "EXPIRED", "card_scheme": "VISA", "card_masked": "411111******1111", "auth_code": "test", "arn": "reference"} {
		if data[key] != want {
			t.Fatalf("%s=%v", key, data[key])
		}
	}
	var method PaymentMethod
	if err := json.Unmarshal([]byte(`{"card":{"issuer_country_code":"SG"},"card_present":{"issuer_country_code":""}}`), &method); err != nil {
		t.Fatal(err)
	}
	if method.Card.IssuerCountryCode != "SG" || method.CardPresent.IssuerCountryCode != "" {
		t.Fatal(method)
	}
	var attempt PaymentAttemptData
	if err := json.Unmarshal([]byte(`{"attempt_status":"EXPIRED"}`), &attempt); err != nil {
		t.Fatal(err)
	}
	if attempt.CapturedAmount != "" {
		t.Fatal("attempt captured_amount must be optional")
	}
}
