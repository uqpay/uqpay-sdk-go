package webhook

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"testing"
	"time"
)

func TestVerifierAcceptsPublishedSignatureAndUnknownEventType(t *testing.T) {
	const secret = "whsec_test_secret"
	payload := []byte(`{"version":"V1.6.0","event_name":"ISSUING","event_type":"future.event.added","event_id":"evt_123","data":{"new_field":true}}`)
	timestamp := time.Now().Unix()
	timestampHeader := stringInt(timestamp)
	signature := signWebhook(secret, payload, timestampHeader)

	event, err := NewVerifier(secret).ConstructEvent(payload, signature, timestampHeader)
	if err != nil {
		t.Fatalf("ConstructEvent returned an error: %v", err)
	}
	if event.EventType != "future.event.added" || event.EventName != "ISSUING" {
		t.Errorf("event = %#v, want forward-compatible envelope", event)
	}
}

func TestVerifierRejectsInvalidSignatureAndStaleTimestamp(t *testing.T) {
	const secret = "whsec_test_secret"
	payload := []byte(`{"event_type":"deposit.completed"}`)
	now := time.Now()
	verifier := NewVerifier(secret)
	verifier.now = func() time.Time { return now }

	currentTimestamp := stringInt(now.Unix())
	if _, err := verifier.ConstructEvent(payload, "deadbeef", currentTimestamp); err == nil {
		t.Fatal("ConstructEvent accepted an invalid signature")
	}

	staleTimestamp := stringInt(now.Add(-10 * time.Minute).Unix())
	staleSignature := signWebhook(secret, payload, staleTimestamp)
	if _, err := verifier.ConstructEvent(payload, staleSignature, staleTimestamp); err == nil {
		t.Fatal("ConstructEvent accepted a stale timestamp")
	}
}

func signWebhook(secret string, payload []byte, timestamp string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	_, _ = mac.Write(payload)
	_, _ = mac.Write([]byte(timestamp))
	return hex.EncodeToString(mac.Sum(nil))
}

func stringInt(value int64) string {
	return strconv.FormatInt(value, 10)
}
