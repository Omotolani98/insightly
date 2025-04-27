package summarizer

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Omotolani98/insightly/internal/config"
	"github.com/Omotolani98/insightly/internal/llm"
	"github.com/Omotolani98/insightly/internal/storage"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Summarizer struct {
	rdb       *redis.Client
	db        *gorm.DB
	llmClient *llm.Client
}

func NewSummarizer(rdb *redis.Client, db *gorm.DB) *Summarizer {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed loading config: %v", err)
	}

	llmBaseURL := "http://" + cfg.LLMHost + ":" + cfg.LLMPort
	engineID := cfg.EngineID
	modelName := cfg.ModelName

	llmClient := llm.NewClient(llmBaseURL, engineID, modelName)

	return &Summarizer{
		rdb:       rdb,
		db:        db,
		llmClient: llmClient,
	}
}

func (s *Summarizer) Run(ctx context.Context) error {
	err := s.rdb.XGroupCreateMkStream(ctx, "log_stream", "summarizers", "$").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return fmt.Errorf("could not create consumer group: %w", err)
	}

	log.Println("🚀 Summarizer worker started...")

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.processBatch(ctx); err != nil {
				log.Printf("Failed processing batch: %v", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Summarizer) processBatch(ctx context.Context) error {
	entries, err := s.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    "summarizers",
		Consumer: "worker-1",
		Streams:  []string{"log_stream", ">"},
		Count:    1000,
		Block:    0,
	}).Result()

	if err != nil && err != redis.Nil {
		return err
	}

	if len(entries) == 0 {
		log.Println("No logs to summarize...")
		return nil
	}

	log.Println("📦 Processing logs batch...")

	// Batch by stream
	batches := make(map[string][]string)

	var ackIDs []string

	for _, stream := range entries {
		for _, msg := range stream.Messages {
			streamName := msg.Values["stream"].(string)
			timestamp := fmt.Sprintf("%v", msg.Values["timestamp"])
			level := fmt.Sprintf("%v", msg.Values["level"])
			message := fmt.Sprintf("%v", msg.Values["message"])

			formatted := fmt.Sprintf("[%s][%s] %s", level, timestamp, message)
			batches[streamName] = append(batches[streamName], formatted)

			ackIDs = append(ackIDs, msg.ID)
		}
	}

	// Summarize each stream separately
	for streamName, logs := range batches {
		prompt := strings.Join(logs, "\n")
		summaryText, err := s.llmClient.Summarize(prompt)
		if err != nil {
			log.Printf("❌ LLM summarization failed: %v", err)
			continue
		}

		now := time.Now().UTC()

		summary := storage.Summary{
			Stream:      streamName,
			WindowStart: now.Add(-1 * time.Minute),
			WindowEnd:   now,
			Text:        summaryText,
			CreatedAt:   now,
		}

		if err := s.db.Create(&summary).Error; err != nil {
			log.Printf("❌ Failed saving summary to DB: %v", err)
			continue
		}

		log.Printf("✅ Summarized %d logs from stream '%s'", len(logs), streamName)
	}

	if len(ackIDs) > 0 {
		if err := s.rdb.XAck(ctx, "log_stream", "summarizers", ackIDs...).Err(); err != nil {
			log.Printf("❌ Failed to ack logs: %v", err)
		} else {
			log.Printf("✅ Acknowledged %d logs", len(ackIDs))
		}
	}

	return nil
}
