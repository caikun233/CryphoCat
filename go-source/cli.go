// CLI mode – menu-driven terminal interface.
package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	cc "github.com/caikun233/cryphocat/crypto"
	"github.com/caikun233/cryphocat/filehash"
	"github.com/caikun233/cryphocat/internal/clipboard"
	"github.com/caikun233/cryphocat/internal/locale"
	"github.com/caikun233/cryphocat/keystore"
	"golang.org/x/term"
)

var stdinReader = bufio.NewReader(os.Stdin)
var lastLoadedText string

func readLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := stdinReader.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}

func readPassword(prompt string) string {
	fmt.Print(prompt)
	if fd := int(syscall.Stdin); term.IsTerminal(fd) {
		pw, err := term.ReadPassword(fd)
		fmt.Println()
		if err == nil {
			return string(pw)
		}
	}
	// fallback (piped / non-tty)
	line, _ := stdinReader.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}

func pressEnter(l *locale.Locale) {
	fmt.Print(l.PressEnter)
	_, _ = stdinReader.ReadString('\n')
}

// ---- CLI actions ------------------------------------------------------------

func cliGenKeys(l *locale.Locale, ks keystore.Paths) {
	// Algorithm selection
	fmt.Println(l.MenuSep)
	algos := cc.AllAlgos()
	for i, a := range algos {
		fmt.Printf("  %d. %s\n", i+1, a.String())
	}
	fmt.Println(l.MenuSep)
	raw := readLine(l.KeyLenPrompt)
	idx := 0
	if raw != "" {
		if _, err := fmt.Sscanf(raw, "%d", &idx); err != nil || idx < 1 || idx > len(algos) {
			fmt.Println(l.KeyLenBad)
			return
		}
	} else {
		idx = 3 // default RSA-4096
	}
	algo := algos[idx-1]

	pw := readPassword(l.PromptPWSet)
	var passphrase []byte
	if pw != "" {
		pw2 := readPassword(l.PromptPWConfirm)
		if pw != pw2 {
			fmt.Println(l.PwMismatch)
			return
		}
		passphrase = []byte(pw)
	} else {
		fmt.Println(l.PwWarn)
	}

	if ks.Disk {
		if err := cc.GenerateKeyPair(algo, passphrase, ks.MyPriv, ks.MyPub); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	} else {
		privPEM, pubPEM, err := cc.GenerateKeyPairPEM(algo, passphrase)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := ks.SaveMyKey(privPEM, pubPEM); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}
	fmt.Printf(l.KeyGenOK+"\n", algo.String())
	fmt.Printf(l.KeyPriv+"\n", ks.MyPriv)
	fmt.Printf(l.KeyPub+"\n", ks.MyPub)
}

func cliImportKey(l *locale.Locale, ks keystore.Paths) {
	path := readLine(l.ImportPrompt)
	if !fileExists(path) {
		fmt.Println(l.ImportNoFile)
		return
	}
	if _, _, err := cc.LoadPublicKey(path); err != nil {
		fmt.Printf(l.ImportFail+"\n", err)
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf(l.ImportFail+"\n", err)
		return
	}
	if err := ks.SaveFriendKey(data); err != nil {
		fmt.Printf(l.ImportFail+"\n", err)
		return
	}
	fmt.Println(l.ImportOK)
}

// getInputText returns text from pre-loaded file or prompts the user.
func getInputText(prompt, emptyMsg string) string {
	if lastLoadedText != "" {
		t := lastLoadedText
		lastLoadedText = ""
		fmt.Printf("Using loaded text (%d bytes).\n", len(t))
		return t
	}
	t := readLine(prompt)
	if t == "" {
		fmt.Println(emptyMsg)
	}
	return t
}

func cliEncrypt(l *locale.Locale, ks keystore.Paths) {
	switch {
	case ks.HasFrPub():
		pemBytes, _ := ks.ReadFrPub()
		pub, algo, err := cc.ParsePublicKeyPEM(pemBytes)
		if err != nil {
			fmt.Printf(l.EncLoadFail+"\n", err)
			return
		}
		text := getInputText(l.EncPrompt, l.EncEmpty)
		if text == "" {
			return
		}
		ct, err := cc.Encrypt(pub, algo, []byte(text))
		if err != nil {
			fmt.Printf(l.EncFail+"\n", err)
			return
		}
		packed := cc.Pack(ct, true)
		setClipboard(packed)
		fmt.Printf(l.EncOK+"\n", packed)
	case ks.HasMyPub():
		fmt.Println(l.EncNoKey)
		ans := readLine(l.EncUseMine)
		if strings.ToLower(ans) != "y" {
			fmt.Println(l.EncCancelled)
			return
		}
		pemBytes, _ := ks.ReadMyPub()
		pub, algo, err := cc.ParsePublicKeyPEM(pemBytes)
		if err != nil {
			fmt.Printf(l.EncLoadFail+"\n", err)
			return
		}
		text := getInputText(l.EncPrompt, l.EncEmpty)
		if text == "" {
			return
		}
		ct, err := cc.Encrypt(pub, algo, []byte(text))
		if err != nil {
			fmt.Printf(l.EncFail+"\n", err)
			return
		}
		packed := cc.Pack(ct, true)
		setClipboard(packed)
		fmt.Printf(l.EncOK+"\n", packed)
	default:
		fmt.Println(l.EncNoKeyAny)
	}
}

func cliDecrypt(l *locale.Locale, ks keystore.Paths) {
	if !ks.HasMyPriv() {
		fmt.Println(l.DecNoKey)
		return
	}
	pw := readPassword(l.PromptPWDec)
	var passphrase []byte
	if pw != "" {
		passphrase = []byte(pw)
	}
	pemBytes, _ := ks.ReadMyPriv()
	priv, algo, err := cc.ParsePrivateKeyPEM(pemBytes, passphrase)
	if err != nil {
		fmt.Printf(l.DecLoadFail+"\n", err)
		return
	}

	packed := getInputText(l.DecPrompt, "")
	raw, err := cc.Unpack(packed)
	if err != nil {
		// Try legacy hex
		if ct, e2 := hex.DecodeString(strings.TrimSpace(packed)); e2 == nil {
			raw = ct
		} else {
			fmt.Println(l.DecBadHex)
			return
		}
	}
	pt, err := cc.Decrypt(priv, algo, raw)
	if err != nil {
		fmt.Printf(l.DecFail+"\n", err)
		return
	}
	plain := string(pt)
	if cc.IsImageBase64(plain) {
		fmt.Println("[Image detected. Copy the output and paste into a browser to view.]")
	}
	fmt.Printf(l.DecOK+"\n", plain)
}

func cliCopyPub(l *locale.Locale, ks keystore.Paths) {
	data, err := ks.ReadMyPub()
	if err != nil || len(data) == 0 {
		fmt.Println(l.ErrNoKey)
		return
	}
	setClipboard(string(data))
	fmt.Println(l.DlgPubCopied)
}

func cliImportClip(l *locale.Locale, ks keystore.Paths) {
	text, err := clipboard.Get()
	if err != nil || text == "" {
		fmt.Println(l.DlgClipImportBad)
		return
	}
	if _, _, err := cc.ParsePublicKeyPEM([]byte(text)); err != nil {
		fmt.Println(l.DlgClipImportBad)
		return
	}
	if err := ks.SaveFriendKey([]byte(text)); err != nil {
		fmt.Printf(l.ImportFail+"\n", err)
		return
	}
	fmt.Println(l.ImportOK)
}

func cliLoadFile(l *locale.Locale) string {
	path := readLine(l.ImportPrompt)
	if !fileExists(path) {
		fmt.Println(l.ImportNoFile)
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	return string(data)
}

func cliHash(l *locale.Locale) {
	path := readLine(l.HashPrompt)
	sums, err := filehash.Compute(path)
	if err != nil {
		fmt.Println(l.HashNoFile)
		return
	}
	fmt.Printf(l.HashMD5+"\n", sums.MD5)
	fmt.Printf(l.HashSHA1+"\n", sums.SHA1)
	fmt.Printf(l.HashSHA256+"\n", sums.SHA256)
}

func cliVerify(l *locale.Locale) {
	path := readLine(l.CmpPrompt)
	if !fileExists(path) {
		fmt.Println(l.CmpNoFile)
		return
	}
	expected := readLine(l.CmpHashIn)
	if expected == "" {
		fmt.Println(l.CmpEmpty)
		return
	}
	algo, err := filehash.Compare(path, expected)
	if err != nil {
		fmt.Println(l.CmpNoFile)
		return
	}
	if algo != "" {
		sums, _ := filehash.Compute(path)
		var digest string
		switch algo {
		case "MD5":
			digest = sums.MD5
		case "SHA1":
			digest = sums.SHA1
		case "SHA256":
			digest = sums.SHA256
		}
		fmt.Printf("\033[32m"+l.CmpMatch+"\033[0m\n", algo, path)
		fmt.Printf("\033[32m"+l.CmpMatchVal+"\033[0m\n", digest)
	} else {
		fmt.Printf("\033[31m%s\033[0m\n", l.CmpFail)
	}
}

// ---- CLI menu loop ----------------------------------------------------------

func runCLI(l *locale.Locale) {
	ks := keystore.InitMem() // default: memory only
	for {
		fmt.Println(l.MenuSep)
		fmt.Println("  " + l.BtnGenKey)
		fmt.Println("  " + l.BtnImportKey)
		fmt.Println("  " + l.BtnImportClip)
		fmt.Println("  " + l.BtnCopyPub)
		fmt.Println("  " + l.BtnLoadFile)
		fmt.Println("  " + l.BtnEncrypt)
		fmt.Println("  " + l.BtnDecrypt)
		fmt.Println("  " + l.MenuHash)
		fmt.Println("  " + l.MenuVerify)
		fmt.Println("  " + l.MenuExit)
		fmt.Println(l.MenuSep)

		choice := readLine(l.ChoicePrompt)
		switch choice {
		case "1":
			cliGenKeys(l, ks)
			time.Sleep(time.Second)
		case "2":
			cliImportKey(l, ks)
		case "3":
			cliImportClip(l, ks)
		case "4":
			cliCopyPub(l, ks)
		case "5":
			text := cliLoadFile(l)
			if text != "" {
				fmt.Printf("Loaded %d bytes.\n", len(text))
				// store for later encrypt/decrypt use via a temporary variable
				lastLoadedText = text
			}
		case "6":
			cliEncrypt(l, ks)
			pressEnter(l)
		case "7":
			cliDecrypt(l, ks)
			pressEnter(l)
		case "8":
			cliHash(l)
			pressEnter(l)
		case "9":
			cliVerify(l)
			pressEnter(l)
		case "10":
			fmt.Println(l.Goodbye)
			return
		default:
			fmt.Println(l.ChoiceBad)
			pressEnter(l)
		}
	}
}
