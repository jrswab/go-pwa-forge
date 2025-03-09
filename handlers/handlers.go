package handlers

import (
	"net/http"

	"github.com/jrswab/go-htmx-forge/views/home"
)

type HomeHandler struct{}

func (h HomeHandler) Load(w http.ResponseWriter, r *http.Request) {
	render(w, r, home.Show(), http.StatusOK)
}
