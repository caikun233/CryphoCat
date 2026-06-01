package crypto_test

import (
	"os"
	"path/filepath"
	"testing"

	cc "github.com/caikun233/cryphocat/crypto"
)

func TestGenerateAndLoadKeyPair(t *testing.T) {
	for _, algo := range cc.AllAlgos() {
		t.Run(algo.String(), func(t *testing.T) {
			dir := t.TempDir()
			privPath := filepath.Join(dir, "private.pem")
			pubPath := filepath.Join(dir, "public.pem")

			if err := cc.GenerateKeyPair(algo, nil, privPath, pubPath); err != nil {
				t.Fatalf("GenerateKeyPair: %v", err)
			}
			if _, err := os.Stat(privPath); err != nil {
				t.Fatalf("private key file missing: %v", err)
			}
			if _, err := os.Stat(pubPath); err != nil {
				t.Fatalf("public key file missing: %v", err)
			}
			if _, _, err := cc.LoadPublicKey(pubPath); err != nil {
				t.Fatalf("LoadPublicKey: %v", err)
			}
			if _, _, err := cc.LoadPrivateKey(privPath, nil); err != nil {
				t.Fatalf("LoadPrivateKey: %v", err)
			}
		})
	}
}

func TestKeyPairWithPassphrase(t *testing.T) {
	for _, algo := range cc.AllAlgos() {
		t.Run(algo.String(), func(t *testing.T) {
			dir := t.TempDir()
			privPath := filepath.Join(dir, "private.pem")
			pubPath := filepath.Join(dir, "public.pem")

			passphrase := []byte("s3cr3t!")
			if err := cc.GenerateKeyPair(algo, passphrase, privPath, pubPath); err != nil {
				t.Fatalf("GenerateKeyPair: %v", err)
			}
			if _, _, err := cc.LoadPrivateKey(privPath, []byte("wrong")); err == nil {
				t.Fatal("expected error with wrong passphrase")
			}
			if _, _, err := cc.LoadPrivateKey(privPath, passphrase); err != nil {
				t.Fatalf("LoadPrivateKey: %v", err)
			}
		})
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
	}{
		{"short", "Hello, CryphoCat!"},
		{"unicode", "你好，世界！这是一条加密消息。"},
		{"long", "ABCDEFGHIJ" + string(make([]byte, 300))},
	}

	for _, algo := range cc.AllAlgos() {
		for _, tc := range tests {
			t.Run(algo.String()+"/"+tc.name, func(t *testing.T) {
				dir := t.TempDir()
				privPath := filepath.Join(dir, "private.pem")
				pubPath := filepath.Join(dir, "public.pem")

				if err := cc.GenerateKeyPair(algo, nil, privPath, pubPath); err != nil {
					t.Fatalf("GenerateKeyPair: %v", err)
				}
				pub, pubAlgo, err := cc.LoadPublicKey(pubPath)
				if err != nil {
					t.Fatalf("LoadPublicKey: %v", err)
				}
				priv, privAlgo, err := cc.LoadPrivateKey(privPath, nil)
				if err != nil {
					t.Fatalf("LoadPrivateKey: %v", err)
				}
				if pubAlgo != algo || privAlgo != algo {
					t.Fatalf("algo mismatch: pub=%v priv=%v want=%v", pubAlgo, privAlgo, algo)
				}

				ct, err := cc.Encrypt(pub, pubAlgo, []byte(tc.plaintext))
				if err != nil {
					t.Fatalf("Encrypt: %v", err)
				}
				pt, err := cc.Decrypt(priv, privAlgo, ct)
				if err != nil {
					t.Fatalf("Decrypt: %v", err)
				}
				if string(pt) != tc.plaintext {
					t.Fatalf("round-trip mismatch: got %q, want %q", pt, tc.plaintext)
				}
			})
		}
	}
}
