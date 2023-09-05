package stratego_rpc

import (
	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
)

type planGameState struct {
	Game game.Game
}

func (s *strateGoServer) PlanGame(stream pb.StrateGo_PlanGameServer) error {
	_, err := requireUserIdDecorator[pb.StrateGo_PlanGameServer, dummy](stream.Context(), planGameHandler, &stream)
	return err
}

func planGameHandler(userId string, stream *pb.StrateGo_PlanGameServer) (*dummy, error) {
	state := &planGameState{}

	handler := StreamingRequestHandler[pb.PlanGameRequest, pb.PlanGameResponse, planGameState]{
		stream:    *stream,
		state:     state,
		process:   handlePlanGameRequest,
		terminate: handlePlanGameTermination,
	}

	return nil, handler.Listen()
}

func handlePlanGameRequest(request *pb.PlanGameRequest, state *planGameState) (*pb.PlanGameResponse, error) {
	return nil, nil
}

func handlePlanGameTermination(request *pb.PlanGameRequest, state *planGameState) (*pb.PlanGameResponse, error) {
	// pop from
	return nil, nil
}
