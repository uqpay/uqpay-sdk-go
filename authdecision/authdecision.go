package authdecision

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const defaultMaxBodyBytes int64 = 1 << 20

var errNotConfigured = errors.New("authdecision: not configured; call Configure first")

// Client decrypts UQPAY authorization requests and encrypts customer decisions.
type Client struct {
	mu  sync.RWMutex
	pgp *pgpContext
}

// NewClient creates an unconfigured authorization decision client.
func NewClient() *Client {
	return &Client{}
}

// Configure parses and validates the customer private key and UQPAY public key.
// It is safe to replace the configuration while existing handlers are running;
// each request uses one immutable configuration snapshot.
func (c *Client) Configure(config Config) error {
	pgp, err := newPGPContext(config)
	if err != nil {
		return fmt.Errorf("authdecision: configure: %w", err)
	}
	c.mu.Lock()
	c.pgp = pgp
	c.mu.Unlock()
	return nil
}

// Process decrypts one request, invokes decide, injects the transaction ID, and
// returns an ASCII-armored encrypted response.
func (c *Client) Process(ctx context.Context, encryptedBody []byte, decide DecisionFunc) ([]byte, error) {
	if ctx == nil {
		return nil, fmt.Errorf("authdecision: context is required")
	}
	if decide == nil {
		return nil, fmt.Errorf("authdecision: decision function is required")
	}
	pgp := c.snapshot()
	if pgp == nil {
		return nil, errNotConfigured
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	plaintext, err := pgp.decrypt(string(encryptedBody))
	if err != nil {
		return nil, err
	}
	var transaction Transaction
	if err := json.Unmarshal([]byte(plaintext), &transaction); err != nil {
		return nil, fmt.Errorf("authdecision: decode transaction: %w", err)
	}
	if transaction.TransactionID == "" {
		return nil, fmt.Errorf("authdecision: transaction_id is required")
	}

	result, err := invokeDecision(ctx, decide, transaction)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if result.ResponseCode == "" {
		return nil, fmt.Errorf("authdecision: response_code is required")
	}

	response, err := json.Marshal(struct {
		TransactionID      string `json:"transaction_id"`
		ResponseCode       string `json:"response_code"`
		PartnerReferenceID string `json:"partner_reference_id"`
	}{
		TransactionID:      transaction.TransactionID,
		ResponseCode:       result.ResponseCode,
		PartnerReferenceID: result.PartnerReferenceID,
	})
	if err != nil {
		return nil, fmt.Errorf("authdecision: encode response: %w", err)
	}
	encryptedResponse, err := pgp.encrypt(string(response))
	if err != nil {
		return nil, err
	}
	return []byte(encryptedResponse), nil
}

type decisionOutcome struct {
	result Result
	err    error
}

func invokeDecision(ctx context.Context, decide DecisionFunc, transaction Transaction) (Result, error) {
	if ctx.Done() == nil {
		result, err := decide(ctx, transaction)
		if err != nil {
			return Result{}, fmt.Errorf("authdecision: decide: %w", err)
		}
		return result, nil
	}

	outcome := make(chan decisionOutcome, 1)
	go func() {
		deferred := decisionOutcome{}
		defer func() {
			if recovered := recover(); recovered != nil {
				deferred.err = fmt.Errorf("authdecision: decision panic: %v", recovered)
			}
			outcome <- deferred
		}()
		deferred.result, deferred.err = decide(ctx, transaction)
		if deferred.err != nil {
			deferred.err = fmt.Errorf("authdecision: decide: %w", deferred.err)
		}
	}()

	select {
	case <-ctx.Done():
		return Result{}, ctx.Err()
	case completed := <-outcome:
		return completed.result, completed.err
	}
}

// Handler returns a net/http handler for UQPAY authorization decision requests.
// On processing errors it invokes OnError and aborts the response so UQPAY can
// apply the configured timeout action instead of receiving an accidental empty
// HTTP 200 response.
func (c *Client) Handler(options HandlerOptions) (http.HandlerFunc, error) {
	if c.snapshot() == nil {
		return nil, errNotConfigured
	}
	if options.Decide == nil {
		return nil, fmt.Errorf("authdecision: decision function is required")
	}
	if options.DecisionTimeout < 0 {
		return nil, fmt.Errorf("authdecision: decision timeout cannot be negative")
	}
	maxBodyBytes := options.MaxBodyBytes
	if maxBodyBytes == 0 {
		maxBodyBytes = defaultMaxBodyBytes
	}
	if maxBodyBytes < 0 {
		return nil, fmt.Errorf("authdecision: max body bytes cannot be negative")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body := http.MaxBytesReader(w, r.Body, maxBodyBytes)
		encryptedBody, err := io.ReadAll(body)
		if err != nil {
			abortHTTPResponse(options.OnError, fmt.Errorf("authdecision: read request body: %w", err))
		}

		ctx := r.Context()
		if options.DecisionTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, options.DecisionTimeout)
			defer cancel()
		}
		encryptedResponse, err := c.Process(ctx, encryptedBody, options.Decide)
		if err != nil {
			abortHTTPResponse(options.OnError, err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(encryptedResponse); err != nil && options.OnError != nil {
			options.OnError(fmt.Errorf("authdecision: write response: %w", err))
		}
	}, nil
}

func (c *Client) snapshot() *pgpContext {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pgp
}

func abortHTTPResponse(onError func(error), err error) {
	if onError != nil {
		onError(err)
	}
	panic(http.ErrAbortHandler)
}
