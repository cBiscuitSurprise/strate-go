package pieces

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cBiscuitSurprise/strate-go/internal/core"
)

func TestCreateBomb(t *testing.T) {
	want := Piece{
		rank:     RANK_Bomb,
		maxMoves: 0,
	}

	got := *CreatePiece(COLOR_red, RANK_Bomb)

	assert.Equal(t, want.color, got.color)
	assert.Equal(t, want.rank, got.rank)
	assert.NotEmpty(t, got.id)
}

func TestAttackBomb(t *testing.T) {
	testCases := []struct {
		rank Rank
		want core.Winner
	}{
		{RANK_Bomb, core.WINNER_Draw},
		{RANK_Marshal, core.WINNER_Attackee},
		{RANK_General, core.WINNER_Attackee},
		{RANK_Colonel, core.WINNER_Attackee},
		{RANK_Major, core.WINNER_Attackee},
		{RANK_Captain, core.WINNER_Attackee},
		{RANK_Lieutenant, core.WINNER_Attackee},
		{RANK_Sergent, core.WINNER_Attackee},
		{RANK_Minor, core.WINNER_Attacker}, // bomb loses
		{RANK_Scout, core.WINNER_Attackee},
		{RANK_Spy, core.WINNER_Attackee},
		{RANK_Flag, core.WINNER_Attacker}, // bomb loses
	}

	bomb := CreatePiece(COLOR_red, RANK_Bomb)

	for _, tc := range testCases {
		attacker := &Piece{
			rank: tc.rank,
		}

		got, err := AttackBomb(bomb, attacker)

		assert.Nil(t, err)
		assert.Equal(t, tc.want, got)
	}
}

func TestAttackByBomb(t *testing.T) {
	testCases := []struct {
		rank Rank
		want core.Winner
	}{
		{RANK_Bomb, core.WINNER_Draw},
		{RANK_Marshal, core.WINNER_Attacker},
		{RANK_General, core.WINNER_Attacker},
		{RANK_Colonel, core.WINNER_Attacker},
		{RANK_Major, core.WINNER_Attacker},
		{RANK_Captain, core.WINNER_Attacker},
		{RANK_Lieutenant, core.WINNER_Attacker},
		{RANK_Sergent, core.WINNER_Attacker},
		{RANK_Minor, core.WINNER_Attackee}, // bomb loses
		{RANK_Scout, core.WINNER_Attacker},
		{RANK_Spy, core.WINNER_Attacker},
		{RANK_Flag, core.WINNER_Attackee}, // bomb loses
	}

	bomb := CreatePiece(COLOR_red, RANK_Bomb)

	for _, tc := range testCases {
		attackee := &Piece{
			rank: tc.rank,
		}

		got, err := AttackBomb(attackee, bomb)

		assert.Nil(t, err)
		assert.Equal(t, tc.want, got)
	}
}
