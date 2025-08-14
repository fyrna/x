package screen

import (
	"os"
	"strconv"
)

// writeSeq writes an ANSI sequence "\x1b[nX" where X is the suffix byte.
func writeSeq(n int, suffix byte) {
	const prefix = "\x1b["
	var buf []byte

	if n == 0 {
		buf = []byte(prefix)
	} else {
		buf = make([]byte, 0, 8)
		buf = append(buf, prefix...)
		buf = append(buf, strconv.Itoa(n)...)
	}

	buf = append(buf, suffix)
	os.Stdout.Write(buf)
}

// Clear clears the entire screen and moves cursor to top-left.
func Clear() {
	os.Stdout.WriteString("\x1b[2J\x1b[H")
}

// ClearLine clears the current line according to mode:
// 0 = from cursor to end
// 1 = from cursor to beginning
// 2 = entire line (default)
func ClearLine(mode ...int) {
	m := 2
	if len(mode) > 0 {
		m = mode[0]
	}
	writeSeq(m, 'K')
}

// ClearLines clears n lines starting from the current cursor position.
func ClearLines(n int) {
	for i := range n {
		os.Stdout.WriteString("\x1b[2K")
		if i < n-1 {
			os.Stdout.WriteString("\x1b[1B")
		}
	}
	writeSeq(n, 'A')
}

// ClearFromCursorToEnd clears screen from cursor to end.
func ClearFromCursorToEnd() {
	os.Stdout.WriteString("\x1b[J")
}

// ClearFromCursorToStart clears screen from cursor to beginning.
func ClearFromCursorToStart() {
	os.Stdout.WriteString("\x1b[1J")
}

// ScrollUp scrolls the screen up by n lines.
func ScrollUp(n int) {
	writeSeq(n, 'S')
}

// ScrollDown scrolls the screen down by n lines.
func ScrollDown(n int) {
	writeSeq(n, 'T')
}

// SaveScreen saves the current screen state.
func SaveScreen() {
	os.Stdout.WriteString("\x1b[?47h")
}

// RestoreScreen restores a previously saved screen state.
func RestoreScreen() {
	os.Stdout.WriteString("\x1b[?47l")
}

// just aliases, i think short name is good, but explicit is better!

// ClearAndHome is an alias for Clear.
var ClearAndHome = Clear

// ClearUp is an alias for ClearFromCursorToStart.
var ClearUp = ClearFromCursorToStart

// ClearDown is an alias for ClearFromCursorToEnd.
var ClearDown = ClearFromCursorToEnd
