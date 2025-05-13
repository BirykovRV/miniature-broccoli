package main

import (
	"net/http"

	"github.com/BirykovRV/miniature-broccoli/internal/lib"
)

func (app *application) routes(cfg config) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(lib.NeuteredFileSystem{
		Fs: http.Dir(cfg.staticDir),
	})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
