package moneytree

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestSanitizeURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in, want string
	}{
		{"/?a=b", "/?a=b"},
		{"/?a=b&client_secret=secret", "/?a=b&client_secret=REDACTED"},
		{"/?a=b&client_id=id&client_secret=secret", "/?a=b&client_id=id&client_secret=REDACTED"},
	}

	for _, tt := range tests {
		inURL, _ := url.Parse(tt.in)
		want, _ := url.Parse(tt.want)

		got := sanitizeURL(inURL)
		if got.String() != want.String() {
			t.Errorf("sanitizeURL(%v) returned %v, want %v", tt.in, got, want)
		}
	}
}

func TestNewRequest(t *testing.T) {
	t.Parallel()

	t.Run("success case: body is JSON encoded when provided", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			config: &Config{
				BaseURL: baseURL,
			},
		}

		body := map[string]string{
			"key": "value",
		}

		req, err := client.NewRequest(context.Background(), http.MethodPost, "test/path", body)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if req.Method != http.MethodPost {
			t.Errorf("expected method %s, got %s", http.MethodPost, req.Method)
		}

		expectedURL := "https://test.getmoneytree.com/test/path"
		if req.URL.String() != expectedURL {
			t.Errorf("expected URL %s, got %s", expectedURL, req.URL.String())
		}

		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", req.Header.Get("Content-Type"))
		}

		var buf bytes.Buffer
		_, err = io.Copy(&buf, req.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		bodyStr := buf.String()
		if !strings.Contains(bodyStr, "key") || !strings.Contains(bodyStr, "value") {
			t.Errorf("expected body to contain key and value, got %s", bodyStr)
		}
	})

	t.Run("success case: no body", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			config: &Config{
				BaseURL: baseURL,
			},
		}

		req, err := client.NewRequest(context.Background(), http.MethodGet, "test/path", nil)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if req.Method != http.MethodGet {
			t.Errorf("expected method %s, got %s", http.MethodGet, req.Method)
		}

		expectedURL := "https://test.getmoneytree.com/test/path"
		if req.URL.String() != expectedURL {
			t.Errorf("expected URL %s, got %s", expectedURL, req.URL.String())
		}

		if req.Header.Get("Content-Type") != "" {
			t.Errorf("expected empty Content-Type, got %s", req.Header.Get("Content-Type"))
		}
	})

	t.Run("success case: RequestOption is applied", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			config: &Config{
				BaseURL: baseURL,
			},
		}

		customHeader := "Custom-Header"
		customValue := "custom-value"

		req, err := client.NewRequest(context.Background(), http.MethodPost, "test/path", nil, func(r *http.Request) {
			r.Header.Set(customHeader, customValue)
		})
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if req.Header.Get(customHeader) != customValue {
			t.Errorf("expected %s header to be %s, got %s", customHeader, customValue, req.Header.Get(customHeader))
		}
	})

	t.Run("error case: BaseURL path does not end with slash", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			config: &Config{
				BaseURL: baseURL,
			},
		}

		_, err = client.NewRequest(context.Background(), http.MethodPost, "test/path", nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "baseURL must have a trailing slash") {
			t.Errorf("expected error about trailing slash, got %v", err)
		}
	})

	t.Run("error case: invalid URL", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			config: &Config{
				BaseURL: baseURL,
			},
		}

		// Specify invalid URL path
		_, err = client.NewRequest(context.Background(), http.MethodPost, "://invalid", nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestNewFormRequest(t *testing.T) {
	t.Parallel()

	t.Run("success case: Content-Type is set when body is provided", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			config: &Config{
				BaseURL: baseURL,
			},
		}

		body := strings.NewReader("key=value&foo=bar")
		req, err := client.NewFormRequest(context.Background(), "test/path", body)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if req.Method != http.MethodPost {
			t.Errorf("expected method %s, got %s", http.MethodPost, req.Method)
		}

		expectedURL := "https://test.getmoneytree.com/test/path"
		if req.URL.String() != expectedURL {
			t.Errorf("expected URL %s, got %s", expectedURL, req.URL.String())
		}

		expectedContentType := "application/x-www-form-urlencoded"
		if req.Header.Get("Content-Type") != expectedContentType {
			t.Errorf("expected Content-Type %s, got %s", expectedContentType, req.Header.Get("Content-Type"))
		}
	})

	t.Run("success case: RequestOption is applied", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			config: &Config{
				BaseURL: baseURL,
			},
		}

		customHeader := "Custom-Header"
		customValue := "custom-value"

		body := strings.NewReader("key=value")
		req, err := client.NewFormRequest(context.Background(), "test/path", body, func(r *http.Request) {
			r.Header.Set(customHeader, customValue)
		})
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if req.Header.Get(customHeader) != customValue {
			t.Errorf("expected %s header to be %s, got %s", customHeader, customValue, req.Header.Get(customHeader))
		}

		expectedContentType := "application/x-www-form-urlencoded"
		if req.Header.Get("Content-Type") != expectedContentType {
			t.Errorf("expected Content-Type %s, got %s", expectedContentType, req.Header.Get("Content-Type"))
		}
	})

	t.Run("error case: BaseURL path does not end with slash", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			config: &Config{
				BaseURL: baseURL,
			},
		}

		body := strings.NewReader("key=value")
		_, err = client.NewFormRequest(context.Background(), "test/path", body)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "baseURL must have a trailing slash") {
			t.Errorf("expected error about trailing slash, got %v", err)
		}
	})

	t.Run("error case: invalid URL", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			config: &Config{
				BaseURL: baseURL,
			},
		}

		body := strings.NewReader("key=value")
		// Specify invalid URL path
		_, err = client.NewFormRequest(context.Background(), "://invalid", body)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestWithBearerToken(t *testing.T) {
	t.Parallel()

	t.Run("正常系: Authorizationヘッダーが正しく設定される", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			config: &Config{
				BaseURL: baseURL,
			},
		}

		token := "test-access-token"
		req, err := client.NewRequest(context.Background(), http.MethodGet, "test/path", nil, WithBearerToken(token))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		expectedAuthHeader := "Bearer test-access-token"
		if req.Header.Get("Authorization") != expectedAuthHeader {
			t.Errorf("expected Authorization header %s, got %s", expectedAuthHeader, req.Header.Get("Authorization"))
		}
	})
}

func TestDo_RetryOnRateLimit(t *testing.T) {
	t.Parallel()

	t.Run("success case: retries on HTTP 429 and succeeds", func(t *testing.T) {
		t.Parallel()

		attemptCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			if attemptCount < 2 {
				// Return 429 on first attempt
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(`{"error": "rate_limit_exceeded", "error_description": "Too many requests"}`))
			} else {
				// Return success on retry
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status": "ok"}`))
			}
		}))
		defer server.Close()

		baseURL, err := url.Parse(server.URL + "/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			httpClient: http.DefaultClient,
			config: &Config{
				BaseURL: baseURL,
			},
			retryConfig: RetryConfig{
				MaxRetries: 3,
				BaseDelay:  10 * time.Millisecond, // Short delay for testing
				Enabled:   true,
			},
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		var result map[string]string
		resp, err := client.Do(context.Background(), req, &result)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if attemptCount != 2 {
			t.Errorf("expected 2 attempts, got %d", attemptCount)
		}

		if result["status"] != "ok" {
			t.Errorf("expected status 'ok', got %v", result)
		}
	})

	t.Run("success case: retries exhausted returns error", func(t *testing.T) {
		t.Parallel()

		attemptCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			// Always return 429
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error": "rate_limit_exceeded", "error_description": "Too many requests"}`))
		}))
		defer server.Close()

		baseURL, err := url.Parse(server.URL + "/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			httpClient: http.DefaultClient,
			config: &Config{
				BaseURL: baseURL,
			},
			retryConfig: RetryConfig{
				MaxRetries: 2,
				BaseDelay:  10 * time.Millisecond, // Short delay for testing
				Enabled:   true,
			},
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		resp, err := client.Do(context.Background(), req, nil)
		if err == nil {
			t.Error("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Errorf("expected APIError, got %T", err)
		}

		if apiErr.StatusCode != http.StatusTooManyRequests {
			t.Errorf("expected status code %d, got %d", http.StatusTooManyRequests, apiErr.StatusCode)
		}

		// Should have attempted MaxRetries + 1 times (initial + retries)
		expectedAttempts := 2 + 1 // MaxRetries + initial attempt
		if attemptCount != expectedAttempts {
			t.Errorf("expected %d attempts, got %d", expectedAttempts, attemptCount)
		}

		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	})

	t.Run("success case: retry disabled does not retry", func(t *testing.T) {
		t.Parallel()

		attemptCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			// Return 429
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"error": "rate_limit_exceeded", "error_description": "Too many requests"}`))
		}))
		defer server.Close()

		baseURL, err := url.Parse(server.URL + "/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			httpClient: http.DefaultClient,
			config: &Config{
				BaseURL: baseURL,
			},
			retryConfig: RetryConfig{
				MaxRetries: 3,
				BaseDelay:  10 * time.Millisecond,
				Enabled:   false, // Retry disabled
			},
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		resp, err := client.Do(context.Background(), req, nil)
		if err == nil {
			t.Error("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Errorf("expected APIError, got %T", err)
		}

		if apiErr.StatusCode != http.StatusTooManyRequests {
			t.Errorf("expected status code %d, got %d", http.StatusTooManyRequests, apiErr.StatusCode)
		}

		// Should only attempt once (no retry)
		if attemptCount != 1 {
			t.Errorf("expected 1 attempt, got %d", attemptCount)
		}

		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	})

	t.Run("success case: non-429 errors are not retried", func(t *testing.T) {
		t.Parallel()

		attemptCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attemptCount++
			// Return 401 (not retried)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "invalid_token", "error_description": "Invalid token"}`))
		}))
		defer server.Close()

		baseURL, err := url.Parse(server.URL + "/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			httpClient: http.DefaultClient,
			config: &Config{
				BaseURL: baseURL,
			},
			retryConfig: RetryConfig{
				MaxRetries: 3,
				BaseDelay:  10 * time.Millisecond,
				Enabled:   true,
			},
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		resp, err := client.Do(context.Background(), req, nil)
		if err == nil {
			t.Error("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Errorf("expected APIError, got %T", err)
		}

		if apiErr.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, apiErr.StatusCode)
		}

		// Should only attempt once (non-429 errors are not retried)
		if attemptCount != 1 {
			t.Errorf("expected 1 attempt, got %d", attemptCount)
		}

		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	})
}
