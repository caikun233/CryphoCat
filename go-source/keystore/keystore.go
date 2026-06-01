// Package keystore manages key storage on disk or in memory.
package keystore

import (
	"os"
	"path/filepath"
)

const (
	Base   = ".cryphocat_keys"
	MyDir  = "my"
	FrDir  = "friend"
	PriPEM = "private.pem"
	PubPEM = "public.pem"
)

// Paths holds the resolved key paths (disk mode) or in-memory key bytes.
type Paths struct {
	Disk bool // true when using filesystem storage

	MyDir  string
	FrDir  string
	MyPriv string
	MyPub  string
	FrPub  string

	// In-memory keys (set when Disk == false).
	MemPriv  []byte
	MemPub   []byte
	MemFrPub []byte
}

// Init creates the key directory tree on disk and returns resolved paths.
func Init() (Paths, error) {
	myDir := filepath.Join(Base, MyDir)
	frDir := filepath.Join(Base, FrDir)
	for _, d := range []string{myDir, frDir} {
		if err := os.MkdirAll(d, 0o700); err != nil {
			return Paths{}, err
		}
	}
	return Paths{
		Disk:   true,
		MyDir:  myDir,
		FrDir:  frDir,
		MyPriv: filepath.Join(myDir, PriPEM),
		MyPub:  filepath.Join(myDir, PubPEM),
		FrPub:  filepath.Join(frDir, PubPEM),
	}, nil
}

// InitMem returns an in-memory-only Paths (no files written).
func InitMem() Paths {
	return Paths{Disk: false}
}

// SaveMyKey stores the private/public key pair (disk or memory).
func (p *Paths) SaveMyKey(privPEM, pubPEM []byte) error {
	if p.Disk {
		if err := os.WriteFile(p.MyPriv, privPEM, 0o600); err != nil {
			return err
		}
		return os.WriteFile(p.MyPub, pubPEM, 0o644)
	}
	p.MemPriv = privPEM
	p.MemPub = pubPEM
	return nil
}

// SaveFriendKey stores the friend's public key (disk or memory).
func (p *Paths) SaveFriendKey(pubPEM []byte) error {
	if p.Disk {
		return os.WriteFile(p.FrPub, pubPEM, 0o644)
	}
	p.MemFrPub = pubPEM
	return nil
}

// ReadMyPriv returns the private key bytes.
func (p *Paths) ReadMyPriv() ([]byte, error) {
	if p.Disk {
		return os.ReadFile(p.MyPriv)
	}
	return p.MemPriv, nil
}

// ReadMyPub returns the public key bytes.
func (p *Paths) ReadMyPub() ([]byte, error) {
	if p.Disk {
		return os.ReadFile(p.MyPub)
	}
	return p.MemPub, nil
}

// ReadFrPub returns the friend's public key bytes.
func (p *Paths) ReadFrPub() ([]byte, error) {
	if p.Disk {
		return os.ReadFile(p.FrPub)
	}
	return p.MemFrPub, nil
}

// HasMyPriv reports whether a private key exists.
func (p *Paths) HasMyPriv() bool {
	if p.Disk {
		_, err := os.Stat(p.MyPriv)
		return err == nil
	}
	return len(p.MemPriv) > 0
}

// HasMyPub reports whether a public key exists.
func (p *Paths) HasMyPub() bool {
	if p.Disk {
		_, err := os.Stat(p.MyPub)
		return err == nil
	}
	return len(p.MemPub) > 0
}

// HasFrPub reports whether a friend's public key exists.
func (p *Paths) HasFrPub() bool {
	if p.Disk {
		_, err := os.Stat(p.FrPub)
		return err == nil
	}
	return len(p.MemFrPub) > 0
}
