package main

import (
	"net/http"

	"github.com/almacitunaberk/goforweb/pkg/helpers"
	"github.com/justinas/nosurf"
)

// Adds CSRFToken for every POST request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// Path "/" means that this cookie is applied to the whole web app
	// Secure is false right now because we are not using HTTPS, BUT IN PRODUCTION, we will use Secure: true

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// It is a middleware to load the created session. Routes use it for every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w,r,"/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w,r)
	})
}