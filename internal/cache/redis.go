package cache

import (
	"os"

	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
