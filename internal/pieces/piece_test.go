package pieces

import (
	"testing"

	"github.com/cBiscuitSurprise/strate-go/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestCreatePiece(t *testing.T) {
	testCases := []struct {
		rank  Rank
		moves int
		color Color
	}{
		{RANK_Bomb, 0, COLOR_red},
		{RANK_Marshal, 1, COLOR_blue},
		{RANK_General, 1, COLOR_red},
		{RANK_Colonel, 1, COLOR_blue},
		{RANK_Major, 1, COLOR_red},
		{RANK_Captain, 1, COLOR_red},
		{RANK_Lieutenant, 1, COLOR_blue},
		{RANK_Sergent, 1, COLOR_red},
		{RANK_Minor, 1, COLOR_red},
		{RANK_Scout, 9, COLOR_blue},
		{RANK_Spy, 1, COLOR_red},
		{RANK_Flag, 0, COLOR_red},
	}

	for _, tc := range testCases {
		piece := CreatePiece(0, tc.color, tc.rank)

		assert.Equal(t, tc.rank, piece.GetRank())
		assert.Equal(t, tc.moves, piece.GetMaxMoves())
		assert.Equal(t, tc.color, piece.GetColor())
	}
}

func TestPiece(t *testing.T) {
	got := Piece{
		rank:     RANK_Bomb,
		maxMoves: 42,
	}

	assert.Equal(t, RANK_Bomb, got.GetRank())
	assert.Equal(t, "Bomb", got.GetName())
	assert.Equal(t, 42, got.GetMaxMoves())
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
		attacker := CreatePiece(0, COLOR_red, tc.attacker)
		attackee := CreatePiece(1, COLOR_red, tc.attackee)

		winner := attacker.Attack(attackee)
		assert.Equal(t, tc.want, winner, "Attacker %v, Attackee %v", attacker, attackee)
	}
}
