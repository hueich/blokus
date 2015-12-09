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
