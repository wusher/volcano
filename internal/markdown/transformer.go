// Package markdown provides markdown parsing and content transformation.
package markdown

// ContentTransformer applies a series of transformations to HTML content.
// This consolidates the content enhancement pipeline into a single reusable component.
type ContentTransformer struct {
	siteURL string
}

// NewContentTransformer creates a new ContentTransformer with the given site URL.
// The site URL is used for external link processing.
func NewContentTransformer(siteURL string) *ContentTransformer {
	return &ContentTransformer{
		siteURL: siteURL,
	}
}

// Transform applies all content transformations to HTML content.
// This includes:
// - Adding heading anchors for linkable sections
// - Processing external links (adding target="_blank" and icons)
// - Wrapping code blocks with copy buttons
func (t *ContentTransformer) Transform(htmlContent string) string {
	// Add heading anchors
	htmlContent = AddHeadingAnchors(htmlContent)

	// Process external links
	htmlContent = ProcessExternalLinks(htmlContent, t.siteURL)

	// Wrap code blocks with copy button
	htmlContent = WrapCodeBlocks(htmlContent)

	return htmlContent
}

// TransformMarkdown processes raw markdown content through the full pipeline:
// 1. Process admonitions in markdown
// 2. Parse markdown to HTML
// 3. Apply HTML transformations
func (t *ContentTransformer) TransformMarkdown(mdContent []byte, sourceDir, sourcePath, outputPath, urlPath, fallbackTitle string) (*Page, error) {
	// Process admonitions before parsing
	mdContent = []byte(ProcessAdmonitions(string(mdContent)))

	// Parse the preprocessed content
	page, err := ParseContent(
		mdContent,
		sourcePath,
		outputPath,
		urlPath,
		sourceDir,
		fallbackTitle,
	)
	if err != nil {
		return nil, err
	}

	// Apply HTML transformations
	page.Content = t.Transform(page.Content)

	return page, nil
}
