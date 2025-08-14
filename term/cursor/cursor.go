package cursor

import (
	"os"
	"strconv"
)

func write(seq string) {
	os.Stdout.WriteString(seq)
}

// writeN writes "\x1b[n<suffix>" where n > 0
func writeN(n int, suffix byte) {
	if n < 1 {
		n = 1
	}
	write("\x1b[" + strconv.Itoa(n) + string(suffix))
}

// getN get first value of slice (for variadic, default 1)
func getN(n []int) int {
	if len(n) > 0 {
		return n[0]
	}
	return 1
}

// --- Movement (relative) ---
func Up(n ...int)    { writeN(getN(n), 'A') }
func Down(n ...int)  { writeN(getN(n), 'B') }
func Right(n ...int) { writeN(getN(n), 'C') }
func Left(n ...int)  { writeN(getN(n), 'D') }

// --- Absolute position (1-based coordinates) ---
func To(row, col int) {
	write("\x1b[" + strconv.Itoa(row) + ";" + strconv.Itoa(col) + "H")
}

// --- Save and restore position
func SavePos()    { write("\x1b7") }
func RestorePos() { write("\x1b8") }

// --- Visibility ---

func Hide() { write("\x1b[?25l") }
func Show() { write("\x1b[?25h") }

// --- utility ---

// Home moves cursor to top-left (1,1 position)
func Home() { write("\x1b[H") }

// Move to beginning of line (column 1)
func StartOfLine() { write("\r") }

// Move to specific column (1-based)
func ToColumn(col int) { write("\x1b[" + strconv.Itoa(col) + "G") }

// NextLine moves to beginning of next line
func NextLine() { write("\x1b[E") }

// PrevLine moves to beginning of previous line
func PrevLine() { write("\x1b[F") }
