package pieces

import (
	"fmt"

	"github.com/cBiscuitSurprise/strate-go/internal/core"
)

func AttackBomb(attackee Piece, attacker Piece) (core.Winner, error) {
	if attacker.Rank == attackee.Rank {
		return core.WINNER_Draw, nil
	} else if attackee.Rank == RANK_Bomb {
		return attackBombWith(attacker), nil
	} else if attacker.Rank == RANK_Bomb {
		return attackBombWith(attackee).Invert(), nil
	} else {
		return core.WINNER_Draw, fmt.Errorf("Invalid arguments for AttackBomb, want at least one Bomb, got (%v, %v)", attackee, attacker)
	}
}

func attackBombWith(attacker Piece) core.Winner {
	if attacker.Rank == RANK_Minor || attacker.Rank == RANK_Flag {
		return core.WINNER_Attacker
	} else {
		return core.WINNER_Attackee
	}
}
