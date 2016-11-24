package blokus

type Color uint8

const (
	Blue Color = iota + 1
	Yellow
	Red
	Green
)

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

type Corner uint8

const (
	TopLeft Corner = iota + 1
	TopRight
	BottomRight
	BottomLeft
)

type Player struct {
	// Unique name of the player.
	name   string
	color  Color
	pieces []*Piece
	// The corner the player starts from, e.g. [-1,-1], or [-1,20] for a size 20 board.
	corner Coord
}

func (p *Player) GetColor() Color {
	return p.color
}

// Coord represents a 2D coordinate, where X increases downward and Y increases rightward.
type Coord struct {
	X, Y int
}

// Board represents the game board.
type Board struct {
	grid [][]*Piece
}

func NewBoard(size int) *Board {
	b := Board{
		grid: make([][]*Piece, size),
	}
	for i := range b.grid {
		b.grid[i] = make([]*Piece, size)
	}
	return &b
}

// Piece represents a puzzle piece, made up of one or more square blocks.
type Piece struct {
	id int
	// The player who owns this piece
	player *Player
	// The coordinate this piece was placed in, or nil if it's not placed yet.
	// This is the coordinate where the (0,0) block is located.
	location *Coord
	// The square blocks this piece consists of. First block must be at (0,0) with other blocks relative to it.
	blocks []Coord
	// The corner squares of this piece, which was calculated from blocks and cached here.
	corners []Coord
	// Number of 90 degree clockwise rotations, between 0-3, where 0 is no rotation, i.e. original orientation.
	rot int
	// True if the piece is flipped horizontally, i.e. around the X-axis.
	flip bool
}

func NewPiece(id int, player *Player, blocks []Coord) *Piece {
	return &Piece{
		id:      id,
		player:  player,
		blocks:  blocks,
		corners: getCorners(blocks),
	}
}

func NewTemplatePiece(blocks []Coord) *Piece {
	return NewPiece(0, nil, blocks)
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

func (p *Piece) GetColor() Color {
	return p.player.GetColor()
}

// Rotate piece clockwise 90 degrees.
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

// Flip piece horizontally, around the X-axis.
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
