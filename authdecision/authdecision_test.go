package authdecision

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const authorizationRequest = `{
  "transaction_id": "550e8400-e29b-41d4-a716-446655440000",
  "transaction_type": 1000,
  "card_id": "card-001",
  "processing_code": "00",
  "billing_amount": "2.31",
  "transaction_amount": 2.31,
  "auth_amount": "0",
  "date_of_transaction": "2025-11-14 15:07:25",
  "billing_currency_code": "SGD",
  "transaction_currency_code": "CAD",
  "auth_currency_code": "USD",
  "card_balance": "90085.59",
  "merchant_id": "CARD ACCEPTOR",
  "merchant_name": "Example Store",
  "merchant_category_code": "5972",
  "merchant_city": "CITY NAME",
  "merchant_country": "US",
  "terminal_id": "TERMID01",
  "pos_entry_mode": "01",
  "pos_condition_code": "08",
  "pos_env": "R",
  "eci": "02",
  "pin_entry_capability": "2",
  "retrieval_reference_number": "529430718653",
  "system_trace_audit_number": "000653",
  "acquiring_institution_country_code": "TK",
  "acquiring_institution_id": "30954284708",
  "wallet_type": "GOOGLE ECOMMERCE"
}`

func TestProcessAuthorizationDecision(t *testing.T) {
	client, uqpayContext := configuredTestClient(t)
	encryptedRequest, err := uqpayContext.encrypt(authorizationRequest)
	if err != nil {
		t.Fatalf("encrypt request: %v", err)
	}

	encryptedResponse, err := client.Process(
		context.Background(),
		[]byte(encryptedRequest),
		func(_ context.Context, transaction Transaction) (Result, error) {
			if transaction.TransactionID != "550e8400-e29b-41d4-a716-446655440000" {
				t.Fatalf("transaction ID = %q", transaction.TransactionID)
			}
			if transaction.BillingAmount != "2.31" || transaction.TransactionAmount != "2.31" {
				t.Fatalf("amounts = %q/%q", transaction.BillingAmount, transaction.TransactionAmount)
			}
			if transaction.PosEnv != "R" || transaction.ECI != "02" {
				t.Fatalf("pos_env/eci = %q/%q", transaction.PosEnv, transaction.ECI)
			}
			return Result{ResponseCode: "00", PartnerReferenceID: "ref-001"}, nil
		},
	)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	plaintext, err := uqpayContext.decrypt(string(encryptedResponse))
	if err != nil {
		t.Fatalf("decrypt response: %v", err)
	}
	var response map[string]string
	if err := json.Unmarshal([]byte(plaintext), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if response["transaction_id"] != "550e8400-e29b-41d4-a716-446655440000" {
		t.Fatalf("response transaction ID = %q", response["transaction_id"])
	}
	if response["response_code"] != "00" || response["partner_reference_id"] != "ref-001" {
		t.Fatalf("response = %#v", response)
	}
}

func TestTransactionPreservesArbitraryPrecisionDecimalStrings(t *testing.T) {
	var transaction Transaction
	data := []byte(`{
		"billing_amount":"123456789012345678901234567890.12345678901234567890",
		"transaction_amount":1.25e3,
		"auth_amount":"0",
		"card_balance":"90085.59"
	}`)
	if err := json.Unmarshal(data, &transaction); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if transaction.BillingAmount != "123456789012345678901234567890.12345678901234567890" {
		t.Fatalf("billing amount = %q", transaction.BillingAmount)
	}
	if transaction.TransactionAmount != "1.25e3" {
		t.Fatalf("transaction amount = %q", transaction.TransactionAmount)
	}

	for _, value := range []string{`"NaN"`, `"Infinity"`, `"01.2"`, `"1."`} {
		payload := []byte(`{"billing_amount":` + value + `}`)
		if err := json.Unmarshal(payload, &transaction); err == nil {
			t.Fatalf("expected invalid decimal %s to fail", value)
		}
	}
}

func TestProcessRequiresConfigurationAndValidDecision(t *testing.T) {
	client := NewClient()
	if _, err := client.Process(context.Background(), []byte("body"), func(context.Context, Transaction) (Result, error) {
		return Result{ResponseCode: "00"}, nil
	}); err == nil {
		t.Fatal("expected unconfigured error")
	}

	configured, uqpayContext := configuredTestClient(t)
	encrypted, err := uqpayContext.encrypt(authorizationRequest)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if _, err := configured.Process(context.Background(), []byte(encrypted), nil); err == nil {
		t.Fatal("expected nil decision error")
	}
	if _, err := configured.Process(context.Background(), []byte(encrypted), func(context.Context, Transaction) (Result, error) {
		return Result{}, nil
	}); err == nil {
		t.Fatal("expected empty response code error")
	}
}

func TestProcessPropagatesContextAndDecisionErrors(t *testing.T) {
	client, uqpayContext := configuredTestClient(t)
	encrypted, err := uqpayContext.encrypt(authorizationRequest)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	expected := errors.New("decision failed")
	if _, err := client.Process(context.Background(), []byte(encrypted), func(context.Context, Transaction) (Result, error) {
		return Result{}, expected
	}); !errors.Is(err, expected) {
		t.Fatalf("Process error = %v, want wrapped decision error", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := client.Process(ctx, []byte(encrypted), func(ctx context.Context, _ Transaction) (Result, error) {
		return Result{}, ctx.Err()
	}); !errors.Is(err, context.Canceled) {
		t.Fatalf("Process context error = %v, want context.Canceled", err)
	}
}

func TestHandlerReturnsEncryptedResponse(t *testing.T) {
	client, uqpayContext := configuredTestClient(t)
	encrypted, err := uqpayContext.encrypt(authorizationRequest)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	handler, err := client.Handler(HandlerOptions{
		DecisionTimeout: 1500 * time.Millisecond,
		Decide: func(context.Context, Transaction) (Result, error) {
			return Result{ResponseCode: "51"}, nil
		},
	})
	if err != nil {
		t.Fatalf("Handler: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader(encrypted))
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
	plaintext, err := uqpayContext.decrypt(recorder.Body.String())
	if err != nil {
		t.Fatalf("decrypt response: %v", err)
	}
	if !strings.Contains(plaintext, `"response_code":"51"`) {
		t.Fatalf("response = %s", plaintext)
	}
}

func TestHandlerAbortsResponseOnError(t *testing.T) {
	client, _ := configuredTestClient(t)
	var received error
	handler, err := client.Handler(HandlerOptions{
		DecisionTimeout: 10 * time.Millisecond,
		Decide: func(ctx context.Context, _ Transaction) (Result, error) {
			<-ctx.Done()
			return Result{}, ctx.Err()
		},
		OnError: func(err error) { received = err },
	})
	if err != nil {
		t.Fatalf("Handler: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader("not encrypted"))
	assertAbortHandler(t, func() { handler.ServeHTTP(recorder, request) })
	if received == nil {
		t.Fatal("OnError was not called")
	}
	if recorder.Body.Len() != 0 {
		t.Fatalf("unexpected response body %q", recorder.Body.String())
	}
}

func TestHandlerEnforcesBodyLimitAndDecisionTimeout(t *testing.T) {
	client, uqpayContext := configuredTestClient(t)
	handler, err := client.Handler(HandlerOptions{
		MaxBodyBytes:    8,
		DecisionTimeout: time.Second,
		Decide: func(context.Context, Transaction) (Result, error) {
			return Result{ResponseCode: "00"}, nil
		},
	})
	if err != nil {
		t.Fatalf("Handler: %v", err)
	}
	assertAbortHandler(t, func() {
		handler.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader("too large")),
		)
	})

	encrypted, err := uqpayContext.encrypt(authorizationRequest)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	timedOut, err := client.Handler(HandlerOptions{
		DecisionTimeout: 5 * time.Millisecond,
		Decide: func(ctx context.Context, _ Transaction) (Result, error) {
			<-ctx.Done()
			return Result{}, ctx.Err()
		},
	})
	if err != nil {
		t.Fatalf("Handler: %v", err)
	}
	assertAbortHandler(t, func() {
		timedOut.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader(encrypted)),
		)
	})
}

func TestHandlerTimeoutDoesNotRequireCallbackToHonorContext(t *testing.T) {
	client, uqpayContext := configuredTestClient(t)
	encrypted, err := uqpayContext.encrypt(authorizationRequest)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	release := make(chan struct{})
	defer close(release)
	handler, err := client.Handler(HandlerOptions{
		DecisionTimeout: 10 * time.Millisecond,
		Decide: func(context.Context, Transaction) (Result, error) {
			<-release
			return Result{ResponseCode: "00"}, nil
		},
	})
	if err != nil {
		t.Fatalf("Handler: %v", err)
	}

	started := time.Now()
	assertAbortHandler(t, func() {
		handler.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader(encrypted)),
		)
	})
	if elapsed := time.Since(started); elapsed > 500*time.Millisecond {
		t.Fatalf("handler took %s after decision timeout", elapsed)
	}
}

func configuredTestClient(t *testing.T) (*Client, *pgpContext) {
	t.Helper()
	customer := mustGenerateKeyPair(t, "Customer", "customer@example.com")
	uqpay := mustGenerateKeyPair(t, "UQPAY", "issuing.tech@uqpay.com")
	client := NewClient()
	if err := client.Configure(Config{
		PrivateKey:     customer.PrivateKey,
		UQPayPublicKey: uqpay.PublicKey,
	}); err != nil {
		t.Fatalf("Configure: %v", err)
	}
	uqpayContext := mustNewPGPContext(t, Config{
		PrivateKey:     uqpay.PrivateKey,
		UQPayPublicKey: customer.PublicKey,
	})
	return client, uqpayContext
}

func assertAbortHandler(t *testing.T, call func()) {
	t.Helper()
	deferred := false
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				if recovered != http.ErrAbortHandler {
					t.Fatalf("panic = %v, want http.ErrAbortHandler", recovered)
				}
				deferred = true
			}
		}()
		call()
	}()
	if !deferred {
		t.Fatal("handler did not abort the response")
	}
}
