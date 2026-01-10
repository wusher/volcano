package color

import (
	"strings"
	"testing"
)

func TestParseHex(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		wantR   uint8
		wantG   uint8
		wantB   uint8
		wantErr bool
	}{
		{
			name:  "full hex orange",
			hex:   "#ff6600",
			wantR: 255, wantG: 102, wantB: 0,
		},
		{
			name:  "full hex black",
			hex:   "#000000",
			wantR: 0, wantG: 0, wantB: 0,
		},
		{
			name:  "full hex white",
			hex:   "#ffffff",
			wantR: 255, wantG: 255, wantB: 255,
		},
		{
			name:  "full hex uppercase",
			hex:   "#FF6600",
			wantR: 255, wantG: 102, wantB: 0,
		},
		{
			name:  "shorthand hex red",
			hex:   "#f00",
			wantR: 255, wantG: 0, wantB: 0,
		},
		{
			name:  "shorthand hex white",
			hex:   "#fff",
			wantR: 255, wantG: 255, wantB: 255,
		},
		{
			name:    "invalid - no hash",
			hex:     "ff6600",
			wantErr: true,
		},
		{
			name:    "invalid - wrong length",
			hex:     "#ffff",
			wantErr: true,
		},
		{
			name:    "invalid - non-hex chars",
			hex:     "#gggggg",
			wantErr: true,
		},
		{
			name:    "invalid - empty",
			hex:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, err := ParseHex(tt.hex)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseHex(%q) expected error, got nil", tt.hex)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseHex(%q) unexpected error: %v", tt.hex, err)
				return
			}
			if r != tt.wantR || g != tt.wantG || b != tt.wantB {
				t.Errorf("ParseHex(%q) = (%d, %d, %d), want (%d, %d, %d)",
					tt.hex, r, g, b, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

func TestRGBToHex(t *testing.T) {
	tests := []struct {
		r, g, b uint8
		want    string
	}{
		{255, 102, 0, "#ff6600"},
		{0, 0, 0, "#000000"},
		{255, 255, 255, "#ffffff"},
		{128, 128, 128, "#808080"},
	}

	for _, tt := range tests {
		got := RGBToHex(tt.r, tt.g, tt.b)
		if got != tt.want {
			t.Errorf("RGBToHex(%d, %d, %d) = %q, want %q", tt.r, tt.g, tt.b, got, tt.want)
		}
	}
}

func TestRGBToHSL(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
		wantH   float64 // approximate
		wantS   float64 // approximate
		wantL   float64 // approximate
	}{
		{"red", 255, 0, 0, 0, 100, 50},
		{"green", 0, 255, 0, 120, 100, 50},
		{"blue", 0, 0, 255, 240, 100, 50},
		{"white", 255, 255, 255, 0, 0, 100},
		{"black", 0, 0, 0, 0, 0, 0},
		{"gray", 128, 128, 128, 0, 0, 50},
		{"orange", 255, 165, 0, 39, 100, 50}, // roughly 39 degrees
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hsl := RGBToHSL(tt.r, tt.g, tt.b)

			// Allow some tolerance for floating point
			tolerance := 2.0

			// For achromatic colors (white, black, gray), hue is undefined
			if tt.wantS > 0 {
				if diff := abs(hsl.H - tt.wantH); diff > tolerance && diff < 360-tolerance {
					t.Errorf("RGBToHSL H = %f, want ~%f", hsl.H, tt.wantH)
				}
			}
			if abs(hsl.S-tt.wantS) > tolerance {
				t.Errorf("RGBToHSL S = %f, want ~%f", hsl.S, tt.wantS)
			}
			if abs(hsl.L-tt.wantL) > tolerance {
				t.Errorf("RGBToHSL L = %f, want ~%f", hsl.L, tt.wantL)
			}
		})
	}
}

func TestHSLToRGB(t *testing.T) {
	tests := []struct {
		name    string
		hsl     HSL
		wantR   uint8
		wantG   uint8
		wantB   uint8
		epsilon uint8 // allow small rounding differences
	}{
		{"red", HSL{0, 100, 50}, 255, 0, 0, 1},
		{"green", HSL{120, 100, 50}, 0, 255, 0, 1},
		{"blue", HSL{240, 100, 50}, 0, 0, 255, 1},
		{"white", HSL{0, 0, 100}, 255, 255, 255, 1},
		{"black", HSL{0, 0, 0}, 0, 0, 0, 1},
		{"gray", HSL{0, 0, 50}, 128, 128, 128, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b := HSLToRGB(tt.hsl)
			if absDiff(r, tt.wantR) > tt.epsilon ||
				absDiff(g, tt.wantG) > tt.epsilon ||
				absDiff(b, tt.wantB) > tt.epsilon {
				t.Errorf("HSLToRGB(%v) = (%d, %d, %d), want (%d, %d, %d)",
					tt.hsl, r, g, b, tt.wantR, tt.wantG, tt.wantB)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	// Test that RGB -> HSL -> RGB produces the same values
	colors := [][3]uint8{
		{255, 0, 0},     // red
		{0, 255, 0},     // green
		{0, 0, 255},     // blue
		{255, 255, 0},   // yellow
		{0, 255, 255},   // cyan
		{255, 0, 255},   // magenta
		{255, 165, 0},   // orange
		{128, 64, 192},  // purple-ish
		{100, 150, 200}, // blue-gray
	}

	for _, c := range colors {
		r, g, b := c[0], c[1], c[2]
		hsl := RGBToHSL(r, g, b)
		r2, g2, b2 := HSLToRGB(hsl)

		// Allow small rounding differences
		if absDiff(r, r2) > 1 || absDiff(g, g2) > 1 || absDiff(b, b2) > 1 {
			t.Errorf("Round trip (%d,%d,%d) -> HSL -> (%d,%d,%d) mismatch",
				r, g, b, r2, g2, b2)
		}
	}
}

func TestGenerateAccentVariants(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		wantErr bool
	}{
		{"orange", "#ff6600", false},
		{"blue", "#0066cc", false},
		{"purple", "#8b00ff", false},
		{"red", "#ff0000", false},
		{"green", "#00ff00", false},
		{"invalid", "not-a-color", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variants, err := GenerateAccentVariants(tt.hex)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify all variants are valid hex colors
			for _, hex := range []string{variants.Accent, variants.AccentDark, variants.AccentLight} {
				if !strings.HasPrefix(hex, "#") || len(hex) != 7 {
					t.Errorf("invalid hex format: %s", hex)
				}
				_, _, _, err := ParseHex(hex)
				if err != nil {
					t.Errorf("generated hex is invalid: %s", hex)
				}
			}

			// Verify lightness values by converting back
			// Accent should be around 50% lightness
			r, g, b, _ := ParseHex(variants.Accent)
			hsl := RGBToHSL(r, g, b)
			if abs(hsl.L-50) > 2 {
				t.Errorf("Accent lightness = %f, want ~50", hsl.L)
			}

			// Dark should be around 10% lightness
			r, g, b, _ = ParseHex(variants.AccentDark)
			hsl = RGBToHSL(r, g, b)
			if abs(hsl.L-10) > 2 {
				t.Errorf("AccentDark lightness = %f, want ~10", hsl.L)
			}

			// Light should be around 95% lightness
			r, g, b, _ = ParseHex(variants.AccentLight)
			hsl = RGBToHSL(r, g, b)
			if abs(hsl.L-95) > 2 {
				t.Errorf("AccentLight lightness = %f, want ~95", hsl.L)
			}
		})
	}
}

func TestGenerateAccentCSS(t *testing.T) {
	t.Run("empty color", func(t *testing.T) {
		css, err := GenerateAccentCSS("")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if css != "" {
			t.Errorf("expected empty string, got %q", css)
		}
	})

	t.Run("valid color", func(t *testing.T) {
		css, err := GenerateAccentCSS("#ff6600")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !strings.Contains(css, "--accent:") {
			t.Errorf("CSS missing --accent: %q", css)
		}
		if !strings.Contains(css, "--accent-dark:") {
			t.Errorf("CSS missing --accent-dark: %q", css)
		}
		if !strings.Contains(css, "--accent-light:") {
			t.Errorf("CSS missing --accent-light: %q", css)
		}
		if !strings.Contains(css, ":root") {
			t.Errorf("CSS missing :root: %q", css)
		}
	})

	t.Run("invalid color", func(t *testing.T) {
		_, err := GenerateAccentCSS("invalid")
		if err == nil {
			t.Error("expected error for invalid color")
		}
	})
}

// Helper functions
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}
