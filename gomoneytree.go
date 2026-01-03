package moneytree

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// RequestOption configures a request.
type RequestOption func(*http.Request)

// RetryConfig configures retry behavior for rate-limited requests.
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts for rate-limited requests.
	// Default is 3.
	MaxRetries int
	// BaseDelay is the base delay in milliseconds for exponential backoff.
	// Default is 3000ms as recommended by Moneytree LINK API documentation.
	BaseDelay time.Duration
	// Enabled enables automatic retry for rate-limited requests (HTTP 429).
	// Default is true.
	Enabled bool
}

// Client is the main client for interacting with the Moneytree LINK API.
type Client struct {
	httpClient  *http.Client
	config      *Config
	retryConfig RetryConfig
	token       *OauthToken
	tokenMutex  *sync.Mutex
	getTokenErr error
}

// newHTTPClient creates a new HTTP client with appropriate timeouts and connection pool settings.
// This function addresses the issues with the default HTTP client:
// 1. Sets timeouts to prevent indefinite waiting
// 2. Increases MaxIdleConnsPerHost to improve connection reuse
func newHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	// Increase connection pool settings
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100

	// Set timeouts for each step of the HTTP request
	transport.DialContext = (&net.Dialer{
		Timeout: 5 * time.Second,
	}).DialContext
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.ResponseHeaderTimeout = 10 * time.Second
	transport.IdleConnTimeout = 90 * time.Second

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
}

// NewClientOption configures options for creating a new Client.
type NewClientOption func(*Client)

// WithRetryConfig configures retry behavior for rate-limited requests.
// This option allows you to customize retry settings according to your needs.
//
// Example:
//
//	client, err := moneytree.NewClient("jp-api-staging",
//		moneytree.WithRetryConfig(moneytree.RetryConfig{
//			MaxRetries: 5,
//			BaseDelay: 5000 * time.Millisecond,
//			Enabled:   true,
//		}),
//	)
//
// Reference: https://docs.link.getmoneytree.com/docs/faq-rate-limiting
func WithRetryConfig(config RetryConfig) NewClientOption {
	return func(c *Client) {
		c.retryConfig = config
	}
}

func NewClient(accountName string, opts ...NewClientOption) (*Client, error) {
	if accountName == "" {
		return nil, fmt.Errorf("account name is required")
	}

	c := &Client{
		httpClient: newHTTPClient(),
		config: &Config{
			BaseURL: &url.URL{
				Scheme: "https",
				Host:   fmt.Sprintf("%s.getmoneytree.com", accountName),
			},
		},
		retryConfig: RetryConfig{
			MaxRetries: 3,
			BaseDelay:  3000 * time.Millisecond,
			Enabled:    true,
		},
		tokenMutex: &sync.Mutex{},
	}

	for _, opt := range opts {
		opt(c)
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

// isRateLimitError checks if the error is a rate limit error (HTTP 429).
func isRateLimitError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// calculateBackoffDelay calculates the exponential backoff delay with jitter.
// Formula: wait_interval = base * 2^n +/- jitter
// Reference: https://docs.link.getmoneytree.com/docs/faq-rate-limiting
func calculateBackoffDelay(baseDelay time.Duration, retryCount int) time.Duration {
	// Limit retryCount to prevent integer overflow
	if retryCount < 0 {
		retryCount = 0
	}
	if retryCount > 30 {
		retryCount = 30 // Limit to prevent overflow
	}

	// Calculate exponential backoff: base * 2^n
	// nolint:gosec // G115: retryCount is limited to 30, preventing overflow
	delay := baseDelay * time.Duration(1<<uint(retryCount))

	// Add jitter: random value between 0 and baseDelay
	// nolint:gosec // G404: Using math/rand is acceptable for jitter calculation (not security-sensitive)
	jitter := time.Duration(rand.Int63n(int64(baseDelay)))

	// Randomly add or subtract jitter
	// nolint:gosec // G404: Using math/rand is acceptable for jitter calculation (not security-sensitive)
	if rand.Intn(2) == 0 {
		delay += jitter
	} else {
		delay -= jitter
		if delay < baseDelay {
			delay = baseDelay
		}
	}

	return delay
}

// cloneRequest creates a clone of the HTTP request with a fresh body.
// This is necessary for retrying requests since the body can only be read once.
// The bodyBytes parameter should contain the original request body bytes.
// nolint:unparam // error return is kept for consistency with potential future error cases
func cloneRequest(req *http.Request, bodyBytes []byte) (*http.Request, error) {
	cloned := req.Clone(req.Context())
	if len(bodyBytes) > 0 {
		cloned.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	return cloned, nil
}

// setAuthorizationHeader sets the Authorization header on the request if a valid token is available.
// This method is thread-safe and checks if the token exists and has a valid access token.
func (c *Client) setAuthorizationHeader(req *http.Request) {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()
	if c.token != nil && c.token.AccessToken != nil && *c.token.AccessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *c.token.AccessToken))
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*http.Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	// Check if this is an OAuth token endpoint that doesn't require authentication
	requiresAuth := !c.isOAuthTokenEndpoint(req.URL)

	// Refresh token if authentication is required
	if requiresAuth {
		if err := c.refreshToken(ctx); err != nil {
			return nil, fmt.Errorf("refresh token: %w", err)
		}
		// Set Authorization header if token is available
		c.setAuthorizationHeader(req)
	}

	// Read the request body once and store it for potential retries
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		_ = req.Body.Close()
		// Restore the body for the first request
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	var lastErr error
	var lastResp *http.Response

	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		// Clone the request for retries (body can only be read once)
		var currentReq *http.Request
		if attempt == 0 {
			currentReq = req
		} else {
			var err error
			currentReq, err = cloneRequest(req, bodyBytes)
			if err != nil {
				return lastResp, fmt.Errorf("failed to clone request for retry: %w", err)
			}
			// Re-set Authorization header for retries if authentication is required
			if requiresAuth {
				c.setAuthorizationHeader(currentReq)
			}
		}

		resp, err := c.httpClient.Do(currentReq)
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

		// Check for rate limit errors
		if err := checkResponseError(resp); err != nil {
			lastErr = err
			lastResp = resp

			// If it's a rate limit error and retry is enabled, attempt retry
			if isRateLimitError(err) && c.retryConfig.Enabled && attempt < c.retryConfig.MaxRetries {
				// Close the response body before retrying
				_ = resp.Body.Close()

				// Calculate backoff delay
				delay := calculateBackoffDelay(c.retryConfig.BaseDelay, attempt)

				// Wait before retrying
				select {
				case <-ctx.Done():
					return resp, ctx.Err()
				case <-time.After(delay):
					// Continue to retry
					continue
				}
			}

			// Not a rate limit error, or retries exhausted, or retry disabled
			defer func() {
				if resp != nil && resp.Body != nil {
					_ = resp.Body.Close()
				}
			}()
			return resp, err
		}

		// Success - process the response
		defer func() {
			if resp != nil && resp.Body != nil {
				_ = resp.Body.Close()
			}
		}()

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

	// All retries exhausted
	if lastResp != nil && lastResp.Body != nil {
		_ = lastResp.Body.Close()
	}
	return lastResp, lastErr
}

// isOAuthTokenEndpoint checks if the URL is an OAuth token endpoint that doesn't require authentication.
func (c *Client) isOAuthTokenEndpoint(u *url.URL) bool {
	if u == nil {
		return false
	}
	path := u.Path
	return strings.HasSuffix(path, oauthTokenPath) || strings.HasSuffix(path, oauthRevokePath)
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

// validateDateFormat validates that the date string is in the format "2006-01-02" (YYYY-MM-DD).
func validateDateFormat(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf("date must be in format YYYY-MM-DD (e.g., 2020-11-08), got: %s", date)
	}
	return nil
}
