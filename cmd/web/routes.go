package main

import (
	"net/http"

	"github.com/BirykovRV/miniature-broccoli/internal/lib"
	"github.com/justinas/alice"
)

func (app *application) routes(cfg config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(lib.NeuteredFileSystem{
		Fs: http.Dir(cfg.staticDir),
	})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("/", dynamic.ThenFunc(app.home))
	mux.Handle("/snippet/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("DELETE /snippet/{id}", dynamic.ThenFunc(app.snippetDelete))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	standart := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standart.Then(mux)
}
