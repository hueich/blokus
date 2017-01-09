package blokus

import (
	"fmt"
)

// GameID is the ID of a created game.
type GameID int64

// Coord represents a 2D coordinate, where X increases downward and Y increases rightward.
type Coord struct {
	X, Y int
}

// Color is an enum of available colors.
type Color uint8

const (
	Blue Color = iota + 1
	Yellow
	Red
	Green

	colorEnd
)

func (c Color) IsValid() bool {
	return c > 0 && c < colorEnd
}

func (c Color) String() string {
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
	return ""
}

type Player struct {
	// Unique name of the player.
	name   string
	color  Color
	// The position the player starts from, e.g. [0,0], or [0,19] for a size 20 board.
	startPos Coord
	// True at an index means the corresponding piece has been placed on the board.
	placedPieces []bool
}

func (p *Player) Color() Color {
	return p.color
}

func (p *Player) RemovePiece(index int) (*Piece, error) {
	if index < 0 || index >= len(p.pieces) {
		return nil, fmt.Errorf("Piece index out of range: %v", index)
	}
	piece := p.pieces[index]
	if piece == nil {
		return nil, fmt.Errorf("Piece at index %v is already placed", index)
	}
	p.pieces[index] = nil
	return piece, nil
}

// Board represents the game board.
type Board struct {
	grid [][]Color
}

func NewBoard(size int) *Board {
	b := Board{
		grid: make([][]Color, size),
	}
	for i := range b.grid {
		b.grid[i] = make([]Color, size)
	}
	return &b
}

func (b *Board) isOutOfBounds(c Coord) bool {
	return c.X < 0 || c.Y < 0 || c.X >= len(b.grid) || c.Y >= len(b.grid[0])
}

// Piece represents a puzzle piece, made up of one or more square blocks.
type Piece struct {
	// The square blocks this piece consists of. First block must be at (0,0) with other blocks relative to it.
	// The blocks are stored in their original coordinates with no rotation or flipping. Orientation is used to calcuate the actual coordinates.
	blocks []Coord
	// The corner squares of this piece, which was calculated from blocks and cached here.
	corners []Coord
	// The squares that border blocks of this piece.
	sides []Coord
}

func NewPiece(blocks []Coord) (*Piece, error) {
	p := &Piece{
		blocks:  make([]Coord, len(blocks)),
		corners: getCorners(blocks),
	}
	if copied := copy(p.blocks, blocks); copied != len(blocks) {
		return nil, fmt.Errorf("Copied %d blocks, but should've copied %d blocks", copied, len(blocks))
	}
	return p, nil
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

type Move struct {
	// The player who made the move. Cannot be nil.
	player *Player
	// The index of the piece that was played. Negative if the turn was passed.
	pieceIndex int
	// Orientation of the piece when played.
	orient Orientation
	// Location on the board where the piece was played.
	// This is the coordinate where the (0,0) block of the piece is located.
	loc Coord
}
