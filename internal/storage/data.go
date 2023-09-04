package storage

import "github.com/cBiscuitSurprise/strate-go/internal/game"

var games *TtlCache = NewTtlCache(10, 86400)

func GetGame(id string) (*game.Game, error) {
	return games.Get(id), nil
}

func SaveGame(id string, game *game.Game) error {
	games.Put(id, game)
	return nil
}
