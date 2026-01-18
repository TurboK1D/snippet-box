package main

import (
	"log/slog"
	"net/http"
	"os"
)

// Routing in Go:
//
// Go's standard library provides http.ServeMux as a basic HTTP request multiplexer.
// It's not as feature-rich as third-party routers (chi, gorilla/mux) but:
//   - No external dependencies
//   - Thread-safe
//   - Sufficient for most applications
//
// Pattern Syntax:
//   "/path"          - Exact match (only /path)
//   "/path/{id}"     - Named parameter (captures one segment)
//   "/path/{id...}"  - Catch-all (captures multiple segments)
//   "/{$}"           - Exact match with optional trailing slash handling
//   "GET /path"      - Method-specific route (Go 1.22+)

// routes returns an http.ServeMux with all application routes registered.
// The ServeMux acts as both a router and a handler - it receives all requests
// and dispatches them to the appropriate handler based on the URL pattern.
func (app *application) routes() *http.ServeMux {
	// Create a new ServeMux (request multiplexer).
	// ServeMux matches URL patterns against registered handlers.
	mux := http.NewServeMux()

	// Register handlers with method + path patterns.
	// Go 1.22+ syntax: "METHOD /path". Earlier versions use mux.HandleFunc separately.
	//
	// The "GET /{$}" pattern:
	//   - GET: only matches GET requests
	//   - /{$}: root path, {$} handles optional trailing slash
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create/{path...}", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Serve static files from the embedded filesystem.
	// Middleware wraps the FileServer to add cache headers.
	//
	// http.FileServer returns a Handler that serves files from a FileSystem.
	// http.FS(staticFiles) converts our embed.FS to http.FileSystem.
	// The "/static/" prefix is stripped before looking up the file.
	mux.Handle("GET /static/", cacheControl(http.FileServer(http.FS(staticFiles))))

	// Load templates on startup.
	// If templates fail to load, we can't serve HTML - exit with an error.
	// This is "fail fast" behavior: better to crash at startup than serve errors.
	if err := loadTemplates(); err != nil {
		app.logger.Error("failed to load templates", slog.Any("err", err))
		os.Exit(1)
	}

	return mux
}
