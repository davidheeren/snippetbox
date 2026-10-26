package main

import (
	"net/http"

	"github.com/davidheeren/snippetbox/ui"
	"github.com/justinas/alice"
)

func (cfg *config) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", middlewareCacheStatic(http.FileServerFS(ui.Files)))

	dynamic := alice.New(cfg.sessionManager.LoadAndSave, middlewarePreventCSRF, cfg.middlewareAuthenticate)

	mux.Handle("GET /users/signup", dynamic.ThenFunc(cfg.handlerUserSignup))
	mux.Handle("POST /users", dynamic.ThenFunc(cfg.handlerUserSignupPost))
	mux.Handle("GET /users/login", dynamic.ThenFunc(cfg.handlerUserLogin))
	mux.Handle("POST /users/login", dynamic.ThenFunc(cfg.handlerUserLoginPost))
	mux.Handle("POST /users/logout", dynamic.ThenFunc(cfg.handlerUserLogoutPost))

	protected := dynamic.Append(cfg.middlewareRequireAuthentication)

	mux.Handle("GET /{$}", protected.ThenFunc(cfg.handlerHome))
	mux.Handle("GET /snippets/view/{id}", protected.ThenFunc(cfg.handlerSnippetView))
	mux.Handle("GET /snippets/create", protected.ThenFunc(cfg.handlerSnippetCreate))
	mux.Handle("POST /snippets", protected.ThenFunc(cfg.handlerSnippetPost))
	mux.Handle("DELETE /snippets/{id}", protected.ThenFunc(cfg.handlerSnippetDelete))

	standardMiddleware := alice.New(cfg.middlewareRecoverPanic, cfg.middlewareLogRequest, middlewareCommonHeaders)
	return standardMiddleware.Then(mux)
}
