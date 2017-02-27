package blokusWebAPI

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"github.com/hueich/blokus"
)

type gameInfo struct {
	ID int64
}

func (s *APIService) getGamesHandler(w http.ResponseWriter, r *http.Request) {
	q := datastore.NewQuery("Game")
	q = q.KeysOnly()
	keys, err := s.client.GetAll(r.Context(), q, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not get list of games: %v\n", err)
		return
	}

	games := make([]*gameInfo, 0, len(keys))
	for _, k := range keys {
		games = append(games, &gameInfo{ID: k.ID})
	}

	b, err := json.Marshal(games)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not marshal list of games: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(b)
}

func (s *APIService) newGameHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Error checking, player validation, etc.

	g, err := blokus.NewGame(blokus.DefaultBoardSize, blokus.DefaultPieces())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not instantiate new game: %v\n", err)
		return
	}
	// TODO: Also add current player to game.
	gameKey := datastore.IncompleteKey("Game", nil)
	gameKey, err = s.client.Put(r.Context(), gameKey, g)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not put new game: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusCreated)

	b, err := json.Marshal(gameInfo{ID: gameKey.ID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not marshal created game info: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(b)
}

func (s *APIService) getGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gid, err := strconv.ParseInt(vars["gid"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid game ID"))
		return
	}
	g := &blokus.Game{}
	gameKey := datastore.IDKey("Game", gid, nil)
	if err := s.client.Get(r.Context(), gameKey, g); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No game found"))
		return
	}
	w.Write([]byte(fmt.Sprintf("Got game: %v", *g)))
}

func (s *APIService) getGameStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handled by getGameStateHandler()"))
}

func (s *APIService) newPlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handled by newPlayerHandler()"))
}

func (s *APIService) newMoveHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handled by newMoveHandler()"))
}
