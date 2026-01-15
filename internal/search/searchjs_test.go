package search

import (
	"strings"
	"testing"
)

func TestGenerateSearchJS(t *testing.T) {
	tests := []struct {
		name            string
		baseURL         string
		expectedStrings []string
	}{
		{
			name:    "empty base URL",
			baseURL: "",
			expectedStrings: []string{
				"const baseURL = ''",
				"command-palette",
				"search-index.json",
				"openCommandPalette",
				"closeCommandPalette",
				"ArrowDown",
				"ArrowUp",
				"Escape",
			},
		},
		{
			name:    "with base URL",
			baseURL: "/docs",
			expectedStrings: []string{
				"const baseURL = '/docs'",
				"baseURL + '/search-index.json'",
			},
		},
		{
			name:    "with trailing slash base URL",
			baseURL: "/site/",
			expectedStrings: []string{
				"const baseURL = '/site/'",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := GenerateSearchJS(tc.baseURL)

			if result == "" {
				t.Error("GenerateSearchJS() should return non-empty string")
			}

			for _, expected := range tc.expectedStrings {
				if !strings.Contains(result, expected) {
					t.Errorf("GenerateSearchJS() should contain %q", expected)
				}
			}
		})
	}
}

func TestGenerateSearchJS_ContainsEventListeners(t *testing.T) {
	result := GenerateSearchJS("")

	// Check for key event listeners
	if !strings.Contains(result, "addEventListener") {
		t.Error("GenerateSearchJS() should contain event listeners")
	}

	// Check for keyboard shortcuts
	if !strings.Contains(result, "metaKey") && !strings.Contains(result, "ctrlKey") {
		t.Error("GenerateSearchJS() should handle Cmd/Ctrl+K shortcut")
	}

	// Check for search functionality
	if !strings.Contains(result, "doSearch") {
		t.Error("GenerateSearchJS() should contain doSearch function")
	}
}

func TestGenerateSearchJS_ContainsHTMLInjection(t *testing.T) {
	result := GenerateSearchJS("")

	// Check for HTML structure injection
	if !strings.Contains(result, "command-palette-modal") {
		t.Error("GenerateSearchJS() should inject modal HTML")
	}

	if !strings.Contains(result, "command-palette-input") {
		t.Error("GenerateSearchJS() should inject input element")
	}

	if !strings.Contains(result, "command-palette-results") {
		t.Error("GenerateSearchJS() should inject results container")
	}
}

func TestGenerateSearchJS_ContainsEscapeFunction(t *testing.T) {
	result := GenerateSearchJS("")

	// Check for HTML escape function
	if !strings.Contains(result, "escapeHtml") {
		t.Error("GenerateSearchJS() should contain escapeHtml function for XSS prevention")
	}

	// Check escape patterns
	if !strings.Contains(result, "&amp;") {
		t.Error("GenerateSearchJS() should escape ampersands")
	}

	if !strings.Contains(result, "&lt;") {
		t.Error("GenerateSearchJS() should escape less-than signs")
	}
}

func TestGenerateSearchJS_IsValidJavaScript(t *testing.T) {
	result := GenerateSearchJS("")

	// Basic syntax checks
	if !strings.HasPrefix(result, "(function()") {
		t.Error("GenerateSearchJS() should be an IIFE")
	}

	if !strings.HasSuffix(result, "})();") {
		t.Error("GenerateSearchJS() should end with IIFE closure")
	}

	// Check balanced braces (simple check)
	openBraces := strings.Count(result, "{")
	closeBraces := strings.Count(result, "}")
	if openBraces != closeBraces {
		t.Errorf("GenerateSearchJS() has unbalanced braces: %d open, %d close", openBraces, closeBraces)
	}

	// Check balanced parentheses
	openParens := strings.Count(result, "(")
	closeParens := strings.Count(result, ")")
	if openParens != closeParens {
		t.Errorf("GenerateSearchJS() has unbalanced parentheses: %d open, %d close", openParens, closeParens)
	}
}
