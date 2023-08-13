package game

import (
	"testing"

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
	flag := pieces.CreatePiece(pieces.RANK_Flag)
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
