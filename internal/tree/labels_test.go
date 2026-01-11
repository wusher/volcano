package tree

import (
	"testing"
)

func TestCleanLabel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"getting-started.md", "Getting Started"},
		{"api_reference.md", "Api Reference"},
		{"FAQ.md", "FAQ"},
		{"README.md", "README"},
		{"01-introduction.md", "Introduction"},
		{"001_setup.md", "Setup"},
		{"2024-01-01-my-post.md", "My Post"},
		{"hello-world.md", "Hello World"},
		{"simple.md", "Simple"},
		{"multiple---dashes.md", "Multiple Dashes"},
		{"under__scores.md", "Under Scores"},
		{"MixedCase.md", "MixedCase"},
		{"ALLCAPS.md", "ALLCAPS"},
		{"API.md", "API"},
		{"file.markdown", "File"},
		{"no-extension", "No Extension"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := CleanLabel(tt.input)
			if result != tt.expected {
				t.Errorf("CleanLabel(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsMarkdownFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"file.md", true},
		{"file.MD", true},
		{"file.Md", true},
		{"file.markdown", true},
		{"file.MARKDOWN", true},
		{"file.txt", false},
		{"file.html", false},
		{"file", false},
		{"readme.md", true},
		{"README.MD", true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := IsMarkdownFile(tt.filename)
			if result != tt.expected {
				t.Errorf("IsMarkdownFile(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestIsHidden(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{".git", true},
		{".hidden", true},
		{".DS_Store", true},
		{"visible", false},
		{"file.md", false},
		{".", true},
		{"..", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsHidden(tt.name)
			if result != tt.expected {
				t.Errorf("IsHidden(%q) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestIsIndexFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"index.md", true},
		{"INDEX.MD", true},
		{"Index.md", true},
		{"index.markdown", true},
		{"readme.md", true},
		{"README.MD", true},
		{"readme.markdown", true},
		{"other.md", false},
		{"index.html", false},
		{"indexmd", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := IsIndexFile(tt.filename)
			if result != tt.expected {
				t.Errorf("IsIndexFile(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestRemoveLeadingNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"01-intro", "intro"},
		{"001_setup", "setup"},
		{"2024-01-01-post", "post"},
		{"2023-06-19 California Real Estate Law", "California Real Estate Law"}, // Date with space separator
		{"2023-06-19 Law sections", "Law sections"},                             // Date with space separator
		{"no-numbers", "no-numbers"},
		{"123", "123"},
		{"12-34-56-test", "34-56-test"},
		// Space-dash-space number prefixes (Obsidian style)
		{"6 - Sign Up for Ramp", "Sign Up for Ramp"},
		{"1 - Bamboo Onboarding", "Bamboo Onboarding"},
		{"10 - New Hire Tasks", "New Hire Tasks"},
		// Space-underscore-space number prefixes
		{"6 _ Sign Up for Ramp", "Sign Up for Ramp"},
		// Dot-space number prefixes
		{"0. Inbox", "Inbox"},
		{"1. Projects", "Projects"},
		// Years should NOT be stripped (4 digits, no leading zeros, space separator)
		{"2023 Goals", "2023 Goals"},
		{"1999 Report", "1999 Report"},
		// Ordering prefixes with space SHOULD be stripped (leading zeros or short numbers)
		{"01 getting-started", "getting-started"},
		{"001 chapter one", "chapter one"},
		{"1 introduction", "introduction"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := removeLeadingNumbers(tt.input)
			if result != tt.expected {
				t.Errorf("removeLeadingNumbers(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"API Reference", "api-reference"},
		{"0. Inbox", "inbox"},
		{"01-introduction", "introduction"},
		{"2024-01-01-my-post", "my-post"},
		// Space-dash-space number prefixes (Obsidian style)
		{"6 - Sign Up for Ramp", "sign-up-for-ramp"},
		{"1 - Bamboo Onboarding", "bamboo-onboarding"},
		{"10 - New Hire Tasks", "new-hire-tasks"},
		{"1. Projects", "projects"},
		// Years should NOT be stripped
		{"2023 Goals", "2023-goals"},
		{"1999 Report", "1999-report"},
		// Edge cases
		{"simple", "simple"},
		{"UPPERCASE", "uppercase"},
		{"with   multiple   spaces", "with-multiple-spaces"},
		{"special!@#chars", "specialchars"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Slugify(tt.input)
			if result != tt.expected {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsDatePrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"2024-01-01-", true},
		{"2023-12-31-", true},
		{"1999-05-15-", true},
		{"2024-1-01-", false},  // wrong format
		{"2024-01-1-", false},  // wrong format
		{"2024/01/01-", false}, // wrong separator
		{"2024-01-01", false},  // missing trailing dash
		{"short", false},       // too short
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isDatePrefix(tt.input)
			if result != tt.expected {
				t.Errorf("isDatePrefix(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "Hello World"},
		{"HELLO WORLD", "HELLO WORLD"},
		{"hello WORLD", "Hello WORLD"},
		{"api reference", "Api Reference"},
		{"FAQ", "FAQ"},
		{"", ""},
		{"single", "Single"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := titleCase(tt.input)
			if result != tt.expected {
				t.Errorf("titleCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsAllUppercase(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"FAQ", true},
		{"API", true},
		{"ABC", true},
		{"Abc", false},
		{"abc", false},
		{"123", false},
		{"ABC123", true},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isAllUppercase(tt.input)
			if result != tt.expected {
				t.Errorf("isAllUppercase(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
