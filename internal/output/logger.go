package output

import (
	"fmt"
	"io"
)

// Logger provides colored output with quiet/verbose modes
type Logger struct {
	writer  io.Writer
	colored bool
	quiet   bool
	verbose bool
}

// NewLogger creates a new Logger
func NewLogger(w io.Writer, colored, quiet, verbose bool) *Logger {
	return &Logger{
		writer:  w,
		colored: colored,
		quiet:   quiet,
		verbose: verbose,
	}
}

// Print prints a message (suppressed in quiet mode)
func (l *Logger) Print(format string, args ...interface{}) {
	if l.quiet {
		return
	}
	_, _ = fmt.Fprintf(l.writer, format, args...)
}

// Println prints a message with newline (suppressed in quiet mode)
func (l *Logger) Println(format string, args ...interface{}) {
	if l.quiet {
		return
	}
	_, _ = fmt.Fprintf(l.writer, format+"\n", args...)
}

// Verbose prints a message only in verbose mode
func (l *Logger) Verbose(format string, args ...interface{}) {
	if !l.verbose {
		return
	}
	if l.colored {
		_, _ = fmt.Fprintf(l.writer, ColorGray+format+ColorReset+"\n", args...)
	} else {
		_, _ = fmt.Fprintf(l.writer, format+"\n", args...)
	}
}

// Success prints a success message in green
func (l *Logger) Success(format string, args ...interface{}) {
	if l.quiet {
		return
	}
	if l.colored {
		_, _ = fmt.Fprintf(l.writer, ColorGreen+format+ColorReset+"\n", args...)
	} else {
		_, _ = fmt.Fprintf(l.writer, format+"\n", args...)
	}
}

// Warning prints a warning message in yellow
func (l *Logger) Warning(format string, args ...interface{}) {
	if l.quiet {
		return
	}
	if l.colored {
		_, _ = fmt.Fprintf(l.writer, ColorYellow+"Warning: "+format+ColorReset+"\n", args...)
	} else {
		_, _ = fmt.Fprintf(l.writer, "Warning: "+format+"\n", args...)
	}
}

// Error prints an error message in red
func (l *Logger) Error(format string, args ...interface{}) {
	if l.colored {
		_, _ = fmt.Fprintf(l.writer, ColorRed+"Error: "+format+ColorReset+"\n", args...)
	} else {
		_, _ = fmt.Fprintf(l.writer, "Error: "+format+"\n", args...)
	}
}

// FileSuccess prints a file success message with checkmark
func (l *Logger) FileSuccess(path string) {
	if l.quiet {
		return
	}
	if l.colored {
		_, _ = fmt.Fprintf(l.writer, "  "+ColorGreen+"✓"+ColorReset+" %s\n", path)
	} else {
		_, _ = fmt.Fprintf(l.writer, "  ✓ %s\n", path)
	}
}

// FileError prints a file error message with X
func (l *Logger) FileError(path string, err error) {
	if l.colored {
		_, _ = fmt.Fprintf(l.writer, "  "+ColorRed+"✗"+ColorReset+" %s: %v\n", path, err)
	} else {
		_, _ = fmt.Fprintf(l.writer, "  ✗ %s: %v\n", path, err)
	}
}
