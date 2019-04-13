package models

import (
	"database/sql"
)

// 1. Declare a Database type (struct in this case)
// 2. Anonymously embed the sql.DB connection pool in our Database struct, so we can
// later access its methods from GetSnippet().
type Database struct {
	*sql.DB // Can be empty if testing database with hard-coded data...
}

// Implement a GetSnippet() method on the Database type. For now, this just returns
// some dummy data, but later we'll update it to query our MySQL database for a
// snippet with a specific ID. In particular, it returns a dummy snippet if the id
// passed to the method equals 123, or returns nil otherwise.
func (db *Database) GetSnippet(id int) (*Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND id = ?` // ? --> placeholder parameter

	// This returns a pointer to a sql.Row object which holds the result returned
	// by the database.
	row := db.QueryRow(stmt, id) // 1. Prepares the statement, 2. Passes parameter, 3. Close

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippet object.
	return s, nil
}

func (db *Database) LatestSnippets() (Snippets, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := Snippets{}
	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

func (db *Database) InsertSnippet(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? SECOND))`

	result, err := db.Exec(stmt, title, content, expires)
	// db.Exec will result sql.Result

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned is of type int64, so we convert it to an int for returning purposes.
	return int(id), nil
}

// NOTE: It's important realize that calls to db.Exec(), db.QueryRow() and db.Query() can use
//any connection from the pool. Even if you have two calls to db.Exec() immediately next to
// each other in your code, there is no guarantee that they will use the same database connection.

// To guarantee that the same connection is used you can wrap multiple statements in a transaction
// tx, err := db.Begin()
