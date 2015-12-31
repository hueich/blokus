package blokus

import (
	"reflect"
	"testing"
)

func TestColorName(t *testing.T) {
	if want, got := "red", ColorName(Red); want != got {
		t.Errorf("ColorName(Red): want %v, got %v", want, got)
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
	if want, got := Yellow, piece.GetColor(); want != got {
		t.Errorf("Piece.GetColor(): want %v, got %v", want, got)
	}
}

func coordSliceToMap(s []Coord) map[Coord]bool {
	m := map[Coord]bool{}
	for _, c := range s {
		m[c] = true
	}
	return m
}

func testBoardCoordTranslation(t *testing.T, b *Board, c Coord, p *Piece) {
	b.SetCoord(c.X, c.Y, p)
	if s := b.GetCoord(c.X, c.Y); s != p {
		t.Errorf("GetCoord(%v): got %v, want %v", c, s, p)
	}
	b.SetCoord(c.X, c.Y, nil)
	if s := b.GetCoord(c.X, c.Y); s != nil {
		t.Fatalf("GetCoord(%v): got %v, want nil", c, s)
	}
}

func TestBoardCoordTranslation(t *testing.T) {
	b := NewBoard(3)
	if got, want := len(b.grid), 3*3; got != want {
		t.Fatalf("NewBoard grid slice length: got %v, want %v", got, want)
	}
	p := NewPiece(123, nil, nil)

	testBoardCoordTranslation(t, b, Coord{0, 0}, p)
	testBoardCoordTranslation(t, b, Coord{0, 1}, p)
	testBoardCoordTranslation(t, b, Coord{1, 0}, p)
	testBoardCoordTranslation(t, b, Coord{1, 2}, p)
	testBoardCoordTranslation(t, b, Coord{2, 2}, p)
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

func TestRotateCoord(t *testing.T) {
	var c Coord

	c = Coord{0, 0}
	if got, want := rotateCoord(c), (Coord{0, 0}); got != want {
		t.Errorf("Rotating %v, got %v, want %v", c, got, want)
	}

	c = Coord{1, 2}
	if got, want := rotateCoord(c), (Coord{2, -1}); got != want {
		t.Errorf("Rotating %v, got %v, want %v", c, got, want)
	}

	c = Coord{2, -1}
	if got, want := rotateCoord(c), (Coord{-1, -2}); got != want {
		t.Errorf("Rotating %v, got %v, want %v", c, got, want)
	}

	c = Coord{-1, -2}
	if got, want := rotateCoord(c), (Coord{-2, 1}); got != want {
		t.Errorf("Rotating %v, got %v, want %v", c, got, want)
	}

	c = Coord{-2, 1}
	if got, want := rotateCoord(c), (Coord{1, 2}); got != want {
		t.Errorf("Rotating %v, got %v, want %v", c, got, want)
	}
}

func TestRotatePiece(t *testing.T) {
	p := NewPiece(123, nil, []Coord{
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

func TestFlipCoord(t *testing.T) {
	var c Coord

	c = Coord{0, 0}
	if got, want := flipCoord(c), (Coord{0, 0}); got != want {
		t.Errorf("Flipping %v, got %v, want %v", c, got, want)
	}

	c = Coord{1, 2}
	if got, want := flipCoord(c), (Coord{1, -2}); got != want {
		t.Errorf("Flipping %v, got %v, want %v", c, got, want)
	}

	c = Coord{1, -2}
	if got, want := flipCoord(c), (Coord{1, 2}); got != want {
		t.Errorf("Flipping %v, got %v, want %v", c, got, want)
	}

	c = Coord{-1, -2}
	if got, want := flipCoord(c), (Coord{-1, 2}); got != want {
		t.Errorf("Flipping %v, got %v, want %v", c, got, want)
	}

	c = Coord{-1, 2}
	if got, want := flipCoord(c), (Coord{-1, -2}); got != want {
		t.Errorf("Flipping %v, got %v, want %v", c, got, want)
	}
}

func TestFlipPiece(t *testing.T) {
	p := NewPiece(123, nil, []Coord{
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
