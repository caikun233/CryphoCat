// CryphoCat – compact RSA encryption tool.
// One binary, two modes: CLI when launched from a terminal,
// Windows-native GUI when double-clicked on Windows.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/caikun233/cryphocat/internal/clipboard"
	"github.com/caikun233/cryphocat/internal/locale"
)

func main() {
	// ---- flags ----
	guiFlag := flag.Bool("gui", false, "Force GUI mode")
	cliFlag := flag.Bool("cli", false, "Force CLI mode")
	langFlag := flag.String("lang", "", "Interface language: en, zhcn (auto-detect if empty)")
	flag.Parse()

	loc := locale.SelectLocale(*langFlag)

	// ---- mode selection ----
	switch {
	case *guiFlag:
		runGUI(loc)
	case *cliFlag:
		runCLI(loc)
	case len(os.Args) > 1:
		// User passed arguments – assume CLI.
		runCLI(loc)
	default:
		// No arguments: auto-detect.  On Windows → GUI; other → CLI.
		if autoGUI() {
			runGUI(loc)
		} else {
			runCLI(loc)
		}
	}
}

// ---- small helpers shared across modes ----

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

func setClipboard(text string) {
	_ = clipboard.Set(text)
}

func showError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
