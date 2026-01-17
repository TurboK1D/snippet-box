package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("server", "go")

	err := ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || i < 1 {
		http.NotFound(w, r)
		return
	}
	if _, err := fmt.Fprintf(w, "Display specific snippet with ID: %d", i); err != nil {
		log.Printf("write error: %v", err)
	}
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	if _, err := fmt.Fprintf(w, "Captured paths: %s", path); err != nil {
		log.Printf("write error: %v", err)
	}
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := fmt.Fprintf(w, "Snippet created successfully"); err != nil {
		log.Printf("write error: %v", err)
	}
}
