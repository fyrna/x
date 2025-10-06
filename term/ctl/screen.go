// Package ctl provides low-level ANSI escape sequences for terminal control.
package ctl

import "strconv"

// Screen and scroll-back buffer.
const (
	// Reset resets all character attributes and styles to their defaults.
	Reset = ESC + "0m"

	// ClearScreen clears the entire visible screen and moves the cursor to the
	// home position (1,1).
	ClearScreen = ESC + "2J"

	// ClearScreenDown clears the screen from the cursor to the bottom
	// (inclusive).
	ClearScreenDown = ESC + "0J"

	// ClearScreenUp clears the screen from the top to the cursor (inclusive).
	ClearScreenUp = ESC + "1J"
)

// Line-oriented erasure.
const (
	// ClearLine erases the complete line the cursor is currently on.
	ClearLine = ESC + "2K"

	// ClearLineRight erases from the cursor to the end of the line (inclusive).
	ClearLineRight = ESC + "0K"

	// ClearLineLeft erases from the start of the line to the cursor (inclusive).
	ClearLineLeft = ESC + "1K"
)

// Scrolling.
const (
	// ScrollUpOne scrolls the viewport up by one line.
	ScrollUpOne = ESC + "S"

	// ScrollDownOne scrolls the viewport down by one line.
	ScrollDownOne = ESC + "T"
)

// Line wrapping.
const (
	// EnableWrap enables automatic line wrapping when the cursor reaches the
	// right border.
	EnableWrap = ESC + "?7h"

	// DisableWrap disables automatic line wrapping; characters at the right
	// border overwrite the last column.
	DisableWrap = ESC + "?7l"
)

// Line insertion / deletion.
const (
	// InsertLine inserts a blank line at the cursor row, shifting existing
	// lines downward.
	InsertLine = ESC + "L"

	// DeleteLine deletes the line at the cursor, shifting existing lines
	// upward.
	DeleteLine = ESC + "M"
)

// ScrollUpN returns the escape sequence that scrolls the viewport up by n
// lines.  n must be positive.
func ScrollUpN(n int) string {
	return ESC + strconv.Itoa(n) + "S"
}

// ScrollDownN returns the escape sequence that scrolls the viewport down by n
// lines.  n must be positive.
func ScrollDownN(n int) string {
	return ESC + strconv.Itoa(n) + "T"
}

// InsertLineN returns the escape sequence that inserts n blank lines at the
// cursor row.  n must be positive.
func InsertLineN(n int) string {
	return ESC + strconv.Itoa(n) + "L"
}

// DeleteLineN returns the escape sequence that deletes n lines starting at the
// cursor row.  n must be positive.
func DeleteLineN(n int) string {
	return ESC + strconv.Itoa(n) + "M"
}
