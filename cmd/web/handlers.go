package main

import (
	"database/sql"
	"errors"
	"net/http"
	"text/template"

	"github.com/google/uuid"
)

func (cfg *config) handlerHome(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles(
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html")
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}
}

func (cfg *config) handlerSnippetView(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.clientError(w, http.StatusNotFound)
		return
	}

	snippet, err := cfg.dbQueries.GetSnippet(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cfg.clientError(w, http.StatusNotFound)
		} else {
			cfg.serverError(w, r, err)
		}
		return
	}

	data := templateData{
		Snippet: snippet,
	}

	ts, err := template.ParseFiles(
		"./ui/html/base.html",
		"./ui/html/pages/view.html",
		"./ui/html/partials/nav.html")
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}
}

func (cfg *config) handlerSnippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (cfg *config) handlerSnippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}
