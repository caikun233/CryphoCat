//go:build windows

package main

import (
	"os/exec"
	"strings"
)

func trySetClipboard(text string) error {
	cmd := exec.Command("clip") //nolint:gosec
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
