package main

import (
	"context"
	_ "log"
	"os"
	"os/signal"
	"syscall"

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
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

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

	go func() {
		if err := srv.Run(ctx); err != nil {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	logger.Info("Server is running", zap.String("addr", cfg.ServerPort))

	select {
	case <-ctx.Done():
	case sig := <-signChan:
		logger.Info("Received signal, shutting down", zap.String("signal", sig.String()))
	}
}
