package blokus

import (
	"google.golang.org/cloud/datastore"
)

// Colors
const (
	_ = iota
	Blue
	Yellow
	Red
	Green
)

func ColorName(c int) string {
	switch c {
	case Blue:
		return "blue"
	case Yellow:
		return "yellow"
	case Red:
		return "red"
	case Green:
		return "green"
	}
	return "unknown"
}

type Player struct {
	// Unique name of the player.
	name   string
	color  int
	// pieces []*datastore.Key  // []*Piece
	// The corner the player starts from, e.g. [-1,-1], or [-1,20] for a size 20 board.
	corner Coord
}

func (p *Player) GetColor() int {
	return p.color
}

// Coord represents a 2D coordinate.
type Coord struct {
	X, Y int
}

// Board represents the game board.
type Board struct {
	size int
	grid []*datastore.Key  // []*Piece
}

func NewBoard(size int) *Board {
	b := Board{
		size: size,
		grid: make([]*Piece, size*size),
	}
	return &b
}

func (b *Board) GetCoord(x, y int) *Piece {
	return b.grid[x*b.size+y]
}

func (b *Board) SetCoord(x, y int, p *Piece) {
	b.grid[x*b.size+y] = p
}

type PieceTemplate struct {
	blocks []Coord
}

// Piece represents a puzzle piece, made up of one or more square blocks.
type Piece struct {
	id int
	// The player who owns this piece
	player *datastore.Key  // *Player
	// The coordinate this piece was placed in, or nil if it's not placed yet.
	// This is the coordinate where the (0,0) block is located.
	location *Coord
	// The square blocks this piece consists of. First block must be at (0,0) with other blocks relative to it.
	blocks []Coord
	// The corner squares of this piece, which was calculated from blocks and cached here.
	corners []Coord
	// Number of 90 degree clockwise rotations, between 0-3, where 0 is no rotation, i.e. original orientation.
	rot int
	// True if the piece is flipped.
	flip bool
}

func NewPiece(id int, player *Player, t *PieceTemplate) *Piece {
	p := &Piece{
		id:      id,
		player:  player,
		blocks:  make([]Coord, len(t.blocks)),
		corners: getCorners(t.blocks),
	}
	copy(p.blocks, t.blocks)
	return p
}

func getCorners(blocks []Coord) []Coord {
	corners := map[Coord]bool{}
	// Add corners of all blocks
	for _, block := range blocks {
		corners[Coord{block.X - 1, block.Y - 1}] = true
		corners[Coord{block.X - 1, block.Y + 1}] = true
		corners[Coord{block.X + 1, block.Y - 1}] = true
		corners[Coord{block.X + 1, block.Y + 1}] = true
	}
	// Remove corners that touch a block
	for _, block := range blocks {
		delete(corners, Coord{block.X, block.Y})
		delete(corners, Coord{block.X + 1, block.Y})
		delete(corners, Coord{block.X - 1, block.Y})
		delete(corners, Coord{block.X, block.Y + 1})
		delete(corners, Coord{block.X, block.Y - 1})
	}
	c := []Coord{}
	for coord := range corners {
		c = append(c, coord)
	}
	return c
}

func (p *Piece) GetColor() int {
	return p.player.GetColor()
}

// Rotate piece clockwise 90 degrees
func (p *Piece) Rotate() {
	for i, c := range p.blocks {
		p.blocks[i] = rotateCoord(c)
	}
	for i, c := range p.corners {
		p.corners[i] = rotateCoord(c)
	}
	p.rot = (p.rot + 1) % 4
}

func rotateCoord(c Coord) Coord {
	return Coord{c.Y, -c.X}
}

// Flip piece along X axis
func (p *Piece) Flip() {
	for i, c := range p.blocks {
		p.blocks[i] = flipCoord(c)
	}
	for i, c := range p.corners {
		p.corners[i] = flipCoord(c)
	}
	p.flip = !p.flip
}

func flipCoord(c Coord) Coord {
	return Coord{c.X, -c.Y}
}
