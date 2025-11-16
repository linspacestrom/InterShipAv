package main

import (
	"context"

	"github.com/linspacestrom/InterShipAv/internal/config"
	"github.com/linspacestrom/InterShipAv/internal/db"
	"github.com/linspacestrom/InterShipAv/internal/handlers"
	"github.com/linspacestrom/InterShipAv/internal/repositories"
	"github.com/linspacestrom/InterShipAv/internal/services"
	"github.com/linspacestrom/InterShipAv/internal/transaction"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// Подключаемся к базе
	pool, err := db.NewPoolPostgres(ctx, cfg.DbConfig)
	if err != nil {
		logger.Fatal("DB connection error", zap.Error(err))
	}
	defer pool.Close()

	teamRepo := repositories.NewTeamRepository(pool)
	userRepo := repositories.NewUserRepository(pool)
	prRepo := repositories.NewPullRequestRepository(pool)
	tm := transaction.NewManager(pool)

	teamSvc := services.NewTeamService(teamRepo, userRepo, tm)
	userSvc := services.NewUserService(userRepo, prRepo, tm)
	prSvc := services.NewPRService(prRepo, userRepo, teamRepo, tm)

	srv := handlers.NewServer(cfg, &teamSvc, &userSvc, &prSvc, logger)

	if err := srv.Run(ctx); err != nil {
		logger.Error("server stopped with error", zap.Error(err))
	} else {
		logger.Info("server stopped gracefully")
	}
}
