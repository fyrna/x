package color

import (
	"fmt"
	"strconv"
	"strings"
)

// HEXtoRGB converts hex to true-color ASCII (ANSI escape code)
// kind is:
// 0 = foreground color (default)
// 1 = background color
func HEXtoRGB(t int, hex string) string {
	hex = strings.TrimSpace(hex)
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) != 6 {
		return "" // empty string if invalid
	}

	// Parse each hex
	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return ""
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return ""
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return ""
	}

	// Return in ANSI true-color escape code format
	// Format: \x1b[38;2;r;g;bm (foreground color)
	if t != 1 {
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
	}

	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}
