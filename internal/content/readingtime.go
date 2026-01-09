// Package content provides content analysis utilities like reading time calculation.
package content

import (
	"regexp"
	"strings"
	"unicode"
)

// ReadingTime holds reading time calculation results
type ReadingTime struct {
	Minutes int
	Words   int
}

const (
	wordsPerMinute     = 225
	codeWordsPerMinute = 100
)

// stripHTMLTags removes HTML tags from content
var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)

// codeBlockRegex matches code blocks
var codeBlockRegex = regexp.MustCompile(`(?s)<pre[^>]*>.*?</pre>|<code[^>]*>.*?</code>`)

// CalculateReadingTime estimates reading time from HTML content
func CalculateReadingTime(htmlContent string) ReadingTime {
	// Extract code blocks first
	codeBlocks := codeBlockRegex.FindAllString(htmlContent, -1)
	codeWords := 0
	for _, block := range codeBlocks {
		text := htmlTagRegex.ReplaceAllString(block, " ")
		codeWords += countWords(text)
	}

	// Remove code blocks from content for regular word count
	contentWithoutCode := codeBlockRegex.ReplaceAllString(htmlContent, " ")

	// Strip HTML tags
	plainText := htmlTagRegex.ReplaceAllString(contentWithoutCode, " ")

	// Count words in regular content
	regularWords := countWords(plainText)

	// Total words
	totalWords := regularWords + codeWords

	// Calculate reading time
	// Regular content at 225 wpm, code at 100 wpm
	regularMinutes := float64(regularWords) / float64(wordsPerMinute)
	codeMinutes := float64(codeWords) / float64(codeWordsPerMinute)
	totalMinutes := regularMinutes + codeMinutes

	// Round up, minimum 1 minute
	minutes := int(totalMinutes + 0.5)
	if minutes < 1 {
		minutes = 1
	}

	return ReadingTime{
		Minutes: minutes,
		Words:   totalWords,
	}
}

// countWords counts words in plain text
func countWords(text string) int {
	words := 0
	inWord := false

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if !inWord {
				words++
				inWord = true
			}
		} else {
			inWord = false
		}
	}

	return words
}

// FormatReadingTime returns a human-readable reading time string
func FormatReadingTime(rt ReadingTime) string {
	if rt.Minutes == 1 {
		return "1 min read"
	}
	return strings.TrimSpace(strings.Join([]string{itoa(rt.Minutes), "min read"}, " "))
}

// itoa converts int to string without importing strconv
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
