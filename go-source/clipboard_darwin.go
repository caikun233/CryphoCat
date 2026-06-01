//go:build darwin

package main

import (
	"os/exec"
	"strings"
)

func trySetClipboard(text string) error {
	cmd := exec.Command("pbcopy") //nolint:gosec
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
