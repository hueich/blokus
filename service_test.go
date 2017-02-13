package blokusService

import (
	"context"
	"testing"
)

func TestNewServiceNoPrefix(t *testing.T) {
	if s := New(""); s == nil {
		t.Error("New(''): got service==nil, want service")
	}
}

func TestNewServiceWithPrefix(t *testing.T) {
	if s := New("/foo"); s == nil {
		t.Error("New('/foo'): got service==nil, want service")
	}
}

func TestEndToEnd(t *testing.T) {
	s := New("/foo")
	if s == nil {
		t.Fatal("New('/foo'): got service==nil, want service")
	}
	if err := s.InitClient(context.Background(), ""); err != nil {
		t.Fatalf("InitClient(''): got %v, want no error", err)
	}
	defer s.Close()

	count, err := s.numGames(context.Background())
	if err != nil {
		t.Fatalf("numGames(): got %v, want no error", err)
	}
	t.Errorf("Got numGames(): %v", count)
}
