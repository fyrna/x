// this is for key event, i don't want to you know, hard code this
// and maybe i guess my implementation is stupid
// i mean yeah i do this stuff for learning though.

package term

import (
	"errors"
	"os"
	"unicode/utf8"

	"golang.org/x/term"
)

type Key uint16

const (
	KeyUnknown Key = iota
	KeyRune        // Normal character, human stuff i guess?

	KeyEscape
	KeyEnter
	KeyBackspace
	KeyTab
	KeySpace

	// Arrows
	KeyUp
	KeyDown
	KeyLeft
	KeyRight

	// Navigation
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDown

	// Function keys
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
)

// Modifier bitset
type Mod uint8

const (
	KeyModNone  Mod = 0
	KeyModShift Mod = 1 << iota
	KeyModAlt
	KeyModCtrl
)

// Event represents a key event
type Event struct {
	Key   Key
	Mod   Mod
	Rune  rune
	Width int // Visual width of rune (for CJK)
}

type reader struct{ fd int }

var (
	// singleton after Init()
	stdinReader *reader
	oldState    *term.State
)

// Init raw mode, call once.
func Init() error {
	if stdinReader != nil {
		return nil
	}

	var err error

	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	stdinReader = &reader{fd: int(os.Stdin.Fd())}
	return nil
}

// Restore()
func Restore() error {
	if oldState == nil {
		return nil
	}

	err := term.Restore(int(os.Stdin.Fd()), oldState)
	oldState = nil
	stdinReader = nil

	return err
}

// Read() blocking high-level API
func Read() (Event, error) {
	if stdinReader == nil {
		return Event{}, errors.New("call term.Init() first")
	}

	b, err := stdinReader.readByte()
	if err != nil {
		return Event{}, err
	}

	return stdinReader.parse(b)
}

func (r *reader) readByte() (byte, error) {
	var buf [1]byte
	_, err := os.Stdin.Read(buf[:])
	return buf[0], err
}

func (r *reader) parse(b byte) (Event, error) {
	switch b {
	case 0x1b: // ESC
		return r.parseEsc()
	case 0x7f: // DEL
		return Event{Key: KeyBackspace}, nil
	case 0x09: // TAB
		return Event{Key: KeyTab}, nil
	case 0x0d, 0x0a: // CR / LF
		return Event{Key: KeyEnter}, nil
	case 0x20:
		return Event{Key: KeySpace}, nil
	default:
		// Ctrl-letter (0x01–0x1a)
		if b >= 0x01 && b <= 0x1a {
			return Event{Key: KeyRune, Rune: rune(b + 96), Mod: KeyModCtrl}, nil
		}
		// printable
		if b >= 0x20 && b < 0x7f {
			r, _ := utf8.DecodeRune([]byte{b})
			w := 1
			if r >= 0x1100 { // crude CJK width
				w = 2
			}
			return Event{Key: KeyRune, Rune: r, Width: w}, nil
		}
		return Event{Key: KeyUnknown}, nil
	}
}

func (r *reader) parseEsc() (Event, error) {
	seq, err := r.readByte()

	if err != nil {
		return Event{Key: KeyEscape}, nil // lone ESC
	}

	if seq != '[' && seq != 'O' {
		// Alt+<key> (ESC + key)
		return Event{Key: KeyRune, Rune: rune(seq), Mod: KeyModAlt}, nil
	}

	// CSI / SS3
	b2, _ := r.readByte()
	switch string([]byte{seq, b2}) {
	case "[A":
		return Event{Key: KeyUp}, nil
	case "[B":
		return Event{Key: KeyDown}, nil
	case "[C":
		return Event{Key: KeyRight}, nil
	case "[D":
		return Event{Key: KeyLeft}, nil
	case "[H":
		return Event{Key: KeyHome}, nil
	case "[F":
		return Event{Key: KeyEnd}, nil
	case "[5~":
		return Event{Key: KeyPgUp}, nil
	case "[6~":
		return Event{Key: KeyPgDown}, nil
	case "OP":
		return Event{Key: KeyF1}, nil
	case "OQ":
		return Event{Key: KeyF2}, nil
	case "OR":
		return Event{Key: KeyF3}, nil
	case "OS":
		return Event{Key: KeyF4}, nil
	case "[15~":
		return Event{Key: KeyF5}, nil
	case "[17~":
		return Event{Key: KeyF6}, nil
	case "[18~":
		return Event{Key: KeyF7}, nil
	case "[19~":
		return Event{Key: KeyF8}, nil
	case "[20~":
		return Event{Key: KeyF9}, nil
	case "[21~":
		return Event{Key: KeyF10}, nil
	case "[23~":
		return Event{Key: KeyF11}, nil
	case "[24~":
		return Event{Key: KeyF12}, nil
	}
	return Event{Key: KeyUnknown}, nil
}
