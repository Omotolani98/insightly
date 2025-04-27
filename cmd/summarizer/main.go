package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Omotolani98/insightly/internal/cache"
	"github.com/Omotolani98/insightly/internal/db"
	"github.com/Omotolani98/insightly/internal/storage"
	"github.com/Omotolani98/insightly/internal/summarizer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("⚡ Shutdown signal received, exiting...")
		cancel()
	}()

	rdb := cache.NewRedisClient()

	db, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := storage.AutoMigrate(db); err != nil {
		log.Fatalf("❌ failed to run migrations: %v", err)
	}

	summarizerService := summarizer.NewSummarizer(rdb, db)

	if err := summarizerService.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("❌ Summarizer error: %v", err)
	}
}
