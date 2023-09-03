package game

import (
	"github.com/cBiscuitSurprise/strate-go/internal/core"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
)

type GamePlayer struct {
	color  pieces.Color
	player *core.Player
}

func NewGamePlayer(color pieces.Color, player *core.Player) *GamePlayer {
	return &GamePlayer{
		color:  color,
		player: player,
	}
}

func (p *GamePlayer) GetColor() pieces.Color {
	return p.color
}
