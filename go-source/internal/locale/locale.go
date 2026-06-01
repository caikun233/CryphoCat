// Package locale provides i18n strings for CryphoCat.
package locale

import (
	"os"
	"strings"
)

// Locale holds all translatable strings used by CLI and GUI.
type Locale struct {
	// App
	AppTitle string

	// Buttons / Menu items
	BtnGenKey     string
	BtnImportKey  string
	BtnImport     string
	BtnImportClip string
	BtnCopyPub    string
	BtnEncrypt    string
	BtnDecrypt    string
	BtnClear      string
	BtnCopyOutput string
	BtnCalcHash   string
	BtnVerifyHash string
	BtnLoadImage  string
	BtnLoadFile   string
	BtnSaveFile   string
	BtnEncFile    string
	BtnDecFile    string
	MenuHash      string
	MenuVerify    string
	MenuExit      string

	// Labels
	LabelInput     string
	LabelOutput    string
	LabelKeyInfo   string
	LabelAlgorithm string
	LabelFriendKey string
	LabelKeyMgmt   string
	LabelTools     string
	LabelMemory    string
	LabelCompress  string

	// Status
	StatusReady    string
	StatusEncOK    string
	StatusDecOK    string
	StatusCleared  string
	StatusCopied   string
	StatusKeyGen   string
	StatusKeyImp   string
	StatusHashOK   string
	StatusHashFail string

	// Dialogs
	DlgPassTitle       string
	DlgPassMsg         string
	DlgPassConfirm     string
	DlgPassBadMsg      string
	DlgDecPassTitle    string
	DlgDecPassMsg      string
	DlgNoFriendTitle   string
	DlgNoFriendMsg     string
	DlgImageViewTitle  string
	DlgImageAskMsg     string
	DlgClipImportTitle string
	DlgClipImportBad   string
	DlgPubCopied       string

	// Key generation
	KeyLenPrompt    string
	KeyLenBad       string
	KeyLenNotNum    string
	PromptPWSet     string
	PromptPWConfirm string
	PwMismatch      string
	PwWarn          string
	KeyGenOK        string
	KeyPriv         string
	KeyPub          string

	// Import
	ImportPrompt string
	ImportOK     string
	ImportFail   string
	ImportNoFile string

	// Encrypt
	EncNoKey     string
	EncUseMine   string
	EncCancelled string
	EncNoKeyAny  string
	EncLoadFail  string
	EncPrompt    string
	EncEmpty     string
	EncFail      string
	EncOK        string

	// Decrypt
	DecNoKey    string
	PromptPWDec string
	DecLoadFail string
	DecPrompt   string
	DecBadHex   string
	DecFail     string
	DecOK       string

	// Hash
	HashPrompt string
	HashNoFile string
	HashMD5    string
	HashSHA1   string
	HashSHA256 string

	// Compare
	CmpPrompt      string
	CmpHashIn      string
	CmpNoFile      string
	CmpEmpty       string
	CmpMatch       string
	CmpMatchVal    string
	CmpFail        string
	DlgVerifyTitle string
	DlgVerifyFile  string
	DlgVerifyHash  string

	// Menu
	MenuSep      string
	ChoicePrompt string
	ChoiceBad    string
	PressEnter   string
	Goodbye      string

	// Errors
	ErrGenKeys   string
	ErrImport    string
	ErrEncrypt   string
	ErrDecrypt   string
	ErrNoKey     string
	ErrBadHex    string
	ErrNoPrivKey string

	// Misc
	NotePassphrase string
	NoteNoPass     string
	GUIOnlyWindows string
	Version        string
}

var localeEN = Locale{
	AppTitle: "CryphoCat",

	BtnGenKey:     "Gen Keys",
	BtnImportKey:  "Import Friend Key",
	BtnImport:     "Import...",
	BtnImportClip: "Import from Clipboard",
	BtnCopyPub:    "Copy My Public Key",
	BtnEncrypt:    "Encrypt",
	BtnDecrypt:    "Decrypt",
	BtnClear:      "Clear",
	BtnCopyOutput: "Copy Output",
	BtnCalcHash:   "Calculate Hash",
	BtnVerifyHash: "Verify Hash",
	BtnLoadImage:  "Load Image",
	BtnLoadFile:   "Load File",
	BtnSaveFile:   "Save Output",
	BtnEncFile:    "Encrypt File",
	BtnDecFile:    "Decrypt File",
	MenuHash:      "Compute File Hashes",
	MenuVerify:    "Verify File Hash",
	MenuExit:      "Exit",

	LabelInput:     "INPUT:",
	LabelOutput:    "OUTPUT:",
	LabelKeyInfo:   "RSA / OAEP-SHA256",
	LabelAlgorithm: "Algorithm:",
	LabelFriendKey: "Friend's key:",
	LabelKeyMgmt:   "[Key Management]",
	LabelTools:     "[Tools]",
	LabelMemory:    "Store keys in memory only",
	LabelCompress:  "Compress output",

	StatusReady:    "Ready.",
	StatusEncOK:    "Encrypted - copied to clipboard.",
	StatusDecOK:    "Decrypted.",
	StatusCleared:  "Cleared.",
	StatusCopied:   "Output copied to clipboard.",
	StatusKeyGen:   "Keys generated.",
	StatusKeyImp:   "Friend's key imported.",
	StatusHashOK:   "Hash computed.",
	StatusHashFail: "Hash verification failed.",

	DlgPassTitle:       "Set Passphrase",
	DlgPassMsg:         "Passphrase for private key (blank = no encryption):",
	DlgPassConfirm:     "Confirm passphrase:",
	DlgPassBadMsg:      "Passphrases do not match.",
	DlgDecPassTitle:    "Passphrase",
	DlgDecPassMsg:      "Private key passphrase (blank if none):",
	DlgNoFriendTitle:   "No friend key",
	DlgNoFriendMsg:     "Friend's public key not found.\nUse your own public key?",
	DlgImageViewTitle:  "Image Preview",
	DlgImageAskMsg:     "Decryption result appears to be an image.\nView it?",
	DlgClipImportTitle: "Import from Clipboard",
	DlgClipImportBad:   "Clipboard does not contain a valid public key.",
	DlgPubCopied:       "My public key copied to clipboard.",

	KeyLenPrompt:    "Algorithm (1-7, default 3=RSA-4096): ",
	KeyLenBad:       "Choose one of: 2048, 4096",
	KeyLenNotNum:    "Please enter a number.",
	PromptPWSet:     "Set passphrase (blank = no encryption, NOT recommended): ",
	PromptPWConfirm: "Confirm passphrase: ",
	PwMismatch:      "Passphrases do not match.",
	PwWarn:          "[WARNING] Private key will be stored without a passphrase.",
	KeyGenOK:        "RSA-%d key pair generated.",
	KeyPriv:         "  Private: %s",
	KeyPub:          "  Public : %s",

	ImportPrompt: "Path to friend's public key: ",
	ImportOK:     "Public key imported successfully.",
	ImportFail:   "Import failed: %v",
	ImportNoFile: "File not found. Check the path.",

	EncNoKey:     "No friend key found.",
	EncUseMine:   "Use your own public key? (y/n): ",
	EncCancelled: "Cancelled. Import a friend's key or generate your key pair first.",
	EncNoKeyAny:  "No public key available. Generate a key pair first.",
	EncLoadFail:  "Cannot load public key: %v",
	EncPrompt:    "Text to encrypt: ",
	EncEmpty:     "Nothing to encrypt.",
	EncFail:      "Encryption failed: %v",
	EncOK:        "Encrypted (copied to clipboard): %s",

	DecNoKey:    "Private key not found. Generate a key pair first.",
	PromptPWDec: "Private key passphrase (blank if none): ",
	DecLoadFail: "Cannot load private key: %v",
	DecPrompt:   "Ciphertext (hex): ",
	DecBadHex:   "Invalid hex ciphertext.",
	DecFail:     "Decryption failed: %v",
	DecOK:       "Decrypted: %s",

	HashPrompt: "File path: ",
	HashNoFile: "File not found. Check the path.",
	HashMD5:    "MD5   : %s",
	HashSHA1:   "SHA1  : %s",
	HashSHA256: "SHA256: %s",

	CmpPrompt:      "File path: ",
	CmpHashIn:      "Expected hash (MD5/SHA1/SHA256): ",
	CmpNoFile:      "File not found. Check the path.",
	CmpEmpty:       "No hash provided.",
	CmpMatch:       "%s MATCH: %s",
	CmpMatchVal:    "Value : %s",
	CmpFail:        "HASH CHECK FAILED - file may have been tampered with.",
	DlgVerifyTitle: "Verify File Hash",
	DlgVerifyFile:  "File path:",
	DlgVerifyHash:  "Expected hash (MD5/SHA1/SHA256):",

	MenuSep:      "----------------",
	ChoicePrompt: "Choice: ",
	ChoiceBad:    "Invalid choice.",
	PressEnter:   "Press Enter to continue...",
	Goodbye:      "Goodbye.",

	ErrGenKeys:   "Key generation failed: %v",
	ErrImport:    "Import failed: %v",
	ErrEncrypt:   "Encryption failed: %v",
	ErrDecrypt:   "Decryption failed: %v",
	ErrNoKey:     "No public key found. Generate a key pair first.",
	ErrBadHex:    "Input is not valid hex ciphertext.",
	ErrNoPrivKey: "Private key not found. Generate a key pair first.",

	NotePassphrase: " [passphrase protected]",
	NoteNoPass:     " [WARNING: no passphrase]",
	GUIOnlyWindows: "GUI is only available on Windows.",
	Version:        "v2.0.0",
}

var localeZHCN = Locale{
	AppTitle: "CryphoCat",

	BtnGenKey:     "生成密钥对",
	BtnImportKey:  "导入对方公钥",
	BtnImport:     "导入...",
	BtnImportClip: "从剪贴板导入",
	BtnCopyPub:    "复制我的公钥",
	BtnEncrypt:    "加密",
	BtnDecrypt:    "解密",
	BtnClear:      "清空",
	BtnCopyOutput: "复制输出",
	BtnCalcHash:   "计算文件哈希",
	BtnVerifyHash: "校验文件哈希",
	BtnLoadImage:  "加载图片",
	BtnLoadFile:   "从文件加载",
	BtnSaveFile:   "保存输出",
	BtnEncFile:    "加密文件",
	BtnDecFile:    "解密文件",
	MenuHash:      "计算文件哈希",
	MenuVerify:    "校验文件哈希",
	MenuExit:      "退出",

	LabelInput:     "输入：",
	LabelOutput:    "输出：",
	LabelKeyInfo:   "RSA / OAEP-SHA256",
	LabelAlgorithm: "算法：",
	LabelFriendKey: "对方公钥：",
	LabelKeyMgmt:   "[密钥管理]",
	LabelTools:     "[工具]",
	LabelMemory:    "仅存储在内存中",
	LabelCompress:  "压缩输出",

	StatusReady:    "就绪。",
	StatusEncOK:    "加密完成 - 已自动复制到剪贴板。",
	StatusDecOK:    "解密完成。",
	StatusCleared:  "已清空。",
	StatusCopied:   "输出已复制到剪贴板。",
	StatusKeyGen:   "密钥对已生成。",
	StatusKeyImp:   "对方公钥已导入。",
	StatusHashOK:   "哈希已计算。",
	StatusHashFail: "哈希校验失败。",

	DlgPassTitle:       "设置口令",
	DlgPassMsg:         "设置私钥保护口令（留空 = 不加密存储，不推荐）：",
	DlgPassConfirm:     "确认口令：",
	DlgPassBadMsg:      "两次输入的口令不一致。",
	DlgDecPassTitle:    "输入口令",
	DlgDecPassMsg:      "私钥口令（留空表示无口令）：",
	DlgNoFriendTitle:   "未找到对方公钥",
	DlgNoFriendMsg:     "未找到对方公钥。\n是否使用自己的公钥加密？",
	DlgImageViewTitle:  "图片预览",
	DlgImageAskMsg:     "解密结果似乎是图片。\n是否查看？",
	DlgClipImportTitle: "从剪贴板导入",
	DlgClipImportBad:   "剪贴板中不含有效的公钥。",
	DlgPubCopied:       "我的公钥已复制到剪贴板。",

	KeyLenPrompt:    "算法 (1-7，默认 3=RSA-4096)：",
	KeyLenBad:       "请从 2048、4096 中选择一个。",
	KeyLenNotNum:    "请输入数字。",
	PromptPWSet:     "设置私钥口令（留空 = 不加密存储，不推荐）：",
	PromptPWConfirm: "确认口令：",
	PwMismatch:      "两次输入的口令不一致。",
	PwWarn:          "[警告] 私钥将以明文形式存储，安全性较低。",
	KeyGenOK:        "RSA-%d 密钥对已生成。",
	KeyPriv:         "  私钥：%s",
	KeyPub:          "  公钥：%s",

	ImportPrompt: "请输入对方公钥文件路径：",
	ImportOK:     "公钥导入成功。",
	ImportFail:   "导入失败：%v",
	ImportNoFile: "文件不存在，请检查路径。",

	EncNoKey:     "未找到对方公钥。",
	EncUseMine:   "是否使用自己的公钥加密？(y/n)：",
	EncCancelled: "已取消。请先导入对方公钥或生成密钥对。",
	EncNoKeyAny:  "未找到公钥，请先生成密钥对。",
	EncLoadFail:  "无法加载公钥：%v",
	EncPrompt:    "请输入要加密的文本：",
	EncEmpty:     "输入为空。",
	EncFail:      "加密失败：%v",
	EncOK:        "加密结果（已复制到剪贴板）：%s",

	DecNoKey:    "未找到私钥，请先生成密钥对。",
	PromptPWDec: "私钥口令（无口令则留空）：",
	DecLoadFail: "无法加载私钥：%v",
	DecPrompt:   "请输入十六进制密文：",
	DecBadHex:   "十六进制密文无效。",
	DecFail:     "解密失败：%v",
	DecOK:       "解密结果：%s",

	HashPrompt: "请输入文件路径：",
	HashNoFile: "文件不存在，请检查路径。",
	HashMD5:    "MD5   ：%s",
	HashSHA1:   "SHA1  ：%s",
	HashSHA256: "SHA256：%s",

	CmpPrompt:      "请输入文件路径：",
	CmpHashIn:      "请输入哈希值（MD5/SHA1/SHA256）：",
	CmpNoFile:      "文件不存在，请检查路径。",
	CmpEmpty:       "未输入哈希值。",
	CmpMatch:       "%s 校验通过：%s",
	CmpMatchVal:    "哈希值：%s",
	CmpFail:        "哈希校验失败 —— 文件可能已被篡改。",
	DlgVerifyTitle: "校验文件哈希",
	DlgVerifyFile:  "文件路径：",
	DlgVerifyHash:  "期望哈希值（MD5/SHA1/SHA256）：",

	MenuSep:      "----------------",
	ChoicePrompt: "请输入选项：",
	ChoiceBad:    "选项无效，请重新输入。",
	PressEnter:   "按下回车继续...",
	Goodbye:      "再见。",

	ErrGenKeys:   "密钥生成失败：%v",
	ErrImport:    "导入失败：%v",
	ErrEncrypt:   "加密失败：%v",
	ErrDecrypt:   "解密失败：%v",
	ErrNoKey:     "未找到公钥，请先生成密钥对。",
	ErrBadHex:    "输入的十六进制密文无效。",
	ErrNoPrivKey: "未找到私钥，请先生成密钥对。",

	NotePassphrase: " [已设置口令保护]",
	NoteNoPass:     " [警告：未设置口令]",
	GUIOnlyWindows: "GUI 仅支持 Windows 系统。",
	Version:        "v2.0.0",
}

var localeJA = Locale{
	AppTitle: "CryphoCat",

	BtnGenKey: "鍵生成", BtnImportKey: "相手の公開鍵をインポート", BtnImport: "インポート...",
	BtnImportClip: "クリップボードから", BtnCopyPub: "自分の公開鍵をコピー",
	BtnEncrypt: "暗号化", BtnDecrypt: "復号", BtnClear: "クリア", BtnCopyOutput: "出力をコピー",
	BtnCalcHash: "ハッシュ計算", BtnVerifyHash: "ハッシュ検証", BtnLoadImage: "画像読込",
	BtnLoadFile: "ファイル読込",
	MenuHash:    "ファイルハッシュ計算", MenuVerify: "ファイルハッシュ検証", MenuExit: "終了",

	LabelInput: "入力：", LabelOutput: "出力：", LabelKeyInfo: "RSA / OAEP-SHA256",
	LabelAlgorithm: "アルゴリズム：", LabelFriendKey: "相手の公開鍵：",
	LabelKeyMgmt: "[鍵管理]", LabelTools: "[ツール]",
	LabelMemory: "メモリのみに保存", LabelCompress: "出力を圧縮",

	StatusReady: "準備完了。", StatusEncOK: "暗号化完了 - クリップボードにコピーしました。",
	StatusDecOK: "復号完了。", StatusCleared: "クリアしました。",
	StatusCopied: "出力をクリップボードにコピーしました。",
	StatusKeyGen: "鍵を生成しました。", StatusKeyImp: "相手の公開鍵をインポートしました。",
	StatusHashOK: "ハッシュ計算完了。", StatusHashFail: "ハッシュ検証に失敗しました。",

	DlgPassTitle: "パスフレーズ設定", DlgPassMsg: "秘密鍵のパスフレーズ（空白＝暗号化なし、非推奨）：",
	DlgPassConfirm: "パスフレーズ確認：", DlgPassBadMsg: "パスフレーズが一致しません。",
	DlgDecPassTitle: "パスフレーズ", DlgDecPassMsg: "秘密鍵のパスフレーズ（空白＝なし）：",
	DlgNoFriendTitle:  "相手の公開鍵なし",
	DlgNoFriendMsg:    "相手の公開鍵が見つかりません。\n自分の公開鍵で暗号化しますか？",
	DlgImageViewTitle: "画像プレビュー", DlgImageAskMsg: "復号結果は画像のようです。\n表示しますか？",
	DlgClipImportTitle: "クリップボードからインポート",
	DlgClipImportBad:   "クリップボードに有効な公開鍵がありません。",
	DlgPubCopied:       "自分の公開鍵をクリップボードにコピーしました。",

	KeyLenPrompt: "アルゴリズム (1-7、デフォルト 3=RSA-4096)：",
	KeyLenBad:    "1-7 の数字を入力してください。", KeyLenNotNum: "数字を入力してください。",
	PromptPWSet:     "秘密鍵のパスフレーズを設定（空白＝暗号化なし、非推奨）：",
	PromptPWConfirm: "パスフレーズ確認：",
	PwMismatch:      "パスフレーズが一致しません。",
	PwWarn:          "[警告] 秘密鍵がパスフレーズなしで保存されます。",
	KeyGenOK:        "%s 鍵ペアを生成しました。",
	KeyPriv:         "  秘密鍵：%s", KeyPub: "  公開鍵：%s",

	ImportPrompt: "相手の公開鍵ファイルのパス：", ImportOK: "公開鍵をインポートしました。",
	ImportFail: "インポート失敗：%v", ImportNoFile: "ファイルが見つかりません。",

	EncNoKey: "相手の公開鍵が見つかりません。", EncUseMine: "自分の公開鍵を使用しますか？(y/n)：",
	EncCancelled: "キャンセルされました。", EncNoKeyAny: "公開鍵がありません。先に鍵を生成してください。",
	EncLoadFail: "公開鍵の読み込みに失敗：%v", EncPrompt: "暗号化するテキスト：",
	EncEmpty: "入力が空です。", EncFail: "暗号化失敗：%v",
	EncOK: "暗号化（クリップボードにコピー）：%s",

	DecNoKey:    "秘密鍵が見つかりません。先に鍵を生成してください。",
	PromptPWDec: "秘密鍵のパスフレーズ（空白＝なし）：",
	DecLoadFail: "秘密鍵の読み込みに失敗：%v", DecPrompt: "暗号文（base64/hex）：",
	DecBadHex: "無効な暗号文です。", DecFail: "復号失敗：%v", DecOK: "復号：%s",

	HashPrompt: "ファイルパス：", HashNoFile: "ファイルが見つかりません。",
	HashMD5: "MD5   ：%s", HashSHA1: "SHA1  ：%s", HashSHA256: "SHA256：%s",

	CmpPrompt: "ファイルパス：", CmpHashIn: "期待するハッシュ値（MD5/SHA1/SHA256）：",
	CmpNoFile: "ファイルが見つかりません。", CmpEmpty: "ハッシュ値が入力されていません。",
	CmpMatch: "%s 一致：%s", CmpMatchVal: "値：%s", CmpFail: "ハッシュ検証失敗 — ファイルが改ざんされた可能性があります。",
	DlgVerifyTitle: "ハッシュ検証", DlgVerifyFile: "ファイルパス：",
	DlgVerifyHash: "期待するハッシュ値（MD5/SHA1/SHA256）：",

	MenuSep: "----------------", ChoicePrompt: "選択：", ChoiceBad: "無効な選択です。",
	PressEnter: "Enter を押して続行...", Goodbye: "さようなら。",

	ErrGenKeys: "鍵生成失敗：%v", ErrImport: "インポート失敗：%v",
	ErrEncrypt: "暗号化失敗：%v", ErrDecrypt: "復号失敗：%v",
	ErrNoKey:     "公開鍵がありません。先に鍵を生成してください。",
	ErrBadHex:    "無効な暗号文です。",
	ErrNoPrivKey: "秘密鍵がありません。先に鍵を生成してください。",

	NotePassphrase: " [パスフレーズ保護]", NoteNoPass: " [警告：パスフレーズなし]",
	GUIOnlyWindows: "GUI は Windows のみ対応しています。", Version: "v2.0.0",
}

var localeKO = Locale{
	AppTitle: "CryphoCat",

	BtnGenKey: "키 생성", BtnImportKey: "상대 공개키 가져오기", BtnImport: "가져오기...",
	BtnImportClip: "클립보드에서 가져오기", BtnCopyPub: "내 공개키 복사",
	BtnEncrypt: "암호화", BtnDecrypt: "복호화", BtnClear: "지우기", BtnCopyOutput: "출력 복사",
	BtnCalcHash: "해시 계산", BtnVerifyHash: "해시 검증", BtnLoadImage: "이미지 불러오기",
	BtnLoadFile: "파일 불러오기",
	MenuHash:    "파일 해시 계산", MenuVerify: "파일 해시 검증", MenuExit: "종료",

	LabelInput: "입력：", LabelOutput: "출력：", LabelKeyInfo: "RSA / OAEP-SHA256",
	LabelAlgorithm: "알고리즘：", LabelFriendKey: "상대 공개키：",
	LabelKeyMgmt: "[키 관리]", LabelTools: "[도구]",
	LabelMemory: "메모리에만 저장", LabelCompress: "출력 압축",

	StatusReady: "준비됨.", StatusEncOK: "암호화 완료 - 클립보드에 복사됨.",
	StatusDecOK: "복호화 완료.", StatusCleared: "지워짐.",
	StatusCopied: "출력을 클립보드에 복사함.", StatusKeyGen: "키 생성됨.", StatusKeyImp: "공개키 가져오기 완료.",
	StatusHashOK: "해시 계산 완료.", StatusHashFail: "해시 검증 실패.",

	DlgPassTitle: "암호 설정", DlgPassMsg: "개인키 암호 (공백=암호화 안 함, 권장하지 않음)：",
	DlgPassConfirm: "암호 확인：", DlgPassBadMsg: "암호가 일치하지 않습니다.",
	DlgDecPassTitle: "암호", DlgDecPassMsg: "개인키 암호 (공백=없음)：",
	DlgNoFriendTitle:  "상대 공개키 없음",
	DlgNoFriendMsg:    "상대 공개키를 찾을 수 없습니다.\n내 공개키로 암호화할까요?",
	DlgImageViewTitle: "이미지 미리보기", DlgImageAskMsg: "복호화 결과가 이미지인 것 같습니다.\n보시겠습니까?",
	DlgClipImportTitle: "클립보드에서 가져오기",
	DlgClipImportBad:   "클립보드에 유효한 공개키가 없습니다.",
	DlgPubCopied:       "내 공개키가 클립보드에 복사되었습니다.",

	KeyLenPrompt: "알고리즘 (1-7, 기본 3=RSA-4096)：",
	KeyLenBad:    "1-7 사이의 숫자를 입력하세요.", KeyLenNotNum: "숫자를 입력하세요.",
	PromptPWSet:     "개인키 암호 설정 (공백=암호화 안 함, 권장하지 않음)：",
	PromptPWConfirm: "암호 확인：", PwMismatch: "암호가 일치하지 않습니다.",
	PwWarn:   "[경고] 개인키가 암호 없이 저장됩니다.",
	KeyGenOK: "%s 키 쌍 생성됨.", KeyPriv: "  개인키：%s", KeyPub: "  공개키：%s",

	ImportPrompt: "상대 공개키 파일 경로：", ImportOK: "공개키 가져오기 완료.",
	ImportFail: "가져오기 실패：%v", ImportNoFile: "파일을 찾을 수 없습니다.",

	EncNoKey: "상대 공개키를 찾을 수 없습니다.", EncUseMine: "내 공개키를 사용할까요? (y/n)：",
	EncCancelled: "취소됨.", EncNoKeyAny: "공개키가 없습니다. 먼저 키를 생성하세요.",
	EncLoadFail: "공개키 로드 실패：%v", EncPrompt: "암호화할 텍스트：",
	EncEmpty: "입력이 비어 있습니다.", EncFail: "암호화 실패：%v",
	EncOK: "암호화됨 (클립보드에 복사)：%s",

	DecNoKey:    "개인키를 찾을 수 없습니다. 먼저 키를 생성하세요.",
	PromptPWDec: "개인키 암호 (공백=없음)：", DecLoadFail: "개인키 로드 실패：%v",
	DecPrompt: "암호문 (base64/hex)：", DecBadHex: "유효하지 않은 암호문입니다.",
	DecFail: "복호화 실패：%v", DecOK: "복호화：%s",

	HashPrompt: "파일 경로：", HashNoFile: "파일을 찾을 수 없습니다.",
	HashMD5: "MD5   ：%s", HashSHA1: "SHA1  ：%s", HashSHA256: "SHA256：%s",

	CmpPrompt: "파일 경로：", CmpHashIn: "예상 해시값 (MD5/SHA1/SHA256)：",
	CmpNoFile: "파일을 찾을 수 없습니다.", CmpEmpty: "해시값이 입력되지 않았습니다.",
	CmpMatch: "%s 일치：%s", CmpMatchVal: "값：%s",
	CmpFail:        "해시 검증 실패 — 파일이 변조되었을 수 있습니다.",
	DlgVerifyTitle: "해시 검증", DlgVerifyFile: "파일 경로：",
	DlgVerifyHash: "예상 해시값 (MD5/SHA1/SHA256)：",

	MenuSep: "----------------", ChoicePrompt: "선택：", ChoiceBad: "잘못된 선택입니다.",
	PressEnter: "Enter를 눌러 계속...", Goodbye: "안녕히 가세요.",

	ErrGenKeys: "키 생성 실패：%v", ErrImport: "가져오기 실패：%v",
	ErrEncrypt: "암호화 실패：%v", ErrDecrypt: "복호화 실패：%v",
	ErrNoKey:     "공개키가 없습니다. 먼저 키를 생성하세요.",
	ErrBadHex:    "유효하지 않은 암호문입니다.",
	ErrNoPrivKey: "개인키가 없습니다. 먼저 키를 생성하세요.",

	NotePassphrase: " [암호 보호됨]", NoteNoPass: " [경고：암호 없음]",
	GUIOnlyWindows: "GUI는 Windows만 지원합니다.", Version: "v2.0.0",
}

var localeRU = Locale{
	AppTitle: "CryphoCat",

	BtnGenKey: "Создать ключи", BtnImportKey: "Импорт ключа друга", BtnImport: "Импорт...",
	BtnImportClip: "Из буфера обмена", BtnCopyPub: "Копировать мой ключ",
	BtnEncrypt: "Шифровать", BtnDecrypt: "Расшифровать", BtnClear: "Очистить",
	BtnCopyOutput: "Копировать вывод", BtnCalcHash: "Хеш файла", BtnVerifyHash: "Проверить хеш",
	BtnLoadImage: "Загрузить изображение",
	BtnLoadFile:  "Загрузить файл",
	MenuHash:     "Хеш файла", MenuVerify: "Проверить хеш", MenuExit: "Выход",

	LabelInput: "ВВОД：", LabelOutput: "ВЫВОД：", LabelKeyInfo: "RSA / OAEP-SHA256",
	LabelAlgorithm: "Алгоритм：", LabelFriendKey: "Ключ друга：",
	LabelKeyMgmt: "[Управление ключами]", LabelTools: "[Инструменты]",
	LabelMemory: "Хранить только в памяти", LabelCompress: "Сжать вывод",

	StatusReady: "Готово.", StatusEncOK: "Зашифровано — скопировано в буфер.",
	StatusDecOK: "Расшифровано.", StatusCleared: "Очищено.",
	StatusCopied: "Вывод скопирован в буфер.", StatusKeyGen: "Ключи созданы.",
	StatusKeyImp: "Ключ друга импортирован.", StatusHashOK: "Хеш вычислен.",
	StatusHashFail: "Проверка хеша не удалась.",

	DlgPassTitle: "Установить пароль", DlgPassMsg: "Пароль для приватного ключа (пусто=без шифрования)：",
	DlgPassConfirm: "Подтвердите пароль：", DlgPassBadMsg: "Пароли не совпадают.",
	DlgDecPassTitle: "Пароль", DlgDecPassMsg: "Пароль приватного ключа (пусто=без пароля)：",
	DlgNoFriendTitle:   "Нет ключа друга",
	DlgNoFriendMsg:     "Ключ друга не найден.\nИспользовать свой публичный ключ?",
	DlgImageViewTitle:  "Просмотр изображения",
	DlgImageAskMsg:     "Результат расшифровки похож на изображение.\nПоказать?",
	DlgClipImportTitle: "Импорт из буфера",
	DlgClipImportBad:   "В буфере нет действительного публичного ключа.",
	DlgPubCopied:       "Мой публичный ключ скопирован в буфер.",

	KeyLenPrompt: "Алгоритм (1-7, по умолч. 3=RSA-4096)：",
	KeyLenBad:    "Введите число от 1 до 7.", KeyLenNotNum: "Введите число.",
	PromptPWSet:     "Установите пароль (пусто=без шифрования, НЕ рекомендуется)：",
	PromptPWConfirm: "Подтвердите пароль：", PwMismatch: "Пароли не совпадают.",
	PwWarn:   "[ПРЕДУПРЕЖДЕНИЕ] Приватный ключ будет сохранён без пароля.",
	KeyGenOK: "Пара ключей %s создана.", KeyPriv: "  Приватный：%s", KeyPub: "  Публичный：%s",

	ImportPrompt: "Путь к публичному ключу друга：", ImportOK: "Ключ импортирован.",
	ImportFail: "Ошибка импорта：%v", ImportNoFile: "Файл не найден.",

	EncNoKey: "Ключ друга не найден.", EncUseMine: "Использовать свой ключ? (y/n)：",
	EncCancelled: "Отменено.", EncNoKeyAny: "Нет публичного ключа. Сначала создайте пару ключей.",
	EncLoadFail: "Ошибка загрузки ключа：%v", EncPrompt: "Текст для шифрования：",
	EncEmpty: "Пустой ввод.", EncFail: "Ошибка шифрования：%v",
	EncOK: "Зашифровано (скопировано в буфер)：%s",

	DecNoKey:    "Приватный ключ не найден. Сначала создайте пару ключей.",
	PromptPWDec: "Пароль приватного ключа (пусто=без пароля)：",
	DecLoadFail: "Ошибка загрузки ключа：%v", DecPrompt: "Шифротекст (base64/hex)：",
	DecBadHex: "Неверный шифротекст.", DecFail: "Ошибка расшифровки：%v", DecOK: "Расшифровано：%s",

	HashPrompt: "Путь к файлу：", HashNoFile: "Файл не найден.",
	HashMD5: "MD5   ：%s", HashSHA1: "SHA1  ：%s", HashSHA256: "SHA256：%s",

	CmpPrompt: "Путь к файлу：", CmpHashIn: "Ожидаемый хеш (MD5/SHA1/SHA256)：",
	CmpNoFile: "Файл не найден.", CmpEmpty: "Хеш не указан.",
	CmpMatch: "%s СОВПАДАЕТ：%s", CmpMatchVal: "Значение：%s",
	CmpFail:        "ПРОВЕРКА ХЕША НЕ ПРОЙДЕНА — файл мог быть изменён.",
	DlgVerifyTitle: "Проверка хеша", DlgVerifyFile: "Путь к файлу：",
	DlgVerifyHash: "Ожидаемый хеш (MD5/SHA1/SHA256)：",

	MenuSep: "----------------", ChoicePrompt: "Выбор：", ChoiceBad: "Неверный выбор.",
	PressEnter: "Нажмите Enter...", Goodbye: "До свидания.",

	ErrGenKeys: "Ошибка создания ключей：%v", ErrImport: "Ошибка импорта：%v",
	ErrEncrypt: "Ошибка шифрования：%v", ErrDecrypt: "Ошибка расшифровки：%v",
	ErrNoKey:     "Нет публичного ключа. Сначала создайте пару ключей.",
	ErrBadHex:    "Неверный шифротекст.",
	ErrNoPrivKey: "Нет приватного ключа. Сначала создайте пару ключей.",

	NotePassphrase: " [защищено паролем]", NoteNoPass: " [ПРЕДУПРЕЖДЕНИЕ：без пароля]",
	GUIOnlyWindows: "GUI доступен только в Windows.", Version: "v2.0.0",
}

var localeFR = Locale{
	AppTitle: "CryphoCat",

	BtnGenKey: "Générer clés", BtnImportKey: "Importer clé ami", BtnImport: "Importer...",
	BtnImportClip: "Depuis presse-papiers", BtnCopyPub: "Copier ma clé publique",
	BtnEncrypt: "Chiffrer", BtnDecrypt: "Déchiffrer", BtnClear: "Effacer",
	BtnCopyOutput: "Copier sortie", BtnCalcHash: "Calculer hachage", BtnVerifyHash: "Vérifier hachage",
	BtnLoadImage: "Charger image",
	BtnLoadFile:  "Charger fichier",
	MenuHash:     "Hachage fichier", MenuVerify: "Vérifier hachage", MenuExit: "Quitter",

	LabelInput: "ENTRÉE：", LabelOutput: "SORTIE：", LabelKeyInfo: "RSA / OAEP-SHA256",
	LabelAlgorithm: "Algorithme：", LabelFriendKey: "Clé ami：",
	LabelKeyMgmt: "[Gestion des clés]", LabelTools: "[Outils]",
	LabelMemory: "Stocker en mémoire seulement", LabelCompress: "Compresser sortie",

	StatusReady: "Prêt.", StatusEncOK: "Chiffré — copié dans le presse-papiers.",
	StatusDecOK: "Déchiffré.", StatusCleared: "Effacé.",
	StatusCopied: "Sortie copiée.", StatusKeyGen: "Clés générées.",
	StatusKeyImp: "Clé ami importée.", StatusHashOK: "Hachage calculé.",
	StatusHashFail: "Échec vérification hachage.",

	DlgPassTitle: "Définir phrase secrète", DlgPassMsg: "Phrase secrète (vide=pas de chiffrement)：",
	DlgPassConfirm: "Confirmer：", DlgPassBadMsg: "Les phrases ne correspondent pas.",
	DlgDecPassTitle: "Phrase secrète", DlgDecPassMsg: "Phrase secrète clé privée (vide=aucune)：",
	DlgNoFriendTitle:   "Pas de clé ami",
	DlgNoFriendMsg:     "Clé publique ami introuvable.\nUtiliser votre propre clé ?",
	DlgImageViewTitle:  "Aperçu image",
	DlgImageAskMsg:     "Le résultat semble être une image.\nAfficher ?",
	DlgClipImportTitle: "Importer du presse-papiers",
	DlgClipImportBad:   "Le presse-papiers ne contient pas de clé publique valide.",
	DlgPubCopied:       "Ma clé publique copiée dans le presse-papiers.",

	KeyLenPrompt: "Algorithme (1-7, défaut 3=RSA-4096)：",
	KeyLenBad:    "Choisissez un nombre entre 1 et 7.", KeyLenNotNum: "Veuillez entrer un nombre.",
	PromptPWSet:     "Phrase secrète (vide=pas de chiffrement, DÉCONSEILLÉ)：",
	PromptPWConfirm: "Confirmer：", PwMismatch: "Les phrases ne correspondent pas.",
	PwWarn:   "[ATTENTION] La clé privée sera stockée sans phrase secrète.",
	KeyGenOK: "Paire de clés %s générée.", KeyPriv: "  Privée：%s", KeyPub: "  Publique：%s",

	ImportPrompt: "Chemin clé publique ami：", ImportOK: "Clé importée.",
	ImportFail: "Échec import：%v", ImportNoFile: "Fichier introuvable.",

	EncNoKey: "Clé ami introuvable.", EncUseMine: "Utiliser votre clé ? (o/n)：",
	EncCancelled: "Annulé.", EncNoKeyAny: "Aucune clé publique. Générez d'abord une paire.",
	EncLoadFail: "Échec chargement clé：%v", EncPrompt: "Texte à chiffrer：",
	EncEmpty: "Entrée vide.", EncFail: "Échec chiffrement：%v",
	EncOK: "Chiffré (copié)：%s",

	DecNoKey:    "Clé privée introuvable. Générez d'abord une paire.",
	PromptPWDec: "Phrase secrète (vide=aucune)：", DecLoadFail: "Échec chargement clé：%v",
	DecPrompt: "Texte chiffré (base64/hex)：", DecBadHex: "Texte chiffré invalide.",
	DecFail: "Échec déchiffrement：%v", DecOK: "Déchiffré：%s",

	HashPrompt: "Chemin fichier：", HashNoFile: "Fichier introuvable.",
	HashMD5: "MD5   ：%s", HashSHA1: "SHA1  ：%s", HashSHA256: "SHA256：%s",

	CmpPrompt: "Chemin fichier：", CmpHashIn: "Hachage attendu (MD5/SHA1/SHA256)：",
	CmpNoFile: "Fichier introuvable.", CmpEmpty: "Hachage non fourni.",
	CmpMatch: "%s CORRESPOND：%s", CmpMatchVal: "Valeur：%s",
	CmpFail:        "ÉCHEC VÉRIFICATION — le fichier a peut-être été altéré.",
	DlgVerifyTitle: "Vérifier hachage", DlgVerifyFile: "Chemin fichier：",
	DlgVerifyHash: "Hachage attendu (MD5/SHA1/SHA256)：",

	MenuSep: "----------------", ChoicePrompt: "Choix：", ChoiceBad: "Choix invalide.",
	PressEnter: "Appuyez sur Entrée...", Goodbye: "Au revoir.",

	ErrGenKeys: "Échec génération：%v", ErrImport: "Échec import：%v",
	ErrEncrypt: "Échec chiffrement：%v", ErrDecrypt: "Échec déchiffrement：%v",
	ErrNoKey:     "Aucune clé publique. Générez d'abord une paire.",
	ErrBadHex:    "Texte chiffré invalide.",
	ErrNoPrivKey: "Aucune clé privée. Générez d'abord une paire.",

	NotePassphrase: " [protégé]", NoteNoPass: " [ATTENTION：pas de phrase secrète]",
	GUIOnlyWindows: "GUI disponible seulement sous Windows.", Version: "v2.0.0",
}

var localeDE = Locale{
	AppTitle: "CryphoCat",

	BtnGenKey: "Schlüssel erstellen", BtnImportKey: "Freundes-Key importieren",
	BtnImport: "Importieren...", BtnImportClip: "Aus Zwischenablage",
	BtnCopyPub: "Eigenen Key kopieren", BtnEncrypt: "Verschlüsseln", BtnDecrypt: "Entschlüsseln",
	BtnClear: "Leeren", BtnCopyOutput: "Ausgabe kopieren",
	BtnCalcHash: "Hash berechnen", BtnVerifyHash: "Hash prüfen", BtnLoadImage: "Bild laden",
	BtnLoadFile: "Datei laden",
	MenuHash:    "Datei-Hash", MenuVerify: "Hash prüfen", MenuExit: "Beenden",

	LabelInput: "EINGABE：", LabelOutput: "AUSGABE：", LabelKeyInfo: "RSA / OAEP-SHA256",
	LabelAlgorithm: "Algorithmus：", LabelFriendKey: "Freundes-Key：",
	LabelKeyMgmt: "[Schlüsselverwaltung]", LabelTools: "[Werkzeuge]",
	LabelMemory: "Nur im Speicher", LabelCompress: "Ausgabe komprimieren",

	StatusReady: "Bereit.", StatusEncOK: "Verschlüsselt — in Zwischenablage kopiert.",
	StatusDecOK: "Entschlüsselt.", StatusCleared: "Geleert.",
	StatusCopied: "Ausgabe kopiert.", StatusKeyGen: "Schlüssel erstellt.",
	StatusKeyImp: "Freundes-Key importiert.", StatusHashOK: "Hash berechnet.",
	StatusHashFail: "Hash-Prüfung fehlgeschlagen.",

	DlgPassTitle: "Passphrase setzen", DlgPassMsg: "Passphrase für privaten Schlüssel (leer=keine Verschlüsselung)：",
	DlgPassConfirm: "Bestätigen：", DlgPassBadMsg: "Passphrasen stimmen nicht überein.",
	DlgDecPassTitle: "Passphrase", DlgDecPassMsg: "Passphrase privater Schlüssel (leer=keine)：",
	DlgNoFriendTitle:   "Kein Freundes-Key",
	DlgNoFriendMsg:     "Freundes-Key nicht gefunden.\nEigenen öffentlichen Schlüssel verwenden?",
	DlgImageViewTitle:  "Bildvorschau",
	DlgImageAskMsg:     "Ergebnis scheint ein Bild zu sein.\nAnzeigen?",
	DlgClipImportTitle: "Aus Zwischenablage",
	DlgClipImportBad:   "Zwischenablage enthält keinen gültigen öffentlichen Schlüssel.",
	DlgPubCopied:       "Eigener öffentlicher Schlüssel in Zwischenablage kopiert.",

	KeyLenPrompt: "Algorithmus (1-7, Standard 3=RSA-4096)：",
	KeyLenBad:    "Zahl zwischen 1 und 7 eingeben.", KeyLenNotNum: "Bitte Zahl eingeben.",
	PromptPWSet:     "Passphrase setzen (leer=keine Verschlüsselung, NICHT empfohlen)：",
	PromptPWConfirm: "Bestätigen：", PwMismatch: "Passphrasen stimmen nicht überein.",
	PwWarn:   "[WARNUNG] Privater Schlüssel wird ohne Passphrase gespeichert.",
	KeyGenOK: "%s Schlüsselpaar erstellt.", KeyPriv: "  Privat：%s", KeyPub: "  Öffentlich：%s",

	ImportPrompt: "Pfad zum Freundes-Key：", ImportOK: "Schlüssel importiert.",
	ImportFail: "Import fehlgeschlagen：%v", ImportNoFile: "Datei nicht gefunden.",

	EncNoKey: "Freundes-Key nicht gefunden.", EncUseMine: "Eigenen Key verwenden? (j/n)：",
	EncCancelled: "Abgebrochen.", EncNoKeyAny: "Kein öffentlicher Schlüssel. Zuerst Schlüsselpaar erstellen.",
	EncLoadFail: "Laden fehlgeschlagen：%v", EncPrompt: "Text zum Verschlüsseln：",
	EncEmpty: "Leere Eingabe.", EncFail: "Verschlüsselung fehlgeschlagen：%v",
	EncOK: "Verschlüsselt (kopiert)：%s",

	DecNoKey:    "Privater Schlüssel nicht gefunden. Zuerst Schlüsselpaar erstellen.",
	PromptPWDec: "Passphrase (leer=keine)：", DecLoadFail: "Laden fehlgeschlagen：%v",
	DecPrompt: "Chiffretext (base64/hex)：", DecBadHex: "Ungültiger Chiffretext.",
	DecFail: "Entschlüsselung fehlgeschlagen：%v", DecOK: "Entschlüsselt：%s",

	HashPrompt: "Dateipfad：", HashNoFile: "Datei nicht gefunden.",
	HashMD5: "MD5   ：%s", HashSHA1: "SHA1  ：%s", HashSHA256: "SHA256：%s",

	CmpPrompt: "Dateipfad：", CmpHashIn: "Erwarteter Hash (MD5/SHA1/SHA256)：",
	CmpNoFile: "Datei nicht gefunden.", CmpEmpty: "Kein Hash angegeben.",
	CmpMatch: "%s STIMMT ÜBEREIN：%s", CmpMatchVal: "Wert：%s",
	CmpFail:        "HASH-PRÜFUNG FEHLGESCHLAGEN — Datei könnte manipuliert sein.",
	DlgVerifyTitle: "Hash prüfen", DlgVerifyFile: "Dateipfad：",
	DlgVerifyHash: "Erwarteter Hash (MD5/SHA1/SHA256)：",

	MenuSep: "----------------", ChoicePrompt: "Auswahl：", ChoiceBad: "Ungültige Auswahl.",
	PressEnter: "Enter drücken...", Goodbye: "Auf Wiedersehen.",

	ErrGenKeys: "Schlüsselerstellung fehlgeschlagen：%v", ErrImport: "Import fehlgeschlagen：%v",
	ErrEncrypt: "Verschlüsselung fehlgeschlagen：%v", ErrDecrypt: "Entschlüsselung fehlgeschlagen：%v",
	ErrNoKey:     "Kein öffentlicher Schlüssel. Zuerst Schlüsselpaar erstellen.",
	ErrBadHex:    "Ungültiger Chiffretext.",
	ErrNoPrivKey: "Kein privater Schlüssel. Zuerst Schlüsselpaar erstellen.",

	NotePassphrase: " [passwortgeschützt]", NoteNoPass: " [WARNUNG：keine Passphrase]",
	GUIOnlyWindows: "GUI nur unter Windows verfügbar.", Version: "v2.0.0",
}

var localeES = Locale{
	AppTitle: "CryphoCat",

	BtnGenKey: "Generar claves", BtnImportKey: "Importar clave amigo", BtnImport: "Importar...",
	BtnImportClip: "Desde portapapeles", BtnCopyPub: "Copiar mi clave pública",
	BtnEncrypt: "Cifrar", BtnDecrypt: "Descifrar", BtnClear: "Limpiar",
	BtnCopyOutput: "Copiar salida", BtnCalcHash: "Calcular hash", BtnVerifyHash: "Verificar hash",
	BtnLoadImage: "Cargar imagen",
	BtnLoadFile:  "Cargar archivo",
	MenuHash:     "Hash archivo", MenuVerify: "Verificar hash", MenuExit: "Salir",

	LabelInput: "ENTRADA：", LabelOutput: "SALIDA：", LabelKeyInfo: "RSA / OAEP-SHA256",
	LabelAlgorithm: "Algoritmo：", LabelFriendKey: "Clave amigo：",
	LabelKeyMgmt: "[Gestión claves]", LabelTools: "[Herramientas]",
	LabelMemory: "Solo en memoria", LabelCompress: "Comprimir salida",

	StatusReady: "Listo.", StatusEncOK: "Cifrado — copiado al portapapeles.",
	StatusDecOK: "Descifrado.", StatusCleared: "Limpiado.",
	StatusCopied: "Salida copiada.", StatusKeyGen: "Claves generadas.",
	StatusKeyImp: "Clave amigo importada.", StatusHashOK: "Hash calculado.",
	StatusHashFail: "Verificación hash fallida.",

	DlgPassTitle: "Establecer contraseña", DlgPassMsg: "Contraseña clave privada (vacío=sin cifrar)：",
	DlgPassConfirm: "Confirmar：", DlgPassBadMsg: "Las contraseñas no coinciden.",
	DlgDecPassTitle: "Contraseña", DlgDecPassMsg: "Contraseña clave privada (vacío=ninguna)：",
	DlgNoFriendTitle:   "Sin clave amigo",
	DlgNoFriendMsg:     "Clave pública del amigo no encontrada.\n¿Usar tu propia clave pública?",
	DlgImageViewTitle:  "Vista previa",
	DlgImageAskMsg:     "El resultado parece una imagen.\n¿Mostrar?",
	DlgClipImportTitle: "Importar del portapapeles",
	DlgClipImportBad:   "El portapapeles no contiene una clave pública válida.",
	DlgPubCopied:       "Mi clave pública copiada al portapapeles.",

	KeyLenPrompt: "Algoritmo (1-7, predet. 3=RSA-4096)：",
	KeyLenBad:    "Elija un número del 1 al 7.", KeyLenNotNum: "Introduzca un número.",
	PromptPWSet:     "Contraseña (vacío=sin cifrar, NO recomendado)：",
	PromptPWConfirm: "Confirmar：", PwMismatch: "Las contraseñas no coinciden.",
	PwWarn:   "[AVISO] La clave privada se guardará sin contraseña.",
	KeyGenOK: "Par de claves %s generado.", KeyPriv: "  Privada：%s", KeyPub: "  Pública：%s",

	ImportPrompt: "Ruta clave pública amigo：", ImportOK: "Clave importada.",
	ImportFail: "Error importación：%v", ImportNoFile: "Archivo no encontrado.",

	EncNoKey: "Clave amigo no encontrada.", EncUseMine: "¿Usar tu clave? (s/n)：",
	EncCancelled: "Cancelado.", EncNoKeyAny: "Sin clave pública. Genere un par primero.",
	EncLoadFail: "Error al cargar clave：%v", EncPrompt: "Texto a cifrar：",
	EncEmpty: "Entrada vacía.", EncFail: "Error cifrado：%v",
	EncOK: "Cifrado (copiado)：%s",

	DecNoKey:    "Clave privada no encontrada. Genere un par primero.",
	PromptPWDec: "Contraseña (vacío=ninguna)：", DecLoadFail: "Error al cargar clave：%v",
	DecPrompt: "Texto cifrado (base64/hex)：", DecBadHex: "Texto cifrado no válido.",
	DecFail: "Error descifrado：%v", DecOK: "Descifrado：%s",

	HashPrompt: "Ruta archivo：", HashNoFile: "Archivo no encontrado.",
	HashMD5: "MD5   ：%s", HashSHA1: "SHA1  ：%s", HashSHA256: "SHA256：%s",

	CmpPrompt: "Ruta archivo：", CmpHashIn: "Hash esperado (MD5/SHA1/SHA256)：",
	CmpNoFile: "Archivo no encontrado.", CmpEmpty: "Hash no proporcionado.",
	CmpMatch: "%s COINCIDE：%s", CmpMatchVal: "Valor：%s",
	CmpFail:        "VERIFICACIÓN FALLIDA — el archivo puede haber sido alterado.",
	DlgVerifyTitle: "Verificar hash", DlgVerifyFile: "Ruta archivo：",
	DlgVerifyHash: "Hash esperado (MD5/SHA1/SHA256)：",

	MenuSep: "----------------", ChoicePrompt: "Opción：", ChoiceBad: "Opción no válida.",
	PressEnter: "Presione Enter...", Goodbye: "Adiós.",

	ErrGenKeys: "Error generación：%v", ErrImport: "Error importación：%v",
	ErrEncrypt: "Error cifrado：%v", ErrDecrypt: "Error descifrado：%v",
	ErrNoKey:     "Sin clave pública. Genere un par primero.",
	ErrBadHex:    "Texto cifrado no válido.",
	ErrNoPrivKey: "Sin clave privada. Genere un par primero.",

	NotePassphrase: " [protegido]", NoteNoPass: " [AVISO：sin contraseña]",
	GUIOnlyWindows: "GUI solo disponible en Windows.", Version: "v2.0.0",
}

var localeZHTW = Locale{
	AppTitle: "CryphoCat",

	BtnGenKey: "產生金鑰", BtnImportKey: "匯入對方公鑰", BtnImport: "匯入...",
	BtnImportClip: "從剪貼簿匯入", BtnCopyPub: "複製我的公鑰",
	BtnEncrypt: "加密", BtnDecrypt: "解密", BtnClear: "清除", BtnCopyOutput: "複製輸出",
	BtnCalcHash: "計算雜湊", BtnVerifyHash: "驗證雜湊", BtnLoadImage: "載入圖片",
	BtnLoadFile: "從檔案載入",
	MenuHash:    "計算檔案雜湊", MenuVerify: "驗證檔案雜湊", MenuExit: "離開",

	LabelInput: "輸入：", LabelOutput: "輸出：", LabelKeyInfo: "RSA / OAEP-SHA256",
	LabelAlgorithm: "演算法：", LabelFriendKey: "對方公鑰：",
	LabelKeyMgmt: "[金鑰管理]", LabelTools: "[工具]",
	LabelMemory: "僅儲存在記憶體", LabelCompress: "壓縮輸出",

	StatusReady: "就緒。", StatusEncOK: "加密完成 - 已複製到剪貼簿。",
	StatusDecOK: "解密完成。", StatusCleared: "已清除。", StatusCopied: "輸出已複製到剪貼簿。",
	StatusKeyGen: "金鑰已產生。", StatusKeyImp: "對方公鑰已匯入。",
	StatusHashOK: "雜湊計算完成。", StatusHashFail: "雜湊驗證失敗。",

	DlgPassTitle: "設定密碼", DlgPassMsg: "私鑰保護密碼（留空 = 不加密儲存，不建議）：",
	DlgPassConfirm: "確認密碼：", DlgPassBadMsg: "兩次輸入的密碼不一致。",
	DlgDecPassTitle: "輸入密碼", DlgDecPassMsg: "私鑰密碼（留空表示無密碼）：",
	DlgNoFriendTitle:  "未找到對方公鑰",
	DlgNoFriendMsg:    "未找到對方公鑰。\n是否使用自己的公鑰加密？",
	DlgImageViewTitle: "圖片預覽", DlgImageAskMsg: "解密結果似乎是圖片。\n是否檢視？",
	DlgClipImportTitle: "從剪貼簿匯入", DlgClipImportBad: "剪貼簿中無有效的公鑰。",
	DlgPubCopied: "我的公鑰已複製到剪貼簿。",

	KeyLenPrompt: "演算法 (1-7，預設 3=RSA-4096)：",
	KeyLenBad:    "請輸入 1-7 之間的數字。", KeyLenNotNum: "請輸入數字。",
	PromptPWSet:     "設定私鑰密碼（留空 = 不加密儲存，不建議）：",
	PromptPWConfirm: "確認密碼：", PwMismatch: "兩次輸入的密碼不一致。",
	PwWarn:   "[警告] 私鑰將以明文形式儲存，安全性較低。",
	KeyGenOK: "%s 金鑰對已產生。", KeyPriv: "  私鑰：%s", KeyPub: "  公鑰：%s",

	ImportPrompt: "請輸入對方公鑰檔案路徑：", ImportOK: "公鑰匯入成功。",
	ImportFail: "匯入失敗：%v", ImportNoFile: "檔案不存在，請檢查路徑。",

	EncNoKey: "未找到對方公鑰。", EncUseMine: "是否使用自己的公鑰加密？(y/n)：",
	EncCancelled: "已取消。", EncNoKeyAny: "未找到公鑰，請先產生金鑰對。",
	EncLoadFail: "無法載入公鑰：%v", EncPrompt: "請輸入要加密的文字：",
	EncEmpty: "輸入為空。", EncFail: "加密失敗：%v",
	EncOK: "加密結果（已複製到剪貼簿）：%s",

	DecNoKey:    "未找到私鑰，請先產生金鑰對。",
	PromptPWDec: "私鑰密碼（無密碼則留空）：", DecLoadFail: "無法載入私鑰：%v",
	DecPrompt: "密文（base64/hex）：", DecBadHex: "密文無效。",
	DecFail: "解密失敗：%v", DecOK: "解密結果：%s",

	HashPrompt: "請輸入檔案路徑：", HashNoFile: "檔案不存在。",
	HashMD5: "MD5   ：%s", HashSHA1: "SHA1  ：%s", HashSHA256: "SHA256：%s",

	CmpPrompt: "請輸入檔案路徑：", CmpHashIn: "請輸入雜湊值（MD5/SHA1/SHA256）：",
	CmpNoFile: "檔案不存在。", CmpEmpty: "未輸入雜湊值。",
	CmpMatch: "%s 驗證通過：%s", CmpMatchVal: "雜湊值：%s",
	CmpFail:        "雜湊驗證失敗 —— 檔案可能已被竄改。",
	DlgVerifyTitle: "驗證檔案雜湊", DlgVerifyFile: "檔案路徑：",
	DlgVerifyHash: "期望雜湊值（MD5/SHA1/SHA256）：",

	MenuSep: "----------------", ChoicePrompt: "請輸入選項：", ChoiceBad: "選項無效。",
	PressEnter: "按下 Enter 繼續...", Goodbye: "再見。",

	ErrGenKeys: "金鑰產生失敗：%v", ErrImport: "匯入失敗：%v",
	ErrEncrypt: "加密失敗：%v", ErrDecrypt: "解密失敗：%v",
	ErrNoKey: "未找到公鑰，請先產生金鑰對。", ErrBadHex: "密文無效。",
	ErrNoPrivKey: "未找到私鑰，請先產生金鑰對。",

	NotePassphrase: " [已設定密碼保護]", NoteNoPass: " [警告：未設定密碼]",
	GUIOnlyWindows: "GUI 僅支援 Windows 系統。", Version: "v2.0.0",
}

var localeZHHK = Locale{
	AppTitle: "CryphoCat",

	BtnGenKey: "產生金鑰", BtnImportKey: "匯入對方公鑰", BtnImport: "匯入...",
	BtnImportClip: "從剪貼簿匯入", BtnCopyPub: "複製我的公鑰",
	BtnEncrypt: "加密", BtnDecrypt: "解密", BtnClear: "清除", BtnCopyOutput: "複製輸出",
	BtnCalcHash: "計算雜湊", BtnVerifyHash: "驗證雜湊", BtnLoadImage: "載入圖片",
	BtnLoadFile: "從檔案載入",
	MenuHash:    "計算檔案雜湊", MenuVerify: "驗證檔案雜湊", MenuExit: "離開",

	LabelInput: "輸入：", LabelOutput: "輸出：", LabelKeyInfo: "RSA / OAEP-SHA256",
	LabelAlgorithm: "演算法：", LabelFriendKey: "對方公鑰：",
	LabelKeyMgmt: "[金鑰管理]", LabelTools: "[工具]",
	LabelMemory: "僅儲存在記憶體", LabelCompress: "壓縮輸出",

	StatusReady: "就緒。", StatusEncOK: "加密完成 - 已複製到剪貼簿。",
	StatusDecOK: "解密完成。", StatusCleared: "已清除。", StatusCopied: "輸出已複製到剪貼簿。",
	StatusKeyGen: "金鑰已產生。", StatusKeyImp: "對方公鑰已匯入。",
	StatusHashOK: "雜湊計算完成。", StatusHashFail: "雜湊驗證失敗。",

	DlgPassTitle: "設定密碼", DlgPassMsg: "私鑰保護密碼（留空 = 不加密儲存，不建議）：",
	DlgPassConfirm: "確認密碼：", DlgPassBadMsg: "兩次輸入的密碼不一致。",
	DlgDecPassTitle: "輸入密碼", DlgDecPassMsg: "私鑰密碼（留空表示無密碼）：",
	DlgNoFriendTitle:  "未找到對方公鑰",
	DlgNoFriendMsg:    "未找到對方公鑰。\n是否使用自己的公鑰加密？",
	DlgImageViewTitle: "圖片預覽", DlgImageAskMsg: "解密結果似乎是圖片。\n是否檢視？",
	DlgClipImportTitle: "從剪貼簿匯入", DlgClipImportBad: "剪貼簿中無有效的公鑰。",
	DlgPubCopied: "我的公鑰已複製到剪貼簿。",

	KeyLenPrompt: "演算法 (1-7，預設 3=RSA-4096)：",
	KeyLenBad:    "請輸入 1-7 之間的數字。", KeyLenNotNum: "請輸入數字。",
	PromptPWSet:     "設定私鑰密碼（留空 = 不加密儲存，不建議）：",
	PromptPWConfirm: "確認密碼：", PwMismatch: "兩次輸入的密碼不一致。",
	PwWarn:   "[警告] 私鑰將以明文形式儲存，安全性較低。",
	KeyGenOK: "%s 金鑰對已產生。", KeyPriv: "  私鑰：%s", KeyPub: "  公鑰：%s",

	ImportPrompt: "請輸入對方公鑰檔案路徑：", ImportOK: "公鑰匯入成功。",
	ImportFail: "匯入失敗：%v", ImportNoFile: "檔案不存在，請檢查路徑。",

	EncNoKey: "未找到對方公鑰。", EncUseMine: "是否使用自己的公鑰加密？(y/n)：",
	EncCancelled: "已取消。", EncNoKeyAny: "未找到公鑰，請先產生金鑰對。",
	EncLoadFail: "無法載入公鑰：%v", EncPrompt: "請輸入要加密的文字：",
	EncEmpty: "輸入為空。", EncFail: "加密失敗：%v",
	EncOK: "加密結果（已複製到剪貼簿）：%s",

	DecNoKey:    "未找到私鑰，請先產生金鑰對。",
	PromptPWDec: "私鑰密碼（無密碼則留空）：", DecLoadFail: "無法載入私鑰：%v",
	DecPrompt: "密文（base64/hex）：", DecBadHex: "密文無效。",
	DecFail: "解密失敗：%v", DecOK: "解密結果：%s",

	HashPrompt: "請輸入檔案路徑：", HashNoFile: "檔案不存在。",
	HashMD5: "MD5   ：%s", HashSHA1: "SHA1  ：%s", HashSHA256: "SHA256：%s",

	CmpPrompt: "請輸入檔案路徑：", CmpHashIn: "請輸入雜湊值（MD5/SHA1/SHA256）：",
	CmpNoFile: "檔案不存在。", CmpEmpty: "未輸入雜湊值。",
	CmpMatch: "%s 驗證通過：%s", CmpMatchVal: "雜湊值：%s",
	CmpFail:        "雜湊驗證失敗 —— 檔案可能已被竄改。",
	DlgVerifyTitle: "驗證檔案雜湊", DlgVerifyFile: "檔案路徑：",
	DlgVerifyHash: "期望雜湊值（MD5/SHA1/SHA256）：",

	MenuSep: "----------------", ChoicePrompt: "請輸入選項：", ChoiceBad: "選項無效。",
	PressEnter: "按下 Enter 繼續...", Goodbye: "再見。",

	ErrGenKeys: "金鑰產生失敗：%v", ErrImport: "匯入失敗：%v",
	ErrEncrypt: "加密失敗：%v", ErrDecrypt: "解密失敗：%v",
	ErrNoKey: "未找到公鑰，請先產生金鑰對。", ErrBadHex: "密文無效。",
	ErrNoPrivKey: "未找到私鑰，請先產生金鑰對。",

	NotePassphrase: " [已設定密碼保護]", NoteNoPass: " [警告：未設定密碼]",
	GUIOnlyWindows: "GUI 僅支援 Windows 系統。", Version: "v2.0.0",
}

// SelectLocale picks the locale based on the -lang flag or system detection.
func SelectLocale(flagLang string) *Locale {
	lang := flagLang
	if lang == "" {
		lang = detectSystemLang()
	}
	switch {
	case lang == "zhcn" || lang == "zh" || strings.HasPrefix(lang, "zh-CN"):
		return &localeZHCN
	case lang == "zhtw" || strings.HasPrefix(lang, "zh-TW"):
		return &localeZHTW
	case lang == "zhhk" || strings.HasPrefix(lang, "zh-HK"):
		return &localeZHHK
	case lang == "ja" || strings.HasPrefix(lang, "ja"):
		return &localeJA
	case lang == "ko" || strings.HasPrefix(lang, "ko"):
		return &localeKO
	case lang == "ru" || strings.HasPrefix(lang, "ru"):
		return &localeRU
	case lang == "fr" || strings.HasPrefix(lang, "fr"):
		return &localeFR
	case lang == "de" || strings.HasPrefix(lang, "de"):
		return &localeDE
	case lang == "es" || strings.HasPrefix(lang, "es"):
		return &localeES
	}
	return &localeEN
}

// detectSystemLang tries to detect the OS language.
// On Windows it checks GetUserDefaultUILanguage; elsewhere it checks LANG env.
func detectSystemLang() string {
	if lang := os.Getenv("LANG"); lang != "" {
		return lang
	}
	if lang := os.Getenv("LC_ALL"); lang != "" {
		return lang
	}
	// Windows-specific detection is in a platform file.
	return detectWindowsLang()
}
