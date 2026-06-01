# CryphoCat
Encrypt your chatting with me, meow~

Currently only Simplified Chinese and English are supported.

![](https://img.shields.io/badge/go-1.25-blue)
![GitHub](https://img.shields.io/github/license/caikun233/CryphoCat)

**中文README请参照**[README.zhcn.md](https://github.com/caikun233/CryphoCat/blob/main/README.zhcn.md)

* Support asymmetric encryption of arbitrary length text using RSA (1024/2048/4096-bit, selectable in CLI; 2048-bit in GUI).
* The key directory generation rules in the "two-person chat" scenario have been built in. You can Pull Request to add more.
* All offline processing, and open source, do not upload any data, decoupled from the chat software.
* As long as your chat software can guarantee that the information sent and received is the same, then the content of your communication will never be revealed.
* **Not Support Audio/Video At All, Unless there is a way to convert audio/video to text**.
* You can base64 encode the picture and send it out, and I may add the base64 encoding function to the software later.

## INSTALLATION & USAGE

### Linux / Windows / macOS CLI 🔨

1. Download the CLI binary for your platform from **Actions** artifacts, or build from source (see below).

2. Run it:

   ```
   ./CryphoCat_CLI_linux_amd64
   ```

   Use `--lang zhcn` for Chinese interface:

   ```
   ./CryphoCat_CLI_linux_amd64 --lang zhcn
   ```

3. The program will automatically create a folder named `RSAkeys` in the current directory, with `my/` and `friend/` sub-folders to store your key pair and your chat partner's public key.

### Linux / Windows GUI 🔨

1. Download the GUI binary from **Actions** artifacts. **x64 only**.
2. Double-click to run (Windows), or `./CryphoCat_GUI_linux_amd64` on Linux.
3. Use `--lang zhcn` flag for Chinese interface:

   ```
   ./CryphoCat_GUI_linux_amd64 --lang zhcn
   ```

### Build from Source 🛠

Requires **Go 1.25+**. On Linux also install:
```
sudo apt-get install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev \
  libxi-dev libgl1-mesa-dev libxxf86vm-dev pkg-config
```

```bash
cd go-source

# GUI build (default):
go build -o CryphoCat .

# CLI-only build (no GUI, no OpenGL dependency):
go build -tags cli -o CryphoCat_CLI .

# Run tests:
go test ./...
```

## Development Plan

- [x] Add English support for releases and source code.
- [x] Finish GUI version development.
- [x] Make RSA key length optional (CLI supports 1024/2048/4096).
- [x] Let the encrypted text be copied automatically.
- [x] Calculate and compare file hashes (MD5 / SHA1 / SHA256).
- [x] Rewrite in Go with Fyne GUI (single self-contained binary, no Python required).
- [ ] Make GUI more beautiful.
- [ ] Try to add images base64 encode.
