# V1/V2 Documentation

Both approaches mock the server using a `TestClient` helper.

`mocks/v1/clientVX_test.go`

```go
client, mux, teardown := stripe.TestClient(t)
```

`mocks/v1/testing.go`

```go
func TestClient(t *testing.T) (*Client, *http.ServeMux, func()) {
    mux := http.NewServeMux()
    server := httptest.NewServer(mux)
    c := &Client{
        baseURL: server.URL,
    }
    return c, mux, func() {
        server.Close()
    }
}
```

The first version, does not use dependency injection, just modifies the client mux with the desired data:

```go
mux.HandleFunc("/v1/charges", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, `{"id":"ch_1DEjEH2eZvKYlo2CxOmkZL4D","amount":2000,
    "description":"Charge for demo purposes.","status":"failed"}`)
})
```

The second one, in addition to modifying the mux that is acting as the server, also injects the client into the app and runs the tests. They will use the local test server using the mux.

```go
//...

mux.HandleFunc("/v1/charges", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, `{"id":"ch_1DEjEH2eZvKYlo2CxOmkZL4D","amount":2000,
    "description":"Charge for demo purposes.","status":"failed"}`)
})

app := App{
    Stripe: client, // injects client
}
app.Run()

charge, err := app.Stripe.Charge(123, "doesnt_matter", "something else")

// ...
```
