package screen

import "fmt"

// Clear entire screen
func Clear() {
	fmt.Print("\x1b[2J")
}

// Clear current line (mode: 0=to end, 1=to start, 2=entire line)
func ClearLine(mode ...int) {
	m := 2 // default: entire line
	if len(mode) > 0 {
		m = mode[0]
	}
	fmt.Printf("\x1b[%dK", m)
}

// Clear part of screen
func ClearFromCursorToEnd() {
	fmt.Print("\x1b[J")
}
func ClearFromCursorToStart() {
	fmt.Print("\x1b[1J")
}

// Combined helpers
func ClearAndHome() {
	fmt.Print("\x1b[2J\x1b[H")
}
