package generator

import (
	"os"
	"path/filepath"

	"github.com/wusher/volcano/internal/autoindex"
	"github.com/wusher/volcano/internal/navigation"
	"github.com/wusher/volcano/internal/seo"
	"github.com/wusher/volcano/internal/templates"
	"github.com/wusher/volcano/internal/tree"
)

// generateAutoIndex generates an auto-index page for a folder without an index.md
func (g *Generator) generateAutoIndex(node *tree.Node, root *tree.Node) error {
	index := autoindex.BuildWithBaseURL(node, g.config.SiteURL)
	fullOutputPath := filepath.Join(g.config.OutputDir, index.OutputPath)

	// Build content
	htmlContent := autoindex.RenderContent(index)

	// Build breadcrumbs (with base URL prefixing)
	breadcrumbs := navigation.BuildBreadcrumbsWithBaseURL(node, g.config.Title, g.config.SiteURL)
	breadcrumbsHTML := navigation.RenderBreadcrumbs(breadcrumbs)

	// Generate SEO meta tags
	seoConfig := seo.Config{
		SiteURL:   g.config.SiteURL,
		SiteTitle: g.config.Title,
		Author:    g.config.Author,
		OGImage:   g.ogImageURL, // Use processed URL, not raw path
	}
	pageMeta := seo.GeneratePageMeta(index.Title, string(htmlContent), index.URLPath, seoConfig)
	metaTagsHTML := seo.RenderMetaTags(pageMeta)

	// Render navigation (with base URL prefixing)
	nav := templates.RenderNavigationWithBaseURL(root, index.URLPath, g.config.SiteURL)

	// Prepare template data
	data := templates.PageData{
		SiteTitle:       g.config.Title,
		PageTitle:       index.Title,
		Content:         htmlContent,
		Navigation:      nav,
		CurrentPath:     index.URLPath,
		Breadcrumbs:     breadcrumbsHTML,
		MetaTags:        metaTagsHTML,
		ShowSearch:      true,
		BaseURL:         g.baseURL,
		InstantNavJS:    g.instantNavJS,
		ViewTransitions: g.viewTransitions,
	}

	// Create output directory
	outputDir := filepath.Dir(fullOutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Write file
	f, err := os.Create(fullOutputPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	return g.renderer.Render(f, data)
}
