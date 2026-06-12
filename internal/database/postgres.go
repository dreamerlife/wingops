package database

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, errors.New("postgres dsn is required")
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Migrate(ctx context.Context, db *gorm.DB, dir string) error {
	if db == nil {
		return errors.New("postgres db is required")
	}
	if dir == "" {
		dir = "migrations"
	}
	if err := db.WithContext(ctx).Exec(`
CREATE TABLE IF NOT EXISTS schema_migrations (
  version TEXT PRIMARY KEY,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
)`).Error; err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		var applied int64
		if err := db.WithContext(ctx).Raw("SELECT count(*) FROM schema_migrations WHERE version = ?", name).Scan(&applied).Error; err != nil {
			return err
		}
		if applied > 0 {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if execErr := tx.Exec(string(content)).Error; execErr != nil {
				return execErr
			}
			return tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", name).Error
		})
		if err != nil {
			return err
		}
	}
	return nil
}
