package main

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/davidheeren/snippetbox/internal/database"
	"github.com/davidheeren/snippetbox/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         database.GetSnippetWithUserRow
	Snippets        []database.GetSnippetsWithUserRow
	Form            any
	Flash           string
	IsAuthenticated bool
	AuthUser        database.User
}

func (cfg *config) newTemplateData(r *http.Request) templateData {
	user, ok := r.Context().Value(authenticatedUserContextKey).(database.User)
	data := templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           cfg.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: cfg.isAuthenticated(r),
	}
	if ok {
		data.AuthUser = user
	}
	return data
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	funcs := template.FuncMap{
		"iso": func(t time.Time) string { return t.UTC().Format(time.RFC3339) },
	}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(funcs).ParseFS(ui.Files, "html/base.html")
		if err != nil {
			return nil, err
		}

		partials, err := fs.Glob(ui.Files, "html/partials/*.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFS(ui.Files, partials...)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFS(ui.Files, page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
