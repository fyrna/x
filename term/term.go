package term

import "os"

type Terminal struct {
	fd    int
	state *State
}

func NewTerminal(fd int) *Terminal {
	return &Terminal{fd: fd}
}

func DefaultFd() int {
	return int(os.Stdin.Fd())
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
