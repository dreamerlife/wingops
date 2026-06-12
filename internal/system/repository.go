package system

import (
	"context"
	"sort"
	"sync"
)

type Repository interface {
	List(ctx context.Context) ([]Config, error)
	Save(ctx context.Context, config Config) (Config, error)
}

type MemoryRepository struct {
	mu      sync.RWMutex
	configs map[string]Config
}

func NewMemoryRepository(configs ...Config) *MemoryRepository {
	repo := &MemoryRepository{configs: make(map[string]Config, len(configs))}
	for _, config := range configs {
		repo.configs[config.Key] = config
	}
	return repo
}

func (r *MemoryRepository) List(_ context.Context) ([]Config, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	configs := make([]Config, 0, len(r.configs))
	for _, config := range r.configs {
		configs = append(configs, config)
	}
	sort.Slice(configs, func(i, j int) bool {
		return configs[i].Key < configs[j].Key
	})
	return configs, nil
}

func (r *MemoryRepository) Save(_ context.Context, config Config) (Config, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.configs[config.Key] = config
	return config, nil
}
