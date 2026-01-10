package navigation

import (
	"strings"
	"testing"

	"github.com/wusher/volcano/internal/tree"
)

func TestBuildBreadcrumbs(t *testing.T) {
	// Create a tree structure: root -> guides -> intro
	root := tree.NewNode("Root", "", true)
	guides := tree.NewNode("Guides", "guides", true)
	intro := tree.NewNode("Introduction", "guides/intro.md", false)

	root.AddChild(guides)
	guides.AddChild(intro)

	crumbs := BuildBreadcrumbs(intro, "My Site")

	// Should have: Home, Guides, Introduction
	if len(crumbs) != 3 {
		t.Fatalf("expected 3 breadcrumbs, got %d", len(crumbs))
	}

	if crumbs[0].Label != "My Site" {
		t.Errorf("first crumb label = %q, want %q", crumbs[0].Label, "My Site")
	}
	if crumbs[0].URL != "/" {
		t.Errorf("first crumb URL = %q, want %q", crumbs[0].URL, "/")
	}

	if crumbs[1].Label != "Guides" {
		t.Errorf("second crumb label = %q, want %q", crumbs[1].Label, "Guides")
	}

	if !crumbs[2].Current {
		t.Error("last crumb should be current")
	}
}

func TestBuildBreadcrumbsRootFile(t *testing.T) {
	root := tree.NewNode("Root", "", true)
	index := tree.NewNode("Home", "index.md", false)
	root.AddChild(index)

	crumbs := BuildBreadcrumbs(index, "My Site")

	// Should have: Home (site), Home (page)
	if len(crumbs) < 1 {
		t.Fatalf("expected at least 1 breadcrumb, got %d", len(crumbs))
	}
}

func TestRenderBreadcrumbs(t *testing.T) {
	crumbs := []Breadcrumb{
		{Label: "Home", URL: "/", Current: false},
		{Label: "Guides", URL: "/guides/", Current: false},
		{Label: "Intro", URL: "/guides/intro/", Current: true},
	}

	html := string(RenderBreadcrumbs(crumbs))

	if !strings.Contains(html, "breadcrumbs") {
		t.Error("should contain breadcrumbs class")
	}
	if !strings.Contains(html, "Home") {
		t.Error("should contain Home")
	}
	if !strings.Contains(html, "Guides") {
		t.Error("should contain Guides")
	}
	if !strings.Contains(html, `aria-current="page"`) {
		t.Error("should mark current page")
	}
}

func TestRenderBreadcrumbsEmpty(t *testing.T) {
	html := RenderBreadcrumbs(nil)
	if html != "" {
		t.Error("empty breadcrumbs should return empty string")
	}
}

func TestBuildPageNavigation(t *testing.T) {
	// Need to set SourcePath since BuildPageNavigation compares by SourcePath
	page1 := &tree.Node{Name: "Page 1", Path: "page1.md", SourcePath: "/test/page1.md"}
	page2 := &tree.Node{Name: "Page 2", Path: "page2.md", SourcePath: "/test/page2.md"}
	page3 := &tree.Node{Name: "Page 3", Path: "page3.md", SourcePath: "/test/page3.md"}

	allPages := []*tree.Node{page1, page2, page3}

	// Test middle page
	nav := BuildPageNavigation(page2, allPages)
	if nav.Previous == nil {
		t.Error("page2 should have previous")
	}
	if nav.Next == nil {
		t.Error("page2 should have next")
	}

	// Test first page
	nav = BuildPageNavigation(page1, allPages)
	if nav.Previous != nil {
		t.Error("page1 should not have previous")
	}
	if nav.Next == nil {
		t.Error("page1 should have next")
	}

	// Test last page
	nav = BuildPageNavigation(page3, allPages)
	if nav.Previous == nil {
		t.Error("page3 should have previous")
	}
	if nav.Next != nil {
		t.Error("page3 should not have next")
	}
}

func TestRenderPageNavigation(t *testing.T) {
	nav := PageNavigation{
		Previous: &NavLink{Title: "Prev", URL: "/prev/"},
		Next:     &NavLink{Title: "Next", URL: "/next/"},
	}

	html := string(RenderPageNavigation(nav))

	if !strings.Contains(html, "page-nav") {
		t.Error("should contain page-nav class")
	}
	if !strings.Contains(html, "Prev") {
		t.Error("should contain Prev")
	}
	if !strings.Contains(html, "Next") {
		t.Error("should contain Next")
	}
	if !strings.Contains(html, "/prev/") {
		t.Error("should contain prev URL")
	}
}

func TestRenderPageNavigationEmpty(t *testing.T) {
	nav := PageNavigation{}
	html := RenderPageNavigation(nav)
	if html != "" {
		t.Error("empty nav should return empty string")
	}
}
