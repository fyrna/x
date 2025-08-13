// golang.org/x/term wrapper

package term

import (
	"golang.org/x/term"
	"os"
)

var (
	// IsTerminal returns whether the given file descriptor is a terminal.
	IsTerminal = term.IsTerminal

	// MakeRaw puts the terminal connected to the given file descriptor into raw mode and returns the previous state of the terminal so that it can be restored.
	MakeRaw = term.MakeRaw

	// Restore restores the terminal connected to the given file descriptor to a previous state.
	Restore = term.Restore

	// 	GetSize returns the visible dimensions of the given terminal.
	//
	// These dimensions don't include any scrollback buffer height.
	GetSize = term.GetSize
)

type State = term.State

var DefaultFd int = int(os.Stdin.Fd())

func IsStdinTerminal() bool {
	return IsTerminal(DefaultFd)
}

func MakeStdinRaw() (*State, error) {
	return MakeRaw(DefaultFd)
}

func RestoreStdin(state *State) error {
	return Restore(DefaultFd, state)
}

func GetStdinSize() (width, height int, err error) {
	return GetSize(DefaultFd)
}

type Terminal struct {
	fd    int
	state *State
}

func NewTerminal(fd int) *Terminal {
	return &Terminal{fd: fd}
}

func NewStdinTerminal() *Terminal {
	return NewTerminal(DefaultFd)
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
