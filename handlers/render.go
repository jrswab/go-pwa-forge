package handlers

import (
	"net/http"

	"github.com/a-h/templ"
)

// Render writes Content-Type and status then the Templ component. On render
// failure it emits an internal-server-error via http.Error (headers may
// already have been sent; the intended status remains committed).
//
// Content-Type is set before WriteHeader so standard Chi compression can
// classify the response as compressible text/html.
func Render(w http.ResponseWriter, r *http.Request, component templ.Component, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
