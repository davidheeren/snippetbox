package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (cfg *config) serverError(w http.ResponseWriter, r *http.Request, err error) {
	// cfg.logger.Error(err.Error(), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()), slog.String("trace", string(debug.Stack())))
	cfg.logger.Error(err.Error(), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (cfg *config) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (cfg *config) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := cfg.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		cfg.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		cfg.serverError(w, r, err)
	}
}
