package main

import (
	"net/http" // Debug purposes
	"strconv"
)

func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}

	// fmt.Printf("[DEBUG]: app.HTMLDir = " + app.HTMLDir + "\n") // Renders path to file
	// app.RenderHTML(w, "homepage.html", nil) // Use the app.RenderHTML() helper.

	// Fetch a slice of the latest snippets from the database.
	snippets, err := app.Database.LatestSnippets()
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// Pass the slice of snippets to the "home.page.html" templates.
	// Include the *http.Request parameter.
	app.RenderHTML(w, r, "homepage.html", &HTMLData{
		Snippets: snippets,
	})
}

func (app *App) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}
	snippet, err := app.Database.GetSnippet(id)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	if snippet == nil {
		app.NotFound(w)
		return
	}
	// Render the showpage.html template, passing in the snippet data wrapped in
	// our HTMLData struct.
	// Include the *http.Request parameter.
	app.RenderHTML(w, r, "homepage.html", &HTMLData{
		Snippet: snippet,
	})
}

func (app *App) NewSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the new snippet form..."))
}
