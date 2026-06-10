package main

import (
	"os"

	"go.uber.org/zap"

	"wingops/internal/config"
	apphttp "wingops/internal/http"
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

	router := apphttp.NewRouter()
	logger.Info("starting wingops server", zap.String("addr", cfg.Server.Addr))
	if err := router.Run(cfg.Server.Addr); err != nil {
		logger.Error("server stopped", zap.Error(err))
		os.Exit(1)
	}
}
