package cursor

import "fmt"

func step(n []int) int {
	if len(n) > 0 {
		return n[0]
	}
	return 1
}

// Movement (relative)
func Up(n ...int) {
	fmt.Printf("\x1b[%dA", step(n))
}
func Down(n ...int) {
	fmt.Printf("\x1b[%dB", step(n))
}
func Right(n ...int) {
	fmt.Printf("\x1b[%dC", step(n))
}
func Left(n ...int) {
	fmt.Printf("\x1b[%dD", step(n))
}

// Absolute position (row, col)
func To(row, col int) {
	fmt.Printf("\x1b[%d;%dH", row, col)
}

// Save & restore position
func SavePos() {
	fmt.Print("\x1b7")
}
func RestorePos() {
	fmt.Print("\x1b8")
}

// Cursor visibility
func Hide() {
	fmt.Print("\x1b[?25l")
}
func Show() {
	fmt.Print("\x1b[?25h")
}

// Home (top-left)
func Home() {
	fmt.Print("\x1b[H")
}
