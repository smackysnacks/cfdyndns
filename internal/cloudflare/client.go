package cloudflare

import (
	"fmt"
	"net/http"
)

const BaseUrl = "https://api.cloudflare.com/client/v4/"

type Client struct {
	client *http.Client

	apiToken string
}

func New(apiToken string) *Client {
	return &Client{
		client:   http.DefaultClient,
		apiToken: apiToken,
	}
}

func (c *Client) addStandardHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
}
