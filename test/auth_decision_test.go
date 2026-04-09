package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	uqpay "github.com/uqpay/uqpay-sdk-go"
	"github.com/uqpay/uqpay-sdk-go/authdecision"
	"github.com/uqpay/uqpay-sdk-go/configuration"
)

// encryptForCustomer encrypts plaintext with the customer's public key (simulates UQPAY sending).
func encryptForCustomer(t *testing.T, customerPubKeyArmored, plaintext string) string {
	t.Helper()
	pubKey, err := crypto.NewKeyFromArmored(customerPubKeyArmored)
	if err != nil {
		t.Fatalf("parse customer public key: %v", err)
	}
	encHandle, err := crypto.PGP().Encryption().Recipient(pubKey).New()
	if err != nil {
		t.Fatalf("create encryption handle: %v", err)
	}
	pgpMsg, err := encHandle.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	armored, err := pgpMsg.Armor()
	if err != nil {
		t.Fatalf("armor: %v", err)
	}
	return armored
}

// decryptAsUqpayServer decrypts ciphertext with UQPAY's private key (simulates UQPAY receiving).
func decryptAsUqpayServer(t *testing.T, uqpayPrivKeyArmored, ciphertext string) map[string]interface{} {
	t.Helper()
	privKey, err := crypto.NewKeyFromArmored(uqpayPrivKeyArmored)
	if err != nil {
		t.Fatalf("parse uqpay private key: %v", err)
	}
	privKeyCopy, err := privKey.Copy()
	if err != nil {
		t.Fatalf("copy private key: %v", err)
	}
	decHandle, err := crypto.PGP().Decryption().DecryptionKey(privKeyCopy).New()
	if err != nil {
		t.Fatalf("create decryption handle: %v", err)
	}
	defer decHandle.ClearPrivateParams()

	result, err := decHandle.Decrypt([]byte(ciphertext), crypto.Armor)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(result.Bytes(), &m); err != nil {
		t.Fatalf("unmarshal decrypted response: %v", err)
	}
	return m
}

func TestAuthDecisionIntegration(t *testing.T) {
	// 1. Generate key pairs
	customerKeys, err := authdecision.GenerateKeyPair("Test Customer", "test@customer.com")
	if err != nil {
		t.Fatalf("generate customer keys: %v", err)
	}
	uqpayKeys, err := authdecision.GenerateKeyPair("UQPAY", "issuing.tech@uqpay.com")
	if err != nil {
		t.Fatalf("generate uqpay keys: %v", err)
	}

	// 2. Create SDK client and configure
	client, err := uqpay.NewClient("test", "test", configuration.Sandbox())
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	err = client.Issuing.AuthDecision.Configure(authdecision.Config{
		PrivateKey:     customerKeys.PrivateKey,
		UQPayPublicKey: uqpayKeys.PublicKey,
	})
	if err != nil {
		t.Fatalf("configure: %v", err)
	}

	// 3. Create handler: decline if amount > 5000
	handler := client.Issuing.AuthDecision.Handler(authdecision.HandlerOptions{
		Decide: func(ctx context.Context, tx authdecision.Transaction) (authdecision.Result, error) {
			t.Logf("Received: tx=%s amount=%.2f %s merchant=%s",
				tx.TransactionID, tx.BillingAmount, tx.BillingCurrencyCode, tx.MerchantName)
			if tx.BillingAmount > 5000 {
				t.Log("Decision: DECLINE")
				return authdecision.Result{ResponseCode: "51"}, nil
			}
			t.Log("Decision: APPROVE")
			return authdecision.Result{ResponseCode: "00", PartnerReferenceID: "test-ref-001"}, nil
		},
		OnError: func(err error) {
			t.Logf("OnError: %v", err)
		},
	})

	baseTx := map[string]interface{}{
		"transaction_id":                    "tx-int-001",
		"transaction_type":                  1000,
		"card_id":                           "card-001",
		"processing_code":                   "00",
		"billing_amount":                    100.50,
		"transaction_amount":                100.50,
		"auth_amount":                       100.50,
		"date_of_transaction":               "2026-04-09 15:00:00",
		"billing_currency_code":             "SGD",
		"transaction_currency_code":         "SGD",
		"auth_currency_code":                "SGD",
		"card_balance":                      500.00,
		"merchant_id":                       "MERCHANT123",
		"merchant_name":                     "Coffee Shop",
		"merchant_category_code":            "5814",
		"merchant_city":                     "Singapore",
		"merchant_country":                  "SG",
		"terminal_id":                       "TERM001",
		"pos_entry_mode":                    "07",
		"pos_condition_code":                "00",
		"pin_entry_capability":              "1",
		"retrieval_reference_number":        "123456789012",
		"system_trace_audit_number":         "123456",
		"acquiring_institution_country_code": "SG",
		"acquiring_institution_id":          "ACQ001",
		"wallet_type":                       "APPLE",
	}

	// Test: Approve (amount=100.50)
	t.Run("Approve", func(t *testing.T) {
		data, _ := json.Marshal(baseTx)
		encrypted := encryptForCustomer(t, customerKeys.PublicKey, string(data))

		req := httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader(encrypted))
		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != 200 {
			t.Fatalf("status: %d", w.Code)
		}

		result := decryptAsUqpayServer(t, uqpayKeys.PrivateKey, w.Body.String())
		if result["transaction_id"] != "tx-int-001" {
			t.Errorf("transaction_id: %v", result["transaction_id"])
		}
		if result["response_code"] != "00" {
			t.Errorf("response_code: %v (want 00)", result["response_code"])
		}
		if result["partner_reference_id"] != "test-ref-001" {
			t.Errorf("partner_reference_id: %v", result["partner_reference_id"])
		}
	})

	// Test: Decline (amount=9999.99)
	t.Run("Decline", func(t *testing.T) {
		tx := make(map[string]interface{})
		for k, v := range baseTx {
			tx[k] = v
		}
		tx["transaction_id"] = "tx-int-002"
		tx["billing_amount"] = 9999.99

		data, _ := json.Marshal(tx)
		encrypted := encryptForCustomer(t, customerKeys.PublicKey, string(data))

		req := httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader(encrypted))
		w := httptest.NewRecorder()
		handler(w, req)

		if w.Code != 200 {
			t.Fatalf("status: %d", w.Code)
		}

		result := decryptAsUqpayServer(t, uqpayKeys.PrivateKey, w.Body.String())
		if result["transaction_id"] != "tx-int-002" {
			t.Errorf("transaction_id: %v", result["transaction_id"])
		}
		if result["response_code"] != "51" {
			t.Errorf("response_code: %v (want 51)", result["response_code"])
		}
	})

	// Test: Malformed body — no response
	t.Run("MalformedBody", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/auth-decision", strings.NewReader("not encrypted"))
		w := httptest.NewRecorder()
		handler(w, req)

		if w.Body.Len() > 0 {
			t.Error("expected empty response on decryption failure")
		}
	})
}
