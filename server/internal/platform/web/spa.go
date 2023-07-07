package web

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

// spaFS implements http.Handler to serve an SPA
type spaFS struct {
	fileSystem fs.FS
	httpFS     http.Handler
}

// ServeHTTP allows the spaFS struct to implement the http.Handler interface for serving an SPA
func (s *spaFS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the path and convert it to a filepath
	upath := r.URL.Path
	upath = strings.TrimPrefix(upath, "/")
	upath = path.Clean(upath)
	// If we are accessing the root, skip SPA routing and handle normally
	if upath != "/" && upath != "" {
		// Open the file
		f, err := s.fileSystem.Open(upath)
		if err != nil {
			// If the error was not an os.IsNotExist error, a real error occurred
			if !os.IsNotExist(err) {
				HandleError(r.Context(), w, err)
				return
			}
			// If the file does not exist, return index.html and let the SPA client router handle routing
			// because this route may be a route on the client
			indexFile, err := s.fileSystem.Open("index.html")
			if err != nil {
				HandleError(r.Context(), w, err)
				return
			}
			defer indexFile.Close()
			iBytes, err := io.ReadAll(indexFile)
			if err != nil {
				HandleError(r.Context(), w, err)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(iBytes)
			return
		}
		f.Close()
	}
	// If the path was "/" or the file was found, let the file server handle it
	s.httpFS.ServeHTTP(w, r)
}
