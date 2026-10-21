package main

import "github.com/davidheeren/snippetbox/internal/database"

type templateData struct {
	Snippet database.Snippet
	Snippets []database.Snippet
}
