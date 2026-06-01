//go:build !windows && !darwin

package main

import (
	"os/exec"
	"strings"
)

func trySetClipboard(text string) error {
	// Try xclip, then xsel as fallbacks on Linux.
	for _, args := range [][]string{
		{"xclip", "-selection", "clipboard"},
		{"xsel", "--clipboard", "--input"},
	} {
		cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	return nil
}
