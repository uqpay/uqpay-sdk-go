package authdecision

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	maxBodySize   = 1 << 20 // 1 MB
	decideTimeout = 4500 * time.Millisecond
)

// AuthDecisionClient handles PGP-encrypted card authorization decision requests.
type AuthDecisionClient struct {
	pgp *pgpContext
}

// NewAuthDecisionClient creates a new unconfigured AuthDecisionClient.
func NewAuthDecisionClient() *AuthDecisionClient {
	return &AuthDecisionClient{}
}

// Configure initializes the PGP context from the given Config.
func (c *AuthDecisionClient) Configure(config Config) error {
	pgp, err := newPgpContext(config)
	if err != nil {
		return fmt.Errorf("authdecision: configure failed: %w", err)
	}
	c.pgp = pgp
	return nil
}

// Handler returns an http.HandlerFunc that processes PGP-encrypted authorization
// decision requests. It panics if Configure has not been called or if
// opts.Decide is nil.
func (c *AuthDecisionClient) Handler(opts HandlerOptions) http.HandlerFunc {
	if c.pgp == nil {
		panic("authdecision: Handler called before Configure")
	}
	if opts.Decide == nil {
		panic("authdecision: Handler called with nil Decide function")
	}

	// Capture pgp at handler creation time to avoid race if Configure is called again later.
	pgp := c.pgp

	return func(w http.ResponseWriter, r *http.Request) {
		if err := handleRequest(pgp, w, r, opts); err != nil {
			if opts.OnError != nil {
				opts.OnError(err)
			}
			// Do NOT write any response on error — let UQPAY timeout strategy decide.
		}
	}
}

func handleRequest(pgp *pgpContext, w http.ResponseWriter, r *http.Request, opts HandlerOptions) error {
	// 1. Read body with size limit
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("authdecision: failed to read request body: %w", err)
	}

	// 2. PGP decrypt
	plaintext, err := pgp.decrypt(string(body))
	if err != nil {
		return fmt.Errorf("authdecision: failed to decrypt request: %w", err)
	}

	// 3. JSON unmarshal to Transaction
	var tx Transaction
	if err := json.Unmarshal([]byte(plaintext), &tx); err != nil {
		return fmt.Errorf("authdecision: failed to unmarshal transaction: %w", err)
	}

	// 4. Call Decide with 4.5s context timeout
	ctx, cancel := context.WithTimeout(r.Context(), decideTimeout)
	defer cancel()

	result, err := opts.Decide(ctx, tx)
	if err != nil {
		return fmt.Errorf("authdecision: decide failed: %w", err)
	}

	// 5. Build response: auto-inject transaction_id
	resp := struct {
		TransactionID      string `json:"transaction_id"`
		ResponseCode       string `json:"response_code"`
		PartnerReferenceID string `json:"partner_reference_id"`
	}{
		TransactionID:      tx.TransactionID,
		ResponseCode:       result.ResponseCode,
		PartnerReferenceID: result.PartnerReferenceID,
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("authdecision: failed to marshal response: %w", err)
	}

	// 6. PGP encrypt response
	encrypted, err := pgp.encrypt(string(respJSON))
	if err != nil {
		return fmt.Errorf("authdecision: failed to encrypt response: %w", err)
	}

	// 7. Write HTTP 200
	// UQPAY sends requests with application/json despite PGP-armored body; match their convention.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(encrypted))

	return nil
}
