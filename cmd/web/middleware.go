package main

import (
	"net/http"

	"github.com/aidisapp/musiqcity_v2/internal/helpers"
	"github.com/justinas/nosurf"
)

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	// LoadAndSave provides middleware which automatically loads and saves session data for the current request, and communicates the session token to and from the client in a cookie.
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAdmin(r) {
			session.Put(r.Context(), "error", "You do not have access to this page!!!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
