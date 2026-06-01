//go:build !cli

// Package main – CryphoCat GUI (Fyne).
// This file is compiled into the same main package.

package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	cc "github.com/caikun233/cryphocat/crypto"
	"github.com/caikun233/cryphocat/keystore"
)

// ---- GUI i18n ---------------------------------------------------------------

type guiLang struct {
	title string

	btnGenKey    string
	btnImportKey string
	btnKeyInfo   string

	labelInput  string
	labelOutput string

	btnEncrypt    string
	btnDecrypt    string
	btnClear      string
	btnCopyOutput string

	statusReady  string
	statusEncOK  string
	statusDecOK  string
	statusCleared string
	statusCopied string
	statusKeyImp string
	statusKeyGen string

	dlgPassTitle    string
	dlgPassMsg      string
	dlgPassConfirm  string
	dlgPassBadMsg   string
	dlgDecPassTitle string
	dlgDecPassMsg   string

	dlgNoFriendTitle string
	dlgNoFriendMsg   string

	errGenKeys string
	errImport  string
	errEncrypt string
	errDecrypt string
	errNoKey   string
	errBadHex  string
}

var guiEN = guiLang{
	title:          "CryphoCat",
	btnGenKey:      "Gen Keys",
	btnImportKey:   "Import Friend Key",
	btnKeyInfo:     "RSA / OAEP-SHA256",
	labelInput:     "INPUT:",
	labelOutput:    "OUTPUT:",
	btnEncrypt:     "Encrypt",
	btnDecrypt:     "Decrypt",
	btnClear:       "Clear",
	btnCopyOutput:  "Copy Output",
	statusReady:    "Ready.",
	statusEncOK:    "Encrypted - copied to clipboard.",
	statusDecOK:    "Decrypted.",
	statusCleared:  "Cleared.",
	statusCopied:   "Output copied to clipboard.",
	statusKeyImp:   "Friend's key imported.",
	statusKeyGen:   "Keys generated.",
	dlgPassTitle:   "Set Passphrase",
	dlgPassMsg:     "Passphrase for private key (blank = no encryption):",
	dlgPassConfirm: "Confirm passphrase:",
	dlgPassBadMsg:  "Passphrases do not match.",
	dlgDecPassTitle: "Passphrase",
	dlgDecPassMsg:   "Private key passphrase (blank if none):",
	dlgNoFriendTitle: "No friend key",
	dlgNoFriendMsg:   "Friend's public key not found.\nUse your own public key?",
	errGenKeys: "Key generation failed: %v",
	errImport:  "Import failed: %v",
	errEncrypt: "Encryption failed: %v",
	errDecrypt: "Decryption failed: %v",
	errNoKey:   "No public key found. Generate a key pair first.",
	errBadHex:  "Input is not valid hex ciphertext.",
}

var guiZhCN = guiLang{
	title:          "CryphoCat",
	btnGenKey:      "生成密钥对",
	btnImportKey:   "导入对方公钥",
	btnKeyInfo:     "RSA / OAEP-SHA256",
	labelInput:     "输入：",
	labelOutput:    "输出：",
	btnEncrypt:     "加密",
	btnDecrypt:     "解密",
	btnClear:       "清空",
	btnCopyOutput:  "复制输出",
	statusReady:    "就绪。",
	statusEncOK:    "加密完成 - 已自动复制到剪贴板。",
	statusDecOK:    "解密完成。",
	statusCleared:  "已清空。",
	statusCopied:   "输出已复制到剪贴板。",
	statusKeyImp:   "对方公钥已导入。",
	statusKeyGen:   "密钥对已生成。",
	dlgPassTitle:   "设置口令",
	dlgPassMsg:     "设置私钥保护口令（留空 = 不加密存储，不推荐）：",
	dlgPassConfirm: "确认口令：",
	dlgPassBadMsg:  "两次输入的口令不一致。",
	dlgDecPassTitle: "输入口令",
	dlgDecPassMsg:   "私钥口令（留空表示无口令）：",
	dlgNoFriendTitle: "未找到对方公钥",
	dlgNoFriendMsg:   "未找到对方公钥。\n是否使用自己的公钥加密？",
	errGenKeys: "密钥生成失败：%v",
	errImport:  "导入失败：%v",
	errEncrypt: "加密失败：%v",
	errDecrypt: "解密失败：%v",
	errNoKey:   "未找到公钥，请先生成密钥对。",
	errBadHex:  "输入的十六进制密文无效。",
}

// ---- GUI state ---------------------------------------------------------------

type guiState struct {
	l          *guiLang
	ks         keystore.Paths
	win        fyne.Window
	inputEntry *widget.Entry
	outputEntry *widget.Entry
	statusLbl  *widget.Label
	importBtn  *widget.Button
}

func (s *guiState) setStatus(msg string) { s.statusLbl.SetText(msg) }
func (s *guiState) setOutput(text string) {
	s.outputEntry.Enable()
	s.outputEntry.SetText(text)
	s.outputEntry.Disable()
}

func (s *guiState) askPassphrase(title, msg string, cb func([]byte, bool)) {
	pw := widget.NewPasswordEntry()
	dlg := dialog.NewForm(title, "OK", "Cancel",
		[]*widget.FormItem{widget.NewFormItem(msg, pw)},
		func(ok bool) {
			if !ok {
				cb(nil, false)
				return
			}
			if pw.Text == "" {
				cb(nil, true)
			} else {
				cb([]byte(pw.Text), true)
			}
		}, s.win)
	dlg.Show()
}

func (s *guiState) askNewPassphrase(cb func([]byte, bool)) {
	pw1 := widget.NewPasswordEntry()
	pw2 := widget.NewPasswordEntry()
	dlg := dialog.NewForm(s.l.dlgPassTitle, "OK", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem(s.l.dlgPassMsg, pw1),
			widget.NewFormItem(s.l.dlgPassConfirm, pw2),
		},
		func(ok bool) {
			if !ok {
				cb(nil, false)
				return
			}
			if pw1.Text != pw2.Text {
				dialog.ShowError(fmt.Errorf("%s", s.l.dlgPassBadMsg), s.win)
				cb(nil, false)
				return
			}
			if pw1.Text == "" {
				cb(nil, true)
			} else {
				cb([]byte(pw1.Text), true)
			}
		}, s.win)
	dlg.Show()
}

// ---- actions ----------------------------------------------------------------

func (s *guiState) doGenKeys() {
	s.askNewPassphrase(func(passphrase []byte, ok bool) {
		if !ok {
			return
		}
		if err := cc.GenerateKeyPair(2048, passphrase, s.ks.MyPriv, s.ks.MyPub); err != nil {
			dialog.ShowError(fmt.Errorf(s.l.errGenKeys, err), s.win)
			return
		}
		note := " [passphrase protected]"
		if len(passphrase) == 0 {
			note = " [WARNING: no passphrase]"
		}
		s.setOutput(fmt.Sprintf("[OK] RSA-2048 key pair generated%s.\n     Private : %s\n     Public  : %s",
			note, s.ks.MyPriv, s.ks.MyPub))
		s.setStatus(s.l.statusKeyGen)
	})
}

func (s *guiState) doImportKey() {
	dlg := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
		if err != nil || f == nil {
			return
		}
		defer f.Close()
		path := f.URI().Path()
		data, err := os.ReadFile(path)
		if err != nil {
			dialog.ShowError(fmt.Errorf(s.l.errImport, err), s.win)
			return
		}
		if _, err := cc.LoadPublicKey(path); err != nil {
			dialog.ShowError(fmt.Errorf(s.l.errImport, err), s.win)
			return
		}
		if err := os.WriteFile(s.ks.FrPub, data, 0o644); err != nil {
			dialog.ShowError(fmt.Errorf(s.l.errImport, err), s.win)
			return
		}
		s.setOutput("[OK] Friend's public key imported.")
		s.setStatus(s.l.statusKeyImp)
		s.importBtn.SetText("✓ " + s.l.btnImportKey)
	}, s.win)
	dlg.Show()
}

func (s *guiState) encryptWith(keyPath string) {
	pub, err := cc.LoadPublicKey(keyPath)
	if err != nil {
		dialog.ShowError(fmt.Errorf(s.l.errEncrypt, err), s.win)
		return
	}
	text := strings.TrimSpace(s.inputEntry.Text)
	if text == "" {
		return
	}
	ct, err := cc.Encrypt(pub, []byte(text))
	if err != nil {
		dialog.ShowError(fmt.Errorf(s.l.errEncrypt, err), s.win)
		return
	}
	hexCT := hex.EncodeToString(ct)
	s.setOutput(hexCT)
	setClipboard(hexCT)
	s.setStatus(s.l.statusEncOK)
}

func (s *guiState) doEncrypt() {
	switch {
	case fileExists(s.ks.FrPub):
		s.encryptWith(s.ks.FrPub)
	case fileExists(s.ks.MyPub):
		dialog.ShowConfirm(s.l.dlgNoFriendTitle, s.l.dlgNoFriendMsg, func(yes bool) {
			if yes {
				s.encryptWith(s.ks.MyPub)
			}
		}, s.win)
	default:
		dialog.ShowError(fmt.Errorf("%s", s.l.errNoKey), s.win)
	}
}

func (s *guiState) doDecrypt() {
	if !fileExists(s.ks.MyPriv) {
		dialog.ShowError(fmt.Errorf("Private key not found. Generate a key pair first."), s.win)
		return
	}
	s.askPassphrase(s.l.dlgDecPassTitle, s.l.dlgDecPassMsg, func(passphrase []byte, ok bool) {
		if !ok {
			return
		}
		priv, err := cc.LoadPrivateKey(s.ks.MyPriv, passphrase)
		if err != nil {
			dialog.ShowError(fmt.Errorf(s.l.errDecrypt, err), s.win)
			return
		}
		hexCT := strings.TrimSpace(s.inputEntry.Text)
		ct, err := hex.DecodeString(hexCT)
		if err != nil {
			dialog.ShowError(fmt.Errorf("%s", s.l.errBadHex), s.win)
			return
		}
		pt, err := cc.Decrypt(priv, ct)
		if err != nil {
			dialog.ShowError(fmt.Errorf(s.l.errDecrypt, err), s.win)
			return
		}
		s.setOutput(string(pt))
		s.setStatus(s.l.statusDecOK)
	})
}

func (s *guiState) doClear() {
	s.inputEntry.SetText("")
	s.setOutput("")
	s.setStatus(s.l.statusCleared)
}

func (s *guiState) doCopyOutput() {
	text := s.outputEntry.Text
	if text != "" {
		setClipboard(text)
		s.setStatus(s.l.statusCopied)
	}
}

// ---- build & run UI ---------------------------------------------------------

func runGUI(langFlag string) {
	l := &guiEN
	if langFlag == "zhcn" || langFlag == "zh" {
		l = &guiZhCN
	}

	ks, err := keystore.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "keystore init: %v\n", err)
		os.Exit(1)
	}

	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	win := a.NewWindow(l.title)
	win.Resize(fyne.NewSize(680, 500))

	st := &guiState{l: l, ks: ks, win: win}

	// -- top bar (key management) --
	genBtn := widget.NewButton(l.btnGenKey, st.doGenKeys)
	st.importBtn = widget.NewButton(l.btnImportKey, st.doImportKey)
	infoLbl := widget.NewLabel(l.btnKeyInfo)
	topBar := container.NewHBox(genBtn, st.importBtn, layout.NewSpacer(), infoLbl)

	// -- input --
	inputLbl := widget.NewLabel(l.labelInput)
	st.inputEntry = widget.NewMultiLineEntry()
	st.inputEntry.Wrapping = fyne.TextWrapWord
	inputScroll := container.NewVScroll(st.inputEntry)
	inputScroll.SetMinSize(fyne.NewSize(0, 160))
	inputArea := container.NewBorder(inputLbl, nil, nil, nil, inputScroll)

	// -- action buttons --
	encBtn := widget.NewButton(l.btnEncrypt, st.doEncrypt)
	decBtn := widget.NewButton(l.btnDecrypt, st.doDecrypt)
	clrBtn := widget.NewButton(l.btnClear, st.doClear)
	cpyBtn := widget.NewButton(l.btnCopyOutput, st.doCopyOutput)
	actBar := container.NewHBox(encBtn, clrBtn, decBtn, layout.NewSpacer(), cpyBtn)

	// -- output --
	outputLbl := widget.NewLabel(l.labelOutput)
	st.outputEntry = widget.NewMultiLineEntry()
	st.outputEntry.Wrapping = fyne.TextWrapWord
	st.outputEntry.Disable()
	outputScroll := container.NewVScroll(st.outputEntry)
	outputScroll.SetMinSize(fyne.NewSize(0, 120))
	outputArea := container.NewBorder(outputLbl, nil, nil, nil, outputScroll)

	// -- status bar --
	st.statusLbl = widget.NewLabel(l.statusReady)

	// -- main layout --
	upper := container.NewVBox(topBar, inputArea, actBar)
	content := container.NewBorder(upper, st.statusLbl, nil, nil, outputArea)

	win.SetContent(content)
	win.ShowAndRun()
}
