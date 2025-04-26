package main

import (
	"log"
	"net"

	"github.com/Omotolani98/insightly/internal/db"
	"github.com/Omotolani98/insightly/internal/query"
	"github.com/Omotolani98/insightly/internal/storage"
	querypb "github.com/Omotolani98/insightly/proto/query"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := storage.AutoMigrate(db); err != nil {
		log.Fatalf("failed migrations: %v", err)
	}

	grpcServer := grpc.NewServer()
	queryServer := query.NewQueryServer(db)

	querypb.RegisterLogQueryServer(grpcServer, queryServer)

	log.Println("🚀 Query server listening on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
