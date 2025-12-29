package moneytree

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_GetAccessToken(t *testing.T) {
	t.Parallel()

	t.Run("正常系: 認可コードからアクセストークンを取得できる", func(t *testing.T) {
		t.Parallel()

		expectedToken := &TokenResponse{
			AccessToken:  "test-access-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			RefreshToken: "test-refresh-token",
			Scope:        "read write",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/oauth/token" {
				t.Errorf("expected /oauth/token, got %s", r.URL.Path)
			}
			if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
				t.Errorf("expected application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
			}

			body := make([]byte, r.ContentLength)
			_, _ = r.Body.Read(body)
			bodyStr := string(body)

			if !strings.Contains(bodyStr, "grant_type=authorization_code") {
				t.Error("expected grant_type=authorization_code in body")
			}
			if !strings.Contains(bodyStr, "code=test-code") {
				t.Error("expected code=test-code in body")
			}
			if !strings.Contains(bodyStr, "client_id=test-client-id") {
				t.Error("expected client_id=test-client-id in body")
			}
			if !strings.Contains(bodyStr, "client_secret=test-client-secret") {
				t.Error("expected client_secret=test-client-secret in body")
			}
			if !strings.Contains(bodyStr, "redirect_uri=https%3A%2F%2Fexample.com%2Fcallback") {
				t.Error("expected redirect_uri in body")
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(expectedToken)
		}))
		defer server.Close()

		config := &Config{
			BaseURL:      server.URL,
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		}
		client, err := NewClient(config, nil)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		req := &GetAccessTokenRequest{
			Code:        "test-code",
			RedirectURI: "https://example.com/callback",
		}

		token, err := client.GetAccessToken(context.Background(), req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if token.AccessToken != expectedToken.AccessToken {
			t.Errorf("expected access token %s, got %s", expectedToken.AccessToken, token.AccessToken)
		}
		if token.TokenType != expectedToken.TokenType {
			t.Errorf("expected token type %s, got %s", expectedToken.TokenType, token.TokenType)
		}
		if token.ExpiresIn != expectedToken.ExpiresIn {
			t.Errorf("expected expires in %d, got %d", expectedToken.ExpiresIn, token.ExpiresIn)
		}
		if token.RefreshToken != expectedToken.RefreshToken {
			t.Errorf("expected refresh token %s, got %s", expectedToken.RefreshToken, token.RefreshToken)
		}
		if token.Scope != expectedToken.Scope {
			t.Errorf("expected scope %s, got %s", expectedToken.Scope, token.Scope)
		}
	})

	t.Run("エラーケース: リクエストがnilの場合、エラーを返す", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			BaseURL:      "https://example.com",
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		}
		client, err := NewClient(config, nil)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		_, err = client.GetAccessToken(context.Background(), nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "request cannot be nil") {
			t.Errorf("expected 'request cannot be nil' error, got %v", err)
		}
	})

	t.Run("エラーケース: codeが空の場合、エラーを返す", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			BaseURL:      "https://example.com",
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		}
		client, err := NewClient(config, nil)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		req := &GetAccessTokenRequest{
			Code:        "",
			RedirectURI: "https://example.com/callback",
		}

		_, err = client.GetAccessToken(context.Background(), req)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "code is required") {
			t.Errorf("expected 'code is required' error, got %v", err)
		}
	})

	t.Run("エラーケース: redirect_uriが空の場合、エラーを返す", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			BaseURL:      "https://example.com",
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		}
		client, err := NewClient(config, nil)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		req := &GetAccessTokenRequest{
			Code:        "test-code",
			RedirectURI: "",
		}

		_, err = client.GetAccessToken(context.Background(), req)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "redirect_uri is required") {
			t.Errorf("expected 'redirect_uri is required' error, got %v", err)
		}
	})

	t.Run("エラーケース: APIがエラーレスポンスを返した場合、APIErrorを返す", func(t *testing.T) {
		t.Parallel()

		expectedError := "invalid_grant"
		expectedErrorDescription := "The provided authorization grant is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client."

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "invalid_grant", "error_description": "The provided authorization grant is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client."}`))
		}))
		defer server.Close()

		config := &Config{
			BaseURL:      server.URL,
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		}
		client, err := NewClient(config, nil)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		req := &GetAccessTokenRequest{
			Code:        "test-code",
			RedirectURI: "https://example.com/callback",
		}

		_, err = client.GetAccessToken(context.Background(), req)
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
		if apiErr.ErrorType != expectedError {
			t.Errorf("expected error type %s, got %s", expectedError, apiErr.ErrorType)
		}
		if apiErr.ErrorDescription != expectedErrorDescription {
			t.Errorf("expected error description %s, got %s", expectedErrorDescription, apiErr.ErrorDescription)
		}
		if !strings.Contains(apiErr.Error(), expectedError) {
			t.Errorf("expected error message to contain %s, got %s", expectedError, apiErr.Error())
		}
		if !strings.Contains(apiErr.Error(), expectedErrorDescription) {
			t.Errorf("expected error message to contain %s, got %s", expectedErrorDescription, apiErr.Error())
		}
	})

	t.Run("エラーケース: 無効なJSONレスポンスの場合、エラーを返す", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		config := &Config{
			BaseURL:      server.URL,
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		}
		client, err := NewClient(config, nil)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		req := &GetAccessTokenRequest{
			Code:        "test-code",
			RedirectURI: "https://example.com/callback",
		}

		_, err = client.GetAccessToken(context.Background(), req)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to unmarshal response") {
			t.Errorf("expected unmarshal error, got %v", err)
		}
	})
}
