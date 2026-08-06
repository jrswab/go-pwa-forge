package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jrswab/go-htmx-forge/config"
	"github.com/jrswab/go-htmx-forge/internal/httpapp"
)

// To include dot (hidden) files use static/*
//
//go:embed static
var static embed.FS

// @anchor:main:lifecycle flags, config, logging, FS rooting, listener
func main() {
	var (
		port    = ":80"
		isLocal = flag.Bool("l", false, "set local development variables")
	)

	flag.Parse()

	if *isLocal {
		port = ":3000"
	}

	file, err := openLogFile("logs/site.log")
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	// Set the output of the log package to the log file
	log.SetOutput(file)

	cfg := new(config.Config)
	if err := cfg.Load(); err != nil {
		log.Fatalf("config: %v", err)
	}

	// Root filesystems outside request-serving composition.
	staticFS, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatal(err)
	}
	mediaFS := os.DirFS("media")

	handler := httpapp.New(staticFS, mediaFS, cfg.CORSAllowedOrigins)

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal(err)
	}
}

// openLogFile creates the parent directory if needed, then opens/creates the log file.
func openLogFile(path string) (*os.File, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
}
