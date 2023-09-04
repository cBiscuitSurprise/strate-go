package game

import (
	"testing"

	"github.com/cBiscuitSurprise/strate-go/internal/core"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	p1 := core.NewPlayer()
	p2 := core.NewPlayer()

	g, err := NewTwoPlayerGame(p1, p2)
	assert.NoError(t, err)

	assert.NotEmpty(t, g.GetId())

	assert.Equal(t, g.GetPlayer(p1.GetId()).player.GetId(), p1.GetId())
	assert.Equal(t, g.GetPlayer(p1.GetId()).color, pieces.COLOR_red)
	assert.Equal(t, g.GetPlayer(p2.GetId()).player.GetId(), p2.GetId())
	assert.Equal(t, g.GetPlayer(p2.GetId()).color, pieces.COLOR_blue)
}

func TestGamePlacement(t *testing.T) {
	p1 := core.NewPlayer()
	p2 := core.NewPlayer()

	g, err := NewTwoPlayerGame(p1, p2)
	assert.NoError(t, err)

	g.SetMode(GAMEMODE_Plan)

	_, err = placePieces(g.Board, g.pieces[p1.GetId()], 0)
	assert.NoError(t, err)

	_, err = placePieces(g.Board, g.pieces[p2.GetId()], 6)
	assert.NoError(t, err)
}
