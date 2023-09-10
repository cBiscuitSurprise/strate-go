package storage

import "fmt"

const (
	strategoNamespace = "sg"
	gameNamespace     = "g"
)

func GameNamespace(gameId string) string {
	return fmt.Sprintf("%s:%s:%s", strategoNamespace, gameNamespace, gameId)
}

func GameMoveStreamKey(gameId string) string {
	return fmt.Sprintf("%s:moves", GameNamespace(gameId))
}
