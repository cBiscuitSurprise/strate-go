package game

import (
	"testing"

	"github.com/cBiscuitSurprise/strate-go/internal/core"
	game_errors "github.com/cBiscuitSurprise/strate-go/internal/errors"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
	"github.com/stretchr/testify/assert"
)

func TestBoardUnInitialized(t *testing.T) {
	board := Board{}

	assert.Len(t, board.squares, 10)
	assert.Len(t, board.squares[0], 10)
	assert.Nil(t, board.squares[0][0])

	err := board.placePiece(nil, Position{0, 0})
	assert.Equal(t, game_errors.ERROR_Board_Uninitialized, err.Code)
}

func TestBoardInitialized(t *testing.T) {
	board := Board{}

	unplayable := []Position{{5, 6}}
	board.initialize(unplayable)

	for x, row := range board.squares {
		for y, cell := range row {
			assert.NotNil(t, cell)
			assert.Nil(t, cell.piece)

			if x == 5 && y == 6 {
				assert.False(t, cell.playable)
			} else {
				assert.True(t, cell.playable)
			}
		}
	}
}

func TestBoardPlacePiece(t *testing.T) {
	flag := pieces.CreatePiece(nil, pieces.RANK_Flag)
	board := Board{}

	unplayable := []Position{{5, 6}}

	// Error: Index Out of Range
	err := board.placePiece(nil, Position{20, 0})
	assert.Equal(t, game_errors.ERROR_Board_IndexOutOfRange, err.Code)

	err = board.placePiece(nil, Position{0, 20})
	assert.Equal(t, game_errors.ERROR_Board_IndexOutOfRange, err.Code)

	// Error: Uninitialized Board
	err = board.placePiece(nil, Position{0, 0})
	assert.Equal(t, game_errors.ERROR_Board_Uninitialized, err.Code)

	board.initialize(unplayable)

	// Success
	validPosition := Position{2, 9}
	err = board.placePiece(flag, validPosition)

	assert.Nil(t, err)
	assert.Equal(t, flag, board.squares[2][9].piece)

	// Error: Occupied
	err = board.placePiece(flag, validPosition)

	assert.Equal(t, game_errors.ERROR_Board_OccupiedSquare, err.Code)

	// Error: Unplayable Square
	invalidPosition := Position{5, 6}
	err = board.placePiece(flag, invalidPosition)

	assert.NotNil(t, err)
	assert.Equal(t, game_errors.ERROR_Board_UnplayableSquare, err.Code)
}

func TestBoardMovePiece(t *testing.T) {
	board := Board{}

	/* Setup Board
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |1M|2M|  |  |  |  |
	|  |  |  |  |  |2C|xx|  |  |  |
	|  |  |  |  |  |  |1G|  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	*/
	unplayable := []Position{{5, 6}}
	board.initialize(unplayable)

	playerOne := &core.Player{Number: 1}
	playerTwo := &core.Player{Number: 2}

	playerOneGeneral := pieces.CreatePiece(playerOne, pieces.RANK_General)
	err := board.placePiece(playerOneGeneral, Position{6, 6})
	assert.Nil(t, err)

	playerOneMarshal := pieces.CreatePiece(playerOne, pieces.RANK_Marshal)
	err = board.placePiece(playerOneMarshal, Position{4, 4})
	assert.Nil(t, err)

	playerTwoColonel := pieces.CreatePiece(playerTwo, pieces.RANK_Colonel)
	err = board.placePiece(playerTwoColonel, Position{5, 5})
	assert.Nil(t, err)

	playerTwoMarshal := pieces.CreatePiece(playerTwo, pieces.RANK_Marshal)
	err = board.placePiece(playerTwoMarshal, Position{4, 5})
	assert.Nil(t, err)

	// Move player-one-general up one space (Error: unplayable)
	removed, err := board.MovePiece(Position{6, 6}, Position{5, 6})
	assert.Equal(t, game_errors.ERROR_Board_UnplayableSquare, err.Code)
	assert.Len(t, removed, 0)
	assert.Equal(t, playerOneGeneral, board.getSquare(Position{6, 6}).piece)

	// Move player-one-general left one space
	removed, err = board.MovePiece(Position{6, 6}, Position{6, 5})
	assert.Nil(t, err)
	assert.Len(t, removed, 0)
	assert.Equal(t, playerOneGeneral, board.getSquare(Position{6, 5}).piece)
	assert.Nil(t, board.getSquare(Position{6, 6}).piece)

	// Move player-one-general up one space (take Player Two Colonel)
	removed, err = board.MovePiece(Position{6, 5}, Position{5, 5})
	assert.Nil(t, err)
	assert.Len(t, removed, 1)
	assert.Equal(t, playerTwoColonel, removed[0])
	assert.Equal(t, playerOneGeneral, board.getSquare(Position{5, 5}).piece)
	assert.Nil(t, board.getSquare(Position{6, 5}).piece)

	// Move player-one-general up one space (lose Player One General)
	removed, err = board.MovePiece(Position{5, 5}, Position{4, 5})
	assert.Nil(t, err)
	assert.Len(t, removed, 1)
	assert.Equal(t, playerOneGeneral, removed[0])
	assert.Equal(t, playerTwoMarshal, board.getSquare(Position{4, 5}).piece)
	assert.Nil(t, board.getSquare(Position{5, 5}).piece)

	// Move player-one-marshal right one space (lose both Marshals)
	removed, err = board.MovePiece(Position{4, 4}, Position{4, 5})
	assert.Nil(t, err)
	assert.Len(t, removed, 2)
	assert.Contains(t, removed, playerOneMarshal)
	assert.Contains(t, removed, playerTwoMarshal)
	assert.Nil(t, board.getSquare(Position{4, 4}).piece)
	assert.Nil(t, board.getSquare(Position{4, 5}).piece)
}
