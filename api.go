package blokus

import (
	"context"
)

type GameID int64

type GameOptions struct {
	// Name of the set of starting pieces to use for the game. Leave empty for default set.
	PieceSetName string
}

type PlayerOptions struct {
	Color  Color
	Corner Corner
}

type Service interface {

	// CreateGame returns the ID of the created game.
	// Username is the user who owns the created game.
	// BoardSize is the length of one edge of the board for a square board.
	CreateGame(ctx context.Context, username, gamename string, boardSize int, opts *GameOptions) (GameID, error)

	// AddPlayer adds a new player to the game. It's an error to add an existing player or if the color or corner is taken.
	AddPlayer(ctx context.Context, id GameID, username string, opts *PlayerOptions) error

	// StartGame starts the game with the specified player. It's an error to call this if the game already started.
	StartGame(ctx context.Context, id GameID, username string) error

	// PlacePiece places a piece on the board, given the position and orientation.
	// Coordinate [0,0] is the top left corner of the board. X increases downward and Y increases rightward.
	// Rotation starts from 0 and increases clockwise in 90-degree increments. Flip, if true, flips the piece horizontally, i.e. around the X-axis.
	PlacePiece(ctx context.Context, id GameID, pieceID int, x, y, rot int, flip bool) error

	// GetGameState gets the current state of the game, from which a client may construct a view of the game.
	// TODO: Figure out what to return to concisely represent game state
	GetGameState(ctx context.Context, id GameID) error
}
