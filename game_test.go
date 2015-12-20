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

	if got, want := g.id, 123; got != want {
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

	if err := g.AddPlayer("foo", Red, Coord{-1, -1}); err != nil {
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
		t.Errorf("Player color: got %v, want %v", ColorName(got), ColorName(want))
	}
	if got, want := p.corner, (Coord{-1, -1}); got != want {
		t.Errorf("Player corner: got %v, want %v", got, want)
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
	err := g.AddPlayer("foo", 999, Coord{1, 1})
	if err == nil {
		t.Fatalf("AddPlayer() with invalid color: got no error, want error")
	}
	if colorValue := "999"; !strings.Contains(err.Error(), colorValue) {
		t.Errorf("AddPlayer() with invalid color: got error %v, want error to contain %v", err, colorValue)
	}
}

func TestAddPlayerDupeName(t *testing.T) {
	g := newGameOrDie(t)
	name := "foo"
	if err := g.AddPlayer(name, Blue, Coord{1, 1}); err != nil {
		t.Fatalf("Add first player: got error %v, want no error", err)
	}

	err := g.AddPlayer(name, Yellow, Coord{-1, -1})
	if err == nil {
		t.Fatalf("AddPlayer duplicate name: got no error, want error")
	}
	if !strings.Contains(err.Error(), name) {
		t.Errorf("AddPlayer duplicate name: got error %v, want duplicate name error to contain %v", err, name)
	}
}

func TestAddPlayerDupeColor(t *testing.T) {
	g := newGameOrDie(t)
	color := Blue
	if err := g.AddPlayer("foo", color, Coord{1, 1}); err != nil {
		t.Fatalf("Add first player: got error %v, want no error", err)
	}

	err := g.AddPlayer("bar", color, Coord{-1, -1})
	if err == nil {
		t.Fatalf("AddPlayer duplicate color: got no error, want error")
	}
	if colorName := ColorName(color); !strings.Contains(err.Error(), colorName) {
		t.Errorf("AddPlayer duplicate color: got error %v, want duplicate color error to contain %v", err, colorName)
	}
}

func TestAddPlayerDupeCorner(t *testing.T) {
	g := newGameOrDie(t)
	if err := g.AddPlayer("foo", Blue, Coord{1, 1}); err != nil {
		t.Fatalf("Add first player: got error %v, want no error", err)
	}

	err := g.AddPlayer("bar", Yellow, Coord{1, 1})
	if err == nil {
		t.Fatalf("AddPlayer duplicate corner: got no error, want error")
	}
	if !strings.Contains(err.Error(), "Corner") {
		t.Errorf("AddPlayer duplicate corner: got error %v, want duplicate corner error", err)
	}
}

func newGameWithTwoPlayersAndTwoPieces(t *testing.T, size int) *Game {
	tps := []*Piece{
		NewTemplatePiece([]Coord{{0,0}, {1,0}, {2,0}}),
		NewTemplatePiece([]Coord{{0,0}, {1,0}, {1,1}, {2,1}}),
	}
	g, err := NewGame(123, size, tps)
	if err != nil {
		t.Fatalf("NewGame(): got %v, want no error", err)
	}
	if err := g.AddPlayer("foo", Blue, Coord{-1,-1}); err != nil {
		t.Fatalf("AddPlayer(foo): got %v, want no error", err)
	}
	if err := g.AddPlayer("bar", Yellow, Coord{size,size}); err != nil {
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
		err := g.PlacePiece(c, 1, 0, false)
		if err == nil {
			t.Errorf("PlacePiece(loc:%v): got no error, want out of bounds error", c)
		} else if !strings.Contains(err.Error(), "out of bounds") {
			t.Errorf("PlacePiece(loc:%v): got %v, want out of bounds error", c, err)
		}
	}
}

func TestPlacePieceInvalidRotation(t *testing.T) {
	g := newGameOrDie(t)
	rots := []int{-10, -1, 4, 100}
	for _, r := range rots {
		err := g.PlacePiece(Coord{1, 1}, 1, r, false)
		if err == nil {
			t.Errorf("PlacePiece(rot:%v): got no error, want invalid rotation error", r)
		} else if !strings.Contains(err.Error(), "rotation") {
			t.Errorf("PlacePiece(rot:%v): got %v, want invalid rotation error", r, err)
		}
	}
}

func TestPlacePieceOutOfTurn(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)

	p := g.players[1].pieces[0]
	err := g.PlacePiece(Coord{0, 0}, p.id, 0, false)
	if err == nil {
		t.Fatal("PlacePiece(): got no error, want out of turn error")
	} else if !strings.Contains(err.Error(), "turn") {
		t.Errorf("PlacePiece(): got error %v, want out of turn error")
	}
}

func TestPlacePieceValid(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)

	p := g.players[0].pieces[0]
	if err := g.PlacePiece(Coord{0, 0}, p.id, 3, false); err != nil {
		t.Fatalf("PlacePiece(): got %v, want no error", err)
	}

	spaces := []Coord{
		{0,0},
		{0,1},
		{0,2},
	}
	for _, s := range spaces {
		if got := g.board.grid[s.X][s.Y]; got != p {
			t.Errorf("Board space [%v]: got %v, want occupied by piece %v", got, p.id)
		}
	}
	if p.location == nil {
		t.Fatalf("Piece location: got nil, want not nil")
	}
	if got, want := *p.location, (Coord{0,0}); got != want {
		t.Errorf("Piece location: got %v, want %v", got, want)
	}
	if got, want := g.curPlayerIndex, 1; got != want {
		t.Errorf("Player turn after placing piece: got %v, want %v", got, want)
	}
}

func TestCheckPiecePlacementNilPlayer(t *testing.T) {
	g := newGameOrDie(t)
	p := NewPiece(123, nil, []Coord{})
	err := g.checkPiecePlacement(Coord{8,9}, p)
	if err == nil {
		t.Fatalf("checkPiecePlacement(): got no error, want error")
	}
	if !strings.Contains(err.Error(), "no owning player") {
		t.Errorf("checkPiecePlacement(): got error %v, want no owning player error")
	}
}

func TestCheckPiecePlacementAtStartingCorner(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	p := g.players[0].pieces[1]
	if err := g.checkPiecePlacement(Coord{0,0}, p); err != nil {
		t.Fatalf("checkPiecePlacement(): got %v, want no error", err)
	}
}

func TestCheckPiecePlacementWithValidTouchingCorner(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.PlacePiece(Coord{0,0}, g.players[0].pieces[0].id, 0, false); err != nil {
		t.Fatalf("PlacePiece(blue 1st piece): got %v, want no error", err)
	}
	if err := g.PlacePiece(Coord{7,9}, g.players[1].pieces[0].id, 0, false); err != nil {
		t.Fatalf("PlacePiece(yellow 1st piece): got %v, want no error", err)
	}

	p := g.players[0].pieces[1]
	if err := g.checkPiecePlacement(Coord{3,1}, p); err != nil {
		t.Fatalf("checkPiecePlacement(blue 2nd piece): got %v, want no error", err)
	}
}

func TestCheckPiecePlacementNoValidCorner(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	err := g.checkPiecePlacement(Coord{2,2}, g.players[0].pieces[0])
	if err == nil {
		t.Fatalf("checkPiecePlacement(): got no error, want error")
	}
	if !strings.Contains(err.Error(), "no corner touching") {
		t.Errorf("checkPiecePlacement(): got error %v, want no corner touching error")
	}
}

func TestCheckPiecePlacementOutOfBounds(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	p := g.players[0].pieces[0]
	err := g.checkPiecePlacement(Coord{19, 19}, p)
	if err == nil {
		t.Fatalf("checkPiecePlacement(): got no error, want error")
	}
	if !strings.Contains(err.Error(), "out of bounds") {
		t.Errorf("checkPiecePlacement(): got error %v, want out of bounds error", err)
	}
}

func TestCheckPiecePlacementOverlapExistingPiece(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 3)
	// Place an existing piece
	if err := g.PlacePiece(Coord{0,0}, g.players[0].pieces[1].id, 0, false); err != nil {
		t.Fatalf("Place initial piece: got %v, want no error", err)
	}

	p := g.players[1].pieces[1]
	err := g.checkPiecePlacement(Coord{0, 1}, p)
	if err == nil {
		t.Fatalf("checkPiecePlacement(): got no error, want error")
	}
	if !strings.Contains(err.Error(), "occupied") {
		t.Errorf("checkPiecePlacement(): got error %v, want occupied error", err)
	}
}

func TestCheckPiecePlacementTouchingSameColor(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.PlacePiece(Coord{0,0}, g.players[0].pieces[0].id, 0, false); err != nil {
		t.Fatalf("PlacePiece(blue 1st piece): got %v, want no error", err)
	}
	if err := g.PlacePiece(Coord{7,9}, g.players[1].pieces[0].id, 0, false); err != nil {
		t.Fatalf("PlacePiece(yellow 1st piece): got %v, want no error", err)
	}

	p := g.players[0].pieces[1]
	err := g.checkPiecePlacement(Coord{2,1}, p)
	if err == nil {
		t.Fatalf("checkPiecePlacement(blue 2nd piece): got no error, want error")
	}
	if !strings.Contains(err.Error(), "another blue piece") {
		t.Errorf("checkPiecePlacement(blue 2nd piece): got error %v, want next to another blue piece error", err)
	}
}

func TestCheckPiecePlacementTouchingAnotherColor(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 3)
	if err := g.PlacePiece(Coord{0,0}, g.players[0].pieces[0].id, 0, false); err != nil {
		t.Fatalf("PlacePiece(blue 1st piece): got %v, want no error", err)
	}
	if err := g.checkPiecePlacement(Coord{0, 2}, g.players[1].pieces[0]); err != nil {
		t.Fatalf("PlacePiece(yellow 1st piece): got %v, want no error", err)
	}
}

func TestAdvanceTurnNoPlayers(t *testing.T) {
	g := newGameOrDie(t)
	err := g.advanceTurn()
	if err == nil {
		t.Errorf("advanceTurn() with no players: got no error, want error")
	} else if !strings.Contains(err.Error(), "no players") {
		t.Errorf("advanceTurn() with no players: got %v, want no players error")
	}
	if got := g.curPlayerIndex; got != 0 {
		t.Fatalf("After advanceTurn() returned error, curPlayerIndex: got %v, want 0", got)
	}
}

func TestAdvanceTurn(t *testing.T) {
	g := newGameOrDie(t)
	g.AddPlayer("foo", Red, Coord{-1, -1})
	g.AddPlayer("bar", Blue, Coord{-1, 1})

	if got, want := g.players[g.curPlayerIndex].name, "foo"; got != want {
		t.Fatalf("Before advanceTurn(): got player %v, want player %v", got, want)
	}
	if err := g.advanceTurn(); err != nil {
		t.Fatalf("advanceTurn() with players: got %v, want no error", err)
	}
	if got, want := g.players[g.curPlayerIndex].name, "bar"; got != want {
		t.Errorf("After advanceTurn(): got player %v, want player %v", got, want)
	}
}
