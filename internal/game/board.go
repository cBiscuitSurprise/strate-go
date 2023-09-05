package game

import (
	"github.com/cBiscuitSurprise/strate-go/internal/core"
	game_errors "github.com/cBiscuitSurprise/strate-go/internal/errors"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
)

const BOARD_SIZE int = 10
const BOARD_SIZE_X int = BOARD_SIZE
const BOARD_SIZE_Y int = BOARD_SIZE

type Position struct {
	R int
	C int
}

type Square struct {
	piece    *pieces.Piece
	playable bool
}

func (s *Square) GetPiece() *pieces.Piece {
	return s.piece
}

func (s *Square) SetPiece(piece *pieces.Piece) {
	s.piece = piece
}

func (s *Square) RemovePiece() {
	s.piece = nil
}

func (s *Square) IsPlayable() bool {
	return s.playable
}

type BoardSize struct {
	Rows    int
	Columns int
}

type Board struct {
	squares [BOARD_SIZE_X][BOARD_SIZE_Y]*Square
}

func (b *Board) GetSize() BoardSize {
	return BoardSize{Rows: BOARD_SIZE_X, Columns: BOARD_SIZE_Y}
}

func (b *Board) Initialize(unplayable []Position) {
	for y, row := range b.squares {
		for x := range row {
			b.squares[x][y] = &Square{playable: true}
		}
	}

	for _, position := range unplayable {
		b.squares[position.R][position.C].playable = false
	}
}

func (b *Board) GetSquare(position Position) *Square {
	return b.squares[position.R][position.C]
}

func (b *Board) Rows() [BOARD_SIZE_X][BOARD_SIZE_Y]*Square {
	return b.squares
}

func (b *Board) RemovePiece(position Position) {
	b.squares[position.R][position.C].RemovePiece()
}

func (b *Board) CheckNeighboringSquares(position Position, stop func(d int, p Position, s *Square) (stop bool)) {
	for i := 1; i < b.GetSize().Rows; i++ {
		p := b.lookForward(position, i)
		if p == nil || stop(i, *p, b.GetSquare(*p)) {
			break
		}
	}
	for i := 1; i < b.GetSize().Rows; i++ {
		p := b.lookBackward(position, i)
		if p == nil || stop(i, *p, b.GetSquare(*p)) {
			break
		}
	}
	for i := 1; i < b.GetSize().Columns; i++ {
		p := b.lookRight(position, i)
		if p == nil || stop(i, *p, b.GetSquare(*p)) {
			break
		}
	}
	for i := 1; i < b.GetSize().Columns; i++ {
		p := b.lookLeft(position, i)
		if p == nil || stop(i, *p, b.GetSquare(*p)) {
			break
		}
	}
}

func (b *Board) PlacePiece(piece *pieces.Piece, position Position) *game_errors.GameError {
	if position.R >= BOARD_SIZE_X || position.C >= BOARD_SIZE_Y {
		return game_errors.GameErrorf(
			game_errors.ERROR_Board_IndexOutOfRange,
			"invalid position: %v for board size: %d rows x %d columns", position, BOARD_SIZE_X, BOARD_SIZE_Y,
		)
	}

	if b.GetSquare(position) == nil {
		return game_errors.GameErrorf(
			game_errors.ERROR_Board_Uninitialized,
			"invalid position: %v, board has not been initialized", position,
		)
	}

	if !b.GetSquare(position).playable {
		return game_errors.GameErrorf(
			game_errors.ERROR_Board_UnplayableSquare,
			"invalid position: %v, square is not playable", position,
		)
	}

	if b.GetSquare(position).GetPiece() != nil {
		return game_errors.GameErrorf(
			game_errors.ERROR_Board_OccupiedSquare,
			"invalid position: %v, square is occupied", position,
		)
	}

	b.squares[position.R][position.C].SetPiece(piece)

	return nil
}

func (b *Board) MovePiece(from Position, to Position) ([]*pieces.Piece, *game_errors.GameError) {
	fromSquare := b.GetSquare(from)
	toSquare := b.GetSquare(to)

	if fromSquare == nil || toSquare == nil {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Board_Uninitialized,
			"invalid positions: from %v, to %v, board has not been initialized", from, to,
		)
	}

	if fromSquare.GetPiece() == nil {
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
	if toSquare.GetPiece() == nil {
		// no contest
		b.squares[to.R][to.C].SetPiece(fromSquare.GetPiece())
		b.squares[from.R][from.C].RemovePiece()
	} else {
		winner, err := fromSquare.GetPiece().Attack(toSquare.GetPiece())

		if err != nil {
			return nil, game_errors.GameErrorf(
				game_errors.ERROR_Board_UnplayableSquare,
				"invalid to position: %v, square is not playable", to,
			)
		}

		if winner == core.WINNER_Attacker {
			losingPieces = append(losingPieces, toSquare.GetPiece())
			b.squares[to.R][to.C].SetPiece(fromSquare.GetPiece())
		} else if winner == core.WINNER_Attackee {
			losingPieces = append(losingPieces, fromSquare.GetPiece())
		} else {
			// both pieces removed from board
			losingPieces = append(losingPieces, toSquare.GetPiece())
			losingPieces = append(losingPieces, fromSquare.GetPiece())
			b.squares[to.R][to.C].RemovePiece()
		}
		b.squares[from.R][from.C].RemovePiece()
	}

	return losingPieces, nil
}

func (b *Board) lookForward(from Position, count int) *Position {
	n := from.R + count
	if n >= b.GetSize().Rows {
		return nil
	}
	return &Position{R: n, C: from.C}
}

func (b *Board) lookBackward(from Position, count int) *Position {
	n := from.R - count
	if n < 0 {
		return nil
	}
	return &Position{R: n, C: from.C}
}

func (b *Board) lookRight(from Position, count int) *Position {
	n := from.C + count
	if n >= b.GetSize().Columns {
		return nil
	}
	return &Position{R: from.R, C: n}
}

func (b *Board) lookLeft(from Position, count int) *Position {
	n := from.C - count
	if n < 0 {
		return nil
	}
	return &Position{R: from.R, C: n}
}
