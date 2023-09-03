package game_errors

import "fmt"

type ErrorCode int64

const (
	ERROR_Board_IndexOutOfRange ErrorCode = iota
	ERROR_Board_Uninitialized
	ERROR_Board_UnplayableSquare
	ERROR_Board_OccupiedSquare

	ERROR_Contest_InvalidContest

	ERROR_Game_InvalidMode
	ERROR_Game_InvalidPiece
)

type GameError struct {
	msg  string    // description of error
	Code ErrorCode // error code
}

func (e *GameError) Error() string { return e.msg }

func GameErrorf(code ErrorCode, format string, a ...any) *GameError {
	return &GameError{Code: code, msg: fmt.Sprintf(format, a...)}
}
