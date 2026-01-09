package output

import (
	"testing"
)

func TestColorConstants(t *testing.T) {
	// Verify color constants are not empty
	colors := []struct {
		name  string
		value string
	}{
		{"ColorReset", ColorReset},
		{"ColorRed", ColorRed},
		{"ColorGreen", ColorGreen},
		{"ColorYellow", ColorYellow},
		{"ColorBlue", ColorBlue},
		{"ColorGray", ColorGray},
	}

	for _, c := range colors {
		if c.value == "" {
			t.Errorf("%s should not be empty", c.name)
		}
		// All ANSI codes should start with escape character
		if c.value[0] != '\033' {
			t.Errorf("%s should start with escape character", c.name)
		}
	}
}

func TestIsTTY(t *testing.T) {
	// Test with invalid file descriptor - should return false
	result := IsTTY(-1)
	if result {
		t.Error("IsTTY should return false for invalid fd")
	}
}

func TestIsStdoutTTY(_ *testing.T) {
	// Just verify the function doesn't panic
	// The result depends on whether we're running in a terminal
	_ = IsStdoutTTY()
}

func TestIsStderrTTY(_ *testing.T) {
	// Just verify the function doesn't panic
	// The result depends on whether we're running in a terminal
	_ = IsStderrTTY()
}
