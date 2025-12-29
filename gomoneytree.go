package moneytree

import (
	"fmt"
	"net/http"
)

// Client is the main client for interacting with the Moneytree LINK API.
type Client struct {
	httpClient *http.Client
	config     *Config
}

// NewClient creates a new Client with the given configuration.
//
// If httpClient is nil, a default http.Client will be used.
//
// Example:
//
//	client, err := moneytree.NewClient(&moneytree.Config{
//		BaseURL:      "https://myaccount-staging.getmoneytree.com",
//		ClientID:     "your-client-id",
//		ClientSecret: "your-client-secret",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
func NewClient(config *Config, httpClient *http.Client) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if config.ClientID == "" {
		return nil, fmt.Errorf("client ID is required")
	}
	if config.ClientSecret == "" {
		return nil, fmt.Errorf("client secret is required")
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Client{
		httpClient: httpClient,
		config:     config,
	}, nil
}
