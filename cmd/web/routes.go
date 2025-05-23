package main

import (
	"net/http"

	"github.com/BirykovRV/miniature-broccoli/internal/lib"
	"github.com/BirykovRV/miniature-broccoli/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// fileServer := http.FileServer(lib.NeuteredFileSystem{
	// 	Fs: http.Dir(cfg.staticDir),
	// })
	fileServer := http.FileServer(http.FS(ui.Files))

	mux.Handle("/static", http.NotFoundHandler())
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/static/", fileServer)

	baseChain := lib.Chain{
		app.recoverPanic,
		app.logRequest,
		secureHeaders,
	}
	// health check handler
	mux.Handle("/ping", baseChain.ThenFunc(ping))

	dynamic := append(baseChain, app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	authChain := append(dynamic, app.requireAuthentication)

	mux.Handle("/", dynamic.ThenFunc(app.home))
	mux.Handle("/snippet/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("/about", dynamic.ThenFunc(app.aboutView))

	mux.Handle("/user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("/user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	// protected app routes
	mux.Handle("DELETE /snippet/{id}", authChain.ThenFunc(app.snippetDelete))
	mux.Handle("GET /snippet/create", authChain.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", authChain.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", authChain.ThenFunc(app.userLogoutPost))

	return baseChain.Then(mux)
}
