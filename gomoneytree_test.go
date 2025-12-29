package moneytree

import (
	"net/http"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("正常系: 有効な設定でクライアントを作成できる", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			BaseURL:      "https://example.com",
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		}

		client, err := NewClient(config, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.config != config {
			t.Error("expected config to be set")
		}
		if client.httpClient == nil {
			t.Error("expected httpClient to be set")
		}
	})

	t.Run("正常系: カスタムHTTPクライアントを指定できる", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			BaseURL:      "https://example.com",
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		}

		customHTTPClient := &http.Client{}
		client, err := NewClient(config, customHTTPClient)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.httpClient != customHTTPClient {
			t.Error("expected custom httpClient to be set")
		}
	})

	t.Run("エラーケース: configがnilの場合、エラーを返す", func(t *testing.T) {
		t.Parallel()

		_, err := NewClient(nil, nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "config cannot be nil") {
			t.Errorf("expected 'config cannot be nil' error, got %v", err)
		}
	})

	t.Run("エラーケース: BaseURLが空の場合、エラーを返す", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			BaseURL:      "",
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		}

		_, err := NewClient(config, nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "base URL is required") {
			t.Errorf("expected 'base URL is required' error, got %v", err)
		}
	})

	t.Run("エラーケース: ClientIDが空の場合、エラーを返す", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			BaseURL:      "https://example.com",
			ClientID:     "",
			ClientSecret: "test-client-secret",
		}

		_, err := NewClient(config, nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "client ID is required") {
			t.Errorf("expected 'client ID is required' error, got %v", err)
		}
	})

	t.Run("エラーケース: ClientSecretが空の場合、エラーを返す", func(t *testing.T) {
		t.Parallel()

		config := &Config{
			BaseURL:      "https://example.com",
			ClientID:     "test-client-id",
			ClientSecret: "",
		}

		_, err := NewClient(config, nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "client secret is required") {
			t.Errorf("expected 'client secret is required' error, got %v", err)
		}
	})
}

