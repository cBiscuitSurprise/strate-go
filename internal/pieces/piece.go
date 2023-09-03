package pieces

import (
	"github.com/cBiscuitSurprise/strate-go/internal/core"
	"github.com/cBiscuitSurprise/strate-go/internal/util"
)

type AttackHandler func(by Piece) core.Winner

type Piece struct {
	id       string
	color    Color
	rank     Rank
	maxMoves int
}

func CreatePiece(color Color, rank Rank) *Piece {
	return &Piece{
		id:       util.NewId(),
		color:    color,
		maxMoves: getMovesForRank(rank),
		rank:     rank,
	}
}

func (p *Piece) GetId() string {
	return p.id
}

func (p *Piece) GetName() string {
	return p.rank.String()
}

func (p *Piece) GetRank() Rank {
	return p.rank
}

func (p *Piece) GetMaxMoves() int {
	return p.maxMoves
}

func (p *Piece) GetColor() Color {
	return p.color
}

func (p *Piece) Attack(attackee *Piece) (core.Winner, error) {
	if p.GetRank() == RANK_Bomb || attackee.GetRank() == RANK_Bomb {
		return AttackBomb(attackee, p)
	} else {
		return genericAttackPiece(attackee, p)
	}
}

func genericAttackPiece(attackee *Piece, attacker *Piece) (core.Winner, error) {
	if attackee.GetRank() == attacker.GetRank() {
		return core.WINNER_Draw, nil
	} else if attacker.GetRank() > attackee.GetRank() {
		return core.WINNER_Attacker, nil
	} else {
		return core.WINNER_Attackee, nil
	}
}

func getMovesForRank(rank Rank) int {
	if rank == RANK_Bomb || rank == RANK_Flag {
		return 0
	} else if rank == RANK_Scout {
		return 9
	}
	return 1
}
