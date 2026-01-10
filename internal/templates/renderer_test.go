package templates

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"github.com/wusher/volcano/internal/tree"
)

func TestNewRenderer(t *testing.T) {
	r, err := NewRenderer("body { color: black; }")
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}
	if r == nil {
		t.Fatal("NewRenderer() returned nil")
	}
	if r.tmpl == nil {
		t.Fatal("Renderer.tmpl is nil")
	}
}

func TestRendererRender(t *testing.T) {
	r, err := NewRenderer("body { color: black; }")
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	var buf bytes.Buffer
	data := PageData{
		SiteTitle:   "Test Site",
		PageTitle:   "Test Page",
		Content:     template.HTML("<p>Hello World</p>"),
		Navigation:  template.HTML("<ul><li><a href=\"/\">Home</a></li></ul>"),
		CurrentPath: "/test/",
	}

	err = r.Render(&buf, data)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	html := buf.String()

	// Check for basic structure
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("Should contain DOCTYPE")
	}
	if !strings.Contains(html, "<html lang=\"en\">") {
		t.Error("Should contain html tag with lang")
	}
	if !strings.Contains(html, "<title>Test Page - Test Site</title>") {
		t.Error("Should contain title with page and site title")
	}
	if !strings.Contains(html, "Hello World") {
		t.Error("Should contain content")
	}
	if !strings.Contains(html, "body { color: black; }") {
		t.Error("Should contain CSS")
	}
	if !strings.Contains(html, "Test Site") {
		t.Error("Should contain site title")
	}
}

func TestRendererRenderToString(t *testing.T) {
	r, err := NewRenderer("/* CSS */")
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	data := PageData{
		SiteTitle:   "Site",
		PageTitle:   "Page",
		Content:     template.HTML("<p>Content</p>"),
		Navigation:  template.HTML("<nav></nav>"),
		CurrentPath: "/",
	}

	html, err := r.RenderToString(data)
	if err != nil {
		t.Fatalf("RenderToString() error = %v", err)
	}

	if html == "" {
		t.Error("RenderToString() returned empty string")
	}
	if !strings.Contains(html, "Content") {
		t.Error("Should contain content")
	}
}

func TestRenderNavigation(t *testing.T) {
	// Create a test tree
	root := tree.NewNode("", "", true)

	// Add files
	home := tree.NewNode("Home", "index.md", false)
	about := tree.NewNode("About", "about.md", false)
	root.AddChild(home)
	root.AddChild(about)

	// Add folder with children
	guides := tree.NewNode("Guides", "guides", true)
	guides.HasIndex = true
	guides.IndexPath = "guides/index.md"
	intro := tree.NewNode("Introduction", "guides/intro.md", false)
	advanced := tree.NewNode("Advanced", "guides/advanced.md", false)
	guides.AddChild(intro)
	guides.AddChild(advanced)
	root.AddChild(guides)

	html := RenderNavigation(root, "/guides/intro/")

	htmlStr := string(html)

	// Check structure
	if !strings.Contains(htmlStr, "<ul role=\"tree\">") {
		t.Error("Should contain ul with tree role")
	}
	if !strings.Contains(htmlStr, "Home") {
		t.Error("Should contain Home link")
	}
	if !strings.Contains(htmlStr, "About") {
		t.Error("Should contain About link")
	}
	if !strings.Contains(htmlStr, "Guides") {
		t.Error("Should contain Guides folder")
	}
	if !strings.Contains(htmlStr, "Introduction") {
		t.Error("Should contain Introduction link")
	}
	if !strings.Contains(htmlStr, "Advanced") {
		t.Error("Should contain Advanced link")
	}
	if !strings.Contains(htmlStr, "class=\"folder\"") {
		t.Error("Should have folder class")
	}
	if !strings.Contains(htmlStr, "active") {
		t.Error("Should have active class on current page")
	}
}

func TestRenderNavigationEmpty(t *testing.T) {
	root := tree.NewNode("", "", true)
	html := RenderNavigation(root, "/")

	if string(html) != "" {
		t.Error("Empty tree should render empty navigation")
	}
}

func TestRenderNavigationFolderWithoutIndex(t *testing.T) {
	root := tree.NewNode("", "", true)

	// Folder without index - now all folders are clickable (auto-index)
	docs := tree.NewNode("Docs", "docs", true)
	docs.HasIndex = false
	page := tree.NewNode("Page", "docs/page.md", false)
	docs.AddChild(page)
	root.AddChild(docs)

	html := RenderNavigation(root, "/")
	htmlStr := string(html)

	// All folders should have folder-link now (for auto-generated indexes)
	if !strings.Contains(htmlStr, "folder-link") {
		t.Error("All folders should have folder-link for auto-index")
	}
	// Should link to the slugified folder path
	if !strings.Contains(htmlStr, "href=\"/docs/\"") {
		t.Error("Folder link should point to /docs/")
	}
}

func TestRenderNavigationActiveState(t *testing.T) {
	root := tree.NewNode("", "", true)

	page1 := tree.NewNode("Page 1", "page1.md", false)
	page2 := tree.NewNode("Page 2", "page2.md", false)
	root.AddChild(page1)
	root.AddChild(page2)

	// Test with page1 active
	html1 := RenderNavigation(root, "/page1/")
	if !strings.Contains(string(html1), "href=\"/page1/\" class=\"file-link active\"") {
		t.Error("Page 1 should be active")
	}

	// Test with page2 active
	html2 := RenderNavigation(root, "/page2/")
	if !strings.Contains(string(html2), "href=\"/page2/\" class=\"file-link active\"") {
		t.Error("Page 2 should be active")
	}
}

func TestPageDataWithEmptyTitle(t *testing.T) {
	r, err := NewRenderer("")
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	data := PageData{
		SiteTitle:   "",
		PageTitle:   "Only Page Title",
		Content:     template.HTML("<p>Content</p>"),
		Navigation:  template.HTML(""),
		CurrentPath: "/",
	}

	html, err := r.RenderToString(data)
	if err != nil {
		t.Fatalf("RenderToString() error = %v", err)
	}

	// Title should not have " - " when site title is empty
	if !strings.Contains(html, "<title>Only Page Title</title>") {
		t.Errorf("Title should be just page title when site title is empty, got: %s", html)
	}
}

func TestBuildTopNavItems(t *testing.T) {
	root := tree.NewNode("", "", true)

	index := tree.NewNode("Home", "index.md", false)
	index.FileName = "index.md"
	index.SourcePath = "index.md"
	root.AddChild(index)

	readme := tree.NewNode("Readme", "readme.md", false)
	readme.FileName = "readme.md"
	readme.SourcePath = "readme.md"
	root.AddChild(readme)

	alpha := tree.NewNode("Alpha", "2024-01-10-alpha.md", false)
	alpha.FileName = "2024-01-10-alpha.md"
	alpha.SourcePath = "2024-01-10-alpha.md"
	root.AddChild(alpha)

	beta := tree.NewNode("Beta", "2023-12-01-beta.md", false)
	beta.FileName = "2023-12-01-beta.md"
	beta.SourcePath = "2023-12-01-beta.md"
	root.AddChild(beta)

	chapter1 := tree.NewNode("Chapter 1", "01-chapter.md", false)
	chapter1.FileName = "01-chapter.md"
	chapter1.SourcePath = "01-chapter.md"
	root.AddChild(chapter1)

	chapter2 := tree.NewNode("Chapter 2", "02-chapter.md", false)
	chapter2.FileName = "02-chapter.md"
	chapter2.SourcePath = "02-chapter.md"
	root.AddChild(chapter2)

	gamma := tree.NewNode("Gamma", "gamma.md", false)
	gamma.FileName = "gamma.md"
	gamma.SourcePath = "gamma.md"
	root.AddChild(gamma)

	guides := tree.NewNode("Guides", "guides", true)
	guides.SourcePath = "guides"
	root.AddChild(guides)

	items := BuildTopNavItems(root, true)
	if len(items) != 6 {
		t.Fatalf("BuildTopNavItems() length = %d, want 6", len(items))
	}

	gotNames := []string{items[0].Name, items[1].Name, items[2].Name, items[3].Name, items[4].Name, items[5].Name}
	wantNames := []string{"Alpha", "Beta", "Chapter 1", "Chapter 2", "Gamma", "Guides"}
	for i, want := range wantNames {
		if gotNames[i] != want {
			t.Errorf("items[%d].Name = %q, want %q", i, gotNames[i], want)
		}
	}
}

func TestBuildTopNavItemsLimits(t *testing.T) {
	root := tree.NewNode("", "", true)
	for i := 0; i < 9; i++ {
		node := tree.NewNode("Page", "page.md", false)
		node.FileName = "page.md"
		node.SourcePath = "page.md"
		root.AddChild(node)
	}

	if items := BuildTopNavItems(root, false); items != nil {
		t.Error("BuildTopNavItems() should return nil when top nav disabled")
	}

	if items := BuildTopNavItems(root, true); items != nil {
		t.Error("BuildTopNavItems() should return nil when more than 8 items")
	}
}

func TestTopNavNumberForSort(t *testing.T) {
	if got := topNavNumberForSort(nil); got != 999999 {
		t.Errorf("topNavNumberForSort(nil) = %d, want 999999", got)
	}
	val := 3
	if got := topNavNumberForSort(&val); got != 3 {
		t.Errorf("topNavNumberForSort(3) = %d, want 3", got)
	}
}

func TestRenderNavigationWithTopNav(t *testing.T) {
	root := tree.NewNode("", "", true)

	about := tree.NewNode("About", "about.md", false)
	about.FileName = "about.md"
	about.SourcePath = "about.md"
	root.AddChild(about)

	contact := tree.NewNode("Contact", "contact.md", false)
	contact.FileName = "contact.md"
	contact.SourcePath = "contact.md"
	root.AddChild(contact)

	guides := tree.NewNode("Guides", "guides", true)
	guides.SourcePath = "guides"
	intro := tree.NewNode("Intro", "guides/intro.md", false)
	intro.FileName = "intro.md"
	intro.SourcePath = "guides/intro.md"
	guides.AddChild(intro)
	root.AddChild(guides)

	topNavItems := BuildTopNavItems(root, true)
	html := RenderNavigationWithTopNav(root, "/guides/intro/", topNavItems)
	htmlStr := string(html)

	if strings.Contains(htmlStr, "About") || strings.Contains(htmlStr, "Contact") {
		t.Error("RenderNavigationWithTopNav() should filter top nav files from sidebar")
	}
	if !strings.Contains(htmlStr, "Guides") {
		t.Error("RenderNavigationWithTopNav() should include folders in sidebar")
	}
	if !strings.Contains(htmlStr, "Intro") {
		t.Error("RenderNavigationWithTopNav() should include folder children")
	}
}
