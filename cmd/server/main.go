package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/linspacestrom/InterShipAv/internal/config"
	"github.com/linspacestrom/InterShipAv/internal/db"
)

func main() {
	log.Println("ChangeLogger started")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ловим SIGINT и SIGTERM
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)

	// загружаем конфиг
	cfg, errConfig := config.LoadConfig()
	if errConfig != nil {
		log.Fatalf("Failed to load config: %s", errConfig)
	}

	log.Println("Config successfully loaded")

	// подключение к базе
	pgxPool, errPgxPool := db.NewPoolPostgres(ctx, cfg.DbConfig)
	if errPgxPool != nil {
		log.Fatalf("DB connection error: %s", errPgxPool)
	}
	defer pgxPool.Close()

	log.Println("Connection to the database is established")

	// ожидаем сигнал или отмену контекста
	select {
	case <-ctx.Done():
		log.Println("Context canceled, shutting down")
	case sig := <-signChan:
		log.Printf("Received signal: %s, shutting down\n", sig)
	}

	log.Println("PR stopped")
}
