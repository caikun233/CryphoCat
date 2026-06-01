package filehash_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caikun233/cryphocat/filehash"
)

func TestComputeAndCompare(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.bin")

	data := []byte("CryphoCat hash test data 喵")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}

	sums, err := filehash.Compute(path)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	if sums.MD5 == "" || sums.SHA1 == "" || sums.SHA256 == "" {
		t.Fatal("one or more digests are empty")
	}

	// Each digest should match via Compare.
	for _, expected := range []string{sums.MD5, sums.SHA1, sums.SHA256} {
		algo, err := filehash.Compare(path, expected)
		if err != nil {
			t.Fatalf("Compare: %v", err)
		}
		if algo == "" {
			t.Fatalf("Compare returned no match for known digest %s", expected)
		}
	}

	// Wrong hash must not match.
	algo, err := filehash.Compare(path, "000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	if algo != "" {
		t.Fatalf("expected no match for wrong hash, got algo=%s", algo)
	}
}

func TestComputeMissingFile(t *testing.T) {
	_, err := filehash.Compute("/tmp/does_not_exist_cryphocat_test")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
