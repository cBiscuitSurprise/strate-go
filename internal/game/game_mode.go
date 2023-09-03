package game

type GameMode int8

const (
	GAMEMODE_Setup GameMode = 0
	GAMEMODE_Plan  GameMode = 1
	GAMEMODE_Play  GameMode = 2
	GAMEMODE_End   GameMode = 3
	GAMEMODE_Error GameMode = 4
)

func (s GameMode) String() string {
	switch s {
	case GAMEMODE_Setup:
		return "Setup"
	case GAMEMODE_Plan:
		return "Plan"
	case GAMEMODE_Play:
		return "Play"
	case GAMEMODE_End:
		return "End"
	case GAMEMODE_Error:
		return "Error"
	}
	return "Unknown"
}
