package navigation

import (
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
