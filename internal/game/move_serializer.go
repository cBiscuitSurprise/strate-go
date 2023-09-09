package game

import (
	"fmt"
)

func SerializePosition(p *Position) string {
	if p == nil {
		return "-"
	} else {
		return fmt.Sprintf("%d,%d", p.R, p.C)
	}
}

func DeserializePosition(p string) (*Position, error) {
	if p == "-" {
		return nil, nil
	} else {
		var r, c int
		if n, err := fmt.Sscanf(p, "%d,%d", &r, &c); err != nil {
			return nil, err
		} else if n == 2 {
			return &Position{
				R: r,
				C: c,
			}, nil
		} else {
			return nil, fmt.Errorf("failed to parse serialized position %s, invalid string", p)
		}
	}
}

func SerializeMove(move Move) string {
	return fmt.Sprintf(
		"%s %s %s %d",
		move.Id,
		SerializePosition(move.From),
		SerializePosition(move.To),
		move.Result,
	)
}

func DeserializeMove(move string) (Move, error) {
	var id string
	var from string
	var to string
	var result MoveResult

	if n, err := fmt.Sscanf(move, "%s %s %s %d", &id, &from, &to, &result); err != nil {
		return Move{}, err
	} else if n == 4 {
		from, err := DeserializePosition(from)
		if err != nil {
			return Move{}, err
		}
		to, err := DeserializePosition(to)
		if err != nil {
			return Move{}, err
		}

		return Move{
			Id:     id,
			From:   from,
			To:     to,
			Result: result,
		}, nil
	} else {
		return Move{}, fmt.Errorf("failed to parse serialized move %s, invalid string", move)
	}
}
