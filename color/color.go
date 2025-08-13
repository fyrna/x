package color

import "strconv"

// 16 Warna ANSI (Constants)
const (
	// Foreground
	Black   = "\x1b[30m"
	Red     = "\x1b[31m"
	Green   = "\x1b[32m"
	Yellow  = "\x1b[33m"
	Blue    = "\x1b[34m"
	Magenta = "\x1b[35m"
	Cyan    = "\x1b[36m"
	White   = "\x1b[37m"

	// Bold Foreground
	BoldBlack   = "\x1b[30;1m"
	BoldRed     = "\x1b[31;1m"
	BoldGreen   = "\x1b[32;1m"
	BoldYellow  = "\x1b[33;1m"
	BoldBlue    = "\x1b[34;1m"
	BoldMagenta = "\x1b[35;1m"
	BoldCyan    = "\x1b[36;1m"
	BoldWhite   = "\x1b[37;1m"

	// Background
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

// 256

func Fg256(cc int) string {
	return "\x1b[38;5;" + strconv.Itoa(cc) + "m"
}

func Bg256(cc int) string {
	return "\x1b[48;5;" + strconv.Itoa(cc) + "m"
}

// True Color (RGB Helpers)

func FgRGB(r, g, b int) string {
	return "\x1b[38;2;" + strconv.Itoa(r) + ";" +
		strconv.Itoa(g) + ";" + strconv.Itoa(b) + "m"
}

func BgRGB(r, g, b int) string {
	return "\x1b[48;2;" + strconv.Itoa(r) + ";" +
		strconv.Itoa(g) + ";" + strconv.Itoa(b) + "m"
}
