package blokus

// Colors
const (
	blank = iota
	blue = iota
	yellow = iota
	red = iota
	green = iota
)

const (
	boardSize = 20
)

type board struct {
	grid [][]space
}

func NewBoard(size int) *board {
	b := board{
		grid: make([][]space, size),
	}
	for i := range b.grid {
		b.grid[i] = make([]space, size)
	}
	return &b
}

type coord struct {
	X, Y int
}

type piece struct {
	color int
	blocks []coord
}

func (p *piece) Rotate() {
	// TODO: Rotate coordinates clockwise
}

func (p *piece) Flip() {
	// TODO: Flip coordinates
}

type space struct {
	// empty bool
	color int
	parent *piece  // TODO: Do we need this?
}

func NewSpace(parent *piece) space {
	return space{
		color: parent.color,
		parent: parent,
	}
}

func (s *space) IsEmpty() bool {
	return s.color == blank
}

func (s *space) GetColor() int {
	return s.color
}