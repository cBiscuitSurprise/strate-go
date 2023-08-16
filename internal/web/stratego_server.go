package web

import (
	"context"
	"io"
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
		Message:   "poOoOong",
		Games:     []*pb.Game{},
	}, nil
}

func (s *strateGoServer) LongPing(stream pb.StrateGo_LongPingServer) error {
	preprend := "p - "
	preprendLittle := true
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			if err := stream.Send(&pb.Pong{
				Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
				Message:   "ng - bye",
				Games:     []*pb.Game{},
			}); err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.Pong{
			Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
			Message:   preprend + in.Message,
			Games:     []*pb.Game{},
		}); err != nil {
			return err
		}
		if preprendLittle {
			preprend = "o - "
		} else {
			preprend = "O - "
		}
		preprendLittle = !preprendLittle
	}
}

func newServer(opts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(opts...)
	strategoServer := &strateGoServer{}

	pb.RegisterStrateGoServer(grpcServer, strategoServer)

	return grpcServer
}
