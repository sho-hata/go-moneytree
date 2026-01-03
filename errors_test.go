package moneytree

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckResponseError(t *testing.T) {
	t.Parallel()

	t.Run("正常系: ステータスコード200の場合はnilを返す", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"access_token": "test-token"}`))
		}))
		defer server.Close()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		err = checkResponseError(resp)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("正常系: ステータスコード201の場合はnilを返す", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id": "123"}`))
		}))
		defer server.Close()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		err = checkResponseError(resp)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("エラーケース: レスポンスがnilの場合、エラーを返す", func(t *testing.T) {
		t.Parallel()

		err := checkResponseError(nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("エラーケース: ステータスコード400の場合、APIErrorを返す", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "invalid_grant", "error_description": "The provided authorization grant is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client."}`))
		}))
		defer server.Close()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		err = checkResponseError(resp)
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

	t.Run("正常系: ステータスコード500の場合はnilを返す（400-499の範囲のみがエラー）", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "internal_server_error"}`))
		}))
		defer server.Close()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		err = checkResponseError(resp)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("エラーケース: JSONパースできない場合はエラーを返す", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		err = checkResponseError(resp)
		if err == nil {
			t.Error("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Errorf("expected APIError, got %T", err)
		}
		if apiErr.RawMessage != "invalid json" {
			t.Errorf("expected raw message 'invalid json', got %s", apiErr.RawMessage)
		}
	})
}
