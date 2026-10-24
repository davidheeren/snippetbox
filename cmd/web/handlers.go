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
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (cfg *config) handlerSnippetCreate(w http.ResponseWriter, r *http.Request) {
	data := cfg.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 7,
	}
	cfg.render(w, r, http.StatusOK, "create.html", data)
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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (cfg *config) handlerUserSignup(w http.ResponseWriter, r *http.Request) {
	data := cfg.newTemplateData(r)
	data.Form = userSignupForm{}
	cfg.render(w, r, http.StatusOK, "signup.html", data)
}

func (cfg *config) handlerUserSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		cfg.clientError(w, r, http.StatusBadRequest)
		return
	}

	var form userSignupForm
	err = cfg.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		cfg.clientError(w, r, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 4), "password", "Password must be at least 4 characters")

	if !form.Valid() {
		data := cfg.newTemplateData(r)
		data.Form = form
		if r.Header.Get("HX-Request") == "true" {
			cfg.renderPage(w, r, http.StatusOK, "signup.html", data)
		} else {
			cfg.renderPage(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		}
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), 12)
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}

	_, err = cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Name:           form.Name,
		Email:          form.Email,
		HashedPassword: string(hash),
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" && pqErr.Constraint == "users_email_key" {
				form.CheckField(false, "email", "This email is already in use")
			}
		} else {
			cfg.serverError(w, r, err)
			return
		}
	}

	if !form.Valid() {
		data := cfg.newTemplateData(r)
		data.Form = form
		if r.Header.Get("HX-Request") == "true" {
			cfg.renderPage(w, r, http.StatusOK, "signup.html", data)
		} else {
			cfg.renderPage(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		}
		return
	}

	cfg.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	cfg.redirect(w, r, "/users/login")
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (cfg *config) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	data := cfg.newTemplateData(r)
	data.Form = userLoginForm{}
	cfg.render(w, r, http.StatusOK, "login.html", data)
}

func (cfg *config) handlerUserLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		cfg.clientError(w, r, http.StatusBadRequest)
		return
	}

	var form userLoginForm
	err = cfg.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		cfg.clientError(w, r, http.StatusBadRequest)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), form.Email)
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(form.Password))
		if err != nil {
			form.AddNonFieldError("Email or password is incorrect")
		}
	} else {
		form.AddNonFieldError("Email or password is incorrect")
	}

	if !form.Valid() {
		data := cfg.newTemplateData(r)
		data.Form = form
		if r.Header.Get("HX-Request") == "true" {
			cfg.renderPage(w, r, http.StatusOK, "login.html", data)
		} else {
			cfg.renderPage(w, r, http.StatusUnprocessableEntity, "login.html", data)
		}
		return
	}

	err = cfg.sessionManager.RenewToken(r.Context())
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}

	cfg.sessionManager.Put(r.Context(), "authenticatedUserID", user.ID.String())
	cfg.sessionManager.Put(r.Context(), "flash", "You've been logged in successfully!")
	cfg.redirect(w, r, "/")
}

func (cfg *config) handlerUserLogoutPost(w http.ResponseWriter, r *http.Request) {
    err := cfg.sessionManager.RenewToken(r.Context())
    if err != nil {
        cfg.serverError(w, r, err)
        return
    }

    cfg.sessionManager.Remove(r.Context(), "authenticatedUserID")

    cfg.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	cfg.redirect(w, r, "/")
}
