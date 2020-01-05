package app

import (
	"sync"
	"sync/atomic"
)

// Storage Interface Supports Get and Set
type HashStorage interface {
	Get(int64) string
	Set(int64, string)
	IncrementIdentifier() int64
	Identifier() int64
	Map() map[int64]string
}

// MemStorage implements Storage
type InMemoryHashStorage struct {
	hashes map[int64]string
	hashIdentifier int64
	mu    *sync.RWMutex
}

//NewStorage creates a new in memory storage
func NewStorage() *InMemoryHashStorage {
	return &InMemoryHashStorage{
		make(map[int64]string),
		0,
		&sync.RWMutex{},
	}
}

func (s *InMemoryHashStorage) Identifier() int64 {
	return s.hashIdentifier
}

func (s *InMemoryHashStorage) IncrementIdentifier() int64 {
	return atomic.AddInt64(&s.hashIdentifier, 1)
}

//Get a cached content by key
func (s *InMemoryHashStorage) Get(key int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value := s.hashes[key]

	return value
}

//Set a cached content by key
func (s *InMemoryHashStorage) Set(key int64, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.hashes[key] = value
}

//Set a cached content by key
func (s *InMemoryHashStorage) Map() map[int64]string {
	return s.hashes
}