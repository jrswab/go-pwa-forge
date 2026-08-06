package handlers

import (
	"net/http"

	"github.com/jrswab/go-htmx-forge/views/home"
)

// Home serves GET / with the home template.
func Home(w http.ResponseWriter, r *http.Request) {
	Render(w, r, home.Show(), http.StatusOK)
}
