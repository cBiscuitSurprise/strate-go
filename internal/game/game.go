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
	pieces  map[string]map[string]*pieces.Piece
	mode    GameMode
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
		pieces: map[string]map[string]*pieces.Piece{
			red.GetId():  pieces.GenerateStandardPieces(pieces.COLOR_red),
			blue.GetId(): pieces.GenerateStandardPieces(pieces.COLOR_blue),
		},
		mode:  GAMEMODE_Setup,
		Board: CreateStandardBaseBoard(),
	}, nil
}

func (g *Game) GetId() string {
	return g.id
}

func (g *Game) GetPlayer(id string) *GamePlayer {
	return g.players[id]
}

func (g *Game) GetMode() GameMode {
	return g.mode
}

func (g *Game) SetMode(mode GameMode) {
	g.mode = mode
}

func (g *Game) GetValidMovesFromPosition(player_id string, from Position) []Position {
	piece := g.Board.GetSquare(from).GetPiece()

	validMoves := []Position{}
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
				validMoves = append(validMoves, to)
			} else {
				if s.GetPiece().GetColor() != piece.GetColor() {
					validMoves = append(validMoves, to)
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

	if p, ok := g.pieces[player_id][piece_id]; ok {
		if err := g.Board.PlacePiece(p, position); err == nil {
			return nil
		} else {
			return err
		}
	} else {
		return game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidPiece,
			"piece, '%s', does not belong to player, '%s'!", piece_id, player_id,
		)
	}
}

func (g *Game) MovePiece(player_id string, from Position, to Position) ([]*pieces.Piece, *game_errors.GameError) {
	if g.mode != GAMEMODE_Play {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidMode,
			"game is not in playing mode, cannot proceed (%s)", g.mode.String(),
		)
	}

	playerColor := g.players[player_id].GetColor()

	// player needs to own `from` piece
	fromPiece := g.Board.GetSquare(from).GetPiece()
	if fromPiece == nil {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidPiece,
			"no piece at position, %v!", from,
		)
	} else if fromPiece.GetColor() != playerColor {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidPiece,
			"piece, '%v', does not belong to player, '%s'!", fromPiece, player_id,
		)
	}

	// player can't to own `to` piece
	toPiece := g.Board.GetSquare(to).GetPiece()
	if toPiece != nil && toPiece.GetColor() == playerColor {
		return nil, game_errors.GameErrorf(
			game_errors.ERROR_Game_InvalidPiece,
			"piece, '%v', already belongs to player, '%s', invalid move!", toPiece, player_id,
		)
	}

	return g.Board.MovePiece(from, to)
}
