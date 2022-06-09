package main

import (
	"context"
	"sync"

	"github.com/gotd/td/session"
)

type memorySession struct {
	mux  sync.RWMutex
	data []byte
}

// LoadSession loads session from memory.
func (s *memorySession) LoadSession(context.Context) ([]byte, error) {
	if s == nil {
		return nil, session.ErrNotFound
	}

	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(s.data) == 0 {
		return nil, session.ErrNotFound
	}

	cpy := append([]byte(nil), s.data...)

	return cpy, nil
}

// StoreSession stores session to memory.
func (s *memorySession) StoreSession(ctx context.Context, data []byte) error {
	s.mux.Lock()
	s.data = data
	s.mux.Unlock()
	return nil
}
