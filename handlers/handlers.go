package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/jrswab/go-htmx-forge/views/home"
)

func render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type HomeHandler struct{}

func (h HomeHandler) Load(w http.ResponseWriter, r *http.Request) {
	render(w, r, home.Show())
}
