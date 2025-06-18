package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jrswab/go-htmx-forge/config"
	"github.com/jrswab/go-htmx-forge/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// To include dot (hidden) files use static/*
//
//go:embed static
var static embed.FS

func main() {
	var (
		port    = ":80"
		isLocal = flag.Bool("l", false, "set local development variables")
	)

	flag.Parse()

	if *isLocal {
		port = ":3000"
	}

	// Open or create a log file
	file, err := os.OpenFile("logs/site.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	// Set the output of the log package to the log file
	log.SetOutput(file)

	cfg := new(config.Config)
	cfg.Load()

	// Initialize Chi router
	r := chi.NewRouter()

	// Brotli Compression:
	c := middleware.NewCompressor(brotli.DefaultCompression, // default level = 6
		"text/html",
		"application/json",
	)

	c.SetEncoder("br", func(w io.Writer, level int) io.Writer {
		// NewWriterV2 currently supports levels 0-7 only.
		if level > 7 {
			level = brotli.DefaultCompression
		}
		return brotli.NewWriterV2(w, level) // v2 encoder
	})

	// Enable brotli compression for use
	r.Use(c.Handler)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://hda.site.com", "https://app.site.com"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"Origin",
			"hx-request",
		},
	}))

	// Additional middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Get paths for all files in the embeded static directory
	staticFS, err := fs.Sub(static, "static")
	if err != nil {
		log.Fatal(err)
	}
	r.Handle("/static/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make sure .js files have a content type of application/javascript
		if strings.HasSuffix(r.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}

		http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))).ServeHTTP(w, r)
	}))

	// Serve media files
	mediaServer := http.FileServer(http.Dir("media"))
	r.Handle("/media/*", http.StripPrefix("/media/", mediaServer))

	// Routes
	h := handlers.HomeHandler{}
	r.Get("/", h.Load)

	// Start server
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
