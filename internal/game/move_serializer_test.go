package game

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerialMoveFull(t *testing.T) {
	move := Move{
		Id:           "Red:001",
		From:         &Position{R: 1, C: 2},
		To:           &Position{R: 1, C: 3},
		Result:       MOVERESULT_BothCaptured,
		AttackerRank: 2,
		AttackeeRank: 2,
	}

	serialized := SerializeMove(move)

	require.Equal(t, "Red:001 1,2 1,3 3 2 2", serialized)
}

func TestDeserialMoveFull(t *testing.T) {
	serialized := "Red:001 1,2 1,3 3 10 10"

	move, err := DeserializeMove(serialized)
	require.NoError(t, err)

	expected := Move{
		Id:           "Red:001",
		From:         &Position{R: 1, C: 2},
		To:           &Position{R: 1, C: 3},
		Result:       MOVERESULT_BothCaptured,
		AttackerRank: 10,
		AttackeeRank: 10,
	}

	require.Equal(t, expected, move)
}

func TestSerialMoveNoPositions(t *testing.T) {
	move := Move{
		Id:     "Blue:003",
		From:   nil,
		To:     &Position{R: 10, C: 3},
		Result: MOVERESULT_NoContest,
	}

	serialized := SerializeMove(move)

	require.Equal(t, "Blue:003 - 10,3 0 0 0", serialized)
}

func TestRoundTrip(t *testing.T) {
	move := Move{
		Id:           "Blue:003",
		From:         nil,
		To:           &Position{R: 10, C: 3},
		Result:       MOVERESULT_AttackerCaptured,
		AttackerRank: 5,
		AttackeeRank: 10,
	}

	serialized := SerializeMove(move)
	deserialized, err := DeserializeMove(serialized)

	require.NoError(t, err)
	require.Equal(t, move, deserialized)
}

func TestDeserialMoveNoPositions(t *testing.T) {
	serialized := "Red:001 - - 0 0 0"

	move, err := DeserializeMove(serialized)
	require.NoError(t, err)

	expected := Move{
		Id:     "Red:001",
		From:   nil,
		To:     nil,
		Result: MOVERESULT_NoContest,
	}

	require.Equal(t, expected, move)
}
