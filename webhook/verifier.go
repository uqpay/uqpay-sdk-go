package webhook

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const DefaultSignatureTolerance = 5 * time.Minute

// Verifier validates UQPAY webhook signatures and parses the common event
// envelope. EventType remains a string so newly introduced server events stay
// forward-compatible even before a typed payload helper is added.
type Verifier struct {
	secret    []byte
	tolerance time.Duration
	now       func() time.Time
}

// NewVerifier creates a webhook verifier. Pass a positive tolerance to
// override the default five-minute replay window.
func NewVerifier(secret string, tolerance ...time.Duration) *Verifier {
	resolvedTolerance := DefaultSignatureTolerance
	if len(tolerance) > 0 && tolerance[0] > 0 {
		resolvedTolerance = tolerance[0]
	}
	return &Verifier{
		secret:    []byte(secret),
		tolerance: resolvedTolerance,
		now:       time.Now,
	}
}

// ConstructEvent verifies HMAC-SHA512(secret, rawPayload + timestamp) and
// returns the parsed event envelope.
func (v *Verifier) ConstructEvent(payload []byte, signatureHeader, timestampHeader string) (*Event, error) {
	if signatureHeader == "" {
		return nil, fmt.Errorf("webhook header missing: x-wk-signature")
	}
	if timestampHeader == "" {
		return nil, fmt.Errorf("webhook header missing: x-wk-timestamp")
	}

	timestamp, err := strconv.ParseInt(timestampHeader, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid x-wk-timestamp %q: %w", timestampHeader, err)
	}
	if delta := v.now().Sub(time.Unix(timestamp, 0)); delta > v.tolerance || delta < -v.tolerance {
		return nil, fmt.Errorf("webhook timestamp is outside the allowed tolerance of %s", v.tolerance)
	}

	received, err := hex.DecodeString(signatureHeader)
	if err != nil {
		return nil, fmt.Errorf("webhook signature verification failed")
	}
	mac := hmac.New(sha512.New, v.secret)
	_, _ = mac.Write(payload)
	_, _ = mac.Write([]byte(timestampHeader))
	if !hmac.Equal(mac.Sum(nil), received) {
		return nil, fmt.Errorf("webhook signature verification failed")
	}

	var event Event
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("webhook body is not valid JSON: %w", err)
	}
	return &event, nil
}
