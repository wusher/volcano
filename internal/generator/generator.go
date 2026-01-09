// Package generator provides the static site generation engine.
package generator

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"volcano/internal/assets"
	"volcano/internal/content"
	"volcano/internal/markdown"
	"volcano/internal/navigation"
	"volcano/internal/output"
	"volcano/internal/seo"
	"volcano/internal/styles"
	"volcano/internal/templates"
	"volcano/internal/toc"
	"volcano/internal/tree"
)

// Config holds configuration for the generator
type Config struct {
	InputDir    string
	OutputDir   string
	Title       string
	Clean       bool
	Quiet       bool
	Verbose     bool
	Colored     bool
	SiteURL     string // Base URL for canonical links
	Author      string // Site author
	OGImage     string // Default Open Graph image
	FaviconPath string // Path to favicon file
	ShowLastMod bool   // Show last modified date
	TopNav      bool   // Display root files in top navigation bar
}

// Result holds the result of generation
type Result struct {
	PagesGenerated int
	Warnings       []string
}

// Generator handles static site generation
type Generator struct {
	config       Config
	renderer     *templates.Renderer
	parser       *markdown.Parser
	logger       *output.Logger
	faviconLinks template.HTML
	topNavItems  []templates.TopNavItem
}

// New creates a new Generator
func New(config Config, writer io.Writer) (*Generator, error) {
	renderer, err := templates.NewRenderer(styles.GetCSS())
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	return &Generator{
		config:   config,
		renderer: renderer,
		parser:   markdown.NewParser(),
		logger:   output.NewLogger(writer, config.Colored, config.Quiet, config.Verbose),
	}, nil
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

	// Count folders
	folderCount := countFolders(site.Root)
	g.logger.Println("Found %d markdown files in %d folders", len(site.AllPages), folderCount)
	g.logger.Println("")

	// Build top nav items if enabled
	g.topNavItems = templates.BuildTopNavItems(site.Root, g.config.TopNav)
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
	foldersNeedingIndex := collectFoldersNeedingAutoIndex(site.Root)
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

	// Process admonitions before parsing
	mdContent = []byte(markdown.ProcessAdmonitions(string(mdContent)))

	// Parse the preprocessed content
	page, err := markdown.ParseContent(
		mdContent,
		node.SourcePath,
		outputPath,
		urlPath,
		node.Name, // fallback title
	)
	if err != nil {
		return err
	}

	// Process content enhancements
	htmlContent := page.Content

	// Add heading anchors
	htmlContent = markdown.AddHeadingAnchors(htmlContent)

	// Process external links
	htmlContent = markdown.ProcessExternalLinks(htmlContent, g.config.SiteURL)

	// Wrap code blocks with copy button
	htmlContent = markdown.WrapCodeBlocks(htmlContent)

	// Calculate reading time
	rt := content.CalculateReadingTime(htmlContent)
	readingTime := content.FormatReadingTime(rt)

	// Get last modified date if enabled
	var lastModified string
	if g.config.ShowLastMod {
		mod := content.GetLastModified(node.SourcePath)
		lastModified = content.FormatLastModified(mod, false) // Use absolute format
	}

	// Build breadcrumbs
	breadcrumbs := navigation.BuildBreadcrumbs(node, g.config.Title)
	breadcrumbsHTML := navigation.RenderBreadcrumbs(breadcrumbs)

	// Build page navigation
	pageNav := navigation.BuildPageNavigation(node, allPages)
	pageNavHTML := navigation.RenderPageNavigation(pageNav)

	// Extract TOC
	pageTOC := toc.ExtractTOC(htmlContent, 3)
	tocHTML := toc.RenderTOC(pageTOC)
	hasTOC := pageTOC != nil && len(pageTOC.Items) > 0

	// Generate SEO meta tags
	seoConfig := seo.Config{
		SiteURL:   g.config.SiteURL,
		SiteTitle: g.config.Title,
		Author:    g.config.Author,
		OGImage:   g.config.OGImage,
	}
	pageMeta := seo.GeneratePageMeta(page.Title, htmlContent, urlPath, seoConfig)
	metaTagsHTML := seo.RenderMetaTags(pageMeta)

	// Render navigation (filtered when top nav is enabled)
	nav := templates.RenderNavigationWithTopNav(root, urlPath, g.topNavItems)

	// Prepare template data
	data := templates.PageData{
		SiteTitle:    g.config.Title,
		PageTitle:    page.Title,
		Content:      template.HTML(htmlContent),
		Navigation:   nav,
		CurrentPath:  urlPath,
		Breadcrumbs:  breadcrumbsHTML,
		PageNav:      pageNavHTML,
		TOC:          tocHTML,
		MetaTags:     metaTagsHTML,
		FaviconLinks: g.faviconLinks,
		ReadingTime:  readingTime,
		LastModified: lastModified,
		HasTOC:       hasTOC,
		ShowSearch:   true,
		TopNavItems:  g.topNavItems,
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

	return nil
}

// generate404 generates the 404 error page
func (g *Generator) generate404(root *tree.Node) error {
	content := `<h1>404 - Page Not Found</h1>
<p>The page you're looking for doesn't exist.</p>
<p><a href="/">Return to home</a></p>`

	nav := templates.RenderNavigation(root, "")

	data := templates.PageData{
		SiteTitle:   g.config.Title,
		PageTitle:   "Page Not Found",
		Content:     template.HTML(content),
		Navigation:  nav,
		CurrentPath: "",
	}

	fullPath := filepath.Join(g.config.OutputDir, "404.html")
	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create 404.html: %w", err)
	}
	defer func() { _ = f.Close() }()

	return g.renderer.Render(f, data)
}
