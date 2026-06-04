// Package color provides color conversion utilities for accent color generation.
package color

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// HSL represents a color in HSL color space
type HSL struct {
	H float64 // Hue: 0-360
	S float64 // Saturation: 0-100
	L float64 // Lightness: 0-100
}

// AccentVariants contains the generated accent color variants
type AccentVariants struct {
	Accent      string // Original color normalized to 50% lightness
	AccentDark  string // Color at 10% lightness (for dark mode backgrounds)
	AccentLight string // Color at 95% lightness (for light mode backgrounds)
}

// hexPattern matches valid hex color formats: #RGB, #RRGGBB
var hexPattern = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)

// ParseHex parses a hex color string to RGB values (0-255)
func ParseHex(hex string) (r, g, b uint8, err error) {
	if !hexPattern.MatchString(hex) {
		return 0, 0, 0, fmt.Errorf("invalid hex color format: %s (expected #RGB or #RRGGBB)", hex)
	}

	hex = strings.TrimPrefix(hex, "#")

	// Expand shorthand (#RGB -> #RRGGBB)
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}

	rVal, _ := strconv.ParseUint(hex[0:2], 16, 8)
	gVal, _ := strconv.ParseUint(hex[2:4], 16, 8)
	bVal, _ := strconv.ParseUint(hex[4:6], 16, 8)

	return uint8(rVal), uint8(gVal), uint8(bVal), nil
}

// RGBToHex converts RGB values (0-255) to hex string
func RGBToHex(r, g, b uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// RGBToHSL converts RGB values (0-255) to HSL
func RGBToHSL(r, g, b uint8) HSL {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	maxVal := math.Max(math.Max(rf, gf), bf)
	minVal := math.Min(math.Min(rf, gf), bf)
	delta := maxVal - minVal

	var h, s, l float64
	l = (maxVal + minVal) / 2.0

	if delta == 0 {
		h, s = 0, 0
	} else {
		if l < 0.5 {
			s = delta / (maxVal + minVal)
		} else {
			s = delta / (2.0 - maxVal - minVal)
		}

		switch maxVal {
		case rf:
			h = (gf - bf) / delta
			if gf < bf {
				h += 6
			}
		case gf:
			h = ((bf - rf) / delta) + 2
		case bf:
			h = ((rf - gf) / delta) + 4
		}
		h /= 6
	}

	return HSL{H: h * 360, S: s * 100, L: l * 100}
}

// HSLToRGB converts HSL to RGB values (0-255)
func HSLToRGB(hsl HSL) (r, g, b uint8) {
	h := hsl.H / 360.0
	s := hsl.S / 100.0
	l := hsl.L / 100.0

	if s == 0 {
		gray := uint8(math.Round(l * 255))
		return gray, gray, gray
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	r = uint8(math.Round(hueToRGB(p, q, h+1.0/3.0) * 255))
	g = uint8(math.Round(hueToRGB(p, q, h) * 255))
	b = uint8(math.Round(hueToRGB(p, q, h-1.0/3.0) * 255))
	return
}

// hueToRGB is a helper for HSL to RGB conversion
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

// GenerateAccentVariants generates accent color variants from a hex color
// Returns: accent (50% L), dark (10% L), light (95% L)
func GenerateAccentVariants(hexColor string) (*AccentVariants, error) {
	r, g, b, err := ParseHex(hexColor)
	if err != nil {
		return nil, err
	}

	hsl := RGBToHSL(r, g, b)

	// Generate variants by adjusting lightness while keeping hue and saturation
	accent := HSL{H: hsl.H, S: hsl.S, L: 50}
	dark := HSL{H: hsl.H, S: hsl.S, L: 10}
	light := HSL{H: hsl.H, S: hsl.S, L: 95}

	return &AccentVariants{
		Accent:      RGBToHex(HSLToRGB(accent)),
		AccentDark:  RGBToHex(HSLToRGB(dark)),
		AccentLight: RGBToHex(HSLToRGB(light)),
	}, nil
}

// GenerateAccentCSS generates CSS custom properties for accent colors.
// Accepts a single color (Tailwind name or hex) or a two-color gradient spec
// like "lime-sky" / "#444444-#555555". Returns empty string if accentColor is empty.
//
// For gradients the rule emits three variables (--accent, --accent-end,
// --accent-gradient) and applies the gradient to the scroll progress bar and
// the prose H1 so the user-visible effect is immediate. Themes can also opt
// in to --accent-gradient elsewhere.
func GenerateAccentCSS(accentColor string) (string, error) {
	start, end, err := ResolveAccentSpec(accentColor)
	if err != nil {
		return "", err
	}
	if start == "" {
		return "", nil
	}

	if end == "" {
		// Single color — preserve existing behavior.
		return fmt.Sprintf(`:root, [data-theme="dark"] {
  --accent: %s;
}`, start), nil
	}

	// Two-color gradient. Two gradient variables ship:
	//   --accent-gradient            left-to-right — used for big backgrounds + text fills,
	//                                where reading direction dominates the perceived blend
	//   --accent-gradient-vertical   top-to-bottom — used for narrow vertical accents
	//                                (left-borders on blockquotes, admonitions, etc.)
	//
	// A diagonal (135°) gradient looked too "first-color heavy" on wide, short
	// headings because the gradient axis runs diagonally — most of the text
	// bounding box sat in the first half of the gradient. Horizontal direction
	// gives an even, predictable A→B sweep across the line.
	return fmt.Sprintf(`:root, [data-theme="dark"] {
  --accent: %s;
  --accent-end: %s;
  --accent-gradient: linear-gradient(to right, %s, %s);
  --accent-gradient-vertical: linear-gradient(to bottom, %s, %s);
}

.scroll-progress-bar {
  background: var(--accent-gradient);
}

/* Gradient text fill: applied to the page H1 and to in-content links
   (prose + prev/next nav). The text-decoration-color line keeps the
   underline visible since color: transparent would hide it otherwise. */
.prose h1,
.prose a,
.page-nav-prev,
.page-nav-next {
  background-image: var(--accent-gradient);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  -webkit-text-fill-color: transparent;
  text-decoration-color: var(--accent);
}

/* The H1 is block-level and spans the full content column, so the
   gradient (which is sized to the element's box) overshoots — the
   second color never reaches the text. Shrink the box to the rendered
   text width so the gradient stops map cleanly across the characters. */
.prose h1 {
  width: -moz-fit-content;
  width: fit-content;
  max-width: 100%%;
}

/* Vertical gradient on single-edge left-border accents (admonitions,
   blockquotes). border-image paints any side that has width — so we can
   only apply this to elements whose ONLY bordered side is the left. */
.admonition,
.prose blockquote {
  border-left-color: transparent;
  border-image: var(--accent-gradient-vertical) 1;
}

/* TOC sidebar has borders on all four sides plus a thicker accent left.
   border-image would paint all four edges, so paint the left bar as a
   layered background strip instead — keeps the rounded corners and the
   existing fill color intact. */
.toc-sidebar {
  border-left-color: transparent;
  background:
    var(--accent-gradient-vertical) left center / 3px 100%% no-repeat,
    var(--bg-primary);
}

/* Horizontal gradient on bottom-border accents (docs-theme H2 underline).
   Only border-bottom has width, so only the bottom edge paints. */
.prose h2 {
  border-bottom-color: transparent;
  border-image: var(--accent-gradient) 1;
}`, start, end, start, end, start, end), nil
}
