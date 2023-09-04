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

	err := board.PlacePiece(nil, Position{0, 0})
	assert.Equal(t, game_errors.ERROR_Board_Uninitialized, err.Code)
}

func TestBoardInitialized(t *testing.T) {
	board := Board{}

	unplayable := []Position{{5, 6}}
	board.Initialize(unplayable)

	for x, row := range board.squares {
		for y, cell := range row {
			assert.NotNil(t, cell)
			assert.Nil(t, cell.GetPiece())

			if x == 5 && y == 6 {
				assert.False(t, cell.IsPlayable())
			} else {
				assert.True(t, cell.IsPlayable())
			}
		}
	}
}

func TestBoardPlacePiece(t *testing.T) {
	flag := pieces.CreatePiece(pieces.COLOR_red, pieces.RANK_Flag)
	board := Board{}

	unplayable := []Position{{5, 6}}

	// Error: Index Out of Range
	err := board.PlacePiece(nil, Position{20, 0})
	assert.Equal(t, game_errors.ERROR_Board_IndexOutOfRange, err.Code)

	err = board.PlacePiece(nil, Position{0, 20})
	assert.Equal(t, game_errors.ERROR_Board_IndexOutOfRange, err.Code)

	// Error: Uninitialized Board
	err = board.PlacePiece(nil, Position{0, 0})
	assert.Equal(t, game_errors.ERROR_Board_Uninitialized, err.Code)

	board.Initialize(unplayable)

	// Success
	validPosition := Position{2, 9}
	err = board.PlacePiece(flag, validPosition)

	assert.Nil(t, err)
	assert.Equal(t, flag, board.squares[2][9].GetPiece())

	// Error: Occupied
	err = board.PlacePiece(flag, validPosition)

	assert.Equal(t, game_errors.ERROR_Board_OccupiedSquare, err.Code)

	// Error: Unplayable Square
	invalidPosition := Position{5, 6}
	err = board.PlacePiece(flag, invalidPosition)

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
	board.Initialize(unplayable)

	userOne := core.NewPlayer()
	userTwo := core.NewPlayer()

	playerOne := NewGamePlayer(pieces.COLOR_red, userOne)
	playerTwo := NewGamePlayer(pieces.COLOR_blue, userTwo)

	playerOneGeneral := pieces.CreatePiece(playerOne.GetColor(), pieces.RANK_General)
	err := board.PlacePiece(playerOneGeneral, Position{6, 6})
	assert.Nil(t, err)

	playerOneMarshal := pieces.CreatePiece(playerOne.GetColor(), pieces.RANK_Marshal)
	err = board.PlacePiece(playerOneMarshal, Position{4, 4})
	assert.Nil(t, err)

	playerTwoColonel := pieces.CreatePiece(playerTwo.GetColor(), pieces.RANK_Colonel)
	err = board.PlacePiece(playerTwoColonel, Position{5, 5})
	assert.Nil(t, err)

	playerTwoMarshal := pieces.CreatePiece(playerTwo.GetColor(), pieces.RANK_Marshal)
	err = board.PlacePiece(playerTwoMarshal, Position{4, 5})
	assert.Nil(t, err)

	// Move player-one-general up one space (Error: unplayable)
	removed, err := board.MovePiece(Position{6, 6}, Position{5, 6})
	assert.Equal(t, game_errors.ERROR_Board_UnplayableSquare, err.Code)
	assert.Len(t, removed, 0)
	assert.Equal(t, playerOneGeneral, board.GetSquare(Position{6, 6}).GetPiece())

	// Move player-one-general left one space
	removed, err = board.MovePiece(Position{6, 6}, Position{6, 5})
	assert.Nil(t, err)
	assert.Len(t, removed, 0)
	assert.Equal(t, playerOneGeneral, board.GetSquare(Position{6, 5}).GetPiece())
	assert.Nil(t, board.GetSquare(Position{6, 6}).GetPiece())

	// Move player-one-general up one space (take Player Two Colonel)
	removed, err = board.MovePiece(Position{6, 5}, Position{5, 5})
	assert.Nil(t, err)
	assert.Len(t, removed, 1)
	assert.Equal(t, playerTwoColonel, removed[0])
	assert.Equal(t, playerOneGeneral, board.GetSquare(Position{5, 5}).GetPiece())
	assert.Nil(t, board.GetSquare(Position{6, 5}).GetPiece())

	// Move player-one-general up one space (lose Player One General)
	removed, err = board.MovePiece(Position{5, 5}, Position{4, 5})
	assert.Nil(t, err)
	assert.Len(t, removed, 1)
	assert.Equal(t, playerOneGeneral, removed[0])
	assert.Equal(t, playerTwoMarshal, board.GetSquare(Position{4, 5}).GetPiece())
	assert.Nil(t, board.GetSquare(Position{5, 5}).GetPiece())

	// Move player-one-marshal right one space (lose both Marshals)
	removed, err = board.MovePiece(Position{4, 4}, Position{4, 5})
	assert.Nil(t, err)
	assert.Len(t, removed, 2)
	assert.Contains(t, removed, playerOneMarshal)
	assert.Contains(t, removed, playerTwoMarshal)
	assert.Nil(t, board.GetSquare(Position{4, 4}).GetPiece())
	assert.Nil(t, board.GetSquare(Position{4, 5}).GetPiece())
}

func TestCheckNeighbors(t *testing.T) {
	board := Board{}

	radius := 4
	visited := []Position{}
	stopFn := func(distance int, to Position, s *Square) (stop bool) {
		visited = append(visited, to)
		return distance > radius-1
	}

	/* Search 4 spaces in every direction
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |oo|  |  |  |  |  |  |  |
	|  |  |oo|  |  |  |  |  |  |  |
	|  |  |oo|  |  |  |  |  |  |  |
	|  |  |oo|  |  |  |  |  |  |  |
	|oo|oo|st|oo|oo|oo|oo|  |  |  |
	|  |  |oo|  |  |  |  |  |  |  |
	*/
	board.CheckNeighboringSquares(Position{R: 1, C: 2}, stopFn)

	assert.Len(t, visited, 11)
}
