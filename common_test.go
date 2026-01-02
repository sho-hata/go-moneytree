package moneytree

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetAccountBalanceDetails(t *testing.T) {
	t.Parallel()

	t.Run("success case: balance details list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		balanceType1 := 0
		balanceType2 := 2
		expectedResponse := AccountBalanceDetails{
			AccountBalances: []AccountBalanceDetail{
				{
					ID:            123,
					AccountID:     456,
					Date:          "2023-12-01",
					Balance:       1000000.50,
					BalanceInBase: 1000000.50,
					BalanceType:   &balanceType1,
				},
				{
					ID:            124,
					AccountID:     456,
					Date:          "2023-12-02",
					Balance:       1005000.00,
					BalanceInBase: 1005000.00,
					BalanceType:   &balanceType2,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/accounts/account_key_123/balances/details.json" {
				t.Errorf("expected path /link/accounts/account_key_123/balances/details.json, got %s", r.URL.Path)
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

		response, err := client.GetAccountBalanceDetails(context.Background(), "test-access-token", "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.AccountBalances) != 2 {
			t.Fatalf("expected 2 balance details, got %d", len(response.AccountBalances))
		}

		detail1 := response.AccountBalances[0]
		if detail1.ID != expectedResponse.AccountBalances[0].ID {
			t.Errorf("expected ID %d, got %d", expectedResponse.AccountBalances[0].ID, detail1.ID)
		}
		if detail1.AccountID != expectedResponse.AccountBalances[0].AccountID {
			t.Errorf("expected AccountID %d, got %d", expectedResponse.AccountBalances[0].AccountID, detail1.AccountID)
		}
		if detail1.Date != expectedResponse.AccountBalances[0].Date {
			t.Errorf("expected Date %s, got %s", expectedResponse.AccountBalances[0].Date, detail1.Date)
		}
		if detail1.Balance != expectedResponse.AccountBalances[0].Balance {
			t.Errorf("expected Balance %f, got %f", expectedResponse.AccountBalances[0].Balance, detail1.Balance)
		}
		if detail1.BalanceInBase != expectedResponse.AccountBalances[0].BalanceInBase {
			t.Errorf("expected BalanceInBase %f, got %f", expectedResponse.AccountBalances[0].BalanceInBase, detail1.BalanceInBase)
		}
		if detail1.BalanceType == nil || *detail1.BalanceType != *expectedResponse.AccountBalances[0].BalanceType {
			t.Errorf("expected BalanceType %v, got %v", expectedResponse.AccountBalances[0].BalanceType, detail1.BalanceType)
		}

		detail2 := response.AccountBalances[1]
		if detail2.Balance != expectedResponse.AccountBalances[1].Balance {
			t.Errorf("expected Balance %f, got %f", expectedResponse.AccountBalances[1].Balance, detail2.Balance)
		}
		if detail2.BalanceType == nil || *detail2.BalanceType != *expectedResponse.AccountBalances[1].BalanceType {
			t.Errorf("expected BalanceType %v, got %v", expectedResponse.AccountBalances[1].BalanceType, detail2.BalanceType)
		}
	})

	t.Run("success case: empty balance details list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := AccountBalanceDetails{
			AccountBalances: []AccountBalanceDetail{},
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

		response, err := client.GetAccountBalanceDetails(context.Background(), "test-access-token", "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.AccountBalances) != 0 {
			t.Fatalf("expected 0 balance details, got %d", len(response.AccountBalances))
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

		_, err = client.GetAccountBalanceDetails(context.Background(), "", "account_key_123")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error about access token, got %v", err)
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

		_, err = client.GetAccountBalanceDetails(context.Background(), "test-token", "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "account ID is required") {
			t.Errorf("expected error about account ID, got %v", err)
		}
	})

	t.Run("error case: returns error when API returns an error", func(t *testing.T) {
		t.Parallel()

		accountID := "account_key_123"

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

		_, err = client.GetAccountBalanceDetails(context.Background(), "invalid-token", accountID)
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

		accountID := "account_key_123"

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
		_, err = client.GetAccountBalanceDetails(nil, "test-token", accountID)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

