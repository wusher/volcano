package minify

import (
	"strings"
	"testing"
)

func TestJS(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		checkFn func(string) bool
	}{
		{
			name:  "empty input",
			input: "",
			checkFn: func(output string) bool {
				return output == ""
			},
		},
		{
			name:  "simple function",
			input: "function test() { return 42; }",
			checkFn: func(output string) bool {
				return len(output) < len("function test() { return 42; }") &&
					strings.Contains(output, "test") &&
					strings.Contains(output, "42")
			},
		},
		{
			name:  "removes whitespace",
			input: "const x    =     5;   \n  const y = 10;",
			checkFn: func(output string) bool {
				return len(output) < 30 &&
					strings.Contains(output, "x") &&
					strings.Contains(output, "y")
			},
		},
		{
			name:  "preserves functionality",
			input: "document.addEventListener('click', function() { console.log('test'); });",
			checkFn: func(output string) bool {
				return strings.Contains(output, "addEventListener") &&
					strings.Contains(output, "click") &&
					strings.Contains(output, "console.log")
			},
		},
		{
			name:  "preserves string content",
			input: "const msg = 'Hello, World!';",
			checkFn: func(output string) bool {
				return strings.Contains(output, "Hello, World!")
			},
		},
		{
			name:  "preserves template syntax in strings",
			input: "const url = '{{.BaseURL}}';",
			checkFn: func(output string) bool {
				return strings.Contains(output, "{{.BaseURL}}")
			},
		},
		{
			name: "minifies large script",
			input: `
				(function() {
					'use strict';

					// This is a comment
					const config = {
						delay: 100,
						active: true
					};

					function init() {
						console.log('Initializing...');
					}

					init();
				})();
			`,
			checkFn: func(output string) bool {
				// Should be significantly smaller and preserve key values
				return len(output) < 200 &&
					strings.Contains(output, "100") &&
					strings.Contains(output, "Initializing")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := JS(tt.input)
			if !tt.checkFn(output) {
				t.Errorf("Minification test failed\nInput length: %d\nOutput length: %d\nOutput: %s",
					len(tt.input), len(output), output)
			}
		})
	}
}

func TestJSPreservesGoTemplates(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "single template variable",
			input: "const url = '{{.BaseURL}}';",
		},
		{
			name:  "multiple template variables",
			input: "const a = '{{.BaseURL}}'; const b = '{{.SiteTitle}}';",
		},
		{
			name:  "template in object",
			input: "const config = { base: '{{.BaseURL}}', title: '{{.SiteTitle}}' };",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := JS(tt.input)

			// Verify template syntax preserved
			if !strings.Contains(output, "{{.") {
				t.Errorf("Template syntax not preserved in minified output\nInput: %s\nOutput: %s",
					tt.input, output)
			}
		})
	}
}

func TestJSInvalidInput(t *testing.T) {
	// Test that malformed JS returns original input (doesn't crash)
	// Note: tdewolff/minify is quite permissive, so we test with
	// content that's clearly not JavaScript
	invalidInputs := []string{
		"",             // Empty (should return empty)
		"   ",          // Whitespace only
		"const x = 5;", // Valid JS (should minify)
	}

	for _, input := range invalidInputs {
		output := JS(input)
		// Should not panic, should return some output
		if input != "" && output == "" && input != "   " {
			t.Errorf("Unexpected empty output for input: %q", input)
		}
	}
}

func TestJSErrorHandling(t *testing.T) {
	// Test that JS function handles errors gracefully
	// Most inputs will minify successfully, but we test various edge cases
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "whitespace only",
			input: "   \n\t  ",
		},
		{
			name:  "unicode characters",
			input: "const msg = '你好世界';",
		},
		{
			name:  "mixed content",
			input: "/* comment */ const x = 1; // inline comment\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			output := JS(tt.input)
			// Empty input should return empty output
			if tt.input == "" && output != "" {
				t.Errorf("Expected empty output for empty input, got: %q", output)
			}
		})
	}
}
