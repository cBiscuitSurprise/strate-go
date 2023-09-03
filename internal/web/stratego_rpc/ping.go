package stratego_rpc

import (
	"context"
	"time"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *strateGoServer) Ping(ctx context.Context, _ *emptypb.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		Message:   "pong",
	}, nil
}
