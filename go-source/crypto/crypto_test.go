package crypto_test

import (
	"os"
	"path/filepath"
	"testing"

	cc "github.com/caikun233/cryphocat/crypto"
)

func TestGenerateAndLoadKeyPair(t *testing.T) {
	dir := t.TempDir()
	privPath := filepath.Join(dir, "private.pem")
	pubPath := filepath.Join(dir, "public.pem")

	if err := cc.GenerateKeyPair(2048, nil, privPath, pubPath); err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	if _, err := os.Stat(privPath); err != nil {
		t.Fatalf("private key file missing: %v", err)
	}
	if _, err := os.Stat(pubPath); err != nil {
		t.Fatalf("public key file missing: %v", err)
	}

	if _, err := cc.LoadPublicKey(pubPath); err != nil {
		t.Fatalf("LoadPublicKey: %v", err)
	}
	if _, err := cc.LoadPrivateKey(privPath, nil); err != nil {
		t.Fatalf("LoadPrivateKey: %v", err)
	}
}

func TestGenerateKeyPairWithPassphrase(t *testing.T) {
	dir := t.TempDir()
	privPath := filepath.Join(dir, "private.pem")
	pubPath := filepath.Join(dir, "public.pem")

	passphrase := []byte("s3cr3t!")
	if err := cc.GenerateKeyPair(2048, passphrase, privPath, pubPath); err != nil {
		t.Fatalf("GenerateKeyPair with passphrase: %v", err)
	}

	// Wrong passphrase must fail.
	if _, err := cc.LoadPrivateKey(privPath, []byte("wrong")); err == nil {
		t.Fatal("expected error with wrong passphrase, got nil")
	}

	// Correct passphrase must succeed.
	if _, err := cc.LoadPrivateKey(privPath, passphrase); err != nil {
		t.Fatalf("LoadPrivateKey with correct passphrase: %v", err)
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	dir := t.TempDir()
	privPath := filepath.Join(dir, "private.pem")
	pubPath := filepath.Join(dir, "public.pem")

	if err := cc.GenerateKeyPair(2048, nil, privPath, pubPath); err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	pub, err := cc.LoadPublicKey(pubPath)
	if err != nil {
		t.Fatalf("LoadPublicKey: %v", err)
	}
	priv, err := cc.LoadPrivateKey(privPath, nil)
	if err != nil {
		t.Fatalf("LoadPrivateKey: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
	}{
		{"short", "Hello, CryphoCat!"},
		{"unicode", "你好，世界！这是一条加密消息。"},
		// Longer than one RSA-2048 block (>190 bytes) – tests segmentation.
		{"long", "ABCDEFGHIJ" + string(make([]byte, 300))},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ct, err := cc.Encrypt(pub, []byte(tc.plaintext))
			if err != nil {
				t.Fatalf("Encrypt: %v", err)
			}
			pt, err := cc.Decrypt(priv, ct)
			if err != nil {
				t.Fatalf("Decrypt: %v", err)
			}
			if string(pt) != tc.plaintext {
				t.Fatalf("round-trip mismatch: got %q, want %q", pt, tc.plaintext)
			}
		})
	}
}

func TestEncryptDecrypt4096(t *testing.T) {
	dir := t.TempDir()
	privPath := filepath.Join(dir, "private.pem")
	pubPath := filepath.Join(dir, "public.pem")

	if err := cc.GenerateKeyPair(4096, nil, privPath, pubPath); err != nil {
		t.Fatalf("GenerateKeyPair 4096: %v", err)
	}
	pub, _ := cc.LoadPublicKey(pubPath)
	priv, _ := cc.LoadPrivateKey(privPath, nil)

	plain := "RSA-4096 test message 测试"
	ct, err := cc.Encrypt(pub, []byte(plain))
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	pt, err := cc.Decrypt(priv, ct)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if string(pt) != plain {
		t.Fatalf("4096 round-trip mismatch")
	}
}
