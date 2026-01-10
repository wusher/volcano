package styles

import (
	"strings"
	"testing"
)

func TestMinifyCSS(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes comments",
			input:    "/* comment */ body { color: red; }",
			expected: "body{color:red}",
		},
		{
			name:     "removes whitespace",
			input:    "body  {  color:  red;  }",
			expected: "body{color:red}",
		},
		{
			name:     "removes newlines",
			input:    "body {\n  color: red;\n}",
			expected: "body{color:red}",
		},
		{
			name:     "removes trailing semicolon",
			input:    "body { color: red; }",
			expected: "body{color:red}",
		},
		{
			name:     "handles multiple rules",
			input:    "body { color: red; } .foo { margin: 0; }",
			expected: "body{color:red}.foo{margin:0}",
		},
		{
			name:     "preserves important values",
			input:    ".class { display: block !important; }",
			expected: ".class{display:block!important}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MinifyCSS(tt.input)
			if err != nil {
				t.Fatalf("MinifyCSS() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("MinifyCSS() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestMinifyCSS_Themes(t *testing.T) {
	themes := []string{"docs", "blog", "vanilla"}

	for _, theme := range themes {
		t.Run(theme, func(t *testing.T) {
			original := GetCSS(theme)
			minified, err := MinifyCSS(original)
			if err != nil {
				t.Fatalf("MinifyCSS() error = %v", err)
			}

			// Minified should be smaller
			if len(minified) >= len(original) {
				t.Errorf("Minified CSS should be smaller: original=%d, minified=%d",
					len(original), len(minified))
			}

			// Minified should not have excessive whitespace
			if strings.Contains(minified, "  ") {
				t.Error("Minified CSS should not contain double spaces")
			}

			// Minified should still contain essential selectors
			essentials := []string{".sidebar", ".content", ":root"}
			for _, sel := range essentials {
				if !strings.Contains(minified, sel) {
					t.Errorf("Minified CSS should contain %q", sel)
				}
			}
		})
	}
}

func TestMinifyCSS_Empty(t *testing.T) {
	result, err := MinifyCSS("")
	if err != nil {
		t.Fatalf("MinifyCSS() error = %v", err)
	}
	if result != "" {
		t.Errorf("MinifyCSS(\"\") = %q, want empty string", result)
	}
}
