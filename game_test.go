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

func TestPlacePieceLocationOutOfBound(t *testing.T) {
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
			t.Errorf("PlacePiece(loc:%v): got no error, want location out of bounds error", c)
		} else if !strings.Contains(err.Error(), "location out of bounds") {
			t.Errorf("PlacePiece(loc:%v): got %v, want location out of bounds error", c, err)
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
	ps := []*Piece{
		{blocks: []Coord{Coord{3, 4}}},
		{blocks: []Coord{Coord{5, 6}}},
	}
	g, err := NewGame(123, 22, ps)
	if err != nil {
		t.Fatalf("NewGame(): got %v, want no error", err)
	}
	if err := g.AddPlayer("foo", Red, Coord{-1, -1}); err != nil {
		t.Fatalf("AddPlayer(foo, Red): got %v, want no error", err)
	}
	if err := g.AddPlayer("bar", Blue, Coord{-1, 1}); err != nil {
		t.Fatalf("AddPlayer(bar, Blue): got %v, want no error", err)
	}
	bluePiece := g.players[1].pieces[0]
	err = g.PlacePiece(Coord{0, 0}, bluePiece.id, 0, false)
	if err == nil {
		t.Error("PlacePiece(): got no error, want out of turn error")
	} else if !strings.Contains(err.Error(), "turn") {
		t.Errorf("PlacePiece(): got error %v, want out of turn error")
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
