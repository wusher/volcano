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
		return "", fmt.Errorf("unknown accent color %q (expected a Tailwind color name like 'sky', 'rose', 'emerald', or a hex value like '#0ea5e9')", input)
	}
	return hex, nil
}
