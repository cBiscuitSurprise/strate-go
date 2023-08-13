package pieces

import "github.com/cBiscuitSurprise/strate-go/internal/core"

func CreatePiece(player *core.Player, rank Rank) *Piece {
	return &Piece{
		Rank:     rank,
		Player:   player,
		MaxMoves: getMovesForRank(rank),
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
