package toc

import (
	"strings"
	"testing"
)

func TestExtractTOC(t *testing.T) {
	html := `
		<h2 id="intro">Introduction</h2>
		<p>Some content</p>
		<h2 id="setup">Setup</h2>
		<h3 id="install">Installation</h3>
		<h3 id="config">Configuration</h3>
		<h2 id="usage">Usage</h2>
	`

	toc := ExtractTOC(html, 3)

	if toc == nil {
		t.Fatal("TOC should not be nil")
	}

	if len(toc.Items) != 3 {
		t.Errorf("expected 3 top-level items, got %d", len(toc.Items))
	}

	// Check first item
	if toc.Items[0].ID != "intro" {
		t.Errorf("first item ID = %q, want %q", toc.Items[0].ID, "intro")
	}
	if toc.Items[0].Text != "Introduction" {
		t.Errorf("first item Text = %q, want %q", toc.Items[0].Text, "Introduction")
	}

	// Check nested items under Setup
	if len(toc.Items[1].Children) != 2 {
		t.Errorf("Setup should have 2 children, got %d", len(toc.Items[1].Children))
	}
}

func TestExtractTOCMinItems(t *testing.T) {
	html := `
		<h2 id="one">One</h2>
		<h2 id="two">Two</h2>
	`

	// With minItems=3, should return nil
	toc := ExtractTOC(html, 3)
	if toc != nil {
		t.Error("TOC should be nil when fewer than minItems")
	}

	// With minItems=2, should return TOC
	toc = ExtractTOC(html, 2)
	if toc == nil {
		t.Error("TOC should not be nil with enough items")
	}
}

func TestExtractTOCEmpty(t *testing.T) {
	toc := ExtractTOC("<p>No headings here</p>", 1)
	if toc != nil {
		t.Error("TOC should be nil for content without headings")
	}
}

func TestExtractTOCStripsHTML(t *testing.T) {
	html := `<h2 id="test"><strong>Bold</strong> Title</h2><h2 id="test2">Another</h2><h2 id="test3">Third</h2>`
	toc := ExtractTOC(html, 1)

	if toc == nil || len(toc.Items) == 0 {
		t.Fatal("TOC should have items")
	}

	if toc.Items[0].Text != "Bold Title" {
		t.Errorf("Text = %q, want %q", toc.Items[0].Text, "Bold Title")
	}
}

func TestRenderTOC(t *testing.T) {
	toc := &PageTOC{
		Items: []*Item{
			{ID: "intro", Text: "Introduction", Level: 2},
			{
				ID: "setup", Text: "Setup", Level: 2,
				Children: []*Item{
					{ID: "install", Text: "Installation", Level: 3},
				},
			},
		},
	}

	html := string(RenderTOC(toc))

	if !strings.Contains(html, "toc-sidebar") {
		t.Error("should contain toc-sidebar class")
	}
	if !strings.Contains(html, `href="#intro"`) {
		t.Error("should contain intro link")
	}
	if !strings.Contains(html, `href="#install"`) {
		t.Error("should contain nested install link")
	}
	if !strings.Contains(html, "On this page") {
		t.Error("should contain title")
	}
}

func TestRenderTOCEmpty(t *testing.T) {
	html := RenderTOC(nil)
	if html != "" {
		t.Error("nil TOC should return empty string")
	}

	html = RenderTOC(&PageTOC{})
	if html != "" {
		t.Error("empty TOC should return empty string")
	}
}

func TestHasTOC(t *testing.T) {
	html := `
		<h2 id="one">One</h2>
		<h2 id="two">Two</h2>
		<h2 id="three">Three</h2>
	`

	if !HasTOC(html, 3) {
		t.Error("should have TOC with 3 headings")
	}
	if HasTOC(html, 5) {
		t.Error("should not have TOC with only 3 headings when minItems=5")
	}
}
