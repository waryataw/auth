package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/waryataw/auth/pkg/user_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserV1Server
}

// Get ...
func (s *server) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	roles := []desc.Role{
		desc.Role_USER,
		desc.Role_ADMIN,
	}

	randomIndex := gofakeit.Number(0, len(roles)-1)

	return &desc.GetResponse{
		Id:        req.GetId(),
		Name:      gofakeit.Username(),
		Email:     gofakeit.Email(),
		Role:      roles[randomIndex],
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s *server) Create(_ context.Context, _ *desc.CreateRequest) (*desc.CreateResponse, error) {
	return &desc.CreateResponse{Id: gofakeit.Int64()}, nil
}

func (s *server) Update(_ context.Context, _ *desc.UpdateRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(_ context.Context, _ *desc.DeleteRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
