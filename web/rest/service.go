package rest

import (
	"context"
	"errors"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
)

type APIService struct {
	client *datastore.Client
}

type Options struct {
	Router *mux.Router
	Client *datastore.Client
}

func NewService(opts Options) (*APIService, error) {
	if opts.Router == nil {
		return nil, errors.New("REST service: router cannot be nil")
	}
	if opts.Client == nil {
		return nil, errors.New("REST service: datastore client cannot be nil")
	}
	s := &APIService{
		client: opts.Client,
	}
	s.addRoutes(opts.Router)
	return s, nil
}

func (s *APIService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

func (s *APIService) addRoutes(r *mux.Router) {
	r = r.PathPrefix("/games").Subrouter()

	// Gets a webpage with a listing of games.
	r.HandleFunc("", s.getGamesHandler).Methods("GET")
	// Creates a game.
	r.HandleFunc("", s.newGameHandler).Methods("POST")

	g := r.PathPrefix("/{gid:[0-9]+}").Subrouter()
	// Gets a webpage showing the specified game.
	g.HandleFunc("", s.getGameHandler).Methods("GET")
	// Gets various game state data.
	g.HandleFunc("/state", s.getGameStateHandler).Methods("GET")
	// Add a player.
	g.HandleFunc("/players", s.newPlayerHandler).Methods("POST")
	// Make a move in the game.
	g.HandleFunc("/moves", s.newMoveHandler).Methods("POST")
}

func (s *APIService) numGames(ctx context.Context) (int, error) {
	return s.client.Count(ctx, datastore.NewQuery("Game"))
}
