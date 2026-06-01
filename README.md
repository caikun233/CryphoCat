# CryphoCat

Offline asymmetric encryption tool for private messaging. Supports RSA, ECC, Curve25519, and Kyber post-quantum algorithms. Single executable — double-click to launch native Windows GUI, or run in a terminal for interactive CLI mode.

![](https://img.shields.io/badge/go-1.22-blue)
![GitHub](https://img.shields.io/github/license/caikun233/CryphoCat)

[中文说明](README.zhcn.md)

## Highlights

- **10 cryptographic algorithms**: RSA-2048/3072/4096, ECC P-256/P-384/P-521, Curve25519, Kyber-512/768/1024
- **Compact output**: Base64 + zlib compression reduces ciphertext length by ~60% compared to hex encoding
- **Memory-first**: keys default to in-memory storage, disk optional
- **Image support**: load an image and it will be auto-encoded to base64 before encryption; on decryption, images are detected and displayed in a preview window
- **Clipboard integration**: one-click copy of your public key; import a friend's key directly from the clipboard
- **File hashing**: built-in MD5 / SHA1 / SHA256 computation and verification
- **Native Windows GUI**, ~7.7 MB binary
- **10 interface languages**: EN, 简体中文, 繁體中文(臺灣/香港), 日本語, 한국어, Русский, Français, Deutsch, Español
- Fully offline, no telemetry, no data leaves your machine

## Quick Start

Download from [Releases](https://github.com/caikun233/CryphoCat/releases).

```
cryphocat.exe                    # auto-detect: terminal -> CLI, double-click -> GUI
cryphocat.exe --gui              # force GUI
cryphocat.exe --cli --lang zhcn  # force CLI, Chinese interface
```

Keys are stored in `.cryphocat_keys/` by default, or in memory only when the checkbox is ticked.

## Build

Requires Go 1.22+. Install `rsrc` once: `go install github.com/akavel/rsrc@latest`.

```
cd go-source
rsrc -manifest cryphocat.manifest -ico favicon.ico -o rsrc.syso
go build -ldflags="-s -w" -o cryphocat.exe .
go test ./...                     # run tests
```

Non-Windows builds (CLI only): skip `rsrc`, just run `go build`.

## Development Plan

- [x] RSA-OAEP, ECC (ECDH+AES-GCM), Curve25519, Kyber ML-KEM encryption
- [x] File hash computation and verification (MD5/SHA1/SHA256)
- [x] Auto-copy ciphertext to clipboard; import/export keys via clipboard
- [x] 10-language i18n with automatic OS language detection on Windows
- [x] Native Windows GUI (~7.7 MB, down from 42 MB)
- [x] Auto-detect CLI vs GUI mode at startup
- [x] Tag-triggered GitHub Release with auto-generated release notes
- [x] Image loading for base64-encoding before encryption, with preview on decryption
- [x] Base64+zlib compact output encoding
- [x] In-memory key storage
