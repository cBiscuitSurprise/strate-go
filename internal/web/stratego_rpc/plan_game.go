package stratego_rpc

import (
	"io"

	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
)

type planGameState struct {
	Game game.Game
}

func (s *strateGoServer) PlanGame(stream pb.StrateGo_PlanGameServer) error {
	state := &planGameState{}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			if response, err := handlePlanGameTermination(state); err != nil {
				return err
			} else {
				if err := stream.Send(response); err != nil {
					return err
				}
			}
		} else if err != nil {
			return err
		}

		if response, err := handlePlanGameRequest(in, state); err != nil {
			return err
		} else {
			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}

func handlePlanGameRequest(request *pb.PlanGameRequest, state *planGameState) (*pb.PlanGameResponse, error) {
	return nil, nil
}

func handlePlanGameTermination(state *planGameState) (*pb.PlanGameResponse, error) {
	// pop from
	return nil, nil
}
