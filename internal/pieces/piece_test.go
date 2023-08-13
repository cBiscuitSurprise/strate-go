package pieces

import (
	"testing"

	"github.com/cBiscuitSurprise/strate-go/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestPiece(t *testing.T) {
	got := Piece{
		Rank:     RANK_Bomb,
		MaxMoves: 42,
	}

	assert.Equal(t, RANK_Bomb, got.Rank)
	assert.Equal(t, "Bomb", got.getName())
	assert.Equal(t, 42, got.MaxMoves)
}

func TestAttackPiece(t *testing.T) {
	testCases := []struct {
		attacker Rank
		attackee Rank
		want     core.Winner
	}{
		{RANK_Bomb, RANK_Bomb, core.WINNER_Draw},
		{RANK_Bomb, RANK_Marshal, core.WINNER_Attacker},
		{RANK_Marshal, RANK_Bomb, core.WINNER_Attackee},
		{RANK_Marshal, RANK_Marshal, core.WINNER_Draw},
		{RANK_Marshal, RANK_General, core.WINNER_Attacker},
		{RANK_General, RANK_Marshal, core.WINNER_Attackee},
		{RANK_General, RANK_Colonel, core.WINNER_Attacker},
	}

	for _, tc := range testCases {
		attacker := CreatePiece(nil, tc.attacker)
		attackee := CreatePiece(nil, tc.attackee)

		winner, err := attacker.Attack(attackee)

		assert.Nil(t, err)
		assert.Equal(t, tc.want, winner, "Attacker %v, Attackee %v", attacker, attackee)
	}
}
