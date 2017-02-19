package blokusWebUI

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

type UIService struct {
	apiURL string
}

func NewService(r *mux.Router, apiURL string) (*UIService, error) {
	if r == nil {
		return nil, fmt.Errorf("Router cannot be nil")
	}
	if apiURL == "" {
		return nil, fmt.Errorf("API path cannot be empty")
	}
	s := &UIService{apiURL: apiURL}
	s.addRoutes(r)
	return s, nil
}

func (s *UIService) addRoutes(r *mux.Router) {
	r = r.PathPrefix("/games").Subrouter()

	// Gets a webpage with a listing of games.
	r.HandleFunc("", s.getGamesHandler).Methods("GET")

	g := r.PathPrefix("/{gid:[0-9]+}").Subrouter()
	// Gets a webpage showing the specified game.
	g.HandleFunc("", s.getGameHandler).Methods("GET")
}

func (s *UIService) getGamesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Move this to use html templates
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.Write([]byte("<html><head></head><body>\n"))
	w.Write([]byte(fmt.Sprintf(`<form name="myform" method="POST" action="%s">
					  <input type="text" name="gid"/>
					  <button type="submit">Create Game</button>
					</form>`, path.Join(s.apiURL, "games"))))

	// TODO: Add client side game fetching code
	w.Write([]byte("</body></html>"))
}

func (s *UIService) getGameHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Page to show game UI
}
