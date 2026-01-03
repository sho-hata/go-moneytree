package moneytree

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

func TestSubmitAccount2FA(t *testing.T) {
	t.Parallel()

	t.Run("success case: submits OTP successfully", func(t *testing.T) {
		t.Parallel()

		otp := "123456"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("expected method %s, got %s", http.MethodPut, r.Method)
			}
			if r.URL.Path != "/link/accounts/account_key_123/2fa.json" {
				t.Errorf("expected path /link/accounts/account_key_123/2fa.json, got %s", r.URL.Path)
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				t.Errorf("expected Authorization header with Bearer prefix, got %s", authHeader)
			}
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			var req SubmitAccount2FARequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}
			if req.KeyValues.OTP == nil || *req.KeyValues.OTP != otp {
				t.Errorf("expected OTP %s, got %v", otp, req.KeyValues.OTP)
			}
			if req.KeyValues.Captcha != nil {
				t.Errorf("expected Captcha nil, got %v", req.KeyValues.Captcha)
			}

			w.WriteHeader(http.StatusOK)
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
		}

		request := &SubmitAccount2FARequest{
			KeyValues: SubmitAccount2FAKeyValues{
				OTP: &otp,
			},
		}

		setTestToken(client, "test-access-token")
		err = client.SubmitAccount2FA(context.Background(), "account_key_123", request)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("success case: submits CAPTCHA successfully", func(t *testing.T) {
		t.Parallel()

		captcha := "captcha_answer"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("expected method %s, got %s", http.MethodPut, r.Method)
			}
			if r.URL.Path != "/link/accounts/account_key_123/2fa.json" {
				t.Errorf("expected path /link/accounts/account_key_123/2fa.json, got %s", r.URL.Path)
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				t.Errorf("expected Authorization header with Bearer prefix, got %s", authHeader)
			}
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			var req SubmitAccount2FARequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}
			if req.KeyValues.Captcha == nil || *req.KeyValues.Captcha != captcha {
				t.Errorf("expected Captcha %s, got %v", captcha, req.KeyValues.Captcha)
			}
			if req.KeyValues.OTP != nil {
				t.Errorf("expected OTP nil, got %v", req.KeyValues.OTP)
			}

			w.WriteHeader(http.StatusOK)
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
		}

		request := &SubmitAccount2FARequest{
			KeyValues: SubmitAccount2FAKeyValues{
				Captcha: &captcha,
			},
		}

		setTestToken(client, "test-access-token")
		err = client.SubmitAccount2FA(context.Background(), "account_key_123", request)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("error case: returns error when access token is empty", func(t *testing.T) {
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

		otp := "123456"
		request := &SubmitAccount2FARequest{
			KeyValues: SubmitAccount2FAKeyValues{
				OTP: &otp,
			},
		}

		// Token is not set, so refreshToken should fail
		err = client.SubmitAccount2FA(context.Background(), "account_key_123", request)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when account ID is empty", func(t *testing.T) {
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

		otp := "123456"
		request := &SubmitAccount2FARequest{
			KeyValues: SubmitAccount2FAKeyValues{
				OTP: &otp,
			},
		}

		setTestToken(client, "test-token")
		err = client.SubmitAccount2FA(context.Background(), "", request)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when request is nil", func(t *testing.T) {
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

		setTestToken(client, "test-token")
		err = client.SubmitAccount2FA(context.Background(), "account_key_123", nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when both OTP and Captcha are set", func(t *testing.T) {
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

		otp := "123456"
		captcha := "captcha_answer"
		request := &SubmitAccount2FARequest{
			KeyValues: SubmitAccount2FAKeyValues{
				OTP:     &otp,
				Captcha: &captcha,
			},
		}

		setTestToken(client, "test-token")
		err = client.SubmitAccount2FA(context.Background(), "account_key_123", request)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when neither OTP nor Captcha is set", func(t *testing.T) {
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

		request := &SubmitAccount2FARequest{
			KeyValues: SubmitAccount2FAKeyValues{},
		}

		setTestToken(client, "test-token")
		err = client.SubmitAccount2FA(context.Background(), "account_key_123", request)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when OTP exceeds 255 characters", func(t *testing.T) {
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

		longOTP := strings.Repeat("a", 256)
		request := &SubmitAccount2FARequest{
			KeyValues: SubmitAccount2FAKeyValues{
				OTP: &longOTP,
			},
		}

		setTestToken(client, "test-token")
		err = client.SubmitAccount2FA(context.Background(), "account_key_123", request)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when Captcha exceeds 255 characters", func(t *testing.T) {
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

		longCaptcha := strings.Repeat("a", 256)
		request := &SubmitAccount2FARequest{
			KeyValues: SubmitAccount2FAKeyValues{
				Captcha: &longCaptcha,
			},
		}

		setTestToken(client, "test-token")
		err = client.SubmitAccount2FA(context.Background(), "account_key_123", request)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when API returns an error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "invalid_request", "error_description": "Invalid OTP."}`))
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
		}

		otp := "123456"
		request := &SubmitAccount2FARequest{
			KeyValues: SubmitAccount2FAKeyValues{
				OTP: &otp,
			},
		}

		setTestToken(client, "test-token")
		err = client.SubmitAccount2FA(context.Background(), "account_key_123", request)
		if err == nil {
			t.Error("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Errorf("expected APIError, got %T", err)
		}
		if apiErr.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
		}
	})

	t.Run("error case: returns error when context is nil", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			httpClient: http.DefaultClient,
			config: &Config{
				BaseURL: baseURL,
			},
		}

		otp := "123456"
		request := &SubmitAccount2FARequest{
			KeyValues: SubmitAccount2FAKeyValues{
				OTP: &otp,
			},
		}

		// nolint:staticcheck // passing nil context for testing purposes
		setTestToken(client, "test-token")
		err = client.SubmitAccount2FA(nil, "account_key_123", request) //nolint:staticcheck
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetAccountCaptcha(t *testing.T) {
	t.Parallel()

	t.Run("success case: CAPTCHA image URL is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		expectedCaptchaImage := CaptchaImage{
			URL: "https://example.com/captcha/image.png",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/accounts/account_key_123/captcha.json" {
				t.Errorf("expected path /link/accounts/account_key_123/captcha.json, got %s", r.URL.Path)
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				t.Errorf("expected Authorization header with Bearer prefix, got %s", authHeader)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(expectedCaptchaImage); err != nil {
				t.Errorf("failed to encode response: %v", err)
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
		}

		setTestToken(client, "test-access-token")
		captchaImage, err := client.GetAccountCaptcha(context.Background(), "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if captchaImage == nil {
			t.Fatal("expected captchaImage, got nil")
		}
		if captchaImage.URL != expectedCaptchaImage.URL {
			t.Errorf("expected URL %s, got %s", expectedCaptchaImage.URL, captchaImage.URL)
		}
	})

	t.Run("error case: returns error when access token is empty", func(t *testing.T) {
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

		// Token is not set, so refreshToken should fail
		_, err = client.GetAccountCaptcha(context.Background(), "account_key_123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when account ID is empty", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			config: &Config{
				BaseURL: baseURL,
			},
			tokenMutex: &sync.Mutex{},
		}

		setTestToken(client, "test-token")
		_, err = client.GetAccountCaptcha(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when API returns an error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "invalid_request", "error_description": "Account status is not suspended.missing-answer.auth.captcha."}`))
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
		}

		setTestToken(client, "test-token")
		_, err = client.GetAccountCaptcha(context.Background(), "account_key_123")
		if err == nil {
			t.Error("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Errorf("expected APIError, got %T", err)
		}
		if apiErr.StatusCode != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
		}
	})

	t.Run("error case: returns error when context is nil", func(t *testing.T) {
		t.Parallel()

		baseURL, err := url.Parse("https://test.getmoneytree.com/")
		if err != nil {
			t.Fatalf("failed to parse base URL: %v", err)
		}

		client := &Client{
			httpClient: http.DefaultClient,
			config: &Config{
				BaseURL: baseURL,
			},
		}

		// nolint:staticcheck // passing nil context for testing purposes
		setTestToken(client, "test-token")
		_, err = client.GetAccountCaptcha(nil, "account_key_123") //nolint:staticcheck
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
