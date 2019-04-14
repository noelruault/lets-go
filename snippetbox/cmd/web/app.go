package main

import (
	"github.com/alexedwards/scs"
	"github.com/noelruault/lets-go/snippetbox/pkg/models"
)

// Define an App struct to hold the application-wide dependencies and configuration
// settings for our web application. For now we'll only include a HTMLDir field
// for the path to the HTML templates directory, but we'll add more to it as our
// build progresses.
type App struct {
	Addr      string // Add an Addr field
	Database  *models.Database
	HTMLDir   string
	Sessions  *scs.Manager
	StaticDir string
	TLSCert   string // Add a TLSCert field
	TLSKey    string // Add a TLSKey field
}
