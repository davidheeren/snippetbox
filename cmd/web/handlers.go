package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/davidheeren/snippetbox/internal/database"
	"github.com/davidheeren/snippetbox/internal/validator"
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
	data.Form = snippetCreateForm{
		Expires: 7,
	}
	cfg.render(w, r, http.StatusOK, "create.html", data)
}

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (cfg *config) handlerSnippetPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		cfg.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		cfg.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(expires > 0, "expires", "This field cannot be less than one")

	if !form.Valid() {
		data := cfg.newTemplateData(r)
		data.Form = form
		if r.Header.Get("HX-Request") == "true" {
			cfg.renderPage(w, r, http.StatusOK, "create.html", data)
		} else {
			cfg.renderPage(w, r, http.StatusUnprocessableEntity, "create.html", data)
		}
		return
	}

	expiresDuration := time.Duration(int(time.Hour) * 24 * expires)

	snippet, err := cfg.dbQueries.CreateSnippet(r.Context(), database.CreateSnippetParams{
		ExpiresAt: time.Now().UTC().Add(expiresDuration),
		Title:     form.Title,
		Content:   form.Content,
	})
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}

	cfg.redirect(w, r, fmt.Sprintf("/snippets/view/%s", snippet.ID.String()))
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

	cfg.redirect(w, r, "/")
}
