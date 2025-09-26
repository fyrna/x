package textarea

import (
	"errors"
	"os"

	"github.com/fyrna/x/term"
	"github.com/fyrna/x/term/cursor"
	"github.com/fyrna/x/term/key"
	"github.com/fyrna/x/term/screen"
)

var ErrFinished = errors.New("finished")

type TextArea struct {
	title  string
	lines  [][]rune
	cursor struct{ x, y int }
}

func New(title string) *TextArea {
	return &TextArea{
		title: title,
		lines: [][]rune{{}}, // one empty line
	}
}

func (t *TextArea) Run() ([]string, error) {
	terminal := term.NewStdinTerminal()
	if err := terminal.MakeRaw(); err != nil {
		return nil, err
	}
	defer terminal.Restore()

	reader := key.NewReader(terminal)

	cursor.Hide()
	defer cursor.Show()

	for {
		t.redraw()

		ev, err := reader.ReadEvent()
		if err != nil {
			return nil, err
		}

		switch {
		case ev.IsCtrl('d'):
			return t.finish(), nil
		case ev.IsCtrl('c'):
			return nil, errors.New("interrupted")
		case ev.Key == key.Rune:
			t.insertRune(ev.Rune)
		case ev.Key == key.Space:
			t.insertRune(' ')
		case ev.Key == key.Enter:
			t.handleEnter()
		case ev.Key == key.Backspace:
			t.handleBackspace()
		case ev.Key == key.Left:
			t.moveLeft()
		case ev.Key == key.Right:
			t.moveRight()
		case ev.Key == key.Up:
			t.moveUp()
		case ev.Key == key.Down:
			t.moveDown()
		}
	}
}

func (t *TextArea) redraw() {
	screen.Clear()

	os.Stdout.WriteString(t.title + "\n\n")

	for i, line := range t.lines {
		cursor.To(3+i, 1)
		if i == t.cursor.y {
			left := string(line[:t.cursor.x])
			right := string(line[t.cursor.x:])

			os.Stdout.WriteString(left)
			os.Stdout.WriteString("█")
			os.Stdout.WriteString(right)
		} else {
			os.Stdout.WriteString(string(line))
		}
	}
}

func (t *TextArea) insertRune(r rune) {
	line := t.lines[t.cursor.y]
	newLine := append([]rune{}, line[:t.cursor.x]...)
	newLine = append(newLine, r)
	newLine = append(newLine, line[t.cursor.x:]...)

	t.lines[t.cursor.y] = newLine
	t.cursor.x++
}

func (t *TextArea) handleEnter() {
	line := t.lines[t.cursor.y]
	left := append([]rune{}, line[:t.cursor.x]...)
	right := append([]rune{}, line[t.cursor.x:]...)

	t.lines[t.cursor.y] = left
	t.lines = append(t.lines[:t.cursor.y+1],
		append([][]rune{right}, t.lines[t.cursor.y+1:]...)...)

	t.cursor.y++
	t.cursor.x = 0
}

func (t *TextArea) handleBackspace() {
	if t.cursor.x > 0 {
		line := t.lines[t.cursor.y]
		newLine := append([]rune{}, line[:t.cursor.x-1]...)
		newLine = append(newLine, line[t.cursor.x:]...)
		t.lines[t.cursor.y] = newLine
		t.cursor.x--
	} else if t.cursor.y > 0 {
		prev := t.lines[t.cursor.y-1]
		curr := t.lines[t.cursor.y]
		t.lines[t.cursor.y-1] = append(prev, curr...)
		t.lines = append(t.lines[:t.cursor.y], t.lines[t.cursor.y+1:]...)
		t.cursor.y--
		t.cursor.x = len(prev)
	}
}

func (t *TextArea) moveLeft() {
	if t.cursor.x > 0 {
		t.cursor.x--
	} else if t.cursor.y > 0 {
		t.cursor.y--
		t.cursor.x = len(t.lines[t.cursor.y])
	}
}

func (t *TextArea) moveRight() {
	if t.cursor.x < len(t.lines[t.cursor.y]) {
		t.cursor.x++
	} else if t.cursor.y < len(t.lines)-1 {
		t.cursor.y++
		t.cursor.x = 0
	}
}

func (t *TextArea) moveUp() {
	if t.cursor.y > 0 {
		t.cursor.y--
		if t.cursor.x > len(t.lines[t.cursor.y]) {
			t.cursor.x = len(t.lines[t.cursor.y])
		}
	}
}

func (t *TextArea) moveDown() {
	if t.cursor.y < len(t.lines)-1 {
		t.cursor.y++
		if t.cursor.x > len(t.lines[t.cursor.y]) {
			t.cursor.x = len(t.lines[t.cursor.y])
		}
	}
}

func (t *TextArea) finish() []string {
	out := make([]string, len(t.lines))
	for i, line := range t.lines {
		out[i] = string(line)
	}
	return out
}
