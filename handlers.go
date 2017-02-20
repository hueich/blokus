package blokusWebAPI

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"github.com/hueich/blokus"
)

type game struct {
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

	games := make([]*game, 0, len(keys))
	for _, k := range keys {
		games = append(games, &game{ID: k.ID})
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
	gameKey := datastore.IncompleteKey("Game", nil)
	gameKey, err = s.client.Put(r.Context(), gameKey, g)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not put new game: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusCreated)

	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		// TODO: May need to escape URL path.
		p := path.Join(r.URL.Path, strconv.FormatInt(gameKey.ID, 10))
		w.Write([]byte(fmt.Sprintf(`<html><head><meta http-equiv="refresh" content="0;url='%s'"/></head></html>`, p)))
	}
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
