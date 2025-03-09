package handlers

import (
	"net/http"

	"github.com/a-h/templ"
)

// render takes in a Templ component and returns the HTML to the client.
func render(w http.ResponseWriter, r *http.Request, component templ.Component, status int) {
	w.WriteHeader(status)

	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// setCookie creates a cookie which does not expire.
// SameSite is set to strict to avoid cookies being passod from other domains.
// Switch to SameSiteLaxMode for more flexibility.
//
// See https://datatracker.ietf.org/doc/html/draft-ietf-httpbis-cookie-same-site-00 for more information.
func setPersistentCookie(w http.ResponseWriter, key, val string) {
	cookie := http.Cookie{
		Name:     key,
		Value:    val,
		Path:     "/",
		Domain:   "",   // Leave empty for current domain
		Secure:   true, // Set to true for HTTPS only
		HttpOnly: true, // Cannot be accessed by JavaScript
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)
}
