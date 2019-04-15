package main

import (
	"log"
	"net/http"
)

// Our LogRequest middleware is a function that accepts the next handler
// in a chain as a parameter. It executes some logic
// (here, logging the request) and then calls the next handler.
func LogRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		pattern := `%s - "%s %s %s"`
		log.Printf(pattern, r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
	// If you call return before you call next.ServeHTTP()
	// then the chain will stop being executed and control
	// will flow back upstream.
	// (Common use-case for early returns is authentication middleware)
}

// SIMPLER LogRequest:
// return http.HandlerFunc(
// 	func(w http.ResponseWriter, r *http.Request) {
// 		pattern := `%s - "%s %s %s"`
// 		log.Printf(pattern, r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
// 		next.ServeHTTP(w, r)
// 	})

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
		next.ServeHTTP(w, r)
	})
}

func (app *App) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the app.LoggedIn() helper to get the status for the current user.
		loggedIn, err := app.LoggedIn(r)
		if err != nil {
			app.ServerError(w, err)
			return
		}
		// If they are not logged in, redirect them to the login page and return
		// from the middleware chain so that no subsequent handlers in the chain
		// are executed.
		if !loggedIn {
			http.Redirect(w, r, "/user/login", 302)
			return
		}
		// Otherwise call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
