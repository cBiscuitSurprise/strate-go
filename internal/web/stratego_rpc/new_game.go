package stratego_rpc

import (
	"context"
	"fmt"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/core"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
	"github.com/cBiscuitSurprise/strate-go/internal/util"
	"google.golang.org/grpc/metadata"
)

func (s *strateGoServer) NewGame(ctx context.Context, request *pb.NewGameRequest) (*pb.NewGameResponse, error) {
	userId := ""
	if md, ok := metadata.FromIncomingContext(ctx); !ok {
		return nil, fmt.Errorf("no user-id provided, can't play game")
	} else {
		userId = md["x-stratego-user-id"][0]
	}

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

	rows := []*pb.Row{}
	for r := uint8(0); r < g.Board.GetSize().Rows; r++ {
		row := &pb.Row{Columns: []*pb.Square{}}
		for c := uint8(0); c < g.Board.GetSize().Columns; c++ {
			square := g.Board.GetSquare(game.Position{R: r, C: c})

			pbSquare := &pb.Square{Playable: square.IsPlayable()}
			gamePiece := square.GetPiece()
			if gamePiece != nil {
				color := pb.PlayerColor_PlayerColor_BLUE
				if gamePiece.GetColor() == pieces.COLOR_red {
					color = pb.PlayerColor_PlayerColor_RED
				} else {
					color = pb.PlayerColor_PlayerColor_BLUE
				}
				pbSquare.Piece = &pb.Piece{
					Id:   fmt.Sprintf("%02d:%02d", r, c),
					Rank: uint32(gamePiece.GetRank()),
					Player: &pb.GamePlayer{
						Id:    "",
						Color: color,
					},
				}
			}
			row.Columns = append(row.Columns, pbSquare)
		}
		rows = append(rows, row)
	}

	pbBoard := &pb.Board{
		Id:         util.NewId(),
		NumRows:    uint32(g.Board.GetSize().Rows),
		NumColumns: uint32(g.Board.GetSize().Columns),
		Rows:       rows,
	}

	// storage.SaveGame(g.GetId(), g)

	return &pb.NewGameResponse{
		Game: &pb.Game{
			Id:        g.GetId(),
			State:     pb.GameState_GameState_PLAY,
			PlayerIds: []string{"0", "1"},
			Board:     pbBoard,
		},
	}, nil
}
