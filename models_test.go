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
			if got, want := c.IsColored(), true; got != want {
				t.Errorf("%v.IsColored(): got %v, want %v", c, got, want)
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
			if got, want := c.IsColored(), false; got != want {
				t.Errorf("%v.IsColored(): got %v, want %v", c, got, want)
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
				t.Errorf("%v.IsColored(): got %v, want %v", tc.c, got, tc.want)
			}
		})
	}
}

func TestEmptyColorString(t *testing.T) {
	c := colorEmpty
	if got, want := c.String(), "empty"; got != want {
		t.Errorf("%v.String(): got %q, want %q", c, got, want)
	}
}

func TestInvalidColorString(t *testing.T) {
	var c Color
	c = 100
	if got, want := c.String(), "unknown color"; got != want {
		t.Errorf("%v.String(): got %q, want %q", c, got, want)
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
			if got, want := b.IsOutOfBounds(c), false; got != want {
				t.Errorf("Board.IsOutOfBounds(%v): got %v, want %v", c, got, want)
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
			if got, want := b.IsOutOfBounds(c), true; got != want {
				t.Errorf("Board.IsOutOfBounds(%v): got %v, want %v", c, got, want)
			}
		})
	}
}

func TestNewPlayer(t *testing.T) {
	p, err := NewPlayer("brown beard", Red, Coord{0, 0}, 2)
	if err != nil {
		t.Errorf("NewPlayer: got %v, want no error", err)
	}
	if p == nil {
		t.Errorf("NewPlayer: got nil player, want player")
	}
}

func TestNewPlayerNoName(t *testing.T) {
	if _, err := NewPlayer("", Red, Coord{0, 0}, 2); err == nil || !strings.Contains(err.Error(), "name") {
		t.Errorf("NewPlayer with no name: got %v, want error about name", err)
	}
}

func TestNewPlayerInvalidColor(t *testing.T) {
	if _, err := NewPlayer("brown beard", 0, Coord{0, 0}, 2); err == nil || !strings.Contains(err.Error(), "color") {
		t.Errorf("NewPlayer with invalid color: got %v, want error about color", err)
	}
}

func TestNewPlayerNoPieces(t *testing.T) {
	if _, err := NewPlayer("brown beard", Red, Coord{0, 0}, 0); err == nil || !strings.Contains(err.Error(), "pieces") {
		t.Errorf("NewPlayer with no pieces: got %v, want error about pieces", err)
	}
}

func TestPlayerPlacePiece(t *testing.T) {
	player, err := NewPlayer("brown beard", Red, Coord{0, 0}, 2)
	if err != nil {
		t.Fatalf("NewPlayer: got %v, want no error", err)
	}

	// Remove piece at invalid index
	if err := player.placePiece(-1); err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Errorf("placePiece(-1): got %v, want index out of range error", err)
	}
	if err := player.placePiece(2); err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Errorf("placePiece(2): got %v, want index out of range error", err)
	}
	// Remove piece at valid index
	if err := player.placePiece(0); err != nil {
		t.Errorf("placePiece(0): got %v, want no error", err)
	}
	// Remove already removed piece
	if err := player.placePiece(0); err == nil || !strings.Contains(err.Error(), "already") {
		t.Errorf("placePiece(0) again: got %v, want already used error", err)
	}
}

func TestNewPieceNil(t *testing.T) {
	_, err := NewPiece(nil)
	if err == nil {
		t.Error("NewPiece(nil): got no error, want error")
	}
}

func TestNewPieceEmpty(t *testing.T) {
	_, err := NewPiece([]Coord{})
	if err == nil {
		t.Error("NewPiece({}): got no error, want error")
	}
}

func TestNewPiece(t *testing.T) {
	blocks := []Coord{{0, 0}, {0, 1}}
	p, err := NewPiece(blocks)
	if err != nil {
		t.Fatalf("NewPiece(): got %v, want no error", err)
	}
	if p == nil {
		t.Fatal("Piece: got nil, want not nil")
	}
	if !reflect.DeepEqual(p.Blocks, blocks) {
		t.Errorf("Piece blocks: got %v, want %v", p.Blocks, blocks)
	}
	if len(p.Corners()) == 0 {
		t.Error("Piece Corners(): got no corners, want at least one corner")
	}
}

func TestNewPieceOrNilReturnNil(t *testing.T) {
	p := NewPieceOrNil(nil)
	if p != nil {
		t.Errorf("NewPieceOrNil(nil): got %v, want nil", p)
	}
}

func TestNewPieceOrNil(t *testing.T) {
	p := NewPieceOrNil([]Coord{{0, 0}, {0, 1}})
	if p == nil {
		t.Error("NewPieceOrNil(): got nil, want not nil")
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

/*
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
*/
