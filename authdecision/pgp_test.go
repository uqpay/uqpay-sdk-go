package authdecision

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

func TestGenerateKeyPairCreatesRSA2048WithEncryptionSubkey(t *testing.T) {
	keys, err := GenerateKeyPair("Test Customer", "test@example.com")
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	if !strings.Contains(keys.PublicKey, "BEGIN PGP PUBLIC KEY BLOCK") {
		t.Fatal("public key is not armored")
	}
	if !strings.Contains(keys.PrivateKey, "BEGIN PGP PRIVATE KEY BLOCK") {
		t.Fatal("private key is not armored")
	}

	entities, err := openpgp.ReadArmoredKeyRing(strings.NewReader(keys.PublicKey))
	if err != nil {
		t.Fatalf("parse public key: %v", err)
	}
	if len(entities) != 1 {
		t.Fatalf("entities = %d, want 1", len(entities))
	}
	primary, ok := entities[0].PrimaryKey.PublicKey.(*rsa.PublicKey)
	if !ok {
		t.Fatalf("primary key type = %T, want *rsa.PublicKey", entities[0].PrimaryKey.PublicKey)
	}
	if got := primary.N.BitLen(); got != 2048 {
		t.Fatalf("RSA bits = %d, want 2048", got)
	}
	if len(entities[0].Subkeys) == 0 {
		t.Fatal("generated key has no encryption subkey")
	}
}

func TestPGPContextRoundTrip(t *testing.T) {
	customer := mustGenerateKeyPair(t, "Customer", "customer@example.com")
	uqpay := mustGenerateKeyPair(t, "UQPAY", "issuing.tech@uqpay.com")

	customerContext := mustNewPGPContext(t, Config{
		PrivateKey:     customer.PrivateKey,
		UQPayPublicKey: uqpay.PublicKey,
	})
	uqpayContext := mustNewPGPContext(t, Config{
		PrivateKey:     uqpay.PrivateKey,
		UQPayPublicKey: customer.PublicKey,
	})

	request := `{"transaction_id":"tx-123","billing_amount":"2.31"}`
	encryptedRequest, err := uqpayContext.encrypt(request)
	if err != nil {
		t.Fatalf("encrypt request: %v", err)
	}
	decryptedRequest, err := customerContext.decrypt(encryptedRequest)
	if err != nil {
		t.Fatalf("decrypt request: %v", err)
	}
	if decryptedRequest != request {
		t.Fatalf("decrypted request = %q, want %q", decryptedRequest, request)
	}

	response := `{"transaction_id":"tx-123","response_code":"00"}`
	encryptedResponse, err := customerContext.encrypt(response)
	if err != nil {
		t.Fatalf("encrypt response: %v", err)
	}
	decryptedResponse, err := uqpayContext.decrypt(encryptedResponse)
	if err != nil {
		t.Fatalf("decrypt response: %v", err)
	}
	if decryptedResponse != response {
		t.Fatalf("decrypted response = %q, want %q", decryptedResponse, response)
	}
}

func TestPGPContextRejectsTamperedCiphertext(t *testing.T) {
	customer := mustGenerateKeyPair(t, "Customer", "customer@example.com")
	uqpay := mustGenerateKeyPair(t, "UQPAY", "issuing.tech@uqpay.com")
	context := mustNewPGPContext(t, Config{
		PrivateKey:     customer.PrivateKey,
		UQPayPublicKey: uqpay.PublicKey,
	})
	sender := mustNewPGPContext(t, Config{
		PrivateKey:     uqpay.PrivateKey,
		UQPayPublicKey: customer.PublicKey,
	})
	encrypted, err := sender.encrypt("integrity protected message")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	block, err := armor.Decode(strings.NewReader(encrypted))
	if err != nil {
		t.Fatalf("decode armor: %v", err)
	}
	payload, err := io.ReadAll(block.Body)
	if err != nil {
		t.Fatalf("read payload: %v", err)
	}
	if len(payload) < 32 {
		t.Fatalf("encrypted payload too short: %d", len(payload))
	}
	payload[len(payload)-16] ^= 0x01

	var tampered bytes.Buffer
	armoredWriter, err := armor.Encode(&tampered, "PGP MESSAGE", nil)
	if err != nil {
		t.Fatalf("armor tampered payload: %v", err)
	}
	if _, err := armoredWriter.Write(payload); err != nil {
		t.Fatalf("write tampered payload: %v", err)
	}
	if err := armoredWriter.Close(); err != nil {
		t.Fatalf("close tampered armor: %v", err)
	}
	if _, err := context.decrypt(tampered.String()); err == nil {
		t.Fatal("expected integrity error for tampered ciphertext")
	}
}

func TestPGPContextAcceptsPassphraseProtectedPrivateKey(t *testing.T) {
	customer := mustGenerateKeyPair(t, "Customer", "customer@example.com")
	uqpay := mustGenerateKeyPair(t, "UQPAY", "issuing.tech@uqpay.com")
	passphrase := []byte("test-passphrase")
	protected := protectPrivateKey(t, customer.PrivateKey, passphrase)

	ctx, err := newPGPContext(Config{
		PrivateKey:     protected,
		UQPayPublicKey: uqpay.PublicKey,
		Passphrase:     string(passphrase),
	})
	if err != nil {
		t.Fatalf("newPGPContext: %v", err)
	}

	uqpayContext := mustNewPGPContext(t, Config{
		PrivateKey:     uqpay.PrivateKey,
		UQPayPublicKey: customer.PublicKey,
	})
	encrypted, err := uqpayContext.encrypt("protected message")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	decrypted, err := ctx.decrypt(encrypted)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if decrypted != "protected message" {
		t.Fatalf("decrypted = %q", decrypted)
	}
}

func TestPGPContextAcceptsKeyFilePaths(t *testing.T) {
	customer := mustGenerateKeyPair(t, "Customer", "customer@example.com")
	uqpay := mustGenerateKeyPair(t, "UQPAY", "issuing.tech@uqpay.com")
	dir := t.TempDir()
	privatePath := filepath.Join(dir, "customer-private.asc")
	publicPath := filepath.Join(dir, "uqpay-public.pgp")
	if err := os.WriteFile(privatePath, []byte(customer.PrivateKey), 0o600); err != nil {
		t.Fatalf("write private key: %v", err)
	}
	if err := os.WriteFile(publicPath, []byte(uqpay.PublicKey), 0o600); err != nil {
		t.Fatalf("write public key: %v", err)
	}

	if _, err := newPGPContext(Config{
		PrivateKey:     privatePath,
		UQPayPublicKey: publicPath,
	}); err != nil {
		t.Fatalf("newPGPContext with files: %v", err)
	}
}

func TestPGPContextRejectsInvalidKeys(t *testing.T) {
	keys := mustGenerateKeyPair(t, "UQPAY", "issuing.tech@uqpay.com")
	if _, err := newPGPContext(Config{
		PrivateKey:     "not-a-key",
		UQPayPublicKey: keys.PublicKey,
	}); err == nil {
		t.Fatal("expected invalid private key error")
	}
	if _, err := newPGPContext(Config{
		PrivateKey:     keys.PrivateKey,
		UQPayPublicKey: "not-a-key",
	}); err == nil {
		t.Fatal("expected invalid public key error")
	}
}

func TestPGPContextRejectsRSAKeysBelow2048Bits(t *testing.T) {
	weak, err := openpgp.NewEntity("Weak", "", "weak@example.com", &packet.Config{
		DefaultHash: crypto.SHA256,
		RSABits:     1024,
	})
	if err != nil {
		t.Fatalf("generate weak key: %v", err)
	}
	weakPrivate, err := serializePrivateKey(weak)
	if err != nil {
		t.Fatalf("serialize weak private key: %v", err)
	}
	weakPublic, err := serializePublicKey(weak)
	if err != nil {
		t.Fatalf("serialize weak public key: %v", err)
	}
	strong := mustGenerateKeyPair(t, "Strong", "strong@example.com")

	if _, err := newPGPContext(Config{
		PrivateKey:     weakPrivate,
		UQPayPublicKey: strong.PublicKey,
	}); err == nil {
		t.Fatal("expected weak private key to be rejected")
	}
	if _, err := newPGPContext(Config{
		PrivateKey:     strong.PrivateKey,
		UQPayPublicKey: weakPublic,
	}); err == nil {
		t.Fatal("expected weak public key to be rejected")
	}
}

func mustGenerateKeyPair(t *testing.T, name, email string) *KeyPair {
	t.Helper()
	keys, err := GenerateKeyPair(name, email)
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	return keys
}

func mustNewPGPContext(t *testing.T, config Config) *pgpContext {
	t.Helper()
	ctx, err := newPGPContext(config)
	if err != nil {
		t.Fatalf("newPGPContext: %v", err)
	}
	return ctx
}

func protectPrivateKey(t *testing.T, armored string, passphrase []byte) string {
	t.Helper()
	entities, err := openpgp.ReadArmoredKeyRing(strings.NewReader(armored))
	if err != nil {
		t.Fatalf("parse private key: %v", err)
	}
	for _, entity := range entities {
		if entity.PrivateKey != nil {
			if err := entity.PrivateKey.Encrypt(passphrase); err != nil {
				t.Fatalf("encrypt primary key: %v", err)
			}
		}
		for _, subkey := range entity.Subkeys {
			if subkey.PrivateKey != nil {
				if err := subkey.PrivateKey.Encrypt(passphrase); err != nil {
					t.Fatalf("encrypt subkey: %v", err)
				}
			}
		}
	}

	var out bytes.Buffer
	armoredWriter, err := armor.Encode(&out, openpgp.PrivateKeyType, nil)
	if err != nil {
		t.Fatalf("armor private key: %v", err)
	}
	for _, entity := range entities {
		if err := entity.SerializePrivateWithoutSigning(armoredWriter, nil); err != nil {
			t.Fatalf("serialize private key: %v", err)
		}
	}
	if err := armoredWriter.Close(); err != nil {
		t.Fatalf("close armor: %v", err)
	}
	return out.String()
}
