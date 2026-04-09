package authdecision

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// testEnv holds a fully configured test environment with key pairs and pgp contexts.
type testEnv struct {
	customerKP *KeyPair
	uqpayKP    *KeyPair
	// uqpayPgp encrypts to customer's public key (simulates UQPAY sending request)
	uqpayPgp *pgpContext
	// customerClient is the customer's AuthDecisionClient (decrypts with customer private key, encrypts to UQPAY public key)
	customerClient *AuthDecisionClient
}

var (
	envOnce    sync.Once
	sharedEnv  *testEnv
	envInitErr error
)

func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()
	envOnce.Do(func() {
		customerKP, err := GenerateKeyPair("Customer", "customer@test.com")
		if err != nil {
			envInitErr = fmt.Errorf("generate customer keys: %w", err)
			return
		}
		uqpayKP, err := GenerateKeyPair("UQPay", "uqpay@test.com")
		if err != nil {
			envInitErr = fmt.Errorf("generate uqpay keys: %w", err)
			return
		}

		// UQPAY encrypts to customer's public key
		uqpayPgp, err := newPgpContext(Config{
			PrivateKey:     uqpayKP.PrivateKey,
			UQPayPublicKey: customerKP.PublicKey,
		})
		if err != nil {
			envInitErr = fmt.Errorf("create uqpay pgpContext: %w", err)
			return
		}

		// Customer's AuthDecisionClient: decrypts with own private key, encrypts to UQPAY's public key
		client := NewAuthDecisionClient()
		if err := client.Configure(Config{
			PrivateKey:     customerKP.PrivateKey,
			UQPayPublicKey: uqpayKP.PublicKey,
		}); err != nil {
			envInitErr = fmt.Errorf("configure customer client: %w", err)
			return
		}

		sharedEnv = &testEnv{
			customerKP:     customerKP,
			uqpayKP:        uqpayKP,
			uqpayPgp:       uqpayPgp,
			customerClient: client,
		}
	})
	if envInitErr != nil {
		t.Fatalf("test environment setup failed: %v", envInitErr)
	}
	return sharedEnv
}

func TestHandlerPanicsWithoutConfigure(t *testing.T) {
	client := NewAuthDecisionClient()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic, got none")
		}
	}()

	client.Handler(HandlerOptions{
		Decide: func(ctx context.Context, tx Transaction) (Result, error) {
			return Result{ResponseCode: "00"}, nil
		},
	})
}

func TestHandlerApproveFlow(t *testing.T) {
	env := setupTestEnv(t)

	handler := env.customerClient.Handler(HandlerOptions{
		Decide: func(ctx context.Context, tx Transaction) (Result, error) {
			return Result{ResponseCode: "00", PartnerReferenceID: "ref-123"}, nil
		},
	})

	tx := Transaction{
		TransactionID:   "tx-001",
		TransactionType: 1,
		CardID:          "card-abc",
		BillingAmount:   100.50,
		MerchantName:    "Test Merchant",
	}

	txJSON, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("marshal tx: %v", err)
	}

	// UQPAY encrypts the transaction for the customer
	encryptedBody, err := env.uqpayPgp.encrypt(string(txJSON))
	if err != nil {
		t.Fatalf("encrypt tx: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader(encryptedBody))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("expected Content-Type application/json; charset=utf-8, got %q", ct)
	}

	// UQPAY decrypts the response (encrypted to UQPAY's public key)
	uqpayDecryptCtx, err := newPgpContext(Config{
		PrivateKey:     env.uqpayKP.PrivateKey,
		UQPayPublicKey: env.customerKP.PublicKey,
	})
	if err != nil {
		t.Fatalf("create uqpay decrypt context: %v", err)
	}

	decryptedResp, err := uqpayDecryptCtx.decrypt(rec.Body.String())
	if err != nil {
		t.Fatalf("decrypt response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedResp), &result); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if result["response_code"] != "00" {
		t.Fatalf("expected response_code 00, got %v", result["response_code"])
	}
	if result["partner_reference_id"] != "ref-123" {
		t.Fatalf("expected partner_reference_id ref-123, got %v", result["partner_reference_id"])
	}
	if result["transaction_id"] != "tx-001" {
		t.Fatalf("expected transaction_id tx-001, got %v", result["transaction_id"])
	}
}

func TestHandlerAutoInjectTransactionID(t *testing.T) {
	env := setupTestEnv(t)

	handler := env.customerClient.Handler(HandlerOptions{
		Decide: func(ctx context.Context, tx Transaction) (Result, error) {
			// Return without setting TransactionID or PartnerReferenceID
			return Result{ResponseCode: "05"}, nil
		},
	})

	tx := Transaction{
		TransactionID:   "tx-auto-inject",
		TransactionType: 2,
	}

	txJSON, _ := json.Marshal(tx)
	encryptedBody, err := env.uqpayPgp.encrypt(string(txJSON))
	if err != nil {
		t.Fatalf("encrypt tx: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader(encryptedBody))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// Decrypt response
	uqpayDecryptCtx, _ := newPgpContext(Config{
		PrivateKey:     env.uqpayKP.PrivateKey,
		UQPayPublicKey: env.customerKP.PublicKey,
	})

	decryptedResp, err := uqpayDecryptCtx.decrypt(rec.Body.String())
	if err != nil {
		t.Fatalf("decrypt response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedResp), &result); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if result["transaction_id"] != "tx-auto-inject" {
		t.Fatalf("expected transaction_id tx-auto-inject, got %v", result["transaction_id"])
	}
	if result["partner_reference_id"] != "" {
		t.Fatalf("expected empty partner_reference_id, got %v", result["partner_reference_id"])
	}
}

func TestHandlerDecideError(t *testing.T) {
	env := setupTestEnv(t)

	var capturedErr error
	handler := env.customerClient.Handler(HandlerOptions{
		Decide: func(ctx context.Context, tx Transaction) (Result, error) {
			return Result{}, fmt.Errorf("decision engine unavailable")
		},
		OnError: func(err error) {
			capturedErr = err
		},
	})

	tx := Transaction{TransactionID: "tx-err"}
	txJSON, _ := json.Marshal(tx)
	encryptedBody, _ := env.uqpayPgp.encrypt(string(txJSON))

	req := httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader(encryptedBody))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// No response body should be written on error
	body, _ := io.ReadAll(rec.Body)
	if len(body) != 0 {
		t.Fatalf("expected empty body on error, got %d bytes", len(body))
	}

	if capturedErr == nil {
		t.Fatal("expected OnError to be called")
	}
}

func TestHandlerMalformedBody(t *testing.T) {
	env := setupTestEnv(t)

	var capturedErr error
	handler := env.customerClient.Handler(HandlerOptions{
		Decide: func(ctx context.Context, tx Transaction) (Result, error) {
			return Result{ResponseCode: "00"}, nil
		},
		OnError: func(err error) {
			capturedErr = err
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader("not a pgp message"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body, _ := io.ReadAll(rec.Body)
	if len(body) != 0 {
		t.Fatalf("expected empty body on malformed input, got %d bytes", len(body))
	}

	if capturedErr == nil {
		t.Fatal("expected OnError to be called for malformed body")
	}
}

func TestHandlerTimeout(t *testing.T) {
	env := setupTestEnv(t)

	var capturedErr error
	handler := env.customerClient.Handler(HandlerOptions{
		Decide: func(ctx context.Context, tx Transaction) (Result, error) {
			// Block until context is cancelled (4.5s timeout)
			<-ctx.Done()
			return Result{}, ctx.Err()
		},
		OnError: func(err error) {
			capturedErr = err
		},
	})

	tx := Transaction{TransactionID: "tx-timeout"}
	txJSON, _ := json.Marshal(tx)
	encryptedBody, _ := env.uqpayPgp.encrypt(string(txJSON))

	req := httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader(encryptedBody))
	rec := httptest.NewRecorder()

	start := time.Now()
	handler.ServeHTTP(rec, req)
	elapsed := time.Since(start)

	// Should take approximately 4.5 seconds
	if elapsed < 4*time.Second || elapsed > 6*time.Second {
		t.Fatalf("expected ~4.5s timeout, got %v", elapsed)
	}

	body, _ := io.ReadAll(rec.Body)
	if len(body) != 0 {
		t.Fatalf("expected empty body on timeout, got %d bytes", len(body))
	}

	if capturedErr == nil {
		t.Fatal("expected OnError to be called on timeout")
	}
}
