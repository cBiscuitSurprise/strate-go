package pieces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePiece(t *testing.T) {
	testCases := []struct {
		rank  Rank
		moves int
	}{
		{RANK_Bomb, 0},
		{RANK_Marshal, 1},
		{RANK_General, 1},
		{RANK_Colonel, 1},
		{RANK_Major, 1},
		{RANK_Captain, 1},
		{RANK_Lieutenant, 1},
		{RANK_Sergent, 1},
		{RANK_Minor, 1},
		{RANK_Scout, 9},
		{RANK_Spy, 1},
		{RANK_Flag, 0},
	}

	for _, tc := range testCases {
		piece := CreatePiece(nil, tc.rank)

		assert.Equal(t, tc.rank, piece.Rank)
		assert.Equal(t, tc.moves, piece.MaxMoves)
	}
}
