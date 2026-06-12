package main

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"

	"wingops/internal/audit"
	"wingops/internal/auth"
	"wingops/internal/cmdb"
	"wingops/internal/config"
	"wingops/internal/database"
	apphttp "wingops/internal/http"
	"wingops/internal/system"
)

func main() {
	cfg := config.Load()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	ctx := context.Background()
	db, err := database.Open(cfg.Postgres.DSN)
	if err != nil {
		logger.Error("connect postgres failed", zap.Error(err))
		os.Exit(1)
	}
	if err := database.Migrate(ctx, db, "migrations"); err != nil {
		logger.Error("migrate postgres failed", zap.Error(err))
		os.Exit(1)
	}
	if err := database.Seed(ctx, db); err != nil {
		logger.Error("seed postgres failed", zap.Error(err))
		os.Exit(1)
	}

	tokenTTL := time.Duration(cfg.JWT.AccessTokenTTLMinutes) * time.Minute
	router := apphttp.NewRouterWithDependencies(apphttp.Dependencies{
		AuthRepository:   auth.NewPostgresRepository(db),
		AuditRepository:  audit.NewPostgresRepository(db),
		SystemRepository: system.NewPostgresRepository(db),
		CMDBRepository:   cmdb.NewPostgresRepository(db),
		JWTSecret:        cfg.JWT.Secret,
		TokenTTL:         tokenTTL,
	})
	logger.Info("starting wingops server", zap.String("addr", cfg.Server.Addr))
	if err := router.Run(cfg.Server.Addr); err != nil {
		logger.Error("server stopped", zap.Error(err))
		os.Exit(1)
	}
}
