package pieces

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateStandardPieces(t *testing.T) {
	p := GenerateStandardPieces(COLOR_red)

	assert.Len(t, p, 40)
}
