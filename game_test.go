package blokus

import (
	"strings"
	"testing"
)

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

func TestPlacePieceLocationOutOfBound(t *testing.T) {
	g, err := NewGame(123, 10, []*Piece{})
	if err != nil {
		t.Fatalf("NewGame(): got %v, want no error", err)
	}
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
	g, err := NewGame(123, 10, []*Piece{})
	if err != nil {
		t.Fatalf("NewGame(): got %v, want no error", err)
	}
	rots := []int{-10, -1, 4, 100}
	for _, r := range rots {
		err := g.PlacePiece(Coord{1,1}, 1, r, false)
		if err == nil {
			t.Errorf("PlacePiece(rot:%v): got no error, want invalid rotation error", r)
		} else if !strings.Contains(err.Error(), "rotation") {
			t.Errorf("PlacePiece(rot:%v): got %v, want invalid rotation error", r, err)
		}
	}
}
