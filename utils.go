package blokus

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"
)

func getPlayerPieces(ctx context.Context, c *datastore.Client, playerKey *datastore.Key) ([]*Piece, error) {
	if c == nil || playerKey == nil {
		return nil, fmt.Errorf("Client and player key cannot be nil")
	}
	if playerKey.Kind() != "Player" {
		return nil, fmt.Errorf("Player key must be of kind Player")
	}
	p := make([]*Piece)
	q := datastore.NewQuery("Piece").Ancestor(playerKey)
	if _, err := c.GetAll(ctx, q, p); err != nil {
		return nil, err
	}
	return p, nil
}
