package rest

import (
	"context"
	"testing"

	"github.com/gorilla/mux"
)

func TestNewServiceNilRouter(t *testing.T) {
	if _, err := NewService(nil, nil); err == nil {
		t.Errorf("NewService(nil, nil): got no error, want error")
	}
}

func TestNewService(t *testing.T) {
	r := mux.NewRouter()
	if s, err := NewService(r, nil); err != nil {
		t.Fatalf("NewService(router): got error %v, want no error", err)
	} else if s == nil {
		t.Fatal("NewService(router): got service==nil, want service")
	}
	// Make sure routes were added.
	count := 0
	if err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		count = count + 1
		return nil
	}); err != nil {
		t.Errorf("Router.Walk(): got error %v, want no error", err)
	}
	if count <= 0 {
		t.Errorf("Router.Walk() count: got %v, want > 0", count)
	}
}

func TestEndToEnd(t *testing.T) {
	r := mux.NewRouter()
	s, err := NewService(r, nil)
	if err != nil {
		t.Fatalf("NewService(router): got error %v, want no error", err)
	}
	if s == nil {
		t.Fatal("NewService(router): got service==nil, want service")
	}
	defer s.Close()

	count, err := s.numGames(context.Background())
	if err != nil {
		t.Fatalf("numGames(): got %v, want no error", err)
	}
	t.Logf("Got numGames(): %v", count)
}
