package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
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
		cfg.clientError(w, r, http.StatusNotFound)
		return
	}

	snippet, err := cfg.dbQueries.GetSnippet(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cfg.clientError(w, r, http.StatusNotFound)
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
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (cfg *config) handlerSnippetPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		cfg.clientError(w, r, http.StatusBadRequest)
		return
	}

	var form snippetCreateForm
	err = cfg.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		cfg.clientError(w, r, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.MinChars(form.Title, 3), "title", "This field cannot be less than 3 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(form.Expires > 0, "expires", "This field cannot be less than one")

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

	expiresDuration := time.Duration(int(time.Hour) * 24 * form.Expires)

	snippet, err := cfg.dbQueries.CreateSnippet(r.Context(), database.CreateSnippetParams{
		ExpiresAt: time.Now().UTC().Add(expiresDuration),
		Title:     form.Title,
		Content:   form.Content,
	})
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}

	cfg.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	cfg.redirect(w, r, fmt.Sprintf("/snippets/view/%s", snippet.ID.String()))
}

func (cfg *config) handlerSnippetDelete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.clientError(w, r, http.StatusNotFound)
		return
	}

	err = cfg.dbQueries.DeleteSnippet(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cfg.clientError(w, r, http.StatusNotFound)
		} else {
			cfg.serverError(w, r, err)
		}
		return
	}

	cfg.sessionManager.Put(r.Context(), "flash", "Snippet successfully deleted!")
	cfg.redirect(w, r, "/")
}

func (cfg *config) handlerUserSignup(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Display a form for signing up a new user...")
}

func (cfg *config) handlerUserSignupPost(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Create a new user...")
}

func (cfg *config) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Display a form for logging in a user...")
}

func (cfg *config) handlerUserLoginPost(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Authenticate and login the user...")
}

func (cfg *config) handlerUserLogoutPost(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Logout the user...")
}
