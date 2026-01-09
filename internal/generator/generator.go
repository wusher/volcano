// Package generator provides the static site generation engine.
package generator

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"volcano/internal/markdown"
	"volcano/internal/styles"
	"volcano/internal/templates"
	"volcano/internal/tree"
)

// Config holds configuration for the generator
type Config struct {
	InputDir  string
	OutputDir string
	Title     string
	Clean     bool
	Quiet     bool
	Verbose   bool
}

// Result holds the result of generation
type Result struct {
	PagesGenerated int
	Warnings       []string
}

// Generator handles static site generation
type Generator struct {
	config   Config
	renderer *templates.Renderer
	parser   *markdown.Parser
	writer   io.Writer
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
		writer:   writer,
	}, nil
}

// Generate runs the full site generation
func (g *Generator) Generate() (*Result, error) {
	result := &Result{}

	// Print startup info
	g.log("Generating site...")
	g.log("  Input:  %s", g.config.InputDir)
	g.log("  Output: %s", g.config.OutputDir)
	g.log("  Title:  %s", g.config.Title)
	g.log("")

	// Step 1: Prepare output directory
	if err := g.prepareOutputDir(); err != nil {
		return nil, err
	}

	// Step 2: Scan input directory
	g.log("Scanning input directory...")
	site, err := tree.Scan(g.config.InputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to scan input directory: %w", err)
	}

	if len(site.AllPages) == 0 {
		result.Warnings = append(result.Warnings, "No markdown files found")
		return result, nil
	}

	g.log("Found %d markdown files", len(site.AllPages))

	// Step 3: Generate pages
	g.log("Generating pages...")
	for _, node := range site.AllPages {
		if err := g.generatePage(node, site.Root); err != nil {
			return nil, fmt.Errorf("failed to generate %s: %w", node.Path, err)
		}
		result.PagesGenerated++
		g.logVerbose("  ✓ %s", node.Path)
	}

	// Step 4: Generate 404 page
	if err := g.generate404(site.Root); err != nil {
		return nil, fmt.Errorf("failed to generate 404 page: %w", err)
	}

	// Print summary
	g.log("")
	g.log("Generated %d pages in %s", result.PagesGenerated, g.config.OutputDir)

	return result, nil
}

// prepareOutputDir creates or cleans the output directory
func (g *Generator) prepareOutputDir() error {
	if g.config.Clean {
		g.logVerbose("Cleaning output directory...")
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
func (g *Generator) generatePage(node *tree.Node, root *tree.Node) error {
	// Get paths
	outputPath := tree.GetOutputPath(node)
	urlPath := tree.GetURLPath(node)
	fullOutputPath := filepath.Join(g.config.OutputDir, outputPath)

	// Parse markdown file
	page, err := markdown.ParseFile(
		node.SourcePath,
		outputPath,
		urlPath,
		node.Name, // fallback title
	)
	if err != nil {
		return err
	}

	// Render navigation
	nav := templates.RenderNavigation(root, urlPath)

	// Prepare template data
	data := templates.PageData{
		SiteTitle:   g.config.Title,
		PageTitle:   page.Title,
		Content:     template.HTML(page.Content),
		Navigation:  nav,
		CurrentPath: urlPath,
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

// log prints a message if not in quiet mode
func (g *Generator) log(format string, args ...interface{}) {
	if !g.config.Quiet {
		_, _ = fmt.Fprintf(g.writer, format+"\n", args...)
	}
}

// logVerbose prints a message only in verbose mode
func (g *Generator) logVerbose(format string, args ...interface{}) {
	if g.config.Verbose {
		_, _ = fmt.Fprintf(g.writer, format+"\n", args...)
	}
}
