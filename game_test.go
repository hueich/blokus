package blokus

import (
	"strings"
	"testing"
)

func newGameOrDie(t *testing.T) *Game {
	g, err := NewGame(123, 10, []*Piece{})
	if err != nil {
		t.Fatalf("NewGame(): got %v, want no error", err)
	}
	return g
}

func TestNewGame(t *testing.T) {
	p := &Piece{blocks: []Coord{Coord{3, 4}}}
	g, err := NewGame(123, 22, []*Piece{p})
	if err != nil {
		t.Fatalf("NewGame(): got %v, want no error", err)
	}

	if got, want := g.id, GameID(123); got != want {
		t.Errorf("Game id: got %v, want %v", got, want)
	}
	if g.board == nil {
		t.Fatal("Game board: got nil, want not nil")
	}
	if got, want := len(g.board.grid), 22; got != want {
		t.Errorf("Game board size: got %v, want %v", got, want)
	}
	if got, want := len(g.pieceSet), 1; got != want {
		t.Fatalf("Game pieces len: got %v, want %v", got, want)
	}
	if got, want := len(g.pieceSet[0].blocks), 1; got != want {
		t.Fatalf("Game pieces[0] blocks len: got %v, want %v", got, want)
	}
	if got, want := g.pieceSet[0].blocks[0], (Coord{3, 4}); got != want {
		t.Errorf("Game pieces[0] blocks[0]: got %v, want %v", got, want)
	}
}

func TestAddPlayer(t *testing.T) {
	ps := []*Piece{
		{blocks: []Coord{Coord{3, 4}}},
		{blocks: []Coord{Coord{5, 6}}},
	}
	g, err := NewGame(123, 22, ps)
	if err != nil {
		t.Fatalf("NewGame(): got %v, want no error", err)
	}

	if err := g.AddPlayer("foo", Red, Coord{0, 0}); err != nil {
		t.Fatalf("AddPlayer(): got %v, want no error", err)
	}

	if got, want := len(g.players), 1; got != want {
		t.Fatalf("Game players len: got %v, want %v", got, want)
	}
	p := g.players[0]
	if got, want := p.name, "foo"; got != want {
		t.Errorf("Player name: got %v, want %v", got, want)
	}
	if got, want := p.color, Red; got != want {
		t.Errorf("Player color: got %v, want %v", got, want)
	}
	if got, want := p.startPos, (Coord{0, 0}); got != want {
		t.Errorf("Player start position: got %v, want %v", got, want)
	}
	if got, want := len(p.pieces), 2; got != want {
		t.Fatalf("Player pieces len: got %v, want %v", got, want)
	}
	if id0, id1 := p.pieces[0].id, p.pieces[1].id; id0 == id1 {
		t.Errorf("Player pieces ids: got %v == %v, want not equal", id0, id1)
	}

	if got, want := len(p.pieces[0].blocks), 1; got != want {
		t.Fatalf("Player pieces[0] blocks len: got %v, want %v", got, want)
	}
	if got, want := p.pieces[0].blocks[0], (Coord{3, 4}); got != want {
		t.Errorf("Player pieces[0] blocks[0]: got %v, want %v", got, want)
	}
	p.pieces[0].Rotate()
	if got, want := g.pieceSet[0].blocks[0], (Coord{3, 4}); got != want {
		t.Errorf("Game pieces[0] blocks[0] after rotating player's piece: got %v, want %v", got, want)
	}
}

func TestAddPlayerInvalidColor(t *testing.T) {
	g := newGameOrDie(t)
	if err, msg := g.AddPlayer("foo", 99, Coord{9, 9}), "color"; err == nil || !strings.Contains(err.Error(), msg) {
		t.Errorf("AddPlayer() with invalid color: got %v, want error to contain %v", err, msg)
	}
}

func TestAddPlayerInvalidStartPosition(t *testing.T) {
	g := newGameOrDie(t)
	if err := g.AddPlayer("foo", Blue, Coord{10, 10}); err == nil || !strings.Contains(err.Error(), "position is out of bounds") {
		t.Errorf("AddPlayer() with invalid start position: got %v, want starting position out of bounds error", err)
	}
}

func TestAddPlayerDupeName(t *testing.T) {
	g := newGameOrDie(t)
	name := "foo"
	if err := g.AddPlayer(name, Blue, Coord{9, 9}); err != nil {
		t.Fatalf("Add first player: got error %v, want no error", err)
	}

	if err := g.AddPlayer(name, Yellow, Coord{0, 0}); err == nil || !strings.Contains(err.Error(), name) {
		t.Errorf("AddPlayer() with duplicate name: got %v, want duplicate name error containing %v", err, name)
	}
}

func TestAddPlayerDupeColor(t *testing.T) {
	g := newGameOrDie(t)
	color := Blue
	if err := g.AddPlayer("foo", color, Coord{9, 9}); err != nil {
		t.Fatalf("Add first player: got error %v, want no error", err)
	}

	if err := g.AddPlayer("bar", color, Coord{0, 0}); err == nil || !strings.Contains(err.Error(), color.String()) {
		t.Errorf("AddPlayer() with duplicate color: got %v, want duplicate color error containing %v", err, color)
	}
}

func TestAddPlayerDupeStartPosition(t *testing.T) {
	g := newGameOrDie(t)
	if err := g.AddPlayer("foo", Blue, Coord{9, 9}); err != nil {
		t.Fatalf("Add first player: got error %v, want no error", err)
	}

	if err := g.AddPlayer("bar", Yellow, Coord{9, 9}); err == nil || !strings.Contains(err.Error(), "position already occupied") {
		t.Errorf("AddPlayer() with duplicate start position: got %v, want position already occupied error", err)
	}
}

func newGameWithTwoPlayersAndTwoPieces(t *testing.T, size int) *Game {
	tps := []*Piece{
		NewTemplatePiece([]Coord{{0, 0}, {1, 0}, {2, 0}}),
		NewTemplatePiece([]Coord{{0, 0}, {1, 0}, {1, 1}, {2, 1}}),
	}
	g, err := NewGame(123, size, tps)
	if err != nil {
		t.Fatalf("NewGame(): got %v, want no error", err)
	}
	if err := g.AddPlayer("foo", Blue, Coord{0, 0}); err != nil {
		t.Fatalf("AddPlayer(foo): got %v, want no error", err)
	}
	if err := g.AddPlayer("bar", Yellow, Coord{size - 1, size - 1}); err != nil {
		t.Fatalf("AddPlayer(bar): got %v, want no error", err)
	}
	return g
}

func TestPlacePieceOutOfBound(t *testing.T) {
	g := newGameOrDie(t)
	coords := []Coord{
		{-1, 5},
		{5, -1},
		{100, 5},
		{5, 100},
	}
	for _, c := range coords {
		err := g.PlacePiece(c, 1, Orientation{Rot0, false})
		if err == nil {
			t.Errorf("PlacePiece(loc:%v): got no error, want out of bounds error", c)
		} else if !strings.Contains(err.Error(), "out of bounds") {
			t.Errorf("PlacePiece(loc:%v): got %v, want out of bounds error", c, err)
		}
	}
}

func TestPlacePieceAlreadyPlaced(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)

	p := g.players[0].pieces[0]
	if err := g.PlacePiece(Coord{0, 0}, p.id, Orientation{Rot0, false}); err != nil {
		t.Fatalf("PlacePiece 1st time: got %v, want no error", err)
	}

	if err := g.PlacePiece(Coord{2, 2}, p.id, Orientation{Rot0, false}); err == nil || !strings.Contains(err.Error(), "already placed") {
		t.Errorf("PlacePiece() 2nd time: got %v, want piece already placed error", err)
	}
}

func TestPlacePieceOutOfTurn(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)

	p := g.players[1].pieces[0]
	if err := g.PlacePiece(Coord{0, 0}, p.id, Orientation{Rot0, false}); err == nil || !strings.Contains(err.Error(), "turn") {
		t.Errorf("PlacePiece(): got %v, want out of turn error", err)
	}
}

func TestPlacePieceValid(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	numMovesBefore := len(g.moves)

	p := g.players[0].pieces[0]
	if err := g.PlacePiece(Coord{0, 0}, p.id, Orientation{Rot270, false}); err != nil {
		t.Fatalf("PlacePiece(): got %v, want no error", err)
	}

	spaces := []Coord{
		{0, 0},
		{0, 1},
		{0, 2},
	}
	for _, s := range spaces {
		if got := g.board.grid[s.X][s.Y]; got != p.Color() {
			t.Errorf("Board space [%v]: got %v, want occupied by piece colored %v", s, got, p.Color())
		}
	}
	if p.location == nil {
		t.Fatalf("Piece location: got nil, want not nil")
	}
	if got, want := *p.location, (Coord{0, 0}); got != want {
		t.Errorf("Piece location: got %v, want %v", got, want)
	}

	// Move was recorded in the game.
	if got, want := len(g.moves), numMovesBefore+1; got != want {
		t.Errorf("Number of moves: got %v, want %v", got, want)
	}
	gotMove := g.moves[len(g.moves)-1]
	wantMove := Move{
		player: p.player,
		piece:  p,
		orient: Orientation{Rot270, false},
		loc:    Coord{0, 0},
	}
	if gotMove != wantMove {
		t.Errorf("Move info: got %v, want %v", gotMove, wantMove)
	}

	// Should not advance turn
	if got, want := g.curPlayerIndex, 0; got != want {
		t.Errorf("Player turn after placing piece: got %v, want %v", got, want)
	}
}

func TestCheckPiecePlacementNilPlayer(t *testing.T) {
	g := newGameOrDie(t)
	p := NewPiece(123, nil, []Coord{})
	if err := g.checkPiecePlacement(Coord{8, 9}, p); err == nil || !strings.Contains(err.Error(), "no owning player") {
		t.Errorf("checkPiecePlacement(): got %v, want no owning player error", err)
	}
}

func TestCheckPiecePlacementBlockIndexBoundsCheck(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	p := g.players[0].pieces[0]
	if err := g.checkPiecePlacementAt(Coord{0, 0}, p, -1); err == nil || !strings.Contains(err.Error(), "out of bounds") {
		t.Errorf("checkPiecePlacementAt([0,0]) with block index -1: got %v, want block index out of bounds error", err)
	}
	if err := g.checkPiecePlacementAt(Coord{0, 0}, p, 3); err == nil || !strings.Contains(err.Error(), "out of bounds") {
		t.Errorf("checkPiecePlacementAt([0,0]) with block index 3: got %v, want block index out of bounds error", err)
	}
	if err := g.checkPiecePlacementAt(Coord{0, 0}, p, 1); err == nil || !strings.Contains(err.Error(), "out of bounds") {
		t.Errorf("checkPiecePlacementAt([0,0]) with block index 1: got %v, want block index out of bounds error", err)
	}
	if err := g.checkPiecePlacementAt(Coord{1, 0}, p, 1); err != nil {
		t.Errorf("checkPiecePlacementAt([1,0]) with block index 1: got %v, want no error", err)
	}
}

func TestCheckPiecePlacementAtStartingPosition(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	p := g.players[0].pieces[1]
	if err := g.checkPiecePlacement(Coord{0, 0}, p); err != nil {
		t.Fatalf("checkPiecePlacement(): got %v, want no error", err)
	}
}

func TestCheckPiecePlacementWithValidTouchingCorner(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.PlacePiece(Coord{0, 0}, g.players[0].pieces[0].id, Orientation{Rot0, false}); err != nil {
		t.Fatalf("PlacePiece(blue 1st piece): got %v, want no error", err)
	}

	p := g.players[0].pieces[1]
	if err := g.checkPiecePlacement(Coord{3, 1}, p); err != nil {
		t.Fatalf("checkPiecePlacement(blue 2nd piece): got %v, want no error", err)
	}
}

func TestCheckPiecePlacementNoValidCorner(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.checkPiecePlacement(Coord{2, 2}, g.players[0].pieces[0]); err == nil || !strings.Contains(err.Error(), "no corner touching") {
		t.Errorf("checkPiecePlacement(): got %v, want no corner touching error", err)
	}
}

func TestCheckPiecePlacementOutOfBounds(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	p := g.players[0].pieces[0]
	if err := g.checkPiecePlacement(Coord{19, 19}, p); err == nil || !strings.Contains(err.Error(), "out of bounds") {
		t.Errorf("checkPiecePlacement(): got %v, want out of bounds error", err)
	}
}

func TestCheckPiecePlacementOverlapExistingPiece(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 3)
	// Place an existing piece
	if err := g.PlacePiece(Coord{0, 0}, g.players[0].pieces[1].id, Orientation{Rot0, false}); err != nil {
		t.Fatalf("Place initial piece: got %v, want no error", err)
	}

	p := g.players[1].pieces[1]
	if err := g.checkPiecePlacement(Coord{0, 1}, p); err == nil || !strings.Contains(err.Error(), "occupied") {
		t.Errorf("checkPiecePlacement(): got %v, want occupied error", err)
	}
}

func TestCheckPiecePlacementTouchingSameColor(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.PlacePiece(Coord{0, 0}, g.players[0].pieces[0].id, Orientation{Rot0, false}); err != nil {
		t.Fatalf("PlacePiece(blue 1st piece): got %v, want no error", err)
	}

	p := g.players[0].pieces[1]
	if err := g.checkPiecePlacement(Coord{2, 1}, p); err == nil || !strings.Contains(err.Error(), "another blue piece") {
		t.Errorf("checkPiecePlacement(blue 2nd piece): got error %v, want next to another blue piece error", err)
	}
}

func TestCheckPiecePlacementTouchingAnotherColor(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 3)
	if err := g.PlacePiece(Coord{0, 0}, g.players[0].pieces[0].id, Orientation{Rot0, false}); err != nil {
		t.Fatalf("PlacePiece(blue 1st piece): got %v, want no error", err)
	}
	if err := g.checkPiecePlacement(Coord{0, 2}, g.players[1].pieces[0]); err != nil {
		t.Fatalf("PlacePiece(yellow 1st piece): got %v, want no error", err)
	}
}

func TestAdvanceTurnNoPlayers(t *testing.T) {
	g := newGameOrDie(t)
	if err := g.AdvanceTurn(); err == nil || !strings.Contains(err.Error(), "no players") {
		t.Errorf("AdvanceTurn() with no players: got %v, want no players error", err)
	}
	if got := g.curPlayerIndex; got != 0 {
		t.Fatalf("After AdvanceTurn() returned error, curPlayerIndex: got %v, want 0", got)
	}
}

func TestAdvanceTurn(t *testing.T) {
	g := newGameOrDie(t)
	if err := g.AddPlayer("foo", Red, Coord{0, 0}); err != nil {
		t.Fatalf("Error adding player 'foo': %v", err)
	}
	if err := g.AddPlayer("bar", Blue, Coord{0, 9}); err != nil {
		t.Fatalf("Error adding player 'bar': %v", err)
	}

	if got, want := g.players[g.curPlayerIndex].name, "foo"; got != want {
		t.Fatalf("Before AdvanceTurn(): got player %v, want player %v", got, want)
	}
	if err := g.AdvanceTurn(); err != nil {
		t.Fatalf("AdvanceTurn() with players: got %v, want no error", err)
	}
	if got, want := g.players[g.curPlayerIndex].name, "bar"; got != want {
		t.Errorf("After AdvanceTurn(): got player %v, want player %v", got, want)
	}
}
