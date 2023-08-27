package util

import (
	"testing"

	"github.com/cBiscuitSurprise/strate-go/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestGenerateStandardPieces(t *testing.T) {
	player := &core.Player{Number: 1}
	p := generateStandardPiecesForPlayer(player)

	assert.Len(t, p, 40)
}
