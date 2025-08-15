package store

import (
	"bytes"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("key not found")

// Backend is the lowest-level contract: []byte ↔ []byte.
type Backend interface {
	Put(key, val []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Scan(prefix []byte, fn func(k, v []byte) error) error
	Close() error
}

type Mem struct {
	mu sync.RWMutex
	m  map[string][]byte
}

// NewMem returns a fresh in-memory store.
func NewMem() *Mem {
	return &Mem{m: make(map[string][]byte)}
}

func (m *Mem) Put(k, v []byte) error {
	m.mu.Lock()
	m.m[string(k)] = v
	m.mu.Unlock()
	return nil
}

func (m *Mem) Get(k []byte) ([]byte, error) {
	m.mu.RLock()
	v, ok := m.m[string(k)]
	m.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func (m *Mem) Delete(k []byte) error {
	m.mu.Lock()
	delete(m.m, string(k))
	m.mu.Unlock()
	return nil
}

func (m *Mem) Scan(prefix []byte, fn func(k, v []byte) error) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.m {
		// TODO: improve with better logic
		if !bytes.HasPrefix([]byte(k), prefix) {
			continue
		}
		if err := fn([]byte(k), v); err != nil {
			return err
		}
	}

	return nil
}

// wip, it does nothing for now
func (m *Mem) Close() error {
	return nil
}
