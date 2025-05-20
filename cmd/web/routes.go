package main

import (
	"net/http"

	"github.com/BirykovRV/miniature-broccoli/internal/lib"
)

func (app *application) routes(cfg config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(lib.NeuteredFileSystem{
		Fs: http.Dir(cfg.staticDir),
	})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	baseChain := lib.Chain{app.sessionManager.LoadAndSave,
		app.recoverPanic,
		app.logRequest,
		secureHeaders,
	}

	authChain := append(baseChain, app.requireAuthentication)

	mux.Handle("/", baseChain.ThenFunc(app.home))
	mux.Handle("/snippet/{id}", baseChain.ThenFunc(app.snippetView))

	mux.Handle("/user/signup", baseChain.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", baseChain.ThenFunc(app.userSignupPost))
	mux.Handle("/user/login", baseChain.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", baseChain.ThenFunc(app.userLoginPost))
	// protected app routes
	mux.Handle("DELETE /snippet/{id}", authChain.ThenFunc(app.snippetDelete))
	mux.Handle("GET /snippet/create", authChain.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", authChain.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", authChain.ThenFunc(app.userLogoutPost))

	return baseChain.Then(mux)
}
