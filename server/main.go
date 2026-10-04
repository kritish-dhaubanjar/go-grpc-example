package main

import (
	"context"
	"log"
	"net"

	pb "github.com/kritish-dhaubanjar/go-grpc-example/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedServiceServer
}

func (s *Server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.User, error) {
	return &pb.User{
		Id:    in.Id,
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	pb.RegisterServiceServer(server, &Server{})

	log.Println("gRPC server listening on :50051")

	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
