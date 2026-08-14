package webhook

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/uqpay/uqpay-sdk-go/v2/banking"
)

func TestVirtualAccountApplicationWebhookVersionsAndOrderingFields(t *testing.T) {
	for _, version := range []string{"V1.5.1", "V1.5.2", "V1.6.0"} {
		for _, eventType := range []string{EventTypeVirtualAccountCreate, EventTypeVirtualAccountUpdate, EventTypeVirtualAccountClosed} {
			payload := `{"account_id":"account-id","direct_id":"direct-id","application_id":"app-id","public_version":3,"country":"BH","currency":"USD","status":"CLOSED","results":[{"payment_method":"SWIFT","status":"CLOSED","virtual_accounts":[{"account_bank_id":"bank-id","account_holder":"Merchant","account_number":"123","country_code":"BH","currency":"USD","bank_name":"Bank","bank_address":"Address","clearing_system":{"type":"bic_swift","value":"BANKBHBM"},"status":"CLOSED","close_reason":""}],"error":null}]}`
			event := Event{Version: version, EventName: EventNameVirtual, EventType: eventType, EventID: "event-id", SourceID: "app-id", Data: json.RawMessage(payload)}
			data, err := event.ParseVirtualAccountApplicationData()
			if err != nil {
				t.Fatalf("%s %s: %v", version, eventType, err)
			}
			if event.SourceID != data.ApplicationID || data.PublicVersion != 3 {
				t.Fatalf("source/version mismatch: %+v %+v", event, data)
			}
			if data.AccountID != "account-id" || data.DirectID != "direct-id" {
				t.Fatalf("account context mismatch: %+v", data)
			}
			if data.Results[0].VirtualAccounts[0].CloseReason != "" || data.Results[0].VirtualAccounts[0].ClearingSystem.Type != "bic_swift" {
				t.Fatalf("bank detail mismatch: %+v", data)
			}
		}
	}
}

func TestVirtualAccountApplicationParserRejectsMissingAccountContext(t *testing.T) {
	for _, missing := range []string{"account_id", "direct_id"} {
		payload := map[string]any{
			"account_id": "account-id", "direct_id": "direct-id",
			"application_id": "app-id", "public_version": 1,
			"country": "BH", "currency": "USD", "status": "SUBMITTED",
			"results": []any{},
		}
		delete(payload, missing)
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		event := Event{Version: "V1.6.0", EventName: EventNameVirtual, EventType: EventTypeVirtualAccountCreate, SourceID: "app-id", Data: raw}
		if _, err := event.ParseVirtualAccountApplicationData(); err == nil {
			t.Fatalf("missing %s must fail typed application parsing", missing)
		}
	}
}

func TestUnknownOldVirtualEventRemainsGeneric(t *testing.T) {
	payload := json.RawMessage(`{"account_bank_id":"legacy-bank-id"}`)
	event := Event{Version: "V1.5.0", EventName: EventNameVirtual, EventType: EventTypeVirtualAccountCreate, EventID: "legacy-event-id", SourceID: "legacy-bank-id", Data: payload}
	if _, err := event.ParseVirtualAccountApplicationData(); err == nil {
		t.Fatal("old event must not be reclassified as an application event")
	}
	if string(event.Data) != string(payload) {
		t.Fatalf("generic event data changed: %s", event.Data)
	}
}

func TestVirtualAccountApplicationParserRejectsSourceMismatch(t *testing.T) {
	payload := json.RawMessage(`{"account_id":"account-id","direct_id":"direct-id","application_id":"app-id","public_version":1,"country":"BH","currency":"USD","status":"SUBMITTED","results":[]}`)
	event := Event{Version: "V1.6.0", EventName: EventNameVirtual, EventType: EventTypeVirtualAccountCreate, SourceID: "different-id", Data: payload}
	if _, err := event.ParseVirtualAccountApplicationData(); err == nil {
		t.Fatal("source_id mismatch must fail typed application parsing")
	}
}

func TestRESTApplicationPublicTypesIncludeAccountContext(t *testing.T) {
	application := banking.VirtualAccountApplication{
		AccountID:     "connected-account-uuid",
		DirectID:      "main-account-uuid",
		ApplicationID: "app-id",
		PublicVersion: 1,
		Country:       "BH",
		Currency:      "USD",
		Status:        banking.VirtualAccountApplicationSubmitted,
		Results:       []banking.VirtualAccountApplicationResult{},
	}
	for name, value := range map[string]any{
		"application": application,
		"response":    banking.VirtualAccountApplicationResponse{Data: application},
		"summary": banking.VirtualAccountApplicationSummary{
			AccountID: "main-account-uuid", DirectID: "0", ApplicationID: "app-id",
			PublicVersion: 1, Country: "BH", Currency: "USD",
			Status: banking.VirtualAccountApplicationSubmitted, CreatedAt: "2026-08-12T00:00:00Z",
		},
	} {
		raw, err := json.Marshal(value)
		if err != nil {
			t.Fatalf("marshal %s: %v", name, err)
		}
		jsonText := string(raw)
		if !strings.Contains(jsonText, `"account_id"`) || !strings.Contains(jsonText, `"direct_id"`) {
			t.Fatalf("REST public %s type is missing required account context: %s", name, jsonText)
		}
	}
}
