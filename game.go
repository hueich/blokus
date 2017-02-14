package blokus

import (
	"fmt"
)

const (
	maxBoardSize = 100
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
	Players []*Player
	Board   *Board
	// Set of pieces every player starts with.
	Pieces []*Piece
	// Index of the player whose turn it is.
	CurPlayerIndex int
	// Moves that have been played.
	Moves []*Move
}

func NewGame(size int, pieces []*Piece) (*Game, error) {
	if size <= 0 || size > maxBoardSize {
		return nil, fmt.Errorf("Board size must be between 1 and %v. Provided: %v", maxBoardSize, size)
	}
	if len(pieces) == 0 {
		return nil, fmt.Errorf("Cannot create game with no pieces")
	}
	return &Game{
		Board:  NewBoard(size),
		Pieces: pieces,
	}, nil
}

func (g *Game) CurrentPlayer() *Player {
	return g.Players[g.CurPlayerIndex]
}

func (g *Game) GetNextFreeColor() (Color, error) {
	allColors := make([]bool, int(colorEnd))
	for _, p := range g.Players {
		if !p.Color.IsColored() {
			return 0, fmt.Errorf("Player %v has invalid color: %v", p.Name, p.Color)
		}
		allColors[int(p.Color)] = true
	}
	for i := 1; i < int(colorEnd); i++ {
		if !allColors[i] {
			return Color(i), nil
		}
	}
	return 0, fmt.Errorf("No more free colors")
}

func (g *Game) AddPlayer(name string, color Color, startPos Coord) error {
	if len(name) == 0 {
		return fmt.Errorf("Player name cannot be empty")
	}
	if color == colorEmpty {
		var err error
		color, err = g.GetNextFreeColor()
		if err != nil {
			return err
		}
	}
	if !color.IsColored() {
		return fmt.Errorf("Invalid color %v", color)
	}
	if g.Board.IsOutOfBounds(startPos) {
		return fmt.Errorf("Starting position is out of bounds: %v", startPos)
	}

	for _, p := range g.Players {
		if p.Name == name {
			return fmt.Errorf("Player %v already in the game", name)
		}
		if p.Color == color {
			return fmt.Errorf("Color %v already taken by player %v", color, p.Name)
		}
		if p.StartPos == startPos {
			return fmt.Errorf("Starting position already occupied by player %v", p.Name)
		}
	}
	p, err := NewPlayer(name, color, startPos, len(g.Pieces))
	if err != nil {
		return fmt.Errorf("Error adding new player: %v", err)
	}
	g.Players = append(g.Players, p)
	return nil
}

func (g *Game) PassTurn(player *Player) error {
	if player == nil {
		return fmt.Errorf("Invalid player")
	}
	// Check if it's this player's turn.
	if player != g.Players[g.CurPlayerIndex] {
		return fmt.Errorf("Turn belongs to player %v, not player %v", g.Players[g.CurPlayerIndex].Name, player.Name)
	}
	// Record the move.
	g.Moves = append(g.Moves, &Move{
		Player:     player,
		PieceIndex: -1,
	})
	return nil
}

// Place the piece on the board and record the move, unless there's an error.
// This does not check for winner nor advance player turn.
func (g *Game) PlacePiece(player *Player, pieceIndex int, orient Orientation, loc Coord) error {
	if player == nil {
		return fmt.Errorf("Invalid player")
	}
	// Check if it's this player's turn.
	if player != g.Players[g.CurPlayerIndex] {
		return fmt.Errorf("Turn belongs to player %v, not player %v", g.Players[g.CurPlayerIndex].Name, player.Name)
	}

	if pieceIndex < 0 || pieceIndex >= len(g.Pieces) {
		return fmt.Errorf("Piece index is out of range: %d", pieceIndex)
	}
	if err := player.CheckPiecePlaceability(pieceIndex); err != nil {
		return err
	}

	piece := g.Pieces[pieceIndex]
	if piece == nil {
		return fmt.Errorf("Piece at index %d is inexplicably nil", pieceIndex)
	}
	orientedPiece := &Piece{
		Blocks:  orient.TransformCoords(piece.Blocks),
		corners: orient.TransformCoords(piece.corners),
	}

	if err := g.checkPiecePlacement(player, orientedPiece, loc); err != nil {
		return err
	}

	// Actually place the piece.
	if err := player.placePiece(pieceIndex); err != nil {
		return err
	}
	for _, b := range orientedPiece.Blocks {
		g.Board.Grid[loc.X+b.X][loc.Y+b.Y] = player.Color
	}

	// Record the move.
	g.Moves = append(g.Moves, &Move{
		Player:     player,
		PieceIndex: pieceIndex,
		Orient:     orient,
		Loc:        loc,
	})

	return nil
}

// Checks whether piece placement is valid. Returns error if invalid.
// The piece should already be oriented.
func (g *Game) checkPiecePlacement(player *Player, piece *Piece, loc Coord) error {
	coversStartPos := false
	for _, b := range piece.Blocks {
		// Change from relative to absolute coordinate.
		b = Coord{b.X + loc.X, b.Y + loc.Y}
		// Check that every block is inside the board
		if g.Board.IsOutOfBounds(b) {
			return fmt.Errorf("Piece placement out of bounds")
		}
		// Check that every block is on an empty space
		if g.Board.Grid[b.X][b.Y].IsColored() {
			return fmt.Errorf("Cell (%v,%v) is occupied by color %v", b.X, b.Y, g.Board.Grid[b.X][b.Y])
		}
		// Check that every block is not next to a piece of same color
		for _, n := range neighbors {
			n = Coord{b.X + n.X, b.Y + n.Y}
			if g.Board.IsOutOfBounds(n) {
				continue
			}
			if g.Board.Grid[n.X][n.Y] == player.Color {
				return fmt.Errorf("Piece is next to another %v piece", player.Color)
			}
		}
		// Check if this is the player's starting position.
		if b == player.StartPos {
			coversStartPos = true
		}
	}
	if coversStartPos {
		// Should never reach here for non-first moves, since the first move necessarily must cover the starting position,
		// so subsequent moves that cover the starting position will fail the earlier check for occupied cells.
		return nil
	}

	hasValidCorner := false
	for _, c := range piece.corners {
		// Change from relative to absolute coordinate.
		c = Coord{c.X + loc.X, c.Y + loc.Y}

		if g.Board.IsOutOfBounds(c) {
			continue
		}
		// Check that at least one corner is touching a block of same color.
		if g.Board.Grid[c.X][c.Y] == player.Color {
			hasValidCorner = true
			break
		}
	}
	if !hasValidCorner {
		// TODO: Make error message more specific.
		return fmt.Errorf("Piece has no corner touching another %v piece, and it doesn't cover the player's starting position %v", player.Color, player.StartPos)
	}
	return nil
}

// Advances the game turn to the next player.
func (g *Game) AdvanceTurn() error {
	if len(g.Players) == 0 {
		return fmt.Errorf("Cannot advance turn with no players")
	}
	g.CurPlayerIndex = (g.CurPlayerIndex + 1) % len(g.Players)
	return nil
}

// Game ends when all players passed for a round.
func (g *Game) IsGameEnd() bool {
	if len(g.Moves) == 0 || len(g.Moves) < len(g.Players) {
		return false
	}
	for _, m := range g.Moves[len(g.Moves)-len(g.Players):] {
		// TODO: Refine logic so players who finished early don't have to keep passing.
		if !m.IsPass() {
			return false
		}
	}
	return true
}
