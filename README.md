# Go gRPC Example

A minimal example of building a **gRPC service in Go** using Protocol Buffers, together with a lightweight HTTP/REST gateway.

The project demonstrates the flow:

```text
HTTP Client
    │
    │ POST /v1/users/{id}
    ▼
HTTP Gateway
    │
    │ gRPC
    ▼
gRPC Server
    │
    ▼
GetUser()
```

The gRPC service is defined once using a `.proto` file, and Go client/server code is generated from that definition.

## Features

- Go gRPC server
- Protocol Buffers service definition
- Generated Go gRPC client and server interfaces
- Unary gRPC RPC
- HTTP/REST gateway using Go's standard `net/http`
- JSON response from the HTTP endpoint
- Minimal dependencies and project structure

## Project Structure

```text
.
├── client/
│   └── main.go
│
├── proto/
│   ├── service.proto
│   ├── service.pb.go
│   └── service_grpc.pb.go
│
├── server/
│   └── main.go
│
├── go.mod
└── go.sum
```

### `proto/`

Contains the Protocol Buffers service definition and generated Go code.

The service exposes a single RPC:

```protobuf
service Service {
  rpc GetUser(GetUserRequest) returns (User);
}
```

### `server/`

Contains the gRPC server implementation.

The server listens on:

```text
localhost:50051
```

### `client/`

Despite the directory name, this application acts as an **HTTP gateway**.

It:

1. Starts an HTTP server on port `8080`.
2. Creates a gRPC client connection to `localhost:50051`.
3. Receives HTTP requests.
4. Calls the gRPC `GetUser` RPC.
5. Returns the gRPC response as JSON.

## Requirements

- [Go](https://go.dev/) 1.27+
- Protocol Buffers compiler (`protoc`)
- Go Protocol Buffers plugins

The project uses:

- `google.golang.org/grpc`
- `google.golang.org/protobuf`

## Getting Started

Clone the repository:

```bash
git clone https://github.com/kritish-dhaubanjar/go-grpc-example.git
cd go-grpc-example
```

Download dependencies:

```bash
go mod download
```

## Generate gRPC Code

The gRPC client and server code is generated from:

```text
proto/service.proto
```

Install the Go Protocol Buffers plugins if they are not already installed:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Make sure the Go binary directory is in your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then generate the Go code:

```bash
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/service.proto
```

This generates:

```text
proto/
├── service.proto
├── service.pb.go
└── service_grpc.pb.go
```

The generated files contain the Protocol Buffers message types and the gRPC client/server interfaces.

## Running the Example

The example consists of two processes: the gRPC server and the HTTP gateway.

### 1. Start the gRPC Server

In the first terminal:

```bash
go run ./server
```

You should see:

```text
gRPC server listening on :50051
```

### 2. Start the HTTP Gateway

In a second terminal:

```bash
go run ./client
```

You should see:

```text
REST gateway listening on :8080
```

The gateway is now available at:

```text
http://localhost:8080
```

## Calling the API

The HTTP gateway exposes:

```http
POST /v1/users/{id}
```

For example:

```bash
curl -X POST http://localhost:8080/v1/users/123
```

Response:

```json
{
  "id": "123",
  "name": "John Doe",
  "email": "johndoe@example.com"
}
```

The request travels through the following path:

```text
curl
 │
 │ POST /v1/users/123
 ▼
HTTP Gateway :8080
 │
 │ GetUser({ id: "123" })
 ▼
gRPC Server :50051
 │
 │ User
 ▼
HTTP Gateway
 │
 │ JSON
 ▼
curl
```

## gRPC Service

The service is defined in [`proto/service.proto`](proto/service.proto):

```protobuf
syntax = "proto3";

package service;

option go_package = "github.com/kritish-dhaubanjar/go-grpc-example/proto";

service Service {
  rpc GetUser(GetUserRequest) returns (User);
}

message GetUserRequest {
  string id = 1;
}

message User {
  string id = 1;
  string name = 2;
  string email = 3;
}
```

### Request

```protobuf
message GetUserRequest {
  string id = 1;
}
```

### Response

```protobuf
message User {
  string id = 1;
  string name = 2;
  string email = 3;
}
```

## gRPC Server

The server implements the generated `ServiceServer` interface:

```go
type Server struct {
    pb.UnimplementedServiceServer
}

func (s *Server) GetUser(
    ctx context.Context,
    in *pb.GetUserRequest,
) (*pb.User, error) {
    return &pb.User{
        Id:    in.Id,
        Name:  "John Doe",
        Email: "johndoe@example.com",
    }, nil
}
```

It then registers the implementation with the gRPC server:

```go
pb.RegisterServiceServer(server, &Server{})
```

and listens on port `50051`.

## HTTP Gateway

The HTTP gateway creates a gRPC connection:

```go
conn, err := grpc.NewClient(
    "localhost:50051",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
```

It creates a generated gRPC client:

```go
grpcClient = pb.NewServiceClient(conn)
```

and exposes the HTTP endpoint:

```text
POST /v1/users/{id}
```

When the endpoint is called, the gateway converts the HTTP request into a gRPC request:

```go
resp, err := grpcClient.GetUser(
    context.Background(),
    &pb.GetUserRequest{Id: id},
)
```

The resulting protobuf message is then encoded as JSON for the HTTP client.

## Why gRPC?

gRPC allows the service contract to be defined using Protocol Buffers and then generates strongly typed client and server code from that contract.

In this example:

```text
service.proto
     │
     │ protoc
     ▼
Generated Go code
     │
     ├── gRPC Server interface
     ├── gRPC Client
     └── Protocol Buffer types
```

This avoids manually defining request/response serialization and client/server interfaces.

For more information, see the official [gRPC Go documentation](https://grpc.io/docs/languages/go/) and [gRPC Go basics guide](https://grpc.io/docs/languages/go/basics/).

## Development

Run the server:

```bash
go run ./server
```

Run the HTTP gateway:

```bash
go run ./client
```

Run all tests:

```bash
go test ./...
```

Format the code:

```bash
gofmt -w .
```

Update dependencies:

```bash
go mod tidy
```
