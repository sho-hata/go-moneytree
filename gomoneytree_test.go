package moneytree

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
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
