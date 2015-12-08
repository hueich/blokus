package blokus

// Colors
const (
	blank = iota
	blue = iota
	yellow = iota
	red = iota
	green = iota
)

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

type Coord struct {
	X, Y int
}

type piece struct {
	// TODO: Add ID/name if we need to identify equivalent pieces in different orientations.
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
