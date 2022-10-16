package memory

import (
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
	"webframe/session"
)

type Session struct {
	mu     sync.RWMutex
	id     string
	Values map[string]string
}

func (s *Session) Get(ctx context.Context, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.Values[key]
	if !ok {
		return "", errors.New("not found")
	}
	return val, nil
}

func (s *Session) Set(ctx context.Context, key string, val string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Values[key] = val
	return nil
}

func (s *Session) ID() string {
	return s.id
}

type Mem struct {
	mu         sync.RWMutex
	c          *cache.Cache
	expiration time.Duration
}

func (m *Mem) Get(ctx context.Context, id string) (session.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.c.Get(id)
	if !ok {
		return nil, errors.New("not found")
	}
	return val.(*Session), nil
}

func (m *Mem) Genrate(ctx context.Context, id string) (session.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sess := &Session{
		id:     id,
		Values: make(map[string]string, 16),
	}
	m.c.Set(id, sess, m.expiration)
	return sess, nil
}

func (m *Mem) Remove(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.c.Delete(id)
	return nil
}

func (m *Mem) Reflash(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.c.Get(id)
	if !ok {
		return errors.New("not found")
	}
	m.c.Set(id, s, m.expiration)
	return nil
}
