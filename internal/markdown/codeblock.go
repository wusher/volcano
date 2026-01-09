package markdown

import (
	"regexp"
	"strconv"
	"strings"
)

// preCodeRegex matches <pre><code> blocks
var preCodeRegex = regexp.MustCompile(`(?s)<pre([^>]*)><code([^>]*)>(.*?)</code></pre>`)

// WrapCodeBlocks adds copy button wrapper to code blocks
func WrapCodeBlocks(htmlContent string) string {
	result := preCodeRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		matches := preCodeRegex.FindStringSubmatch(match)
		if len(matches) < 4 {
			return match
		}

		preAttrs := matches[1]
		codeAttrs := matches[2]
		code := matches[3]

		var sb strings.Builder
		sb.WriteString(`<div class="code-block">`)
		sb.WriteString("\n")
		sb.WriteString(`  <button class="copy-button" aria-label="Copy code to clipboard">`)
		sb.WriteString("\n")
		sb.WriteString(`    <svg class="copy-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>`)
		sb.WriteString("\n")
		sb.WriteString(`    <svg class="check-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>`)
		sb.WriteString("\n")
		sb.WriteString(`    <span class="copy-text">Copy</span>`)
		sb.WriteString("\n")
		sb.WriteString(`  </button>`)
		sb.WriteString("\n")
		sb.WriteString(`  <pre`)
		sb.WriteString(preAttrs)
		sb.WriteString(`><code`)
		sb.WriteString(codeAttrs)
		sb.WriteString(`>`)
		sb.WriteString(code)
		sb.WriteString(`</code></pre>`)
		sb.WriteString("\n")
		sb.WriteString(`</div>`)

		return sb.String()
	})

	return result
}

// LineSpec represents line highlighting specification
type LineSpec struct {
	Lines []int // Individual lines to highlight
}

// ParseLineSpec parses a line specification string like "3,5-7,10"
func ParseLineSpec(spec string) LineSpec {
	var lines []int

	if spec == "" {
		return LineSpec{Lines: lines}
	}

	parts := strings.Split(spec, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Check for range
		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			if len(rangeParts) == 2 {
				start, err1 := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
				end, err2 := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
				if err1 == nil && err2 == nil && start <= end {
					for i := start; i <= end; i++ {
						lines = append(lines, i)
					}
				}
			}
		} else {
			// Single line
			if n, err := strconv.Atoi(part); err == nil {
				lines = append(lines, n)
			}
		}
	}

	return LineSpec{Lines: lines}
}

// highlightInfoRegex matches language info with highlight spec like "go {3,5-7}"
var highlightInfoRegex = regexp.MustCompile(`^(\w+)\s*\{([^}]+)\}$`)

// ParseCodeBlockInfo parses the info string from a fenced code block
// Returns language and highlight spec
func ParseCodeBlockInfo(info string) (string, string) {
	info = strings.TrimSpace(info)
	if info == "" {
		return "", ""
	}

	matches := highlightInfoRegex.FindStringSubmatch(info)
	if len(matches) == 3 {
		return matches[1], matches[2]
	}

	// No highlight spec, just language
	return info, ""
}

// ApplyLineHighlighting wraps lines in spans with highlight class where specified
func ApplyLineHighlighting(code string, highlightLines []int) string {
	if len(highlightLines) == 0 {
		return code
	}

	// Create a set of lines to highlight
	highlightSet := make(map[int]bool)
	for _, line := range highlightLines {
		highlightSet[line] = true
	}

	lines := strings.Split(code, "\n")
	var result strings.Builder

	for i, line := range lines {
		lineNum := i + 1
		if highlightSet[lineNum] {
			result.WriteString(`<span class="line highlight">`)
			result.WriteString(line)
			result.WriteString(`</span>`)
		} else {
			result.WriteString(`<span class="line">`)
			result.WriteString(line)
			result.WriteString(`</span>`)
		}
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
