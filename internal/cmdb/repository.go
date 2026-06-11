package cmdb

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"sync"
)

var (
	ErrModelGroupNotFound = errors.New("model group not found")
	ErrModelNotFound      = errors.New("model not found")
)

type Repository interface {
	ListModelGroups(ctx context.Context) ([]ModelGroup, error)
	CreateModelGroup(ctx context.Context, group ModelGroup) (ModelGroup, error)
	ListModels(ctx context.Context, groupID string) ([]Model, error)
	CreateModel(ctx context.Context, model Model) (Model, error)
	GetModel(ctx context.Context, id string) (Model, error)
	UpdateModel(ctx context.Context, model Model) (Model, error)
	DeleteModel(ctx context.Context, id string) error
}

type MemoryRepository struct {
	mu          sync.RWMutex
	nextGroupID int
	nextModelID int
	groups      map[string]ModelGroup
	models      map[string]Model
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextGroupID: 1,
		nextModelID: 1,
		groups:      make(map[string]ModelGroup),
		models:      make(map[string]Model),
	}
}

func (r *MemoryRepository) ListModelGroups(_ context.Context) ([]ModelGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	groups := make([]ModelGroup, 0, len(r.groups))
	for _, group := range r.groups {
		groups = append(groups, group)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].ID < groups[j].ID
	})
	return groups, nil
}

func (r *MemoryRepository) CreateModelGroup(_ context.Context, group ModelGroup) (ModelGroup, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	group.ID = strconv.Itoa(r.nextGroupID)
	r.nextGroupID++
	r.groups[group.ID] = group
	return group, nil
}

func (r *MemoryRepository) ListModels(_ context.Context, groupID string) ([]Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.groups[groupID]; !ok {
		return nil, ErrModelGroupNotFound
	}
	models := make([]Model, 0)
	for _, model := range r.models {
		if model.GroupID == groupID {
			models = append(models, model)
		}
	}
	sort.Slice(models, func(i, j int) bool {
		return models[i].ID < models[j].ID
	})
	return models, nil
}

func (r *MemoryRepository) CreateModel(_ context.Context, model Model) (Model, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.groups[model.GroupID]; !ok {
		return Model{}, ErrModelGroupNotFound
	}
	model.ID = strconv.Itoa(r.nextModelID)
	r.nextModelID++
	r.models[model.ID] = model
	return model, nil
}

func (r *MemoryRepository) GetModel(_ context.Context, id string) (Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	model, ok := r.models[id]
	if !ok {
		return Model{}, ErrModelNotFound
	}
	return model, nil
}

func (r *MemoryRepository) UpdateModel(_ context.Context, model Model) (Model, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.models[model.ID]; !ok {
		return Model{}, ErrModelNotFound
	}
	r.models[model.ID] = model
	return model, nil
}

func (r *MemoryRepository) DeleteModel(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.models[id]; !ok {
		return ErrModelNotFound
	}
	delete(r.models, id)
	return nil
}
