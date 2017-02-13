package blokusService

import (
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

func (s *BlokusService) getGamesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.Write([]byte("<html><head></head><body>\n"))
	w.Write([]byte(`<form name="myform" method="POST">
					  <input type="text" name="gid"/>
					  <button type="submit">Create Game</button>
					</form>`))

	q := datastore.NewQuery("Game")
	q = q.KeysOnly()
	keys, err := s.client.GetAll(r.Context(), q, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not get list of games: %v\n", err)
		return
	}
	w.Write([]byte(fmt.Sprintf("<div>There are %d games:</div>\n<ul>\n", len(keys))))
	for _, k := range keys {
		w.Write([]byte(fmt.Sprintf("<li><a href=\"games/%v\">%v</a></li>", k.ID, k.ID)))
	}
	w.Write([]byte("</ul></body></html>"))
}

func (s *BlokusService) newGameHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Error checking, player validation, etc.

	g, err := blokus.NewGame(blokus.GameID(0), blokus.DefaultBoardSize, blokus.DefaultPieces())
	log.Printf("Created new game in memory: %v\n", g)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	gameKey := datastore.IncompleteKey("Game", nil)
	gameKey, err = s.client.Put(r.Context(), gameKey, g)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not create new game: %v", err)
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

func (s *BlokusService) getGameHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *BlokusService) getGameStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handled by getGameStateHandler()"))
}

func (s *BlokusService) newPlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handled by newPlayerHandler()"))
}

func (s *BlokusService) newMoveHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handled by newMoveHandler()"))
}
