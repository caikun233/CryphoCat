package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	cc "github.com/caikun233/cryphocat/crypto"
	"github.com/caikun233/cryphocat/filehash"
	"github.com/caikun233/cryphocat/keystore"
	"golang.org/x/term"
)

// ---- i18n strings -----------------------------------------------------------

type lang struct {
	keyLenPrompt   string
	keyLenBad      string
	keyLenNotNum   string
	promptPWSet    string
	promptPWConfirm string
	pwMismatch     string
	pwWarn         string
	keyGenOK       string
	keyPriv        string
	keyPub         string

	importPrompt string
	importOK     string
	importFail   string
	importNoFile string

	encNoKey     string
	encUseMine   string
	encCancelled string
	encNoKeyAny  string
	encLoadFail  string
	encPrompt    string
	encEmpty     string
	encFail      string
	encOK        string

	decNoKey       string
	promptPWDec    string
	decLoadFail    string
	decPrompt      string
	decBadHex      string
	decFail        string
	decOK          string

	hashPrompt   string
	hashNoFile   string
	hashMD5      string
	hashSHA1     string
	hashSHA256   string

	cmpPrompt    string
	cmpHashIn    string
	cmpNoFile    string
	cmpEmpty     string
	cmpMatch     string
	cmpMatchVal  string
	cmpFail      string

	menu1 string
	menu2 string
	menu3 string
	menu4 string
	menu5 string
	menu6 string
	menu7 string
	sep   string

	choicePrompt string
	choiceBad    string
	pressEnter   string
	goodbye      string
}

var en = lang{
	keyLenPrompt:       "Key length (2048/4096, default 4096): ",
	keyLenBad:          "Choose one of: 2048, 4096",
	keyLenNotNum:       "Please enter a number.",
	promptPWSet:      "Set passphrase (blank = no encryption, NOT recommended): ",
	promptPWConfirm:  "Confirm passphrase: ",
	pwMismatch: "Passphrases do not match.",
	pwWarn:     "[WARNING] Private key will be stored without a passphrase.",
	keyGenOK:           "RSA-%d key pair generated.",
	keyPriv:            "  Private: %s",
	keyPub:             "  Public : %s",

	importPrompt: "Path to friend's public key: ",
	importOK:     "Public key imported successfully.",
	importFail:   "Import failed: %v",
	importNoFile: "File not found. Check the path.",

	encNoKey:     "No friend key found.",
	encUseMine:   "Use your own public key? (y/n): ",
	encCancelled: "Cancelled. Import a friend's key or generate your key pair first.",
	encNoKeyAny:  "No public key available. Generate a key pair first.",
	encLoadFail:  "Cannot load public key: %v",
	encPrompt:    "Text to encrypt: ",
	encEmpty:     "Nothing to encrypt.",
	encFail:      "Encryption failed: %v",
	encOK:        "Encrypted (copied to clipboard): %s",

	decNoKey:      "Private key not found. Generate a key pair first.",
	promptPWDec: "Private key passphrase (blank if none): ",
	decLoadFail:   "Cannot load private key: %v",
	decPrompt:     "Ciphertext (hex): ",
	decBadHex:     "Invalid hex ciphertext.",
	decFail:       "Decryption failed: %v",
	decOK:         "Decrypted: %s",

	hashPrompt:  "File path: ",
	hashNoFile:  "File not found. Check the path.",
	hashMD5:     "MD5   : %s",
	hashSHA1:    "SHA1  : %s",
	hashSHA256:  "SHA256: %s",

	cmpPrompt:   "File path: ",
	cmpHashIn:   "Expected hash (MD5/SHA1/SHA256): ",
	cmpNoFile:   "File not found. Check the path.",
	cmpEmpty:    "No hash provided.",
	cmpMatch:    "%s MATCH: %s",
	cmpMatchVal: "Value : %s",
	cmpFail:     "HASH CHECK FAILED - file may have been tampered with.",

	menu1: "1. Generate RSA key pair",
	menu2: "2. Import friend's public key",
	menu3: "3. Encrypt text",
	menu4: "4. Decrypt text",
	menu5: "5. Calculate file hashes",
	menu6: "6. Verify file hash",
	menu7: "7. Exit",
	sep:   "----------------",

	choicePrompt: "Choice: ",
	choiceBad:    "Invalid choice.",
	pressEnter:   "Press Enter to continue...",
	goodbye:      "Goodbye.",
}

var zhcn = lang{
	keyLenPrompt:       "密钥长度（2048/4096，默认 4096）：",
	keyLenBad:          "请从 2048、4096 中选择一个。",
	keyLenNotNum:       "请输入数字。",
	promptPWSet:      "设置私钥口令（留空 = 不加密存储，不推荐）：",
	promptPWConfirm:  "确认口令：",
	pwMismatch: "两次输入的口令不一致。",
	pwWarn:     "[警告] 私钥将以明文形式存储，安全性较低。",
	keyGenOK:           "RSA-%d 密钥对已生成。",
	keyPriv:            "  私钥：%s",
	keyPub:             "  公钥：%s",

	importPrompt: "请输入对方公钥文件路径：",
	importOK:     "公钥导入成功。",
	importFail:   "导入失败：%v",
	importNoFile: "文件不存在，请检查路径。",

	encNoKey:     "未找到对方公钥。",
	encUseMine:   "是否使用自己的公钥加密？(y/n)：",
	encCancelled: "已取消。请先导入对方公钥或生成密钥对。",
	encNoKeyAny:  "未找到公钥，请先生成密钥对。",
	encLoadFail:  "无法加载公钥：%v",
	encPrompt:    "请输入要加密的文本：",
	encEmpty:     "输入为空。",
	encFail:      "加密失败：%v",
	encOK:        "加密结果（已复制到剪贴板）：%s",

	decNoKey:      "未找到私钥，请先生成密钥对。",
	promptPWDec: "私钥口令（无口令则留空）：",
	decLoadFail:   "无法加载私钥：%v",
	decPrompt:     "请输入十六进制密文：",
	decBadHex:     "十六进制密文无效。",
	decFail:       "解密失败：%v",
	decOK:         "解密结果：%s",

	hashPrompt: "请输入文件路径：",
	hashNoFile: "文件不存在，请检查路径。",
	hashMD5:    "MD5   ：%s",
	hashSHA1:   "SHA1  ：%s",
	hashSHA256: "SHA256：%s",

	cmpPrompt:   "请输入文件路径：",
	cmpHashIn:   "请输入哈希值（MD5/SHA1/SHA256）：",
	cmpNoFile:   "文件不存在，请检查路径。",
	cmpEmpty:    "未输入哈希值。",
	cmpMatch:    "%s 校验通过：%s",
	cmpMatchVal: "哈希值：%s",
	cmpFail:     "哈希校验失败 —— 文件可能已被篡改。",

	menu1: "1. 生成 RSA 密钥对",
	menu2: "2. 导入对方公钥",
	menu3: "3. 加密文本",
	menu4: "4. 解密文本",
	menu5: "5. 计算文件哈希",
	menu6: "6. 校验文件哈希",
	menu7: "7. 退出",
	sep:   "----------------",

	choicePrompt: "请输入选项：",
	choiceBad:    "选项无效，请重新输入。",
	pressEnter:   "按下回车继续...",
	goodbye:      "再见。",
}

// ---- helpers ----------------------------------------------------------------

var reader = bufio.NewReader(os.Stdin)

func readLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
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
	// Fallback (piped input / CI)
	line, _ := reader.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}

func setClipboard(text string) {
	// Try to write to clipboard; silently ignore errors (headless environments)
	_ = trySetClipboard(text)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func pressEnter(l *lang) {
	fmt.Print(l.pressEnter)
	_, _ = reader.ReadString('\n')
}

// ---- actions ----------------------------------------------------------------

func generateKeyPair(l *lang, ks keystore.Paths) {
	raw := readLine(l.keyLenPrompt)
	keySize := 4096
	if raw != "" {
		n := 0
		_, err := fmt.Sscanf(raw, "%d", &n)
		if err != nil || (n != 2048 && n != 4096) {
			if err != nil {
				fmt.Println(l.keyLenNotNum)
			} else {
				fmt.Println(l.keyLenBad)
			}
			return
		}
		keySize = n
	}

	pw := readPassword(l.promptPWSet)
	var passphrase []byte
	if pw != "" {
		pw2 := readPassword(l.promptPWConfirm)
		if pw != pw2 {
			fmt.Println(l.pwMismatch)
			return
		}
		passphrase = []byte(pw)
	} else {
		fmt.Println(l.pwWarn)
	}

	if err := cc.GenerateKeyPair(keySize, passphrase, ks.MyPriv, ks.MyPub); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf(l.keyGenOK+"\n", keySize)
	fmt.Printf(l.keyPriv+"\n", ks.MyPriv)
	fmt.Printf(l.keyPub+"\n", ks.MyPub)
}

func importPublicKey(l *lang, ks keystore.Paths) {
	path := readLine(l.importPrompt)
	if !fileExists(path) {
		fmt.Println(l.importNoFile)
		return
	}
	// Validate it's a real RSA public key.
	if _, err := cc.LoadPublicKey(path); err != nil {
		fmt.Printf(l.importFail+"\n", err)
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf(l.importFail+"\n", err)
		return
	}
	if err := os.WriteFile(ks.FrPub, data, 0o644); err != nil {
		fmt.Printf(l.importFail+"\n", err)
		return
	}
	fmt.Println(l.importOK)
}

func encryptData(l *lang, ks keystore.Paths) {
	var keyPath string
	switch {
	case fileExists(ks.FrPub):
		keyPath = ks.FrPub
	case fileExists(ks.MyPub):
		fmt.Println(l.encNoKey)
		ans := readLine(l.encUseMine)
		if strings.ToLower(ans) != "y" {
			fmt.Println(l.encCancelled)
			return
		}
		keyPath = ks.MyPub
	default:
		fmt.Println(l.encNoKeyAny)
		return
	}

	pub, err := cc.LoadPublicKey(keyPath)
	if err != nil {
		fmt.Printf(l.encLoadFail+"\n", err)
		return
	}

	text := readLine(l.encPrompt)
	if text == "" {
		fmt.Println(l.encEmpty)
		return
	}

	ct, err := cc.Encrypt(pub, []byte(text))
	if err != nil {
		fmt.Printf(l.encFail+"\n", err)
		return
	}
	hexCT := hex.EncodeToString(ct)
	setClipboard(hexCT)
	fmt.Printf(l.encOK+"\n", hexCT)
}

func decryptData(l *lang, ks keystore.Paths) {
	if !fileExists(ks.MyPriv) {
		fmt.Println(l.decNoKey)
		return
	}
	pw := readPassword(l.promptPWDec)
	var passphrase []byte
	if pw != "" {
		passphrase = []byte(pw)
	}
	priv, err := cc.LoadPrivateKey(ks.MyPriv, passphrase)
	if err != nil {
		fmt.Printf(l.decLoadFail+"\n", err)
		return
	}

	hexCT := readLine(l.decPrompt)
	ct, err := hex.DecodeString(strings.TrimSpace(hexCT))
	if err != nil {
		fmt.Println(l.decBadHex)
		return
	}
	pt, err := cc.Decrypt(priv, ct)
	if err != nil {
		fmt.Printf(l.decFail+"\n", err)
		return
	}
	fmt.Printf(l.decOK+"\n", string(pt))
}

func calculateHashes(l *lang) {
	path := readLine(l.hashPrompt)
	sums, err := filehash.Compute(path)
	if err != nil {
		fmt.Println(l.hashNoFile)
		return
	}
	fmt.Printf(l.hashMD5+"\n", sums.MD5)
	fmt.Printf(l.hashSHA1+"\n", sums.SHA1)
	fmt.Printf(l.hashSHA256+"\n", sums.SHA256)
}

func compareHashes(l *lang) {
	path := readLine(l.cmpPrompt)
	if !fileExists(path) {
		fmt.Println(l.cmpNoFile)
		return
	}
	expected := readLine(l.cmpHashIn)
	if expected == "" {
		fmt.Println(l.cmpEmpty)
		return
	}
	algo, err := filehash.Compare(path, expected)
	if err != nil {
		fmt.Println(l.cmpNoFile)
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
		fmt.Printf("\033[32m"+l.cmpMatch+"\033[0m\n", algo, path)
		fmt.Printf("\033[32m"+l.cmpMatchVal+"\033[0m\n", digest)
	} else {
		fmt.Printf("\033[31m%s\033[0m\n", l.cmpFail)
	}
}

// ---- main menu --------------------------------------------------------------

func runCLI(l *lang) {
	ks, err := keystore.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialise key store: %v\n", err)
		os.Exit(1)
	}
	for {
		fmt.Println(l.sep)
		fmt.Println(l.menu1)
		fmt.Println(l.menu2)
		fmt.Println(l.menu3)
		fmt.Println(l.menu4)
		fmt.Println(l.menu5)
		fmt.Println(l.menu6)
		fmt.Println(l.menu7)
		fmt.Println(l.sep)

		choice := readLine(l.choicePrompt)
		switch choice {
		case "1":
			generateKeyPair(l, ks)
			time.Sleep(time.Second)
		case "2":
			importPublicKey(l, ks)
		case "3":
			encryptData(l, ks)
			pressEnter(l)
		case "4":
			decryptData(l, ks)
			pressEnter(l)
		case "5":
			calculateHashes(l)
			pressEnter(l)
		case "6":
			compareHashes(l)
			pressEnter(l)
		case "7":
			fmt.Println(l.goodbye)
			return
		default:
			fmt.Println(l.choiceBad)
			pressEnter(l)
		}
	}
}

// ---- entry ------------------------------------------------------------------

func main() {
	guiFlag := flag.Bool("gui", false, "launch graphical user interface")
	langFlag := flag.String("lang", "en", "interface language: en or zhcn")
	flag.Parse()

	if *guiFlag {
		runGUI(*langFlag)
		return
	}

	if *langFlag == "zhcn" || *langFlag == "zh" {
		runCLI(&zhcn)
	} else {
		runCLI(&en)
	}
}

// trySetClipboard is defined in clipboard_*.go (platform build-tag shims).
