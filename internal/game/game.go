package game

import (
	"github.com/cBiscuitSurprise/strate-go/internal/core"
	game_errors "github.com/cBiscuitSurprise/strate-go/internal/errors"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
	"github.com/cBiscuitSurprise/strate-go/internal/util"
)

type Game struct {
	Board *Board

	id      string
	players map[string]*GamePlayer
	mode    GameMode
	nonce   int
}

func NewTwoPlayerGame(red *core.Player, blue *core.Player) (*Game, error) {
	redGamePlayer := NewGamePlayer(pieces.COLOR_red, red)
	blueGamePlayer := NewGamePlayer(pieces.COLOR_blue, blue)

	return &Game{
		id: util.NewId(),
		players: map[string]*GamePlayer{
			red.GetId():  redGamePlayer,
			blue.GetId(): blueGamePlayer,
		},
		mode:  GAMEMODE_Setup,
		Board: CreateStandardBaseBoard(CreateStandardPieceSet()),
	}, nil
}

func (g *Game) GetId() string {
	return g.id
}

func (g *Game) GetMode() GameMode {
	return g.mode
}

func (g *Game) SetMode(mode GameMode) {
	g.mode = mode
}

func (g *Game) GetNonce() int {
	return g.nonce
}

func (g *Game) GetPlayerWithId(id string) *GamePlayer {
	for _, p := range g.players {
		if p.player.GetId() == id {
			return p
		}
	}
	return nil
}

func (g *Game) ApplyMove(nonce int, move Move) (ok bool) {
	return false
}

func (g *Game) GetValidMovesFromPosition(playerId string, from Position) []*Position {
	piece := g.Board.GetSquare(from).GetPiece()

	validMoves := []*Position{}
	if piece == nil {
		return validMoves
	}

	if piece.GetMaxMoves() == 0 {
		return validMoves
	}

	findValidMove := func(distance int, to Position, s *Square) (stop bool) {
		if distance > piece.GetMaxMoves()-1 {
			stop = true
		}
		if s.IsPlayable() {
			if s.GetPiece() == nil {
				validMoves = append(validMoves, &to)
			} else {
				if s.GetPiece().GetColor() != piece.GetColor() {
					validMoves = append(validMoves, &to)
				}
				stop = true
			}
		} else {
			stop = true
		}
		return stop
	}

	g.Board.CheckNeighboringSquares(from, findValidMove)

	return validMoves
}

func (g *Game) PlacePiece(player_id string, piece_id string, position Position) *game_errors.GameError {
	if g.mode != GAMEMODE_Plan {
		return game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidMode,
			"game is not in planning mode, cannot proceed (%s)", g.mode.String(),
		)
	}

	if ok := g.checkOwnership(player_id, piece_id); ok {
		return g.Board.PlacePiece(piece_id, position)
	} else {
		return game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidMove,
			"piece, '%s', does not belong to player, '%s'!", piece_id, player_id,
		)
	}
}

func (g *Game) MovePiece(playerId string, from Position, to Position) (*MovePieceResponse, *game_errors.GameError) {
	if g.mode != GAMEMODE_Play {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidMode,
			"game is not in playing mode, cannot proceed (%s)", g.mode.String(),
		)
	}

	fromPiece := g.Board.GetSquare(from).GetPiece()

	// fromPiece needs to exists
	if fromPiece == nil {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidMove,
			"no piece to move from %v!", from,
		)
	}

	// fromPiece needs to be movable
	if fromPiece.GetMaxMoves() == 0 {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidMove,
			"can't move piece, %s, at position (max-moves == 0), %v!", fromPiece.GetRank().String(), from,
		)
	}

	// player needs to own `from` piece
	if fromPiece == nil {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidMove,
			"no piece at position, %v!", from,
		)
	} else if !g.checkOwnership(playerId, fromPiece.GetId()) {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidMove,
			"piece, '%v', does not belong to player, '%s'!", fromPiece, playerId,
		)
	}

	// `to` Square must be playable
	toSquare := g.Board.GetSquare(to)
	if !toSquare.IsPlayable() {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidMove,
			"can't move to square '%v', unplayable square!", to,
		)
	}

	// `to` square must be a valid move
	if !g.validateMoveTo(playerId, from, to) {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidMove,
			"piece, '%v', can't move to, '%v', invalid move (max distance = %d)!", fromPiece, to, fromPiece.GetMaxMoves(),
		)
	}

	// player can't to own `to` piece
	toPiece := toSquare.GetPiece()
	if toPiece != nil && g.checkOwnership(playerId, toPiece.GetId()) {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidMove,
			"piece, '%v', already belongs to player, '%s', invalid move!", toPiece, playerId,
		)
	}

	response, err := g.Board.MovePiece(from, to)
	if err == nil {
		g.nonce += 1
	}
	return response, err
}

func (g *Game) checkOwnership(playerId string, pieceId string) bool {
	if p, onBoard := g.Board.GetPiece(pieceId); onBoard {
		if player, inGame := g.players[playerId]; inGame {
			return p.GetColor() == player.GetColor()
		}
	}
	return false
}

func (g *Game) validateMoveTo(playerId string, from Position, to Position) bool {
	for _, validTo := range g.GetValidMovesFromPosition(playerId, from) {
		if to == *validTo {
			return true
		}
	}
	return false
}
