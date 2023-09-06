package stratego_rpc

import (
	"fmt"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
	"github.com/cBiscuitSurprise/strate-go/internal/storage"
	"github.com/cBiscuitSurprise/strate-go/internal/web/apiadapter"
	"github.com/rs/zerolog/log"
)

type playGameState struct {
	game            *game.Game
	userId          string
	redPlayerActive bool
}

func (s *playGameState) GetGame() *game.Game {
	return s.game
}

func (s *playGameState) ResolveGame(gameId string) error {
	if g, err := storage.GetGame(gameId); err == nil {
		if g == nil {
			return fmt.Errorf("no game found with id, '%s'", gameId)
		}
		// TODO: this needs to be moved to the game itself
		g.SetMode(game.GAMEMODE_Play)
		s.game = g
		return nil
	} else {
		return err
	}
}

func (s *playGameState) IsRedPlayerActive() bool {
	return s.redPlayerActive
}

func (s *playGameState) SwitchPlayers() {
	s.redPlayerActive = !s.redPlayerActive
}

func (s *strateGoServer) PlayGame(stream pb.StrateGo_PlayGameServer) error {
	_, err := requireUserIdDecorator[pb.StrateGo_PlayGameServer, dummy](stream.Context(), playGameHandler, &stream)
	return err
}

func playGameHandler(userId string, stream *pb.StrateGo_PlayGameServer) (*dummy, error) {
	state := &playGameState{
		userId: userId,
	}

	// TODO: subscribe to opponent's move-piece events

	handler := StreamingRequestHandler[pb.PlayGameRequest, pb.PlayGameResponse, playGameState]{
		stream:    *stream,
		state:     state,
		process:   handlePlayGameRequest,
		terminate: handlePlayGameTermination,
	}

	return nil, handler.Listen()
}

func handlePlayGameRequest(request *pb.PlayGameRequest, state *playGameState) (*pb.PlayGameResponse, error) {
	if state.game == nil {
		if err := state.ResolveGame(request.GetGameId()); err != nil {
			log.Warn().
				Err(err).
				Msgf("failed to resolve game with id, '%s'", request.GetGameId())
			return &pb.PlayGameResponse{
				RedPlayerActive: state.IsRedPlayerActive(),
				Error:           fmt.Sprintf("failed to resolve game with id, '%s'", request.GetGameId()),
			}, nil
		}
	}

	switch request.GetCommand() {
	case pb.PlayGameRequestCommand_PlayGameRequestCommand_PICK_PIECE:
		return handlePlayGamePickPiece(request, state)
	case pb.PlayGameRequestCommand_PlayGameRequestCommand_MOVE_PIECE:
		return handlePlayGameMovePiece(request, state)
	}

	return &pb.PlayGameResponse{
		RedPlayerActive: state.IsRedPlayerActive(),
	}, nil
}

func handlePlayGameTermination(request *pb.PlayGameRequest, state *playGameState) (*pb.PlayGameResponse, error) {
	return &pb.PlayGameResponse{}, nil
}

func handlePlayGamePickPiece(request *pb.PlayGameRequest, state *playGameState) (*pb.PlayGameResponse, error) {
	validPlacements := state.GetGame().GetValidMovesFromPosition(
		state.userId,
		apiadapter.ApiPositionToGamePosition(request.GetSelectedPiecePosition()),
	)

	return &pb.PlayGameResponse{
		RedPlayerActive: state.IsRedPlayerActive(),
		ValidPlacements: apiadapter.MapConvert[game.Position, pb.Position](validPlacements, apiadapter.GamePositionToApiPosition),
	}, nil
}

func handlePlayGameMovePiece(request *pb.PlayGameRequest, state *playGameState) (*pb.PlayGameResponse, error) {
	if response, err := state.game.MovePiece(
		state.userId,
		apiadapter.ApiPositionToGamePosition(request.GetSelectedPiecePosition()),
		apiadapter.ApiPositionToGamePosition(request.GetSelectedPlacement()),
	); err != nil {
		log.Error().
			Err(err).
			Msg("failed to move piece")

		return &pb.PlayGameResponse{
			RedPlayerActive: state.IsRedPlayerActive(),
			Error:           err.Error(),
		}, nil
	} else {
		state.SwitchPlayers()

		var attackEvent *pb.AttackEvent
		if len(response.RemovedPieces) > 0 {
			removed := make([]string, len(response.RemovedPieces))
			for i, r := range response.RemovedPieces {
				removed[i] = r.GetId()
			}

			var attackerRank int32
			if response.Attacker != nil {
				attackerRank = int32(response.Attacker.GetRank())
			}

			var attackeeRank int32
			if response.Attackee != nil {
				attackeeRank = int32(response.Attackee.GetRank())
			}

			attackEvent = &pb.AttackEvent{
				RemovedPieceIds: removed,
				AttackerRank:    attackerRank,
				AttackeeRank:    attackeeRank,
			}

		}

		// TODO: emit piece-moved-event
		moveEvent := &pb.PieceMovedEvent{
			Nonce:         uint32(state.game.GetNonce()),
			PieceId:       response.Attacker.GetId(),
			From:          request.GetSelectedPiecePosition(),
			To:            request.GetSelectedPlacement(),
			PieceAttacked: attackEvent,
		}
		return &pb.PlayGameResponse{
			RedPlayerActive: state.IsRedPlayerActive(),
			PieceMoved:      moveEvent,
		}, nil
	}
}
