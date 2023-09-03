package game

import (
	"fmt"

	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
	"github.com/rs/zerolog/log"
)

func CreateRandomlyPlannedBoard() (*Board, error) {
	board := CreateStandardBaseBoard()

	playerOnePieces := pieces.GenerateStandardPieces(pieces.COLOR_red)
	playerTwoPieces := pieces.GenerateStandardPieces(pieces.COLOR_blue)

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

func placePieces(board *Board, unplaced map[string]*pieces.Piece, startRow uint8) (*Board, error) {
	numColumns := int(board.GetSize().Columns)
	inx := 0
	for _, p := range unplaced {
		row := uint8(inx / numColumns)
		column := uint8(inx - int(row)*numColumns)
		row += startRow
		err := board.PlacePiece(p, Position{R: row, C: column})

		if err != nil {
			log.Error().
				Err(err).
				Int("number", inx).
				Uint8("row", row).
				Uint8("number", column).
				Str("piece", p.GetRank().String()).
				Str("position", fmt.Sprintf("(%d, %d)", row, column)).
				Msg("failed to place piece")
			return nil, err
		}
		inx++
	}

	return board, nil
}
