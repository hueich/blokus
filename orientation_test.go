package blokus

import (
	"reflect"
	"testing"
)

func TestRotationNormalize(t *testing.T) {
	var rot Rotation
	rot = Rot180
	if got, want := Normalize(rot), Rot180; got != want {
		t.Errorf("Normalize(Rot180): got %v, want %v", got, want)
	}

	rot = Rot180 + 4
	if got, want := Normalize(rot), Rot180; got != want {
		t.Errorf("Normalize(Rot180 + 4): got %v, want %v", got, want)
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

func TestTransformCoords(t *testing.T) {
	cs := []Coord{
		{0, 0},
		{1, 0},
		{1, 1},
		{2, 1},
	}

	o := Orientation{Rot0, false}
	if got, want := o.TransformCoords(cs), []Coord{
		{0, 0},
		{1, 0},
		{1, 1},
		{2, 1},
	}; !reflect.DeepEqual(got, want) {
		t.Errorf("TransformCoords(%v): got %v, want %v", o, got, want)
	}
}
