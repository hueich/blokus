package blokus

import (
	"fmt"
)

const (
	// Default game board size.
	DefaultBoardSize = 20
	maxBoardSize     = 100
)

var (
	neighbors = [4]Coord{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	}
)

type Game struct {
	id      GameID
	players []*Player
	board   *Board
	// Set of pieces every player starts with.
	pieceSet []*Piece
	// ID to use when generating the next new piece.
	nextPieceID int
	// Index of the player whose turn it is.
	curPlayerIndex int
	// Moves that have been played.
	moves []Move
}

func NewGame(id GameID, size int, pieces []*Piece) (*Game, error) {
	if size <= 0 || size > maxBoardSize {
		return nil, fmt.Errorf("Board size must be between 1 and %v. Provided: %v", maxBoardSize, size)
	}
	return &Game{
		id:          id,
		board:       NewBoard(size),
		pieceSet:    pieces,
		nextPieceID: 1,
	}, nil
}

func (g *Game) genPieceID() int {
	// TODO: Implement locking or use database to keep track of IDs
	id := g.nextPieceID
	g.nextPieceID++
	return id
}

func (g *Game) AddPlayer(name string, color Color, startPos Coord) error {
	if !color.IsValid() {
		return fmt.Errorf("Unknown color value: %v", color)
	}
	if g.board.isOutOfBounds(startPos) {
		return fmt.Errorf("Starting position is out of bounds: %v", startPos)
	}
	for _, p := range g.players {
		if p.name == name {
			return fmt.Errorf("Player %v already in the game", name)
		}
		if p.color == color {
			return fmt.Errorf("Color %v already taken by player %v", color, p.name)
		}
		if p.startPos == startPos {
			return fmt.Errorf("Starting position already occupied by player %v", p.name)
		}
	}
	p := &Player{
		name:     name,
		color:    color,
		startPos: startPos,
	}
	// Make a copy of the set of starting pieces of the game.
	for _, ps := range g.pieceSet {
		b := make([]Coord, len(ps.blocks))
		copy(b, ps.blocks)
		p.pieces = append(p.pieces, NewPiece(g.genPieceID(), p, b))
	}
	g.players = append(g.players, p)
	return nil
}

// Place the piece on the board and record the move, unless there's an error.
// This does not check for winner nor advance player turn.
func (g *Game) PlacePiece(loc Coord, pieceID int, orient Orientation) error {
	// Preliminary input validation.
	if g.board.isOutOfBounds(loc) {
		return fmt.Errorf("Piece placement out of bounds: %v,%v", loc.X, loc.Y)
	}

	// TODO: Integrate with database for more efficient piece lookup
	var p *Piece
	for _, player := range g.players {
		for _, piece := range player.pieces {
			if piece.id == pieceID {
				p = piece
				break
			}
		}
	}
	if p == nil {
		return fmt.Errorf("Could not find piece with ID: %v", pieceID)
	}

	// Check if piece is already placed.
	if p.location != nil {
		return fmt.Errorf("Piece is already placed at [%v]", p.location)
	}

	// Check if it's this player's turn.
	if p.player != g.players[g.curPlayerIndex] {
		return fmt.Errorf("It's player %v's turn, but piece belongs to player %v", g.players[g.curPlayerIndex].name, p.player.name)
	}

	// Rotate/flip piece to specified orientation.
	orient.Rot = orient.Rot.Normalize()
	for p.rot != int(orient.Rot) {
		p.Rotate()
	}
	if p.flip != orient.Flip {
		p.Flip()
	}

	if err := g.checkPiecePlacement(loc, p); err != nil {
		return err
	}

	// Actually place the piece.
	p.location = &Coord{loc.X, loc.Y}
	for _, b := range p.blocks {
		g.board.grid[loc.X+b.X][loc.Y+b.Y] = p.Color()
	}

	// Record the move.
	g.moves = append(g.moves, Move{
		player: p.player,
		piece:  p,
		orient: orient,
		loc:    loc,
	})

	return nil
}

// Checks whether piece placement is valid. Returns error if invalid.
func (g *Game) checkPiecePlacement(loc Coord, p *Piece) error {
	return g.checkPiecePlacementAt(loc, p, 0)
}

func (g *Game) checkPiecePlacementAt(loc Coord, p *Piece, block int) error {
	if p == nil {
		return fmt.Errorf("Piece cannot be nil")
	}
	if p.player == nil {
		return fmt.Errorf("Piece has no owning player")
	}
	if block < 0 || block >= len(p.blocks) {
		return fmt.Errorf("Specified block index %v for piece of %v blocks is out of bounds", block, len(p.blocks))
	}

	// Offset the placement location with block's offset.
	loc = Coord{loc.X - p.blocks[block].X, loc.Y - p.blocks[block].Y}

	isStartPos := false
	for _, b := range p.blocks {
		// Change from relative to absolute coordinate.
		b = Coord{b.X + loc.X, b.Y + loc.Y}
		// Check that every block is inside the board
		if g.board.isOutOfBounds(b) {
			return fmt.Errorf("Piece placement out of bounds")
		}
		// Check that every block is on an empty space
		if g.board.grid[b.X][b.Y].IsValid() {
			return fmt.Errorf("Cell %v,%v is occupied", b.X, b.Y)
		}
		// Check that every block is not next to a piece of same color
		for _, n := range neighbors {
			n = Coord{b.X + n.X, b.Y + n.Y}
			if g.board.isOutOfBounds(n) {
				continue
			}
			if s := g.board.grid[n.X][n.Y]; s.IsValid() && s == p.player.color {
				return fmt.Errorf("Piece is next to another %v piece", p.player.color)
			}
		}
		// Check if this is the player's starting position.
		if b == p.player.startPos {
			isStartPos = true
		}
	}
	if isStartPos {
		return nil
	}

	validCorner := false
	for _, c := range p.corners {
		// Change from relative to absolute coordinate.
		c = Coord{c.X + loc.X, c.Y + loc.Y}

		if g.board.isOutOfBounds(c) {
			continue
		}
		// Check that at least one corner is touching a block of same color.
		if s := g.board.grid[c.X][c.Y]; s.IsValid() && s == p.player.color {
			validCorner = true
			break
		}
	}
	if !validCorner {
		return fmt.Errorf("Piece has no corner touching another %v piece, and it doesn't cover the player's starting position %v", p.player.color, p.player.startPos)
	}
	return nil
}

// Advances the game turn to the next player.
func (g *Game) AdvanceTurn() error {
	if len(g.players) == 0 {
		return fmt.Errorf("Cannot advance turn with no players")
	}
	g.curPlayerIndex = (g.curPlayerIndex + 1) % len(g.players)
	return nil
}
