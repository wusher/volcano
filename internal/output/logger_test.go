package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("Println", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, false, false, false)
		logger.Println("Hello %s", "World")

		if !strings.Contains(buf.String(), "Hello World") {
			t.Errorf("Expected 'Hello World', got %q", buf.String())
		}
	})

	t.Run("Println quiet mode", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, false, true, false)
		logger.Println("Hello")

		if buf.String() != "" {
			t.Errorf("Expected empty output in quiet mode, got %q", buf.String())
		}
	})

	t.Run("Verbose only in verbose mode", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, false, false, false)
		logger.Verbose("Debug info")

		if buf.String() != "" {
			t.Errorf("Expected no output when not verbose, got %q", buf.String())
		}

		buf.Reset()
		logger = NewLogger(&buf, false, false, true)
		logger.Verbose("Debug info")

		if !strings.Contains(buf.String(), "Debug info") {
			t.Errorf("Expected 'Debug info' in verbose mode, got %q", buf.String())
		}
	})

	t.Run("Success message", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, false, false, false)
		logger.Success("Task completed")

		if !strings.Contains(buf.String(), "Task completed") {
			t.Errorf("Expected 'Task completed', got %q", buf.String())
		}
	})

	t.Run("Success message colored", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, true, false, false)
		logger.Success("Task completed")

		if !strings.Contains(buf.String(), ColorGreen) {
			t.Errorf("Expected green color code, got %q", buf.String())
		}
	})

	t.Run("Warning message", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, false, false, false)
		logger.Warning("Something might be wrong")

		if !strings.Contains(buf.String(), "Warning:") {
			t.Errorf("Expected 'Warning:', got %q", buf.String())
		}
	})

	t.Run("Warning message colored", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, true, false, false)
		logger.Warning("Something might be wrong")

		if !strings.Contains(buf.String(), ColorYellow) {
			t.Errorf("Expected yellow color code, got %q", buf.String())
		}
	})

	t.Run("Error message", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, false, false, false)
		logger.Error("Something went wrong")

		if !strings.Contains(buf.String(), "Error:") {
			t.Errorf("Expected 'Error:', got %q", buf.String())
		}
	})

	t.Run("Error message colored", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, true, false, false)
		logger.Error("Something went wrong")

		if !strings.Contains(buf.String(), ColorRed) {
			t.Errorf("Expected red color code, got %q", buf.String())
		}
	})

	t.Run("FileSuccess", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, false, false, false)
		logger.FileSuccess("index.md")

		output := buf.String()
		if !strings.Contains(output, "✓") {
			t.Errorf("Expected checkmark, got %q", output)
		}
		if !strings.Contains(output, "index.md") {
			t.Errorf("Expected 'index.md', got %q", output)
		}
	})

	t.Run("FileSuccess colored", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, true, false, false)
		logger.FileSuccess("index.md")

		if !strings.Contains(buf.String(), ColorGreen) {
			t.Errorf("Expected green color code, got %q", buf.String())
		}
	})

	t.Run("FileError", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, false, false, false)
		logger.FileError("broken.md", errors.New("parse error"))

		output := buf.String()
		if !strings.Contains(output, "✗") {
			t.Errorf("Expected X mark, got %q", output)
		}
		if !strings.Contains(output, "broken.md") {
			t.Errorf("Expected 'broken.md', got %q", output)
		}
		if !strings.Contains(output, "parse error") {
			t.Errorf("Expected 'parse error', got %q", output)
		}
	})

	t.Run("FileError colored", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, true, false, false)
		logger.FileError("broken.md", errors.New("parse error"))

		if !strings.Contains(buf.String(), ColorRed) {
			t.Errorf("Expected red color code, got %q", buf.String())
		}
	})

	t.Run("Verbose colored", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, true, false, true)
		logger.Verbose("Debug info")

		if !strings.Contains(buf.String(), ColorGray) {
			t.Errorf("Expected gray color code, got %q", buf.String())
		}
	})

	t.Run("Print", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, false, false, false)
		logger.Print("No newline")

		if buf.String() != "No newline" {
			t.Errorf("Expected 'No newline', got %q", buf.String())
		}
	})

	t.Run("Print quiet mode", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger(&buf, false, true, false)
		logger.Print("No output")

		if buf.String() != "" {
			t.Errorf("Expected empty output in quiet mode, got %q", buf.String())
		}
	})
}

func TestNewLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, true, true, true)

	if logger == nil {
		t.Fatal("NewLogger should not return nil")
	}
	if logger.writer != &buf {
		t.Error("Writer should be set correctly")
	}
	if !logger.colored {
		t.Error("Colored should be true")
	}
	if !logger.quiet {
		t.Error("Quiet should be true")
	}
	if !logger.verbose {
		t.Error("Verbose should be true")
	}
}
