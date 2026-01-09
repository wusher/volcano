package styles

import (
	"strings"
	"testing"
)

func TestCSS(t *testing.T) {
	if CSS == "" {
		t.Error("CSS should not be empty")
	}

	// Check for essential CSS content
	checks := []string{
		":root",
		"--bg-primary",
		"--text-primary",
		"[data-theme=\"dark\"]",
		".sidebar",
		".content",
		".prose",
		".tree-nav",
		".theme-toggle",
		".mobile-menu-btn",
		".drawer-backdrop",
		"@media",
	}

	for _, check := range checks {
		if !strings.Contains(CSS, check) {
			t.Errorf("CSS should contain %q", check)
		}
	}
}

func TestGetCSS(t *testing.T) {
	css := GetCSS()
	if css == "" {
		t.Error("GetCSS() should not return empty string")
	}
	if css != CSS {
		t.Error("GetCSS() should return same value as CSS variable")
	}
}
