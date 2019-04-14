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

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
		next.ServeHTTP(w, r)
	})
}

// SIMPLER LogRequest:
// return http.HandlerFunc(
// 	func(w http.ResponseWriter, r *http.Request) {
// 		pattern := `%s - "%s %s %s"`
// 		log.Printf(pattern, r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
// 		next.ServeHTTP(w, r)
// 	})
