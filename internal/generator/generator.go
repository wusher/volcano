// Package generator provides the static site generation engine.
package generator

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wusher/volcano/internal/assets"
	"github.com/wusher/volcano/internal/autoindex"
	"github.com/wusher/volcano/internal/content"
	"github.com/wusher/volcano/internal/instant"
	"github.com/wusher/volcano/internal/markdown"
	"github.com/wusher/volcano/internal/navigation"
	"github.com/wusher/volcano/internal/output"
	"github.com/wusher/volcano/internal/pwa"
	"github.com/wusher/volcano/internal/search"
	"github.com/wusher/volcano/internal/seo"
	"github.com/wusher/volcano/internal/styles"
	"github.com/wusher/volcano/internal/templates"
	"github.com/wusher/volcano/internal/toc"
	"github.com/wusher/volcano/internal/tree"
)

// Config holds configuration for the generator
type Config struct {
	InputDir         string
	OutputDir        string
	Title            string
	Clean            bool
	Quiet            bool
	Verbose          bool
	Colored          bool
	SiteURL          string // Base URL for canonical links
	Author           string // Site author
	OGImage          string // Path to local OG image file (copied to output)
	FaviconPath      string // Path to favicon file
	ShowLastMod      bool   // Show last modified date
	TopNav           bool   // Display root files in top navigation bar
	ShowPageNav      bool   // Show previous/next page navigation
	ShowBreadcrumbs  bool   // Show breadcrumb navigation
	Theme            string // Theme name (docs, blog, vanilla)
	CSSPath          string // Path to custom CSS file
	AccentColor      string // Custom accent color in hex format (e.g., "#ff6600")
	InstantNav       bool   // Enable instant navigation with hover prefetching
	ViewTransitions  bool   // Enable browser view transitions API
	InlineAssets     bool   // Embed CSS/JS inline instead of external files
	PWA              bool   // Enable PWA manifest and service worker generation
	Search           bool   // Enable search index generation
	AllowBrokenLinks bool   // Don't fail build on broken internal links
}

// Result holds the result of generation
type Result struct {
	PagesGenerated int
	Warnings       []string
}

// generatedPage tracks a page and its content for link validation
type generatedPage struct {
	urlPath     string
	sourceFile  string
	mdContent   string
	htmlContent string
}

// Generator handles static site generation
type Generator struct {
	config          Config
	renderer        *templates.Renderer
	transformer     *markdown.ContentTransformer
	logger          *output.Logger
	faviconLinks    template.HTML
	ogImageURL      string // Processed OG image URL (absolute if BaseURL provided)
	topNavItems     []templates.TopNavItem
	generatedPages  []generatedPage // Track pages for link validation
	baseURL         string          // Base URL path prefix extracted from SiteURL
	instantNavJS    template.JS     // Instant navigation JavaScript (if enabled)
	viewTransitions bool            // Enable browser view transitions API
	cssURL          string          // External CSS file URL (hashed)
	jsURL           string          // External JS file URL (hashed)
	css             string          // CSS content (for writing to file)
	pwaEnabled      bool            // Whether PWA support is enabled
	searchEnabled   bool            // Whether search is enabled
	searchIndex     *search.Index   // Search index data
}

// New creates a new Generator
func New(config Config, writer io.Writer) (*Generator, error) {
	// Get CSS content using the shared CSSLoader
	cssConfig := styles.CSSConfig{
		Theme:       config.Theme,
		CSSPath:     config.CSSPath,
		AccentColor: config.AccentColor,
	}
	cssLoader := styles.NewCSSLoader(cssConfig, os.ReadFile)
	css, err := cssLoader.LoadCSS()
	if err != nil {
		return nil, fmt.Errorf("failed to load CSS: %w", err)
	}

	// Create renderer (CSS will be passed per-page or via external file)
	renderer, err := templates.NewRenderer("")
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	// Extract base URL path from SiteURL for prefixing all links
	baseURL := tree.PrefixURL(config.SiteURL, "/")
	if baseURL == "/" {
		baseURL = ""
	} else {
		// Remove trailing slash from base URL path
		baseURL = baseURL[:len(baseURL)-1]
	}

	gen := &Generator{
		config:          config,
		renderer:        renderer,
		transformer:     markdown.NewContentTransformer(config.SiteURL),
		logger:          output.NewLogger(writer, config.Colored, config.Quiet, config.Verbose),
		baseURL:         baseURL,
		viewTransitions: config.ViewTransitions,
		css:             css,
		pwaEnabled:      config.PWA,
		searchEnabled:   config.Search,
	}

	// Initialize search index if enabled
	if config.Search {
		gen.searchIndex = &search.Index{Pages: []search.PageEntry{}}
	}

	// Initialize instant navigation JS if enabled
	if config.InstantNav {
		gen.instantNavJS = template.JS(instant.InstantNavJS)
	}

	return gen, nil
}

// Generate runs the full site generation
func (g *Generator) Generate() (*Result, error) {
	result := &Result{}

	// Print startup info
	g.logger.Println("Generating site...")
	g.logger.Println("  Input:  %s", g.config.InputDir)
	g.logger.Println("  Output: %s", g.config.OutputDir)
	g.logger.Println("  Title:  %s", g.config.Title)
	g.logger.Println("")

	// Step 1: Prepare output directory
	if err := g.prepareOutputDir(); err != nil {
		return nil, err
	}

	// Process favicon if configured
	if g.config.FaviconPath != "" {
		faviconConfig := assets.FaviconConfig{IconPath: g.config.FaviconPath}
		links, err := assets.ProcessFavicon(faviconConfig, g.config.OutputDir)
		if err != nil {
			g.logger.Warning("Failed to process favicon: %v", err)
		} else {
			g.faviconLinks = assets.RenderFaviconLinks(links)
		}
	}

	// Process OG image if configured
	if g.config.OGImage != "" {
		ogConfig := assets.OGImageConfig{
			ImagePath: g.config.OGImage,
			BaseURL:   g.config.SiteURL,
		}
		ogURL, err := assets.ProcessOGImage(ogConfig, g.config.OutputDir)
		if err != nil {
			g.logger.Warning("Failed to process og-image: %v", err)
		} else {
			g.ogImageURL = ogURL
		}
	}

	// Step 2: Scan input directory
	g.logger.Println("Scanning input directory...")
	site, err := tree.Scan(g.config.InputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to scan input directory: %w", err)
	}

	if len(site.AllPages) == 0 {
		g.logger.Warning("No markdown files found in %s", g.config.InputDir)
		result.Warnings = append(result.Warnings, "No markdown files found")
		return result, nil
	}

	// Write hashed CSS and JS assets (after confirming we have pages)
	// Skip when inline assets are requested
	if !g.config.InlineAssets {
		if err := g.writeHashedAssets(); err != nil {
			return nil, err
		}
	}

	// Count folders
	folderCount := countFolders(site.Root)
	g.logger.Println("Found %d markdown files in %d folders", len(site.AllPages), folderCount)
	g.logger.Println("")

	// Build top nav items if enabled (with base URL prefixing)
	g.topNavItems = templates.BuildTopNavItemsWithBaseURL(site.Root, g.config.TopNav, g.config.SiteURL)
	if len(g.topNavItems) > 0 {
		g.logger.Verbose("Using top navigation bar with %d items", len(g.topNavItems))
	}

	// Step 3: Generate pages
	g.logger.Println("Generating pages...")
	for _, node := range site.AllPages {
		if err := g.generatePage(node, site.Root, site.AllPages); err != nil {
			return nil, fmt.Errorf("failed to generate %s: %w", node.Path, err)
		}
		result.PagesGenerated++
		g.logger.FileSuccess(node.Path)
	}

	// Step 4: Generate auto-index pages for folders without index.md
	foldersNeedingIndex := autoindex.CollectFoldersNeedingAutoIndex(site.Root)
	if len(foldersNeedingIndex) > 0 {
		g.logger.Verbose("Generating auto-index pages for %d folders...", len(foldersNeedingIndex))
		for _, folder := range foldersNeedingIndex {
			if err := g.generateAutoIndex(folder, site.Root); err != nil {
				return nil, fmt.Errorf("failed to generate auto-index for %s: %w", folder.Path, err)
			}
			g.logger.Verbose("  Auto-indexed: %s", folder.Path)
		}
	}

	// Step 5: Generate 404 page
	if err := g.generate404(site.Root); err != nil {
		return nil, fmt.Errorf("failed to generate 404 page: %w", err)
	}

	// Step 6: Verify all navigation links resolve
	g.logger.Verbose("Verifying navigation links...")
	brokenLinks := g.verifyLinks(site.AllPages)
	if len(brokenLinks) > 0 {
		g.logger.Println("")
		if g.config.AllowBrokenLinks {
			g.logger.Warning("Found %d broken navigation links (continuing due to --allow-broken-links):", len(brokenLinks))
		} else {
			g.logger.Error("Found %d broken navigation links:", len(brokenLinks))
		}
		for _, link := range brokenLinks {
			if g.config.AllowBrokenLinks {
				g.logger.Warning("  %s", link)
			} else {
				g.logger.Error("  %s", link)
			}
		}
		if !g.config.AllowBrokenLinks {
			return nil, fmt.Errorf("build failed: %d broken navigation links found", len(brokenLinks))
		}
	}

	// Step 7: Verify all internal links in content resolve
	g.logger.Verbose("Verifying internal links in content...")
	validURLs := tree.BuildValidURLMapWithAutoIndex(site.AllPages, foldersNeedingIndex, g.config.SiteURL)
	brokenContentLinks := g.verifyContentLinks(validURLs)
	if len(brokenContentLinks) > 0 {
		g.logger.Println("")
		if g.config.AllowBrokenLinks {
			g.logger.Warning("Found %d broken internal link(s) (continuing due to --allow-broken-links):", len(brokenContentLinks))
		} else {
			g.logger.Error("Found %d broken internal link(s):", len(brokenContentLinks))
		}
		for i, bl := range brokenContentLinks {
			if g.config.AllowBrokenLinks {
				g.logger.Warning("")
				g.logger.Warning("Link #%d:", i+1)
			} else {
				g.logger.Error("")
				g.logger.Error("Link #%d:", i+1)
			}
			if bl.SourceFile != "" {
				if bl.LineNumber > 0 {
					if g.config.AllowBrokenLinks {
						g.logger.Warning("  File: %s:%d", bl.SourceFile, bl.LineNumber)
					} else {
						g.logger.Error("  File: %s:%d", bl.SourceFile, bl.LineNumber)
					}
				} else {
					if g.config.AllowBrokenLinks {
						g.logger.Warning("  File: %s", bl.SourceFile)
					} else {
						g.logger.Error("  File: %s", bl.SourceFile)
					}
				}
			}
			if bl.OriginalSyntax != "" {
				if g.config.AllowBrokenLinks {
					g.logger.Warning("  Syntax: %s", bl.OriginalSyntax)
				} else {
					g.logger.Error("  Syntax: %s", bl.OriginalSyntax)
				}
			}
			if bl.LinkText != "" && bl.LinkText != bl.LinkURL {
				if g.config.AllowBrokenLinks {
					g.logger.Warning("  Text: %s", bl.LinkText)
				} else {
					g.logger.Error("  Text: %s", bl.LinkText)
				}
			}
			if g.config.AllowBrokenLinks {
				g.logger.Warning("  Broken URL: %s", bl.LinkURL)
			} else {
				g.logger.Error("  Broken URL: %s", bl.LinkURL)
			}
			if len(bl.Suggestions) > 0 {
				if g.config.AllowBrokenLinks {
					g.logger.Warning("  Suggestions:")
					for _, suggestion := range bl.Suggestions {
						g.logger.Warning("    - %s", suggestion)
					}
				} else {
					g.logger.Error("  Suggestions:")
					for _, suggestion := range bl.Suggestions {
						g.logger.Error("    - %s", suggestion)
					}
				}
			}
		}
		if !g.config.AllowBrokenLinks {
			return nil, fmt.Errorf("build failed: %d broken internal links found", len(brokenContentLinks))
		}
	}

	// Step 8: Generate PWA assets if enabled
	if g.pwaEnabled {
		if err := g.generatePWA(site.AllPages, foldersNeedingIndex); err != nil {
			return nil, fmt.Errorf("failed to generate PWA assets: %w", err)
		}
	}

	// Step 9: Generate search assets if enabled
	if g.searchEnabled && g.searchIndex != nil {
		if err := search.GenerateSearchIndex(g.config.OutputDir, g.searchIndex); err != nil {
			return nil, fmt.Errorf("failed to generate search index: %w", err)
		}
		g.logger.Verbose("  search-index.json")

		// Write search.js
		searchJS := search.GenerateSearchJS(g.baseURL)
		searchJSPath := filepath.Join(g.config.OutputDir, "search.js")
		if err := os.WriteFile(searchJSPath, []byte(searchJS), 0644); err != nil {
			return nil, fmt.Errorf("failed to write search.js: %w", err)
		}
		g.logger.Verbose("  search.js")
	}

	// Print summary
	g.logger.Println("")
	g.logger.Success("Generated %d pages in %s", result.PagesGenerated, g.config.OutputDir)

	return result, nil
}

// countFolders counts the number of folders in the tree
func countFolders(node *tree.Node) int {
	if node == nil {
		return 0
	}
	count := 0
	if node.IsFolder {
		count = 1
	}
	for _, child := range node.Children {
		count += countFolders(child)
	}
	return count
}

// verifyLinks checks that all generated pages have corresponding output files
func (g *Generator) verifyLinks(allPages []*tree.Node) []string {
	var broken []string
	for _, node := range allPages {
		outputPath := tree.GetOutputPath(node)
		if outputPath == "" {
			continue
		}
		fullPath := filepath.Join(g.config.OutputDir, outputPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			urlPath := tree.GetURLPath(node)
			broken = append(broken, fmt.Sprintf("%s -> %s (expected: %s)", node.Path, urlPath, outputPath))
		}
	}
	return broken
}

// verifyContentLinks checks all internal links in generated page content
func (g *Generator) verifyContentLinks(validURLs map[string]bool) []markdown.BrokenLink {
	var allBroken []markdown.BrokenLink

	for _, page := range g.generatedPages {
		broken := markdown.ValidateLinksWithSource(page.htmlContent, page.urlPath, page.sourceFile, page.mdContent, validURLs)
		allBroken = append(allBroken, broken...)
	}

	return allBroken
}

// prepareOutputDir creates or cleans the output directory
func (g *Generator) prepareOutputDir() error {
	if g.config.Clean {
		g.logger.Verbose("Cleaning output directory...")
		if err := os.RemoveAll(g.config.OutputDir); err != nil {
			return fmt.Errorf("failed to clean output directory: %w", err)
		}
	}

	if err := os.MkdirAll(g.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return nil
}

// generatePage generates a single page
func (g *Generator) generatePage(node *tree.Node, root *tree.Node, allPages []*tree.Node) error {
	// Get paths
	outputPath := tree.GetOutputPath(node)
	urlPath := tree.GetURLPath(node)
	fullOutputPath := filepath.Join(g.config.OutputDir, outputPath)

	// Read markdown content
	mdContent, err := os.ReadFile(node.SourcePath)
	if err != nil {
		return err
	}

	// Compute source directory for wikilink resolution
	// e.g., "guides/customizing-appearance.md" -> "/guides/"
	relDir := filepath.Dir(node.Path)
	sourceDir := "/"
	if relDir != "." && relDir != "" {
		sourceDir = "/" + tree.SlugifyPath(relDir) + "/"
	}

	// Transform markdown to HTML with all enhancements
	page, err := g.transformer.TransformMarkdown(
		mdContent,
		sourceDir,
		node.SourcePath,
		outputPath,
		urlPath,
		node.Name, // fallback title
	)
	if err != nil {
		return err
	}

	htmlContent := page.Content

	// Calculate reading time
	rt := content.CalculateReadingTime(htmlContent)
	readingTime := content.FormatReadingTime(rt)

	// Get last modified date if enabled
	var lastModified string
	if g.config.ShowLastMod {
		mod := content.GetLastModified(node.SourcePath)
		lastModified = content.FormatLastModified(mod, false) // Use absolute format
	}

	// Build breadcrumbs (with base URL prefixing) - only if enabled
	var breadcrumbsHTML template.HTML
	if g.config.ShowBreadcrumbs {
		breadcrumbs := navigation.BuildBreadcrumbsWithBaseURL(node, g.config.Title, g.config.SiteURL)
		breadcrumbsHTML = navigation.RenderBreadcrumbs(breadcrumbs)
	}

	// Build page navigation (only if enabled, with base URL prefixing)
	var pageNavHTML template.HTML
	if g.config.ShowPageNav {
		pageNav := navigation.BuildPageNavigationWithBaseURL(node, allPages, g.config.SiteURL)
		pageNavHTML = navigation.RenderPageNavigation(pageNav)
	}

	// Extract TOC
	pageTOC := toc.ExtractTOC(htmlContent, 3)
	tocHTML := toc.RenderTOC(pageTOC)
	hasTOC := pageTOC != nil && len(pageTOC.Items) > 0

	// Generate SEO meta tags
	seoConfig := seo.Config{
		SiteURL:   g.config.SiteURL,
		SiteTitle: g.config.Title,
		Author:    g.config.Author,
		OGImage:   g.ogImageURL, // Use processed URL, not raw path
	}
	pageMeta := seo.GeneratePageMeta(page.Title, htmlContent, urlPath, seoConfig)
	metaTagsHTML := seo.RenderMetaTags(pageMeta)

	// Render navigation (filtered when top nav is enabled, with base URL prefixing)
	nav := templates.RenderNavigationWithTopNavAndBaseURL(root, urlPath, g.topNavItems, g.config.SiteURL)

	// Prepare template data
	data := templates.PageData{
		SiteTitle:       g.config.Title,
		PageTitle:       page.Title,
		Content:         template.HTML(htmlContent),
		Navigation:      nav,
		CurrentPath:     urlPath,
		Breadcrumbs:     breadcrumbsHTML,
		PageNav:         pageNavHTML,
		TOC:             tocHTML,
		MetaTags:        metaTagsHTML,
		FaviconLinks:    g.faviconLinks,
		ReadingTime:     readingTime,
		LastModified:    lastModified,
		HasTOC:          hasTOC,
		ShowSearch:      true,
		TopNavItems:     g.topNavItems,
		BaseURL:         g.baseURL,
		CSSURL:          g.cssURL,
		JSURL:           g.jsURL,
		CSS:             g.inlineCSS(),
		InstantNavJS:    g.instantNavJS,
		ViewTransitions: g.viewTransitions,
		PWAEnabled:      g.pwaEnabled,
		SearchEnabled:   g.searchEnabled,
	}

	// Create output directory
	outputDir := filepath.Dir(fullOutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	// Write file
	f, err := os.Create(fullOutputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fullOutputPath, err)
	}
	defer func() { _ = f.Close() }()

	if err := g.renderer.Render(f, data); err != nil {
		return fmt.Errorf("failed to render page: %w", err)
	}

	// Track page for link validation
	g.generatedPages = append(g.generatedPages, generatedPage{
		urlPath:     urlPath,
		sourceFile:  node.SourcePath,
		mdContent:   string(mdContent),
		htmlContent: htmlContent,
	})

	// Collect search index data if enabled
	if g.searchEnabled && g.searchIndex != nil {
		entry := search.PageEntry{
			Title:    page.Title,
			URL:      urlPath,
			Headings: search.ExtractHeadings(htmlContent),
		}
		g.searchIndex.Pages = append(g.searchIndex.Pages, entry)
	}

	return nil
}

// generate404 generates the 404 error page
func (g *Generator) generate404(root *tree.Node) error {
	// Build home URL with base URL prefix
	homeURL := "/"
	if g.baseURL != "" {
		homeURL = g.baseURL + "/"
	}
	content := fmt.Sprintf(`<h1>404 - Page Not Found</h1>
<p>The page you're looking for doesn't exist.</p>
<p><a href="%s">Return to home</a></p>`, homeURL)

	nav := templates.RenderNavigationWithBaseURL(root, "", g.config.SiteURL)

	data := templates.PageData{
		SiteTitle:       g.config.Title,
		PageTitle:       "Page Not Found",
		Content:         template.HTML(content),
		Navigation:      nav,
		CurrentPath:     "",
		BaseURL:         g.baseURL,
		CSSURL:          g.cssURL,
		JSURL:           g.jsURL,
		CSS:             g.inlineCSS(),
		InstantNavJS:    g.instantNavJS,
		ViewTransitions: g.viewTransitions,
		PWAEnabled:      g.pwaEnabled,
	}

	fullPath := filepath.Join(g.config.OutputDir, "404.html")
	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create 404.html: %w", err)
	}
	defer func() { _ = f.Close() }()

	return g.renderer.Render(f, data)
}

// writeHashedAssets writes CSS and JS to separate files with content hashes.
// This enables browser caching with automatic cache invalidation on content changes.
func (g *Generator) writeHashedAssets() error {
	// Write CSS file
	cssAsset, err := assets.WriteHashedAsset(g.config.OutputDir, "styles", "css", g.css, g.baseURL)
	if err != nil {
		return fmt.Errorf("failed to write CSS: %w", err)
	}
	g.cssURL = cssAsset.URLPath
	g.logger.Verbose("  CSS: %s", cssAsset.FileName)

	// Write JS file if instant navigation is enabled
	if g.config.InstantNav {
		jsAsset, err := assets.WriteHashedAsset(g.config.OutputDir, "app", "js", string(g.instantNavJS), g.baseURL)
		if err != nil {
			return fmt.Errorf("failed to write JS: %w", err)
		}
		g.jsURL = jsAsset.URLPath
		g.logger.Verbose("  JS: %s", jsAsset.FileName)
	}

	return nil
}

// inlineCSS returns CSS content for inline embedding when InlineAssets is enabled.
// Returns empty string when external assets are used (CSSURL will be set instead).
func (g *Generator) inlineCSS() template.CSS {
	if g.config.InlineAssets {
		return template.CSS(g.css)
	}
	return ""
}

// generatePWA creates PWA assets (manifest.json, sw.js, icons).
func (g *Generator) generatePWA(allPages []*tree.Node, autoIndexFolders []*tree.Node) error {
	g.logger.Verbose("Generating PWA assets...")

	// 1. Generate icons from favicon
	iconResult, err := pwa.GenerateIcons(g.config.FaviconPath, g.config.OutputDir)
	if err != nil {
		return err
	}
	if iconResult.Warning != "" {
		g.logger.Warning(iconResult.Warning)
	}

	// 2. Generate manifest.json
	manifestConfig := pwa.ManifestConfig{
		SiteTitle:  g.config.Title,
		ThemeColor: g.config.AccentColor, // Empty = default in manifest.go
		BaseURL:    g.baseURL,
		HasIcons:   iconResult.Generated,
	}
	if err := pwa.GenerateManifest(g.config.OutputDir, manifestConfig); err != nil {
		return err
	}

	// 3. Collect all page URLs
	pageURLs := collectPageURLs(allPages, autoIndexFolders, g.baseURL)

	// 4. Collect asset URLs
	assetURLs := []string{}
	if g.cssURL != "" {
		assetURLs = append(assetURLs, g.cssURL)
	}
	if g.jsURL != "" {
		assetURLs = append(assetURLs, g.jsURL)
	}
	if iconResult.Generated {
		assetURLs = append(assetURLs, pwa.GetIconURLs(g.baseURL)...)
	}
	// Add search files if search is enabled
	if g.searchEnabled {
		searchBase := g.baseURL
		if searchBase != "" && !strings.HasSuffix(searchBase, "/") {
			searchBase += "/"
		}
		if searchBase == "" {
			searchBase = "/"
		}
		assetURLs = append(assetURLs, searchBase+"search.js", searchBase+"search-index.json")
	}
	// Add favicon if configured
	if g.config.FaviconPath != "" {
		faviconName := filepath.Base(g.config.FaviconPath)
		faviconURL := g.baseURL
		if faviconURL != "" && !strings.HasSuffix(faviconURL, "/") {
			faviconURL += "/"
		}
		if faviconURL == "" {
			faviconURL = "/"
		}
		assetURLs = append(assetURLs, faviconURL+faviconName)
	}

	// 5. Generate service worker
	swConfig := pwa.ServiceWorkerConfig{
		BaseURL:   g.baseURL,
		PageURLs:  pageURLs,
		AssetURLs: assetURLs,
	}
	if err := pwa.GenerateServiceWorker(g.config.OutputDir, swConfig); err != nil {
		return err
	}

	g.logger.Verbose("  manifest.json")
	g.logger.Verbose("  sw.js")
	if iconResult.Generated {
		g.logger.Verbose("  icon-192.png, icon-512.png")
	}

	return nil
}

// collectPageURLs gathers all page URLs for service worker precaching.
func collectPageURLs(allPages []*tree.Node, autoIndexFolders []*tree.Node, baseURL string) []string {
	// Use a map to deduplicate URLs
	urlSet := make(map[string]bool)

	// Add all pages
	for _, page := range allPages {
		urlPath := tree.GetURLPath(page)
		if baseURL != "" {
			urlPath = baseURL + urlPath
		}
		urlSet[urlPath] = true
	}

	// Add auto-index folders
	for _, folder := range autoIndexFolders {
		urlPath := "/" + tree.SlugifyPath(folder.Path) + "/"
		if baseURL != "" {
			urlPath = baseURL + urlPath
		}
		urlSet[urlPath] = true
	}

	// Convert map to slice
	urls := make([]string, 0, len(urlSet))
	for url := range urlSet {
		urls = append(urls, url)
	}

	return urls
}
