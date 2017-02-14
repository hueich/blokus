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

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

// Color is an enum of available colors.
type Color uint8

const (
	colorEmpty Color = iota

	Blue
	Yellow
	Red
	Green

	colorEnd
)

func (c Color) IsColored() bool {
	return c > 0 && c < colorEnd
}

func (c Color) String() string {
	switch c {
	case colorEmpty:
		return "empty"

	case Blue:
		return "blue"
	case Yellow:
		return "yellow"
	case Red:
		return "red"
	case Green:
		return "green"
	}
	return "unknown color"
}

type Player struct {
	// Unique name of the player.
	Name  string
	Color Color
	// The position the player starts from, e.g. [0,0], or [0,19] for a size 20 board.
	StartPos Coord
	// True at an index means the corresponding piece has been placed on the board.
	PlacedPieces []bool
}

func NewPlayer(name string, color Color, startPos Coord, numPieces int) (*Player, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("Player name cannot be empty")
	}
	if !color.IsColored() {
		return nil, fmt.Errorf("Invalid player color: %v", color)
	}
	// Start position can be arbitrary, so should check at the Game level.
	if numPieces <= 0 {
		return nil, fmt.Errorf("Number of pieces for the player must be positive: %v", numPieces)
	}
	return &Player{
		Name:         name,
		Color:        color,
		StartPos:     startPos,
		PlacedPieces: make([]bool, numPieces),
	}, nil
}

func (p *Player) CheckPiecePlaceability(index int) error {
	if index < 0 || index >= len(p.PlacedPieces) {
		return fmt.Errorf("Piece index out of range: %v", index)
	}
	if p.PlacedPieces[index] {
		return fmt.Errorf("Piece at index %d is already placed", index)
	}
	return nil
}

func (p *Player) placePiece(index int) error {
	if err := p.CheckPiecePlaceability(index); err != nil {
		return err
	}
	p.PlacedPieces[index] = true
	return nil
}

// Board represents the game board.
type Board struct {
	grid [][]Color
}

func NewBoard(size int) *Board {
	b := &Board{
		grid: make([][]Color, size),
	}
	for i := range b.grid {
		b.grid[i] = make([]Color, size)
	}
	return b
}

func (b *Board) Grid() [][]Color {
	return b.grid
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
}

func NewPiece(blocks []Coord) (*Piece, error) {
	if len(blocks) == 0 {
		return nil, fmt.Errorf("Cannot create a piece with no blocks")
	}
	p := &Piece{
		// Make a copy, in case the same block slice is used to make other pieces.
		blocks:  append([]Coord(nil), blocks...),
		corners: getCorners(blocks),
	}
	return p, nil
}

func NewPieceOrNil(blocks []Coord) *Piece {
	p, err := NewPiece(blocks)
	if err != nil {
		return nil
	}
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

func (m Move) isPass() bool {
	return m.pieceIndex < 0
}
