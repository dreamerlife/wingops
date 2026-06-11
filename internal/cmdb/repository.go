package cmdb

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"sync"
	"time"
)

var (
	ErrModelGroupNotFound = errors.New("model group not found")
	ErrModelNotFound      = errors.New("model not found")
	ErrAssetNotFound      = errors.New("asset not found")
)

type Repository interface {
	ListModelGroups(ctx context.Context) ([]ModelGroup, error)
	CreateModelGroup(ctx context.Context, group ModelGroup) (ModelGroup, error)
	ListModels(ctx context.Context, groupID string) ([]Model, error)
	CreateModel(ctx context.Context, model Model) (Model, error)
	GetModel(ctx context.Context, id string) (Model, error)
	UpdateModel(ctx context.Context, model Model) (Model, error)
	DeleteModel(ctx context.Context, id string) error
	ListAssets(ctx context.Context) ([]Asset, error)
	CreateAsset(ctx context.Context, asset Asset, actorID string) (Asset, error)
	GetAsset(ctx context.Context, id string) (Asset, error)
	UpdateAsset(ctx context.Context, asset Asset, actorID string) (Asset, error)
	DeleteAsset(ctx context.Context, id string) error
	ListAssetChangeLogs(ctx context.Context, assetID string) ([]AssetChangeLog, error)
}

type MemoryRepository struct {
	mu          sync.RWMutex
	nextGroupID int
	nextModelID int
	nextAssetID int
	nextLogID   int
	groups      map[string]ModelGroup
	models      map[string]Model
	assets      map[string]Asset
	changeLogs  map[string][]AssetChangeLog
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextGroupID: 1,
		nextModelID: 1,
		nextAssetID: 1,
		nextLogID:   1,
		groups:      make(map[string]ModelGroup),
		models:      make(map[string]Model),
		assets:      make(map[string]Asset),
		changeLogs:  make(map[string][]AssetChangeLog),
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

func (r *MemoryRepository) ListAssets(_ context.Context) ([]Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	assets := make([]Asset, 0, len(r.assets))
	for _, asset := range r.assets {
		assets = append(assets, asset)
	}
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].ID < assets[j].ID
	})
	return assets, nil
}

func (r *MemoryRepository) CreateAsset(_ context.Context, asset Asset, actorID string) (Asset, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	model, ok := r.models[asset.ModelID]
	if !ok {
		return Asset{}, ErrModelNotFound
	}
	if asset.Status == "" {
		asset.Status = AssetStatusRunning
	}
	if asset.Attributes == nil {
		asset.Attributes = map[string]any{}
	}
	if err := asset.Validate(model); err != nil {
		return Asset{}, err
	}

	asset.ID = strconv.Itoa(r.nextAssetID)
	r.nextAssetID++
	r.assets[asset.ID] = asset
	r.appendChangeLog(asset.ID, actorID, nil, asset.Attributes)
	return asset, nil
}

func (r *MemoryRepository) GetAsset(_ context.Context, id string) (Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	asset, ok := r.assets[id]
	if !ok {
		return Asset{}, ErrAssetNotFound
	}
	return asset, nil
}

func (r *MemoryRepository) UpdateAsset(_ context.Context, asset Asset, actorID string) (Asset, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	before, ok := r.assets[asset.ID]
	if !ok {
		return Asset{}, ErrAssetNotFound
	}
	model, ok := r.models[asset.ModelID]
	if !ok {
		return Asset{}, ErrModelNotFound
	}
	if asset.Status == "" {
		asset.Status = before.Status
	}
	if asset.Attributes == nil {
		asset.Attributes = map[string]any{}
	}
	if err := asset.Validate(model); err != nil {
		return Asset{}, err
	}

	r.assets[asset.ID] = asset
	r.appendChangeLog(asset.ID, actorID, before.Attributes, asset.Attributes)
	return asset, nil
}

func (r *MemoryRepository) DeleteAsset(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.assets[id]; !ok {
		return ErrAssetNotFound
	}
	delete(r.assets, id)
	return nil
}

func (r *MemoryRepository) ListAssetChangeLogs(_ context.Context, assetID string) ([]AssetChangeLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.assets[assetID]; !ok {
		return nil, ErrAssetNotFound
	}
	logs := make([]AssetChangeLog, len(r.changeLogs[assetID]))
	copy(logs, r.changeLogs[assetID])
	return logs, nil
}

func (r *MemoryRepository) appendChangeLog(assetID string, actorID string, before map[string]any, after map[string]any) {
	log := AssetChangeLog{
		ID:          strconv.Itoa(r.nextLogID),
		AssetID:     assetID,
		ActorID:     actorID,
		BeforeValue: cloneMap(before),
		AfterValue:  cloneMap(after),
		CreatedAt:   time.Now(),
	}
	r.nextLogID++
	r.changeLogs[assetID] = append(r.changeLogs[assetID], log)
}

func cloneMap(values map[string]any) map[string]any {
	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}
