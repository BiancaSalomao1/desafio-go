package security

import (
	"context"
	"sync"
	"time"
)

type MemoryTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]time.Time
}

func NewMemoryTokenStore() *MemoryTokenStore {
	return &MemoryTokenStore{
		tokens: make(map[string]time.Time),
	}
}

func (s *MemoryTokenStore) Save(
	ctx context.Context,
	token string,
	ttl time.Duration,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens[token] = time.Now().Add(ttl)

	return nil
}

func (s *MemoryTokenStore) Exists(
	ctx context.Context,
	token string,
) (bool, error) {

	s.mu.RLock()
	expiresAt, exists := s.tokens[token]
	s.mu.RUnlock()

	if !exists {
		return false, nil
	}

	if time.Now().After(expiresAt) {

		s.mu.Lock()
		delete(s.tokens, token)
		s.mu.Unlock()

		return false, nil
	}

	return true, nil
}

func (s *MemoryTokenStore) Delete(
	ctx context.Context,
	token string,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tokens, token)

	return nil
}

var _ TokenStore = (*MemoryTokenStore)(nil)
