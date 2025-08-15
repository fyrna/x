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
	Close() error
}

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

	delete(m.m, string(k))
	return nil
}

func (m *Mem) Scan(prefix []byte, fn func(k, v []byte) error) error {
	m.mu.RLock()
	if m.closed {
		m.mu.RUnlock()
		return ErrClosed
	}

	ps := string(prefix)
	type kv struct{ k, v []byte }
	var out []kv

	for sk, v := range m.m {
		if !strings.HasPrefix(sk, ps) {
			continue
		}
		kb := []byte(sk)
		vb := clone(v)
		out = append(out, kv{kb, vb})
	}

	m.mu.RUnlock()

	for _, p := range out {
		if err := fn(p.k, p.v); err != nil {
			return err
		}
	}

	return nil
}

func (m *Mem) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return nil
	}
	for k := range m.m {
		delete(m.m, k)
	}
	m.m = nil
	m.closed = true
	return nil
}
