package moneytree

import "net/url"

type Config struct {
	BaseURL      *url.URL
	ClientID     string
	ClientSecret string
}
