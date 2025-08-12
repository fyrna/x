package textarea

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/fyrna/x/term"
)

const (
	grey      = "\033[38;5;241m"
	green     = "\033[32m"
	reset     = "\033[0m"
	lightGrey = "\033[38;5;246m"
)

// TextArea represents a multi-line text input area in the terminal
type TextArea struct {
	title    string
	bodyHint string
	lines    [][]rune
	cursor   struct {
		x, y int
	}
}

// NewInput creates a new TextArea instance
func NewInput(title, bodyHint string) *TextArea {
	return &TextArea{
		title:    title,
		bodyHint: bodyHint,
		lines:    [][]rune{{}},
	}
}

// Run starts the text input session
func (t *TextArea) Run() ([]string, error) {
	if err := term.Init(); err != nil {
		return nil, err
	}
	defer term.Restore()

	for {
		t.redraw()

		ev, err := term.Read()
		if err != nil {
			return nil, err
		}

		if err := t.handleInput(ev); err != nil {
			if errors.Is(err, ErrFinished) {
				return t.finish(), nil
			}
			return nil, err
		}
	}
}

// finish returns the current text as a slice of strings
func (t *TextArea) finish() []string {
	result := make([]string, len(t.lines))
	for i, line := range t.lines {
		result[i] = string(line)
	}
	return result
}

// handleInput processes a single input event
func (t *TextArea) handleInput(ev term.Event) error {
	switch {
	case ev.IsCtrl('c'):
		return errors.New("interrupted")
	case ev.IsCtrl('d'):
		return ErrFinished
	case ev.Key == term.KeyRune:
		t.insertRune(ev.Rune)
	case ev.Key == term.KeySpace:
		t.insertRune(' ')
	case ev.Key == term.KeyEnter:
		t.handleEnter()
	case ev.Key == term.KeyBackspace:
		t.handleBackspace()
	case ev.Key == term.KeyUp:
		t.moveUp()
	case ev.Key == term.KeyDown:
		t.moveDown()
	case ev.Key == term.KeyLeft:
		t.moveLeft()
	case ev.Key == term.KeyRight:
		t.moveRight()
	}
	return nil
}

// insertRune inserts a rune at the current cursor position
func (t *TextArea) insertRune(r rune) {
	if t.cursor.y >= len(t.lines) {
		t.lines = append(t.lines, []rune{})
	}

	line := t.lines[t.cursor.y]
	left := line[:t.cursor.x]
	right := line[t.cursor.x:]
	t.lines[t.cursor.y] = append(append(left, r), right...)
	t.cursor.x++
}

// handleBackspace handles backspace key press
func (t *TextArea) handleBackspace() {
	if t.cursor.x > 0 {
		// Delete within current line
		line := t.lines[t.cursor.y]
		t.lines[t.cursor.y] = append(line[:t.cursor.x-1], line[t.cursor.x:]...)
		t.cursor.x--
	} else if t.cursor.y > 0 {
		// Merge with previous line
		prev := t.lines[t.cursor.y-1]
		t.lines[t.cursor.y-1] = append(prev, t.lines[t.cursor.y]...)
		t.lines = append(t.lines[:t.cursor.y], t.lines[t.cursor.y+1:]...)
		t.cursor.y--
		t.cursor.x = len(prev)
	}
}

// handleEnter handles enter/return key press
func (t *TextArea) handleEnter() {
	line := t.lines[t.cursor.y]
	left := line[:t.cursor.x]
	right := line[t.cursor.x:]

	t.lines[t.cursor.y] = left
	t.lines = append(t.lines[:t.cursor.y+1], append([][]rune{right}, t.lines[t.cursor.y+1:]...)...)

	t.cursor.y++
	t.cursor.x = 0
}

// moveUp moves the cursor up
func (t *TextArea) moveUp() {
	if t.cursor.y > 0 {
		t.cursor.y--
		t.snapX()
	}
}

// moveDown moves the cursor down
func (t *TextArea) moveDown() {
	if t.cursor.y < len(t.lines)-1 {
		t.cursor.y++
		t.snapX()
	}
}

// moveLeft moves the cursor left
func (t *TextArea) moveLeft() {
	if t.cursor.x > 0 {
		t.cursor.x--
	}
}

// moveRight moves the cursor right
func (t *TextArea) moveRight() {
	currentLineLen := len(t.lines[t.cursor.y])
	if t.cursor.x < currentLineLen {
		t.cursor.x++
	}
}

// snapX ensures the cursor X position is within bounds of the current line
func (t *TextArea) snapX() {
	currentLineLen := len(t.lines[t.cursor.y])
	if t.cursor.x > currentLineLen {
		t.cursor.x = currentLineLen
	}
}

// redraw renders the text area to the terminal
func (t *TextArea) redraw() {
	var buf bytes.Buffer

	// Clear screen and move to top-left
	buf.WriteString("\033[H\033[2J")

	// Title
	buf.WriteString(t.title)
	buf.WriteString("\n\n")

	// Body hint
	buf.WriteString("  ")
	buf.WriteString(t.bodyHint)
	buf.WriteString("\n")

	// Content lines
	for i, line := range t.lines {
		buf.WriteString(fmt.Sprintf("\033[%d;1H", 4+i)) // Position each line

		if i == t.cursor.y {
			t.renderActiveLine(&buf, line)
		} else {
			t.renderInactiveLine(&buf, line)
		}
	}

	// Footer
	footerRow := 5 + len(t.lines)
	buf.WriteString(fmt.Sprintf("\033[%d;1H", footerRow))
	buf.WriteString(lightGrey)
	buf.WriteString("Ctrl-D to finish  •  Ctrl-C to cancel")
	buf.WriteString(reset)

	fmt.Print(buf.String())
}

// renderActiveLine renders a line with cursor
func (t *TextArea) renderActiveLine(buf *bytes.Buffer, line []rune) {
	buf.WriteString(grey)
	buf.WriteString("  │ ")
	buf.WriteString(green)
	buf.WriteString("> ")
	buf.WriteString(reset)

	left := string(line[:t.cursor.x])
	right := string(line[t.cursor.x:])

	buf.WriteString(left)
	buf.WriteString("\033[7m \033[0m") // Cursor
	buf.WriteString(right)
	buf.WriteString("\033[K") // Clear to end of line
}

// renderInactiveLine renders a normal line
func (t *TextArea) renderInactiveLine(buf *bytes.Buffer, line []rune) {
	buf.WriteString(grey)
	buf.WriteString("  │ ")
	buf.WriteString(reset)
	buf.WriteString("  ")
	buf.WriteString(string(line))
	buf.WriteString("\033[K") // Clear to end of line
}

// ErrFinished is returned when the user finishes input (Ctrl-D)
var ErrFinished = errors.New("finished")
