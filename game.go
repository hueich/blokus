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
	pieceSet []*Piece
	// ID to use when generating the next new piece.
	nextPieceID int
	// Index of the player whose turn it is.
	curPlayerIndex int
}

func NewGame(id int, size int, pieces []*Piece) (*Game, error) {
	if size <= 0 || size > maxBoardSize {
		return nil, fmt.Errorf("Board size must be between 1 and %v. Provided: %v", maxBoardSize, size)
	}
	return &Game{
		id:          id,
		board:       NewBoard(size),
		pieceSet:    pieces,
		nextPieceID: 1,
	}, nil
}

func (g *Game) genPieceID() int {
	// TODO: Implement locking or use database to keep track of IDs
	id := g.nextPieceID
	g.nextPieceID += 1
	return id
}

func (g *Game) AddPlayer(name string, color int, corner Coord) error {
	// TODO: Check valid color
	// TODO: Check valid corner
	for _, p := range g.players {
		if p.name == name {
			return fmt.Errorf("Player %v already in the game", name)
		}
		if p.color == color {
			return fmt.Errorf("Color %v already taken by player %v", ColorName(color), p.name)
		}
		if p.corner == corner {
			return fmt.Errorf("Corner aleady occupied by player %v", p.name)
		}
	}
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
	if p.player != g.players[g.curPlayerIndex] {
		return fmt.Errorf("It's player %v's turn, but piece belongs to player %v", g.players[g.curPlayerIndex].name, p.player.name)
	}
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
	if err := g.advanceTurn(); err != nil {
		return err
	}
	return nil
}

func (g *Game) advanceTurn() error {
	if len(g.players) == 0 {
		return fmt.Errorf("Cannot advance turn with no players")
	}
	g.curPlayerIndex = (g.curPlayerIndex + 1) % len(g.players)
	return nil
}
