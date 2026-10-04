package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	pb "github.com/kritish-dhaubanjar/go-grpc-example/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcClient pb.ServiceClient

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v1/users/"), "/")
	id := parts[0]

	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	resp, err := grpcClient.GetUser(context.Background(), &pb.GetUserRequest{Id: id})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	/*
		data, err := protojson.Marshal(resp)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write()
	*/

	json.NewEncoder(w).Encode(resp)
}

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer conn.Close()

	grpcClient = pb.NewServiceClient(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/users/{id}", getUserHandler)

	log.Println("REST gateway listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
