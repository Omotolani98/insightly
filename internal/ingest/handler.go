package ingest

import (
	"context"
	"log"

	ingest "github.com/Omotolani98/insightly/proto/ingest"
	"github.com/redis/go-redis/v9"
)

// IngestServer implements the gRPC LogIngestServer interface
type IngestServer struct {
	ingest.UnimplementedLogIngestServer
	rdb *redis.Client
}

// Constructor for IngestServer
func NewIngestServer(rdb *redis.Client) *IngestServer {
	return &IngestServer{rdb: rdb}
}

// StreamLogs streams incoming logs and pushes them into Redis Stream
func (s *IngestServer) StreamLogs(stream ingest.LogIngest_StreamLogsServer) error {
	ctx := context.Background()

	for {
		record, err := stream.Recv()
		if err != nil {
			break // End of stream (EOF)
		}

		// Push log record into Redis Stream
		_, err = s.rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: "log_stream",
			Values: map[string]interface{}{
				"stream":    record.Stream,
				"timestamp": record.Timestamp,
				"level":     record.Level,
				"message":   record.Message,
				"metadata":  record.Metadata,
			},
		}).Result()

		if err != nil {
			log.Printf("Failed to push to Redis: %v", err)
			return err
		}
	}

	// After all logs are streamed
	return stream.SendAndClose(&ingest.IngestAck{Success: true})
}
