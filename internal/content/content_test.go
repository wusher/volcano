package content

import (
	"testing"
)

func TestCalculateReadingTime(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int // expected minutes
	}{
		{"empty", "", 1},
		{"short", "Hello world", 1},
		{"one minute", generateWords(200), 1},
		{"two minutes", generateWords(450), 2},
		{"five minutes", generateWords(1100), 5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateReadingTime(tc.content)
			if result.Minutes != tc.expected {
				t.Errorf("Minutes = %d, want %d", result.Minutes, tc.expected)
			}
		})
	}
}

func generateWords(count int) string {
	var content string
	for i := 0; i < count; i++ {
		content += "word "
	}
	return content
}

func TestFormatReadingTime(t *testing.T) {
	tests := []struct {
		rt       ReadingTime
		expected string
	}{
		{ReadingTime{Minutes: 1, Words: 100}, "1 min read"},
		{ReadingTime{Minutes: 5, Words: 1000}, "5 min read"},
		{ReadingTime{Minutes: 10, Words: 2000}, "10 min read"},
	}

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			result := FormatReadingTime(tc.rt)
			if result != tc.expected {
				t.Errorf("FormatReadingTime() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestCalculateReadingTimeWithCode(t *testing.T) {
	content := `<p>Some regular text here.</p>
<pre><code>func main() {
    fmt.Println("Hello, world!")
}</code></pre>
<p>More text here.</p>`

	rt := CalculateReadingTime(content)
	if rt.Words < 5 {
		t.Errorf("Words should be at least 5, got %d", rt.Words)
	}
	if rt.Minutes < 1 {
		t.Errorf("Minutes should be at least 1, got %d", rt.Minutes)
	}
}
