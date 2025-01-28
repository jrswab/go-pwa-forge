package handlers

import (
	"github.com/jrswab/go-htmx-forge/views/home"
	"github.com/labstack/echo/v4"
)

// Visit https://echo.labstack.com/docs for more information.

// HomeHandler is the method receiver for loading the home page.
type HomeHandler struct{}

// Load renders the home page.
func (h HomeHandler) Load(c echo.Context) error {
	return render(c, home.Show())
}
