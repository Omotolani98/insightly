package main

import (
	"log"
	"net"

	"github.com/Omotolani98/insightly/internal/cache"
	"github.com/Omotolani98/insightly/internal/ingest"
	ingestpb "github.com/Omotolani98/insightly/proto/ingest"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rdb := cache.NewRedisClient()

	s := grpc.NewServer()
	ingestSrv := ingest.NewIngestServer(rdb)
	ingestpb.RegisterLogIngestServer(s, ingestSrv)

	log.Println("Ingest server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
