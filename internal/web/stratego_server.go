package web

import (
	"context"
	"time"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type strateGoServer struct {
	pb.UnimplementedStrateGoServer
}

func (s *strateGoServer) Ping(ctx context.Context, _ *emptypb.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		Message:   "pong",
	}, nil
}

func (s *strateGoServer) DeepPing(ctx context.Context, _ *emptypb.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		Message:   "pong",
		Games:     []*pb.Game{},
	}, nil
}

func newServer(opts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(opts...)
	strategoServer := &strateGoServer{}

	pb.RegisterStrateGoServer(grpcServer, strategoServer)

	return grpcServer
}
