package navigation

import (
	"strings"
	"testing"

	"github.com/wusher/volcano/internal/tree"
)

func TestFlattenTreeForPagination(t *testing.T) {
	root := tree.NewNode("", "", true)

	home := tree.NewNode("Home", "home.md", false)
	home.SourcePath = "home.md"
	root.AddChild(home)

	docs := tree.NewNode("Docs", "docs", true)
	docs.SourcePath = "docs"
	docs.HasIndex = true

	index := tree.NewNode("Docs", "docs/index.md", false)
	index.SourcePath = "docs/index.md"
	docs.AddChild(index)

	intro := tree.NewNode("Intro", "docs/intro.md", false)
	intro.SourcePath = "docs/intro.md"
	docs.AddChild(intro)

	root.AddChild(docs)

	pages := FlattenTreeForPagination(root)
	if len(pages) < 3 {
		t.Fatalf("FlattenTreeForPagination() returned %d pages, want at least 3", len(pages))
	}

	if pages[0].Name != "Home" {
		t.Errorf("pages[0].Name = %q, want %q", pages[0].Name, "Home")
	}

	foundIndex := false
	foundIntro := false
	for _, page := range pages {
		if page.Path == "docs/index.md" {
			foundIndex = true
		}
		if page.Path == "docs/intro.md" {
			foundIntro = true
		}
	}

	if !foundIndex {
		t.Error("FlattenTreeForPagination() should include folder index")
	}
	if !foundIntro {
		t.Error("FlattenTreeForPagination() should include folder pages")
	}
}

func TestBuildPageNavigationWithSection(_ *testing.T) {
	// Create a tree with grandparent -> parent -> child structure
	root := tree.NewNode("", "", true)

	docs := tree.NewNode("Docs", "docs", true)
	docs.HasIndex = true
	root.AddChild(docs)

	// Create an index for docs
	index := tree.NewNode("Docs", "docs/index.md", false)
	index.SourcePath = "docs/index.md"
	index.Parent = docs
	docs.AddChild(index)

	// Create subfolder with pages
	subfolder := tree.NewNode("Guide", "docs/guide", true)
	subfolder.Parent = docs
	docs.AddChild(subfolder)

	intro := tree.NewNode("Intro", "docs/guide/intro.md", false)
	intro.SourcePath = "docs/guide/intro.md"
	intro.Parent = subfolder
	subfolder.AddChild(intro)

	advanced := tree.NewNode("Advanced", "docs/guide/advanced.md", false)
	advanced.SourcePath = "docs/guide/advanced.md"
	advanced.Parent = subfolder
	subfolder.AddChild(advanced)

	// Flatten and build navigation for intro page
	pages := FlattenTreeForPagination(root)

	// Find intro page in flattened list
	var introIdx int
	for i, p := range pages {
		if p.Path == "docs/guide/intro.md" {
			introIdx = i
			break
		}
	}

	if introIdx > 0 {
		nav := BuildPageNavigation(pages[introIdx], pages)
		// Navigation should be built
		_ = nav
	}
}

func TestRenderPageNavigation_Empty(t *testing.T) {
	nav := PageNavigation{}
	result := RenderPageNavigation(nav)
	if result != "" {
		t.Errorf("Expected empty result for empty navigation, got %q", result)
	}
}

func TestRenderPageNavigation_OnlyPrevious(t *testing.T) {
	nav := PageNavigation{
		Previous: &NavLink{Title: "Previous", URL: "/prev/", Section: ""},
	}
	result := RenderPageNavigation(nav)
	if !strings.Contains(string(result), "Previous") {
		t.Error("Expected result to contain 'Previous'")
	}
	if !strings.Contains(string(result), "/prev/") {
		t.Error("Expected result to contain previous URL")
	}
}

func TestRenderPageNavigation_OnlyNext(t *testing.T) {
	nav := PageNavigation{
		Next: &NavLink{Title: "Next Page", URL: "/next/", Section: ""},
	}
	result := RenderPageNavigation(nav)
	if !strings.Contains(string(result), "Next Page") {
		t.Error("Expected result to contain 'Next Page'")
	}
	if !strings.Contains(string(result), "/next/") {
		t.Error("Expected result to contain next URL")
	}
}

func TestRenderPageNavigation_WithSection(t *testing.T) {
	nav := PageNavigation{
		Previous: &NavLink{Title: "Previous", URL: "/prev/", Section: "Guide"},
		Next:     &NavLink{Title: "Next", URL: "/next/", Section: "Guide"},
	}
	result := RenderPageNavigation(nav)
	// Section is stored but not rendered in current implementation
	// Just verify the nav renders successfully with sections set
	if !strings.Contains(string(result), "Previous") {
		t.Error("Expected result to contain 'Previous'")
	}
	if !strings.Contains(string(result), "Next") {
		t.Error("Expected result to contain 'Next'")
	}
}
