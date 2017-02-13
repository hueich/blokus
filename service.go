package blokusService

import (
	"context"
	"net/http"
	"path"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
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
// For both projectID and credsFile, they can alternatively be provided through environment variables
// DATASTORE_PROJECT_ID and GOOGLE_APPLICATION_CREDENTIALS, respectively, in which case the cooresponding params can be left empty.
func (s *BlokusService) InitClient(ctx context.Context, projectID, credsFile string) error {
	opts := []option.ClientOption{}
	if credsFile != "" {
		opts = append(opts, option.WithServiceAccountFile(credsFile))
	}
	c, err := datastore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return err
	}
	s.client = c
	return nil
}

func (s *BlokusService) numGames(ctx context.Context) (int, error) {
	return s.client.Count(ctx, datastore.NewQuery("Game"))
}
