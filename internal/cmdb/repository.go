package cmdb

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrModelGroupNotFound = errors.New("model group not found")
	ErrModelNotFound      = errors.New("model not found")
	ErrAssetNotFound      = errors.New("asset not found")
	ErrAssetGroupNotFound = errors.New("asset group not found")
	ErrAPIKeyNotFound     = errors.New("api key not found")
	ErrModelHasAssets     = errors.New("model has assets")
)

type Repository interface {
	ListModelGroups(ctx context.Context) ([]ModelGroup, error)
	CreateModelGroup(ctx context.Context, group ModelGroup) (ModelGroup, error)
	ListModels(ctx context.Context, groupID string) ([]Model, error)
	CreateModel(ctx context.Context, model Model) (Model, error)
	GetModel(ctx context.Context, id string) (Model, error)
	UpdateModel(ctx context.Context, model Model) (Model, error)
	DeleteModel(ctx context.Context, id string) error
	ListAssetGroups(ctx context.Context) ([]AssetGroup, error)
	CreateAssetGroup(ctx context.Context, group AssetGroup) (AssetGroup, error)
	ListAssets(ctx context.Context, filter AssetListFilter) (AssetListResult, error)
	CreateAsset(ctx context.Context, asset Asset, actorID string) (Asset, error)
	GetAsset(ctx context.Context, id string) (Asset, error)
	UpdateAsset(ctx context.Context, asset Asset, actorID string) (Asset, error)
	DeleteAsset(ctx context.Context, id string) error
	ListAssetChangeLogs(ctx context.Context, assetID string) ([]AssetChangeLog, error)
	UpsertAsset(ctx context.Context, asset Asset, actorID string) (Asset, error)
	ListAPIKeys(ctx context.Context) ([]APIKey, error)
	CreateAPIKey(ctx context.Context, key APIKey) (APIKey, error)
	GetAPIKeyByKeyID(ctx context.Context, keyID string) (APIKey, error)
	RevokeAPIKey(ctx context.Context, id string) error
}

type MemoryRepository struct {
	mu               sync.RWMutex
	nextGroupID      int
	nextModelID      int
	nextAssetID      int
	nextAssetGroupID int
	nextRelationID   int
	nextLogID        int
	groups           map[string]ModelGroup
	models           map[string]Model
	assets           map[string]Asset
	assetGroups      map[string]AssetGroup
	changeLogs       map[string][]AssetChangeLog
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextGroupID:      1,
		nextModelID:      1,
		nextAssetID:      1,
		nextAssetGroupID: 1,
		nextRelationID:   1,
		nextLogID:        1,
		groups:           make(map[string]ModelGroup),
		models:           make(map[string]Model),
		assets:           make(map[string]Asset),
		assetGroups:      make(map[string]AssetGroup),
		changeLogs:       make(map[string][]AssetChangeLog),
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
	r.prepareModelRelations(&model)
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
	r.prepareModelRelations(&model)
	r.models[model.ID] = model
	return model, nil
}

func (r *MemoryRepository) DeleteModel(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.models[id]; !ok {
		return ErrModelNotFound
	}
	for _, asset := range r.assets {
		if asset.ModelID == id {
			return ErrModelHasAssets
		}
	}
	delete(r.models, id)
	return nil
}

func (r *MemoryRepository) ListAssetGroups(_ context.Context) ([]AssetGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	groups := make([]AssetGroup, 0, len(r.assetGroups))
	for _, group := range r.assetGroups {
		groups = append(groups, group)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].ID < groups[j].ID
	})
	return groups, nil
}

func (r *MemoryRepository) CreateAssetGroup(_ context.Context, group AssetGroup) (AssetGroup, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	group.ID = strconv.Itoa(r.nextAssetGroupID)
	r.nextAssetGroupID++
	r.assetGroups[group.ID] = group
	return group, nil
}

func (r *MemoryRepository) ListAssets(_ context.Context, filter AssetListFilter) (AssetListResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	assets := make([]Asset, 0, len(r.assets))
	for _, asset := range r.assets {
		if !assetMatchesFilter(asset, filter) {
			continue
		}
		assets = append(assets, asset)
	}
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].ID < assets[j].ID
	})
	total := len(assets)
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	start := (page - 1) * pageSize
	if start > total {
		assets = []Asset{}
	} else {
		end := start + pageSize
		if end > total {
			end = total
		}
		assets = assets[start:end]
	}
	return AssetListResult{Items: assets, Total: total, Page: page, PageSize: pageSize}, nil
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
	if err := r.validateAssetGroups(asset.GroupIDs); err != nil {
		return Asset{}, err
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
	if err := r.validateAssetGroups(asset.GroupIDs); err != nil {
		return Asset{}, err
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

func (r *MemoryRepository) UpsertAsset(ctx context.Context, asset Asset, actorID string) (Asset, error) {
	r.mu.RLock()
	for _, existing := range r.assets {
		if existing.ModelID == asset.ModelID && existing.UniqueKey == asset.UniqueKey {
			asset.ID = existing.ID
			if asset.Status == "" {
				asset.Status = existing.Status
			}
			r.mu.RUnlock()
			return r.UpdateAsset(ctx, asset, actorID)
		}
	}
	r.mu.RUnlock()

	return r.CreateAsset(ctx, asset, actorID)
}

func (r *MemoryRepository) ListAPIKeys(_ context.Context) ([]APIKey, error) {
	return []APIKey{NewDevelopmentAPIKey()}, nil
}

func (r *MemoryRepository) CreateAPIKey(_ context.Context, key APIKey) (APIKey, error) {
	if key.KeyID == "" {
		key.KeyID = "dev-sync-key"
	}
	if key.Secret == "" {
		key.Secret = "dev-sync-secret"
	}
	if key.Status == "" {
		key.Status = "active"
	}
	return key, nil
}

func (r *MemoryRepository) GetAPIKeyByKeyID(_ context.Context, keyID string) (APIKey, error) {
	key := NewDevelopmentAPIKey()
	if key.KeyID != keyID {
		return APIKey{}, ErrAPIKeyNotFound
	}
	return key, nil
}

func (r *MemoryRepository) RevokeAPIKey(_ context.Context, id string) error {
	if id == "" {
		return ErrAPIKeyNotFound
	}
	return nil
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

func (r *MemoryRepository) prepareModelRelations(model *Model) {
	for index := range model.Relations {
		if model.Relations[index].ID == "" {
			model.Relations[index].ID = strconv.Itoa(r.nextRelationID)
			r.nextRelationID++
		}
		model.Relations[index].SourceModelID = model.ID
	}
}

func (r *MemoryRepository) validateAssetGroups(groupIDs []string) error {
	for _, groupID := range groupIDs {
		if _, ok := r.assetGroups[groupID]; !ok {
			return ErrAssetGroupNotFound
		}
	}
	return nil
}

func assetMatchesFilter(asset Asset, filter AssetListFilter) bool {
	if filter.ModelID != "" && asset.ModelID != filter.ModelID {
		return false
	}
	if filter.Status != "" && asset.Status != filter.Status {
		return false
	}
	if filter.GroupID != "" && !containsString(asset.GroupIDs, filter.GroupID) {
		return false
	}
	if filter.Keyword != "" {
		keyword := strings.ToLower(filter.Keyword)
		if strings.Contains(strings.ToLower(asset.UniqueKey), keyword) {
			return true
		}
		for _, value := range asset.Attributes {
			if strings.Contains(strings.ToLower(toString(value)), keyword) {
				return true
			}
		}
		return false
	}
	return true
}

func normalizePage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}
	return page, pageSize
}

func containsString(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func toString(value any) string {
	return fmt.Sprint(value)
}

func cloneMap(values map[string]any) map[string]any {
	cloned := make(map[string]any, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}
