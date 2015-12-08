package blokus

// Colors
const (
	UnknownColor = iota
	Blue = iota
	Yellow = iota
	Red = iota
	Green = iota
)

type Coord struct {
	X, Y int
}

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
