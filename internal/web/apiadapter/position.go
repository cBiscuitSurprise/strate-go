package apiadapter

import (
	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
)

func ApiPositionToGamePosition(value *pb.Position) game.Position {
	return game.Position{
		R: int(value.GetRow()),
		C: int(value.GetColumn()),
	}
}

func GamePositionToApiPosition(value *game.Position) *pb.Position {
	return &pb.Position{
		Row:    uint32(value.R),
		Column: uint32(value.C),
	}
}
