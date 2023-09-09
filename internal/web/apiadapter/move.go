package apiadapter

import (
	pb "github.com/cBiscuitSurprise/strate-go/api/go/strategopb"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
	"github.com/rs/zerolog/log"
)

func GameMoveResultToApiMoveResult(result game.MoveResult) pb.AttackResult {
	switch result {
	case game.MOVERESULT_NoContest:
		return pb.AttackResult_AttackResult_NO_CONTEST
	case game.MOVERESULT_AttackeeCaptured:
		return pb.AttackResult_AttackResult_ATTACKEE_CAPTURED
	case game.MOVERESULT_AttackerCaptured:
		return pb.AttackResult_AttackResult_ATTACKER_CAPTURED
	case game.MOVERESULT_BothCaptured:
		return pb.AttackResult_AttackResult_BOTH_CAPTURED
	default:
		log.Warn().
			Msgf("failed to convert game-move-result, '%v', to api-move-result", result)
		return pb.AttackResult_AttackResult_NO_CONTEST
	}
}
