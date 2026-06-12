package system

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

type configRow struct {
	Key   string
	Value []byte
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context) ([]Config, error) {
	var rows []configRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT key, value
FROM system_configs
ORDER BY key`).Scan(&rows).Error; err != nil {
		return nil, err
	}
	configs := make([]Config, 0, len(rows))
	for _, row := range rows {
		configs = append(configs, Config{Key: row.Key, Value: stringValue(row.Value)})
	}
	return configs, nil
}

func (r *PostgresRepository) Save(ctx context.Context, config Config) (Config, error) {
	value, err := json.Marshal(config.Value)
	if err != nil {
		return Config{}, err
	}
	if err := r.db.WithContext(ctx).Exec(`
INSERT INTO system_configs (key, value, updated_at)
VALUES (?, ?::jsonb, now())
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = now()`,
		config.Key, string(value)).Error; err != nil {
		return Config{}, err
	}
	return config, nil
}

func stringValue(raw []byte) string {
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
	}
	return string(raw)
}
