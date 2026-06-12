package audit

import (
	"context"
	"sync"
)

type Repository interface {
	Append(ctx context.Context, log Log) error
	List(ctx context.Context) ([]Log, error)
}

type MemoryRepository struct {
	mu   sync.RWMutex
	logs []Log
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) Append(_ context.Context, log Log) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logs = append(r.logs, log)
	return nil
}

func (r *MemoryRepository) List(_ context.Context) ([]Log, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	logs := make([]Log, len(r.logs))
	copy(logs, r.logs)
	return logs, nil
}
