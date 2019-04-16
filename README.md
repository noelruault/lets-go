1. [Organization and structure for a project](#organization-and-structure-for-a-project)
2. [Foundations](#foundations)

    21. [Introduction.](#introduction)

        211. [Super-simple web application with one route.](#super-simple-web-application-with-one-route)
        212. [Handler Functions](#handler-functions)

    22. [Url resolver can be sent to handler](#url-resolver-can-be-sent-to-handler)
    23. [URL Query Strings](#url-query-strings)
    24. [Basic HTML Templates](#basic-html-templates)
        [The http.FileServer Handler](#the-httpfileserver-handler)

3. [Configuration and Error Handling](#configuration-and-error-handling)
    31. [Command-line Flags](#command-line-flags)
    32. [Dependency Injection](#dependency-injection)

4. [Database-Driven Responses](#database-driven-responses)
5. [Dynamic HTML Templates](#dynamic-html-templates)
6. [RESTful Routing](#restful-routing)
7. [Processing Forms](#processing-forms)

# Organization and structure for a project

Before even start a new project, it's a good practice to think about the organization and structure for our project.


For our project we'll use a [popular](https://github.com/thockin/go-build-template) and [tried-and-tested](https://peter.bourgon.org/go-best-practices-2016/#repository-structure) approach which should be a good fit for a wide range of applications. Check also this well-known [proyect-layout](https://github.com/golang-standards/project-layout)

Therefore, the organization of a finished project should be similar to the following:

```
├── README.md
├── cmd
│   └── web
│       ├── app.go
│       ├── errors.go
│       ├── handlers.go
│       ├── helpers.go
│       ├── main.go
│       ├── middleware.go
│       ├── routes.go
│       ├── server.go
│       └── views.go
│
├── test
│   └── *  // https://github.com/golang-standards/project-layout
│
├── pkg
│   ├── forms
│   │   └── forms.go
│   └── models
│       ├── database.go
│       └── models.go
│
├── tls
│   ├── cert.pem    // Only development environment.
│   └── key.pem     // $_
│
└── ui
    ├── html
    │   ├── base.html
    │   ├── homepage.html
    │   ├── loginpage.html
    │   ├── newpage.html
    │   ├── showpage.html
    │   └── signuppage.html
    └── static
        ├── css
        │   └── main.css
        └── img
            ├── favicon.ico
            └── logo.png
```

# Foundations

## Introduction.
### Super-simple web application with one route.

The **first thing** we need is a handler. If you're coming from an MVC-background, you can think of handlers as being a bit like controllers.

```go
import "net/http"
func Home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello from Snippetbox"))
}
```

<details>
<summary>1. `html.ResponseWriter` and `html.Request`</summary>

- The `http.ResponseWriter` parameter provides methods for assembling a HTTP response and sending it to the user. So far we've used its `w.Write()` method to send a byte slice containing `"Hello from Snippetbox"` as the response body.
- The `*http.Request` parameter holds information about the current request, such as the HTTP method and the URL being requested.
</details>

<details>
<summary>2. What is a handler? `ServeHTTP(http.ResponseWriter, *http.Request)` </summary>
Basically to be a handler an object must have a `ServeHTTP()` method with the exact signature:

`ServeHTTP(http.ResponseWriter, *http.Request)`

In it's simplest form a handler might look something like this:

```go
type Home struct {}
func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is my homepage"))
}
```

Here we have an object (in this case it's a Home struct, but it could equally be a string or function or anything else), and we've implemented a method with the signature `ServeHTTP(http.ResponseWriter, *http.Request)` on it. That's all we need to make a handler.

If we set up the handler with a struct like we did, we could then register this with our serve mux like so:

```go
mux := http.NewServeMux()
mux.Handle("/", &Home{})
```

### Handler Functions

In practice it's far more common to write your handlers as a normal function, its less confusing.

```go
func Home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is my homepage"))
}
```

But this Home function is just a normal function; it doesn't have a `ServeHTTP()` method. So in itself it isn't a handler. Instead we need to transform it into a handler using the `http.HandlerFunc()` adapter, like so:

```go
mux := http.NewServeMux()
mux.Handle("/", http.HandlerFunc(Home))
```

The http.HandlerFunc() adapter works by automatically adding a ServeHTTP() method to the Home function. When executed, this ServeHTTP() method then simply calls the content of the original Home function. It's a roundabout but convenient way of coercing a normal function to satisfy the http.Handler interface.

With a bit of magic in between, `Handlefunc()` transforms a function to a handler and registers it in one step. So the code above is functionality equivalent to the code in the following section.

</details>

The **second component** is a router (or serve mux in Go terminology). This stores a mapping between the URL patterns for your application and the corresponding handlers.

```go
mux := http.NewServeMux()
mux.HandleFunc("/", Home)
```

A ServeMux is essentially a HTTP request router (or multiplexor). It compares incoming requests against a list of predefined URL paths, and calls the associated handler for the path whenever a match is found.

**The last thing** we need is a running web server. One of the great things about Go is that you can establish a web server and listen for incoming requests as part of your application itself. You don't need a third-party server like Nginx or Apache.

Maybe it simplifies things to think of the serve mux as just being a special kind of handler, which instead of providing a response itself passes the request on to a second handler.

In fact, what exactly is happening is this: When our server receives a new HTTP request, it calls our serve mux's `ServeHTTP()` method. This looks up the relevant handler based on the request URL path, and in turn calls that handler’s `ServeHTTP()` method.

Basically, you can think of a Go web application as a chain of ServeHTTP() methods being called (concurrently) one after another.

```go
http.ListenAndServe(":4000", mux)
log.Fatal(err)
```


## Url resolver can be sent to handler

[Source]: https://www.alexedwards.net/blog/a-recap-of-request-handling
Processing HTTP requests with Go is primarily about two things: ServeMuxes and Handlers.
We can refactor the given code to use a handler like this:


<details>
<summary>Example lvl 1</summary>

> cmd/web/main.go

```go
package main
import (
    "log"
    "net/http"
)
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", Home) // <-- where the magic happens
    log.Println("Starting server on :4000")
    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
```

> cmd/web/handlers.go

```go
package main
import (
"net/http"
)
func Home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    w.Write([]byte("Hello from Snippetbox"))
}
```

</details>

<details>
<summary>Example lvl 2</summary>

> cmd/web/main.go

```go
package main
import (
    "log"
    "net/http"
)
func main() {
    // Register the two new routes with the serve mux. Notice how we use
    // the http.HandlerFunc() adapter to convert the two functions to handlers.
    mux := http.NewServeMux()
    mux.HandleFunc("/", Home)                   // <-- where the magic happens
    mux.HandleFunc("/snippet", ShowSnippet)     // <-- $_
    mux.HandleFunc("/snippet/new", NewSnippet)  // <-- $_
    log.Println("Starting server on :4000")
    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
```

> cmd/web/handlers.go

```go
package main
import (
    "net/http"
)
func Home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    w.Write([]byte("Hello from Snippetbox"))
}
// Add a placeholder ShowSnippet handler function.
func ShowSnippet(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a specific snippet..."))
}
// Add a placeholder NewSnippet handler function.
func NewSnippet(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display the new snippet form..."))
}
```

</details>

## URL Query Strings

In order to add a bit of dynamic behavior to our application, for example the "/snippet" route accepts a query string parameter, like so:

`/snippet?id=1`

1. Retrieve the value of the id parameter from the URL query string, which we can do using the r.URL.Query().Get() method
2. Because the id parameter is untrusted user input, we must validate the value to make sure it is sane and sensible. In this case we want to check that it contains a natural number.

<details>
<summary>Example</summary>

> cmd/web/handlers.go

```go
package main
import (
    "fmt" // New import
    "net/http"
    "strconv" // New import
)

···

func ShowSnippet(w http.ResponseWriter, r *http.Request) {
    // Extract the value of the id parameter from the query string and try to
    // convert it to an integer using the strconv.Atoi() function. If it couldn't
    // be converted to an integer, or the value is less than 1, we return a 404
    // Not Found response.
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }
    // Use the fmt.Fprintf() function to interpolate the id value with our response
    // and write it to the http.ResponseWriter.
    fmt.Fprintf(w, "Display a specific snippet (ID %d)...", id)
}

···

```

e.g: http://localhost:4000/snippet?id=123

</details>


## Basic HTML Templates

The next step in building our application is to make the Home function render a proper HTML homepage.
To make this happen we'll use Go's [html/template](https://golang.org/pkg/html/template/) package, which provides a family of functions for safely parsing and rendering HTML templates.


<details>
<summary>HTML templates</summary>

> ui/html/base.html

```html
{{define "base"}}
<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <title>{{template "page-title" .}} - Snippetbox</title>
    </head>
    <body>
        <header>
            <h1><a href="/">Snippetbox</a></h1>
        </header>
        <nav>
            <a href="/">Home</a>
            <a href="/snippet/new">New snippet</a>
        </nav>
        <section>
            {{template "page-body" .}}
        </section>
    </body>
</html>
{{end}}
```

> ui/html/homepage.html

```html
{{define "page-title"}}Home{{end}}
{{define "page-body"}}
<h2>Latest Snippets</h2>
<p>There's nothing to see here yet!</p>
{{end}}

```

</details>

<details>
<summary>Handler</summary>

> cmd/web/handlers.go

```go
package main
import (
    "fmt"
    "html/template" // New import
    "log" // New import
    "net/http"
    "strconv"
)

func Home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    // Initialize a slice containing the paths to the two files.
    files := []string{
        "./ui/html/base.html",
        "./ui/html/home.page.html",
    }
    // Use the template.ParseFiles() function to read the files and store the
    // templates in a template set (notice that we can pass the slice of file paths // as a variadic parameter). If there's an error, we log the detailed error
    // message and use the http.Error() function to send a generic 500 Internal
    // Server Error response.
    ts, err := template.ParseFiles(files...)
    if err != nil {
        log.Println(err.Error())
        http.Error(w, "Internal Server Error", 500) return
    }
    // Our template set contains three named templates: base, page-title and // page-body (note that every template in your template set must have a
    // unique name). We use the ExecuteTemplate() method to execute the "base" // template and write its content to our http.RespsonseWriter. The last
    // parameter to ExecuteTemplate() represents any dynamic data that we want to // pass in, which for now we'll leave as nil.
    err = ts.ExecuteTemplate(w, "base", nil)
    if err != nil {
        log.Println(err.Error())
        http.Error(w, "Internal Server Error", 500)
    }
}

···

```

</details>


### The http.FileServer Handler

The key to serving these static files from our web application is the http.FileServer() function. This lets us create a http.FileServer handler which serves files from a specific directory

<details>
<summary>Implementation</summary>

> cmd/web/main.go

```go

func main() {

    mux := http.NewServeMux()
    ...

    // Create a file server which serves files out of the "./ui/static" directory.
    // As before, the path given to the http.Dir function is relative to our project
    // repository root.
    fileServer := http.FileServer(http.Dir("./ui/static"))

    // Use the mux.Handle() function to register the file server as the
    // handler for all URL paths that start with "/static/". For matching
    // paths, we strip the "/static" prefix before the request reaches the file
    // server.
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))

    ...
    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}
```

Using static files:

> ui/html/base.html

```html
{{define "base"}}
<!doctype html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <title>{{template "page-title" .}} - Snippetbox</title>
        <!-- Link to the CSS stylesheet and favicon -->
        <link rel="stylesheet" href="/static/css/main.css">
        <link rel="shortcut icon" href="/static/img/favicon.ico" type="image/x-icon">
    </head>
    <body>
        <header>
            <h1><a href="/">Snippetbox</a></h1>
        </header>
        <nav>
            <a href="/">Home</a>
            <a href="/snippet/new">New snippet</a>
        </nav>
        <section>
            {{template "page-body" .}}
        </section>
    </body>
</html>
{{end}}
```

</details>

# Configuration and Error Handling
## Command-line Flags

In Go, a common and idiomatic way to manage configuration settings is to use command-line flags when starting an application.

`$ go run cmd/web/* -addr=":80" -static-dir="/var/www/static"`

<details>
<summary>How and where configure flags?</summary>

> cmd/web/main.go

```go
package main
import "flag" // New import

func main() {
    // Define command-line flags for the network address and location of the static
    // files directory.
    addr := flag.String("addr", ":4000", "HTTP network address")
    staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")

    // Importantly, we use the flag.Parse() function to parse the command-line flags.
    // This reads in the command-line flag values and assigns them to the addr and
    // staticDir variables. You need to parse the flags *before* you use the addr
    // or staticDir variables, otherwise they will always contain the default value.
    // If any errors are encountered during parsing the application will be
    // terminated.
    flag.Parse()

    ...

    // The value returned from the flag.String() function is a pointer to the flag
    // value, not the value itself. So we need to dereference the pointer (i.e.
    // prefix it with the * symbol) before we use it as the path for our static file
    // server.
    fileServer := http.FileServer(http.Dir(*staticDir)) mux.Handle("/static/", http.StripPrefix("/static", fileServer))
    // Again, we dereference the addr variable and use it as the network address
    // to listen on. Notice that we also use the log.Printf() function to interpolate
    // the correct address in the log message.
    log.Printf("Starting server on %s", *addr)
    err := http.ListenAndServe(*addr, mux)
    log.Fatal(err)
}

```

Use cases:

```bash
$ go run cmd/web/* -addr=":9999"
2017/08/20 13:23:00 Starting server on :9999
```

```bash
$ export SNIPPETBOX_ADDR=":9999"
$ go run cmd/web/* -addr=$SNIPPETBOX_ADDR
2017/08/20 14:53:46 Starting server on :9999
```

A great feature of Go is that you can use the -help flag to list all the available command- line flags for an application and the accompanying help text.

`$ go run cmd/web/* -help`

</details>

## Dependency Injection

In order to **avoid hard-coded location for the HTML templates**, we can easily add a new `-html-dir` command-line flag to our application so that we can configure it at runtime like our other settings.

But this raises a good question: **how can we make the flag value available to our Home function from main()?**

There are a [few different ways](https://www.alexedwards.net/blog/organising-database-access). Injecting dependencies into the handlers makes your code more explicit, less error-prone and easier to unit test than if you use global variables.

For applications where all your code is in the same package, like ours, a neat way to inject dependencies is to put them into a custom App struct:

<details>
<summary>How to inject dependencies into handlers?</summary>

> cmd/web/app.go

```go
package main
// Define an App struct to hold the application-wide dependencies and configuration
// settings for our web application. For now we'll only include a HTMLDir field
// for the path to the HTML templates directory, but we'll add more to it as our
// build progresses.
type App struct {
       HTMLDir string
}
```

> cmd/web/handlers.go

```go
import "path/filepath" // New import

// Change the signature of our handlers so it is defined as a method against // *App.
func (app *App) <HandlerName>(w http.ResponseWriter, r *http.Request) {
    ...

    // Because the Home handler function is a now method against App it can access // its fields. So we can build the paths to the HTML template files using the // HTMLDir value in the App instance.
    files := []string{
        filepath.Join(app.HTMLDir, "base.html"),
        filepath.Join(app.HTMLDir, "home.page.html"),
    }
    ...
}

```

> cmd/web/main.go

```go
func main() {
    // Define a new command-line flag for the path to the HTML template directory.
    htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")

    // Initialize a new instance of App containing the dependencies.
    app := &App{
        HTMLDir: *htmlDir,
    }

    // Swap our route declarations to use the App object's methods as the handler // functions.
    mux := http.NewServeMux()
    mux.HandleFunc("path/to/wherever", app.<HandlerName>)

    ...

```

This pattern that we're using to inject dependencies won't work if your handlers are spread across multiple packages. In that case, an alternative approach is to create a ```config``` package exporting an ```App``` struct and have your handler functions close over this to form a closure. Very roughly:

> cmd/web/main.go

```go
func main() {
    app := &config.App{...}
    mux.Handle("/", handlers.Home(app))
}
```

> cmd/web/handlers.go

```go
func Home(app *config.App) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" { http.NotFound(w, r)
            return
        }
        files := []string{
            filepath.Join(app.HTMLDir, "base.html"),
            filepath.Join(app.HTMLDir, "home.page.html"),
        }
        ···
    }
}
```

You can find a more complete and concrete example of how to use the closure pattern in this [Gist](https://gist.github.com/alexedwards/5cd712192b4831058b21).

</details>


# Database-Driven Responses

To use MySQL from our Go web application we need to [install a database driver](https://github.com/golang/go/wiki/SQLDrivers). This essentially acts as a middleman, translating commands between Go and the actual database itself.


<details>
<summary>Sketch out a database model for our project (get single row).</summary>

> pkg/models/models.go

```go
package models
import "time"

// Define a Snippet type to hold the information about an individual snippet.
type Snippet struct { ID int
       Title   string
       Content string
       Created time.Time
       Expires time.Time
}
// For convenience we also define a Snippets type, which is a slice for holding // multiple Snippet objects.
type Snippets []*Snippet
```


> pkg/models/database.go

```go
package models
import "database/sql"

type Database struct {
    *sql.DB
}

func (db *Database) GetSnippet(id int) (*Snippet, error) {
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is why it's surrounded with backticks instead
	// of normal double quotes).
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// Use the QueryRow() method on the embedded connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result returned by the database.
	row := db.QueryRow(stmt, id)
	// Initialize a pointer to a new zeroed Snippet struct.
	s := &Snippet{}
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement. If our query returned no rows, then
	// row.Scan() will return a sql.ErrNoRows error. We check for that and return
	// nil instead of a Snippet object.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
	return nil, nil
	} else if err != nil {
	return nil, err }
	// If everything went OK then return the Snippet object.
	return s, nil
}
```

> cmd/web/main.go

```go
package main
import (
	...
	"database/sql" // New import
	"</path/to/project>/pkg/models" // New import
	_ "github.com/go-sql-driver/mysql" // New import
)
func main() {
	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "sb:pass@/snippetbox?parseTime=true", "MySQL DSN")
	...

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate connect() function below. We pass connect() the DSN
	// from the command-line flag.
	db := connect(*dsn)
	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	app := &App{
		Database: &models.Database{},
		HTMLDir: *htmlDir,
		StaticDir: *staticDir,
	}

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, app.Routes()) log.Fatal(err)
}

// The connect() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func connect(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
			log.Fatal(err)
		}
	if err := db.Ping(); err != nil {
			log.Fatal(err)
		}
	return db
}
```


> cmd/web/app.go

```go
package main
import "snippetbox.org/pkg/models"

type App struct {
    Database *models.Database HTMLDir string
    StaticDir string
}
```



> cmd/web/handlers.go

```go
···
func (app *App) ShowSnippet(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.NotFound(w)
        return
    }
    snippet, err := app.Database.GetSnippet(id) if err != nil {
        app.ServerError(w, err)
        return
    }
    if snippet == nil {
        app.NotFound(w)
        return
    }
    fmt.Fprint(w, snippet)
}
···
```

</details>

<details>
<summary>Sketch out a database model for our project (get multiple rows).</summary>

> pkg/models/database.go

```go
 package models
···
func (db *Database) LatestSnippets() (Snippets, error) {
    // Write the SQL statement we want to execute.
    stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

    // Use the QueryRow() method on the embedded connection pool to execute our // SQL statement. This results a sql.Rows resultset containing the result of // our query.
    rows, err := db.Query(stmt)
    if err != nil {
        return nil, err
    }

    // IMPORTANTLY we defer rows.Close() to ensure the sql.Rows resultset is
    // always properly closed before LatestSnippets() returns. Closing a
    // resultset is really important. As long as a resultset is open it will
    // keep the underlying database connection open. So if something goes wrong // in this method and the resultset isn't closed, it can rapidly lead to all // the connections in your pool being used up. Another gotcha is that the
    // defer statement should come *after* you check for an error from
    // db.Query(). Otherwise, if db.Query() returns an error, you'll get a panic // trying to close a nil resultset.
    defer rows.Close()

    // Initialize an empty Snippets object (remember that this is just a slice of // the type []*Snippet).
    snippets := Snippets{}
    // Use rows.Next to iterate through the rows in the resultset. This
    // prepares the first (and then each subsequent) row to be acted on by the // rows.Scan() method. If iteration over all of the rows completes then the // resultset automatically closes itself and frees-up the underlying
    // database connection.

    for rows.Next() {
        // Create a pointer to a new zeroed Snippet object.
        s := &Snippet{}
        // Use rows.Scan() to copy the values from each field in the row to the // new Snippet object that we created. Again, the arguments to row.Scan() // must be pointers to the place you want to copy the data into, and the // number of arguments must be exactly the same as the number of
        // columns returned by your statement.
        err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
        if err != nil {
            return nil, err
        }
        // Append it to the slice of snippets.
        snippets = append(snippets, s)
    }

    // When the rows.Next() loop has finished we call rows.Err() to retrieve any // error that was encountered during the iteration. It's important to
    // call this - don't assume that a successful iteration was completed
    // over the whole resultset.
    if err = rows.Err(); err != nil {
        return nil, err
    }
    // If everything went OK then return the Snippets slice.
    return snippets, nil
}
```


</details>

<details>
<summary>Sketch out a database model for our project (insert single row).</summary>

1. Insert a new record into the snippets table, containing a given title, content and expiry time (in seconds).
2. Return the id for the new record.

> pkg/models/database.go

```go
package models
···
func (db *Database) InsertSnippet(title, content, expires string) (int, error) {
    // Write the SQL statement we want to execute.
    stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? SECOND))`
    // Use the db.Exec() method to execute the statement snippet, passing in values
    // for our (untrusted) title, content and expiry placeholder parameters in
    // exactly the same way that we did with the QueryRow() method. This returns
    // a sql.Result object, which contains some basic information about what
    // happened when the statement was executed.
    result, err := db.Exec(stmt, title, content, expires)
    if err != nil {
        return 0, err
    }
    // Use the LastInsertId() method on the result object to get the ID of our
    // newly inserted record in the snippets table.
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
    // The ID returned is of type int64, so we convert it to an int type for
    // returning from our Insert function.
    return int(id), nil
}
```


</details>

# Dynamic HTML Templates
In this section we're going to concentrate on displaying the dynamic data from our MySQL database in some proper HTML pages.

The html/template package provides [some template functions](https://golang.org/pkg/text/template/#hdr-Functions)


<details>
<summary>Views, handlers and HTML Templates example.</summary>

> cmd/web/handlers.go

```go
package main

...

func (app *App) Home(w http.ResponseWriter, r *http.Request) {
    ...
    // Include the *http.Request parameter.
    app.RenderHTML(w, r, "homepage.html", &HTMLData{
        Snippets: snippets,
    })
    ...
}

func (app *App) ShowSnippet(w http.ResponseWriter, r *http.Request) {
    ...
    // Include the *http.Request parameter.
    app.RenderHTML(w, r, "showpage.html", &HTMLData{
        Snippets: snippets,
    })
    ...
}

...

```

> cmd/web/views.go

```go
package main
···
// Add a Path field to the struct.
type HTMLData struct {
    Path string
    Snippet *models.Snippet
    Snippets []*models.Snippet
}

// Change the signature of the RenderHTML() method so that it accepts *http.Request
// as the second parameter.
func (app *App) RenderHTML(w http.ResponseWriter, r *http.Request, page string, data *HTMLData) {
    // If no data has been passed in, initialize a new empty HTMLData object.
    if data == nil {
        data = &HTMLData{}
    }
    // Add the current request URL path to the data.
    data.Path = r.URL.Path

    files := []string{
        filepath.Join(app.HTMLDir, "base.html"),
        filepath.Join(app.HTMLDir, page),
    }
    funcs := template.FuncMap{
        "humanDate": humanDate,
    }
    ts, err := template.New("").Funcs(funcs).ParseFiles(files...)

    if err != nil {
        app.ServerError(w, err)
        return
    }

    // Initialize a new buffer.
    buf := new(bytes.Buffer)

    // Write the template to the buffer, instead of straight to the
    // http.ResponseWriter. If there's an error, call our error handler and then return.
    err = ts.ExecuteTemplate(buf, "base", data)
    if err != nil {
        app.ServerError(w, err)
        return
    }

    // Write the contents of the buffer to the http.ResponseWriter. Again, this
    // is another time where we pass our http.ResponseWriter to a function that
    // takes an io.Writer.
    buf.WriteTo(w)
}

```

> ui/html/base.html

```html
{{define "base"}}
<!doctype html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <title>{{template "page-title" .}} - Snippetbox</title>
    <link rel="stylesheet" href="/static/css/main.css">
    <link rel="shortcut icon" href="/static/img/favicon.ico" type="image/x-icon">
</head>

<body>
    <header>
        <h1><a href="/">Snippetbox</a></h1>
    </header>
    <nav>
        <a href="/" {{if eq .Path "/" }} class="live" {{end}}> Home
        </a>
        <a href="/snippet/new" {{if eq .Path "/snippet/new" }} class="live" {{end}}>
            New snippet
        </a> </nav>
    <section>
        {{template "page-body" .}}
    </section>
</body>

</html>
{{end}}
```

> ui/html/showpage.html

```html
{{define "page-title"}}Snippet #{{.Snippet.ID}}{{end}}
{{define "page-body"}} {{with .Snippet}}
<div class="snippet">
    <div class="metadata">
        <strong>{{.Title}}</strong>
        <span>#{{.ID}}</span>
    </div>
    <pre><code>{{.Content}}</code></pre>
    <div class="metadata">
        <time>Created: {{humanDate .Created}}</time>
        <time>Expires: {{humanDate .Expires}}</time> </div>
</div>
{{end}}
{{end}}
```

> ui/html/homepage.html

```html
{{define "page-title"}}Home{{end}}
{{define "page-body"}} <h2>Latest Snippets</h2> {{if .Snippets}}
<table>
    <tr>
        <th>Title</th>
        <th>Created</th>
        <th>ID</th>
    </tr>
    {{range .Snippets}}
    <tr>
        <td><a href="/snippet?id={{.ID}}">{{.Title}}</a></td>
        <td>{{humanDate .Created}}</td>
        <td>#{{.ID}}</td>
    </tr>
    {{end}}
</table>
{{else}}
<p>There's nothing to see here yet!</p>
{{end}}
{{end}}
```




</details>



# RESTful Routing

Port our application from using Go's inbuilt serve mux to using Pat.

The basic syntax for creating a router and registering a route with Pat looks like this:

```go
mux := pat.New()
mux.Get("/snippet/:id", http.HandlerFunc(app.ShowSnippet))
```

- The "/snippet/:id" pattern includes a named capture :id.
- We use the `mux.Get()` method to register a URL pattern and handler which will be called only if the request has a `GET` HTTP method.
- Get(), Post(), Put(), Delete() and other methods are provided.
- Pat doesn't allow us to register handler functions directly, so we need to convert them using the `http.HandlerFunc()` adapter – just like we were doing at the start of the book.

<details>
<summary>Let's head over to the routes.go file.</summary>

> cmd/web/routes.go

```go
package main import (
    "net/http"
    "github.com/bmizerany/pat" // New import
)
// Change the signature so we're returning a http.Handler instead of a
// *http.ServeMux.
func (app *App) Routes() http.Handler {
    mux := pat.New()
    mux.Get("/", http.HandlerFunc(app.Home))
    mux.Get("/snippet/new", http.HandlerFunc(app.NewSnippet)) mux.Post("/snippet/new", http.HandlerFunc(app.CreateSnippet)) mux.Get("/snippet/:id", http.HandlerFunc(app.ShowSnippet)) // Moved downwards
    fileServer := http.FileServer(http.Dir(app.StaticDir)) mux.Get("/static/", http.StripPrefix("/static", fileServer))
    return mux
}
```

- Pat matches patterns in the order that they were registered.
- Refered to URL patterns which end in a trailing slash (like `"/static/"` in our code above). Any request which matches the start of the pattern will be dispatched to the corresponding handler.
- The pattern `"/"` is a special case. It will only match requests where the URL path is exactly `"/"`.


> cmd/web/handlers.go

```go
// Because Pat matches the "/" path exactly, we can now remove the manual check
// of r.URL.Path != "/" from the Home function.
func (app *App) Home(w http.ResponseWriter, r *http.Request) {
    ...
}

func (app *App) ShowSnippet(w http.ResponseWriter, r *http.Request) {
    // Pat doesn't strip the colon from the named capture key, so we need to
    // get the value of ":id" from the query string instead of "id".
    id, err := strconv.Atoi(r.URL.Query().Get(":id"))
    if err != nil || id < 1 {
        app.NotFound(w)
        return
    }
    ...
}
```

> ui/html/homepage.html

```html
{{define "page-title"}}Home{{end}}
{{define "page-body"}} <h2>Latest Snippets</h2>
{{if .Snippets}}
    <table>
        <tr>
            <th>...</th>
        </tr>
        {{range .Snippets}}
        <tr>
            <!-- Interpolate the ID, instead of appending it in a query string. -->
            <td><a href="/snippet/{{.ID}}">{{.Title}}</a></td>
            <td>...</td>
        </tr>
        {{end}}
    </table>
{{else}}
    <p>There's nothing to see here yet!</p>
{{end}}
{{end}}
```


</details>

# Processing Forms

Validate the form data when it's submitted.


1. We first need to use the [r.ParseForm()](https://golang.org/pkg/net/http/#Request.ParseForm) method to parse the request body. This checks that the request body is well-formed, and stores the form data in the request's [r.PostForm](https://golang.org/pkg/net/http/#Request) map. If there are any errors encountered when parsing the body (like there is no body, or it's too large to process) then it will return an error. The `r.ParseForm()` method is also idempotent; it can safely be called multiple times on the same request without any side-effects.

2. We can then get to the form data contained in `r.PostForm` by using the `r.PostForm.Get()` method. For example, we can retrieve the value of the `title` field with `r.PostForm.Get("title")`. If there is no matching field name in the form this will return the empty string "". This is similar to the way that query string parameters worked earlier in the book.

3. We can then validate the individual form values using the various functions in the [strings](https://golang.org/pkg/strings/) and [unicode/utf8](https://golang.org/pkg/unicode/utf8/) packages.


<details>
<summary>Validation</summary>

> pkg/forms/forms.go

```go
package forms
// Declare a struct to hold the form values (and also a map to hold any validation
// failure messages).
type NewSnippet struct {
    Title string
    Content string
    Expires string
    Failures map[string]string
}
// Implement an Valid() method which carries out validation checks on the form
// fields and returns true if there are no failures.
func (f *NewSnippet) Valid() bool {
    f.Failures = make(map[string]string)
    // We will validate the form fields here...
    return len(f.Failures) == 0
}
```

> cmd/web/handlers.go

```go
package main
import (
    "fmt" // New import
    "net/http"
    "strconv"

    "snippetbox.org/pkg/forms" // New import
)
···
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
        Title: r.PostForm.Get("title"),
        Content: r.PostForm.Get("content"),
        Expires: r.PostForm.Get("expires"),
    }
    // Check if the form passes the validation checks. If not, then use the
    // fmt.Fprint function to dump the failure messages to the response body.
    if !form.Valid() {
        fmt.Fprint(w, form.Failures)
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
    // If successful, send a 303 See Other response redirecting the user to the // page with their new snippet.
    http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
```

So that's the basic pattern, but it's not much use without any actual validation rules. Let's head back to pkg/forms/forms.go and add some validation rules:

> pkg/forms/forms.go

```go
package forms
import (
    "strings" // New import
    "unicode/utf8" // New import
)
type NewSnippet struct {
    Title string
    Content string
    Expires string
    Failures map[string]string
}
func (f *NewSnippet) Valid() bool { f.Failures = make(map[string]string)
    // Check that the Title field is not blank and is not more than 100 characters
    // long. If it fails either of those checks, add a message to the f.Failures
    // map using the field name as the key.
    if strings.TrimSpace(f.Title) == "" {
        f.Failures["Title"] = "Title is required"
    } else if utf8.RuneCountInString(f.Title) > 100 {
        f.Failures["Title"] = "Title cannot be longer than 100 characters"
    }

    // Validate the Content and Expires fields aren't blank in a similar way.
    if strings.TrimSpace(f.Content) == "" {
        f.Failures["Content"] = "Content is required"
    }

    // Check that the Expires field isn't blank and is one of a fixed list. Using
    // a lookup on a map keyed with the permitted options and values of true is a
    // neat trick which saves you looping over the permitted values.
    permitted := map[string]bool{"3600": true, "86400": true, "31536000": true}
    if strings.TrimSpace(f.Expires) == "" {
        f.Failures["Expires"] = "Expiry time is required"
    } else if !permitted[f.Expires] {
        .Failures["Expires"] = "Expiry time must be 3600, 86400 or 31536000 seconds"
    }

    // If there are no failure messages, return true.
    return len(f.Failures) == 0
}
```

Find code patterns for processing and validating different types of inputs in [this blog post](http://www.alexedwards.net/blog/validation-snippets-for-go).

</details>


<details>
<summary>Validation failures</summary>

> cmd/web/main.go

```go
package main
···
// Add a Form field to the struct.
type HTMLData struct {
    Form interface{}
    Path string
    Snippet *models.Snippet
    Snippets []*models.Snippet
}
···
```

> cmd/web/handlers.go

```go
package main
···
func (app *App) NewSnippet(w http.ResponseWriter, r *http.Request) {
    // Pass an empty *forms.NewSnippet object to the new.page.html template. Because
    // it's empty, it won't contain any previously submitted data or validation
    // failure messages.
    app.RenderHTML(w, r, "new.page.html", &HTMLData{
        Form: &forms.NewSnippet{},
    })
}

func (app *App) CreateSnippet(w http.ResponseWriter, r *http.Request) { err := r.ParseForm()
    if err != nil {
        app.ClientError(w, http.StatusBadRequest)
        return
    }
    form := &forms.NewSnippet{
        Title: r.PostForm.Get("title"),
        Content: r.PostForm.Get("content"),
        Expires: r.PostForm.Get("expires"),
    }
    if !form.Valid() {
    // Re-display the new.page.html template passing in the *forms.NewSnippet
    // object (which contains the validation failure messages and previously
    // submitted data).
        app.RenderHTML(w, r, "new.page.html", &HTMLData{Form: form})
        return
    }
    id, err := app.Database.InsertSnippet(form.Title, form.Content, form.Expires)
    if err != nil {
        app.ServerError(w, err)
        return
    }
    http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
```

> ui/html/newpage.html

```go
{{define "page-title"}}Add a New Snippet{{end}}
{{define "page-body"}}
<form action="/snippet/new" method="POST"> {{with .Form}}
    <div>
        <label>Title:</label> {{with .Failures.Title}}
        <label class="error">{{.}}</label> {{end}}
        <input type="text" name="title" value="{{.Title}}"> </div>
    <div>
        <label>Content:</label> {{with .Failures.Content}}
        <label class="error">{{.}}</label> {{end}}
        <textarea name="content">{{.Content}}</textarea> </div>
    <div>
        <label>Delete in:</label> {{with .Failures.Expires}}
        <label class="error">{{.}}</label> {{end}}
        {{$expires := or .Expires "31536000"}}
        <input type="radio" name="expires" value="31536000" {{if (eq $expires "31536000")}} checked{{end}}> One Year
        <input type="radio" name="expires" value="86400" {{if (eq $expires "86400")}} checked{{end}}> One Day
        <input type="radio" name="expires" value="3600" {{if (eq $expires "3600")}} checked{{end}}> One Hour
    </div>
    <div>
        <input type="submit" value="Publish snippet"> </div>
    {{end}}
</form>
{{end}}

```

</details>