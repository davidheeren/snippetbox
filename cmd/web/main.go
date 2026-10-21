package main

import (
	"flag"
	"log"
	"net/http"
)

type arguments struct {
	addr string
}

func main() {
	var args arguments
	flag.StringVar(&args.addr, "addr", ":4000", "HTTP network address")

	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handlerHome)
	mux.HandleFunc("GET /snippet/view/{id}", handlerSnippetView)
	mux.HandleFunc("GET /snippet/create", handlerSnippetCreate)
	mux.HandleFunc("POST /snippet/create", handlerSnippetCreatePost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:    args.addr,
		Handler: mux,
	}

	log.Printf("Serving on: http://localhost%s\n", args.addr)
	log.Fatal(srv.ListenAndServe())
}
