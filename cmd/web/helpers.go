package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/davidheeren/snippetbox/internal/database"
)

func (cfg *config) serverError(w http.ResponseWriter, r *http.Request, err error) {
	// cfg.logger.Error(err.Error(), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()), slog.String("trace", string(debug.Stack())))
	cfg.logger.Error(err.Error(), slog.String("method", r.Method), slog.String("uri", r.URL.RequestURI()))
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (cfg *config) clientError(w http.ResponseWriter, r *http.Request, status int) {
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(status)))
		return
	}
	http.Error(w, http.StatusText(status), status)
}

func (cfg *config) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := cfg.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		cfg.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (cfg *config) renderPage(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := cfg.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		cfg.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "main", data)
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (cfg *config) redirect(w http.ResponseWriter, r *http.Request, url string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", url)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (cfg *config) isAuthenticated(r *http.Request) bool {
	_, ok := r.Context().Value(authenticatedUserContextKey).(database.User)
	return ok
}
