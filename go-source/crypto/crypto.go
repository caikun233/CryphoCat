// Package crypto provides RSA key generation, encryption, and decryption.
package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// GenerateKeyPair creates a new RSA key pair of the given bit size.
// If passphrase is non-empty the private key PEM is encrypted with AES-256-CBC
// using the legacy pem.EncryptPEMBlock mechanism (compatible with OpenSSL 1.x).
// The unencrypted form uses PKCS#8 (header: PRIVATE KEY).
func GenerateKeyPair(keySize int, passphrase []byte, privPath, pubPath string) error {
	if keySize < 2048 {
		return fmt.Errorf("key size must be at least 2048 bits (got %d)", keySize)
	}
	priv, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("key generation: %w", err)
	}

	// Marshal private key as PKCS#8 DER.
	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("marshal private key: %w", err)
	}

	var privBlock *pem.Block
	if len(passphrase) > 0 {
		privBlock, err = x509.EncryptPEMBlock( //nolint:staticcheck // intentional legacy use
			rand.Reader,
			"ENCRYPTED PRIVATE KEY",
			privDER,
			passphrase,
			x509.PEMCipherAES256,
		)
		if err != nil {
			return fmt.Errorf("encrypt private key: %w", err)
		}
	} else {
		privBlock = &pem.Block{Type: "PRIVATE KEY", Bytes: privDER}
	}

	// Marshal public key as PKIX/SubjectPublicKeyInfo DER.
	pubDER, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return fmt.Errorf("marshal public key: %w", err)
	}
	pubBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}

	if err := writeFile(privPath, pem.EncodeToMemory(privBlock), 0o600); err != nil {
		return err
	}
	return writeFile(pubPath, pem.EncodeToMemory(pubBlock), 0o644)
}

// LoadPublicKey parses a PEM-encoded RSA public key from the given file.
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read public key: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found in %s", path)
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}
	pub, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return pub, nil
}

// LoadPrivateKey parses a PEM-encoded RSA private key (PKCS#8, optionally
// encrypted) from the given file.
func LoadPrivateKey(path string, passphrase []byte) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found in %s", path)
	}
	var der []byte
	if x509.IsEncryptedPEMBlock(block) { //nolint:staticcheck
		der, err = x509.DecryptPEMBlock(block, passphrase) //nolint:staticcheck
		if err != nil {
			return nil, fmt.Errorf("decrypt private key: %w", err)
		}
	} else {
		der = block.Bytes
	}
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}
	return priv, nil
}

// Encrypt encrypts plaintext with the given RSA public key using OAEP-SHA256.
// The plaintext is split into blocks that fit within the key modulus.
// Returns the ciphertext as a concatenation of fixed-size RSA blocks.
func Encrypt(pub *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	hash := sha256.New()
	// max plaintext per block = modulus size - 2*hashLen - 2
	maxBlock := pub.Size() - 2*hash.Size() - 2
	var ct []byte
	for len(plaintext) > 0 {
		chunk := plaintext
		if len(chunk) > maxBlock {
			chunk = plaintext[:maxBlock]
		}
		plaintext = plaintext[len(chunk):]
		block, err := rsa.EncryptOAEP(hash, rand.Reader, pub, chunk, nil)
		if err != nil {
			return nil, fmt.Errorf("encrypt block: %w", err)
		}
		ct = append(ct, block...)
	}
	return ct, nil
}

// Decrypt decrypts ciphertext produced by Encrypt using the given RSA private key.
func Decrypt(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	hash := sha256.New()
	segSize := priv.Size()
	var pt []byte
	for len(ciphertext) > 0 {
		if len(ciphertext) < segSize {
			return nil, fmt.Errorf("ciphertext length %d is not a multiple of block size %d",
				len(ciphertext), segSize)
		}
		block := ciphertext[:segSize]
		ciphertext = ciphertext[segSize:]
		plain, err := rsa.DecryptOAEP(hash, rand.Reader, priv, block, nil)
		if err != nil {
			return nil, fmt.Errorf("decrypt block: %w", err)
		}
		pt = append(pt, plain...)
	}
	return pt, nil
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
