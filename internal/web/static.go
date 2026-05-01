package web

import (
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

type SPAOptions struct {
	APIPrefixes []string
}

func NewSPAHandler(options *SPAOptions) (http.Handler, error) {
	assets := FS()
	if _, err := fs.Stat(assets, "index.html"); err != nil {
		return nil, err
	}
	prefixes := []string{"/api", "/auth"}
	if options != nil && len(options.APIPrefixes) > 0 {
		prefixes = options.APIPrefixes
	}
	fileServer := http.FileServer(http.FS(assets))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, prefix := range prefixes {
			if r.URL.Path == prefix || strings.HasPrefix(r.URL.Path, prefix+"/") {
				http.NotFound(w, r)
				return
			}
		}
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "." || name == "" {
			name = "index.html"
		}
		if stat, err := fs.Stat(assets, name); err == nil && !stat.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
			http.Error(w, "static asset error", http.StatusInternalServerError)
			return
		}
		serveIndex(w, r, assets)
	}), nil
}

func serveIndex(w http.ResponseWriter, r *http.Request, assets fs.FS) {
	contents, err := fs.ReadFile(assets, "index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(contents)
}
