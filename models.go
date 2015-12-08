package blokus

// Colors
const (
	UnknownColor = iota
	Blue = iota
	Yellow = iota
	Red = iota
	Green = iota
)

// Coord represents a 2D coordinate.
type Coord struct {
	X, Y int
}

// Board represents the game board.
type Board struct {
	grid [][]space
}

func NewBoard(size int) *Board {
	b := Board{
		grid: make([][]space, size),
	}
	for i := range b.grid {
		b.grid[i] = make([]space, size)
	}
	return &b
}

// Piece represents a puzzle piece, made up of one or more square blocks.
type piece struct {
	id int
	color int
	blocks []Coord
	corners []Coord
}

func (p *piece) Rotate() {
	// TODO: Rotate coordinates clockwise.
}

func (p *piece) Flip() {
	// TODO: Flip coordinates.
}

// Space represents a space on the board.
type space struct {
	// TODO: Evaluate if we can just use *piece directly in the Board.
	parent *piece
}

func NewSpace(parent *piece) space {
	return space{
		parent: parent,
	}
}

func (s *space) IsEmpty() bool {
	return s.parent == nil
}

func (s *space) GetColor() int {
	// TODO: Gracefully handle nil pointer.
	return s.parent.color
}
