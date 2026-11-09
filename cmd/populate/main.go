package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davidheeren/snippetbox/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/loremipsum.v1"
)

func main() {
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
	const testEmail string = "user@test.com"
	const snippetCount int = 10

	// get test user if already created
	user, err := dbQueries.GetUserByEmail(context.Background(), testEmail)
	createUser := false
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			createUser = true
		} else {
			log.Fatal(err)
		}
	}

	// if there is no test user, then make a new one
	if createUser {
		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
		if err != nil {
			log.Fatal(err)
		}

		user, err = dbQueries.CreateUser(context.Background(), database.CreateUserParams{
			Name:           "Test User",
			Email:          testEmail,
			HashedPassword: string(hash),
			IsAdmin:        false,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(os.Args) < 2 {
		createSnippets(dbQueries, user, snippetCount)
	} else if os.Args[1] == "delete" {
		deleteSnippets(dbQueries, user)
	} else if count, err := strconv.Atoi(os.Args[1]); err == nil {
		createSnippets(dbQueries, user, count)
	} else {
		fmt.Printf("options: 'delete', 'snippetCount {%d}'\n", snippetCount)
	}
}

func deleteSnippets(dbQueries *database.Queries, user database.User) {
	_, err := dbQueries.DeleteSnippetsByUser(context.Background(), user.ID)
	if err != nil {
		log.Fatal(err)
	}
}

func createSnippets(dbQueries *database.Queries, user database.User, count int) {
	lorem := loremipsum.New()

	for _ = range count {
		// 2-6 words keeps the title well under the 100 char limit
		title := lorem.Words(rand.IntN(5) + 2)
		title = strings.ToUpper(title[:1]) + title[1:]

		// Paragraph() has no newlines, so join several with blank lines
		content := strings.Join(lorem.SentenceList(rand.IntN(3)+2), "\n")

		_, err := dbQueries.CreateSnippet(context.Background(), database.CreateSnippetParams{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
			Title:     title,
			Content:   content,
			UserID:    user.ID,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Created %d new snippets\n", count)
}
