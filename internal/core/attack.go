package core

type Winner int8

const (
	WINNER_Attackee Winner = -1
	WINNER_Draw     Winner = 0
	WINNER_Attacker Winner = 1
)

func (w Winner) String() string {
	switch w {
	case WINNER_Attackee:
		return "attackee"
	case WINNER_Draw:
		return "draw"
	case WINNER_Attacker:
		return "attacker"
	}
	return "unknown"
}

func (w Winner) Invert() Winner {
	switch w {
	case WINNER_Attackee:
		return WINNER_Attacker
	case WINNER_Attacker:
		return WINNER_Attackee
	case WINNER_Draw:
		return WINNER_Draw
	}
	return w
}
