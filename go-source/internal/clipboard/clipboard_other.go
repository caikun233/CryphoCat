//go:build !windows

package clipboard

// Set copies text to the system clipboard (no-op on this platform).
func Set(_ string) error { return nil }

// Get returns the current clipboard text (no-op on this platform).
func Get() (string, error) { return "", nil }
