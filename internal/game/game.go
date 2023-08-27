package game

import "github.com/cBiscuitSurprise/strate-go/internal/core"

type Game struct {
	Players []*core.Player
	Board   *Board
}
