// this is for key event, i don't want to you know, hard code this
// and maybe i guess my implementation is stupid
// i mean yeah i do this stuff for learning though.

package key

import (
	"errors"
	"os"
	"unicode/utf8"

	"github.com/fyrna/x/term"
)

type Key uint16

const (
	Unknown Key = iota
	Rune        // Normal character, human stuff i guess?

	Escape
	Enter
	Backspace
	Tab
	Space

	// Navigation
	Up
	Down
	Left
	Right
	Home
	End
	PgUp
	PgDown

	// Function keys
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
)

// Modifier bitset
type Mod uint8

const (
	ModNone  Mod = 0
	ModShift Mod = 1 << iota
	ModAlt
	ModCtrl
)

// Event represents a key event
type Event struct {
	Key   Key
	Mod   Mod
	Rune  rune
	Width int // Visual width of rune (for CJK)
}

// Reader handles keyboard input
type Reader struct {
	term *term.Terminal
}

// NewReader creates a new keyboard reader for the given terminal
func NewReader(t *term.Terminal) *Reader {
	return &Reader{term: t}
}

// ReadEvent reads keyboard event (blocking)
// Requires MakeKeyboardRaw() to be called first
func (r *Reader) ReadEvent() (Event, error) {
	if r.term == nil {
		return Event{}, errors.New("terminal not initialized")
	}

	b, err := r.readByte()
	if err != nil {
		return Event{}, err
	}

	return r.parse(b)
}

func (ev Event) IsCtrl(r rune) bool {
	return ev.Key == Rune && ev.Mod&ModCtrl != 0 && ev.Rune == r
}

func (ev Event) IsAlt(r rune) bool {
	return ev.Key == Rune && ev.Mod&ModAlt != 0 && ev.Rune == r
}
func (ev Event) IsShift(r rune) bool {
	return ev.Key == Rune && ev.Mod&ModShift != 0 && ev.Rune == r
}

// String returns string representation of event
func (ev Event) String() string {
	var mod string
	if ev.Mod&ModCtrl != 0 {
		mod += "Ctrl+"
	}
	if ev.Mod&ModAlt != 0 {
		mod += "Alt+"
	}
	if ev.Mod&ModShift != 0 {
		mod += "Shift+"
	}
	if ev.Key == Rune {
		return mod + string(ev.Rune)
	}

	keyNames := map[Key]string{
		Escape:    "Escape",
		Enter:     "Enter",
		Backspace: "Backspace",
		Tab:       "Tab",
		Space:     "Space",
		Up:        "Up",
		Down:      "Down",
		Left:      "Left",
		Right:     "Right",
		Home:      "Home",
		End:       "End",
		PgUp:      "PgUp",
		PgDown:    "PgDown",
		F1:        "F1",
		F2:        "F2",
		F3:        "F3",
		F4:        "F4",
		F5:        "F5",
		F6:        "F6",
		F7:        "F7",
		F8:        "F8",
		F9:        "F9",
		F10:       "F10",
		F11:       "F11",
		F12:       "F12",
	}

	if name, ok := keyNames[ev.Key]; ok {
		return mod + name
	}

	return mod + "Unknown"
}

func (r *Reader) readByte() (byte, error) {
	var buf [1]byte
	_, err := os.NewFile(uintptr(r.term.Fd()), "terminal").Read(buf[:])
	return buf[0], err
}

func (r *Reader) parse(b byte) (Event, error) {
	switch b {
	case 0x1b: // ESC
		return r.parseEsc()
	case 0x7f: // DEL
		return Event{Key: Backspace}, nil
	case 0x09: // TAB
		return Event{Key: Tab}, nil
	case 0x0d, 0x0a: // CR / LF
		return Event{Key: Enter}, nil
	case 0x20:
		return Event{Key: Space}, nil
	default:
		// Ctrl-letter (0x01–0x1a)
		if b >= 0x01 && b <= 0x1a {
			return Event{Key: Rune, Rune: rune(b + 96), Mod: ModCtrl}, nil
		}

		// printable
		if b >= 0x20 && b < 0x7f {
			r, _ := utf8.DecodeRune([]byte{b})
			w := 1
			if r >= 0x1100 { // crude CJK width
				w = 2
			}
			return Event{Key: Rune, Rune: r, Width: w}, nil
		}
		return Event{Key: Unknown}, nil
	}
}

func (r *Reader) parseEsc() (Event, error) {
	seq, err := r.readByte()
	if err != nil {
		return Event{Key: Escape}, nil // lone ESC
	}

	if seq != '[' && seq != 'O' {
		// Alt+<key> (ESC + key)
		return Event{Key: Rune, Rune: rune(seq), Mod: ModAlt}, nil
	}

	// Read the next byte for sequence
	b2, err := r.readByte()
	if err != nil {
		return Event{Key: Escape}, nil
	}

	// Handle multi-byte sequences
	switch string([]byte{seq, b2}) {
	case "[A":
		return Event{Key: Up}, nil
	case "[B":
		return Event{Key: Down}, nil
	case "[C":
		return Event{Key: Right}, nil
	case "[D":
		return Event{Key: Left}, nil
	case "[H":
		return Event{Key: Home}, nil
	case "[F":
		return Event{Key: End}, nil
	case "OP":
		return Event{Key: F1}, nil
	case "OQ":
		return Event{Key: F2}, nil
	case "OR":
		return Event{Key: F3}, nil
	case "OS":
		return Event{Key: F4}, nil
	}

	// Handle longer sequences
	if seq == '[' && b2 >= '0' && b2 <= '9' {
		b3, err := r.readByte()
		if err != nil {
			return Event{Key: Unknown}, err
		}

		switch string([]byte{seq, b2, b3}) {
		case "[5~":
			return Event{Key: PgUp}, nil
		case "[6~":
			return Event{Key: PgDown}, nil
		case "[15~":
			return Event{Key: F5}, nil
		case "[17~":
			return Event{Key: F6}, nil
		case "[18~":
			return Event{Key: F7}, nil
		case "[19~":
			return Event{Key: F8}, nil
		case "[20~":
			return Event{Key: F9}, nil
		case "[21~":
			return Event{Key: F10}, nil
		case "[23~":
			return Event{Key: F11}, nil
		case "[24~":
			return Event{Key: F12}, nil
		}
	}

	return Event{Key: Unknown}, nil
}
