package httpapp

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/jrswab/go-htmx-forge/config"
	"github.com/jrswab/go-htmx-forge/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// @anchor:httpapp:composition request-serving router, middleware policy, routes
// New builds the fully composed application handler.
// staticFS and mediaFS must already be rooted by the caller.
// corsOrigins are exact allowed origins for app routes only; empty disables CORS.
// Invalid or wildcard origins panic (fail closed; same class as config startup reject).
func New(staticFS fs.FS, mediaFS fs.FS, corsOrigins []string) http.Handler {
	return newApp(staticFS, mediaFS, corsOrigins)
}

// newApp is the package-local constructor so tests can register probe routes
// without asserting that the public seam is a *chi.Mux.
func newApp(staticFS fs.FS, mediaFS fs.FS, corsOrigins []string) *chi.Mux {
	if err := config.ValidateCORSOrigins(corsOrigins); err != nil {
		// deliberate: constructor invariant; panic matches config hard-fail at startup
		panic("httpapp.New: " + err.Error())
	}

	r := chi.NewRouter()

	// Global middleware (order owned here).
	r.Use(middleware.Compress(5, "text/html", "text/css", "application/javascript", "application/json", "text/javascript"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS on app routes only. Applied at mux level with path exclusion so
	// OPTIONS preflight is handled (group-scoped CORS misses method-not-allowed).
	// Never pass origins to chi with '*' — chi would treat that as allow-all.
	if len(corsOrigins) > 0 {
		// Copy so callers cannot mutate policy after construction.
		allowed := append([]string(nil), corsOrigins...)
		appCORS := cors.Handler(cors.Options{
			AllowedOrigins: allowed,
			AllowedMethods: []string{"GET", "OPTIONS"},
			AllowedHeaders: []string{
				"Accept",
				"Content-Type",
				"Origin",
				"hx-request",
			},
			// Exact list only; do not fall back to chi's default ["*"].
			AllowOriginFunc: nil,
		})
		r.Use(func(next http.Handler) http.Handler {
			withCORS := appCORS(next)
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if isFileRoute(req.URL.Path) {
					next.ServeHTTP(w, req)
					return
				}
				withCORS.ServeHTTP(w, req)
			})
		})
	}

	r.Get("/", handlers.Home)

	// Static / media: FileServer semantics; no CORS policy.
	r.Handle("/static/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		}
		http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))).ServeHTTP(w, req)
	}))
	r.Handle("/media/*", http.StripPrefix("/media/", http.FileServer(http.FS(mediaFS))))

	return r
}

func isFileRoute(path string) bool {
	return strings.HasPrefix(path, "/static/") || strings.HasPrefix(path, "/media/") ||
		path == "/static" || path == "/media"
}
