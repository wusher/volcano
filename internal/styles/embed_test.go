package styles

import (
	"strings"
	"testing"
)

func TestDocsCSS(t *testing.T) {
	if DocsCSS == "" {
		t.Error("DocsCSS should not be empty")
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
		if !strings.Contains(DocsCSS, check) {
			t.Errorf("DocsCSS should contain %q", check)
		}
	}
}

func TestBlogCSS(t *testing.T) {
	if BlogCSS == "" {
		t.Error("BlogCSS should not be empty")
	}

	// Blog theme should have elegant serif font stack
	if !strings.Contains(BlogCSS, "Charter") {
		t.Error("BlogCSS should contain Charter font family")
	}

	// Blog theme should have sidebar
	if !strings.Contains(BlogCSS, "--sidebar-width: 240px") {
		t.Error("BlogCSS should have 240px sidebar width")
	}

	// Blog theme should have refined table styling
	if !strings.Contains(BlogCSS, ".prose thead") {
		t.Error("BlogCSS should have refined table styling")
	}
}

func TestVanillaCSS(t *testing.T) {
	if VanillaCSS == "" {
		t.Error("VanillaCSS should not be empty")
	}

	// Vanilla theme should have layout but minimal styling
	if !strings.Contains(VanillaCSS, ".sidebar") {
		t.Error("VanillaCSS should contain .sidebar")
	}

	// Vanilla theme should not have many color declarations
	// (it uses inherit/defaults)
	colorCount := strings.Count(VanillaCSS, "color:")
	if colorCount > 5 {
		t.Logf("VanillaCSS has %d 'color:' declarations - should be minimal", colorCount)
	}
}

func TestGetCSS(t *testing.T) {
	tests := []struct {
		theme    string
		expected *string
	}{
		{"docs", &DocsCSS},
		{"blog", &BlogCSS},
		{"vanilla", &VanillaCSS},
		{"", &DocsCSS},       // Default
		{"invalid", &DocsCSS}, // Unknown falls back to docs
	}

	for _, tt := range tests {
		css := GetCSS(tt.theme)
		if css == "" {
			t.Errorf("GetCSS(%q) should not return empty string", tt.theme)
		}
		if css != *tt.expected {
			t.Errorf("GetCSS(%q) returned unexpected CSS", tt.theme)
		}
	}
}

func TestValidateTheme(t *testing.T) {
	tests := []struct {
		theme   string
		wantErr bool
	}{
		{"docs", false},
		{"blog", false},
		{"vanilla", false},
		{"", false},         // Empty is valid (defaults to docs)
		{"invalid", true},
		{"DOCS", true},      // Case sensitive
	}

	for _, tt := range tests {
		err := ValidateTheme(tt.theme)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateTheme(%q) error = %v, wantErr %v", tt.theme, err, tt.wantErr)
		}
	}
}

func TestCSSBackwardCompat(t *testing.T) {
	// CSS variable should still work for backward compatibility
	if CSS == "" {
		t.Error("CSS variable should not be empty")
	}
	// CSS should be the same as DocsCSS
	if CSS != DocsCSS {
		t.Error("CSS should equal DocsCSS for backward compatibility")
	}
}
