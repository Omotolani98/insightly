package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "github.com/Omotolani98/insightly/proto/query" // Update path as needed
)

type Summary struct {
	Stream      string `json:"stream"`
	Text        string `json:"text"`
	WindowStart string `json:"window_start"`
	WindowEnd   string `json:"window_end"`
}

func main() {
	// Connect to gRPC Query service
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ Failed to connect to Query service: %v", err)
	}
	defer conn.Close()

	client := pb.NewLogQueryClient(conn)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,OPTIONS",
	}))

	app.Get("/summaries", func(c *fiber.Ctx) error {
		streamName := c.Query("stream", "")
		limit := int32(5)
		if c.Query("limit") != "" {
			parsedLimit := c.QueryInt("limit")
			limit = int32(parsedLimit)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		stream, err := client.GetSummaries(ctx, &pb.GetReq{
			Stream: streamName,
			Limit:  limit,
		})
		if err != nil {
			log.Printf("❌ Failed to start summary stream: %v", err)
			return c.Status(500).SendString("Failed to start summary stream")
		}

		var summaries []Summary

		for {
			summaryResp, err := stream.Recv()
			if err != nil {
				if err.Error() == "EOF" || status.Code(err) == codes.OutOfRange {
					break
				}
				log.Printf("❌ Error reading stream: %v", err)
				return c.Status(500).SendString("Error reading summaries")
			}

			summaries = append(summaries, Summary{
				Stream:      summaryResp.Stream,
				Text:        summaryResp.Text,
				WindowStart: summaryResp.WindowStart,
				WindowEnd:   summaryResp.WindowEnd,
			})
		}

		return c.JSON(summaries)
	})

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🌐 Insightly API Server listening on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
