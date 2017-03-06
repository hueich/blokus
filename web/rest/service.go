package rest

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

type APIService struct {
	client *datastore.Client
}

type Options struct {
	Router    *mux.Router
	ProjectID string
	CredsFile string
}

func NewService(opts *Options) (*APIService, error) {
	if opts.Router == nil {
		return nil, fmt.Errorf("Router cannot be nil")
	}
	c, err := newDBClient(context.Background(), opts.ProjectID, opts.CredsFile)
	if err != nil {
		return nil, err
	}
	s := &APIService{
		client: c,
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

// newDBClient creates a new Google Datastore client.
// For both projectID and credsFile, they can alternatively be provided through environment variables
// DATASTORE_PROJECT_ID and GOOGLE_APPLICATION_CREDENTIALS, respectively, in which case the cooresponding params can be left empty.
func newDBClient(ctx context.Context, projectID, credsFile string) (*datastore.Client, error) {
	opts := []option.ClientOption{}
	if credsFile != "" {
		opts = append(opts, option.WithServiceAccountFile(credsFile))
	}
	c, err := datastore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *APIService) numGames(ctx context.Context) (int, error) {
	return s.client.Count(ctx, datastore.NewQuery("Game"))
}
