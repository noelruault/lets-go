package stripe

import (
	"encoding/json"
	"fmt"
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
	baseURL string
	// interface httpClient that defines Do
}

func (c *Client) BaseURL() string {
	if c.baseURL == "" {
		return "https://api.stripe.com"
	}
	return c.baseURL
}

type StripeClient interface {
	Charge(amount int, source, desc string) (*Charge, error)
}

type App struct {
	sc StripeClient
}

func (a *App) Run() error {
	charge, err := a.sc.Charge(0, "", "")
	if err != nil {
		return err
	}
	fmt.Println(charge)
	return nil
}

// Mocking the http.client would be like this:
// func (c *Client) Charge(cli httpClient, amount int, source, desc string) (*Charge, error) {
// or might come inside of *Client
func (c *Client) Charge(amount int, source, desc string) (*Charge, error) {
	v := url.Values{}
	v.Set("amount", strconv.Itoa(amount))
	v.Set("currency", "usd")
	v.Set("source", source)
	v.Set("description", desc)

	req, err := http.NewRequest(http.MethodPost, c.BaseURL()+"/v1/charges",
		strings.NewReader(v.Encode()))
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
