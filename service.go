package blokusWebApp

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
	"github.com/gorilla/sessions"
)

const (
	importDir    = "github.com/hueich/blokus-web-app"
	templatesDir = "templates"
)

type AppService struct {
	apiURL   string
	tmplsDir string
	tmpls    *template.Template
	store    sessions.Store
}

type Options struct {
	Router       *mux.Router
	ApiURL       string
	TemplatesDir string
	Store        sessions.Store
}

func NewService(opt *Options) (*AppService, error) {
	if opt.Router == nil {
		return nil, fmt.Errorf("Router cannot be nil")
	}
	if opt.ApiURL == "" {
		return nil, fmt.Errorf("API URL cannot be empty")
	}
	if opt.Store == nil {
		return nil, fmt.Errorf("Session store cannot be nil")
	}

	templatesDir := opt.TemplatesDir
	if templatesDir == "" {
		d, err := getTemplatesDir()
		if err != nil {
			return nil, err
		}
		templatesDir = d
	} else if !isReadableDir(templatesDir) {
		return nil, fmt.Errorf("Could not read templates directory: %v", templatesDir)
	}
	log.Println("Using templates directory:", templatesDir)

	t, err := template.ParseGlob(filepath.Join(templatesDir, "*.gohtml"))
	if err != nil {
		return nil, err
	}

	s := &AppService{
		apiURL:   opt.ApiURL,
		tmplsDir: templatesDir,
		tmpls:    t,
		store:    opt.Store,
	}
	s.addRoutes(opt.Router)
	return s, nil
}

func (s *AppService) addRoutes(r *mux.Router) {
	r = r.PathPrefix("/games").Subrouter()

	// Gets a webpage with a listing of games.
	r.HandleFunc("", s.getGamesHandler).Methods("GET")

	g := r.PathPrefix("/{gid:[0-9]+}").Subrouter()
	// Gets a webpage showing the specified game.
	g.HandleFunc("", s.getGameHandler).Methods("GET")
}

func (s *AppService) getGamesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")

	//DEBUG
	var err error
	s.tmpls, err = template.ParseGlob(filepath.Join(s.tmplsDir, "*.gohtml"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not glob templates: %v\n", err)
		return
	}

	if err := s.tmpls.ExecuteTemplate(w, "main-games-view", map[string]interface{}{
		"NewGameURL":  path.Join(s.apiURL, "games"),
		"GetGamesURL": path.Join(s.apiURL, "games"),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not execute template: %v\n", err)
		return
	}

	// TODO: Add client side game fetching code
}

func (s *AppService) getGameHandler(w http.ResponseWriter, r *http.Request) {
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
