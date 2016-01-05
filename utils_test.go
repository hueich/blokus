package blokus

import (
	"testing"

	"golang.org/x/net/context"
)

func TestGetPlayerPiecesNilClientAndKey(t *testing.T) {
	if _, err := getPlayerPieces(context.Background(), nil, nil); err == nil {
		t.Fatalf("GetPlayerPieces(): got nil, want error")
	}
}
