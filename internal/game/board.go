package game

import (
	game_errors "github.com/cBiscuitSurprise/strate-go/internal/errors"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
)

const BOARD_SIZE uint8 = 10
const BOARD_SIZE_X uint8 = BOARD_SIZE
const BOARD_SIZE_Y uint8 = BOARD_SIZE

type Position struct {
	x uint8
	y uint8
}

type Square struct {
	piece    *pieces.Piece
	playable bool
}

type Board struct {
	squares [BOARD_SIZE_X][BOARD_SIZE_Y]*Square
}

func (b *Board) initialize(unplayable []Position) {
	for y, row := range b.squares {
		for x := range row {
			b.squares[x][y] = &Square{playable: true}
		}
	}

	for _, position := range unplayable {
		b.squares[position.x][position.y].playable = false
	}
}

func (b *Board) getSquare(position Position) *Square {
	return b.squares[position.x][position.y]
}

func (b *Board) placePiece(piece *pieces.Piece, position Position) *game_errors.GameError {
	if position.x >= BOARD_SIZE_X || position.y >= BOARD_SIZE_Y {
		return game_errors.GameErrorf(
			game_errors.ERROR_Board_IndexOutOfRange,
			"invalid position: %v for board size: %d rows x %d columns", position, BOARD_SIZE_X, BOARD_SIZE_Y,
		)
	}

	if b.getSquare(position) == nil {
		return game_errors.GameErrorf(
			game_errors.ERROR_Board_Uninitialized,
			"invalid position: %v, board has not been initialized", position,
		)
	}

	if !b.getSquare(position).playable {
		return game_errors.GameErrorf(
			game_errors.ERROR_Board_UnplayableSquare,
			"invalid position: %v, square is not playable", position,
		)
	}

	if b.getSquare(position).piece != nil {
		return game_errors.GameErrorf(
			game_errors.ERROR_Board_OccupiedSquare,
			"invalid position: %v, square is occupied", position,
		)
	}

	b.squares[position.x][position.y].piece = piece

	return nil
}
