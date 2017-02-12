package blokusService

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

type BlokusService struct {
	router *mux.Router
}

func New(prefix string) *BlokusService {
	return &BlokusService{router: newRouter(prefix)}
}

func newRouter(prefix string) *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix(path.Join(prefix, "/games")).Subrouter()

	// Gets a webpage with a listing of games.
	s.HandleFunc("", getGamesHandler).Methods("GET")
	// Creates a game.
	s.HandleFunc("", newGameHandler).Methods("POST")

	g := s.PathPrefix("/{gid:[0-9]+}").Subrouter()
	// Gets a webpage showing the specified game.
	g.HandleFunc("", getGameHandler).Methods("GET")
	// Gets various game state data.
	g.HandleFunc("/state", getGameStateHandler).Methods("GET")
	// Add a player.
	g.HandleFunc("/players", newPlayerHandler).Methods("POST")
	// Make a move in the game.
	g.HandleFunc("/moves", newMoveHandler).Methods("POST")

	return r
}

func (s *BlokusService) Router() http.Handler {
	return s.router
}
