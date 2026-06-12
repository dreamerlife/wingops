package cmdb

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

type fieldRow struct {
	Name        string
	DisplayName string
	FieldType   string
	Required    bool
	UniqueValue bool
	Options     []byte
	SortOrder   int
}

type assetRow struct {
	ID         string
	ModelID    string
	UniqueKey  string
	Status     string
	Attributes []byte
	GroupIDs   []byte
}

type changeLogRow struct {
	ID          string
	AssetID     string
	ActorID     string
	BeforeValue []byte
	AfterValue  []byte
	CreatedAt   time.Time
}

type relationRow struct {
	ID            string
	SourceModelID string
	TargetModelID string
	RelationType  string
	DisplayName   string
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListModelGroups(ctx context.Context) ([]ModelGroup, error) {
	var groups []ModelGroup
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, name, display_name, description
FROM model_groups
ORDER BY name`).Scan(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *PostgresRepository) CreateModelGroup(ctx context.Context, group ModelGroup) (ModelGroup, error) {
	var created ModelGroup
	if err := r.db.WithContext(ctx).Raw(`
INSERT INTO model_groups (name, display_name, description)
VALUES (?, ?, ?)
RETURNING id::text, name, display_name, description`,
		group.Name, group.DisplayName, group.Description).Scan(&created).Error; err != nil {
		return ModelGroup{}, err
	}
	return created, nil
}

func (r *PostgresRepository) ListModels(ctx context.Context, groupID string) ([]Model, error) {
	if exists, err := r.modelGroupExists(ctx, groupID); err != nil || !exists {
		if err != nil {
			return nil, err
		}
		return nil, ErrModelGroupNotFound
	}
	var rows []Model
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, group_id::text, name, display_name, description
FROM models
WHERE group_id = ?::uuid
ORDER BY name`, groupID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	for index := range rows {
		fields, err := r.fieldsForModel(ctx, rows[index].ID)
		if err != nil {
			return nil, err
		}
		rows[index].Fields = fields
	}
	return rows, nil
}

func (r *PostgresRepository) CreateModel(ctx context.Context, model Model) (Model, error) {
	if exists, err := r.modelGroupExists(ctx, model.GroupID); err != nil || !exists {
		if err != nil {
			return Model{}, err
		}
		return Model{}, ErrModelGroupNotFound
	}
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw(`
INSERT INTO models (group_id, name, display_name, description)
VALUES (?::uuid, ?, ?, ?)
RETURNING id::text, group_id::text, name, display_name, description`,
			model.GroupID, model.Name, model.DisplayName, model.Description).Scan(&model).Error; err != nil {
			return err
		}
		if err := saveFields(tx, model.ID, model.Fields); err != nil {
			return err
		}
		return saveRelations(tx, model.ID, model.Relations)
	})
	if err != nil {
		return Model{}, err
	}
	return r.GetModel(ctx, model.ID)
}

func (r *PostgresRepository) GetModel(ctx context.Context, id string) (Model, error) {
	var model Model
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, group_id::text, name, display_name, description
FROM models
WHERE id = ?::uuid`, id).Scan(&model).Error; err != nil {
		return Model{}, err
	}
	if model.ID == "" {
		return Model{}, ErrModelNotFound
	}
	fields, err := r.fieldsForModel(ctx, model.ID)
	if err != nil {
		return Model{}, err
	}
	model.Fields = fields
	relations, err := r.relationsForModel(ctx, model.ID)
	if err != nil {
		return Model{}, err
	}
	model.Relations = relations
	return model, nil
}

func (r *PostgresRepository) UpdateModel(ctx context.Context, model Model) (Model, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(`
UPDATE models
SET name = ?, display_name = ?, description = ?, updated_at = now()
WHERE id = ?::uuid`,
			model.Name, model.DisplayName, model.Description, model.ID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrModelNotFound
		}
		if err := tx.Exec("DELETE FROM model_fields WHERE model_id = ?::uuid", model.ID).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM model_relations WHERE source_model_id = ?::uuid", model.ID).Error; err != nil {
			return err
		}
		if err := saveFields(tx, model.ID, model.Fields); err != nil {
			return err
		}
		return saveRelations(tx, model.ID, model.Relations)
	})
	if err != nil {
		return Model{}, err
	}
	return r.GetModel(ctx, model.ID)
}

func (r *PostgresRepository) DeleteModel(ctx context.Context, id string) error {
	var count int64
	if err := r.db.WithContext(ctx).Raw("SELECT count(*) FROM assets WHERE model_id = ?::uuid", id).Scan(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrModelHasAssets
	}
	result := r.db.WithContext(ctx).Exec("DELETE FROM models WHERE id = ?::uuid", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrModelNotFound
	}
	return nil
}

func (r *PostgresRepository) ListAssetGroups(ctx context.Context) ([]AssetGroup, error) {
	var groups []AssetGroup
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, name, display_name, dimension
FROM asset_groups
ORDER BY dimension, display_name`).Scan(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *PostgresRepository) CreateAssetGroup(ctx context.Context, group AssetGroup) (AssetGroup, error) {
	var created AssetGroup
	if err := r.db.WithContext(ctx).Raw(`
INSERT INTO asset_groups (name, display_name, dimension)
VALUES (?, ?, ?)
RETURNING id::text, name, display_name, dimension`,
		group.Name, group.DisplayName, group.Dimension).Scan(&created).Error; err != nil {
		return AssetGroup{}, err
	}
	return created, nil
}

func (r *PostgresRepository) ListAssets(ctx context.Context, filter AssetListFilter) (AssetListResult, error) {
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	offset := (page - 1) * pageSize
	args := []any{}
	where := "WHERE true"
	if filter.ModelID != "" {
		args = append(args, filter.ModelID)
		where += " AND a.model_id = ?::uuid"
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		where += " AND a.status = ?"
	}
	if filter.GroupID != "" {
		args = append(args, filter.GroupID)
		where += " AND EXISTS (SELECT 1 FROM asset_group_members agm WHERE agm.asset_id = a.id AND agm.group_id = ?::uuid)"
	}
	if filter.Keyword != "" {
		args = append(args, "%"+strings.ToLower(filter.Keyword)+"%")
		where += " AND (lower(a.unique_key) LIKE ? OR lower(a.attributes::text) LIKE ?)"
		args = append(args, "%"+strings.ToLower(filter.Keyword)+"%")
	}
	var total int
	if err := r.db.WithContext(ctx).Raw("SELECT count(*) FROM assets a "+where, args...).Scan(&total).Error; err != nil {
		return AssetListResult{}, err
	}
	var rows []assetRow
	args = append(args, pageSize, offset)
	if err := r.db.WithContext(ctx).Raw(`
SELECT a.id::text, a.model_id::text, a.unique_key, a.status, a.attributes,
  COALESCE(jsonb_agg(agm.group_id::text) FILTER (WHERE agm.group_id IS NOT NULL), '[]'::jsonb) AS group_ids
FROM assets a
LEFT JOIN asset_group_members agm ON agm.asset_id = a.id
`+where+`
GROUP BY a.id
ORDER BY a.updated_at DESC, a.unique_key
LIMIT ? OFFSET ?`, args...).Scan(&rows).Error; err != nil {
		return AssetListResult{}, err
	}
	assets, err := assetsFromRows(rows)
	if err != nil {
		return AssetListResult{}, err
	}
	return AssetListResult{Items: assets, Total: total, Page: page, PageSize: pageSize}, nil
}

func (r *PostgresRepository) CreateAsset(ctx context.Context, asset Asset, actorID string) (Asset, error) {
	model, err := r.GetModel(ctx, asset.ModelID)
	if err != nil {
		return Asset{}, err
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
	groupIDs := append([]string(nil), asset.GroupIDs...)
	attributes, err := json.Marshal(asset.Attributes)
	if err != nil {
		return Asset{}, err
	}
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw(`
INSERT INTO assets (model_id, unique_key, status, attributes)
VALUES (?::uuid, ?, ?, ?::jsonb)
RETURNING id::text, model_id::text, unique_key, status, attributes`,
			asset.ModelID, asset.UniqueKey, asset.Status, string(attributes)).Scan(&assetRow{
			ID: asset.ID,
		}).Error; err != nil {
			return err
		}
		var saved assetRow
		if err := tx.Raw(`
SELECT id::text, model_id::text, unique_key, status, attributes
FROM assets
WHERE model_id = ?::uuid AND unique_key = ?`, asset.ModelID, asset.UniqueKey).Scan(&saved).Error; err != nil {
			return err
		}
		asset = assetFromRow(saved)
		asset.GroupIDs = groupIDs
		if err := saveAssetGroups(tx, asset.ID, asset.GroupIDs); err != nil {
			return err
		}
		return appendChangeLog(tx, asset.ID, actorID, nil, asset.Attributes)
	})
	if err != nil {
		return Asset{}, err
	}
	return asset, nil
}

func (r *PostgresRepository) GetAsset(ctx context.Context, id string) (Asset, error) {
	var row assetRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, model_id::text, unique_key, status, attributes
FROM assets
WHERE id = ?::uuid`, id).Scan(&row).Error; err != nil {
		return Asset{}, err
	}
	if row.ID == "" {
		return Asset{}, ErrAssetNotFound
	}
	asset := assetFromRow(row)
	groupIDs, err := r.groupIDsForAsset(ctx, asset.ID)
	if err != nil {
		return Asset{}, err
	}
	asset.GroupIDs = groupIDs
	return asset, nil
}

func (r *PostgresRepository) UpdateAsset(ctx context.Context, asset Asset, actorID string) (Asset, error) {
	before, err := r.GetAsset(ctx, asset.ID)
	if err != nil {
		return Asset{}, err
	}
	modelID := asset.ModelID
	if modelID == "" {
		modelID = before.ModelID
	}
	model, err := r.GetModel(ctx, modelID)
	if err != nil {
		return Asset{}, err
	}
	if asset.Status == "" {
		asset.Status = before.Status
	}
	if asset.Attributes == nil {
		asset.Attributes = map[string]any{}
	}
	asset.ModelID = modelID
	if err := asset.Validate(model); err != nil {
		return Asset{}, err
	}
	attributes, err := json.Marshal(asset.Attributes)
	if err != nil {
		return Asset{}, err
	}
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(`
UPDATE assets
SET model_id = ?::uuid, unique_key = ?, status = ?, attributes = ?::jsonb, updated_at = now()
WHERE id = ?::uuid`,
			asset.ModelID, asset.UniqueKey, asset.Status, string(attributes), asset.ID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrAssetNotFound
		}
		if err := saveAssetGroups(tx, asset.ID, asset.GroupIDs); err != nil {
			return err
		}
		return appendChangeLog(tx, asset.ID, actorID, before.Attributes, asset.Attributes)
	})
	if err != nil {
		return Asset{}, err
	}
	return r.GetAsset(ctx, asset.ID)
}

func (r *PostgresRepository) DeleteAsset(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Exec("DELETE FROM assets WHERE id = ?::uuid", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAssetNotFound
	}
	return nil
}

func (r *PostgresRepository) ListAssetChangeLogs(ctx context.Context, assetID string) ([]AssetChangeLog, error) {
	if _, err := r.GetAsset(ctx, assetID); err != nil {
		return nil, err
	}
	var rows []changeLogRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, asset_id::text, COALESCE(actor_id::text, '') AS actor_id, before_value, after_value, created_at
FROM asset_change_logs
WHERE asset_id = ?::uuid
ORDER BY created_at DESC`, assetID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	logs := make([]AssetChangeLog, 0, len(rows))
	for _, row := range rows {
		logs = append(logs, AssetChangeLog{
			ID:          row.ID,
			AssetID:     row.AssetID,
			ActorID:     row.ActorID,
			BeforeValue: mapFromJSON(row.BeforeValue),
			AfterValue:  mapFromJSON(row.AfterValue),
			CreatedAt:   row.CreatedAt,
		})
	}
	return logs, nil
}

func (r *PostgresRepository) UpsertAsset(ctx context.Context, asset Asset, actorID string) (Asset, error) {
	var row assetRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, model_id::text, unique_key, status, attributes
FROM assets
WHERE model_id = ?::uuid AND unique_key = ?`, asset.ModelID, asset.UniqueKey).Scan(&row).Error; err != nil {
		return Asset{}, err
	}
	if row.ID == "" {
		return r.CreateAsset(ctx, asset, actorID)
	}
	asset.ID = row.ID
	if asset.Status == "" {
		asset.Status = row.Status
	}
	return r.UpdateAsset(ctx, asset, actorID)
}

func (r *PostgresRepository) ListAPIKeys(ctx context.Context) ([]APIKey, error) {
	var keys []APIKey
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, name, key_id, status
FROM api_keys
ORDER BY created_at DESC`).Scan(&keys).Error; err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *PostgresRepository) CreateAPIKey(ctx context.Context, key APIKey) (APIKey, error) {
	if key.KeyID == "" {
		key.KeyID = "sync-" + randomHex(6)
	}
	if key.Secret == "" {
		key.Secret = randomHex(16)
	}
	if key.Status == "" {
		key.Status = "active"
	}
	if key.Name == "" {
		key.Name = key.KeyID
	}
	if err := r.db.WithContext(ctx).Raw(`
INSERT INTO api_keys (name, key_id, secret_hash, status)
VALUES (?, ?, ?, ?)
RETURNING id::text, name, key_id, status`,
		key.Name, key.KeyID, key.Secret, key.Status).Scan(&key).Error; err != nil {
		return APIKey{}, err
	}
	return key, nil
}

func (r *PostgresRepository) GetAPIKeyByKeyID(ctx context.Context, keyID string) (APIKey, error) {
	var key APIKey
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, name, key_id, secret_hash AS secret, status
FROM api_keys
WHERE key_id = ?`, keyID).Scan(&key).Error; err != nil {
		return APIKey{}, err
	}
	if key.ID == "" {
		return APIKey{}, ErrAPIKeyNotFound
	}
	return key, nil
}

func (r *PostgresRepository) RevokeAPIKey(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Exec(`
UPDATE api_keys
SET status = 'revoked', revoked_at = now()
WHERE id = ?::uuid`, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAPIKeyNotFound
	}
	return nil
}

func (r *PostgresRepository) modelGroupExists(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Raw("SELECT count(*) FROM model_groups WHERE id = ?::uuid", id).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgresRepository) fieldsForModel(ctx context.Context, modelID string) ([]FieldDefinition, error) {
	var rows []fieldRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT name, display_name, field_type, required, unique_value, options, sort_order
FROM model_fields
WHERE model_id = ?::uuid
ORDER BY sort_order, name`, modelID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	fields := make([]FieldDefinition, 0, len(rows))
	for _, row := range rows {
		var options []string
		_ = json.Unmarshal(row.Options, &options)
		fields = append(fields, FieldDefinition{
			Name:        row.Name,
			DisplayName: row.DisplayName,
			Type:        FieldType(row.FieldType),
			Required:    row.Required,
			UniqueValue: row.UniqueValue,
			Options:     options,
			SortOrder:   row.SortOrder,
		})
	}
	return fields, nil
}

func saveFields(tx *gorm.DB, modelID string, fields []FieldDefinition) error {
	for index, field := range fields {
		options, err := json.Marshal(field.Options)
		if err != nil {
			return err
		}
		sortOrder := field.SortOrder
		if sortOrder == 0 {
			sortOrder = index + 1
		}
		if err := tx.Exec(`
INSERT INTO model_fields (model_id, name, display_name, field_type, required, unique_value, options, sort_order)
VALUES (?::uuid, ?, ?, ?, ?, ?, ?::jsonb, ?)`,
			modelID, field.Name, field.DisplayName, string(field.Type), field.Required, field.UniqueValue, string(options), sortOrder).Error; err != nil {
			return err
		}
	}
	return nil
}

func saveRelations(tx *gorm.DB, modelID string, relations []ModelRelation) error {
	for _, relation := range relations {
		sourceID := relation.SourceModelID
		if sourceID == "" {
			sourceID = modelID
		}
		if relation.TargetModelID == "" || relation.RelationType == "" {
			continue
		}
		if err := tx.Exec(`
INSERT INTO model_relations (source_model_id, target_model_id, relation_type, display_name)
VALUES (?::uuid, ?::uuid, ?, ?)`,
			sourceID, relation.TargetModelID, relation.RelationType, relation.DisplayName).Error; err != nil {
			return err
		}
	}
	return nil
}

func saveAssetGroups(tx *gorm.DB, assetID string, groupIDs []string) error {
	if err := tx.Exec("DELETE FROM asset_group_members WHERE asset_id = ?::uuid", assetID).Error; err != nil {
		return err
	}
	for _, groupID := range groupIDs {
		if err := tx.Exec(`
INSERT INTO asset_group_members (asset_id, group_id)
VALUES (?::uuid, ?::uuid)
ON CONFLICT DO NOTHING`, assetID, groupID).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresRepository) relationsForModel(ctx context.Context, modelID string) ([]ModelRelation, error) {
	var rows []relationRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT id::text, source_model_id::text, target_model_id::text, relation_type, display_name
FROM model_relations
WHERE source_model_id = ?::uuid
ORDER BY display_name, relation_type`, modelID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	relations := make([]ModelRelation, 0, len(rows))
	for _, row := range rows {
		relations = append(relations, ModelRelation{
			ID:            row.ID,
			SourceModelID: row.SourceModelID,
			TargetModelID: row.TargetModelID,
			RelationType:  row.RelationType,
			DisplayName:   row.DisplayName,
		})
	}
	return relations, nil
}

func (r *PostgresRepository) groupIDsForAsset(ctx context.Context, assetID string) ([]string, error) {
	var groupIDs []string
	if err := r.db.WithContext(ctx).Raw(`
SELECT group_id::text
FROM asset_group_members
WHERE asset_id = ?::uuid
ORDER BY group_id`, assetID).Scan(&groupIDs).Error; err != nil {
		return nil, err
	}
	return groupIDs, nil
}

func assetsFromRows(rows []assetRow) ([]Asset, error) {
	assets := make([]Asset, 0, len(rows))
	for _, row := range rows {
		assets = append(assets, assetFromRow(row))
	}
	return assets, nil
}

func assetFromRow(row assetRow) Asset {
	groupIDs := make([]string, 0)
	if len(row.GroupIDs) > 0 {
		_ = json.Unmarshal(row.GroupIDs, &groupIDs)
	}
	return Asset{
		ID:         row.ID,
		ModelID:    row.ModelID,
		UniqueKey:  row.UniqueKey,
		Status:     row.Status,
		GroupIDs:   groupIDs,
		Attributes: mapFromJSON(row.Attributes),
	}
}

func appendChangeLog(tx *gorm.DB, assetID string, actorID string, before map[string]any, after map[string]any) error {
	beforeJSON, err := json.Marshal(cloneMap(before))
	if err != nil {
		return err
	}
	afterJSON, err := json.Marshal(cloneMap(after))
	if err != nil {
		return err
	}
	return tx.Exec(`
INSERT INTO asset_change_logs (asset_id, actor_id, before_value, after_value)
VALUES (?::uuid, NULLIF(?, '')::uuid, ?::jsonb, ?::jsonb)`,
		assetID, uuidOrEmpty(actorID), string(beforeJSON), string(afterJSON)).Error
}

func mapFromJSON(raw []byte) map[string]any {
	values := make(map[string]any)
	_ = json.Unmarshal(raw, &values)
	return values
}

func randomHex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return time.Now().Format("20060102150405")
	}
	return hex.EncodeToString(bytes)
}

func uuidOrEmpty(value string) string {
	if len(value) == 36 && strings.Count(value, "-") == 4 {
		return value
	}
	return ""
}
