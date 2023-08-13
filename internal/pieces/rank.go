package pieces

type Rank int8

const (
	RANK_Bomb       Rank = 11
	RANK_Marshal    Rank = 10
	RANK_General    Rank = 9
	RANK_Colonel    Rank = 8
	RANK_Major      Rank = 7
	RANK_Captain    Rank = 6
	RANK_Lieutenant Rank = 5
	RANK_Sergent    Rank = 4
	RANK_Minor      Rank = 3
	RANK_Scout      Rank = 2
	RANK_Spy        Rank = 1
	RANK_Flag       Rank = 0
)

func (s Rank) String() string {
	switch s {
	case RANK_Bomb:
		return "Bomb"
	case RANK_Marshal:
		return "Marshal"
	case RANK_General:
		return "General"
	case RANK_Colonel:
		return "Colonel"
	case RANK_Major:
		return "Major"
	case RANK_Captain:
		return "Captain"
	case RANK_Lieutenant:
		return "Lieutenant"
	case RANK_Sergent:
		return "Sergent"
	case RANK_Minor:
		return "Minor"
	case RANK_Scout:
		return "Scount"
	case RANK_Spy:
		return "Spy"
	case RANK_Flag:
		return "Flag"
	}
	return "unknown"
}
