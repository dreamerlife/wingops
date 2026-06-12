package audit

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

type logRow struct {
	ActorID    string
	Method     string
	Path       string
	StatusCode int
	Resource   string
	CreatedAt  time.Time
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Append(ctx context.Context, log Log) error {
	return r.db.WithContext(ctx).Exec(`
INSERT INTO audit_logs (actor_id, method, path, status_code, resource, created_at)
VALUES (NULLIF(?, '')::uuid, ?, ?, ?, ?, ?)`,
		log.ActorID, log.Method, log.Path, log.StatusCode, log.Resource, log.CreatedAt).Error
}

func (r *PostgresRepository) List(ctx context.Context) ([]Log, error) {
	var rows []logRow
	if err := r.db.WithContext(ctx).Raw(`
SELECT COALESCE(actor_id::text, '') AS actor_id, method, path, status_code, resource, created_at
FROM audit_logs
ORDER BY created_at DESC
LIMIT 200`).Scan(&rows).Error; err != nil {
		return nil, err
	}
	logs := make([]Log, 0, len(rows))
	for _, row := range rows {
		logs = append(logs, Log{
			ActorID:    row.ActorID,
			Method:     row.Method,
			Path:       row.Path,
			StatusCode: row.StatusCode,
			Resource:   row.Resource,
			CreatedAt:  row.CreatedAt,
		})
	}
	return logs, nil
}
