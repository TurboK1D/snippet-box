package main

import (
	"embed"
	"html/template"
)

//go:embed ui/html
var htmlFiles embed.FS

//go:embed static
var staticFiles embed.FS

var ts *template.Template

func loadTemplates() error {
	var err error
	ts, err = template.ParseFS(htmlFiles, "ui/html/**/*.html")
	if err != nil {
		return err
	}
	return nil
}
