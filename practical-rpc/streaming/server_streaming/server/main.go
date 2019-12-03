package main

import (
	"google.golang.org/grpc"

	"log"
	"net"

	pb "practical_grpc/server_streaming/server/proto"
)

func main() {
	s := grpc.NewServer()
	pb.RegisterDatabaseServer(s, new(DatabaseService))

	log.Print("Starting RPC server on port 8080...")
	list, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to setup tcp listener: %v", err)
	}

	if err := s.Serve(list); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
