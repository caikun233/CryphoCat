//go:build windows

package locale

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procGetUserDefaultUILanguage = kernel32.NewProc("GetUserDefaultUILanguage")
)

// detectWindowsLang detects the Windows UI language via GetUserDefaultUILanguage.
func detectWindowsLang() string {
	langID, _, _ := procGetUserDefaultUILanguage.Call()
	if langID == 0 {
		return "en"
	}
	switch uint16(langID) {
	case 0x0004, 0x0804: // zh-CHS, zh-CN
		return "zhcn"
	case 0x0404: // zh-TW
		return "zhtw"
	case 0x0C04: // zh-HK
		return "zhhk"
	case 0x0011: // ja
		return "ja"
	case 0x0012: // ko
		return "ko"
	case 0x0019: // ru
		return "ru"
	case 0x000C: // fr
		return "fr"
	case 0x0007: // de
		return "de"
	case 0x000A: // es
		return "es"
	}
	return "en"
}

// ensure syscall and unsafe are used (avoid "imported and not used" error).
var _ = unsafe.Sizeof(0)
