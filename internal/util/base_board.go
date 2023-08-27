package util

import (
	"github.com/cBiscuitSurprise/strate-go/internal/game"
)

// create an empty standard board
func createStandardBaseBoard() *game.Board {
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
	board := &game.Board{}

	unplayable := []game.Position{
		{R: 4, C: 2},
		{R: 5, C: 2},
		{R: 4, C: 3},
		{R: 5, C: 3},
		{R: 4, C: 6},
		{R: 5, C: 6},
		{R: 4, C: 7},
		{R: 5, C: 7},
	}
	board.Initialize(unplayable)

	return board
}
