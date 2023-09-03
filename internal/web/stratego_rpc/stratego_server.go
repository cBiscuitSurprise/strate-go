package stratego_rpc

import (
	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"google.golang.org/grpc"
)

type strateGoServer struct {
	pb.UnimplementedStrateGoServer
}

func NewServer(opts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(opts...)
	strategoServer := &strateGoServer{}

	pb.RegisterStrateGoServer(grpcServer, strategoServer)

	return grpcServer
}
