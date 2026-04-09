package authdecision

import (
	"fmt"
	"os"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

// GenerateKeyPair generates a PGP key pair with an encryption subkey and returns
// the armored public and private keys. Compatible with UQPAY's server-side Go openpgp.
func GenerateKeyPair(name, email string) (*KeyPair, error) {
	pgpHandle := crypto.PGP()
	key, err := pgpHandle.KeyGeneration().
		AddUserId(name, email).
		New().
		GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("authdecision: failed to generate key: %w", err)
	}

	privateArmored, err := key.Armor()
	if err != nil {
		return nil, fmt.Errorf("authdecision: failed to armor private key: %w", err)
	}

	publicArmored, err := key.GetArmoredPublicKey()
	if err != nil {
		return nil, fmt.Errorf("authdecision: failed to armor public key: %w", err)
	}

	return &KeyPair{
		PublicKey:  publicArmored,
		PrivateKey: privateArmored,
	}, nil
}

// pgpContext holds parsed PGP keys for encrypt/decrypt operations.
type pgpContext struct {
	privateKey *crypto.Key
	publicKey  *crypto.Key
}

// newPgpContext creates a new pgpContext from the given Config.
func newPgpContext(config Config) (*pgpContext, error) {
	privArmored, err := resolveKey(config.PrivateKey)
	if err != nil {
		return nil, err
	}
	pubArmored, err := resolveKey(config.UQPayPublicKey)
	if err != nil {
		return nil, err
	}

	var privKey *crypto.Key
	if config.Passphrase != "" {
		privKey, err = crypto.NewPrivateKeyFromArmored(privArmored, []byte(config.Passphrase))
	} else {
		privKey, err = crypto.NewKeyFromArmored(privArmored)
	}
	if err != nil {
		return nil, fmt.Errorf("authdecision: failed to parse private key: %w", err)
	}

	pubKey, err := crypto.NewKeyFromArmored(pubArmored)
	if err != nil {
		return nil, fmt.Errorf("authdecision: failed to parse public key: %w", err)
	}

	return &pgpContext{
		privateKey: privKey,
		publicKey:  pubKey,
	}, nil
}

// decrypt decrypts an armored PGP message using the private key.
func (p *pgpContext) decrypt(ciphertext string) (string, error) {
	pgpHandle := crypto.PGP()
	// Copy the private key to avoid ClearPrivateParams destroying the original.
	privKeyCopy, err := p.privateKey.Copy()
	if err != nil {
		return "", fmt.Errorf("authdecision: failed to copy private key: %w", err)
	}

	decHandle, err := pgpHandle.Decryption().
		DecryptionKey(privKeyCopy).
		New()
	if err != nil {
		return "", fmt.Errorf("authdecision: failed to create decryption handle: %w", err)
	}
	defer decHandle.ClearPrivateParams()

	result, err := decHandle.Decrypt([]byte(ciphertext), crypto.Armor)
	if err != nil {
		return "", fmt.Errorf("authdecision: failed to decrypt message: %w", err)
	}

	return result.String(), nil
}

// encrypt encrypts plaintext to an armored PGP message using the public key.
func (p *pgpContext) encrypt(plaintext string) (string, error) {
	pgpHandle := crypto.PGP()
	encHandle, err := pgpHandle.Encryption().
		Recipient(p.publicKey).
		New()
	if err != nil {
		return "", fmt.Errorf("authdecision: failed to create encryption handle: %w", err)
	}

	pgpMessage, err := encHandle.Encrypt([]byte(plaintext))
	if err != nil {
		return "", fmt.Errorf("authdecision: failed to encrypt message: %w", err)
	}

	armored, err := pgpMessage.Armor()
	if err != nil {
		return "", fmt.Errorf("authdecision: failed to armor encrypted message: %w", err)
	}

	return armored, nil
}

// resolveKey returns the key value as-is, or reads it from a file if the value
// ends with .asc, .pgp, or .gpg.
func resolveKey(value string) (string, error) {
	lower := strings.ToLower(value)
	if strings.HasSuffix(lower, ".asc") || strings.HasSuffix(lower, ".pgp") || strings.HasSuffix(lower, ".gpg") {
		data, err := os.ReadFile(value)
		if err != nil {
			return "", fmt.Errorf("authdecision: failed to read key file %s: %w", value, err)
		}
		return string(data), nil
	}
	return value, nil
}
