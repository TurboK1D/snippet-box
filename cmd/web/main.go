package main

import (
	"log"
	"net/http"
)

func main() {
	if err := loadTemplates(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create/{path...}", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// Server embeded static files
	mux.Handle("GET /static/", http.FileServer(http.FS(staticFiles)))

	log.Print("starting server on :4000")

	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}
