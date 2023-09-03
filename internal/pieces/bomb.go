package pieces

import (
	"github.com/cBiscuitSurprise/strate-go/internal/core"
	game_errors "github.com/cBiscuitSurprise/strate-go/internal/errors"
)

func AttackBomb(attackee *Piece, attacker *Piece) (core.Winner, error) {
	if attacker.GetRank() == attackee.GetRank() {
		return core.WINNER_Draw, nil
	} else if attackee.GetRank() == RANK_Bomb {
		return attackBombWith(attacker), nil
	} else if attacker.GetRank() == RANK_Bomb {
		return attackBombWith(attackee).Invert(), nil
	} else {
		return core.WINNER_Draw, game_errors.GameErrorf(
			game_errors.ERROR_Contest_InvalidContest,
			"invalid Bomb contest, want at least one Bomb, got (%v, %v)", attackee, attacker,
		)
	}
}

func attackBombWith(attacker *Piece) core.Winner {
	if attacker.GetRank() == RANK_Minor || attacker.GetRank() == RANK_Flag {
		return core.WINNER_Attacker
	} else {
		return core.WINNER_Attackee
	}
}
