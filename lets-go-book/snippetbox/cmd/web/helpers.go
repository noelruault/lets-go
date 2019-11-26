package main

import (
	"net/http"
)

func (app *App) LoggedIn(r *http.Request) (bool, error) {
	// Load the session data for the current request, and use the Exists() method
	// to check if it contains a currentUserID key. This returns true if the
	// key is in the session data; false otherwise.
	session := app.Sessions.Load(r)
	loggedIn, err := session.Exists("currentUserID")
	if err != nil {
		return false, err
	}
	return loggedIn, nil
}
