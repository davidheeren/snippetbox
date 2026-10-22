package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/davidheeren/snippetbox/internal/database"
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

func (cfg *config) handlerSnippetPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		cfg.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	fieldErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field cannot be blank"
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		cfg.clientError(w, http.StatusBadRequest)
		return
	}

	if expires < 0 {
		fieldErrors["expires"] = "This field cannot be less than one"
	}

	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

	expiresDuration := time.Duration(int(time.Hour) * 24 * expires)

	snippet, err := cfg.dbQueries.CreateSnippet(r.Context(), database.CreateSnippetParams{
		ExpiresAt: time.Now().UTC().Add(expiresDuration),
		Title:     title,
		Content:   content,
	})
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets/view/%s", snippet.ID.String()), http.StatusSeeOther)
}

func (cfg *config) handlerSnippetDelete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.clientError(w, http.StatusNotFound)
		return
	}

	err = cfg.dbQueries.DeleteSnippet(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cfg.clientError(w, http.StatusNotFound)
		} else {
			cfg.serverError(w, r, err)
		}
		return
	}

	// regular redirects don't work with htmx
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
