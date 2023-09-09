package pieces

import (
	"fmt"

	"github.com/cBiscuitSurprise/strate-go/internal/core"
)

type AttackHandler func(by Piece) core.Winner

type Piece struct {
	id       string
	color    Color
	rank     Rank
	maxMoves int
}

func CreatePiece(index int, color Color, rank Rank) *Piece {
	return &Piece{
		id:       fmt.Sprintf("%s:%02d:%02d", color.String(), rank, index),
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

func (p *Piece) Attack(attackee *Piece) core.Winner {
	if p.GetRank() == RANK_Bomb || attackee.GetRank() == RANK_Bomb {
		winner, _ := AttackBomb(attackee, p)
		return winner
	} else {
		return genericAttackPiece(attackee, p)
	}
}

func genericAttackPiece(attackee *Piece, attacker *Piece) core.Winner {
	if attackee.GetRank() == attacker.GetRank() {
		return core.WINNER_Draw
	} else if attacker.GetRank() > attackee.GetRank() {
		return core.WINNER_Attacker
	} else {
		return core.WINNER_Attackee
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
