package instant

import (
	"strings"
	"testing"
)

func TestInstantNavJSIsMinified(t *testing.T) {
	// The raw JS contains function names like 'init', 'handleMouseOver', etc.
	// When minified, whitespace is removed and variable names are shortened

	// Check that InstantNavJS is shorter than instantNavJSRaw
	if len(InstantNavJS) >= len(instantNavJSRaw) {
		t.Errorf("InstantNavJS should be minified (shorter than raw)\nRaw length: %d\nMinified length: %d",
			len(instantNavJSRaw), len(InstantNavJS))
	}

	// Check that it doesn't contain excessive whitespace (multiple spaces or tabs)
	if strings.Contains(InstantNavJS, "    ") {
		t.Error("InstantNavJS contains indentation, should be minified")
	}

	// Check that key functionality is preserved (strings are kept)
	if !strings.Contains(InstantNavJS, ".content") {
		t.Error("InstantNavJS should contain '.content' selector")
	}

	if !strings.Contains(InstantNavJS, "instant:navigated") {
		t.Error("InstantNavJS should contain 'instant:navigated' event name")
	}
}

func TestInstantNavJSPreservesSelectors(t *testing.T) {
	// Verify that CSS selectors are preserved (they're in strings)
	selectors := []string{".content", ".tree-nav", ".toc", ".breadcrumbs", "title"}

	for _, sel := range selectors {
		if !strings.Contains(InstantNavJS, sel) {
			t.Errorf("InstantNavJS should contain selector %q", sel)
		}
	}
}

func TestInstantNavJSNotEmpty(t *testing.T) {
	if InstantNavJS == "" {
		t.Error("InstantNavJS should not be empty")
	}

	// Should be at least 500 bytes (minified version of the script)
	if len(InstantNavJS) < 500 {
		t.Errorf("InstantNavJS seems too short: %d bytes", len(InstantNavJS))
	}
}
