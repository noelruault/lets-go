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

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

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

// Mocking the http.client would be like this:
func (c *Client) Charge(cli httpClient, amount int, source, desc string) (*Charge, error) {
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
	// var client http.Client
	res, err := cli.Do(req)
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
