package stratego_rpc

import (
	"context"
	"fmt"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
	"github.com/cBiscuitSurprise/strate-go/internal/storage"
	"github.com/cBiscuitSurprise/strate-go/internal/web/apiadapter"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type playGameWebHandlerInput struct {
	request *pb.PlayGameWebListenerRequest
	stream  pb.StrateGo_PlayGameWebListenerServer
}

func (s *strateGoServer) PlayGameWeb(ctx context.Context, request *pb.PlayGameRequest) (*pb.PlayGameWebResponse, error) {
	return requireUserIdDecorator[pb.PlayGameRequest, pb.PlayGameWebResponse](ctx, playGameWebRequestHandler, request)
}

func (s *strateGoServer) PlayGameWebListener(request *pb.PlayGameWebListenerRequest, stream pb.StrateGo_PlayGameWebListenerServer) error {
	_, err := requireUserIdDecorator[playGameWebHandlerInput, dummy](stream.Context(), playGameWebListener, &playGameWebHandlerInput{
		request: request,
		stream:  stream,
	})
	return err
}

func playGameWebRequestHandler(userId string, request *pb.PlayGameRequest) (*pb.PlayGameWebResponse, error) {
	state := &playGameState{
		userId: userId,
	}

	if state.game == nil {
		if err := state.ResolveGame(request.GetGameId()); err != nil {
			log.Warn().
				Err(err).
				Str("gameId", request.GetGameId()).
				Msgf("failed to resolve game with id, '%s'", request.GetGameId())
			return nil, status.Errorf(codes.NotFound, "failed to resolve game with id, '%s'", request.GetGameId())
		}
	}

	if response, err := handlePlayGameRequest(request, state); err == nil {
		return &pb.PlayGameWebResponse{
			GameId:          response.GetGameId(),
			Error:           response.GetError(),
			ValidPlacements: response.GetValidPlacements(),
		}, nil
	} else {
		return nil, err
	}
}

func playGameWebListener(userId string, input *playGameWebHandlerInput) (*dummy, error) {
	state := &playGameState{
		userId: userId,
	}

	stream := input.stream

	if err := state.ResolveGame(input.request.GetGameId()); err != nil {
		log.Warn().
			Err(err).
			Str("gameId", input.request.GetGameId()).
			Msgf("failed to resolve game with id, '%s'", input.request.GetGameId())
		return nil, status.Errorf(codes.NotFound, "failed to resolve game with id, '%s'", input.request.GetGameId())
	}

	client := storage.NewStrategoRedisClient(fmt.Sprintf("Listener:Game:%s", state.game.GetId()))
	client.Connect()
	moves := make(chan game.Move)
	quit := make(chan bool)

	go client.ListenForPieceMoveEvent(state.game.GetId(), "$", moves, quit)

	for move := range moves {
		var attackEvent *pb.AttackEvent
		if move.Result != game.MOVERESULT_NoContest {
			attackEvent = &pb.AttackEvent{
				AttackerRank: int32(move.AttackerRank),
				AttackeeRank: int32(move.AttackeeRank),
				Result:       apiadapter.GameMoveResultToApiMoveResult(move.Result),
			}
		}

		moveEvent := &pb.PieceMovedEvent{
			Nonce:         0,
			PieceId:       move.Id,
			From:          apiadapter.GamePositionToApiPosition(move.From),
			To:            apiadapter.GamePositionToApiPosition(move.To),
			PieceAttacked: attackEvent,
		}

		stream.Send(&pb.PlayGameResponse{
			RedPlayerActive: false,
			PieceMoved:      moveEvent,
		})
	}
	return nil, nil
}
