package blokus

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestColorIsValid(t *testing.T) {
	colors := []Color{
		Blue,
		Yellow,
		Red,
		Green,
	}
	for _, c := range colors {
		t.Run(fmt.Sprintf("Color(%v)", c), func(t *testing.T) {
			if got, want := c.IsValid(), true; got != want {
				t.Errorf("%v.IsValid(): got %v, want %v", c, got, want)
			}
		})
	}
}

func TestColorIsNotValid(t *testing.T) {
	colors := []Color{
		0,
		100,
	}
	for _, c := range colors {
		t.Run(fmt.Sprintf("Color(%v)", c), func(t *testing.T) {
			if got, want := c.IsValid(), false; got != want {
				t.Errorf("%v.IsValid(): got %v, want %v", c, got, want)
			}
		})
	}
}

func TestColorString(t *testing.T) {
	testCases := []struct {
		c    Color
		want string
	}{
		{Blue, "blue"},
		{Yellow, "yellow"},
		{Red, "red"},
		{Green, "green"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Color(%v)", tc.c), func(t *testing.T) {
			if got := tc.c.String(); got != tc.want {
				t.Errorf("%v.IsValid(): got %v, want %v", tc.c, got, tc.want)
			}
		})
	}
}

func TestInvalidColorString(t *testing.T) {
	colors := []Color{
		0,
		100,
	}
	for _, c := range colors {
		t.Run(fmt.Sprintf("Color(%v)", c), func(t *testing.T) {
			if got, want := c.String(), ""; got != want {
				t.Errorf("%v.IsValid(): got %q, want %q", c, got, want)
			}
		})
	}
}

func TestBoardCoordsWithinBounds(t *testing.T) {
	coords := []Coord{
		{0, 0},
		{0, 9},
		{9, 0},
		{9, 9},
		{4, 5},
	}
	b := NewBoard(10)
	for _, c := range coords {
		t.Run(fmt.Sprintf("Coord(%v)", c), func(t *testing.T) {
			if got, want := b.isOutOfBounds(c), false; got != want {
				t.Errorf("Board.isOutOfBounds(%v): got %v, want %v", c, got, want)
			}
		})
	}
}

func TestBoardCoordsOutOfBounds(t *testing.T) {
	coords := []Coord{
		{0, -1},
		{-1, 0},
		{0, 10},
		{10, 0},
		{100, 100},
	}
	b := NewBoard(10)
	for _, c := range coords {
		t.Run(fmt.Sprintf("Coord(%v)", c), func(t *testing.T) {
			if got, want := b.isOutOfBounds(c), true; got != want {
				t.Errorf("Board.isOutOfBounds(%v): got %v, want %v", c, got, want)
			}
		})
	}
}

func TestPieceInheritsPlayerColor(t *testing.T) {
	player := Player{
		name:  "foo",
		color: Yellow,
	}
	piece := Piece{
		player: &player,
	}
	if got, want := piece.Color(), Yellow; got != want {
		t.Errorf("Piece.Color(): got %v, want %v", got, want)
	}
}

func TestPlayerRemovePiece(t *testing.T) {
	piece0 := &Piece{}
	piece1 := &Piece{}
	player := &Player{
		name:  "brown beard",
		color: Red,
		pieces: []*Piece{
			piece0,
			piece1,
		},
		startPos: Coord{},
	}
	// Remove piece at invalid index
	if _, err := player.RemovePiece(-1); err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Errorf("RemovePiece(-1): got %v, want index out of range error", err)
	}
	if _, err := player.RemovePiece(2); err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Errorf("RemovePiece(2): got %v, want index out of range error", err)
	}
	// Remove piece at valid index
	p, err := player.RemovePiece(0)
	if err != nil {
		t.Errorf("RemovePiece(0): got %v, want no error", err)
	}
	if p != piece0 {
		t.Errorf("RemovePiece(0): Got piece %v, want piece %v", p, piece0)
	}
	// Remove already removed piece
	if _, err := player.RemovePiece(0); err == nil || !strings.Contains(err.Error(), "already") {
		t.Errorf("RemovePiece(0) again: got %v, want already used error", err)
	}
}

func coordSliceToMap(s []Coord) map[Coord]bool {
	m := map[Coord]bool{}
	for _, c := range s {
		m[c] = true
	}
	return m
}

func TestGetCorners_OneBlock(t *testing.T) {
	blocks := []Coord{
		{0, 0},
	}
	want := map[Coord]bool{
		{-1, -1}: true,
		{-1, 1}:  true,
		{1, -1}:  true,
		{1, 1}:   true,
	}
	if got := coordSliceToMap(getCorners(blocks)); !reflect.DeepEqual(got, want) {
		t.Errorf("getCorners(): got %v, want %v", got, want)
	}
}

func TestGetCorners_TwoBlock(t *testing.T) {
	blocks := []Coord{
		{0, 0},
		{1, 0},
	}
	want := map[Coord]bool{
		{-1, -1}: true,
		{-1, 1}:  true,
		{2, -1}:  true,
		{2, 1}:   true,
	}
	if got := coordSliceToMap(getCorners(blocks)); !reflect.DeepEqual(got, want) {
		t.Errorf("getCorners(): got %v, want %v", got, want)
	}
}

func TestGetCorners_ThreeBlockL(t *testing.T) {
	blocks := []Coord{
		{0, 0},
		{1, 0},
		{1, 1},
	}
	want := map[Coord]bool{
		{-1, -1}: true,
		{-1, 1}:  true,
		{0, 2}:   true,
		{2, -1}:  true,
		{2, 2}:   true,
	}
	if got := coordSliceToMap(getCorners(blocks)); !reflect.DeepEqual(got, want) {
		t.Errorf("getCorners(): got %v, want %v", got, want)
	}
}

func TestRotatePiece(t *testing.T) {
	p := NewPiece(nil, []Coord{
		{0, 0},
		{1, 0},
		{1, 1},
		{2, 1},
	})
	bOrig := p.blocks
	bWant := p.blocks
	for i, b := range bWant {
		bWant[i] = rotateCoord(b)
	}
	cOrig := p.corners
	cWant := p.corners
	for i, c := range cWant {
		cWant[i] = rotateCoord(c)
	}
	p.Rotate()
	if got := p.blocks; !reflect.DeepEqual(got, bWant) {
		t.Errorf("Rotated blocks: got %v, want %v", got, bWant)
	}
	if got := p.corners; !reflect.DeepEqual(got, cWant) {
		t.Errorf("Rotated corners: got %v, want %v", got, cWant)
	}
	if got, want := p.rot, 1; got != want {
		t.Errorf("Rotation state: got %v, want %v", got, want)
	}

	p.Rotate()
	if got, want := p.rot, 2; got != want {
		t.Errorf("Rotation state: got %v, want %v", got, want)
	}

	p.Rotate()
	if got, want := p.rot, 3; got != want {
		t.Errorf("Rotation state: got %v, want %v", got, want)
	}

	p.Rotate()
	if got, want := p.rot, 0; got != want {
		t.Errorf("Rotation state: got %v, want %v", got, want)
	}
	if got := p.blocks; !reflect.DeepEqual(got, bOrig) {
		t.Errorf("4th rotation blocks: got %v, want %v", got, bOrig)
	}
	if got := p.corners; !reflect.DeepEqual(got, cOrig) {
		t.Errorf("4th rotated corners: got %v, want %v", got, cOrig)
	}
}

func TestFlipPiece(t *testing.T) {
	p := NewPiece(nil, []Coord{
		{0, 0},
		{1, 0},
		{1, 1},
		{2, 1},
	})
	bWant := p.blocks
	for i, b := range bWant {
		bWant[i] = flipCoord(b)
	}
	cWant := p.corners
	for i, c := range cWant {
		cWant[i] = flipCoord(c)
	}
	p.Flip()
	if got := p.blocks; !reflect.DeepEqual(got, bWant) {
		t.Errorf("Flipped blocks: got %v, want %v", got, bWant)
	}
	if got := p.corners; !reflect.DeepEqual(got, cWant) {
		t.Errorf("Flipped corners: got %v, want %v", got, cWant)
	}
	if !p.flip {
		t.Error("Flipped state: got false, want true")
	}
}
