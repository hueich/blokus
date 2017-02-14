package blokus

import (
	"fmt"
)

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
	Height, Width int
	Grid          []Color
}

func NewBoard(size int) (*Board, error) {
	return NewRectBoard(size, size)
}

func NewRectBoard(height, width int) (*Board, error) {
	if height <= 0 {
		return nil, fmt.Errorf("Board height must be positive, gave %v", height)
	}
	if width <= 0 {
		return nil, fmt.Errorf("Board width must be positive, gave %v", width)
	}
	b := &Board{
		Height: height,
		Width:  width,
		Grid:   make([]Color, height*width),
	}
	return b, nil
}

func (b *Board) Cell(c Coord) Color {
	if b.IsOutOfBounds(c) {
		return colorEmpty
	}
	return b.Grid[c.X*b.Width+c.Y]
}

func (b *Board) SetCell(coord Coord, color Color) {
	if b.IsOutOfBounds(coord) {
		return
	}
	b.Grid[coord.X*b.Width+coord.Y] = color
}

func (b *Board) IsOutOfBounds(c Coord) bool {
	return c.X < 0 || c.Y < 0 || c.X >= b.Height || c.Y >= b.Width
}

// Piece represents a puzzle piece, made up of one or more square blocks.
type Piece struct {
	// The square blocks this piece consists of. First block must be at (0,0) with other blocks relative to it.
	// The blocks are stored in their original coordinates with no rotation or flipping. Orientation is used to calcuate the actual coordinates.
	Blocks []Coord
	// The corner squares of this piece, which was calculated from blocks and cached here.
	corners []Coord
}

func NewPiece(blocks []Coord) (*Piece, error) {
	if len(blocks) == 0 {
		return nil, fmt.Errorf("Cannot create a piece with no blocks")
	}
	p := &Piece{
		// Make a copy, in case the same block slice is used to make other pieces.
		Blocks: append([]Coord(nil), blocks...),
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

func (p *Piece) Corners() []Coord {
	if len(p.corners) == 0 {
		p.corners = getCorners(p.Blocks)
	}
	return p.corners
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
	Player *Player
	// The index of the piece that was played. Negative if the turn was passed.
	PieceIndex int
	// Orientation of the piece when played.
	Orient Orientation
	// Location on the board where the piece was played.
	// This is the coordinate where the (0,0) block of the piece is located.
	Loc Coord
}

func (m Move) IsPass() bool {
	return m.PieceIndex < 0
}
