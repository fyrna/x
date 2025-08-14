package textarea

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fyrna/x/color"
	"github.com/fyrna/x/term"
	"github.com/fyrna/x/term/cursor"
	"github.com/fyrna/x/term/key"
	"github.com/fyrna/x/term/screen"
)

var ErrFinished = errors.New("finished")

// HighlightRule defines a syntax highlighting rule
type HighlightRule struct {
	Pattern *regexp.Regexp
	Color   string
}

// Highlighter defines a syntax highlighter
type Highlighter interface {
	Rules() []HighlightRule
	Name() string
}

// SyntaxMode represents the current syntax highlighting mode
type SyntaxMode int

const (
	ModePlainText SyntaxMode = iota
	ModeGo
	ModeMarkdown
)

// TextArea represents an enhanced text input area
type TextArea struct {
	title        string
	lines        [][]rune
	cursor       struct{ x, y int }
	showLineNr   bool
	searchMode   bool
	searchTerm   []rune
	syntaxMode   SyntaxMode
	highlighters map[SyntaxMode]Highlighter
}

// NewInput creates a new TextArea instance
func NewInput(title string) *TextArea {
	ta := &TextArea{
		title:        title,
		lines:        [][]rune{{}},
		showLineNr:   true,
		highlighters: make(map[SyntaxMode]Highlighter),
	}

	// Register built-in highlighters
	ta.highlighters[ModeGo] = &GoHighlighter{}
	ta.highlighters[ModeMarkdown] = &MarkdownHighlighter{}

	return ta
}

// Run starts the text input session
func (t *TextArea) Run() ([]string, error) {
	terminal := term.NewStdinTerminal()
	err := terminal.MakeRaw()
	if err != nil {
		return nil, err
	}
	defer terminal.Restore()

	keyReader := key.NewReader(terminal)

	for {
		t.redraw()

		ev, err := keyReader.ReadEvent()
		if err != nil {
			return nil, err
		}

		if t.searchMode {
			if err := t.handleSearchInput(ev); err != nil {
				return nil, err
			}
			continue
		}

		if err := t.handleNormalInput(ev); err != nil {
			if errors.Is(err, ErrFinished) {
				return t.finish(), nil
			}
			return nil, err
		}
	}
}

// handleNormalInput processes input in normal mode
func (t *TextArea) handleNormalInput(ev key.Event) error {
	switch {
	case ev.IsCtrl('c'):
		return errors.New("interrupted")
	case ev.IsCtrl('d'):
		return ErrFinished
	case ev.IsCtrl('s'):
		t.toggleSearch()
		return nil
	case ev.IsCtrl('g'):
		t.toggleSyntaxMode()
		return nil
	case ev.IsCtrl('l'):
		t.toggleLineNumbers()
		return nil
	case ev.Key == key.Rune:
		t.insertRune(ev.Rune)
	case ev.Key == key.Space:
		t.insertRune(' ')
	case ev.Key == key.Enter:
		t.handleEnter()
	case ev.Key == key.Backspace:
		t.handleBackspace()

	case ev.Key == key.Up:
		t.moveUp()
	case ev.Key == key.Down:
		t.moveDown()
	case ev.Key == key.Left:
		t.moveLeft()
	case ev.Key == key.Right:
		t.moveRight()
	}
	return nil
}

// handleSearchInput processes input in search mode
func (t *TextArea) handleSearchInput(ev key.Event) error {
	switch ev.Key {
	case key.Escape:
		t.searchMode = false
	case key.Enter:
		t.searchMode = false
		t.search(string(t.searchTerm))
		t.searchTerm = nil
	case key.Backspace:
		if len(t.searchTerm) > 0 {
			t.searchTerm = t.searchTerm[:len(t.searchTerm)-1]
		}
	case key.Rune:
		t.searchTerm = append(t.searchTerm, ev.Rune)
	}
	return nil
}

// ** Core Editing Functions **

func (t *TextArea) insertRune(r rune) {
	if t.cursor.y >= len(t.lines) {
		t.lines = append(t.lines, []rune{})
	}

	// create new slice WITHOUT modification to real slice
	newLine := make([]rune, 0, len(t.lines[t.cursor.y])+1)
	newLine = append(newLine, t.lines[t.cursor.y][:t.cursor.x]...)
	newLine = append(newLine, r)
	newLine = append(newLine, t.lines[t.cursor.y][t.cursor.x:]...)

	t.lines[t.cursor.y] = newLine
	t.cursor.x++
}

func (t *TextArea) handleBackspace() {
	if t.cursor.x > 0 {
		// Buat line baru tanpa karakter yang dihapus
		newLine := make([]rune, 0, len(t.lines[t.cursor.y])-1)
		newLine = append(newLine, t.lines[t.cursor.y][:t.cursor.x-1]...)
		newLine = append(newLine, t.lines[t.cursor.y][t.cursor.x:]...)

		t.lines[t.cursor.y] = newLine
		t.cursor.x--
	} else if t.cursor.y > 0 {
		prev := make([]rune, len(t.lines[t.cursor.y-1]))
		copy(prev, t.lines[t.cursor.y-1])

		t.lines[t.cursor.y-1] = append(prev, t.lines[t.cursor.y]...)
		t.lines = append(t.lines[:t.cursor.y], t.lines[t.cursor.y+1:]...)
		t.cursor.y--
		t.cursor.x = len(prev)
	}
}

func (t *TextArea) handleEnter() {
	line := t.lines[t.cursor.y]
	left := line[:t.cursor.x]
	right := line[t.cursor.x:]

	t.lines[t.cursor.y] = left
	t.lines = append(t.lines[:t.cursor.y+1], append([][]rune{right}, t.lines[t.cursor.y+1:]...)...)

	t.cursor.y++
	t.cursor.x = 0
}

// ** Cursor Movement **

func (t *TextArea) moveUp() {
	if t.cursor.y > 0 {
		t.cursor.y--
		t.snapX()
	}
}

func (t *TextArea) moveDown() {
	if t.cursor.y < len(t.lines)-1 {
		t.cursor.y++
		t.snapX()
	}
}

func (t *TextArea) moveLeft() {
	if t.cursor.x > 0 {
		t.cursor.x--
	}
}

func (t *TextArea) moveRight() {
	currentLineLen := len(t.lines[t.cursor.y])
	if t.cursor.x < currentLineLen {
		t.cursor.x++
	}
}

func (t *TextArea) snapX() {
	currentLineLen := len(t.lines[t.cursor.y])
	if t.cursor.x > currentLineLen {
		t.cursor.x = currentLineLen
	}
}

// ** Search **
// TODO: better UI

func (t *TextArea) toggleSearch() {
	t.searchMode = !t.searchMode
	if !t.searchMode {
		t.searchTerm = nil
	}
}

func (t *TextArea) search(term string) {
	if term == "" {
		return
	}

	for y := t.cursor.y; y < len(t.lines); y++ {
		line := string(t.lines[y])
		if idx := strings.Index(line, term); idx != -1 {
			t.cursor.y = y
			t.cursor.x = idx
			return
		}
	}

	// Wrap around if not found
	for y := 0; y < t.cursor.y; y++ {
		line := string(t.lines[y])
		if idx := strings.Index(line, term); idx != -1 {
			t.cursor.y = y
			t.cursor.x = idx
			return
		}
	}
}

// ** Syntax Highlighting **
// TODO: basically as far as i can see it works good

func (t *TextArea) toggleSyntaxMode() {
	t.syntaxMode = (t.syntaxMode + 1) % 3
}

func (t *TextArea) highlightLine(line []rune) string {
	if t.syntaxMode == ModePlainText {
		return string(line)
	}

	highlighter, exists := t.highlighters[t.syntaxMode]
	if !exists {
		return string(line)
	}

	strLine := string(line)
	var result strings.Builder

	for _, rule := range highlighter.Rules() {
		matches := rule.Pattern.FindAllStringIndex(strLine, -1)

		if matches == nil {
			continue
		}

		lastPos := 0

		for _, match := range matches {
			// Text before  match
			result.WriteString(strLine[lastPos:match[0]])
			// Text after match with color
			result.WriteString(rule.Color)
			result.WriteString(strLine[match[0]:match[1]])
			result.WriteString(color.Reset)
			lastPos = match[1]
		}

		strLine = strLine[lastPos:]
	}

	// the rest not match
	result.WriteString(strLine)
	return result.String()
}

// ** Line Numbers **

func (t *TextArea) toggleLineNumbers() {
	t.showLineNr = !t.showLineNr
}

// ** Rendering **

func (t *TextArea) redraw() {
	screen.Clear()

	// Title
	os.Stdout.WriteString(t.title + "\n\n")

	// Content lines
	for i, line := range t.lines {
		cursor.To(4+i, 1) // Position each line
		t.renderLine(i, line)
	}

	// Search box
	if t.searchMode {
		t.renderSearchBox()
	}

	// Status bar
	t.renderStatusBar()
}

func (t *TextArea) renderLine(lineNum int, line []rune) {
	// Line number
	if t.showLineNr {
		os.Stdout.WriteString(color.Fg256(241))
		fmt.Fprintf(os.Stdout, " %3d ", lineNum+1)
		os.Stdout.WriteString(color.Reset)
		os.Stdout.WriteString(" ")
	}

	// Line border
	os.Stdout.WriteString(color.Fg256(241))
	os.Stdout.WriteString("│ ")
	os.Stdout.WriteString(color.Reset)

	// Line content
	if lineNum == t.cursor.y {
		t.renderActiveLine(line)
	} else {
		t.renderInactiveLine(line)
	}
}

func (t *TextArea) renderActiveLine(line []rune) {
	// highlighted := t.highlightLine(line)

	left := string(line[:t.cursor.x])
	right := string(line[t.cursor.x:])

	os.Stdout.WriteString(t.highlightLine([]rune(left)))
	os.Stdout.WriteString(color.Green)
	os.Stdout.WriteString("|")
	os.Stdout.WriteString(color.Reset)
	os.Stdout.WriteString(t.highlightLine([]rune(right)))
	screen.ClearLine(0)
}

func (t *TextArea) renderInactiveLine(line []rune) {
	os.Stdout.WriteString(t.highlightLine(line))
	screen.ClearLine(0)
}

func (t *TextArea) renderSearchBox() {
	row := 4 + len(t.lines) + 1

	cursor.To(row, 1)

	os.Stdout.WriteString("\033[7mSearch: ")
	os.Stdout.WriteString(string(t.searchTerm))

	screen.ClearLine(0)
}

func (t *TextArea) renderStatusBar() {
	row := 4 + len(t.lines) + 3

	if t.searchMode {
		row += 2
	}

	cursor.To(row, 1)
	// os.Stdout.WriteString(fmt.Sprintf("\033[%d;1H", row))
	os.Stdout.WriteString(color.Fg256(246))

	// Left status
	status := fmt.Sprintf(" %s | Ln %d, Col %d ", t.getSyntaxName(), t.cursor.y+1, t.cursor.x+1)
	os.Stdout.WriteString(status)

	// Right status
	help := " ^S:Search ^G:Syntax ^L:LineNums ^D:Finish "
	width := 65
	padding := width - len(status) - len(help) - 1

	if padding > 0 {
		os.Stdout.WriteString(strings.Repeat(" ", padding))
	}

	os.Stdout.WriteString(help)
	os.Stdout.WriteString(color.Reset)
}

func (t *TextArea) getSyntaxName() string {
	switch t.syntaxMode {
	case ModePlainText:
		return "Text"
	case ModeGo:
		return "Go"
	case ModeMarkdown:
		return "Markdown"
	default:
		return ""
	}
}

// Built-in Highlighters
// regexp is cute for this :)

// GoHighlighter implements Highlighter for Go syntax
type GoHighlighter struct{}

func (h *GoHighlighter) Name() string { return "Go" }

func (h *GoHighlighter) Rules() []HighlightRule {
	return []HighlightRule{
		{regexp.MustCompile(`\b(func|package|import|var|const|type|struct|interface|return|if|else|for|range|select|case|default|go|defer)\b`), color.Blue},
		{regexp.MustCompile(`\b(true|false|nil)\b`), color.Yellow},
		{regexp.MustCompile(`\b(int|string|bool|float|rune|byte|error)\b`), color.Cyan},
		{regexp.MustCompile(`//.*`), color.Fg256(241)},
		{regexp.MustCompile(`".*?"`), color.Green},
		{regexp.MustCompile("`.*?`"), color.Green},
		{regexp.MustCompile(`'.*?'`), color.Green},
	}
}

// MarkdownHighlighter implements Highlighter for Markdown syntax
type MarkdownHighlighter struct{}

func (h *MarkdownHighlighter) Name() string { return "Markdown" }

func (h *MarkdownHighlighter) Rules() []HighlightRule {
	return []HighlightRule{
		{regexp.MustCompile(`^#+.+`), color.Blue},          // Headers
		{regexp.MustCompile(`\*\*.+?\*\*`), color.Red},     // Bold
		{regexp.MustCompile(`\*.+?\*`), color.Yellow},      // Italic
		{regexp.MustCompile("`.+?`"), color.Green},         // Code
		{regexp.MustCompile(`\[.+?\]\(.+?\)`), color.Cyan}, // Links
		{regexp.MustCompile(`^- .+`), color.Magenta},       // List items
		{regexp.MustCompile(`^\d+\. .+`), color.Magenta},   // Numbered list
	}
}

// Utility

func (t *TextArea) finish() []string {
	result := make([]string, len(t.lines))
	for i, line := range t.lines {
		result[i] = string(line)
	}
	return result
}
