package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
	_ "github.com/go-sql-driver/mysql" // main.go doesn't actually use anything in the mysql package
	"github.com/noelruault/lets-go/snippetbox/pkg/models"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "sb:pass@/snippetbox?parseTime=true", "MySQL DSN")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	secret := flag.String("secret", "s6Nd%+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")
	tlsCert := flag.String("tls-cert", "./tls/cert.pem", "Path to TLS certificate")
	tlsKey := flag.String("tls-key", "./tls/key.pem", "Path to TLS key")

	flag.Parse()

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate connect() function below. We pass connect() the DSN
	// from the command-line flag.
	db := connect(*dsn)
	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	// ... Our application is only ever terminated by a signal interrupt
	//(i.e. Ctrl+c) or by log.Fatal
	defer db.Close()

	// Use the scs.NewCookieManager() function to initialize a new session manager,
	// passing in the secret key as the parameter. Then we configure it so the
	// session always expires after 12 hours and sessions are persisted across
	// browser restarts.
	sessionManager := scs.NewCookieManager(*secret)
	sessionManager.Lifetime(12 * time.Hour)
	sessionManager.Persist(true)
	sessionManager.Secure(true) // Set the Secure flag on our session cookies
	// ... other methods: https://godoc.org/github.com/alexedwards/scs#pkg-index

	app := &App{
		// Pass in the connection pool when initializing the models.Database object.
		Database:  &models.Database{db}, // {}
		HTMLDir:   *htmlDir,
		Sessions:  sessionManager,
		StaticDir: *staticDir,
	}

	// Pass the app.Routes() method (which returns a serve mux) to the
	// http.ListenAndServe() function.
	log.Printf("Starting server on %s", *addr)
	// err := http.ListenAndServe(*addr, app.Routes())
	err := http.ListenAndServeTLS(*addr, *tlsCert, *tlsKey, app.Routes()) // Start the HTTPS server.
	log.Fatal(err)
}

// The connect() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func connect(dsn string) *sql.DB {

	// sql.Open() it's a pool of many connections. Go manages these connections
	// as needed, automatically opening and closing connections to the database
	// via the driver.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// connections to the database are established lazily, as and when needed
	// for the first time. So to verify that everything is set up correctly
	// we use the db.Ping()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
