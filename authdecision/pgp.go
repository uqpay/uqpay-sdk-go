package authdecision

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

const maxDecryptedBodyBytes = 1 << 20

func pgpPacketConfig() *packet.Config {
	return &packet.Config{
		DefaultHash:            crypto.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
		RSABits:                2048,
		MinRSABits:             2048,
	}
}

// GenerateKeyPair generates an ASCII-armored RSA 2048 key pair with an
// encryption subkey, matching UQPAY's authorization decision requirements.
func GenerateKeyPair(name, email string) (*KeyPair, error) {
	if strings.TrimSpace(name) == "" && strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("authdecision: key name or email is required")
	}
	entity, err := openpgp.NewEntity(name, "", email, pgpPacketConfig())
	if err != nil {
		return nil, fmt.Errorf("authdecision: generate RSA key pair: %w", err)
	}

	privateKey, err := serializePrivateKey(entity)
	if err != nil {
		return nil, err
	}
	publicKey, err := serializePublicKey(entity)
	if err != nil {
		return nil, err
	}
	return &KeyPair{PublicKey: publicKey, PrivateKey: privateKey}, nil
}

type pgpContext struct {
	privateKeys openpgp.EntityList
	publicKeys  openpgp.EntityList
}

func newPGPContext(config Config) (*pgpContext, error) {
	privateArmored, err := resolveKey(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("authdecision: resolve private key: %w", err)
	}
	publicArmored, err := resolveKey(config.UQPayPublicKey)
	if err != nil {
		return nil, fmt.Errorf("authdecision: resolve UQPAY public key: %w", err)
	}

	privateKeys, err := openpgp.ReadArmoredKeyRing(strings.NewReader(privateArmored))
	if err != nil {
		return nil, fmt.Errorf("authdecision: parse private key: %w", err)
	}
	if err := unlockPrivateKeys(privateKeys, []byte(config.Passphrase)); err != nil {
		return nil, err
	}
	if !hasRSADecryptionKey(privateKeys) {
		return nil, fmt.Errorf("authdecision: private key has no RSA decryption key of at least 2048 bits")
	}

	publicKeys, err := openpgp.ReadArmoredKeyRing(strings.NewReader(publicArmored))
	if err != nil {
		return nil, fmt.Errorf("authdecision: parse UQPAY public key: %w", err)
	}
	if !hasRSAEncryptionKey(publicKeys, time.Now()) {
		return nil, fmt.Errorf("authdecision: UQPAY public key has no RSA encryption key of at least 2048 bits")
	}

	return &pgpContext{privateKeys: privateKeys, publicKeys: publicKeys}, nil
}

func hasRSADecryptionKey(entities openpgp.EntityList) bool {
	for _, entity := range entities {
		if entity.PrivateKey != nil && entity.PrimaryKey.PubKeyAlgo.CanEncrypt() && isRSA2048(entity.PrimaryKey) {
			return true
		}
		for _, subkey := range entity.Subkeys {
			if subkey.PrivateKey != nil && subkey.PublicKey.PubKeyAlgo.CanEncrypt() && isRSA2048(subkey.PublicKey) {
				return true
			}
		}
	}
	return false
}

func hasRSAEncryptionKey(entities openpgp.EntityList, now time.Time) bool {
	for _, entity := range entities {
		key, ok := entity.EncryptionKey(now)
		if ok && isRSA2048(key.PublicKey) {
			return true
		}
	}
	return false
}

func isRSA2048(publicKey *packet.PublicKey) bool {
	key, ok := publicKey.PublicKey.(*rsa.PublicKey)
	return ok && key.N.BitLen() >= 2048
}

func (p *pgpContext) decrypt(ciphertext string) (string, error) {
	block, err := armor.Decode(strings.NewReader(ciphertext))
	if err != nil {
		return "", fmt.Errorf("authdecision: decode armored request: %w", err)
	}
	if block.Type != "PGP MESSAGE" {
		return "", fmt.Errorf("authdecision: unexpected armor type %q", block.Type)
	}
	message, err := openpgp.ReadMessage(block.Body, p.privateKeys, nil, pgpPacketConfig())
	if err != nil {
		return "", fmt.Errorf("authdecision: decrypt request: %w", err)
	}
	plaintext, err := io.ReadAll(io.LimitReader(message.UnverifiedBody, maxDecryptedBodyBytes+1))
	if err != nil {
		return "", fmt.Errorf("authdecision: read decrypted request: %w", err)
	}
	if len(plaintext) > maxDecryptedBodyBytes {
		return "", fmt.Errorf("authdecision: decrypted request exceeds %d bytes", maxDecryptedBodyBytes)
	}
	return string(plaintext), nil
}

func (p *pgpContext) encrypt(plaintext string) (string, error) {
	var output bytes.Buffer
	armoredWriter, err := armor.Encode(&output, "PGP MESSAGE", nil)
	if err != nil {
		return "", fmt.Errorf("authdecision: create armored response: %w", err)
	}
	plaintextWriter, err := openpgp.Encrypt(armoredWriter, p.publicKeys, nil, nil, pgpPacketConfig())
	if err != nil {
		_ = armoredWriter.Close()
		return "", fmt.Errorf("authdecision: encrypt response: %w", err)
	}
	if _, err := io.WriteString(plaintextWriter, plaintext); err != nil {
		_ = plaintextWriter.Close()
		_ = armoredWriter.Close()
		return "", fmt.Errorf("authdecision: write encrypted response: %w", err)
	}
	if err := plaintextWriter.Close(); err != nil {
		_ = armoredWriter.Close()
		return "", fmt.Errorf("authdecision: finalize encrypted response: %w", err)
	}
	if err := armoredWriter.Close(); err != nil {
		return "", fmt.Errorf("authdecision: finalize armored response: %w", err)
	}
	return output.String(), nil
}

func resolveKey(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("key is required")
	}
	lower := strings.ToLower(trimmed)
	if strings.HasSuffix(lower, ".asc") || strings.HasSuffix(lower, ".pgp") || strings.HasSuffix(lower, ".gpg") {
		data, err := os.ReadFile(trimmed)
		if err != nil {
			return "", fmt.Errorf("read key file %q: %w", trimmed, err)
		}
		return string(data), nil
	}
	return value, nil
}

func unlockPrivateKeys(entities openpgp.EntityList, passphrase []byte) error {
	foundPrivateKey := false
	for _, entity := range entities {
		privateKeys := []*packet.PrivateKey{entity.PrivateKey}
		for _, subkey := range entity.Subkeys {
			privateKeys = append(privateKeys, subkey.PrivateKey)
		}
		for _, privateKey := range privateKeys {
			if privateKey == nil {
				continue
			}
			foundPrivateKey = true
			if !privateKey.Encrypted {
				continue
			}
			if len(passphrase) == 0 {
				return fmt.Errorf("authdecision: private key is passphrase-protected")
			}
			if err := privateKey.Decrypt(passphrase); err != nil {
				return fmt.Errorf("authdecision: unlock private key: %w", err)
			}
		}
	}
	if !foundPrivateKey {
		return fmt.Errorf("authdecision: configured private key contains no private material")
	}
	return nil
}

func serializePrivateKey(entity *openpgp.Entity) (string, error) {
	var output bytes.Buffer
	armoredWriter, err := armor.Encode(&output, openpgp.PrivateKeyType, nil)
	if err != nil {
		return "", fmt.Errorf("authdecision: armor private key: %w", err)
	}
	if err := entity.SerializePrivate(armoredWriter, pgpPacketConfig()); err != nil {
		_ = armoredWriter.Close()
		return "", fmt.Errorf("authdecision: serialize private key: %w", err)
	}
	if err := armoredWriter.Close(); err != nil {
		return "", fmt.Errorf("authdecision: finalize private key: %w", err)
	}
	return output.String(), nil
}

func serializePublicKey(entity *openpgp.Entity) (string, error) {
	var output bytes.Buffer
	armoredWriter, err := armor.Encode(&output, openpgp.PublicKeyType, nil)
	if err != nil {
		return "", fmt.Errorf("authdecision: armor public key: %w", err)
	}
	if err := entity.Serialize(armoredWriter); err != nil {
		_ = armoredWriter.Close()
		return "", fmt.Errorf("authdecision: serialize public key: %w", err)
	}
	if err := armoredWriter.Close(); err != nil {
		return "", fmt.Errorf("authdecision: finalize public key: %w", err)
	}
	return output.String(), nil
}
