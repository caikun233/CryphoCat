//go:build windows

package main

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"

	cc "github.com/caikun233/cryphocat/crypto"
	"github.com/caikun233/cryphocat/filehash"
	"github.com/caikun233/cryphocat/internal/clipboard"
	"github.com/caikun233/cryphocat/internal/locale"
	"github.com/caikun233/cryphocat/keystore"
	"golang.org/x/sys/windows"
)

//go:embed logo.png
var logoPNG []byte

var (
	comctl32               = syscall.NewLazyDLL("comctl32.dll")
	initCommonControlsExFn = comctl32.NewProc("InitCommonControlsEx")
)

type iccex struct{ Size, ICC uint32 }

const iccStandard = 0x00004000

func autoGUI() bool {
	switch parentProcessName() {
	case "explorer.exe":
		return true
	case "cmd.exe", "powershell.exe", "pwsh.exe", "wt.exe", "conhost.exe", "code.exe", "devenv.exe":
		return false
	default:
		return win.GetConsoleWindow() == 0
	}
}

func parentProcessName() string {
	snap, _ := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if snap == 0 {
		return ""
	}
	defer windows.CloseHandle(snap)
	pid := windows.GetCurrentProcessId()
	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	for err := windows.Process32First(snap, &pe); err == nil; err = windows.Process32Next(snap, &pe) {
		if pe.ProcessID == pid {
			ppid := pe.ParentProcessID
			for err2 := windows.Process32First(snap, &pe); err2 == nil; err2 = windows.Process32Next(snap, &pe) {
				if pe.ProcessID == ppid {
					return strings.ToLower(windows.UTF16ToString(pe.ExeFile[:]))
				}
			}
			break
		}
	}
	return ""
}

// ---- app state --------------------------------------------------------------

type guiApp struct {
	loc           *locale.Locale
	ks            keystore.Paths
	mw            *walk.MainWindow
	algoCombo     *walk.ComboBox
	memCheck      *walk.CheckBox
	compressCheck *walk.CheckBox
	friendKeyEdit *walk.LineEdit
	inputTE       *walk.TextEdit
	outputTE      *walk.TextEdit
	statusLabel   *walk.Label
	keyType       cc.KeyType
	useMemory     bool
	useCompress   bool
}

func (app *guiApp) setStatus(s string) { app.statusLabel.SetText(s) }
func (app *guiApp) setOutput(s string) { app.outputTE.SetText(s) }
func (app *guiApp) getInput() string   { return strings.TrimSpace(app.inputTE.Text()) }

func (app *guiApp) hasMyKeys() bool {
	if app.useMemory {
		return app.ks.HasMyPriv() && app.ks.HasMyPub()
	}
	return fileExists(app.ks.MyPriv) && fileExists(app.ks.MyPub)
}

// ---- dialogs ----------------------------------------------------------------

func (app *guiApp) askPassphrase(title, msg string) (string, bool) {
	var dlg *walk.Dialog
	var pw *walk.LineEdit
	result := ""
	ok := false
	Dialog{
		AssignTo: &dlg, Title: title,
		MinSize: Size{Width: 420, Height: 0},
		Layout:  VBox{Margins: Margins{Left: 12, Top: 12, Right: 12, Bottom: 8}},
		Children: []Widget{
			Label{Text: msg},
			LineEdit{AssignTo: &pw, PasswordMode: true},
			Composite{
				Layout: HBox{MarginsZero: true, Spacing: 8},
				Children: []Widget{
					HSpacer{},
					PushButton{Text: "Cancel", OnClicked: func() { dlg.Cancel() }},
					PushButton{Text: "OK", OnClicked: func() { result = pw.Text(); ok = true; dlg.Accept() }},
				},
			},
		},
	}.Create(app.mw)
	dlg.Run()
	return result, ok
}

func (app *guiApp) askNewPassphrase() ([]byte, bool) {
	var dlg *walk.Dialog
	var pw1, pw2 *walk.LineEdit
	var result []byte
	accepted := false
	Dialog{
		AssignTo: &dlg, Title: app.loc.DlgPassTitle,
		MinSize: Size{Width: 420, Height: 0},
		Layout:  VBox{Margins: Margins{Left: 12, Top: 12, Right: 12, Bottom: 8}},
		Children: []Widget{
			Label{Text: app.loc.DlgPassMsg},
			LineEdit{AssignTo: &pw1, PasswordMode: true},
			Label{Text: app.loc.DlgPassConfirm},
			LineEdit{AssignTo: &pw2, PasswordMode: true},
			Composite{
				Layout: HBox{MarginsZero: true, Spacing: 8},
				Children: []Widget{
					HSpacer{},
					PushButton{Text: "Cancel", OnClicked: func() { dlg.Cancel() }},
					PushButton{Text: "OK", OnClicked: func() {
						if pw1.Text() != pw2.Text() {
							walk.MsgBox(dlg, app.loc.DlgPassTitle, app.loc.DlgPassBadMsg, walk.MsgBoxIconWarning)
							return
						}
						if pw1.Text() != "" {
							result = []byte(pw1.Text())
						}
						accepted = true
						dlg.Accept()
					}},
				},
			},
		},
	}.Create(app.mw)
	dlg.Run()
	return result, accepted
}

func sectionLabel(text string) Widget {
	return Label{Text: text, Font: Font{Bold: true, PointSize: 10}}
}

// ---- actions ----------------------------------------------------------------

func (app *guiApp) doGenKeys() {
	passphrase, ok := app.askNewPassphrase()
	if !ok {
		return
	}
	if app.useMemory {
		privPEM, pubPEM, err := cc.GenerateKeyPairPEM(app.keyType, passphrase)
		if err != nil {
			walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrGenKeys, err), walk.MsgBoxIconError)
			return
		}
		if err := app.ks.SaveMyKey(privPEM, pubPEM); err != nil {
			walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrGenKeys, err), walk.MsgBoxIconError)
			return
		}
	} else {
		if err := cc.GenerateKeyPair(app.keyType, passphrase, app.ks.MyPriv, app.ks.MyPub); err != nil {
			walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrGenKeys, err), walk.MsgBoxIconError)
			return
		}
	}
	note := app.loc.NotePassphrase
	if len(passphrase) == 0 {
		note = app.loc.NoteNoPass
	}
	app.setOutput(fmt.Sprintf("[OK] %s key pair generated%s.", app.keyType.String(), note))
	app.setStatus(app.loc.StatusKeyGen)
}

func (app *guiApp) doCopyPub() {
	pubBytes, err := app.ks.ReadMyPub()
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, app.loc.ErrNoKey, walk.MsgBoxIconWarning)
		return
	}
	clipboard.Set(string(pubBytes))
	app.setStatus(app.loc.DlgPubCopied)
}

func (app *guiApp) importKey(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrImport, err), walk.MsgBoxIconError)
		return
	}
	if _, _, err := cc.ParsePublicKeyPEM(data); err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrImport, err), walk.MsgBoxIconError)
		return
	}
	if err := app.ks.SaveFriendKey(data); err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrImport, err), walk.MsgBoxIconError)
		return
	}
	app.friendKeyEdit.SetText(path)
	app.setStatus(app.loc.StatusKeyImp)
}

func (app *guiApp) doImportFile() {
	dlg := new(walk.FileDialog)
	dlg.Title = app.loc.BtnImportKey
	dlg.Filter = "PEM files (*.pem)|*.pem|All files (*.*)|*.*"
	ok, err := dlg.ShowOpen(app.mw)
	if err != nil || !ok {
		return
	}
	app.importKey(dlg.FilePath)
}

func (app *guiApp) doImportClipboard() {
	text, err := clipboard.Get()
	if err != nil || text == "" {
		walk.MsgBox(app.mw, app.loc.DlgClipImportTitle, app.loc.DlgClipImportBad, walk.MsgBoxIconWarning)
		return
	}
	// Validate as PEM public key.
	if _, _, err := cc.ParsePublicKeyPEM([]byte(text)); err != nil {
		walk.MsgBox(app.mw, app.loc.DlgClipImportTitle, app.loc.DlgClipImportBad, walk.MsgBoxIconWarning)
		return
	}
	if err := app.ks.SaveFriendKey([]byte(text)); err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrImport, err), walk.MsgBoxIconError)
		return
	}
	app.friendKeyEdit.SetText("(clipboard)")
	app.setStatus(app.loc.StatusKeyImp)
}

func (app *guiApp) doLoadImage() {
	dlg := new(walk.FileDialog)
	dlg.Title = app.loc.BtnLoadImage
	dlg.Filter = "Images (*.png;*.jpg;*.jpeg;*.gif;*.bmp)|*.png;*.jpg;*.jpeg;*.gif;*.bmp|All files (*.*)|*.*"
	ok, err := dlg.ShowOpen(app.mw)
	if err != nil || !ok {
		return
	}
	data, err := os.ReadFile(dlg.FilePath)
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf("Read image: %v", err), walk.MsgBoxIconError)
		return
	}
	// Re-encode as PNG for consistency.
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf("Decode image: %v", err), walk.MsgBoxIconError)
		return
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf("Encode image: %v", err), walk.MsgBoxIconError)
		return
	}
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	uri := "data:image/png;base64," + b64
	app.inputTE.SetText(uri)
	app.setStatus("Image loaded.")
}

func (app *guiApp) doLoadFile() {
	dlg := new(walk.FileDialog)
	dlg.Title = app.loc.BtnLoadFile
	dlg.Filter = "All files (*.*)|*.*|Text files (*.txt)|*.txt"
	ok, err := dlg.ShowOpen(app.mw)
	if err != nil || !ok {
		return
	}
	data, err := os.ReadFile(dlg.FilePath)
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf("Read file: %v", err), walk.MsgBoxIconError)
		return
	}
	app.inputTE.SetText(string(data))
	app.setStatus(fmt.Sprintf("Loaded %d bytes.", len(data)))
}

// loadFriendPub tries to load the friend's public key only.
func (app *guiApp) loadFriendPub() (key any, algo cc.KeyType, err error) {
	if app.useMemory {
		pemBytes, e := app.ks.ReadFrPub()
		if e != nil || len(pemBytes) == 0 {
			return nil, 0, fmt.Errorf("no friend key")
		}
		return cc.ParsePublicKeyPEM(pemBytes)
	}
	if fileExists(app.ks.FrPub) {
		return cc.LoadPublicKey(app.ks.FrPub)
	}
	return nil, 0, fmt.Errorf("no friend key")
}

// loadMyPub loads my own public key.
func (app *guiApp) loadMyPub() (key any, algo cc.KeyType, err error) {
	if app.useMemory {
		pemBytes, e := app.ks.ReadMyPub()
		if e != nil || len(pemBytes) == 0 {
			return nil, 0, fmt.Errorf("no my key")
		}
		return cc.ParsePublicKeyPEM(pemBytes)
	}
	if fileExists(app.ks.MyPub) {
		return cc.LoadPublicKey(app.ks.MyPub)
	}
	return nil, 0, fmt.Errorf("no my key")
}

func (app *guiApp) doEncrypt() {
	pub, algo, err := app.loadFriendPub()
	if err != nil {
		// No friend key – ask if user wants to use own key.
		pub, algo, err = app.loadMyPub()
		if err != nil {
			walk.MsgBox(app.mw, app.loc.AppTitle, app.loc.ErrNoKey, walk.MsgBoxIconWarning)
			return
		}
		if walk.MsgBox(app.mw, app.loc.DlgNoFriendTitle, app.loc.DlgNoFriendMsg,
			walk.MsgBoxIconQuestion|walk.MsgBoxYesNo) != walk.DlgCmdYes {
			return
		}
	}
	text := app.getInput()
	if text == "" {
		return
	}
	ct, err := cc.Encrypt(pub, algo, []byte(text))
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrEncrypt, err), walk.MsgBoxIconError)
		return
	}
	packed := cc.Pack(ct, app.useCompress)
	app.setOutput(packed)
	clipboard.Set(packed)
	app.setStatus(app.loc.StatusEncOK)
}

func (app *guiApp) doDecrypt() {
	var privKey any
	var algo cc.KeyType
	var err error

	if app.useMemory {
		pemBytes, e := app.ks.ReadMyPriv()
		if e != nil || len(pemBytes) == 0 {
			walk.MsgBox(app.mw, app.loc.AppTitle, app.loc.ErrNoPrivKey, walk.MsgBoxIconWarning)
			return
		}
		pass, ok := app.askPassphrase(app.loc.DlgDecPassTitle, app.loc.DlgDecPassMsg)
		if !ok {
			return
		}
		var passphrase []byte
		if pass != "" {
			passphrase = []byte(pass)
		}
		privKey, algo, err = cc.ParsePrivateKeyPEM(pemBytes, passphrase)
	} else {
		if !fileExists(app.ks.MyPriv) {
			walk.MsgBox(app.mw, app.loc.AppTitle, app.loc.ErrNoPrivKey, walk.MsgBoxIconWarning)
			return
		}
		pass, ok := app.askPassphrase(app.loc.DlgDecPassTitle, app.loc.DlgDecPassMsg)
		if !ok {
			return
		}
		var passphrase []byte
		if pass != "" {
			passphrase = []byte(pass)
		}
		privKey, algo, err = cc.LoadPrivateKey(app.ks.MyPriv, passphrase)
	}
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrDecrypt, err), walk.MsgBoxIconError)
		return
	}

	packed := app.getInput()
	raw, err := cc.Unpack(packed)
	if err != nil {
		// Try legacy hex format
		if ct, e2 := hex.DecodeString(packed); e2 == nil {
			raw = ct
		} else {
			walk.MsgBox(app.mw, app.loc.AppTitle, app.loc.ErrBadHex, walk.MsgBoxIconWarning)
			return
		}
	}
	pt, err := cc.Decrypt(privKey, algo, raw)
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrDecrypt, err), walk.MsgBoxIconError)
		return
	}
	plain := string(pt)

	// Auto-detect image
	if cc.IsImageBase64(plain) {
		if walk.MsgBox(app.mw, app.loc.DlgImageViewTitle, app.loc.DlgImageAskMsg,
			walk.MsgBoxIconQuestion|walk.MsgBoxYesNo) == walk.DlgCmdYes {
			app.showImageViewer(plain)
			return
		}
	}
	app.setOutput(plain)
	app.setStatus(app.loc.StatusDecOK)
}

func (app *guiApp) showImageViewer(dataURI string) {
	b64 := dataURI[len(cc.ImagePrefix):]
	if idx := strings.Index(b64, ";base64,"); idx >= 0 {
		b64 = b64[idx+8:]
	}
	imgData, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, "Failed to decode image.", walk.MsgBoxIconError)
		return
	}
	src, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, "Failed to decode image.", walk.MsgBoxIconError)
		return
	}

	// Scale down to fit screen (max 85% of work area).
	bounds := src.Bounds()
	sw, sh := bounds.Dx(), bounds.Dy()
	maxW := int(win.GetSystemMetrics(win.SM_CXSCREEN)) * 85 / 100
	maxH := int(win.GetSystemMetrics(win.SM_CYSCREEN)) * 80 / 100
	if sw > maxW || sh > maxH {
		scale := 1.0
		if float64(sw)/float64(maxW) > float64(sh)/float64(maxH) {
			scale = float64(maxW) / float64(sw)
		} else {
			scale = float64(maxH) / float64(sh)
		}
		nw := int(float64(sw) * scale)
		nh := int(float64(sh) * scale)
		if nw < 1 {
			nw = 1
		}
		if nh < 1 {
			nh = 1
		}
		dst := image.NewRGBA(image.Rect(0, 0, nw, nh))
		drawApproxScale(dst, src)
		src = dst
	}

	bmp, err := walk.NewBitmapFromImage(src)
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, "Failed to create bitmap.", walk.MsgBoxIconError)
		return
	}

	var dlg *walk.Dialog
	var iv *walk.ImageView
	Dialog{
		AssignTo: &dlg, Title: app.loc.DlgImageViewTitle,
		Layout: VBox{},
		Children: []Widget{
			ImageView{AssignTo: &iv, Image: bmp},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{Text: "Save...", OnClicked: func() {
						sdlg := new(walk.FileDialog)
						sdlg.Title = "Save Image"
						sdlg.Filter = "PNG (*.png)|*.png"
						if ok, _ := sdlg.ShowSave(dlg); ok {
							os.WriteFile(sdlg.FilePath, imgData, 0o644)
						}
					}},
					HSpacer{},
					PushButton{Text: "Close", OnClicked: func() { dlg.Accept() }},
				},
			},
		},
	}.Create(app.mw)
	dlg.Run()
}

// drawApproxScale performs a nearest-neighbour scale of src into dst.
func drawApproxScale(dst *image.RGBA, src image.Image) {
	sb := src.Bounds()
	db := dst.Bounds()
	dw, dh := db.Dx(), db.Dy()
	sw, sh := sb.Dx(), sb.Dy()
	for y := 0; y < dh; y++ {
		sy := y * sh / dh
		for x := 0; x < dw; x++ {
			sx := x * sw / dw
			dst.Set(x, y, src.At(sb.Min.X+sx, sb.Min.Y+sy))
		}
	}
}

func (app *guiApp) doClear() {
	app.inputTE.SetText("")
	app.setOutput("")
	app.setStatus(app.loc.StatusCleared)
}

func (app *guiApp) doCopyOutput() {
	text := app.outputTE.Text()
	if text != "" {
		clipboard.Set(text)
		app.setStatus(app.loc.StatusCopied)
	}
}

func (app *guiApp) doSaveOutput() {
	text := app.outputTE.Text()
	if text == "" {
		return
	}
	dlg := new(walk.FileDialog)
	dlg.Title = app.loc.BtnSaveFile
	ok, err := dlg.ShowSave(app.mw)
	if err != nil || !ok {
		return
	}
	if err := os.WriteFile(dlg.FilePath, []byte(text), 0o644); err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf("Save: %v", err), walk.MsgBoxIconError)
		return
	}
	app.setStatus(fmt.Sprintf("Saved to %s", filepath.Base(dlg.FilePath)))
}

func (app *guiApp) doCalcHash() {
	dlg := new(walk.FileDialog)
	dlg.Title = app.loc.BtnCalcHash
	dlg.Filter = "All files (*.*)|*.*"
	ok, err := dlg.ShowOpen(app.mw)
	if err != nil || !ok {
		return
	}
	sums, err := filehash.Compute(dlg.FilePath)
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, app.loc.HashNoFile, walk.MsgBoxIconWarning)
		return
	}
	app.setOutput(fmt.Sprintf("%s\n  "+app.loc.HashMD5+"\n  "+app.loc.HashSHA1+"\n  "+app.loc.HashSHA256,
		dlg.FilePath, sums.MD5, sums.SHA1, sums.SHA256))
	app.setStatus(app.loc.StatusHashOK)
}

func (app *guiApp) doVerifyHash() {
	var dlg *walk.Dialog
	var fileEdit, hashEdit *walk.LineEdit
	Dialog{
		AssignTo: &dlg, Title: app.loc.DlgVerifyTitle,
		MinSize: Size{Width: 480, Height: 0},
		Layout:  VBox{Margins: Margins{Left: 12, Top: 12, Right: 12, Bottom: 8}},
		Children: []Widget{
			Label{Text: app.loc.DlgVerifyFile},
			Composite{
				Layout: HBox{MarginsZero: true, Spacing: 4},
				Children: []Widget{
					LineEdit{AssignTo: &fileEdit},
					PushButton{Text: "...", MaxSize: Size{Width: 30, Height: 22}, OnClicked: func() {
						fdlg := new(walk.FileDialog)
						fdlg.Title = app.loc.DlgVerifyTitle
						if ok, _ := fdlg.ShowOpen(dlg); ok {
							fileEdit.SetText(fdlg.FilePath)
						}
					}},
				},
			},
			Label{Text: app.loc.DlgVerifyHash},
			LineEdit{AssignTo: &hashEdit},
			Composite{
				Layout: HBox{MarginsZero: true, Spacing: 8},
				Children: []Widget{
					HSpacer{},
					PushButton{Text: "Cancel", OnClicked: func() { dlg.Cancel() }},
					PushButton{Text: "OK", OnClicked: func() {
						path := strings.TrimSpace(fileEdit.Text())
						expected := strings.TrimSpace(hashEdit.Text())
						if path == "" || expected == "" {
							return
						}
						algo, err := filehash.Compare(path, expected)
						if err != nil || algo == "" {
							walk.MsgBox(dlg, app.loc.AppTitle, app.loc.CmpFail, walk.MsgBoxIconWarning)
						} else {
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
							walk.MsgBox(dlg, app.loc.AppTitle,
								fmt.Sprintf(algo+"\n"+app.loc.CmpMatch+"\n"+app.loc.CmpMatchVal, path, digest),
								walk.MsgBoxIconInformation)
						}
						dlg.Accept()
					}},
				},
			},
		},
	}.Create(app.mw)
	dlg.Run()
}

func (app *guiApp) doEncryptFile() {
	// Select input file
	idlg := new(walk.FileDialog)
	idlg.Title = "Select file to encrypt"
	if ok, _ := idlg.ShowOpen(app.mw); !ok {
		return
	}
	data, err := os.ReadFile(idlg.FilePath)
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf("Read: %v", err), walk.MsgBoxIconError)
		return
	}

	pub, algo, err := app.loadFriendPub()
	if err != nil {
		pub, algo, err = app.loadMyPub()
		if err != nil {
			walk.MsgBox(app.mw, app.loc.AppTitle, app.loc.ErrNoKey, walk.MsgBoxIconWarning)
			return
		}
		if walk.MsgBox(app.mw, app.loc.DlgNoFriendTitle, app.loc.DlgNoFriendMsg,
			walk.MsgBoxIconQuestion|walk.MsgBoxYesNo) != walk.DlgCmdYes {
			return
		}
	}

	ct, err := cc.Encrypt(pub, algo, data)
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrEncrypt, err), walk.MsgBoxIconError)
		return
	}

	// Save as raw binary (no base64, minimal overhead).
	sdlg := new(walk.FileDialog)
	sdlg.Title = "Save encrypted file"
	sdlg.Filter = "Encrypted (*.cry)|*.cry|All files (*.*)|*.*"
	if ok, _ := sdlg.ShowSave(app.mw); !ok {
		return
	}
	if err := os.WriteFile(sdlg.FilePath, ct, 0o644); err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf("Save: %v", err), walk.MsgBoxIconError)
		return
	}
	app.setStatus(fmt.Sprintf("Encrypted: %s", filepath.Base(sdlg.FilePath)))
}

func (app *guiApp) doDecryptFile() {
	idlg := new(walk.FileDialog)
	idlg.Title = "Select encrypted file (.cry)"
	if ok, _ := idlg.ShowOpen(app.mw); !ok {
		return
	}
	raw, err := os.ReadFile(idlg.FilePath)
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf("Read: %v", err), walk.MsgBoxIconError)
		return
	}

	var privKey any
	var algo cc.KeyType
	if app.useMemory {
		pemBytes, e := app.ks.ReadMyPriv()
		if e != nil || len(pemBytes) == 0 {
			walk.MsgBox(app.mw, app.loc.AppTitle, app.loc.ErrNoPrivKey, walk.MsgBoxIconWarning)
			return
		}
		pass, ok := app.askPassphrase(app.loc.DlgDecPassTitle, app.loc.DlgDecPassMsg)
		if !ok {
			return
		}
		var passphrase []byte
		if pass != "" {
			passphrase = []byte(pass)
		}
		privKey, algo, err = cc.ParsePrivateKeyPEM(pemBytes, passphrase)
	} else {
		if !fileExists(app.ks.MyPriv) {
			walk.MsgBox(app.mw, app.loc.AppTitle, app.loc.ErrNoPrivKey, walk.MsgBoxIconWarning)
			return
		}
		pass, ok := app.askPassphrase(app.loc.DlgDecPassTitle, app.loc.DlgDecPassMsg)
		if !ok {
			return
		}
		var passphrase []byte
		if pass != "" {
			passphrase = []byte(pass)
		}
		privKey, algo, err = cc.LoadPrivateKey(app.ks.MyPriv, passphrase)
	}
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrDecrypt, err), walk.MsgBoxIconError)
		return
	}

	pt, err := cc.Decrypt(privKey, algo, raw)
	if err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf(app.loc.ErrDecrypt, err), walk.MsgBoxIconError)
		return
	}

	// Auto-detect image and offer preview.
	if _, _, e := image.Decode(bytes.NewReader(pt)); e == nil {
		if walk.MsgBox(app.mw, app.loc.DlgImageViewTitle, app.loc.DlgImageAskMsg,
			walk.MsgBoxIconQuestion|walk.MsgBoxYesNo) == walk.DlgCmdYes {
			b64 := base64.StdEncoding.EncodeToString(pt)
			app.showImageViewer("data:image/png;base64," + b64)
			return // viewer has its own Save button, skip main save dialog
		}
	}

	sdlg := new(walk.FileDialog)
	sdlg.Title = "Save decrypted file"
	if ok, _ := sdlg.ShowSave(app.mw); !ok {
		return
	}
	if err := os.WriteFile(sdlg.FilePath, pt, 0o644); err != nil {
		walk.MsgBox(app.mw, app.loc.AppTitle, fmt.Sprintf("Save: %v", err), walk.MsgBoxIconError)
		return
	}
	app.setStatus(fmt.Sprintf("Decrypted: %s", filepath.Base(sdlg.FilePath)))
}

// ---- helpers ----------------------------------------------------------------

func algoNames() []string {
	algos := cc.AllAlgos()
	names := make([]string, len(algos))
	for i, a := range algos {
		names[i] = a.String()
	}
	return names
}

func algoFromName(name string) cc.KeyType {
	for _, a := range cc.AllAlgos() {
		if a.String() == name {
			return a
		}
	}
	return cc.KeyRSA2048
}

// ---- build & run ------------------------------------------------------------

func runGUI(loc *locale.Locale) {
	console := win.GetConsoleWindow()
	if console != 0 {
		win.ShowWindow(console, win.SW_HIDE)
	}
	icc := iccex{Size: uint32(unsafe.Sizeof(iccex{})), ICC: iccStandard}
	initCommonControlsExFn.Call(uintptr(unsafe.Pointer(&icc)))

	// Use disk store by default; user can switch to memory.
	ks, err := keystore.Init()
	if err != nil {
		walk.MsgBox(nil, "Error", fmt.Sprintf("Keystore init: %v", err), walk.MsgBoxIconError)
		os.Exit(1)
	}
	app := &guiApp{loc: loc, ks: ks, keyType: cc.KeyRSA2048, useMemory: true, useCompress: true}
	ks.Disk = false // default memory mode

	friendPath := ""
	if fileExists(ks.FrPub) {
		friendPath, _ = filepath.Abs(ks.FrPub)
	}

	if err := (MainWindow{
		AssignTo: &app.mw,
		Title:    loc.AppTitle,
		MinSize:  Size{Width: 680, Height: 580},
		Size:     Size{Width: 720, Height: 640},
		Layout:   VBox{Margins: Margins{Left: 10, Top: 6, Right: 10, Bottom: 4}},

		Children: []Widget{
			// [Key Management]
			sectionLabel(loc.LabelKeyMgmt),
			Composite{
				Layout: HBox{Spacing: 8},
				Children: []Widget{
					Label{Text: loc.LabelAlgorithm, MinSize: Size{Width: 70, Height: 0}},
					ComboBox{
						AssignTo: &app.algoCombo, Model: algoNames(), Value: "RSA-2048",
						MinSize:               Size{Width: 130, Height: 0},
						OnCurrentIndexChanged: func() { app.keyType = algoFromName(app.algoCombo.Text()) },
					},
					PushButton{Text: loc.BtnGenKey, MaxSize: Size{Width: 150, Height: 26}, OnClicked: app.doGenKeys},
					PushButton{Text: loc.BtnCopyPub, MaxSize: Size{Width: 150, Height: 26}, OnClicked: app.doCopyPub},
				},
			},
			Composite{
				Layout: HBox{Spacing: 8},
				Children: []Widget{
					Label{Text: loc.LabelFriendKey, MinSize: Size{Width: 70, Height: 0}},
					LineEdit{AssignTo: &app.friendKeyEdit, Text: friendPath, ReadOnly: true},
					PushButton{Text: loc.BtnImport, MaxSize: Size{Width: 80, Height: 26}, OnClicked: app.doImportFile},
					PushButton{Text: loc.BtnImportClip, MaxSize: Size{Width: 150, Height: 26}, OnClicked: app.doImportClipboard},
				},
			},
			CheckBox{AssignTo: &app.memCheck, Text: loc.LabelMemory, Checked: true,
				OnCheckedChanged: func() {
					app.useMemory = app.memCheck.Checked()
					app.ks.Disk = !app.useMemory
				}},

			// [Input]
			Composite{
				Layout: HBox{Spacing: 4},
				Children: []Widget{
					sectionLabel(loc.LabelInput),
					PushButton{Text: loc.BtnLoadFile, MaxSize: Size{Width: 100, Height: 22}, OnClicked: app.doLoadFile},
				},
			},
			TextEdit{AssignTo: &app.inputTE, VScroll: true, MinSize: Size{Width: 0, Height: 140}},
			Composite{
				Layout: HBox{Spacing: 8, Margins: Margins{Top: 4}},
				Children: []Widget{
					PushButton{Text: loc.BtnEncrypt, MaxSize: Size{Width: 110, Height: 28}, OnClicked: app.doEncrypt},
					PushButton{Text: loc.BtnClear, MaxSize: Size{Width: 80, Height: 28}, OnClicked: app.doClear},
					PushButton{Text: loc.BtnDecrypt, MaxSize: Size{Width: 110, Height: 28}, OnClicked: app.doDecrypt},
					PushButton{Text: loc.BtnLoadImage, MaxSize: Size{Width: 100, Height: 28}, OnClicked: app.doLoadImage},
				},
			},

			// [Output]
			sectionLabel(loc.LabelOutput),
			TextEdit{AssignTo: &app.outputTE, ReadOnly: true, VScroll: true, MinSize: Size{Width: 0, Height: 120}},
			Composite{
				Layout: HBox{Margins: Margins{Top: 4}, Spacing: 8},
				Children: []Widget{
					PushButton{Text: loc.BtnCopyOutput, MaxSize: Size{Width: 120, Height: 28}, OnClicked: app.doCopyOutput},
					PushButton{Text: loc.BtnSaveFile, MaxSize: Size{Width: 120, Height: 28}, OnClicked: app.doSaveOutput},
					CheckBox{AssignTo: &app.compressCheck, Text: loc.LabelCompress, Checked: true,
						OnCheckedChanged: func() { app.useCompress = app.compressCheck.Checked() }},
				},
			},

			// [Tools]
			sectionLabel(loc.LabelTools),
			Composite{
				Layout: HBox{Spacing: 8},
				Children: []Widget{
					PushButton{Text: loc.BtnCalcHash, MaxSize: Size{Width: 150, Height: 28}, OnClicked: app.doCalcHash},
					PushButton{Text: loc.BtnVerifyHash, MaxSize: Size{Width: 150, Height: 28}, OnClicked: app.doVerifyHash},
					PushButton{Text: loc.BtnEncFile, MaxSize: Size{Width: 130, Height: 28}, OnClicked: app.doEncryptFile},
					PushButton{Text: loc.BtnDecFile, MaxSize: Size{Width: 130, Height: 28}, OnClicked: app.doDecryptFile},
				},
			},

			VSpacer{},
			Composite{
				Layout: HBox{Margins: Margins{Top: 2}},
				Children: []Widget{
					Label{AssignTo: &app.statusLabel, Text: loc.StatusReady},
					HSpacer{},
					Label{Text: loc.Version, TextColor: walk.RGB(140, 140, 140)},
				},
			},
		},
	}.Create()); err != nil {
		walk.MsgBox(nil, "Error", fmt.Sprintf("Failed to create window: %v", err), walk.MsgBoxIconError)
		os.Exit(1)
	}
	// Set window icon from embedded PNG.
	if img, err := png.Decode(bytes.NewReader(logoPNG)); err == nil {
		if bmp, err := walk.NewBitmapFromImage(img); err == nil {
			if icon, err := walk.NewIconFromBitmap(bmp); err == nil {
				app.mw.SetIcon(icon)
			}
		}
	}
	app.mw.Run()
}
