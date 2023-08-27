package util

import (
	"github.com/cBiscuitSurprise/strate-go/internal/core"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
)

func generateStandardPiecesForPlayer(player *core.Player) []*pieces.Piece {
	standardSet := map[pieces.Rank]int{
		pieces.RANK_Bomb:       6,
		pieces.RANK_Marshal:    1,
		pieces.RANK_General:    1,
		pieces.RANK_Colonel:    2,
		pieces.RANK_Major:      3,
		pieces.RANK_Captain:    4,
		pieces.RANK_Lieutenant: 4,
		pieces.RANK_Sergent:    4,
		pieces.RANK_Minor:      5,
		pieces.RANK_Scout:      8,
		pieces.RANK_Spy:        1,
		pieces.RANK_Flag:       1,
	}
	output := []*pieces.Piece{}

	for rank, count := range standardSet {
		output = appendRank(output, player, rank, count)
	}

	return output
}

func appendRank(currentSet []*pieces.Piece, player *core.Player, rank pieces.Rank, count int) []*pieces.Piece {
	output := currentSet

	for inx := 0; inx < count; inx++ {
		output = append(output, pieces.CreatePiece(player, rank))
	}

	return output
}
