package main

import (
	"embed"
	"html/template"
)

// Go Embed Pattern (Go 1.16+):
//
// The //go:embed directive tells the Go compiler to include files or directories
// in the application binary. This is useful for:
//   - Static assets (CSS, JS, images)
//   - HTML templates
//   - Configuration files
//   - Any files that should be part of the binary
//
// Benefits:
//   - Single binary deployment (no separate static files)
//   - No file system dependencies
//   - Easier distribution
//
// The embed.FS type is a read-only file system exposed by the compiler.

//go:embed ui/html
// htmlFiles contains all HTML template files from the ui/html directory.
// The compiler reads these files at build time and includes them in the binary.
var htmlFiles embed.FS

//go:embed static
// staticFiles contains static assets (CSS, JS, images).
// http.FileServer(http.FS(staticFiles)) serves these files.
var staticFiles embed.FS

// ts (template set) is a parsed collection of templates.
// Using a package-level variable avoids re-parsing on every request.
// In production, you might want to re-parse on startup and cache.
var ts *template.Template

// loadTemplates reads all HTML files from the embedded file system.
// template.ParseFS walks the directory tree, parsing any .html files.
// The second argument is a glob pattern: "ui/html/**/*.html" means
// all .html files in ui/html and its subdirectories.
func loadTemplates() error {
	var err error
	ts, err = template.ParseFS(htmlFiles, "ui/html/**/*.html")
	if err != nil {
		return err
	}
	return nil
}
