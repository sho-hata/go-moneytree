package moneytree

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestGetProfile(t *testing.T) {
	t.Parallel()

	t.Run("success case: profile information is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		expectedProfile := Profile{
			LocaleIdentifier: stringPtr("ja_JP"),
			Email:            stringPtr("user@example.com"),
			MoneytreeID:      stringPtr("1234567890"),
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/profile.json" {
				t.Errorf("expected path /link/profile.json, got %s", r.URL.Path)
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				t.Errorf("expected Authorization header with Bearer prefix, got %s", authHeader)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(expectedProfile); err != nil {
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

		profile, err := client.GetProfile(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if profile == nil {
			t.Fatal("expected profile, got nil")
		}
		if profile.LocaleIdentifier == nil || *profile.LocaleIdentifier != *expectedProfile.LocaleIdentifier {
			t.Errorf("expected LocaleIdentifier %s, got %v", *expectedProfile.LocaleIdentifier, profile.LocaleIdentifier)
		}
		if profile.Email == nil || *profile.Email != *expectedProfile.Email {
			t.Errorf("expected Email %s, got %v", *expectedProfile.Email, profile.Email)
		}
		if profile.MoneytreeID == nil || *profile.MoneytreeID != *expectedProfile.MoneytreeID {
			t.Errorf("expected MoneytreeID %s, got %v", *expectedProfile.MoneytreeID, profile.MoneytreeID)
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

		_, err = client.GetProfile(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error about access token, got %v", err)
		}
	})

	t.Run("error case: returns error when API returns an error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "invalid_token", "error_description": "The access token is invalid or expired."}`))
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

		_, err = client.GetProfile(context.Background(), "invalid-token")
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
		if !strings.Contains(err.Error(), "invalid_token") {
			t.Errorf("expected error about invalid_token, got %v", err)
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
		_, err = client.GetProfile(nil, "test-token")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestRevokeProfile(t *testing.T) {
	t.Parallel()

	t.Run("success case: profile connection is revoked correctly", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected method %s, got %s", http.MethodPost, r.Method)
			}
			if r.URL.Path != "/link/profile/revoke.json" {
				t.Errorf("expected path /link/profile/revoke.json, got %s", r.URL.Path)
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				t.Errorf("expected Authorization header with Bearer prefix, got %s", authHeader)
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

		err = client.RevokeProfile(context.Background(), "test-access-token")
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

		err = client.RevokeProfile(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error about access token, got %v", err)
		}
	})

	t.Run("error case: returns error when API returns an error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "invalid_token", "error_description": "The access token is invalid or expired."}`))
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

		err = client.RevokeProfile(context.Background(), "invalid-token")
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
		if !strings.Contains(err.Error(), "invalid_token") {
			t.Errorf("expected error about invalid_token, got %v", err)
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
		err = client.RevokeProfile(nil, "test-token")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestGetAccountGroups(t *testing.T) {
	t.Parallel()

	t.Run("success case: account groups information is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		lastAggregatedAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		lastAggregatedSuccess := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		id := int64(123)
		accountGroupID := int64(456)

		expectedResponse := AccountGroups{
			AccountGroups: []AccountGroup{
				{
					AggregationState:      "success",
					AggregationStatus:     "success",
					LastAggregatedAt:      lastAggregatedAt,
					LastAggregatedSuccess: &lastAggregatedSuccess,
					ID:                    &id,
					AccountGroup:          accountGroupID,
					InstitutionEntityKey:  "test_institution_key",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/profile/account_groups.json" {
				t.Errorf("expected path /link/profile/account_groups.json, got %s", r.URL.Path)
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				t.Errorf("expected Authorization header with Bearer prefix, got %s", authHeader)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(expectedResponse); err != nil {
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

		response, err := client.GetAccountGroups(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.AccountGroups) != 1 {
			t.Fatalf("expected 1 account group, got %d", len(response.AccountGroups))
		}

		ag := response.AccountGroups[0]
		if ag.AggregationState != expectedResponse.AccountGroups[0].AggregationState {
			t.Errorf("expected AggregationState %s, got %s", expectedResponse.AccountGroups[0].AggregationState, ag.AggregationState)
		}
		if ag.AggregationStatus != expectedResponse.AccountGroups[0].AggregationStatus {
			t.Errorf("expected AggregationStatus %s, got %s", expectedResponse.AccountGroups[0].AggregationStatus, ag.AggregationStatus)
		}
		if !ag.LastAggregatedAt.Equal(expectedResponse.AccountGroups[0].LastAggregatedAt) {
			t.Errorf("expected LastAggregatedAt %v, got %v", expectedResponse.AccountGroups[0].LastAggregatedAt, ag.LastAggregatedAt)
		}
		if ag.LastAggregatedSuccess == nil {
			t.Error("expected LastAggregatedSuccess, got nil")
		} else if !ag.LastAggregatedSuccess.Equal(*expectedResponse.AccountGroups[0].LastAggregatedSuccess) {
			t.Errorf("expected LastAggregatedSuccess %v, got %v", *expectedResponse.AccountGroups[0].LastAggregatedSuccess, *ag.LastAggregatedSuccess)
		}
		if ag.AccountGroup != expectedResponse.AccountGroups[0].AccountGroup {
			t.Errorf("expected AccountGroup %d, got %d", expectedResponse.AccountGroups[0].AccountGroup, ag.AccountGroup)
		}
		if ag.InstitutionEntityKey != expectedResponse.AccountGroups[0].InstitutionEntityKey {
			t.Errorf("expected InstitutionEntityKey %s, got %s", expectedResponse.AccountGroups[0].InstitutionEntityKey, ag.InstitutionEntityKey)
		}
	})

	t.Run("success case: account groups with null last_aggregated_success", func(t *testing.T) {
		t.Parallel()

		lastAggregatedAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		accountGroupID := int64(456)

		expectedResponse := AccountGroups{
			AccountGroups: []AccountGroup{
				{
					AggregationState:      "running",
					AggregationStatus:     "running.data",
					LastAggregatedAt:      lastAggregatedAt,
					LastAggregatedSuccess: nil,
					AccountGroup:          accountGroupID,
					InstitutionEntityKey:  "test_institution_key",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(expectedResponse); err != nil {
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

		response, err := client.GetAccountGroups(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.AccountGroups) != 1 {
			t.Fatalf("expected 1 account group, got %d", len(response.AccountGroups))
		}

		ag := response.AccountGroups[0]
		if ag.LastAggregatedSuccess != nil {
			t.Errorf("expected LastAggregatedSuccess nil, got %v", ag.LastAggregatedSuccess)
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

		_, err = client.GetAccountGroups(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error about access token, got %v", err)
		}
	})

	t.Run("error case: returns error when API returns an error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "invalid_token", "error_description": "The access token is invalid or expired."}`))
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

		_, err = client.GetAccountGroups(context.Background(), "invalid-token")
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
		if !strings.Contains(err.Error(), "invalid_token") {
			t.Errorf("expected error about invalid_token, got %v", err)
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
		_, err = client.GetAccountGroups(nil, "test-token")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestRefreshProfile(t *testing.T) {
	t.Parallel()

	t.Run("success case: refresh request is accepted correctly", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected method %s, got %s", http.MethodPost, r.Method)
			}
			if r.URL.Path != "/link/profile/refresh.json" {
				t.Errorf("expected path /link/profile/refresh.json, got %s", r.URL.Path)
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				t.Errorf("expected Authorization header with Bearer prefix, got %s", authHeader)
			}

			w.WriteHeader(http.StatusAccepted)
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

		err = client.RefreshProfile(context.Background(), "test-access-token")
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

		err = client.RefreshProfile(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error about access token, got %v", err)
		}
	})

	t.Run("error case: returns error when API returns 403 Forbidden", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error": "insufficient_scope", "error_description": "The request requires higher privileges than provided by the access token."}`))
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

		err = client.RefreshProfile(context.Background(), "invalid-token")
		if err == nil {
			t.Error("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Errorf("expected APIError, got %T", err)
		}
		if apiErr.StatusCode != http.StatusForbidden {
			t.Errorf("expected status code %d, got %d", http.StatusForbidden, apiErr.StatusCode)
		}
		if !strings.Contains(err.Error(), "insufficient_scope") {
			t.Errorf("expected error about insufficient_scope, got %v", err)
		}
	})

	t.Run("error case: returns error when API returns an error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "invalid_token", "error_description": "The access token is invalid or expired."}`))
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

		err = client.RefreshProfile(context.Background(), "invalid-token")
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
		if !strings.Contains(err.Error(), "invalid_token") {
			t.Errorf("expected error about invalid_token, got %v", err)
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
		err = client.RefreshProfile(nil, "test-token")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestRefreshAccountGroup(t *testing.T) {
	t.Parallel()

	t.Run("success case: refresh request is accepted correctly", func(t *testing.T) {
		t.Parallel()

		accountGroupID := int64(12345)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected method %s, got %s", http.MethodPost, r.Method)
			}
			expectedPath := fmt.Sprintf("/link/account_groups/%d/refresh.json", accountGroupID)
			if r.URL.Path != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				t.Errorf("expected Authorization header with Bearer prefix, got %s", authHeader)
			}

			w.WriteHeader(http.StatusAccepted)
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

		err = client.RefreshAccountGroup(context.Background(), "test-access-token", accountGroupID)
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

		err = client.RefreshAccountGroup(context.Background(), "", 12345)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error about access token, got %v", err)
		}
	})

	t.Run("error case: returns error when API returns 403 Forbidden", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error": "insufficient_scope", "error_description": "The request requires higher privileges than provided by the access token."}`))
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

		err = client.RefreshAccountGroup(context.Background(), "invalid-token", 12345)
		if err == nil {
			t.Error("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Errorf("expected APIError, got %T", err)
		}
		if apiErr.StatusCode != http.StatusForbidden {
			t.Errorf("expected status code %d, got %d", http.StatusForbidden, apiErr.StatusCode)
		}
		if !strings.Contains(err.Error(), "insufficient_scope") {
			t.Errorf("expected error about insufficient_scope, got %v", err)
		}
	})

	t.Run("error case: returns error when API returns an error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": "invalid_token", "error_description": "The access token is invalid or expired."}`))
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

		err = client.RefreshAccountGroup(context.Background(), "invalid-token", 12345)
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
		if !strings.Contains(err.Error(), "invalid_token") {
			t.Errorf("expected error about invalid_token, got %v", err)
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
		err = client.RefreshAccountGroup(nil, "test-token", 12345)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func stringPtr(s string) *string {
	return &s
}
