package blokus

import (
	"context"
)

type GameOptions struct {
	// Name of the set of starting pieces to use for the game. Leave empty for default set.
	PieceSetName string
}

type PlayerOptions struct {
	Color  Color
	Corner Corner
}

// CreateGame returns the ID of the created game.
// Username is the user who owns the created game.
// BoardSize is the length of one edge of the board for a square board.
func CreateGame(ctx context.Context, username, gamename string, boardSize int, opt *GameOptions) (int64, error) {
	return 123, nil
}

// AddPlayer adds a new player to the game. It's an error to add an existing player or if the color or corner is taken.
func AddPlayer(ctx context.Context, gameID int64, username string, opt *PlayerOptions) error {
	return nil
}

// StartGame starts the game with the specified player. It's an error to call this if the game already started.
func StartGame(ctx context.Context, gameID int64, username string) error {
	return nil
}

// PlacePiece places a piece on the board, given the position and orientation.
// Coordinate [0,0] is the top left corner of the board. X increases downward and Y increases rightward.
// Rotation starts from 0 and increases clockwise in 90-degree increments. Flip, if true, flips the piece horizontally, i.e. around the X-axis.
func PlacePiece(ctx context.Context, gameID, pieceID int64, x, y, rot int, flip bool) error {
	return nil
}

// GetGameState gets the current state of the game, from which a client may construct a view of the game.
// TODO: Figure out what to return to concisely represent game state
func GetGameState(ctx context.Context, gameID int64) error {
	return nil
}
