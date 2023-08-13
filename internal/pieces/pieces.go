package pieces

func CreatePiece(rank Rank) *Piece {
	return &Piece{
		Rank:     rank,
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
