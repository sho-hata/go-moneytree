package moneytree

import (
	"fmt"
	"net/http"
	"net/url"
)

// Client is the main client for interacting with the Moneytree LINK API.
type Client struct {
	httpClient *http.Client
	config     *Config
}

func NewClient(accountName string) (*Client, error) {
	if accountName == "" {
		return nil, fmt.Errorf("account name is required")
	}

	c := &Client{
		httpClient: http.DefaultClient,
		config: &Config{
			BaseURL: &url.URL{
				Scheme: "https",
				Host:   fmt.Sprintf("%s.getmoneytree.com", accountName),
			},
		},
	}
	return c, nil
}
