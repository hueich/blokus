package blokus

import (
	"golang.org/x/net/context"
)

type GameOptions struct {
	// ID of the starting set of pieces to use for the game. 0 is the default set.
	PieceSet int
}

// CreateGame returns the ID of the created game.
func CreateGame(ctx context.Context, username, gamename string, size int, opt GameOptions) (int64, error) {
	return 123, nil
}

func AddPlayer(ctx context.Context, gameID int64, username string, color, corner int) error {
	return nil
}

func StartGame(ctx context.Context, gameID int64) error {
	return nil
}

func PlacePiece(ctx context.Context, gameID, pieceID int64, x, y, rot int, flip bool) error {
	return nil
}

// TODO: Figure out what to return to concisely represent game state
func GetGameState(ctx context.Context, gameID, lastPieceID int64) error {
	return nil
}
