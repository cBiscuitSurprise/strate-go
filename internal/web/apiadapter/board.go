package apiadapter

import (
	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
	"github.com/cBiscuitSurprise/strate-go/internal/util"
)

func GameSquareToApiSquare(s *game.Square) *pb.Square {
	return &pb.Square{
		Playable: s.IsPlayable(),
		Piece:    GamePieceToApiPiece(s.GetPiece()),
	}
}

func GameBoardToApiBoard(b *game.Board) *pb.Board {
	rows := make([]*pb.Row, b.GetSize().Rows)
	for rnx, r := range b.Rows() {
		row := &pb.Row{Columns: make([]*pb.Square, b.GetSize().Columns)}
		for cnx, square := range r {
			row.Columns[cnx] = GameSquareToApiSquare(square)
		}
		rows[rnx] = row
	}

	return &pb.Board{
		Id:         util.NewId(),
		NumRows:    uint32(b.GetSize().Rows),
		NumColumns: uint32(b.GetSize().Columns),
		Rows:       rows,
	}
}
