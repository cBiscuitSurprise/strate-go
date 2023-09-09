package game

import (
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
	"github.com/cBiscuitSurprise/strate-go/internal/util"
)

// create an empty standard board
func CreateStandardBaseBoard(pieceSet map[string]*pieces.Piece) *Board {
	/* Base Board
	Blue
	0 |  |  |  |  |  |  |  |  |  |  |
	2 |  |  |  |  |  |  |  |  |  |  |
	2 |  |  |  |  |  |  |  |  |  |  |
	3 |  |  |  |  |  |  |  |  |  |  |
	4 |  |  |xx|xx|  |  |xx|xx|  |  |
	5 |  |  |xx|xx|  |  |xx|xx|  |  |
	6 |  |  |  |  |  |  |  |  |  |  |
	7 |  |  |  |  |  |  |  |  |  |  |
	8 |  |  |  |  |  |  |  |  |  |  |
	9 |  |  |  |  |  |  |  |  |  |  |
	   0  1  2  3  4  5  6  7  8  9
	Red
	*/
	board := &Board{}

	unplayable := []Position{
		{R: 4, C: 2},
		{R: 5, C: 2},
		{R: 4, C: 3},
		{R: 5, C: 3},
		{R: 4, C: 6},
		{R: 5, C: 6},
		{R: 4, C: 7},
		{R: 5, C: 7},
	}
	board.Initialize(pieceSet, unplayable)

	return board
}

func CreateStandardPieceSet() map[string]*pieces.Piece {
	playerOnePieces := pieces.GenerateStandardPieces(pieces.COLOR_red)
	playerTwoPieces := pieces.GenerateStandardPieces(pieces.COLOR_blue)

	return util.UpdateMap[string, *pieces.Piece](playerOnePieces, playerTwoPieces)
}
