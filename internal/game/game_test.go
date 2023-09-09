package game

import (
	"testing"

	"github.com/cBiscuitSurprise/strate-go/internal/core"
	game_errors "github.com/cBiscuitSurprise/strate-go/internal/errors"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	p1 := core.NewPlayer()
	p2 := core.NewPlayer()

	g, err := NewTwoPlayerGame(p1, p2)
	assert.NoError(t, err)

	assert.NotEmpty(t, g.GetId())

	assert.Equal(t, g.GetPlayerWithId(p1.GetId()).player.GetId(), p1.GetId())
	assert.Equal(t, g.GetPlayerWithId(p1.GetId()).color, pieces.COLOR_red)
	assert.Equal(t, g.GetPlayerWithId(p2.GetId()).player.GetId(), p2.GetId())
	assert.Equal(t, g.GetPlayerWithId(p2.GetId()).color, pieces.COLOR_blue)
}

func TestGamePlacement(t *testing.T) {
	p1 := core.NewPlayer()
	p2 := core.NewPlayer()

	g, err := NewTwoPlayerGame(p1, p2)
	assert.Nil(t, err)

	// Try to Place Spy on 0,0 (game in wrong mode)
	gerr := g.PlacePiece(p1.GetId(), "Red:01:00", Position{R: 0, C: 0})
	if assert.NotNil(t, gerr) {
		assert.Equal(t, game_errors.ERROR_Game_InvalidMode, gerr.Code)
	}

	g.SetMode(GAMEMODE_Plan)

	// Place Spy on 0,0
	gerr = g.PlacePiece(p1.GetId(), "Red:01:00", Position{R: 0, C: 0})
	assert.Nil(t, gerr)

	// Try to Place Same Spy on 0,1 (already placed)
	gerr = g.PlacePiece(p1.GetId(), "Red:01:00", Position{R: 0, C: 0})
	if assert.NotNil(t, gerr) {
		assert.Equal(t, game_errors.ERROR_Board_InvalidMove, gerr.Code)
	}

	// Try to Place Scout on 0,0 (occupied)
	gerr = g.PlacePiece(p1.GetId(), "Red:02:00", Position{R: 0, C: 0})
	if assert.NotNil(t, gerr) {
		assert.Equal(t, game_errors.ERROR_Board_InvalidMove, gerr.Code)
	}

	// Red Player try placing Blue General on 4,1
	gerr = g.PlacePiece(p1.GetId(), "Blue:09:00", Position{R: 4, C: 1})
	if assert.NotNil(t, gerr) {
		assert.Equal(t, game_errors.ERROR_Game_InvalidMove, gerr.Code)
	}
}

func TestGameMove(t *testing.T) {
	// #region setup
	p1 := core.HydratePlayer("TESTREDPLAYER")
	p2 := core.HydratePlayer("TESTBLUEPLAYER")

	g, err := NewTwoPlayerGame(p1, p2)
	assert.Nil(t, err)

	g.SetMode(GAMEMODE_Plan)

	/* Setup Board
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |RM|  |  |  |  |  |  |  |
	|  |BG|xx|xx|  |  |xx|xx|  |  |
	|  |BM|xx|xx|  |  |xx|xx|  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	|  |  |  |  |  |  |  |  |  |  |
	*/

	// Place Red Marshall on 3,2
	gerr := g.PlacePiece(p1.GetId(), "Red:10:00", Position{R: 3, C: 2})
	assert.Nil(t, gerr)

	// Place Blue General on 4,1
	gerr = g.PlacePiece(p2.GetId(), "Blue:09:00", Position{R: 4, C: 1})
	assert.Nil(t, gerr)

	// Place Blue Marshall on 5,1
	gerr = g.PlacePiece(p2.GetId(), "Blue:10:00", Position{R: 5, C: 1})
	assert.Nil(t, gerr)

	// #endregion setup

	// Try moving in Plan Mode
	response, gerr := g.MovePiece(p1.GetId(), Position{R: 0, C: 0}, Position{R: 0, C: 0})
	assert.Nil(t, response)
	if assert.NotNil(t, gerr) {
		assert.Equal(t, game_errors.ERROR_Game_InvalidMode, gerr.Code)
	}

	g.SetMode(GAMEMODE_Play)

	// Try moving to invalid square
	response, gerr = g.MovePiece(p1.GetId(), Position{R: 3, C: 2}, Position{R: 4, C: 2})
	assert.Nil(t, response)
	if assert.NotNil(t, gerr) {
		assert.Equal(t, game_errors.ERROR_Board_UnplayableSquare, gerr.Code)
	}

	// Try moving from empty square
	response, gerr = g.MovePiece(p1.GetId(), Position{R: 0, C: 0}, Position{R: 0, C: 1})
	assert.Nil(t, response)
	if assert.NotNil(t, gerr) {
		assert.Equal(t, game_errors.ERROR_Game_InvalidMove, gerr.Code)
	}

	// Move red marshall one to the left
	response, gerr = g.MovePiece(p1.GetId(), Position{R: 3, C: 2}, Position{R: 3, C: 1})
	assert.Nil(t, gerr)
	assert.Equal(t, "Red:10:00", response.Attacker.GetId())
	assert.Nil(t, response.Attackee)
	assert.Equal(t, "Red:10:00", response.Move.Id)
	assert.Equal(t, Position{R: 3, C: 2}, *response.Move.From)
	assert.Equal(t, Position{R: 3, C: 1}, *response.Move.To)
	assert.Equal(t, MOVERESULT_NoContest, response.Move.Result)

	// Move Blue General up one, losing to Red Marshal
	response, gerr = g.MovePiece(p2.GetId(), Position{R: 4, C: 1}, Position{R: 3, C: 1})
	assert.Nil(t, gerr)
	assert.Equal(t, "Blue:09:00", response.Attacker.GetId())
	assert.Equal(t, "Red:10:00", response.Attackee.GetId())
	assert.Equal(t, "Blue:09:00", response.Move.Id)
	assert.Equal(t, Position{R: 4, C: 1}, *response.Move.From)
	assert.Equal(t, Position{R: 3, C: 1}, *response.Move.To)
	assert.Equal(t, MOVERESULT_AttackerCaptured, response.Move.Result)

	// Move Red Marshal down one
	response, gerr = g.MovePiece(p1.GetId(), Position{R: 3, C: 1}, Position{R: 4, C: 1})
	assert.Nil(t, gerr)
	assert.Equal(t, "Red:10:00", response.Attacker.GetId())
	assert.Nil(t, response.Attackee)
	assert.Equal(t, "Red:10:00", response.Move.Id)
	assert.Equal(t, Position{R: 3, C: 1}, *response.Move.From)
	assert.Equal(t, Position{R: 4, C: 1}, *response.Move.To)
	assert.Equal(t, MOVERESULT_NoContest, response.Move.Result)

	// Move Blue Marshal up one, losing both Marshals
	response, gerr = g.MovePiece(p2.GetId(), Position{R: 5, C: 1}, Position{R: 4, C: 1})
	assert.Nil(t, gerr)
	assert.Equal(t, "Blue:10:00", response.Attacker.GetId())
	assert.Equal(t, "Red:10:00", response.Attackee.GetId())
	assert.Equal(t, "Blue:10:00", response.Move.Id)
	assert.Equal(t, Position{R: 5, C: 1}, *response.Move.From)
	assert.Equal(t, Position{R: 4, C: 1}, *response.Move.To)
	assert.Equal(t, MOVERESULT_BothCaptured, response.Move.Result)
}
