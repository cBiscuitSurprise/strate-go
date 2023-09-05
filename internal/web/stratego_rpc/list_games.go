package stratego_rpc

import (
	"context"
	"fmt"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
	"github.com/cBiscuitSurprise/strate-go/internal/storage"
	"github.com/cBiscuitSurprise/strate-go/internal/web/apiadapter"
	"github.com/rs/zerolog/log"
)

func (s *strateGoServer) ListGames(ctx context.Context, request *pb.ListGamesRequest) (*pb.ListGamesResponse, error) {
	return requireUserIdDecorator[pb.ListGamesRequest, pb.ListGamesResponse](ctx, listGamesHandler, request)
}

func listGamesHandler(userId string, request *pb.ListGamesRequest) (*pb.ListGamesResponse, error) {
	if games, err := storage.ListGamesForUser(userId); err == nil {
		return &pb.ListGamesResponse{
			Games: apiadapter.MapConvert[game.Game, pb.GameInfo](games, apiadapter.GameToApiGameInfo),
		}, nil
	} else {
		log.Error().
			Err(err).
			Msg("failed to list games")
		return nil, fmt.Errorf("failed to list games")
	}
}
