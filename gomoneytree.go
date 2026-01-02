package moneytree

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RequestOption configures a request.
type RequestOption func(*http.Request)

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

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body any, opts ...RequestOption) (*http.Request, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}
	if !strings.HasSuffix(c.config.BaseURL.Path, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", c.config.BaseURL)
	}

	u, err := c.config.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

// NewFormRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
// Body is sent with Content-Type: application/x-www-form-urlencoded.
func (c *Client) NewFormRequest(ctx context.Context, urlStr string, body io.Reader, opts ...RequestOption) (*http.Request, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}
	if !strings.HasSuffix(c.config.BaseURL.Path, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", c.config.BaseURL)
	}

	u, err := c.config.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			if resp != nil && resp.Body != nil {
				_ = resp.Body.Close()
			}
			return resp, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		var e *url.Error
		if errors.As(err, &e) {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(url).String()
				if resp != nil && resp.Body != nil {
					_ = resp.Body.Close()
				}
				return resp, e
			}
		}

		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return resp, err
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if err := checkResponseError(resp); err != nil {
		return resp, err
	}

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}
	return resp, err
}

// WithBearerToken returns a RequestOption that sets the Authorization header
// with the provided bearer token.
func WithBearerToken(token string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

// sanitizeURL redacts sensitive parameters from the URL which may be
// exposed to the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	sensitiveParams := []string{"client_secret", "refresh_token", "access_token"}
	for _, param := range sensitiveParams {
		if len(params.Get(param)) > 0 {
			params.Set(param, "REDACTED")
		}
	}
	uri.RawQuery = params.Encode()
	return uri
}
