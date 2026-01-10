package markdown

import (
	"bytes"
)

// StripFrontMatter removes YAML front matter from the beginning of markdown content.
// Front matter is delimited by --- at the start of the file and closed by another ---.
// Returns the content with front matter removed.
func StripFrontMatter(content []byte) []byte {
	// Must start with ---
	if !bytes.HasPrefix(content, []byte("---")) {
		return content
	}

	// Find the closing ---
	// Skip the first 3 characters (the opening ---)
	rest := content[3:]

	// Skip optional newline after opening ---
	if len(rest) > 0 && rest[0] == '\n' {
		rest = rest[1:]
	} else if len(rest) > 1 && rest[0] == '\r' && rest[1] == '\n' {
		rest = rest[2:]
	}

	// Find the closing --- on its own line
	idx := bytes.Index(rest, []byte("\n---"))
	if idx == -1 {
		// Try Windows line endings
		idx = bytes.Index(rest, []byte("\r\n---"))
		if idx == -1 {
			// No closing ---, return original content
			return content
		}
		// Skip past \r\n---
		afterClose := rest[idx+5:]
		return stripLeadingNewlines(afterClose)
	}

	// Skip past the closing \n---
	afterClose := rest[idx+4:]
	return stripLeadingNewlines(afterClose)
}

// stripLeadingNewlines removes leading newlines (Unix or Windows style)
func stripLeadingNewlines(content []byte) []byte {
	for len(content) > 0 {
		if content[0] == '\n' {
			content = content[1:]
		} else if len(content) > 1 && content[0] == '\r' && content[1] == '\n' {
			content = content[2:]
		} else {
			break
		}
	}
	return content
}
