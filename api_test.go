package blokus

import (
	"context"
	"testing"
)

type DummyService struct {
	gameID GameID
}

func (d *DummyService) CreateGame(ctx context.Context, username, gamename string, boardSize int, opt *GameOptions) (GameID, error) {
	return d.gameID, nil
}

func (d *DummyService) AddPlayer(ctx context.Context, id GameID, username string, opt *PlayerOptions) error {
	return nil
}

func (d *DummyService) StartGame(ctx context.Context, id GameID, username string) error {
	return nil
}

func (d *DummyService) PlacePiece(ctx context.Context, id GameID, pieceID int, x, y, rot int, flip bool) error {
	return nil
}

func (d *DummyService) GetGameState(ctx context.Context, id GameID) error {
	return nil
}

func newDummyService(id GameID) Service {
	return &DummyService{gameID: id}
}

func TestCreateGameWithoutOptions(t *testing.T) {
	s := newDummyService(1234)
	id, err := s.CreateGame(context.Background(), "some_user", "some game", 20, nil)
	if err != nil {
		t.Errorf("CreateGame: got err %v, want no error", err)
	}
	if got, want := id, GameID(1234); got != want {
		t.Errorf("CreateGame: got %v, want %v", got, want)
	}
}

func TestCreateGameWithOptions(t *testing.T) {
	s := newDummyService(1234)
	opts := &GameOptions{
		Description: "This is my game!",
	}
	id, err := s.CreateGame(context.Background(), "some_user", "some game", 20, opts)
	if err != nil {
		t.Errorf("CreateGame: got err %v, want no error", err)
	}
	if got, want := id, GameID(1234); got != want {
		t.Errorf("CreateGame: got %v, want %v", got, want)
	}
}

func TestAddPlayerWithoutOptions(t *testing.T) {
	s := newDummyService(1234)
	if err := s.AddPlayer(context.Background(), GameID(1234), "some_user", nil); err != nil {
		t.Errorf("AddPlayer: got err %v, want no error", err)
	}
}

func TestAddPlayerWithOptions(t *testing.T) {
	s := newDummyService(1234)
	opts := &PlayerOptions{
		Color: Red,
	}
	if err := s.AddPlayer(context.Background(), GameID(1234), "some_user", opts); err != nil {
		t.Errorf("AddPlayer: got err %v, want no error", err)
	}
}
