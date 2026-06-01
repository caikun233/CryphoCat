//go:build !windows

package locale

// detectWindowsLang is a no-op on non-Windows; locale is driven by LANG env.
func detectWindowsLang() string {
	return "en"
}
