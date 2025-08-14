// golang.org/x/term wrapper

package term

import (
	"os"

	"golang.org/x/term"
)

// State contains the state of a terminal.
type State = term.State

// IsTerminal returns whether the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	return term.IsTerminal(fd)
}

// MakeRaw puts the terminal connected to the given file descriptor into raw mode and returns the previous state of the terminal so that it can be restored.
func MakeRaw(fd int) (*State, error) {
	return term.MakeRaw(fd)
}

// Restore restores the terminal connected to the given file descriptor to a previous state.
func Restore(fd int, oldState *State) error {
	return term.Restore(fd, oldState)
}

//	GetSize returns the visible dimensions of the given terminal.
//
// These dimensions don't include any scrollback buffer height.
func GetSize(fd int) (width, height int, err error) {
	return term.GetSize(fd)
}

// our own!

type Terminal struct {
	fd    int
	state *State
}

func DefaultFd() int {
	return int(os.Stdin.Fd())
}

func NewTerminal(fd int) *Terminal {
	return &Terminal{fd: fd}
}

func NewStdinTerminal() *Terminal {
	return NewTerminal(DefaultFd())
}

func (t *Terminal) Fd() int {
	return t.fd
}

func (t *Terminal) IsTerminal() bool {
	return IsTerminal(t.fd)
}

func (t *Terminal) MakeRaw() error {
	state, err := MakeRaw(t.fd)
	if err != nil {
		return err
	}

	t.state = state
	return nil
}

func (t *Terminal) Restore() error {
	if t.state != nil {
		err := Restore(t.fd, t.state)
		t.state = nil // Clear state after restore
		return err
	}
	return nil
}

func (t *Terminal) GetSize() (width, height int, err error) {
	return GetSize(t.fd)
}

func (t *Terminal) State() *State {
	return t.state
}

func (t *Terminal) HasState() bool {
	return t.state != nil
}
