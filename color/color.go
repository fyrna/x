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
	"fmt"
	"strconv"
	"strings"
)

// Enabled controls whether color output is enabled.
// Set to false to disable all color output (useful for non-TTY outputs).
var Enabled = true

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

	BrightBlack   = "\x1b[90m"
	BrightRed     = "\x1b[91m"
	BrightGreen   = "\x1b[92m"
	BrightYellow  = "\x1b[93m"
	BrightBlue    = "\x1b[94m"
	BrightMagenta = "\x1b[95m"
	BrightCyan    = "\x1b[96m"
	BrightWhite   = "\x1b[97m"

	BgBlack   = "\x1b[40m"
	BgRed     = "\x1b[41m"
	BgGreen   = "\x1b[42m"
	BgYellow  = "\x1b[43m"
	BgBlue    = "\x1b[44m"
	BgMagenta = "\x1b[45m"
	BgCyan    = "\x1b[46m"
	BgWhite   = "\x1b[47m"

	BgBrightBlack   = "\x1b[100m"
	BgBrightRed     = "\x1b[101m"
	BgBrightGreen   = "\x1b[102m"
	BgBrightYellow  = "\x1b[103m"
	BgBrightBlue    = "\x1b[104m"
	BgBrightMagenta = "\x1b[105m"
	BgBrightCyan    = "\x1b[106m"
	BgBrightWhite   = "\x1b[107m"

	Reset = "\x1b[0m"
)

// Text style constants
const (
	Bold      = "\x1b[1m"
	Faint     = "\x1b[2m"
	Italic    = "\x1b[3m"
	Underline = "\x1b[4m"
	Blink     = "\x1b[5m"
	Reverse   = "\x1b[7m"
	Conceal   = "\x1b[8m"
	Strike    = "\x1b[9m"
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
	return "\x1b[38;5;" + strconv.Itoa(n) + "m"
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
	return "\x1b[48;5;" + strconv.Itoa(n) + "m"
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
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
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
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
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
		return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
	}
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
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
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
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
			seq += fmt.Sprintf("\x1b[%dm", 40+(c.Value%8))
		} else {
			seq += fmt.Sprintf("\x1b[%dm", 30+(c.Value%8))
		}
	case Mode256:
		if c.Bg {
			seq += fmt.Sprintf("\x1b[48;5;%dm", c.Value)
		} else {
			seq += fmt.Sprintf("\x1b[38;5;%dm", c.Value)
		}
	case ModeRGB:
		if c.Bg {
			seq += fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
		} else {
			seq += fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
		}
	}

	if len(c.Styles) > 0 {
		seq += strings.Join(c.Styles, "")
	}

	return seq
}
