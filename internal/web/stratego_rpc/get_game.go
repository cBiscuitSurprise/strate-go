package stratego_rpc

import (
	"context"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/storage"
	"github.com/cBiscuitSurprise/strate-go/internal/web/apiadapter"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *strateGoServer) GetGame(ctx context.Context, request *pb.GetGameRequest) (*pb.GetGameResponse, error) {
	return requireUserIdDecorator[pb.GetGameRequest, pb.GetGameResponse](ctx, getGameHandler, request)
}

func getGameHandler(userId string, request *pb.GetGameRequest) (*pb.GetGameResponse, error) {
	if g, err := storage.GetGameForUser(userId, request.GetGameId()); err == nil {
		return &pb.GetGameResponse{
			Game: apiadapter.GameToApiGame(g),
		}, nil
	} else {
		log.Error().
			Err(err).
			Msg("failed to get game")
		return nil, status.Errorf(codes.NotFound, "failed to find game, %s", request.GetGameId())
	}
}
