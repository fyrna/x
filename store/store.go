package store

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrNotFound = errors.New("key not found")
	ErrClosed   = errors.New("state is closed")
)

// Backend is the lowest-level contract: []byte ↔ []byte.
type Backend interface {
	Put(key, val []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Scan(prefix []byte, fn func(k, v []byte) error) error
	Range(fn func(k, v []byte) error) error
	Keys(prefix []byte) [][]byte
	Len() int
	Exists(key []byte) bool
	Close() error
}

// Mem is an in-memory Backend implementation.
type Mem struct {
	mu     sync.RWMutex
	m      map[string][]byte
	closed bool
}

// NewMem returns a fresh in-memory store.
func NewMem() *Mem {
	return &Mem{m: make(map[string][]byte)}
}

func clone(b []byte) []byte {
	if b == nil {
		return nil
	}
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (m *Mem) Put(k, v []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrClosed
	}

	m.m[string(k)] = clone(v)
	return nil
}

func (m *Mem) Get(k []byte) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, ErrClosed
	}

	v, ok := m.m[string(k)]
	if !ok {
		return nil, ErrNotFound
	}

	return clone(v), nil
}

func (m *Mem) Delete(k []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrClosed
	}

	sk := string(k)
	if _, ok := m.m[sk]; !ok {
		return ErrNotFound
	}

	delete(m.m, sk)
	return nil
}

func (m *Mem) Scan(prefix []byte, fn func(k, v []byte) error) error {
	m.mu.RLock()
	if m.closed {
		m.mu.RUnlock()
		return ErrClosed
	}

	type kv struct{ k, v []byte }
	ps := string(prefix)
	out := make([]kv, 0, 64)

	for sk, v := range m.m {
		if !strings.HasPrefix(sk, ps) {
			continue
		}
		out = append(out, kv{[]byte(sk), clone(v)})
	}

	m.mu.RUnlock()

	for _, e := range out {
		if err := fn(e.k, e.v); err != nil {
			return err
		}
	}

	return nil
}

func (m *Mem) Range(fn func(k, v []byte) error) error {
	m.mu.RLock()
	if m.closed {
		m.mu.RUnlock()
		return ErrClosed
	}

	type kv struct{ k, v []byte }
	out := make([]kv, 0, len(m.m))

	for sk, v := range m.m {
		out = append(out, kv{[]byte(sk), clone(v)})
	}

	m.mu.RUnlock()

	for _, e := range out {
		if err := fn(e.k, e.v); err != nil {
			return err
		}
	}

	return nil
}

func (m *Mem) Keys(prefix []byte) [][]byte {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil
	}

	ps := string(prefix)
	res := make([][]byte, 0, 64)

	for sk := range m.m {
		if strings.HasPrefix(sk, ps) {
			res = append(res, []byte(sk))
		}
	}

	return res
}

func (m *Mem) Exists(k []byte) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.m[string(k)]
	return ok
}

func (m *Mem) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.m)
}

func (m *Mem) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}

	m.closed = true
	return nil
}
