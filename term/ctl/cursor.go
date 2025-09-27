// Package ctl provides low-level ANSI escape sequences for terminal control.
package ctl

import "fmt"

// Cursor positioning.
const (
	// Home moves the cursor to the upper-left corner (row 1, column 1).
	Home = ESC + "H"

	// SavePos saves the current cursor position in memory.
	SavePos = ESC + "s"

	// RestorePos restores the cursor position previously saved with SavePos.
	RestorePos = ESC + "u"

	// SavePosAlt is an alternative DEC-compatible sequence to save the cursor
	// position.
	SavePosAlt = "\x1b7"

	// RestorePosAlt is an alternative DEC-compatible sequence to restore the
	// cursor position.
	RestorePosAlt = "\x1b8"
)

// Cursor movement by one cell.
const (
	// CursorUp moves the cursor one row up.
	CursorUp = ESC + "1A"

	// CursorDown moves the cursor one row down.
	CursorDown = ESC + "1B"

	// CursorRight moves the cursor one column right.
	CursorRight = ESC + "1C"

	// CursorLeft moves the cursor one column left.
	CursorLeft = ESC + "1D"
)

// Cursor visibility.
const (
	// HideCursor hides the text cursor.
	HideCursor = ESC + "?25l"

	// ShowCursor makes the text cursor visible again.
	ShowCursor = ESC + "?25h"
)

// MoveTo returns the escape sequence that moves the cursor to the specified
// 1-based row and column.
func MoveTo(row, col int) string {
	return fmt.Sprintf(ESC+"%d;%dH", row, col)
}

// MoveCol returns the escape sequence that moves the cursor to the specified
// 1-based column in the current row.
func MoveCol(col int) string {
	return fmt.Sprintf(ESC+"%dG", col)
}

// MoveRow returns the escape sequence that moves the cursor to the specified
// 1-based row in the current column.
func MoveRow(row int) string {
	return fmt.Sprintf(ESC+"%dd", row)
}

// MoveUpN returns the escape sequence that moves the cursor up by n rows.
// n must be positive.
func MoveUpN(n int) string {
	return fmt.Sprintf(ESC+"%dA", n)
}

// MoveDownN returns the escape sequence that moves the cursor down by n rows.
// n must be positive.
func MoveDownN(n int) string {
	return fmt.Sprintf(ESC+"%dB", n)
}

// MoveRightN returns the escape sequence that moves the cursor right by n
// columns.  n must be positive.
func MoveRightN(n int) string {
	return fmt.Sprintf(ESC+"%dC", n)
}

// MoveLeftN returns the escape sequence that moves the cursor left by n
// columns.  n must be positive.
func MoveLeftN(n int) string {
	return fmt.Sprintf(ESC+"%dD", n)
}
