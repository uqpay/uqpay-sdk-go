package webhook

import (
	"encoding/json"
	"testing"
)

func TestVirtualAccountApplicationWebhookVersionsAndOrderingFields(t *testing.T) {
	for _, version := range []string{"V1.5.1", "V1.5.2", "V1.6.0"} {
		for _, eventType := range []string{EventTypeVirtualAccountCreate, EventTypeVirtualAccountUpdate, EventTypeVirtualAccountClosed} {
			payload := `{"application_id":"app-id","public_version":3,"country":"BH","currency":"USD","status":"CLOSED","results":[{"payment_method":"SWIFT","status":"CLOSED","virtual_accounts":[{"account_bank_id":"bank-id","account_holder":"Merchant","account_number":"123","country_code":"BH","currency":"USD","bank_name":"Bank","bank_address":"Address","clearing_system":{"type":"bic_swift","value":"BANKBHBM"},"status":"CLOSED","close_reason":""}],"error":null}]}`
			event := Event{Version: version, EventName: EventNameVirtual, EventType: eventType, EventID: "event-id", SourceID: "app-id", Data: json.RawMessage(payload)}
			data, err := event.ParseVirtualAccountApplicationData()
			if err != nil {
				t.Fatalf("%s %s: %v", version, eventType, err)
			}
			if event.SourceID != data.ApplicationID || data.PublicVersion != 3 {
				t.Fatalf("source/version mismatch: %+v %+v", event, data)
			}
			if data.Results[0].VirtualAccounts[0].CloseReason != "" || data.Results[0].VirtualAccounts[0].ClearingSystem.Type != "bic_swift" {
				t.Fatalf("bank detail mismatch: %+v", data)
			}
		}
	}
}
