package stratego_rpc

import (
	"context"
	"time"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *strateGoServer) DeepPing(ctx context.Context, _ *emptypb.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		Message:   "poOoOong",
		Games:     []*pb.Game{},
	}, nil
}
