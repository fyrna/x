// this is for key event, i don't want to you know, hard code this
// and maybe i guess my implementation is stupid
// i mean yeah i do this stuff for learning though.

package term

import (
	"errors"
	"os"
	"unicode/utf8"
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

	// Navigation
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
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

type keyReader struct {
	fd int
}

// KeyboardState holds keyboard raw mode state
type KeyboardState struct {
	termState *State
	reader    *keyReader
}

func MakeKeyboardRaw() (*KeyboardState, error) {
	return MakeKeyboardRawFd(DefaultFd)
}

func MakeKeyboardRawFd(fd int) (*KeyboardState, error) {
	termState, err := MakeRaw(fd)
	if err != nil {
		return nil, err
	}

	return &KeyboardState{
		termState: termState,
		reader:    &keyReader{fd: fd},
	}, nil
}

// RestoreKeyboard restores keyboard
func RestoreKeyboard(state *KeyboardState) error {
	if state == nil || state.termState == nil {
		return nil
	}

	err := Restore(state.reader.fd, state.termState)
	state.termState = nil // Clear state
	return err
}

// ReadEvent reads keyboard event (blocking)
// Requires MakeKeyboardRaw() to be called first
func ReadEvent(state *KeyboardState) (Event, error) {
	if state == nil || state.reader == nil {
		return Event{}, errors.New("keyboard not in raw mode, call MakeKeyboardRaw() first")
	}

	b, err := state.reader.readByte()
	if err != nil {
		return Event{}, err
	}

	return state.reader.parse(b)
}

func (ev Event) IsCtrl(r rune) bool {
	return ev.Key == KeyRune && ev.Mod&ModCtrl != 0 && ev.Rune == r
}

func (ev Event) IsAlt(r rune) bool {
	return ev.Key == KeyRune && ev.Mod&ModAlt != 0 && ev.Rune == r
}
func (ev Event) IsShift(r rune) bool {
	return ev.Key == KeyRune && ev.Mod&ModShift != 0 && ev.Rune == r
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
	if ev.Key == KeyRune {
		return mod + string(ev.Rune)
	}

	keyNames := map[Key]string{
		KeyEscape:    "Escape",
		KeyEnter:     "Enter",
		KeyBackspace: "Backspace",
		KeyTab:       "Tab",
		KeySpace:     "Space",
		KeyUp:        "Up",
		KeyDown:      "Down",
		KeyLeft:      "Left",
		KeyRight:     "Right",
		KeyHome:      "Home",
		KeyEnd:       "End",
		KeyPgUp:      "PgUp",
		KeyPgDown:    "PgDown",
		KeyF1:        "F1",
		KeyF2:        "F2",
		KeyF3:        "F3",
		KeyF4:        "F4",
		KeyF5:        "F5",
		KeyF6:        "F6",
		KeyF7:        "F7",
		KeyF8:        "F8",
		KeyF9:        "F9",
		KeyF10:       "F10",
		KeyF11:       "F11",
		KeyF12:       "F12",
	}

	if name, ok := keyNames[ev.Key]; ok {
		return mod + name
	}

	return mod + "Unknown"
}

func (r *keyReader) readByte() (byte, error) {
	var buf [1]byte
	file := os.Stdin
	if r.fd != DefaultFd {
		file = os.NewFile(uintptr(r.fd), "terminal")
		defer file.Close()
	}
	_, err := file.Read(buf[:])
	return buf[0], err
}

func (r *keyReader) parse(b byte) (Event, error) {
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
			return Event{Key: KeyRune, Rune: rune(b + 96), Mod: ModCtrl}, nil
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

func (r *keyReader) parseEsc() (Event, error) {
	seq, err := r.readByte()
	if err != nil {
		return Event{Key: KeyEscape}, nil // lone ESC
	}

	if seq != '[' && seq != 'O' {
		// Alt+<key> (ESC + key)
		return Event{Key: KeyRune, Rune: rune(seq), Mod: ModAlt}, nil
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
