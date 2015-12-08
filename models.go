package blokus

// Colors
const (
	UnknownColor = iota
	Blue = iota
	Yellow = iota
	Red = iota
	Green = iota
)

type Player struct {
	name string
	color int
	pieces []*Piece
	// The corner the player starts from, with X and Y coordinates being either -1 or 1.
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
type Piece struct {
	id int
	player *Player
	blocks []Coord
	corners []Coord
}

func (p *Piece) GetColor() int {
	return p.player.GetColor()
}

func (p *Piece) Rotate() {
	// TODO: Rotate coordinates clockwise.
}

func (p *Piece) Flip() {
	// TODO: Flip coordinates.
}

// Space represents a space on the board.
type space struct {
	// TODO: Evaluate if we can just use *Piece directly in the Board.
	piece *Piece
}

func NewSpace(piece *Piece) space {
	return space{
		piece: piece,
	}
}

func (s *space) IsEmpty() bool {
	return s.piece == nil
}

func (s *space) GetColor() int {
	// TODO: Gracefully handle nil pointer.
	return s.piece.GetColor()
}
