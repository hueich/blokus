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
	// TODO: Check color is valid and not taken
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

func (g *Game) PlacePiece(loc Coord, pieceID int, rot int, flip bool) error {
	// Preliminary input validation.
	if loc.X < 0 || loc.Y < 0 || loc.X >= len(g.board.grid) || loc.Y >= len(g.board.grid[0]) {
		return fmt.Errorf("Piece location out of bounds: %v,%v", loc.X, loc.Y)
	}
	if rot < 0 || rot >= 4 {
		return fmt.Errorf("Invalid piece rotation: %v", rot)
	}
	// TODO: Integrate with database for more efficient piece lookup
	var p *Piece
	for _, player := range g.players {
		for _, piece := range player.pieces {
			if piece.id == pieceID {
				p = piece
				break
			}
		}
	}
	if p == nil {
		return fmt.Errorf("Could not find piece with ID: %v", pieceID)
	}
	// TODO: Check if correct player turn
	// Rotate/flip to specified position.
	for p.rot != rot {
		p.Rotate()
	}
	if p.flip != flip {
		p.Flip()
	}
	// TODO: Check if valid position
	p.location = &loc
	g.board.grid[loc.X][loc.Y] = p
	return nil
}