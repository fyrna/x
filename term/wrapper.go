// golang.org/x/term wrapper

package term

import "golang.org/x/term"

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
