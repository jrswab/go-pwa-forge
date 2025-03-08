package handlers

import "net/http"

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
