package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Centralized Error Handling Pattern:
//
// Instead of repeating error handling logic in every handler, we create
// reusable helper methods on the application struct. Benefits:
//
//   1. DRY (Don't Repeat Yourself) - one place to change error handling
//   2. Consistency - all errors handled the same way
//   3. Security - ensure internal errors never leak to users
//   4. Logging - capture request context (method, URI) automatically
//
// See: Let's Go book chapter on centralized error handling

// serverError handles internal server errors (5xx).
// It logs the error with request context and returns a generic 500 response.
//
// Important: We never expose the actual error message to the user.
// This prevents leaking internal implementation details (database queries,
// file paths, stack traces) to potential attackers.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	// Extract request metadata for context.
	// This helps debug which endpoint failed and with what request.
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		// debug.Stack() returns the goroutine's current stack trace.
		// Useful for debugging but can be noisy in production.
		// Consider using runtime.Caller() for lighter file:line info.
		trace  = string(debug.Stack())
	)

	// slog.String() creates an slog.Value for structured logging.
	// Structured logging makes it easy to filter and search logs.
	app.logger.Error(err.Error(),
		slog.String("method", method),
		slog.String("uri", uri),
		slog.String("trace", trace),
	)

	// Send a generic 500 response.
	// http.StatusText(500) returns "Internal Server Error".
	// Never send err.Error() directly - it may contain sensitive info.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError sends an HTTP error response for client errors (4xx).
// Use this for validation errors, not found, unauthorized, etc.
//
// Examples:
//   app.clientError(w, http.StatusBadRequest)      // 400
//   app.clientError(w, http.StatusUnauthorized)    // 401
//   app.clientError(w, http.StatusForbidden)       // 403
//   app.clientError(w, http.StatusNotFound)        // 404
func (app *application) clientError(w http.ResponseWriter, status int) {
	// http.Error is a helper that sets the status, Content-Type, and body.
	// It also calls w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// by default (unless the Content-Type is already set).
	http.Error(w, http.StatusText(status), status)
}
