package game

import (
	"fmt"

	"github.com/cBiscuitSurprise/strate-go/internal/core"
	game_errors "github.com/cBiscuitSurprise/strate-go/internal/errors"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
	"github.com/cBiscuitSurprise/strate-go/internal/util"
	"github.com/rs/zerolog/log"
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

func (s *Square) RemovePiece() *pieces.Piece {
	p := s.piece
	s.piece = nil
	return p
}

func (s *Square) IsPlayable() bool {
	return s != nil && s.playable
}

type BoardSize struct {
	Rows    int
	Columns int
}

type MovePieceResponse struct {
	Attacker *pieces.Piece
	Attackee *pieces.Piece
	Move     Move
}

var emptyMovePieceResponse MovePieceResponse = MovePieceResponse{}

type Board struct {
	squares  [BOARD_SIZE_X][BOARD_SIZE_Y]*Square
	pieces   map[string]*pieces.Piece
	reserves map[string]*pieces.Piece
	history  []Move
}

func (b *Board) GetSize() BoardSize {
	return BoardSize{Rows: BOARD_SIZE_X, Columns: BOARD_SIZE_Y}
}

func (b *Board) GetPiece(id string) (pieces.Piece, bool) {
	if p, onBoard := b.pieces[id]; onBoard {
		return *p, true
	} else {
		return pieces.Piece{}, false
	}
}

func (b *Board) Initialize(pieceSet map[string]*pieces.Piece, unplayable []Position) {
	b.pieces = util.UpdateMap[string, *pieces.Piece](make(map[string]*pieces.Piece, len(pieceSet)), pieceSet)
	b.reserves = pieceSet

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

func (b *Board) PlacePiece(id string, position Position) *game_errors.GameError {
	if ok := b.ApplyMove(Move{
		Id:     id,
		To:     &position,
		Result: MOVERESULT_NoContest,
	}); !ok {
		return game_errors.GameErrorf(
			game_errors.ERROR_Board_InvalidMove,
			"invalid move: failed to move from reserves",
		)
	}

	return nil
}

func (b *Board) MovePiece(from Position, to Position) (*MovePieceResponse, *game_errors.GameError) {
	fromSquare := b.GetSquare(from)
	toSquare := b.GetSquare(to)

	attacker := fromSquare.GetPiece()
	attackee := toSquare.GetPiece()

	if fromSquare == nil || toSquare == nil {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Board_Uninitialized,
			"invalid positions: from %v, to %v, board has not been initialized", from, to,
		)
	}

	if attacker == nil {
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

	response := &MovePieceResponse{
		Attacker: attacker,
		Attackee: attackee,
	}

	response.Move = Move{
		Id:     attacker.GetId(),
		From:   &from,
		To:     &to,
		Result: b.computeMoveResult(attacker, attackee),
	}
	if ok := b.ApplyMove(response.Move); ok {
		return response, nil
	} else {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Board_UnplayableSquare,
			"failed to apply move %v", response.Move,
		)
	}
}

func (b *Board) ApplyMove(move Move) (ok bool) {
	ok = true
	if ok = b.ValidateMove(move); !ok {
		return ok
	}

	if move.From == nil && move.To == nil {
		ok = move.Result == MOVERESULT_NoContest
	} else if move.From == nil {
		b.fromReserves(move.Id, *move.To)
	} else if move.To == nil {
		b.toReserves(*move.From)
	} else {
		switch move.Result {
		case MOVERESULT_NoContest:
			ok = ok && b.transferPiece(*move.From, *move.To)
			break
		case MOVERESULT_AttackeeCaptured:
			ok = ok && b.toReserves(*move.To)
			ok = ok && b.transferPiece(*move.From, *move.To)
			break
		case MOVERESULT_AttackerCaptured:
			ok = ok && b.toReserves(*move.From)
			break
		case MOVERESULT_BothCaptured:
			ok = ok && b.toReserves(*move.To)
			ok = ok && b.toReserves(*move.From)
			break
		}
	}

	if ok {
		b.history = append(b.history, move)
	}
	return ok
}

func (b *Board) ValidateMove(move Move) (ok bool) {
	// Validates that the move can be made on the board (attack logic is handled elsewhere)
	if move.From == nil && move.To == nil {
		// noop
		ok = true
	} else {
		ok = true
		if move.To != nil {
			// `To` shall be playable
			ok = ok && b.hasPosition(*move.To) && b.GetSquare(*move.To).IsPlayable()

			if ok && b.GetSquare(*move.To).GetPiece() != nil {
				// There shall be a contest
				ok = ok && move.Result != MOVERESULT_NoContest
			}
		} else {
			// Piece (`Id`) shall not be in reserves
			_, inReserve := b.reserves[move.Id]
			ok = ok && !inReserve
		}
		if move.From != nil {
			// `From` shall be playable and equal to Move.Id
			ok = (ok &&
				b.hasPosition(*move.From) &&
				b.GetSquare(*move.From).IsPlayable() &&
				b.GetSquare(*move.From).GetPiece() != nil &&
				b.GetSquare(*move.From).GetPiece().GetId() == move.Id)
		} else {
			// Piece (`Id`) shall be in reserves
			_, inReserve := b.reserves[move.Id]
			ok = ok && inReserve
		}
	}
	return ok
}

func (b *Board) hasPosition(position Position) bool {
	if position.R >= BOARD_SIZE_X || position.C >= BOARD_SIZE_Y {
		return false
	}
	return true
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

func (b *Board) fromReserves(id string, to Position) (ok bool) {
	if piece, inReserve := b.reserves[id]; inReserve {
		b.GetSquare(to).SetPiece(piece)
		delete(b.reserves, id)
		return true
	}
	log.Error().
		Err(fmt.Errorf("invalid move from reserves")).
		Msgf("piece, '%s', is not in reserves", id)
	return false
}

func (b *Board) toReserves(from Position) (ok bool) {
	if piece := b.GetSquare(from).RemovePiece(); piece != nil {
		b.reserves[piece.GetId()] = piece
		return true
	}
	log.Error().
		Err(fmt.Errorf("invalid move to reserves")).
		Msgf("no piece located at, position '%v'", from)
	return false
}

func (b *Board) transferPiece(from Position, to Position) (ok bool) {
	if toSquare := b.GetSquare(to); toSquare.GetPiece() == nil {
		if fromPiece := b.GetSquare(from).RemovePiece(); fromPiece != nil {
			toSquare.SetPiece(fromPiece)
			return true
		}
		log.Error().
			Err(fmt.Errorf("invalid move transferring piece")).
			Msgf("no piece located at, position '%v'", from)
		return false
	}
	log.Error().
		Err(fmt.Errorf("invalid move transferring piece")).
		Msgf("piece already located at, position '%v'", to)
	return false
}

func (b *Board) computeMoveResult(attacker, attackee *pieces.Piece) MoveResult {
	if attackee == nil {
		return MOVERESULT_NoContest
	} else {
		result := attacker.Attack(attackee)
		switch result {
		case core.WINNER_Attacker:
			return MOVERESULT_AttackeeCaptured
		case core.WINNER_Attackee:
			return MOVERESULT_AttackerCaptured
		case core.WINNER_Draw:
			return MOVERESULT_BothCaptured
		}
		log.Warn().
			Str("winner", result.String()).
			Msgf("unhandled attack result")
		return MOVERESULT_NoContest
	}
}
