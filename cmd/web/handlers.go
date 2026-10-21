package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
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
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		cfg.serverError(w, r, err)
		return
	}

	msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)
	w.Write([]byte(msg))
}

func (cfg *config) handlerSnippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func (cfg *config) handlerSnippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}
