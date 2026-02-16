package memory

import (
	"context"
	"sync"

	"github.com/francisjdev/urlshortener/internal/model"
	"github.com/francisjdev/urlshortener/internal/repository"
)

type MemoryURLHolder struct {
	data map[string]*model.URL
	mu   sync.Mutex
}

func NewURLMemory() repository.URLRepository { // returns the interface
	return &MemoryURLHolder{
		data: make(map[string]*model.URL),
	}
}

func (m *MemoryURLHolder) Create(ctx context.Context, url *model.URL) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.data[url.Code]
	if ok {
		return repository.ErrCodeAlreadyExists

	} else {
		m.data[url.Code] = url
	}
	return nil
}

func (m *MemoryURLHolder) GetByCode(ctx context.Context, code string) (*model.URL, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.data[code]
	if ok {
		return value, nil
	}
	return nil, repository.ErrNotFound

}

func (m *MemoryURLHolder) IncrementHitCount(ctx context.Context, code string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.data[code]
	if ok {
		value.HitCount++
		return nil
	}
	return repository.ErrNotFound
}
