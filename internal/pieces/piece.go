package pieces

import "github.com/cBiscuitSurprise/strate-go/internal/core"

type AttackHandler func(by Piece) core.Winner

type Piece struct {
	Rank     Rank
	MaxMoves int
}

func (p Piece) getName() string {
	return p.Rank.String()
}

func (p Piece) attack(attackee Piece) (core.Winner, error) {
	if p.Rank == RANK_Bomb || attackee.Rank == RANK_Bomb {
		return AttackBomb(attackee, p)
	} else {
		return genericAttackPiece(attackee, p)
	}
}

func genericAttackPiece(attackee Piece, attacker Piece) (core.Winner, error) {
	if attackee.Rank == attacker.Rank {
		return core.WINNER_Draw, nil
	} else if attacker.Rank > attackee.Rank {
		return core.WINNER_Attacker, nil
	} else {
		return core.WINNER_Attackee, nil
	}
}
