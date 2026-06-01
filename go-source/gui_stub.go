//go:build cli

package main

import "fmt"

// runGUI is a no-op stub when built with the cli tag (no Fyne dependency).
func runGUI(_ string) {
	fmt.Println("GUI is not available in this CLI-only build.")
}
