package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/noelruault/lets-go/snippetbox/pkg/forms"
)

func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	// Fetch a slice of the latest snippets from the database.
	snippets, err := app.Database.LatestSnippets()
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// Pass the slice of snippets to the "homepage.html" templates.
	// Include the *http.Request parameter.
	app.RenderHTML(w, r, "homepage.html", &HTMLData{
		Snippets: snippets,
	})
}

func (app *App) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id".
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
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

	session := app.Sessions.Load(r)
	flash, err := session.PopString(w, "flash") // PopString will delete flash after reading
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Render the showpage.html template, passing in the snippet data
	// wrapped in our HTMLData struct.
	// Include the *http.Request parameter. (r)
	app.RenderHTML(w, r, "showpage.html", &HTMLData{
		Flash:   flash, // Pass the flash message to the template.
		Snippet: snippet,
	})
}

func (app *App) NewSnippet(w http.ResponseWriter, r *http.Request) {
	// Pass an empty *forms.NewSnippet object to the newpage.html template. Because
	// it's empty, it won't contain any previously submitted data or validation
	// failure messages.
	app.RenderHTML(w, r, "newpage.html", &HTMLData{
		Form: &forms.NewSnippet{},
	})
}

func (app *App) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	// First we call r.ParseForm() which adds any POST (also PUT and PATCH) data
	// to the r.PostForm map. If there are any errors we use our
	// app.ClientError helper to send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}
	// We initialize a *forms.NewSnippet object and use the r.PostForm.Get() method
	// to assign the data to the relevant fields.
	form := &forms.NewSnippet{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: r.PostForm.Get("expires"),
	}
	// Check if the form passes the validation checks. If not, then use the
	// fmt.Fprint function to dump the failure messages to the response body.
	if !form.Valid() {
		// fmt.Fprint(w, form.Failures) // Will break the page. Rendering raw html.
		app.RenderHTML(w, r, "newpage.html", &HTMLData{Form: form})
		return
	}
	// If the validation checks have been passed, call our database model's
	// InsertSnippet() method to create a new database record and return it's ID
	// value.
	id, err := app.Database.InsertSnippet(form.Title, form.Content, form.Expires)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	session := app.Sessions.Load(r)
	err = session.PutString(w, "flash", "Your snippet was saved successfully!")
	// ...other methods than PutString: https://godoc.org/github.com/alexedwards/scs#pkg-index
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// If successful, send a 303 See Other response redirecting the user to the
	// page with their new snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
