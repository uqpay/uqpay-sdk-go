package authdecision

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

var (
	testKeysOnce    sync.Once
	testCustomerKP  *KeyPair
	testUQPayKP     *KeyPair
	testKeysInitErr error
)

func initTestKeys(t *testing.T) {
	t.Helper()
	testKeysOnce.Do(func() {
		testCustomerKP, testKeysInitErr = GenerateKeyPair("Customer", "customer@test.com")
		if testKeysInitErr != nil {
			return
		}
		testUQPayKP, testKeysInitErr = GenerateKeyPair("UQPay", "uqpay@test.com")
	})
	if testKeysInitErr != nil {
		t.Fatalf("failed to generate test keys: %v", testKeysInitErr)
	}
}

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair("Test User", "test@example.com")
	if err != nil {
		t.Fatalf("GenerateKeyPair returned error: %v", err)
	}

	if kp.PublicKey == "" {
		t.Fatal("PublicKey is empty")
	}
	if kp.PrivateKey == "" {
		t.Fatal("PrivateKey is empty")
	}

	// Verify keys are parseable
	pubKey, err := crypto.NewKeyFromArmored(kp.PublicKey)
	if err != nil {
		t.Fatalf("failed to parse public key: %v", err)
	}
	if pubKey.IsPrivate() {
		t.Fatal("public key should not be private")
	}

	privKey, err := crypto.NewKeyFromArmored(kp.PrivateKey)
	if err != nil {
		t.Fatalf("failed to parse private key: %v", err)
	}
	if !privKey.IsPrivate() {
		t.Fatal("private key should be private")
	}
}

func TestPgpContextEncryptDecrypt(t *testing.T) {
	initTestKeys(t)

	// Simulate UQPAY encrypting with customer's public key, customer decrypting with their private key
	// UQPAY side: encrypt to customer
	uqpayCtx, err := newPgpContext(Config{
		PrivateKey:     testUQPayKP.PrivateKey,
		UQPayPublicKey: testCustomerKP.PublicKey,
	})
	if err != nil {
		t.Fatalf("failed to create UQPAY pgpContext: %v", err)
	}

	plaintext := `{"transaction_id":"tx-001","response_code":"00"}`
	encrypted, err := uqpayCtx.encrypt(plaintext)
	if err != nil {
		t.Fatalf("UQPAY encrypt failed: %v", err)
	}

	// Customer side: decrypt with their private key
	customerCtx, err := newPgpContext(Config{
		PrivateKey:     testCustomerKP.PrivateKey,
		UQPayPublicKey: testUQPayKP.PublicKey,
	})
	if err != nil {
		t.Fatalf("failed to create customer pgpContext: %v", err)
	}

	decrypted, err := customerCtx.decrypt(encrypted)
	if err != nil {
		t.Fatalf("customer decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf("roundtrip mismatch: got %q, want %q", decrypted, plaintext)
	}

	// Reverse direction: customer encrypts, UQPAY decrypts
	encrypted2, err := customerCtx.encrypt("hello from customer")
	if err != nil {
		t.Fatalf("customer encrypt failed: %v", err)
	}

	decrypted2, err := uqpayCtx.decrypt(encrypted2)
	if err != nil {
		t.Fatalf("UQPAY decrypt failed: %v", err)
	}

	if decrypted2 != "hello from customer" {
		t.Fatalf("reverse roundtrip mismatch: got %q", decrypted2)
	}
}

func TestPgpContextInvalidKeys(t *testing.T) {
	_, err := newPgpContext(Config{
		PrivateKey:     "not a valid key",
		UQPayPublicKey: "also not valid",
	})
	if err == nil {
		t.Fatal("expected error for invalid keys, got nil")
	}
}

func TestResolveKeyFromFile(t *testing.T) {
	initTestKeys(t)

	dir := t.TempDir()

	privPath := filepath.Join(dir, "private.asc")
	pubPath := filepath.Join(dir, "public.pgp")

	if err := os.WriteFile(privPath, []byte(testCustomerKP.PrivateKey), 0600); err != nil {
		t.Fatalf("failed to write private key file: %v", err)
	}
	if err := os.WriteFile(pubPath, []byte(testUQPayKP.PublicKey), 0644); err != nil {
		t.Fatalf("failed to write public key file: %v", err)
	}

	// Create pgpContext using file paths
	ctx, err := newPgpContext(Config{
		PrivateKey:     privPath,
		UQPayPublicKey: pubPath,
	})
	if err != nil {
		t.Fatalf("newPgpContext with file paths failed: %v", err)
	}

	// Verify roundtrip works
	plaintext := "file-based key test"
	encrypted, err := ctx.encrypt(plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// To decrypt we need a context where the private key can decrypt
	// ctx encrypts to UQPay's public key, so we need UQPay's private key to decrypt
	uqpayCtx, err := newPgpContext(Config{
		PrivateKey:     testUQPayKP.PrivateKey,
		UQPayPublicKey: testCustomerKP.PublicKey,
	})
	if err != nil {
		t.Fatalf("failed to create UQPAY pgpContext: %v", err)
	}

	decrypted, err := uqpayCtx.decrypt(encrypted)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf("file-based roundtrip mismatch: got %q, want %q", decrypted, plaintext)
	}
}
