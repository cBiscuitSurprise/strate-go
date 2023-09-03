package core

import (
	"github.com/cBiscuitSurprise/strate-go/internal/util"
)

type Player struct {
	id string
}

func NewPlayer() *Player {
	return &Player{
		id: util.NewId(),
	}
}

func HydratePlayer(id string) *Player {
	return &Player{
		id: id,
	}
}

func (p *Player) GetId() string {
	return p.id
}
