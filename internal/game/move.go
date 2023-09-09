package game

type Move struct {
	Id     string
	From   *Position
	To     *Position
	Result MoveResult
}
