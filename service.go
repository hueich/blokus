package blokusService

import (
	"context"
	"net/http"
	"path"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
)

type BlokusService struct {
	router *mux.Router
	client *datastore.Client
}

func New(prefix string) *BlokusService {
	s := &BlokusService{}
	s.initRouter(prefix)
	return s
}

func (s *BlokusService) Close() {
	if s.client != nil {
		s.client.Close()
		s.client = nil
	}
}

func (s *BlokusService) Router() http.Handler {
	return s.router
}

func (s *BlokusService) initRouter(prefix string) {
	r := mux.NewRouter()
	sr := r.PathPrefix(path.Join(prefix, "/games")).Subrouter()

	// Gets a webpage with a listing of games.
	sr.HandleFunc("", s.getGamesHandler).Methods("GET")
	// Creates a game.
	sr.HandleFunc("", s.newGameHandler).Methods("POST")

	g := sr.PathPrefix("/{gid:[0-9]+}").Subrouter()
	// Gets a webpage showing the specified game.
	g.HandleFunc("", s.getGameHandler).Methods("GET")
	// Gets various game state data.
	g.HandleFunc("/state", s.getGameStateHandler).Methods("GET")
	// Add a player.
	g.HandleFunc("/players", s.newPlayerHandler).Methods("POST")
	// Make a move in the game.
	g.HandleFunc("/moves", s.newMoveHandler).Methods("POST")

	s.router = r
}

// InitClient initializes the Google Datastore client.
// Credentials should be set via the environment variable: GOOGLE_APPLICATION_CREDENTIALS
// Project ID can alternatively be set via the environment variable: DATASTORE_PROJECT_ID
func (s *BlokusService) InitClient(ctx context.Context, projectID string) error {
	c, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	s.client = c
	return nil
}

func (s *BlokusService) numGames(ctx context.Context) (int, error) {
	return s.client.Count(ctx, datastore.NewQuery("Game"))
}
