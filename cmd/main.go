package main

import (
	"flag"
	"log"
	"os"

	"github.com/jrswab/go-htmx-forge/config"
	"github.com/jrswab/go-htmx-forge/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	app := echo.New()
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO: Change site name to match
		AllowOrigins: []string{"https://hda.site.com", "https://app.site.com"},
		AllowMethods: []string{"GET", "OPTIONS"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			"hx-request",
		},
	}))

	// Load directory paths:
	app.Static("/static", "static")
	app.Static("/media", "media")

	// Add handlers here as needed
	h := handlers.HomeHandler{}
	app.GET("/", h.Load)

	err = app.Start(port)
	if err != nil {
		log.Fatal(err)
	}
}
