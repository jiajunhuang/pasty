package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) GetPaste(ctx context.Context, req *PasteRequest) (*PasteResponse, error) {
	return &PasteResponse{Items: []*PasteItem{}}, nil
}

func main() {
	c := NewConfig()

	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listen at: %v", c.Addr)
	s := grpc.NewServer()
	RegisterPastyServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
