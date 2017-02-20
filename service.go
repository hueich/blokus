package blokusWebUI

import (
	"fmt"
	"go/build"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
)

const (
	importDir    = "github.com/hueich/blokus-web-ui"
	templatesDir = "templates"
)

type UIService struct {
	apiURL   string
	tmplsDir string
}

func NewService(r *mux.Router, apiURL, templatesDir string) (*UIService, error) {
	if r == nil {
		return nil, fmt.Errorf("Router cannot be nil")
	}
	if apiURL == "" {
		return nil, fmt.Errorf("API URL cannot be empty")
	}
	if templatesDir == "" {
		if dir, err := getTemplatesDir(); err != nil {
			return nil, err
		} else {
			templatesDir = dir
		}
	} else if !isReadableDir(templatesDir) {
		return nil, fmt.Errorf("Could not read templates directory: %v", templatesDir)
	}

	log.Println("Using templates directory:", templatesDir)
	s := &UIService{
		apiURL:   apiURL,
		tmplsDir: templatesDir,
	}
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
	t, err := template.ParseFiles(filepath.Join(s.tmplsDir, "games-list.gohtml"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not parse template files: %v\n", err)
		return
	}

	// TODO: Move this to use html templates
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.Write([]byte("<html><head></head><body>\n"))
	w.Write([]byte(fmt.Sprintf(`<form name="myform" method="POST" action="%s">
					  <input type="text" name="gid"/>
					  <button type="submit">Create Game</button>
					</form>`, path.Join(s.apiURL, "games"))))

	t.ExecuteTemplate(w, "games-list", []int64{123, 456})

	// TODO: Add client side game fetching code
	w.Write([]byte("</body></html>"))
}

func (s *UIService) getGameHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Page to show game UI
}

func isReadableDir(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()
	if info, err := f.Stat(); err != nil || !info.IsDir() {
		return false
	}
	return true
}

func getTemplatesDir() (string, error) {
	for _, sdir := range build.Default.SrcDirs() {
		tdir := filepath.Join(sdir, importDir, templatesDir)
		if isReadableDir(tdir) {
			return tdir, nil
		}
	}
	return "", fmt.Errorf("Could not find viable templates directory")
}
