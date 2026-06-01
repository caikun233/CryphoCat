// Package crypto provides multi-algorithm key generation, encryption, and decryption.
// Supported: RSA (OAEP-SHA256), ECC (ECDH+AES-GCM), X25519, Kyber (ML-KEM+AES-GCM).
package crypto

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/kyber/kyber1024"
	"github.com/cloudflare/circl/kem/kyber/kyber512"
	"github.com/cloudflare/circl/kem/kyber/kyber768"
	"golang.org/x/crypto/hkdf"
)

// ---- algorithm identifiers --------------------------------------------------

// KeyType identifies a supported cryptographic algorithm.
type KeyType byte

const (
	KeyRSA2048 KeyType = iota + 1
	KeyRSA3072
	KeyRSA4096
	KeyECCP256
	KeyECCP384
	KeyECCP521
	KeyX25519
	KeyKyber512
	KeyKyber768
	KeyKyber1024
)

func (k KeyType) String() string {
	switch k {
	case KeyRSA2048:
		return "RSA-2048"
	case KeyRSA3072:
		return "RSA-3072"
	case KeyRSA4096:
		return "RSA-4096"
	case KeyECCP256:
		return "ECC P-256"
	case KeyECCP384:
		return "ECC P-384"
	case KeyECCP521:
		return "ECC P-521"
	case KeyX25519:
		return "Curve25519"
	case KeyKyber512:
		return "Kyber-512"
	case KeyKyber768:
		return "Kyber-768"
	case KeyKyber1024:
		return "Kyber-1024"
	}
	return fmt.Sprintf("unknown(%d)", k)
}

// IsRSA returns true for RSA-family key types.
func (k KeyType) IsRSA() bool { return k >= KeyRSA2048 && k <= KeyRSA4096 }

// IsKyber returns true for Kyber-family key types.
func (k KeyType) IsKyber() bool { return k >= KeyKyber512 && k <= KeyKyber1024 }

// KyberScheme returns the circl KEM scheme for this Kyber variant (nil otherwise).
func (k KeyType) KyberScheme() kem.Scheme {
	switch k {
	case KeyKyber512:
		return kyber512.Scheme()
	case KeyKyber768:
		return kyber768.Scheme()
	case KeyKyber1024:
		return kyber1024.Scheme()
	}
	return nil
}

// RSABits returns the modulus size in bits (0 for non-RSA).
func (k KeyType) RSABits() int {
	switch k {
	case KeyRSA2048:
		return 2048
	case KeyRSA3072:
		return 3072
	case KeyRSA4096:
		return 4096
	}
	return 0
}

// ECDHCurve returns the ECDH curve for ECC/X25519 types (nil for RSA).
func (k KeyType) ECDHCurve() ecdh.Curve {
	switch k {
	case KeyECCP256:
		return ecdh.P256()
	case KeyECCP384:
		return ecdh.P384()
	case KeyECCP521:
		return ecdh.P521()
	case KeyX25519:
		return ecdh.X25519()
	}
	return nil
}

// detectKeyType inspects a parsed public/private key and returns its type.
func detectKeyType(key any) KeyType {
	switch k := key.(type) {
	case *rsa.PublicKey:
		switch k.Size() {
		case 256:
			return KeyRSA2048
		case 384:
			return KeyRSA3072
		case 512:
			return KeyRSA4096
		}
	case *rsa.PrivateKey:
		return detectKeyType(&k.PublicKey)
	case *ecdsa.PublicKey:
		// Convert to ECDH to identify the curve.
		if e, err := k.ECDH(); err == nil {
			return detectKeyType(e)
		}
	case *ecdsa.PrivateKey:
		if e, err := k.ECDH(); err == nil {
			return detectKeyType(e)
		}
	case *ecdh.PublicKey:
		switch k.Curve() {
		case ecdh.P256():
			return KeyECCP256
		case ecdh.P384():
			return KeyECCP384
		case ecdh.P521():
			return KeyECCP521
		case ecdh.X25519():
			return KeyX25519
		}
	case *ecdh.PrivateKey:
		return detectKeyType(k.PublicKey())
	case kem.PublicKey:
		switch k.Scheme().Name() {
		case "Kyber512":
			return KeyKyber512
		case "Kyber768":
			return KeyKyber768
		case "Kyber1024":
			return KeyKyber1024
		}
	case kem.PrivateKey:
		return detectKeyType(k.Public())
	}
	return 0
}

// encodeKeyPairPEM marshals a key pair to PEM blocks (supports RSA/ECDH/Kyber).
func encodeKeyPairPEM(algo KeyType, privKey, pubKey any, passphrase []byte) (privPEM, pubPEM []byte) {
	if algo.IsKyber() {
		return kyberToPEM(privKey, pubKey, passphrase)
	}
	// Standard PKCS#8 / PKIX for RSA and ECDH.
	privDER, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, nil
	}
	var privBlock *pem.Block
	if len(passphrase) > 0 {
		privBlock, _ = x509.EncryptPEMBlock(rand.Reader, "ENCRYPTED PRIVATE KEY", privDER, passphrase, x509.PEMCipherAES256)
	} else {
		privBlock = &pem.Block{Type: "PRIVATE KEY", Bytes: privDER}
	}
	pubDER, _ := x509.MarshalPKIXPublicKey(pubKey)
	pubBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}
	return pem.EncodeToMemory(privBlock), pem.EncodeToMemory(pubBlock)
}

// parsePublicKeyPEMGeneric handles both standard PKIX and Kyber PEM blocks.
func parsePublicKeyPEMGeneric(pemData []byte) (any, KeyType, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, 0, fmt.Errorf("no PEM block found")
	}
	switch block.Type {
	case "KYBER PUBLIC KEY":
		return parseKyberPubPEM(block.Bytes)
	default:
		k, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, 0, fmt.Errorf("parse public key: %w", err)
		}
		kt := detectKeyType(k)
		if kt == 0 {
			return nil, 0, fmt.Errorf("unsupported public key type: %T", k)
		}
		return k, kt, nil
	}
}

// parsePrivateKeyPEMGeneric handles both standard PKCS#8 and Kyber PEM blocks.
func parsePrivateKeyPEMGeneric(pemData, passphrase []byte) (any, KeyType, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, 0, fmt.Errorf("no PEM block found")
	}
	switch block.Type {
	case "KYBER PRIVATE KEY":
		return parseKyberPrivPEM(block.Bytes, passphrase)
	case "ENCRYPTED KYBER PRIVATE KEY":
		if len(passphrase) == 0 {
			return nil, 0, fmt.Errorf("kyber private key requires a passphrase")
		}
		der, err := x509.DecryptPEMBlock(block, passphrase)
		if err != nil {
			return nil, 0, fmt.Errorf("decrypt kyber key: %w", err)
		}
		return parseKyberPrivPEM(der, nil)
	default:
		var der []byte
		var err error
		if x509.IsEncryptedPEMBlock(block) {
			der, err = x509.DecryptPEMBlock(block, passphrase)
			if err != nil {
				return nil, 0, fmt.Errorf("decrypt private key: %w", err)
			}
		} else {
			der = block.Bytes
		}
		k, err := x509.ParsePKCS8PrivateKey(der)
		if err != nil {
			return nil, 0, fmt.Errorf("parse private key: %w", err)
		}
		kt := detectKeyType(k)
		if kt == 0 {
			return nil, 0, fmt.Errorf("unsupported private key type: %T", k)
		}
		return k, kt, nil
	}
}

// ---- Kyber PEM & encryption -------------------------------------------------

func kyberToPEM(privKey, pubKey any, passphrase []byte) (privPEM, pubPEM []byte) {
	pubBytes, _ := pubKey.(encoding.BinaryMarshaler).MarshalBinary()
	privBytes, _ := privKey.(encoding.BinaryMarshaler).MarshalBinary()

	pubBlock := &pem.Block{Type: "KYBER PUBLIC KEY", Bytes: pubBytes}
	pubPEM = pem.EncodeToMemory(pubBlock)

	if len(passphrase) > 0 {
		privBlock, _ := x509.EncryptPEMBlock(rand.Reader, "ENCRYPTED KYBER PRIVATE KEY", privBytes, passphrase, x509.PEMCipherAES256)
		privPEM = pem.EncodeToMemory(privBlock)
	} else {
		privBlock := &pem.Block{Type: "KYBER PRIVATE KEY", Bytes: privBytes}
		privPEM = pem.EncodeToMemory(privBlock)
	}
	return
}

func parseKyberPubPEM(der []byte) (kem.PublicKey, KeyType, error) {
	for _, a := range AllAlgos() {
		if !a.IsKyber() {
			continue
		}
		scheme := a.KyberScheme()
		if len(der) != scheme.PublicKeySize() {
			continue
		}
		pub, err := scheme.UnmarshalBinaryPublicKey(der)
		if err != nil {
			continue
		}
		return pub, a, nil
	}
	return nil, 0, fmt.Errorf("unsupported kyber public key")
}

func parseKyberPrivPEM(der, _ []byte) (kem.PrivateKey, KeyType, error) {
	for _, a := range AllAlgos() {
		if !a.IsKyber() {
			continue
		}
		scheme := a.KyberScheme()
		if len(der) != scheme.PrivateKeySize() {
			continue
		}
		priv, err := scheme.UnmarshalBinaryPrivateKey(der)
		if err != nil {
			continue
		}
		return priv, a, nil
	}
	return nil, 0, fmt.Errorf("unsupported kyber private key")
}

// pubKeySize returns the byte length of the raw public key for ECDH curves.
func pubKeySize(kt KeyType) int {
	switch kt {
	case KeyECCP256:
		return 65
	case KeyECCP384:
		return 97
	case KeyECCP521:
		return 133
	case KeyX25519:
		return 32
	}
	return 0
}

// ---- key generation ---------------------------------------------------------

// GenerateKeyPair creates a new key pair for the given algorithm.
func GenerateKeyPair(algo KeyType, passphrase []byte, privPath, pubPath string) error {
	var privKey, pubKey any

	switch {
	case algo.IsRSA():
		priv, e := rsa.GenerateKey(rand.Reader, algo.RSABits())
		if e != nil {
			return fmt.Errorf("RSA keygen: %w", e)
		}
		privKey, pubKey = priv, &priv.PublicKey
	case algo.ECDHCurve() != nil:
		priv, e := algo.ECDHCurve().GenerateKey(rand.Reader)
		if e != nil {
			return fmt.Errorf("ECDH keygen: %w", e)
		}
		privKey, pubKey = priv, priv.PublicKey()
	case algo.IsKyber():
		scheme := algo.KyberScheme()
		pub, priv, e := scheme.GenerateKeyPair()
		if e != nil {
			return fmt.Errorf("Kyber keygen: %w", e)
		}
		privKey, pubKey = priv, pub
	default:
		return fmt.Errorf("unsupported algorithm: %s", algo)
	}

	privPEM, pubPEM := encodeKeyPairPEM(algo, privKey, pubKey, passphrase)
	if err := writeFile(privPath, privPEM, 0o600); err != nil {
		return err
	}
	return writeFile(pubPath, pubPEM, 0o644)
}

// ---- key loading ------------------------------------------------------------

// LoadPublicKey parses a PEM-encoded public key from the given file.
// Returns the key, its algorithm type, and any error.
func LoadPublicKey(path string) (key any, algo KeyType, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("read public key: %w", err)
	}
	return parsePublicKeyPEMGeneric(data)
}

// LoadPrivateKey parses a PEM-encoded private key (PKCS#8 or Kyber).
func LoadPrivateKey(path string, passphrase []byte) (key any, algo KeyType, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, fmt.Errorf("read private key: %w", err)
	}
	return parsePrivateKeyPEMGeneric(data, passphrase)
}

// ---- encryption -------------------------------------------------------------

// Encrypt encrypts plaintext with the given public key.
// For RSA: OAEP-SHA256. For ECC/X25519: ECDH+AES-GCM. For Kyber: KEM+AES-GCM.
func Encrypt(pub any, algo KeyType, plaintext []byte) ([]byte, error) {
	switch {
	case algo.IsRSA():
		rsaKey, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("expected *rsa.PublicKey, got %T", pub)
		}
		return rsaEncrypt(rsaKey, plaintext)
	case algo.ECDHCurve() != nil:
		ecdhKey, err := toECDHPublic(pub)
		if err != nil {
			return nil, err
		}
		return ecdhEncrypt(ecdhKey, algo, plaintext)
	case algo.IsKyber():
		kyberKey, ok := pub.(kem.PublicKey)
		if !ok {
			return nil, fmt.Errorf("expected kem.PublicKey, got %T", pub)
		}
		return kyberEncrypt(kyberKey, algo, plaintext)
	}
	return nil, fmt.Errorf("unsupported algorithm: %s", algo)
}

// Decrypt decrypts ciphertext produced by Encrypt.
func Decrypt(priv any, algo KeyType, ciphertext []byte) ([]byte, error) {
	switch {
	case algo.IsRSA():
		rsaKey, ok := priv.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("expected *rsa.PrivateKey, got %T", priv)
		}
		return rsaDecrypt(rsaKey, ciphertext)
	case algo.ECDHCurve() != nil:
		ecdhKey, err := toECDHPrivate(priv)
		if err != nil {
			return nil, err
		}
		return ecdhDecrypt(ecdhKey, algo, ciphertext)
	case algo.IsKyber():
		kyberKey, ok := priv.(kem.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("expected kem.PrivateKey, got %T", priv)
		}
		return kyberDecrypt(kyberKey, algo, ciphertext)
	}
	return nil, fmt.Errorf("unsupported algorithm: %s", algo)
}

// toECDHPublic converts an ECDSA or ECDH public key to *ecdh.PublicKey.
func toECDHPublic(pub any) (*ecdh.PublicKey, error) {
	switch k := pub.(type) {
	case *ecdh.PublicKey:
		return k, nil
	case *ecdsa.PublicKey:
		return k.ECDH()
	}
	return nil, fmt.Errorf("cannot convert %T to ECDH public key", pub)
}

// toECDHPrivate converts an ECDSA or ECDH private key to *ecdh.PrivateKey.
func toECDHPrivate(priv any) (*ecdh.PrivateKey, error) {
	switch k := priv.(type) {
	case *ecdh.PrivateKey:
		return k, nil
	case *ecdsa.PrivateKey:
		return k.ECDH()
	}
	return nil, fmt.Errorf("cannot convert %T to ECDH private key", priv)
}

// ---- RSA OAEP ---------------------------------------------------------------

func rsaEncrypt(pub *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	hash := sha256.New()
	maxBlock := pub.Size() - 2*hash.Size() - 2
	var ct []byte
	ct = append(ct, byte(algoToByte(pub.Size()*8)))
	for len(plaintext) > 0 {
		chunk := plaintext
		if len(chunk) > maxBlock {
			chunk = plaintext[:maxBlock]
		}
		plaintext = plaintext[len(chunk):]
		block, err := rsa.EncryptOAEP(hash, rand.Reader, pub, chunk, nil)
		if err != nil {
			return nil, fmt.Errorf("RSA encrypt: %w", err)
		}
		ct = append(ct, block...)
	}
	return ct, nil
}

func rsaDecrypt(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	hash := sha256.New()
	segSize := priv.Size()
	if len(ciphertext) < 1+segSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	data := ciphertext[1:] // skip algo byte
	var pt []byte
	for len(data) > 0 {
		if len(data) < segSize {
			return nil, fmt.Errorf("ciphertext length not a multiple of block size")
		}
		block := data[:segSize]
		data = data[segSize:]
		plain, err := rsa.DecryptOAEP(hash, rand.Reader, priv, block, nil)
		if err != nil {
			return nil, fmt.Errorf("RSA decrypt: %w", err)
		}
		pt = append(pt, plain...)
	}
	return pt, nil
}

func algoToByte(bits int) byte {
	switch bits {
	case 2048:
		return byte(KeyRSA2048)
	case 3072:
		return byte(KeyRSA3072)
	case 4096:
		return byte(KeyRSA4096)
	}
	return 0
}

// ---- ECIES (ECDH + AES-256-GCM) ---------------------------------------------

func ecdhEncrypt(pub *ecdh.PublicKey, algo KeyType, plaintext []byte) ([]byte, error) {
	ephemeral, err := pub.Curve().GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("ephemeral key: %w", err)
	}
	shared, err := ephemeral.ECDH(pub)
	if err != nil {
		return nil, fmt.Errorf("ECDH: %w", err)
	}
	aesKey := deriveKey(shared, nil)
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ephemeralPubBytes := ephemeral.PublicKey().Bytes()
	ct := make([]byte, 0, 1+len(ephemeralPubBytes)+12+len(plaintext)+16)
	ct = append(ct, byte(algo))
	ct = append(ct, ephemeralPubBytes...)
	ct = append(ct, nonce...)
	ct = aesgcm.Seal(ct, nonce, plaintext, nil)
	return ct, nil
}

func ecdhDecrypt(priv *ecdh.PrivateKey, algo KeyType, ciphertext []byte) ([]byte, error) {
	pkSize := pubKeySize(algo)
	minLen := 1 + pkSize + 12 + 16
	if len(ciphertext) < minLen {
		return nil, fmt.Errorf("ciphertext too short for %s", algo)
	}
	ephemeralPubBytes := ciphertext[1 : 1+pkSize]
	nonce := ciphertext[1+pkSize : 1+pkSize+12]
	encrypted := ciphertext[1+pkSize+12:]

	curve := algo.ECDHCurve()
	ephemeralPub, err := curve.NewPublicKey(ephemeralPubBytes)
	if err != nil {
		return nil, fmt.Errorf("parse ephemeral public key: %w", err)
	}
	shared, err := priv.ECDH(ephemeralPub)
	if err != nil {
		return nil, fmt.Errorf("ECDH: %w", err)
	}
	aesKey := deriveKey(shared, nil)
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesgcm.Open(nil, nonce, encrypted, nil)
}

// ---- Kyber KEM + AES-256-GCM -------------------------------------------------

func kyberEncrypt(pub kem.PublicKey, algo KeyType, plaintext []byte) ([]byte, error) {
	scheme := algo.KyberScheme()
	ctKEM, ss, err := scheme.Encapsulate(pub)
	if err != nil {
		return nil, fmt.Errorf("Kyber encaps: %w", err)
	}
	aesKey := deriveKey(ss, nil)
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	block, _ := aes.NewCipher(aesKey)
	aesgcm, _ := cipher.NewGCM(block)
	out := make([]byte, 0, 1+len(ctKEM)+12+len(plaintext)+16)
	out = append(out, byte(algo))
	out = append(out, ctKEM...)
	out = append(out, nonce...)
	return aesgcm.Seal(out, nonce, plaintext, nil), nil
}

func kyberDecrypt(priv kem.PrivateKey, algo KeyType, ciphertext []byte) ([]byte, error) {
	scheme := algo.KyberScheme()
	kemLen := scheme.CiphertextSize()
	minLen := 1 + kemLen + 12 + 16
	if len(ciphertext) < minLen {
		return nil, fmt.Errorf("ciphertext too short for %s", algo)
	}
	ctKEM := ciphertext[1 : 1+kemLen]
	nonce := ciphertext[1+kemLen : 1+kemLen+12]
	encrypted := ciphertext[1+kemLen+12:]

	ss, err := scheme.Decapsulate(priv, ctKEM)
	if err != nil {
		return nil, fmt.Errorf("Kyber decaps: %w", err)
	}
	aesKey := deriveKey(ss, nil)
	block, _ := aes.NewCipher(aesKey)
	aesgcm, _ := cipher.NewGCM(block)
	return aesgcm.Open(nil, nonce, encrypted, nil)
}

func deriveKey(secret, salt []byte) []byte {
	r := hkdf.New(sha256.New, secret, salt, []byte("cryphocat-ecies-v1"))
	key := make([]byte, 32)
	_, _ = io.ReadFull(r, key)
	return key
}

// ---- helpers ----------------------------------------------------------------

func writeFile(path string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// AllAlgos returns all supported algorithms in display order.
func AllAlgos() []KeyType {
	return []KeyType{
		KeyRSA2048, KeyRSA3072, KeyRSA4096,
		KeyECCP256, KeyECCP384, KeyECCP521,
		KeyX25519,
		KeyKyber512, KeyKyber768, KeyKyber1024,
	}
}

// ---- encoding / compression -------------------------------------------------

const compressFlag byte = 0x01

// Pack encodes raw ciphertext (with algo prefix) as a compact string.
// If compress is true, zlib compresses the payload (after algo byte) before base64.
func Pack(raw []byte, compress bool) string {
	if len(raw) == 0 {
		return ""
	}
	if compress {
		raw[0] |= compressFlag // mark algo byte
		var buf bytes.Buffer
		w := zlib.NewWriter(&buf)
		w.Write(raw[1:]) // compress payload only
		w.Close()
		out := make([]byte, 1+buf.Len())
		out[0] = raw[0]
		copy(out[1:], buf.Bytes())
		return base64.RawStdEncoding.EncodeToString(out)
	}
	return base64.RawStdEncoding.EncodeToString(raw)
}

// Unpack decodes a string produced by Pack, auto-detecting compression.
func Unpack(encoded string) ([]byte, error) {
	raw, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}
	if len(raw) < 1 {
		return nil, fmt.Errorf("empty ciphertext")
	}
	if raw[0]&compressFlag != 0 {
		raw[0] &^= compressFlag // clear flag from algo byte
		r, err := zlib.NewReader(bytes.NewReader(raw[1:]))
		if err != nil {
			return nil, fmt.Errorf("zlib decompress: %w", err)
		}
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			r.Close()
			return nil, fmt.Errorf("zlib read: %w", err)
		}
		r.Close()
		out := make([]byte, 1+buf.Len())
		out[0] = raw[0]
		copy(out[1:], buf.Bytes())
		return out, nil
	}
	return raw, nil
}

// ---- in-memory key generation / parsing ------------------------------------

// GenerateKeyPairPEM generates a key pair and returns the PEM-encoded bytes.
func GenerateKeyPairPEM(algo KeyType, passphrase []byte) (privPEM, pubPEM []byte, err error) {
	var privKey, pubKey any
	switch {
	case algo.IsRSA():
		priv, e := rsa.GenerateKey(rand.Reader, algo.RSABits())
		if e != nil {
			return nil, nil, fmt.Errorf("RSA keygen: %w", e)
		}
		privKey, pubKey = priv, &priv.PublicKey
	case algo.ECDHCurve() != nil:
		priv, e := algo.ECDHCurve().GenerateKey(rand.Reader)
		if e != nil {
			return nil, nil, fmt.Errorf("ECDH keygen: %w", e)
		}
		privKey, pubKey = priv, priv.PublicKey()
	case algo.IsKyber():
		scheme := algo.KyberScheme()
		pub, priv, e := scheme.GenerateKeyPair()
		if e != nil {
			return nil, nil, fmt.Errorf("Kyber keygen: %w", e)
		}
		privKey, pubKey = priv, pub
	default:
		return nil, nil, fmt.Errorf("unsupported algorithm: %s", algo)
	}
	privPEM, pubPEM = encodeKeyPairPEM(algo, privKey, pubKey, passphrase)
	return privPEM, pubPEM, nil
}

// ParsePublicKeyPEM parses a PEM-encoded public key from memory.
func ParsePublicKeyPEM(pemData []byte) (key any, algo KeyType, err error) {
	return parsePublicKeyPEMGeneric(pemData)
}

// ParsePrivateKeyPEM parses a PEM-encoded private key from memory.
func ParsePrivateKeyPEM(pemData, passphrase []byte) (key any, algo KeyType, err error) {
	return parsePrivateKeyPEMGeneric(pemData, passphrase)
}

// imagePrefix is the MIME header prepended to base64-encoded images.
const ImagePrefix = "data:image/"

// IsImageBase64 checks whether text looks like a base64-encoded image.
func IsImageBase64(text string) bool {
	return len(text) > 22 && text[:11] == ImagePrefix
}
