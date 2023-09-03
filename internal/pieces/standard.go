package pieces

func GenerateStandardPieces(color Color) map[string]*Piece {
	standardSet := map[Rank]int{
		RANK_Bomb:       6,
		RANK_Marshal:    1,
		RANK_General:    1,
		RANK_Colonel:    2,
		RANK_Major:      3,
		RANK_Captain:    4,
		RANK_Lieutenant: 4,
		RANK_Sergent:    4,
		RANK_Minor:      5,
		RANK_Scout:      8,
		RANK_Spy:        1,
		RANK_Flag:       1,
	}
	output := map[string]*Piece{}

	for rank, count := range standardSet {
		output = generateRank(output, color, rank, count)
	}

	return output
}

func generateRank(currentSet map[string]*Piece, color Color, rank Rank, count int) map[string]*Piece {
	output := currentSet

	for inx := 0; inx < count; inx++ {
		p := CreatePiece(color, rank)
		output[p.GetId()] = p
	}

	return output
}
