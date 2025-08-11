//go:build linux
// +build linux

package textarea

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	clear = "\033[2J\033[H"
	grey  = "\033[38;5;241m"
	green = "\033[32m"
	reset = "\033[0m"
)

type TextArea struct {
	title    string
	bodyHint string

	lines [][]rune
	curY  int
	curX  int

	oldState *term.State
}

// NewTextArea creates a new editor.
func NewInput(title, bodyHint string) (*TextArea, error) {
	old, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		return nil, err
	}

	return &TextArea{
		title:    title,
		bodyHint: bodyHint,
		lines:    [][]rune{},
		oldState: old,
	}, nil
}

// Run shows the UI and returns the final lines.
// Ctrl-D => lines, nil
// Ctrl-C => nil, interrupted
func (t *TextArea) Run() ([]string, error) {
	defer t.restore()
	t.redraw()

	for {
		// read single byte (raw mode)
		var buf [1]byte

		if _, err := os.Stdin.Read(buf[:]); err != nil {
			return nil, err
		}

		b := buf[0]

		switch b {
		case 3: // Ctrl-C
			return nil, errors.New("interrupted")
		case 4: // Ctrl-D
			out := make([]string, len(t.lines))
			for i, r := range t.lines {
				out[i] = string(r)
			}
			return out, nil
		case 13: // Enter
			t.handleEnter()
		case 127, 8: // Backspace
			t.handleBackspace()
		case 27: // possible arrow key
			seq := make([]byte, 2)
			if _, err := os.Stdin.Read(seq); err != nil {
				continue
			}
			if seq[0] == '[' {
				t.handleArrow(seq[1])
			}
		default:
			if b >= 32 && b < 127 {
				t.insertRune(rune(b))
			}
		}

		t.redraw()
	}
}

// restore terminal to cooked mode
func (t *TextArea) restore() {
	term.Restore(int(os.Stdin.Fd()), t.oldState)
}

func (t *TextArea) insertRune(r rune) {
	if t.curY >= len(t.lines) {
		t.lines = append(t.lines, []rune{})
	}

	line := t.lines[t.curY]
	left := append([]rune{}, line[:t.curX]...)
	right := append([]rune{}, line[t.curX:]...)

	t.lines[t.curY] = append(append(left, r), right...)
	t.curX++
}

func (t *TextArea) handleBackspace() {
	if t.curX > 0 {
		line := t.lines[t.curY]

		t.lines[t.curY] = append(line[:t.curX-1], line[t.curX:]...)
		t.curX--
	} else if t.curY > 0 {
		// merge to previous line
		prev := t.lines[t.curY-1]

		t.lines[t.curY-1] = append(prev, t.lines[t.curY]...)
		t.lines = append(t.lines[:t.curY], t.lines[t.curY+1:]...)
		t.curY--
		t.curX = len(prev)
	}
}

func (t *TextArea) handleEnter() {
	if t.curY >= len(t.lines) {
		t.lines = append(t.lines, []rune{})
	}
	line := t.lines[t.curY]
	left := append([]rune{}, line[:t.curX]...)
	right := append([]rune{}, line[t.curX:]...)

	t.lines[t.curY] = left
	// insert right under cursor
	t.lines = append(t.lines[:t.curY+1], append([][]rune{right}, t.lines[t.curY+1:]...)...)
	t.curY++
	t.curX = 0
}

func (t *TextArea) handleArrow(ch byte) {
	switch ch {
	case 'A': // up
		if t.curY > 0 {
			t.curY--
			if t.curX > len(t.lines[t.curY]) {
				t.curX = len(t.lines[t.curY])
			}
		}
	case 'B': // down
		if t.curY < len(t.lines)-1 {
			t.curY++
			if t.curX > len(t.lines[t.curY]) {
				t.curX = len(t.lines[t.curY])
			}
		}
	case 'C': // right
		if t.curX < len(t.lines[t.curY]) {
			t.curX++
		}
	case 'D': // left
		if t.curX > 0 {
			t.curX--
		}
	}
}

func (t *TextArea) redraw() {
	fmt.Print(clear)

	// header
	fmt.Println(t.title)
	fmt.Println()

	// content
	for i, line := range t.lines {
		prefix := fmt.Sprintf("%s  │ %s", grey, reset)
		if i == t.curY {
			left := string(line[:t.curX])
			right := string(line[t.curX:])
			cursor := "\033[7m \033[0m"

			fmt.Printf("\r%s%s> %s%s%s%s\n", prefix, green, reset, left, cursor, right)
		} else {
			fmt.Printf("\r%s  %s\n", prefix, string(line))
		}
	}

	// empty line below last
	if t.curY == len(t.lines) {
		fmt.Printf("\r%s  │ %s%s> %s\033[7m \033[0m\n", grey, green, reset, reset)
	}

	fmt.Println()
	fmt.Print("\r\033[38;5;246mCtrl-D to finish  •  Ctrl-C to cancel\033[0m\n")
}
