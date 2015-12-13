package blokus

import (
	"fmt"
)

const (
	// Default game board size.
	DefaultBoardSize = 20
	maxBoardSize     = 100
)

type Game struct {
	id      int
	players []*Player
	board   *Board
	// Set of pieces every player starts with.
	pieceSet   []*Piece
	curPieceID int
}

func NewGame(id int, size int, pieces []*Piece) (*Game, error) {
	if size <= 0 || size > maxBoardSize {
		return nil, fmt.Errorf("Board size must be between 1 and %v. Provided: %v", maxBoardSize, size)
	}
	return &Game{
		id:         id,
		board:      NewBoard(size),
		pieceSet:   pieces,
		curPieceID: 1,
	}, nil
}

func (g *Game) genPieceID() int {
	// TODO: Implement locking or use database to keep track of IDs
	id := g.curPieceID
	g.curPieceID += 1
	return id
}

func (g *Game) AddPlayer(name string, color int, corner Coord) error {
	// TODO: Check corner
	p := &Player{
		name:   name,
		color:  color,
		corner: corner,
	}
	// Make a copy of the set of starting pieces of the game.
	for _, ps := range g.pieceSet {
		b := make([]Coord, len(ps.blocks))
		copy(b, ps.blocks)
		p.pieces = append(p.pieces, NewPiece(g.genPieceID(), p, b))
	}
	g.players = append(g.players, p)
	return nil
}
