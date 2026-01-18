package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// HTTP Handler Pattern in Go:
//
// A handler is any function with signature: func(http.ResponseWriter, *http.Request)
// Alternatively, implement the http.Handler interface: type Handler interface { ServeHTTP(ResponseWriter, *Request) }
//
// The http.ResponseWriter:
//   - Is an interface that collects the HTTP response
//   - Headers must be set BEFORE writing the body
//   - Once you Write() or WriteHeader() is called, headers are sent
//
// The *http.Request:
//   - Contains all client information (headers, URL, body, etc.)
//   - Body must be fully read then closed to prevent leaks
//   - PathValue() extracts named path segments (e.g., {id} in /snippet/view/{id})

// home is the handler for the homepage (GET /).
// Notice the method receiver: (app *application) means this is a method on our
// application struct, giving it access to app.logger and other dependencies.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Headers should be set before writing the response body.
	// This identifies the server in response headers.
	w.Header().Add("server", "go")

	// Execute the base template with nil data.
	// ts is a package-level variable (see templates.go) containing parsed templates.
	err := ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		// Centralized error handling - see helper.go for details.
		// We return early to prevent sending a partial/broken response.
		app.serverError(w, r, err)
		return
	}
}

// snippetView displays a single snippet by ID.
// r.PathValue("id") extracts the {id} from the route pattern /snippet/view/{id}
// The return type is string, so we must convert with strconv.Atoi().
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("snippet-view accessed")

	// Extract and validate the ID from the URL path.
	i, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || i < 1 {
		// http.NotFound is a convenient shortcut for 404 responses.
		// It writes the status and a simple body automatically.
		http.NotFound(w, r)
		return
	}

	// fmt.Fprintf formats and writes directly to the ResponseWriter.
	// This is simpler than fmt.Fprint + w.Write() but less flexible.
	if _, err := fmt.Fprintf(w, "Display specific snippet with ID: %d", i); err != nil {
		// We log but don't return a 500 - partial writes are acceptable here.
		log.Printf("write error: %v", err)
	}
}

// snippetCreate demonstrates path parameter capture with {path...}.
// The ellipsis (...) captures multiple path segments as a single value.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("snippet-create accessed")
	w.Header().Set("content-type", "text/plain;charset=utf-8")

	// r.PathValue("path") captures the wildcard path segment(s).
	path := r.PathValue("path")
	if _, err := fmt.Fprintf(w, "Captured paths: %s", path); err != nil {
		log.Printf("write error: %v", err)
	}
}

// snippetCreatePost handles POST requests to create new snippets.
// The path doesn't have wildcards - it matches exactly /snippet/create.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("snippet-create-post accessed")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// WriteHeader before Write - headers are sent when Write() is called,
	// but explicit WriteHeader gives more control over the status code.
	w.WriteHeader(http.StatusCreated) // 201 Created

	if _, err := fmt.Fprintf(w, "Snippet created successfully"); err != nil {
		log.Printf("write error: %v", err)
	}
}
