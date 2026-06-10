package database

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, errors.New("postgres dsn is required")
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
