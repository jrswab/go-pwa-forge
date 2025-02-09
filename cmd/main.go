package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jrswab/go-htmx-forge/config"
	"github.com/jrswab/go-htmx-forge/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

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

	// Static file servers
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

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
