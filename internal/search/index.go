// Package search provides search index generation for static sites.
package search

// Index represents the complete search index for the site.
type Index struct {
	Pages []PageEntry `json:"pages"`
}

// PageEntry represents a single page in the search index.
type PageEntry struct {
	Title    string         `json:"title"`
	URL      string         `json:"url"`
	Headings []HeadingEntry `json:"headings,omitempty"`
}

// HeadingEntry represents a heading within a page.
type HeadingEntry struct {
	Text   string `json:"text"`
	Anchor string `json:"anchor"`
	Level  int    `json:"level"`
}
