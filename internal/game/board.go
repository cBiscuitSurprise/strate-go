package game

import (
	"github.com/cBiscuitSurprise/strate-go/internal/core"
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

func (b *Board) MovePiece(from Position, to Position) ([]*pieces.Piece, *game_errors.GameError) {
	fromSquare := b.getSquare(from)
	toSquare := b.getSquare(to)

	if fromSquare == nil || toSquare == nil {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Board_Uninitialized,
			"invalid positions: from %v, to %v, board has not been initialized", from, to,
		)
	}

	if fromSquare.piece == nil {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Board_Uninitialized,
			"invalid from position: %v, there is no pieces here", from,
		)
	}

	if !toSquare.playable {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Board_UnplayableSquare,
			"invalid to position: %v, square is not playable", to,
		)
	}

	var losingPieces []*pieces.Piece
	if toSquare.piece == nil {
		// no contest
		b.squares[to.x][to.y].piece = fromSquare.piece
		b.squares[from.x][from.y].piece = nil
	} else {
		winner, err := fromSquare.piece.Attack(toSquare.piece)

		if err != nil {
			return nil, game_errors.GameErrorf(
				game_errors.ERROR_Board_UnplayableSquare,
				"invalid to position: %v, square is not playable", to,
			)
		}

		if winner == core.WINNER_Attacker {
			losingPieces = append(losingPieces, toSquare.piece)
			b.squares[to.x][to.y].piece = fromSquare.piece
		} else if winner == core.WINNER_Attackee {
			losingPieces = append(losingPieces, fromSquare.piece)
		} else {
			// both pieces removed from board
			losingPieces = append(losingPieces, toSquare.piece)
			losingPieces = append(losingPieces, fromSquare.piece)
			b.squares[to.x][to.y].piece = nil
		}
		b.squares[from.x][from.y].piece = nil
	}

	return losingPieces, nil
}
