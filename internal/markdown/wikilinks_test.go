package markdown

import (
	"testing"
)

func TestConvertWikiLinks(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		currentURLPath string
		expected       string
	}{
		// Tests with empty currentURLPath (root level) - absolute paths from root
		{"simple page at root", "[[Page Name]]", "", "[Page Name](/page-name/)"},
		{"with display text at root", "[[Page|Display Text]]", "", "[Display Text](/page/)"},
		{"path with numbers at root", "[[5. Guidance/Old Guidance/2023 Goals]]", "", "[2023 Goals](/guidance/old-guidance/2023-goals/)"},
		{"simple name at root", "[[Guidance Readme]]", "", "[Guidance Readme](/guidance-readme/)"},
		{"with .md extension at root", "[[Page.md]]", "", "[Page](/page/)"},
		{"multiple on one line at root", "[[Link1]] and [[Link2]]", "", "[Link1](/link1/) and [Link2](/link2/)"},
		{"in a list at root", "- [[Item One]]\n- [[Item Two]]", "", "- [Item One](/item-one/)\n- [Item Two](/item-two/)"},
		{"no conversion needed", "Normal text", "", "Normal text"},
		{"standard markdown link unchanged", "[text](url)", "", "[text](url)"},

		// Tests with currentURLPath (relative resolution)
		{"simple page relative to guidance", "[[Guidance Readme]]", "/guidance/", "[Guidance Readme](/guidance/guidance-readme/)"},
		{"simple page relative to nested dir", "[[Page]]", "/docs/api/", "[Page](/docs/api/page/)"},
		{"explicit path ignores current dir", "[[other/Page]]", "/guidance/", "[Page](/other/page/)"},
		{"root path ignores current dir", "[[Page]]", "/", "[Page](/page/)"},

		// Display text with relative resolution
		{"display text relative", "[[Life Goals|My Goals]]", "/guidance/", "[My Goals](/guidance/life-goals/)"},

		// Multiple links with relative resolution
		{"multiple relative links", "[[First]] and [[Second]]", "/guidance/", "[First](/guidance/first/) and [Second](/guidance/second/)"},

		// Mixed - some relative, some absolute
		{"mixed relative and absolute", "[[Local]] and [[other/Absolute]]", "/guidance/", "[Local](/guidance/local/) and [Absolute](/other/absolute/)"},

		// Embeds (![[...]]) converted to regular links
		{"embed converted to link", "![[Page Name]]", "", "[Page Name](/page-name/)"},
		{"embed with path", "![[5. Guidance/Old Guidance/2023 Goals]]", "", "[2023 Goals](/guidance/old-guidance/2023-goals/)"},
		{"embed relative", "![[Life Goals]]", "/guidance/", "[Life Goals](/guidance/life-goals/)"},
		{"embed with display text", "![[Page|Custom Text]]", "", "[Custom Text](/page/)"},

		// Index/readme files resolve to parent directory
		{"index resolves to current dir", "[[index]]", "/guidance/", "[index](/guidance/)"},
		{"readme resolves to current dir", "[[readme]]", "/guidance/", "[readme](/guidance/)"},
		{"Index case insensitive", "[[Index]]", "/guidance/", "[Index](/guidance/)"},
		{"README case insensitive", "[[README]]", "/guidance/", "[README](/guidance/)"},
		{"index at root", "[[index]]", "/", "[index](/)"},
		{"folder/index resolves to folder", "[[other/index]]", "/guidance/", "[index](/other/)"},
		{"folder/readme resolves to folder", "[[docs/readme]]", "/", "[readme](/docs/)"},

		// Anchors (fragments) are preserved
		{"page with anchor", "[[faq#permissions]]", "/guides/", "[faq#permissions](/guides/faq/#permissions)"},
		{"page with anchor at root", "[[faq#section]]", "", "[faq#section](/faq/#section)"},
		{"page with anchor and display text", "[[faq#help|Get Help]]", "", "[Get Help](/faq/#help)"},
		{"explicit path with anchor", "[[reference/api#methods]]", "/guides/", "[api#methods](/reference/api/#methods)"},
		{"just anchor (same page)", "[[#section]]", "/guides/", "[#section](#section)"},

		// Attachments preserve file extension and don't get trailing slash
		{"image png", "[[pasted-image-20251226101657.png]]", "", "[pasted-image-20251226101657.png](/pasted-image-20251226101657.png)"},
		{"image with spaces", "[[My Screenshot.png]]", "", "[My Screenshot.png](/my-screenshot.png)"},
		{"image in folder", "[[attachments/image.jpg]]", "", "[image.jpg](/attachments/image.jpg)"},
		{"image relative", "[[screenshot.png]]", "/system/attachments/", "[screenshot.png](/system/attachments/screenshot.png)"},
		{"pdf attachment", "[[documents/report.pdf]]", "", "[report.pdf](/documents/report.pdf)"},
		{"video attachment", "[[videos/demo.mp4]]", "", "[demo.mp4](/videos/demo.mp4)"},
		{"embed image", "![[image.png]]", "/assets/", "[image.png](/assets/image.png)"},
		{"image with display text", "[[photo.jpg|My Photo]]", "", "[My Photo](/photo.jpg)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(ConvertWikiLinks([]byte(tt.input), tt.currentURLPath))
			if result != tt.expected {
				t.Errorf("ConvertWikiLinks(%q, %q) = %q, want %q", tt.input, tt.currentURLPath, result, tt.expected)
			}
		})
	}
}
