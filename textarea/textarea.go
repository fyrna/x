package textarea

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fyrna/x/term"
)

const (
	grey      = "\033[38;5;241m"
	green     = "\033[32m"
	red       = "\033[31m"
	blue      = "\033[34m"
	yellow    = "\033[33m"
	magenta   = "\033[35m"
	cyan      = "\033[36m"
	reset     = "\033[0m"
	lightGrey = "\033[38;5;246m"
)

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
	bodyHint     string
	lines        [][]rune
	cursor       struct{ x, y int }
	showLineNr   bool
	history      []Snapshot
	historyPos   int
	searchMode   bool
	searchTerm   []rune
	syntaxMode   SyntaxMode
	highlighters map[SyntaxMode]Highlighter
}

// Snapshot represents a state snapshot for undo/redo
type Snapshot struct {
	lines  [][]rune
	cursor struct{ x, y int }
}

// NewInput creates a new TextArea instance
func NewInput(title, bodyHint string) *TextArea {
	ta := &TextArea{
		title:        title,
		bodyHint:     bodyHint,
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
func (t *TextArea) handleNormalInput(ev term.Event) error {
	switch {
	case ev.IsCtrl('c'):
		return errors.New("interrupted")
	case ev.IsCtrl('d'):
		return ErrFinished
	case ev.IsCtrl('s'):
		t.toggleSearch()
		return nil
	case ev.IsCtrl('z'):
		t.undo()
		return nil
	case ev.IsCtrl('y'):
		t.redo()
		return nil
	case ev.IsCtrl('g'):
		t.toggleSyntaxMode()
		return nil
	case ev.IsCtrl('l'):
		t.toggleLineNumbers()
		return nil
	case ev.Key == term.KeyRune:
		t.takeSnapshot()
		t.insertRune(ev.Rune)
	case ev.Key == term.KeySpace:
		t.takeSnapshot()
		t.insertRune(' ')
	case ev.Key == term.KeyEnter:
		t.takeSnapshot()
		t.handleEnter()
	case ev.Key == term.KeyBackspace:
		t.takeSnapshot()
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

// handleSearchInput processes input in search mode
func (t *TextArea) handleSearchInput(ev term.Event) error {
	switch ev.Key {
	case term.KeyEscape:
		t.searchMode = false
	case term.KeyEnter:
		t.searchMode = false
		t.search(string(t.searchTerm))
		t.searchTerm = nil
	case term.KeyBackspace:
		if len(t.searchTerm) > 0 {
			t.searchTerm = t.searchTerm[:len(t.searchTerm)-1]
		}
	case term.KeyRune:
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

// ** Undo/Redo **
// TODO: undo/redo currently only a basic implementation, it even very buggy lmao!

func (t *TextArea) takeSnapshot() {
	if t.historyPos < len(t.history)-1 {
		t.history = t.history[:t.historyPos+1]
	}

	t.history = append(t.history, Snapshot{
		lines:  t.copyLines(),
		cursor: t.cursor,
	})

	t.historyPos = len(t.history) - 1
}

// i still can't figured out these two, but keep it for "history" XD
func (t *TextArea) undo() {
	if t.historyPos > 0 {
		t.historyPos--
		t.restoreSnapshot()
	}
}

func (t *TextArea) redo() {
	if t.historyPos < len(t.history)-1 {
		t.historyPos++
		t.restoreSnapshot()
	}
}

func (t *TextArea) restoreSnapshot() {
	snap := t.history[t.historyPos]
	t.lines = t.copyLinesSlice(snap.lines)
	t.cursor = snap.cursor
}

func (t *TextArea) copyLines() [][]rune {
	return t.copyLinesSlice(t.lines)
}

func (t *TextArea) copyLinesSlice(src [][]rune) [][]rune {
	dst := make([][]rune, len(src))
	for i, line := range src {
		dst[i] = append([]rune(nil), line...)
	}
	return dst
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
			result.WriteString(reset)
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
		t.renderLine(&buf, i, line)
	}

	// Search box
	if t.searchMode {
		t.renderSearchBox(&buf)
	}

	// Status bar
	t.renderStatusBar(&buf)

	fmt.Print(buf.String())
}

func (t *TextArea) renderLine(buf *bytes.Buffer, lineNum int, line []rune) {
	// Line number
	if t.showLineNr {
		buf.WriteString(grey)
		fmt.Fprintf(buf, " %3d ", lineNum+1)
		buf.WriteString(reset)
		buf.WriteString(" ")
	}

	// Line border
	buf.WriteString(grey)
	buf.WriteString("│ ")
	buf.WriteString(reset)

	// Line content
	if lineNum == t.cursor.y {
		t.renderActiveLine(buf, line)
	} else {
		t.renderInactiveLine(buf, line)
	}
}

func (t *TextArea) renderActiveLine(buf *bytes.Buffer, line []rune) {
	// highlighted := t.highlightLine(line)

	left := string(line[:t.cursor.x])
	right := string(line[t.cursor.x:])

	buf.WriteString(t.highlightLine([]rune(left)))
	buf.WriteString(green)
	buf.WriteString("|")
	buf.WriteString(reset)
	buf.WriteString(t.highlightLine([]rune(right)))
	buf.WriteString("\033[K") // Clear to end of line
}

func (t *TextArea) renderInactiveLine(buf *bytes.Buffer, line []rune) {
	buf.WriteString(t.highlightLine(line))
	buf.WriteString("\033[K") // Clear to end of line
}

func (t *TextArea) renderSearchBox(buf *bytes.Buffer) {
	row := 4 + len(t.lines) + 1
	fmt.Fprintf(os.Stdout, "\033[%d;1H", row)
	buf.WriteString("\033[7mSearch: ")
	buf.WriteString(string(t.searchTerm))
	buf.WriteString("\033[0m\033[K")
}

func (t *TextArea) renderStatusBar(buf *bytes.Buffer) {
	row := 4 + len(t.lines) + 3

	if t.searchMode {
		row += 2
	}

	fmt.Fprintf(os.Stdout, "\033[%d;1H", row)
	buf.WriteString(lightGrey)

	// Left status
	status := fmt.Sprintf(" %s | Ln %d, Col %d ", t.getSyntaxName(), t.cursor.y+1, t.cursor.x+1)
	buf.WriteString(status)

	// Right status
	help := " ^S:Search ^G:Syntax ^L:LineNums ^Z:Undo ^Y:Redo ^D:Finish "
	width := 80
	padding := width - len(status) - len(help) - 1

	if padding > 0 {
		buf.WriteString(strings.Repeat(" ", padding))
	}

	buf.WriteString(help)
	buf.WriteString(reset)
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
		{regexp.MustCompile(`\b(func|package|import|var|const|type|struct|interface|return|if|else|for|range|select|case|default|go|defer)\b`), blue},
		{regexp.MustCompile(`\b(true|false|nil)\b`), yellow},
		{regexp.MustCompile(`\b(int|string|bool|float|rune|byte|error)\b`), cyan},
		{regexp.MustCompile(`//.*`), grey},
		{regexp.MustCompile(`".*?"`), green},
		{regexp.MustCompile("`.*?`"), green},
		{regexp.MustCompile(`'.*?'`), green},
	}
}

// MarkdownHighlighter implements Highlighter for Markdown syntax
type MarkdownHighlighter struct{}

func (h *MarkdownHighlighter) Name() string { return "Markdown" }

func (h *MarkdownHighlighter) Rules() []HighlightRule {
	return []HighlightRule{
		{regexp.MustCompile(`^#+.+`), blue},          // Headers
		{regexp.MustCompile(`\*\*.+?\*\*`), red},     // Bold
		{regexp.MustCompile(`\*.+?\*`), yellow},      // Italic
		{regexp.MustCompile("`.+?`"), green},         // Code
		{regexp.MustCompile(`\[.+?\]\(.+?\)`), cyan}, // Links
		{regexp.MustCompile(`^- .+`), magenta},       // List items
		{regexp.MustCompile(`^\d+\. .+`), magenta},   // Numbered list
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

// ErrFinished is returned when the user finishes input (Ctrl-D)
var ErrFinished = errors.New("finished")
