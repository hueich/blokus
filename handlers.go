package blokusService

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	// "log"
	"path"

	"github.com/gorilla/mux"
	"github.com/hueich/blokus"
)

var (
	// Our poor man's database. At least it'll be fast!
	gamesDB = make(map[blokus.GameID]*blokus.Game)
)

func getGamesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.Write([]byte("<html><head></head><body>\n"))
	w.Write([]byte(`<form name="myform" method="POST">
					  <input type="text" name="gid"/>
					  <button type="submit">Create Game</button>
					</form>`))
	w.Write([]byte(fmt.Sprintf("<div>There are %d games:</div>\n<ul>\n", len(gamesDB))))
	for k := range gamesDB {
		w.Write([]byte(fmt.Sprintf("<li><a href=\"games/%v\">%v</a></li>", k, k)))
	}
	w.Write([]byte("</ul></body></html>"))
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	gidStr := r.FormValue("gid")
	gidInt, err := strconv.ParseInt(gidStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid game ID"))
		return
	}
	gid := blokus.GameID(gidInt)
	if _, exists := gamesDB[gid]; exists {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Game with ID already exists"))
		return
	}
	// g := blokus.NewGame()
	gamesDB[gid] = nil

	w.WriteHeader(http.StatusCreated)

	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		// TODO: May need to escape URL path.
		p := path.Join(r.URL.Path, strconv.FormatInt(gidInt, 10))
		w.Write([]byte(fmt.Sprintf(`<html><head><meta http-equiv="refresh" content="0;url='%s'"/></head></html>`, p)))
	}
}

func getGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gid, err := strconv.ParseInt(vars["gid"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid game ID"))
		return
	}
	g, ok := gamesDB[blokus.GameID(gid)]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No game found"))
		return
	}
	w.Write([]byte(fmt.Sprintf("Got game: %v", g)))
}

func getGameStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handled by getGameStateHandler()"))
}

func newPlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handled by newPlayerHandler()"))
}

func newMoveHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handled by newMoveHandler()"))
}
