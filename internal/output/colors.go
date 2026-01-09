// Package output provides colored terminal output utilities.
package output

import (
	"os"

	"golang.org/x/term"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorGray   = "\033[90m"
)

// IsTTY checks if the given file descriptor is a terminal
func IsTTY(fd int) bool {
	return term.IsTerminal(fd)
}

// IsStdoutTTY checks if stdout is a terminal
func IsStdoutTTY() bool {
	return IsTTY(int(os.Stdout.Fd()))
}

// IsStderrTTY checks if stderr is a terminal
func IsStderrTTY() bool {
	return IsTTY(int(os.Stderr.Fd()))
}
