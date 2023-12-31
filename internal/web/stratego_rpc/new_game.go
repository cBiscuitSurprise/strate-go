package stratego_rpc

import (
	"context"
	"fmt"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/core"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
	"github.com/cBiscuitSurprise/strate-go/internal/storage"
	"github.com/cBiscuitSurprise/strate-go/internal/web/apiadapter"
)

func (s *strateGoServer) NewGame(ctx context.Context, request *pb.NewGameRequest) (*pb.NewGameResponse, error) {
	return requireUserIdDecorator[pb.NewGameRequest, pb.NewGameResponse](ctx, newGameHandler, request)
}

func newGameHandler(userId string, request *pb.NewGameRequest) (*pb.NewGameResponse, error) {
	red := core.HydratePlayer(userId)
	blue := core.NewPlayer()

	g, err := game.NewTwoPlayerGame(red, blue)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new game, internal error")
	}
	g.Board, err = game.CreateRandomlyPlannedBoard()
	if err != nil {
		return nil, fmt.Errorf("failed to create a new game, internal error")
	}

	storage.SaveGame(g.GetId(), g)

	return &pb.NewGameResponse{
		Game: apiadapter.GameToApiGame(g),
	}, nil
}
