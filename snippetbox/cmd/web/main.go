package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // main.go doesn't actually use anything in the mysql package
	"github.com/noelruault/lets-go/snippetbox/pkg/models"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "sb:pass@/snippetbox?parseTime=true", "MySQL DSN")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")
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

	app := &App{
		// Pass in the connection pool when initializing the models.Database object.
		Database:  &models.Database{db}, // {}
		HTMLDir:   *htmlDir,
		StaticDir: *staticDir,
	}

	// Pass the app.Routes() method (which returns a serve mux) to the
	// http.ListenAndServe() function.
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, app.Routes())
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
