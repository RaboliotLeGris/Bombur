package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type StaticHandler struct {
	StaticPath string
	IndexPath  string
}

// Stolen from mux README
func (s StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("GET ", r.URL.Path)
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(s.StaticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(s.StaticPath, s.IndexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(s.StaticPath)).ServeHTTP(w, r)
}
