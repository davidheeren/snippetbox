package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davidheeren/snippetbox/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("adding snippets to database")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbURL, ok := os.LookupEnv("DB_URL")
	if !ok {
		log.Fatal("DB_URL environment variable must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	duration := time.Hour * 10
	_, err = dbQueries.CreateSnippet(context.Background(), database.CreateSnippetParams{
		ExpiresAt: time.Now().UTC().Add(duration),
		Title:     "An old silent pond",
		Content:   "An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = dbQueries.CreateSnippet(context.Background(), database.CreateSnippetParams{
		ExpiresAt: time.Now().UTC().Add(duration),
		Title:     "Over the wintry forest",
		Content:   "Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = dbQueries.CreateSnippet(context.Background(), database.CreateSnippetParams{
		ExpiresAt: time.Now().UTC().Add(duration),
		Title:     "First autumn morning",
		Content:   "First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo",
	})
	if err != nil {
		log.Fatal(err)
	}
}
