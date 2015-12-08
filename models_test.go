package blokus

import "testing"

func TestColorName(t *testing.T) {
	if want, got := "red", ColorName(Red); want != got {
		t.Errorf("ColorName(Red): want %v, got %v", want, got)
	}
}

func TestPieceInheritsPlayerColor(t *testing.T) {
	player := Player{
		name: "foo",
		color: Yellow,
	}
	piece := Piece {
		player: &player,
	}
	if want, got := Yellow, piece.GetColor(); want != got {
		t.Errorf("Piece.GetColor(): want %v, got %v", want, got)
	}
}
