package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttackInvert(t *testing.T) {
	testCases := []struct {
		original Winner
		inverted Winner
	}{
		{WINNER_Attackee, WINNER_Attacker},
		{WINNER_Attacker, WINNER_Attackee},
		{WINNER_Draw, WINNER_Draw},
		{100, 100},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.original.Invert(), tc.inverted)
	}
}
