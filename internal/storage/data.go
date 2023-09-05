package storage

import (
	"fmt"

	"github.com/cBiscuitSurprise/strate-go/internal/game"
)

var games *TtlCache[string, game.Game] = NewTtlCache[string, game.Game](10, 86400)

func GetGame(id string) (*game.Game, error) {
	return games.Get(id), nil
}

func SaveGame(id string, game *game.Game) error {
	games.Put(id, game)
	return nil
}

func GetGameForUser(userId string, gameId string) (*game.Game, error) {
	if g, _ := GetGame(gameId); g != nil {
		if p := g.GetPlayerWithId(userId); p != nil {
			return g, nil
		}
	}

	return nil, fmt.Errorf("failed to get game, '%s', for user '%s'", gameId, userId)
}

func ListGamesForUser(userId string) ([]*game.Game, error) {
	all := []*game.Game{}

	capture := func(gameId string, g *game.Game) {
		if p := g.GetPlayerWithId(userId); p != nil {
			all = append(all, g)
		}
	}

	games.ForEach(capture)

	return all, nil
}
