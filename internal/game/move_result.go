package game

type MoveResult int8

const (
	MOVERESULT_NoContest MoveResult = iota
	MOVERESULT_AttackeeCaptured
	MOVERESULT_AttackerCaptured
	MOVERESULT_BothCaptured
)
