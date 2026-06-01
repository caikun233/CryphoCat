// Package keystore manages the RSAkeys directory structure on disk.
package keystore

import (
	"os"
	"path/filepath"
)

const (
	Base   = "RSAkeys"
	MyDir  = "my"
	FrDir  = "friend"
	PriPEM = "private.pem"
	PubPEM = "public.pem"
)

// Paths holds the resolved key file paths relative to the working directory.
type Paths struct {
	MyDir  string
	FrDir  string
	MyPriv string
	MyPub  string
	FrPub  string
}

// Init creates the RSAkeys directory tree and returns the resolved paths.
func Init() (Paths, error) {
	myDir := filepath.Join(Base, MyDir)
	frDir := filepath.Join(Base, FrDir)
	for _, d := range []string{myDir, frDir} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			return Paths{}, err
		}
	}
	return Paths{
		MyDir:  myDir,
		FrDir:  frDir,
		MyPriv: filepath.Join(myDir, PriPEM),
		MyPub:  filepath.Join(myDir, PubPEM),
		FrPub:  filepath.Join(frDir, PubPEM),
	}, nil
}
