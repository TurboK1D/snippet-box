package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// application struct groups all dependencies shared across HTTP handlers.
// This is Go's "dependency injection" pattern - instead of passing individual
// loggers, configs, etc. to each handler, we group them here.
//
// Using a pointer receiver (*application) on handler methods allows mutation
// of shared state if needed, and avoids copying the struct on every call.
type application struct {
	logger *slog.Logger
}

func main() {
	// flag.String returns a pointer to the parsed value.
	// This is idiomatic Go: use pointers for optional/parsed values.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Parse flags (must call before using *addr)
	flag.Parse()

	// Create a structured logger writing to stdout.
	// slog is Go 1.21+'s built-in structured logging (replaces log.Printf).
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize application with its dependencies.
	// Note: we're passing a pointer to application, not a copy.
	app := &application{
		logger: logger,
	}

	// http.Server is Go's built-in HTTP server.
	// It handles graceful shutdown, timeouts, and concurrency automatically.
	// Setting Handler to app.routes() (which returns *http.ServeMux)
	// means all routing is delegated to our routes() method.
	server := &http.Server{
		Addr:    *addr,              // Network address to listen on
		Handler: app.routes(),       // Request multiplexer (router)
	}

	// Log before starting - this won't block.
	app.logger.Info("server started", slog.String("addr", *addr))

	// ListenAndServe blocks until the server returns an error.
	// Common errors: address in use, permission denied.
	if err := server.ListenAndServe(); err != nil {
		// Log and exit - Go programs typically use os.Exit(1) for fatal errors.
		// Note: http.ErrServerClosed is expected during graceful shutdown.
		app.logger.Error("server error", slog.Any("err", err))
		os.Exit(1)
	}
}
