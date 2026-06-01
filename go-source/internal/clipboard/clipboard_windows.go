//go:build windows

// Package clipboard provides platform-specific clipboard operations.
package clipboard

import (
	"os/exec"
	"strings"
)

// Set copies text to the system clipboard.
func Set(text string) error {
	cmd := exec.Command("clip")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// Get returns the current text content of the system clipboard.
func Get() (string, error) {
	out, err := exec.Command("powershell", "-NoProfile", "-Command", "Get-Clipboard").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
