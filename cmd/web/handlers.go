package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *config) handlerHome(w http.ResponseWriter, r *http.Request) {
	snippets, err := cfg.dbQueries.GetSnippets(r.Context())
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}

	data := cfg.newTemplateData(r)
	data.Snippets = snippets
	cfg.render(w, r, http.StatusOK, "home.html", data)
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

	data := cfg.newTemplateData(r)
	data.Snippet = snippet
	cfg.render(w, r, http.StatusOK, "view.html", data)
}

func (cfg *config) handlerSnippetCreate(w http.ResponseWriter, r *http.Request) {
	data := cfg.newTemplateData(r)
	cfg.render(w, r, http.StatusOK, "create.html", data)
}

func (cfg *config) handlerSnippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}
