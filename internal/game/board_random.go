package game

import (
	"fmt"

	"github.com/cBiscuitSurprise/strate-go/internal/pieces"
	"github.com/cBiscuitSurprise/strate-go/internal/util"
	"github.com/rs/zerolog/log"
)

func CreateRandomlyPlannedBoard() (*Board, error) {
	playerOnePieces := pieces.GenerateStandardPieces(pieces.COLOR_red)
	playerTwoPieces := pieces.GenerateStandardPieces(pieces.COLOR_blue)

	allPieces := util.UpdateMap[string, *pieces.Piece](map[string]*pieces.Piece{}, playerOnePieces, playerTwoPieces)
	board := CreateStandardBaseBoard(allPieces)

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

func placePieces(board *Board, unplaced map[string]*pieces.Piece, startRow int) (*Board, error) {
	numColumns := int(board.GetSize().Columns)
	inx := 0
	for _, p := range unplaced {
		row := int(inx / numColumns)
		column := int(inx - int(row)*numColumns)
		row += startRow
		err := board.PlacePiece(p.GetId(), Position{R: row, C: column})

		if err != nil {
			log.Error().
				Err(err).
				Str("method", "board_random.placePieces").
				Int("number", inx).
				Int("row", row).
				Int("number", column).
				Str("piece", p.GetId()).
				Str("rank", p.GetRank().String()).
				Str("position", fmt.Sprintf("(%d, %d)", row, column)).
				Msg("failed to generate random board")
			return nil, err
		}
		inx++
	}

	return board, nil
}
