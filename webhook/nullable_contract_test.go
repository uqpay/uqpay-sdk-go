package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

// WH-AQ: raw event bytes preserve absence/null/empty even when typed fields merge them.
func TestAcquiringNullableEventsAndSignatureBytes(t *testing.T) {
	for _, tc := range []struct {
		kind           string
		times, objects []string
	}{
		{"payment_intent.succeeded", []string{"complete_time", "cancel_time"}, []string{"metadata", "next_action", "payment_method"}},
		{"payment_attempt.succeeded", []string{"complete_time", "cancel_time"}, nil},
		{"refund.succeeded", []string{"complete_time"}, []string{"metadata"}},
		{"payout.succeeded", []string{"complete_time"}, nil},
		{"chargeback.alert.created", []string{"appeal_time", "response_time"}, nil},
	} {
		for _, mode := range []string{"missing", "null", "empty", "populated"} {
			data := map[string]interface{}{}
			if mode != "missing" {
				for _, field := range tc.times {
					var value interface{}
					if mode == "empty" {
						value = ""
					}
					if mode == "populated" {
						value = "2026-09-17T00:00:00Z"
					}
					data[field] = value
				}
				for _, field := range tc.objects {
					var value interface{}
					if mode == "empty" {
						value = map[string]string{}
					}
					if mode == "populated" {
						value = map[string]string{"ref": "0001"}
					}
					data[field] = value
				}
			}
			expected, _ := json.Marshal(data)
			raw, _ := json.Marshal(map[string]interface{}{"event_type": "acquiring." + tc.kind, "data": data})
			timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
			mac := hmac.New(sha512.New, []byte("offline-secret"))
			_, _ = mac.Write(append(append([]byte{}, raw...), []byte(timestamp)...))
			sig := hex.EncodeToString(mac.Sum(nil))
			verifier := NewVerifier("offline-secret")
			event, err := verifier.ConstructEvent(raw, sig, timestamp)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(event.Data, expected) {
				t.Fatalf("%s %s: raw changed", tc.kind, mode)
			}
			if _, err = verifier.ConstructEvent(append(raw, ' '), sig, timestamp); err == nil {
				t.Fatal("changed signature bytes accepted")
			}
			var complete, cancel *string
			switch tc.kind {
			case "payment_intent.succeeded":
				v, e := event.ParsePaymentIntentData()
				if e != nil {
					t.Fatal(e)
				}
				complete, cancel = v.CompleteTime, v.CancelTime
			case "payment_attempt.succeeded":
				v, e := event.ParsePaymentAttemptData()
				if e != nil {
					t.Fatal(e)
				}
				complete, cancel = v.CompleteTime, v.CancelTime
			case "refund.succeeded":
				v, e := event.ParseRefundData()
				if e != nil {
					t.Fatal(e)
				}
				complete = v.CompleteTime
			default:
				continue // No dedicated acquiring payout/chargeback helper; raw data is the supported interface.
			}
			for _, ptr := range []*string{complete, cancel} {
				if ptr != nil && mode == "null" {
					t.Fatal("null became non-nil")
				}
			}
			if mode == "empty" && (complete == nil || *complete != "") {
				t.Fatal("empty completion lost")
			}
			if mode == "populated" && (complete == nil || *complete != "2026-09-17T00:00:00Z") {
				t.Fatal("completion lost")
			}
		}
	}
}
