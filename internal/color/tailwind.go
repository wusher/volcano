package color

import (
	"fmt"
	"sort"
	"strings"
)

// TailwindColors maps Tailwind CSS color names to their 500-shade hex values.
// Sourced from Tailwind CSS v3 default palette.
var TailwindColors = map[string]string{
	"slate":   "#64748b",
	"gray":    "#6b7280",
	"zinc":    "#71717a",
	"neutral": "#737373",
	"stone":   "#78716c",
	"red":     "#ef4444",
	"orange":  "#f97316",
	"amber":   "#f59e0b",
	"yellow":  "#eab308",
	"lime":    "#84cc16",
	"green":   "#22c55e",
	"emerald": "#10b981",
	"teal":    "#14b8a6",
	"cyan":    "#06b6d4",
	"sky":     "#0ea5e9",
	"blue":    "#3b82f6",
	"indigo":  "#6366f1",
	"violet":  "#8b5cf6",
	"purple":  "#a855f7",
	"fuchsia": "#d946ef",
	"pink":    "#ec4899",
	"rose":    "#f43f5e",
}

// TailwindColorNames returns the supported Tailwind color names, sorted.
func TailwindColorNames() []string {
	names := make([]string, 0, len(TailwindColors))
	for name := range TailwindColors {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ResolveAccentColor converts an accent color input (Tailwind name or hex) to a hex string.
// Empty input returns empty without error. Hex inputs are validated and returned normalized.
// Tailwind names are case-insensitive.
func ResolveAccentColor(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	if strings.HasPrefix(input, "#") {
		if _, _, _, err := ParseHex(input); err != nil {
			return "", err
		}
		return input, nil
	}

	name := strings.ToLower(strings.TrimSpace(input))
	hex, ok := TailwindColors[name]
	if !ok {
		return "", fmt.Errorf(
			"unknown accent color %q. Expected a hex value like '#0ea5e9', a two-color gradient like 'lime-sky' or '#444444-#555555', or one of the Tailwind color names: %s",
			input,
			strings.Join(TailwindColorNames(), ", "),
		)
	}
	return hex, nil
}

// ResolveAccentSpec parses an accent-color spec that may name a single color
// or a two-color gradient. Supported forms (each side can be a Tailwind name
// or a hex value):
//
//	"sky"             → start="#0ea5e9", end=""
//	"#ff6600"         → start="#ff6600", end=""
//	"lime-sky"        → start="#84cc16", end="#0ea5e9"
//	"#444444-#555555" → start="#444444", end="#555555"
//	"sky-#ff6600"     → start="#0ea5e9", end="#ff6600"
//
// Empty input returns ("", "", nil). On parse error end is "".
func ResolveAccentSpec(input string) (start, end string, err error) {
	if input == "" {
		return "", "", nil
	}

	first, second, hasSecond := splitAccentSpec(input)
	start, err = ResolveAccentColor(first)
	if err != nil {
		return "", "", err
	}
	if !hasSecond {
		return start, "", nil
	}
	end, err = ResolveAccentColor(second)
	if err != nil {
		return "", "", err
	}
	return start, end, nil
}

// splitAccentSpec splits an accent spec on the dash separating the two colors.
// It is hex-aware: the boundary "-#" is preferred so "lime-#ff6600" parses
// correctly. If no "-#" exists it falls back to the first dash, which handles
// "lime-sky" and "#ff6600-sky". A third dash anywhere is rejected.
func splitAccentSpec(input string) (first, second string, ok bool) {
	if idx := strings.Index(input, "-#"); idx > 0 {
		// Second color is a hex. Whatever comes before "-#" is the first color.
		// The remainder ("-#..." minus the leading "-") starts with "#".
		first = input[:idx]
		second = input[idx+1:]
	} else if idx := strings.Index(input, "-"); idx > 0 {
		first = input[:idx]
		second = input[idx+1:]
	} else {
		return input, "", false
	}
	// Reject trailing/extra dashes (e.g. "red-blue-green") to keep parsing
	// unambiguous. Hex shorthand "#abc" has no dashes so this is safe.
	if !strings.HasPrefix(second, "#") && strings.Contains(second, "-") {
		return input, "", false
	}
	return first, second, true
}
