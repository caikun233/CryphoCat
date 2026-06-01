//go:build !windows

// Stub for non-Windows platforms.
// The GUI is Windows-only; on other platforms we always run CLI mode.

package main

import "github.com/caikun233/cryphocat/internal/locale"

func autoGUI() bool {
	return false
}

func runGUI(loc *locale.Locale) {
	runCLI(loc) // fallback
}
