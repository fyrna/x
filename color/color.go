// Package color provides a tiny, zero-dependency helper set for ANSI colors
// (16 standard, 256 extended, 24-bit true color, and HEX).
//
// Quick start
//
//	// 16-color
//	fmt.Println(color.Red + "plain red" + color.Reset)
//
//	// 256-color
//	fmt.Println(color.Fg256(196) + "bright red" + color.Reset)
//
//	// 24-bit RGB
//	fmt.Println(color.FgRGB(255, 105, 180) + "hot pink" + color.Reset)
//
//	// Hex
//	fmt.Println(color.FgHEX("#00ff00") + "green from hex" + color.Reset)
package color

import (
	"fmt"
	"strconv"
	"strings"
)

// 16-color ANSI constants
const (
	Black   = "\x1b[30m"
	Red     = "\x1b[31m"
	Green   = "\x1b[32m"
	Yellow  = "\x1b[33m"
	Blue    = "\x1b[34m"
	Magenta = "\x1b[35m"
	Cyan    = "\x1b[36m"
	White   = "\x1b[37m"

	BoldBlack   = "\x1b[30;1m"
	BoldRed     = "\x1b[31;1m"
	BoldGreen   = "\x1b[32;1m"
	BoldYellow  = "\x1b[33;1m"
	BoldBlue    = "\x1b[34;1m"
	BoldMagenta = "\x1b[35;1m"
	BoldCyan    = "\x1b[36;1m"
	BoldWhite   = "\x1b[37;1m"

	BgBlack   = "\x1b[40m"
	BgRed     = "\x1b[41m"
	BgGreen   = "\x1b[42m"
	BgYellow  = "\x1b[43m"
	BgBlue    = "\x1b[44m"
	BgMagenta = "\x1b[45m"
	BgCyan    = "\x1b[46m"
	BgWhite   = "\x1b[47m"

	Reset = "\x1b[0m"
)

// Fg256 returns the escape sequence for 8-bit foreground color (0-255).
func Fg256(n int) string {
	return "\x1b[38;5;" + strconv.Itoa(n) + "m"
}

// Bg256 returns the escape sequence for 8-bit background color (0-255).
func Bg256(n int) string {
	return "\x1b[48;5;" + strconv.Itoa(n) + "m"
}

// FgRGB builds an RGB foreground escape sequence.
// Each component must be 0-255.
func FgRGB(r, g, b int) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

// BgRGB builds an RGB background escape sequence.
func BgRGB(r, g, b int) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

// FgHEX is syntactic sugar for HEXtoRGB(0, hex).
func FgHEX(hex string) string { return HEXtoRGB(0, hex) }

// BgHEX is syntactic sugar for HEXtoRGB(1, hex).
func BgHEX(hex string) string { return HEXtoRGB(1, hex) }

// HEXtoRGB converts a hex color string to a 24-bit ANSI escape sequence.
// kind: 0 → foreground, 1 → background.
// Returns an empty string on malformed input.
func HEXtoRGB(kind int, hex string) string {
	hex = strings.TrimSpace(hex)
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return ""
	}
	parse := func(s string) int64 {
		v, _ := strconv.ParseInt(s, 16, 64)
		return v
	}
	r := parse(hex[0:2])
	g := parse(hex[2:4])
	b := parse(hex[4:6])

	if kind == 1 {
		return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
	}
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

// Wrap wraps text with the given color escape sequence and a trailing reset.
// This is a convenience function for color codes obtained from Fg256, FgRGB, FgHEX,
// Bg256, BgRGB, BgHEX, or any of the 16-color constants.
//
// Example:
//
//	white := color.FgHEX("#ffffff")
//	blue := color.Blue
//	fmt.Println(color.Wrap(white, "hello") + " " + color.Wrap(blue, "world"))
func Wrap(color, text string) string {
	return color + text + Reset
}

// Color represents a 24-bit RGB color plus an optional human-readable name.
type Color struct {
	R, G, B uint8
	Name    string
}

// ToANSI returns the foreground escape sequence for this color.
func (c Color) ToANSI() string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

// ToHEX returns the color in #rrggbb form.
func (c Color) ToHEX() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// Wrap wraps text with the color’s escape codes and a trailing reset.
func (c Color) Wrap(text string) string {
	return c.ToANSI() + text + Reset
}
