// Package filehash computes and compares MD5, SHA1, and SHA256 digests of files.
package filehash

import (
	"crypto/md5"  //nolint:gosec // MD5 used for hash display only, not security
	"crypto/sha1" //nolint:gosec // SHA1 used for hash display only, not security
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// Sums holds the hex-encoded digests of a file.
type Sums struct {
	MD5    string
	SHA1   string
	SHA256 string
}

// Compute reads the file at path and returns its MD5, SHA1, and SHA256 digests.
func Compute(path string) (Sums, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Sums{}, fmt.Errorf("read file: %w", err)
	}
	md5sum := md5.Sum(data)   //nolint:gosec
	sha1sum := sha1.Sum(data) //nolint:gosec
	sha256sum := sha256.Sum256(data)
	return Sums{
		MD5:    hex.EncodeToString(md5sum[:]),
		SHA1:   hex.EncodeToString(sha1sum[:]),
		SHA256: hex.EncodeToString(sha256sum[:]),
	}, nil
}

// Compare returns the algorithm name that matches expected, or an empty string
// if none match.
func Compare(path, expected string) (string, error) {
	sums, err := Compute(path)
	if err != nil {
		return "", err
	}
	exp := strings.ToLower(strings.TrimSpace(expected))
	switch exp {
	case sums.MD5:
		return "MD5", nil
	case sums.SHA1:
		return "SHA1", nil
	case sums.SHA256:
		return "SHA256", nil
	}
	return "", nil
}
