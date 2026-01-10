package styles

import (
	"strings"
	"testing"
)

func TestLayoutCSS(t *testing.T) {
	if LayoutCSS == "" {
		t.Error("LayoutCSS should not be empty")
	}

	// Check for structural CSS content in layout.css
	checks := []string{
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
		if !strings.Contains(LayoutCSS, check) {
			t.Errorf("LayoutCSS should contain %q", check)
		}
	}
}

func TestDocsCSS(t *testing.T) {
	if DocsCSS == "" {
		t.Error("DocsCSS should not be empty")
	}

	// Check for styling CSS content in docs.css (colors, fonts, theming)
	checks := []string{
		":root",
		"--bg-primary",
		"--text-primary",
		"[data-theme=\"dark\"]",
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

	// Blog theme should have color variables (styling, not layout)
	if !strings.Contains(BlogCSS, "--bg-primary") {
		t.Error("BlogCSS should have --bg-primary color variable")
	}

	// Blog theme should have dark mode styling
	if !strings.Contains(BlogCSS, "[data-theme=\"dark\"]") {
		t.Error("BlogCSS should have dark mode styling")
	}
}

func TestVanillaCSS(t *testing.T) {
	if VanillaCSS == "" {
		t.Error("VanillaCSS should not be empty")
	}

	// Vanilla theme is a customization skeleton with commented examples
	if !strings.Contains(VanillaCSS, "Customization Skeleton") {
		t.Error("VanillaCSS should contain 'Customization Skeleton' header")
	}

	// Vanilla theme should have :root section (even if mostly commented)
	if !strings.Contains(VanillaCSS, ":root") {
		t.Error("VanillaCSS should contain :root selector")
	}

	// Vanilla theme should not have many actual style declarations
	// (it's mostly commented examples)
	colorCount := strings.Count(VanillaCSS, "color:")
	if colorCount > 5 {
		t.Logf("VanillaCSS has %d 'color:' declarations - should be minimal", colorCount)
	}
}

func TestGetCSS(t *testing.T) {
	tests := []struct {
		theme       string
		expectedCSS *string // The theme-specific CSS that should be included
	}{
		{"docs", &DocsCSS},
		{"blog", &BlogCSS},
		{"vanilla", &VanillaCSS},
		{"", &DocsCSS},        // Default
		{"invalid", &DocsCSS}, // Unknown falls back to docs
	}

	for _, tt := range tests {
		css := GetCSS(tt.theme)
		if css == "" {
			t.Errorf("GetCSS(%q) should not return empty string", tt.theme)
		}
		// GetCSS returns LayoutCSS + "\n" + themeCSS
		expectedCombined := LayoutCSS + "\n" + *tt.expectedCSS
		if css != expectedCombined {
			t.Errorf("GetCSS(%q) should return combined layout + theme CSS", tt.theme)
		}
		// Verify it contains both layout and theme content
		if !strings.Contains(css, ".sidebar") {
			t.Errorf("GetCSS(%q) should contain layout CSS (.sidebar)", tt.theme)
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
		{"", false}, // Empty is valid (defaults to docs)
		{"invalid", true},
		{"DOCS", true}, // Case sensitive
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
	// CSS should be combined LayoutCSS + DocsCSS for backward compatibility
	expected := LayoutCSS + "\n" + DocsCSS
	if CSS != expected {
		t.Error("CSS should equal LayoutCSS + DocsCSS for backward compatibility")
	}
}
