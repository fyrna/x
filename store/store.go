package store

import (
	"errors"
	"strings"
	"sync"
	"unsafe"
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

// first time usind unsafe XD
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func s2b(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func (m *Mem) Put(k, v []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrClosed
	}

	m.m[b2s(k)] = clone(v)
	return nil
}

func (m *Mem) Get(k []byte) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, ErrClosed
	}

	v, ok := m.m[b2s(k)]
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

	sk := b2s(k)
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

	ps := b2s(prefix)
	type kv struct{ k, v []byte }
	var out []kv

	for sk, v := range m.m {
		if !strings.HasPrefix(sk, ps) {
			continue
		}
		out = append(out, kv{s2b(sk), clone(v)})
	}

	m.mu.RUnlock()

	for _, p := range out {
		if err := fn(p.k, p.v); err != nil {
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
	var out []kv

	for sk, v := range m.m {
		out = append(out, kv{s2b(sk), clone(v)})
	}

	m.mu.RUnlock()

	for _, p := range out {
		if err := fn(p.k, p.v); err != nil {
			return err
		}
	}

	return nil
}

func (m *Mem) Keys(prefix []byte) [][]byte {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var res [][]byte
	ps := b2s(prefix)

	for sk := range m.m {
		if strings.HasPrefix(sk, ps) {
			res = append(res, s2b(sk))
		}
	}

	return res
}

func (m *Mem) Exists(k []byte) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.m[b2s(k)]
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
