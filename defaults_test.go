package blokus

import (
	"testing"
)

func TestDefaultPieces(t *testing.T) {
	pieces := DefaultPieces()

	if len(pieces) == 0 {
		t.Error("DefaultPieces: got 0 pieces, want more than 0")
	}
	for i, p := range pieces {
		if p == nil {
			t.Errorf("DefaultPieces[%d]: got nil, want not nil", i)
		}
	}
}
