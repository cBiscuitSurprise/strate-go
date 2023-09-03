package web

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aidarkhanov/nanoid"
	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type strateGoServer struct {
	pb.UnimplementedStrateGoServer
}

func (s *strateGoServer) Ping(ctx context.Context, _ *emptypb.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		Message:   "pong",
	}, nil
}

func (s *strateGoServer) DeepPing(ctx context.Context, _ *emptypb.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		Message:   "poOoOong",
		Games:     []*pb.Game{},
	}, nil
}

func (s *strateGoServer) LongPing(stream pb.StrateGo_LongPingServer) error {
	preprend := "[p] "
	preprendLittle := true
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			if err := stream.Send(&pb.Pong{
				Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
				Message:   "[ng] byeeeeee 👋 ...",
				Games:     []*pb.Game{},
			}); err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.Pong{
			Timestamp: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
			Message:   preprend + in.Message,
			Games:     []*pb.Game{},
		}); err != nil {
			return err
		}
		if preprendLittle {
			preprend = "[o] "
		} else {
			preprend = "[O] "
		}
		preprendLittle = !preprendLittle
	}
}

func (s *strateGoServer) NewGame(ctx context.Context, request *pb.NewGameRequest) (*pb.NewGameResponse, error) {
	board, err := game.CreateRandomlyPlannedBoard()
	if err != nil {
		return nil, err
	}

	rows := []*pb.Row{}
	for r := uint8(0); r < board.GetSize().Rows; r++ {
		row := &pb.Row{Columns: []*pb.Square{}}
		for c := uint8(0); c < board.GetSize().Columns; c++ {
			square := board.GetSquare(game.Position{R: r, C: c})

			pbSquare := &pb.Square{Playable: square.IsPlayable()}
			if square.GetPiece() != nil {
				color := pb.PlayerColor_PlayerColor_BLUE
				if square.GetPiece().GetColor() == pieces.COLOR_red {
					color = pb.PlayerColor_PlayerColor_RED
				}
				pbSquare.Piece = &pb.Piece{
					Id:   fmt.Sprintf("%02d:%02d", r, c),
					Rank: uint32(square.GetPiece().GetRank()),
					Player: &pb.GamePlayer{
						Id:    square.GetPiece().GetColor().String(),
						Color: color,
					},
				}
			}
			row.Columns = append(row.Columns, pbSquare)
		}
		rows = append(rows, row)
	}

	pbBoard := &pb.Board{
		Id:         nanoid.New(),
		NumRows:    uint32(board.GetSize().Rows),
		NumColumns: uint32(board.GetSize().Columns),
		Rows:       rows,
	}

	return &pb.NewGameResponse{
		Game: &pb.Game{
			Id:        nanoid.New(),
			State:     pb.GameState_GameState_PLAY,
			PlayerIds: []string{"0", "1"},
			Board:     pbBoard,
		},
	}, nil
}

func newServer(opts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(opts...)
	strategoServer := &strateGoServer{}

	pb.RegisterStrateGoServer(grpcServer, strategoServer)

	return grpcServer
}
