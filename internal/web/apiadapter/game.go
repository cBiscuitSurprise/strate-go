package apiadapter

import (
	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"

	"github.com/cBiscuitSurprise/strate-go/internal/game"
)

func GameToApiGame(g *game.Game) *pb.Game {
	return &pb.Game{
		Id:        g.GetId(),
		State:     pb.GameState_GameState_PLAY,
		PlayerIds: []string{"0", "1"},
		Board:     GameBoardToApiBoard(g.Board),
	}
}

func GameToApiGameInfo(g *game.Game) *pb.GameInfo {
	return &pb.GameInfo{
		Id:        g.GetId(),
		State:     pb.GameState_GameState_PLAY,
		PlayerIds: []string{"0", "1"},
	}
}
