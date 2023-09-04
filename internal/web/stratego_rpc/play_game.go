package stratego_rpc

import (
	"fmt"
	"io"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
	"github.com/cBiscuitSurprise/strate-go/internal/storage"
	"github.com/cBiscuitSurprise/strate-go/internal/web/apiadapter"
	"google.golang.org/grpc/metadata"
)

type playGameState struct {
	game            *game.Game
	userId          string
	redPlayerActive bool
}

func (s *playGameState) GetGame() *game.Game {
	return s.game
}

func (s *playGameState) IsRedPlayerActive() bool {
	return s.redPlayerActive
}

func (s *playGameState) SwitchPlayers() {
	s.redPlayerActive = !s.redPlayerActive
}

func (s *strateGoServer) PlayGame(stream pb.StrateGo_PlayGameServer) error {
	userId := ""
	if md, ok := metadata.FromIncomingContext(stream.Context()); !ok {
		return fmt.Errorf("no user-id provided, can't play game")
	} else {
		userId = md["x-stratego-user-id"][0]
	}

	state := &playGameState{
		userId: userId,
	}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			if response, err := handlePlayGameTermination(state); err != nil {
				return err
			} else {
				if err := stream.Send(response); err != nil {
					return err
				}
			}
		} else if err != nil {
			return err
		}

		if response, err := handlePlayGameRequest(in, state); err != nil {
			return err
		} else {
			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}

func handlePlayGameRequest(request *pb.PlayGameRequest, state *playGameState) (*pb.PlayGameResponse, error) {
	if state.game == nil {
		if g, err := storage.GetGame(request.GetGameId()); err == nil {
			state = &playGameState{
				game:   g,
				userId: state.userId,
			}
		} else {
			return nil, fmt.Errorf("failed to load game: %s", request.GetGameId())
		}
	}

	switch request.GetCommand() {
	case pb.PlayGameRequestCommand_PlayGameRequestCommand_PICK_PIECE:
		return handlePlayGamePickPiece(request, state)
	case pb.PlayGameRequestCommand_PlayGameRequestCommand_PLACE_PIECE:
		return handlePlayGamePlacePiece(request, state)
	}

	return &pb.PlayGameResponse{
		RedPlayerActive: true,
	}, nil
}

func handlePlayGameTermination(state *playGameState) (*pb.PlayGameResponse, error) {
	return &pb.PlayGameResponse{}, nil
}

func handlePlayGamePickPiece(request *pb.PlayGameRequest, state *playGameState) (*pb.PlayGameResponse, error) {
	validPlacements := state.GetGame().GetValidMovesFromPosition(
		state.userId,
		apiadapter.ToGamePosition(request.GetSelectedPiecePosition()),
	)

	var validPlacementsOut []*pb.Position
	for _, p := range validPlacements {
		validPlacementsOut = append(validPlacementsOut, apiadapter.ToApiPosition(p))
	}

	return &pb.PlayGameResponse{
		RedPlayerActive: state.IsRedPlayerActive(),
		ValidPlacements: validPlacementsOut,
	}, nil
}

func handlePlayGamePlacePiece(request *pb.PlayGameRequest, state *playGameState) (*pb.PlayGameResponse, error) {
	return nil, nil
}
