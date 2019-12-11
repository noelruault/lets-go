package stripe_test

import (
	"fmt"
	"net/http"
	"testing"

	stripe "github.com/noelruault/programming-training/interfaces-training/mocks/v2"
)

func TestAppV1(t *testing.T) {
	client, mux, teardown := stripe.TestClient(t)
	defer teardown()

	mux.HandleFunc("/v1/charges", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":"ch_1DEjEH2eZvKYlo2CxOmkZL4D","amount":2000,
		"description":"Charge for demo purposes.","status":"failed"}`)
	})

	charge, err := client.Charge(123, "doesnt_matter", "something else")
	if err != nil {
		t.Errorf("Charge() err = %s; want nil", err)
	}
	if charge.Status != "succeeded" {
		t.Errorf("Charge() status = %s; want %s", charge.Status, "succeeded")
	}
}
