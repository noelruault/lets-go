
0. [Learning about Golang Interfaces, DI and mocking](#learning-golang-interfaces)
1. [Structs and Interfaces](#structs-and-interfaces)
2. [Dependency injection](#dependency-injection)
3. [Mocks](#mocks)

    31. [Mocks: Faking API's](#mocks-faking-apis)

        311. [1. Replacing the server](#1-replacing-the-server)
        312. [2. Replacing the client](#2-replacing-the-client)

# Learning about Golang Interfaces, dependency injection and mocking

Short code-guide to refresh useful stuff about interfaces, dependency injection and mocking.

## Structs and Interfaces

Go's build-in data types. A type that contains named fields can be defined with a struct. For example, we can represent a circle or a Rectangle like this:

```go
    type Circle struct {
        radius, diagonal float64
    }

    type Rectangle struct {
        length, width float64
    }
```

Although we can implement a first version of "calculating area functions" like next:

```go
    func circleArea(radius, diagonal float64) float64 {
        return math.Pi * radius * radius
    }

    func rectangleArea(length, width float64) float64 {
        return length * width
    }
```

We can improve it with an special type of function known as methods. We can use receivers, for different structs to define functionality specific to each type. We no longer need the `&` operator (Go automatically knows to pass a pointer to the circle for this method), and because this function can only be used with specific types such as Circle or Rectangle, we can rename the function to just `area`

```go
    func (c *Circle) area() float64 {
        return math.Pi * c.radius * c.radius
    }

    func (r *Rectangle) area() float64 {
        return r.width * r.length
    }
```

We can embed structs inside another structs.

```go
    type Combo struct {
        circle      Circle      // Named embedding      -> Combo.circle.area
        rectangle   Rectangle   //                      -> Combo.rectangle.area
    }

    type JustRectangle struct {
        Rectangle               // Unnamed embedding    -> JustRectangle.area
    }
```

Because both of our figures have an area method, can be implemented by an interface...
In case that we would like to design a function that calculates the area of multiple variadic parameters for all our figures, we would end-up doing something like this:

```go
    // THIS IS INVALID
    func totalArea(circles ...Circle, rectangles ...Rectangle) float64{}

    func totalArea(circles []Circle, rectangles []Rectangle) float64{}
```

But this is not useful speaking of scalability.

Interfaces were designed to solve this problem. Because both of our figures have an area method, they both implement the Shape interface and we can code a function like this:

```go
    type Shape interface {
        area() float64
        // With a Shape instance, we wouldn't be able to access the struct fields
        // for the figure. Only the area.
    }

    func totalArea(shapes ...Shape) float64 {
        var area float64
        for _, s := range shapes {
            area += s.area()
        }
        return area
    }
```

This allows us to manage multiple types that implement the same method. In this case the "area" method can be used to get the expected result for different types of figures.

```go
    func main () {
        // ...
        circle = Circle{5, 0}
        rectangle = Rectangle{3, 5}
        fmt.Println(totalArea(&circle, &rectangle))
        // ...
    }
```

In Go, **Interfaces define functionality rather than data** so interfaces can also be used as fields...

```go
    type MultiShape struct {
        shapes []Shape
    }

    func (m *MultiShape) area() float64 {
        var area float64
        for _, s := range m.shapes {
            area += s.area()
        }
        return area
    }

    func main () {
        // ...
        // Here, Shape interface is being used as a field.
        multiShape := MultiShape{
            shapes: []Shape{
                &Circle{5, 0},
                &Rectangle{3, 5},
            },
        }

        fmt.Println(multiShape.area())
    }
```

Now a MultiShape can contain any figure.
Interfaces are perticulary useful as software projects grow and become more complex. They allow us to hide the incidental details of implementation (e.g. the fields of our struct), which makes it easeier to reason about software components in isolation.

That can be applied to this example because as long as the area methods we defined continue to produce the same results, we are free to change how a Circle or Rectangle is structured without having to worry about the correct functionality of the totalArea function.

## Dependency injection

Dependency injection enables us to write implementation agnostic code, to write tests easily by simulating specific behaviour or helps us to remove global state.
DI is done by using variables or more commonly interfaces. Through specific but simple design patterns, we can provide the dependecies that are required, by injecting them.

When creating any type of functionality, we can face a situation when we want to use a piece of code that will generate a different result depending of the environment where is run.

The goal of the next function is to return the "current" git version that the "user" is using.

```go
    func Version() string {
        cmd := exec.Command("git", "version")
        stdout, err := cmd.Output()
        if err != nil {
            panic(err)
        }
        n := len("git version ")
        version := string(stdout[n:])
        return strings.TrimSpace(version)
    }


    func TestVersion(t *testing.T) {
        got := Version()
        want := "2.24.0"
        if got != want {
            t.Errorf("Version() = %q; want %q", got, want)
        }
    }
```

The problem with this is that in our code, we are getting the right version of git, at least when this was coded on my computer. But what happens if the user running the tests doesn't have the same version of git? What if doesn't have git at all?

The best way to solve this is by using Dependency Injection.

**Using variables** (global-state limited to the package). Notice that the test need to be an internal test to be able to override the funcionality of execCommand.

```go
    var execCommand = exec.Command

    func Version() string {
        cmd := execCommand("git", "version")
        stdout, err := cmd.Output()
        if err != nil {
            panic(err)
        }
        n := len("git version ")
        version := string(stdout[n:])
        return strings.TrimSpace(version)
    }

    func TestVersion(t *testing.T) {
        execCommand = func(name string, arg ...string) *exec.Cmd {
            return exec.Command("echo", "git version 2.22.2")
        }
        defer func() { // TEAR-DOWN
            execCommand = exec.Command
        }()

        got := Version()
        want := "2.22.2"
        if got != want {
            t.Errorf("Version() = %q; want %q", got, want)
        }
    }
```

**Using Interfaces**. Creating a custom type. We can set a default value by using the `command` method but at the same time provide a custom one.

```go
    type Checker struct {
        execCommand func(name string, arg ...string) *exec.Cmd
    }

    func (gc *Checker) command(name string, arg ...string) *exec.Cmd {
        if gc.execCommand == nil {
            return exec.Command(name, arg...)
        }
        return gc.execCommand(name, arg...)
    }

    func (gc *Checker) Version() string {
        cmd := gc.command("git", "version")
        stdout, err := cmd.Output()
        if err != nil {
            panic(err)
        }
        n := len("git version ")
        version := string(stdout[n:])
        return strings.TrimSpace(version)
    }
```

The big difference with this one is that the tear-down is much simpler, and execCommand won't affect to the global (or package) state, only affects `checker := Checker{}` which is not accessible anywhere else.

```go
    func TestChecker_Version(t *testing.T) {
        checker := Checker{
            execCommand: func(name string, arg ...string) *exec.Cmd {
                return exec.Command("echo", "git version 2.22.2")
            },
        }
        got := checker.Version()
        want := "2.22.2"
        if got != want {
            t.Errorf("checker.Version() = %q; want %q", got, want)
        }
    }
```

## Mocks

The next example shows an easy and stubbed-out implementation of an email client.

```go
    type MailClient struct {
        // stuff here
    }

    func (mc *MailClient) Welcome(name, email string) error {
        // send out a welcome email to the user!
        return nil
    }

    // this is all fake just to make the demo work
    type User struct{}
    type UserStore struct{}

    func (us *UserStore) Create(name, email string) (*User, error) {
        // pretend to add user to DB
        return &User{}, nil
    }

    func Signup(name, email string, ec *MailClient, us *UserStore) (*User, error) {
        email = strings.ToLower(email)
        user, err := us.Create(name, email)
        if err != nil {
            return nil, err
        }
        err = ec.Welcome(name, email)
        if err != nil {
            return nil, err
        }
        return user, nil
    }
```

In order to test this client, a nice way would be introducing a layer for the email client to mock it's functionality. We would replace the implementation with another one and we would mock it out to run our test without using the real implementation:

```go
    type EmailClient interface {
        Welcome(name, email string) error
    }

    func Signup(name, email string, ec EmailClient, us *UserStore) (*User, error) {
        //...
    }
```

We now accept an interface instead of an strict type.

### Mocks: Faking API's

Given an Stripe client like this one, will make a charge of 20 USD to a test client, hitting the real Stripe API.

```go
    package main

    import (
        "encoding/json"
        "fmt"

        stripe "github.com/joncalhoun/twg/stripe/v0"
    )

    // Running this main, will hit the real (testing) stripe end-point.
    func main() {
        // curl https://api.stripe.com/v1/charges \
        //    -u sk_test_4eC39HqLyjWDarjtT1zdp7dc: \
        //    -d amount=2000 \
        //    -d currency=usd \
        //    -d source=tok_mastercard \
        //    -d description="Charge for jenny.rosen@example.com"
        c := stripe.Client{
            Key: "sk_test_4eC39HqLyjWDarjtT1zdp7dc",
        }
        charge, err := c.Charge(2000, "tok_mastercard", "Charge for demo purposes.")
        if err != nil {
            panic(err)
        }
        fmt.Println(charge)
        jsonBytes, err := json.Marshal(charge)
        if err != nil {
            panic(err)
        }
        fmt.Println(string(jsonBytes))
    }
```

To implement the `Charge` method, we would need to start to follow the correct design patterns...

```go
    package stripe

    import (
        "encoding/json"
        "io/ioutil"
        "net/http"
        "net/url"
        "strconv"
        "strings"
    )

    // This is a small subset of the Stripe charge fields
    type Charge struct {
        ID          string `json:"id"`
        Amount      int    `json:"amount"`
        Description string `json:"description"`
        Status      string `json:"status"`
    }

    type Client struct {
        Key     string
    }

    func (c *Client) Charge(amount int, source, desc string) (*Charge, error) {
        v := url.Values{}
        v.Set("amount", strconv.Itoa(amount))
        v.Set("currency", "usd")
        v.Set("source", source)
        v.Set("description", desc)

        req, err := http.NewRequest(http.MethodPost, "https://api.stripe.com/v1/charges",       strings.NewReader(v.Encode()))
        if err != nil {
            return nil, err
        }

        req.SetBasicAuth(c.Key, "")
        var client http.Client
        res, err := client.Do(req)
        if err != nil {
            return nil, err
        }

        defer res.Body.Close()
        resBytes, err := ioutil.ReadAll(res.Body)
        if err != nil {
            return nil, err
        }

        var charge Charge
        err = json.Unmarshal(resBytes, &charge)
        if err != nil {
            return nil, err
        }

        return &charge, nil
    }
```

So first step will be to define a way to handle the base URL, to make it alterable.

```go
    type Client struct {
        Key     string
        baseURL string
    }

    func (c *Client) BaseURL() string {
        if c.baseURL == "" {
            return "https://api.stripe.com"
        }
        return c.baseURL
    }

    func (c *Client) Charge(amount int, source, desc string) (*Charge, error) {
        // ...
        req, err := http.NewRequest(http.MethodPost, c.BaseURL()+"/v1/charges",
            strings.NewReader(v.Encode()))
        // ...
    }
```

What this allows us to do is to point to another web client.
So we can create a "fake client", like this one:

```go
    package stripe

    import (
        "net/http"
        "net/http/httptest"
        "testing"
    )

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

We can mock this in two different ways:
    1. Use the client but replace the server is communicating with.
    2. Replace the client we are communicating with.

#### 1. Replacing the server

```go
    package stripe_test

    import (
        "fmt"
        "net/http"
        "testing"

        stripe "github.com/noelruault/programming-training/interfaces-training/mocks/v2"
    )

    func TestAppV2(t *testing.T) {
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
```

#### 2. Replacing the client

With this option, we are injecting a client that has a base URL that communicates with a local server that we have control over. This approach allows us to write end to end tests and they will like as if all the application were running from start to end.

```go
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

        // Now inject client into the app and run the tests - they will use the
        // local test server using this mux.
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
```
