package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

// Change the signature so we're returning a http.Handler instead of a
// *http.ServeMux.
func (app *App) Routes() http.Handler {
	mux := pat.New()

	// The order of the handler calls matters.
	mux.Get("/", http.HandlerFunc(app.Home))
	mux.Get("/snippet/new", app.RequireLogin(http.HandlerFunc(app.NewSnippet)))
	mux.Post("/snippet/new", app.RequireLogin(http.HandlerFunc(app.CreateSnippet)))
	mux.Get("/snippet/:id", http.HandlerFunc(app.ShowSnippet))
	mux.Get("/user/signup", http.HandlerFunc(app.SignupUser))
	mux.Post("/user/signup", http.HandlerFunc(app.CreateUser))
	mux.Get("/user/login", http.HandlerFunc(app.LoginUser))
	mux.Post("/user/login", http.HandlerFunc(app.VerifyUser))
	mux.Post("/user/logout", app.RequireLogin(http.HandlerFunc(app.LogoutUser)))

	fileServer := http.FileServer(http.Dir(app.StaticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	//return LogRequest(mux) // LogRequest → Router → Application Handler
	return LogRequest(SecureHeaders(mux)) // LogRequest ↔ SecureHeaders ↔ Router ↔ Application Handler
}
