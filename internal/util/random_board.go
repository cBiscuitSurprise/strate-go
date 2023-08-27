package util

import (
	"fmt"
	"math/rand"

	"github.com/cBiscuitSurprise/strate-go/internal/core"
	"github.com/cBiscuitSurprise/strate-go/internal/game"
	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
	"github.com/rs/zerolog/log"
)

func CreateRandomlyPlannedBoard() (*game.Board, error) {
	board := createStandardBaseBoard()

	playerOne := &core.Player{Number: 1}
	playerTwo := &core.Player{Number: 2}

	playerOnePieces := shufflePieces(generateStandardPiecesForPlayer(playerOne))
	playerTwoPieces := shufflePieces(generateStandardPiecesForPlayer(playerTwo))

	board, err := placePieces(board, playerOnePieces, 0)
	if err != nil {
		return nil, err
	}

	board, err = placePieces(board, playerTwoPieces, 6)
	if err != nil {
		return nil, err
	}

	return board, nil
}

func shufflePieces(input []*pieces.Piece) []*pieces.Piece {
	output := input

	rand.Shuffle(len(output), func(i, j int) { output[i], output[j] = output[j], output[i] })

	return output
}

func placePieces(board *game.Board, unplaced []*pieces.Piece, startRow uint8) (*game.Board, error) {
	numColumns := int(board.GetSize().Columns)
	for inx := 0; inx < len(unplaced); inx++ {
		row := uint8(inx / numColumns)
		column := uint8(inx - int(row)*numColumns)
		row += startRow
		err := board.PlacePiece(unplaced[inx], game.Position{R: row, C: column})

		if err != nil {
			log.Error().
				Err(err).
				Int("number", inx).
				Uint8("row", row).
				Uint8("number", column).
				Str("piece", (*unplaced[inx]).Rank.String()).
				Str("position", fmt.Sprintf("(%d, %d)", row, column)).
				Msg("failed to place piece")
			return nil, err
		}
	}

	return board, nil
}
