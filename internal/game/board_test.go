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

	err := board.PlacePiece("", Position{0, 0})
	assert.Equal(t, game_errors.ERROR_Board_InvalidMove, err.Code)
}

func TestBoardInitialized(t *testing.T) {
	board := Board{}

	unplayable := []Position{{5, 6}}
	board.Initialize(map[string]*pieces.Piece{}, unplayable)

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
	flag := pieces.CreatePiece(0, pieces.COLOR_red, pieces.RANK_Flag)
	marshal := pieces.CreatePiece(0, pieces.COLOR_red, pieces.RANK_Marshal)
	pieceSet := map[string]*pieces.Piece{
		flag.GetId():    flag,
		marshal.GetId(): marshal,
	}

	board := Board{}

	unplayable := []Position{{5, 6}}

	// Error: Index Out of Range
	err := board.PlacePiece("", Position{20, 0})
	assert.Equal(t, game_errors.ERROR_Board_InvalidMove, err.Code)

	err = board.PlacePiece("", Position{0, 20})
	assert.Equal(t, game_errors.ERROR_Board_InvalidMove, err.Code)

	// Error: Uninitialized Board
	err = board.PlacePiece("", Position{0, 0})
	assert.Equal(t, game_errors.ERROR_Board_InvalidMove, err.Code)

	board.Initialize(pieceSet, unplayable)

	// Success
	validPosition := Position{2, 9}
	err = board.PlacePiece(flag.GetId(), validPosition)

	assert.Nil(t, err)
	assert.Equal(t, flag, board.GetSquare(validPosition).GetPiece())

	// Error: Occupied
	err = board.PlacePiece(marshal.GetId(), validPosition)
	assert.Equal(t, game_errors.ERROR_Board_InvalidMove, err.Code)

	// Error: Unplayable Square
	invalidPosition := Position{5, 6}
	err = board.PlacePiece(marshal.GetId(), invalidPosition)

	assert.NotNil(t, err)
	assert.Equal(t, game_errors.ERROR_Board_InvalidMove, err.Code)
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
	userOne := core.NewPlayer()
	playerOne := NewGamePlayer(pieces.COLOR_red, userOne)
	playerOneGeneral := pieces.CreatePiece(0, playerOne.GetColor(), pieces.RANK_General)
	playerOneMarshal := pieces.CreatePiece(0, playerOne.GetColor(), pieces.RANK_Marshal)

	userTwo := core.NewPlayer()
	playerTwo := NewGamePlayer(pieces.COLOR_blue, userTwo)
	playerTwoColonel := pieces.CreatePiece(0, playerTwo.GetColor(), pieces.RANK_Colonel)
	playerTwoMarshal := pieces.CreatePiece(0, playerTwo.GetColor(), pieces.RANK_Marshal)

	pieces := map[string]*pieces.Piece{
		playerOneGeneral.GetId(): playerOneGeneral,
		playerOneMarshal.GetId(): playerOneMarshal,
		playerTwoColonel.GetId(): playerTwoColonel,
		playerTwoMarshal.GetId(): playerTwoMarshal,
	}
	unplayable := []Position{{5, 6}}
	board.Initialize(pieces, unplayable)

	err := board.PlacePiece(playerOneGeneral.GetId(), Position{6, 6})
	assert.Nil(t, err)
	_, inReserves := board.reserves[playerOneGeneral.GetId()]
	assert.False(t, inReserves)

	err = board.PlacePiece(playerOneMarshal.GetId(), Position{4, 4})
	assert.Nil(t, err)
	_, inReserves = board.reserves[playerOneMarshal.GetId()]
	assert.False(t, inReserves)

	err = board.PlacePiece(playerTwoColonel.GetId(), Position{5, 5})
	assert.Nil(t, err)
	_, inReserves = board.reserves[playerTwoColonel.GetId()]
	assert.False(t, inReserves)

	err = board.PlacePiece(playerTwoMarshal.GetId(), Position{4, 5})
	assert.Nil(t, err)
	_, inReserves = board.reserves[playerTwoMarshal.GetId()]
	assert.False(t, inReserves)

	// Move player-one-general up one space (Error: unplayable)
	response, err := board.MovePiece(Position{6, 6}, Position{5, 6})
	assert.Equal(t, game_errors.ERROR_Board_UnplayableSquare, err.Code)
	assert.Nil(t, response)
	assert.Equal(t, playerOneGeneral, board.GetSquare(Position{6, 6}).GetPiece())

	// Move player-one-general left one space
	response, err = board.MovePiece(Position{6, 6}, Position{6, 5})
	assert.Nil(t, err)
	assert.Nil(t, response.Attackee)
	assert.Equal(t, playerOneGeneral, board.GetSquare(Position{6, 5}).GetPiece())
	assert.Nil(t, board.GetSquare(Position{6, 6}).GetPiece())

	// Move player-one-general up one space (take Player Two Colonel)
	response, err = board.MovePiece(Position{6, 5}, Position{5, 5})
	assert.Nil(t, err)
	assert.Equal(t, playerOneGeneral, response.Attacker)
	assert.Equal(t, playerTwoColonel, response.Attackee)
	assert.Equal(t, MOVERESULT_AttackeeCaptured, response.Move.Result)
	assert.Equal(t, playerOneGeneral, board.GetSquare(Position{5, 5}).GetPiece())
	assert.Nil(t, board.GetSquare(Position{6, 5}).GetPiece())
	_, inReserves = board.reserves[playerTwoColonel.GetId()]
	assert.True(t, inReserves)

	// Move player-one-general up one space (lose Player One General)
	response, err = board.MovePiece(Position{5, 5}, Position{4, 5})
	assert.Nil(t, err)
	assert.Equal(t, playerOneGeneral, response.Attacker)
	assert.Equal(t, playerTwoMarshal, response.Attackee)
	assert.Equal(t, MOVERESULT_AttackerCaptured, response.Move.Result)
	assert.Equal(t, playerTwoMarshal, board.GetSquare(Position{4, 5}).GetPiece())
	assert.Nil(t, board.GetSquare(Position{5, 5}).GetPiece())
	_, inReserves = board.reserves[playerOneGeneral.GetId()]
	assert.True(t, inReserves)

	// Move player-one-marshal right one space (lose both Marshals)
	response, err = board.MovePiece(Position{4, 4}, Position{4, 5})
	assert.Nil(t, err)
	assert.Equal(t, playerOneMarshal, response.Attacker)
	assert.Equal(t, playerTwoMarshal, response.Attackee)
	assert.Equal(t, MOVERESULT_BothCaptured, response.Move.Result)
	assert.Nil(t, board.GetSquare(Position{4, 4}).GetPiece())
	assert.Nil(t, board.GetSquare(Position{4, 5}).GetPiece())
	_, inReserves = board.reserves[playerOneMarshal.GetId()]
	assert.True(t, inReserves)
	_, inReserves = board.reserves[playerTwoMarshal.GetId()]
	assert.True(t, inReserves)
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
