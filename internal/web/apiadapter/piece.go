package apiadapter

import (
	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
)

func GameColorToApiColor(c pieces.Color) pb.PlayerColor {
	if c == pieces.COLOR_red {
		return pb.PlayerColor_PlayerColor_RED
	} else {
		return pb.PlayerColor_PlayerColor_BLUE
	}
}

func GamePieceToApiPiece(p *pieces.Piece) *pb.Piece {
	if p == nil {
		return nil
	}

	return &pb.Piece{
		Id:   p.GetId(),
		Rank: uint32(p.GetRank()),
		Player: &pb.GamePlayer{
			Color: GameColorToApiColor(p.GetColor()),
		},
	}
}
