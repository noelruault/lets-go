package stripe_test

import (
	"fmt"
	"net/http"
	"testing"

	stripe "github.com/noelruault/programming-training/interfaces-training/mocks/v2"
)

type App struct {
	Stripe *stripe.Client
}

func (a *App) Run() {}

func TestApp(t *testing.T) {
	client, mux, teardown := stripe.TestClient(t)
	defer teardown()

	mux.HandleFunc("/v1/charges", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":"ch_1DEjEH2eZvKYlo2CxOmkZL4D","amount":2000,
		"description":"Charge for demo purposes.","status":"failed"}`)
	})

	// Now inject client into your app and run your tests - they will use your
	// local test server using this mux
	app := App{
		Stripe: client,
	}
	app.Run()

	charge, err := app.Stripe.Charge(123, "doesnt_matter", "something else")
	if err != nil {
		t.Errorf("Charge() err = %s; want nil", err)
	}
	if charge.Status != "succeeded" {
		t.Errorf("Charge() status = %s; want %s", charge.Status, "succeeded")
	}
}
