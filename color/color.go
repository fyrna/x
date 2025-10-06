// Package color provides a tiny, zero-dependency helper set for ANSI colors
// (16 standard, 256 extended, 24-bit true color, and HEX) with additional
// text styles and utilities.
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
//
//	// With wrapper
//	fmt.Println(color.Wrap(color.Red, "red text"))
//
//	// With styles
//	fmt.Println(color.Wrap(color.Red+color.Bold, "bold red text"))
package color

import (
	"strconv"
	"strings"
)

const esc = "\x1b["

// Enabled controls whether color output is enabled.
// Set to false to disable all color output (useful for non-TTY outputs).
var Enabled = true

// 16-color ANSI constants
const (
	Black   = esc + "30m"
	Red     = esc + "31m"
	Green   = esc + "32m"
	Yellow  = esc + "33m"
	Blue    = esc + "34m"
	Magenta = esc + "35m"
	Cyan    = esc + "36m"
	White   = esc + "37m"

	BrightBlack   = esc + "90m"
	BrightRed     = esc + "91m"
	BrightGreen   = esc + "92m"
	BrightYellow  = esc + "93m"
	BrightBlue    = esc + "94m"
	BrightMagenta = esc + "95m"
	BrightCyan    = esc + "96m"
	BrightWhite   = esc + "97m"

	BgBlack   = esc + "40m"
	BgRed     = esc + "41m"
	BgGreen   = esc + "42m"
	BgYellow  = esc + "43m"
	BgBlue    = esc + "44m"
	BgMagenta = esc + "45m"
	BgCyan    = esc + "46m"
	BgWhite   = esc + "47m"

	BgBrightBlack   = esc + "100m"
	BgBrightRed     = esc + "101m"
	BgBrightGreen   = esc + "102m"
	BgBrightYellow  = esc + "103m"
	BgBrightBlue    = esc + "104m"
	BgBrightMagenta = esc + "105m"
	BgBrightCyan    = esc + "106m"
	BgBrightWhite   = esc + "107m"
)

// Text style constants
const (
	Reset      = esc + "0m"
	Bold       = esc + "1m"
	Faint      = esc + "2m"
	Italic     = esc + "3m"
	Underline  = esc + "4m"
	Blink      = esc + "5m"
	BlinkRapid = esc + "6m"
	Reverse    = esc + "7m"
	Conceal    = esc + "8m"
	Strike     = esc + "9m"
	Overline   = esc + "53m"
)

// Fg256 returns the escape sequence for 8-bit foreground color (0-255).
// Returns empty string if color is out of range.
func Fg256(n int) string {
	if !Enabled {
		return ""
	}
	if n < 0 || n > 255 {
		return ""
	}
	return esc + "38;5;" + strconv.Itoa(n) + "m"
}

// Bg256 returns the escape sequence for 8-bit background color (0-255).
// Returns empty string if color is out of range.
func Bg256(n int) string {
	if !Enabled {
		return ""
	}
	if n < 0 || n > 255 {
		return ""
	}
	return esc + "48;5;" + strconv.Itoa(n) + "m"
}

// FgRGB builds an RGB foreground escape sequence.
// Each component must be 0-255. Returns empty string if out of range.
func FgRGB(r, g, b int) string {
	if !Enabled {
		return ""
	}
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return ""
	}
	return esc + "38;2;" +
		strconv.Itoa(r) + ";" +
		strconv.Itoa(g) + ";" +
		strconv.Itoa(b) + "m"
}

// BgRGB builds an RGB background escape sequence.
// Each component must be 0-255. Returns empty string if out of range.
func BgRGB(r, g, b int) string {
	if !Enabled {
		return ""
	}
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return ""
	}
	return esc + "48;2;" +
		strconv.Itoa(r) + ";" +
		strconv.Itoa(g) + ";" +
		strconv.Itoa(b) + "m"
}

// FgHEX converts a hex color string to a foreground ANSI escape sequence.
// Supports 3-digit (#rgb) and 6-digit (#rrggbb) formats.
// Returns empty string on malformed input.
func FgHEX(hex string) string {
	if !Enabled {
		return ""
	}
	return hexToANSISequence(hex, false)
}

// BgHEX converts a hex color string to a background ANSI escape sequence.
// Supports 3-digit (#rgb) and 6-digit (#rrggbb) formats.
// Returns empty string on malformed input.
func BgHEX(hex string) string {
	if !Enabled {
		return ""
	}
	return hexToANSISequence(hex, true)
}

// hexToANSISequence converts a hex color string to an ANSI escape sequence.
// isBg: false → foreground, true → background.
// Returns an empty string on malformed input.
func hexToANSISequence(hex string, isBg bool) string {
	hex = strings.TrimSpace(hex)
	hex = strings.TrimPrefix(hex, "#")

	var r, g, b uint8
	var err error

	switch len(hex) {
	case 3: // #rgb format
		// Double each character: #rgb -> #rrggbb
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
		fallthrough
	case 6: // #rrggbb format
		parse := func(s string) (uint8, error) {
			v, err := strconv.ParseUint(s, 16, 8)
			return uint8(v), err
		}

		r, err = parse(hex[0:2])
		if err != nil {
			return ""
		}

		g, err = parse(hex[2:4])
		if err != nil {
			return ""
		}

		b, err = parse(hex[4:6])
		if err != nil {
			return ""
		}
	default:
		return ""
	}

	if isBg {
		return esc + "48;2;" +
			strconv.Itoa(int(r)) + ";" +
			strconv.Itoa(int(g)) + ";" +
			strconv.Itoa(int(b)) + "m"
	}
	return esc + "38;2;" +
		strconv.Itoa(int(r)) + ";" +
		strconv.Itoa(int(g)) + ";" +
		strconv.Itoa(int(b)) + "m"
}

// Wrap wraps text with the given color escape sequence and a trailing reset.
// If color output is disabled, returns the original text without escape sequences.
// This is a convenience function for color codes obtained from Fg256, FgRGB, FgHEX,
// Bg256, BgRGB, BgHEX, or any of the 16-color constants.
//
// Example:
//
//	white := color.FgHEX("#ffffff")
//	blue := color.Blue
//	fmt.Println(color.Wrap(white, "hello") + " " + color.Wrap(blue, "world"))
func Wrap(color, text string) string {
	if !Enabled || color == "" {
		return text
	}
	return color + text + Reset
}

// Style combines multiple ANSI escape sequences (colors and styles).
// Useful for creating compound styles like bold red text.
func Style(codes ...string) string {
	if !Enabled {
		return ""
	}
	return strings.Join(codes, "")
}

type Mode int

const (
	ModeANSI Mode = iota // 16-color standard
	Mode256              // 0–255
	ModeRGB              // R, G, B
)

type Color struct {
	Name    string
	Mode    Mode     // 16, 256, RGB
	Value   int      // index utk 16/256, atau inline ke RGB
	R, G, B uint8    // RGB Mode
	Bg      bool     // true = background, false = foreground
	Styles  []string // bold, italic, underline, dll
}

// ToHEX returns the color in #rrggbb form.
func (c Color) ToHEX() string {
	return "#" + byteToHex(c.R) + byteToHex(c.G) + byteToHex(c.B)
}

// Wrap wraps text with the color's foreground escape codes and a trailing reset.
func (c Color) Wrap(text string) string {
	return Wrap(c.ToANSI(), text)
}

// ToANSI returns the escape sequence for this color.
func (c Color) ToANSI() string {
	if !Enabled {
		return ""
	}

	seq := ""

	switch c.Mode {
	case ModeANSI:
		if c.Bg {
			seq += esc + strconv.Itoa(40+(c.Value%8)) + "m"
		} else {
			seq += esc + strconv.Itoa(30+(c.Value%8)) + "m"
		}
	case Mode256:
		if c.Bg {
			seq += esc + "48;5;" + strconv.Itoa(c.Value) + "m"
		} else {
			seq += esc + "38;5;" + strconv.Itoa(c.Value) + "m"
		}
	case ModeRGB:
		if c.Bg {
			seq += esc + "48;2;" +
				strconv.Itoa(int(c.R)) + ";" +
				strconv.Itoa(int(c.G)) + ";" +
				strconv.Itoa(int(c.B)) + "m"
		} else {
			seq += esc + "38;2;" +
				strconv.Itoa(int(c.R)) + ";" +
				strconv.Itoa(int(c.G)) + ";" +
				strconv.Itoa(int(c.B)) + "m"
		}
	}

	if len(c.Styles) > 0 {
		seq += strings.Join(c.Styles, "")
	}

	return seq
}

func byteToHex(b uint8) string {
	if b < 16 {
		return "0" + strconv.FormatUint(uint64(b), 16)
	}
	return strconv.FormatUint(uint64(b), 16)
}
