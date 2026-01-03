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
	"time"
)

// setTestToken is a helper function to set a test token on a client.
func setTestToken(client *Client, accessToken string) {
	if accessToken == "" {
		return
	}
	// Initialize tokenMutex if it's nil (for test clients created directly)
	if client.tokenMutex == nil {
		client.tokenMutex = &sync.Mutex{}
	}
	now := int(time.Now().Unix())
	expiresIn := 3600
	refreshToken := "test-refresh-token"
	client.SetToken(&OauthToken{
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
		CreatedAt:    &now,
		ExpiresIn:    &expiresIn,
	})
}

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

		setTestToken(client, "test-access-token")
		response, err := client.GetAccountBalanceDetails(context.Background(), "account_key_123")
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

		setTestToken(client, "test-access-token")
		response, err := client.GetAccountBalanceDetails(context.Background(), "account_key_123")
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

		// Token is not set, so refreshToken should fail
		_, err = client.GetAccountBalanceDetails(context.Background(), "account_key_123")
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

		setTestToken(client, "test-token")
		_, err = client.GetAccountBalanceDetails(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
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

		setTestToken(client, "invalid-token")
		_, err = client.GetAccountBalanceDetails(context.Background(), accountID)
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

		setTestToken(client, "test-token")
		// nolint:staticcheck // passing nil context for testing purposes
		_, err = client.GetAccountBalanceDetails(nil, accountID)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetAccountDueBalances(t *testing.T) {
	t.Parallel()

	t.Run("success case: due balance is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		dueAmount := 50000.00
		expectedResponse := AccountDueBalances{
			DueBalances: AccountDueBalance{
				ID:        123,
				AccountID: 456,
				Date:      "2023-12-01",
				DueAmount: &dueAmount,
				DueDate:   "2023-12-25",
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/accounts/account_key_123/due_balances.json" {
				t.Errorf("expected path /link/accounts/account_key_123/due_balances.json, got %s", r.URL.Path)
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

		setTestToken(client, "test-access-token")
		response, err := client.GetAccountDueBalances(context.Background(), "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.DueBalances.ID != expectedResponse.DueBalances.ID {
			t.Errorf("expected ID %d, got %d", expectedResponse.DueBalances.ID, response.DueBalances.ID)
		}
		if response.DueBalances.AccountID != expectedResponse.DueBalances.AccountID {
			t.Errorf("expected AccountID %d, got %d", expectedResponse.DueBalances.AccountID, response.DueBalances.AccountID)
		}
		if response.DueBalances.Date != expectedResponse.DueBalances.Date {
			t.Errorf("expected Date %s, got %s", expectedResponse.DueBalances.Date, response.DueBalances.Date)
		}
		if response.DueBalances.DueDate != expectedResponse.DueBalances.DueDate {
			t.Errorf("expected DueDate %s, got %s", expectedResponse.DueBalances.DueDate, response.DueBalances.DueDate)
		}
		if response.DueBalances.DueAmount == nil || *response.DueBalances.DueAmount != *expectedResponse.DueBalances.DueAmount {
			t.Errorf("expected DueAmount %v, got %v", expectedResponse.DueBalances.DueAmount, response.DueBalances.DueAmount)
		}
	})

	t.Run("success case: due balance with null due_amount", func(t *testing.T) {
		t.Parallel()

		expectedResponse := AccountDueBalances{
			DueBalances: AccountDueBalance{
				ID:        123,
				AccountID: 456,
				Date:      "2023-12-01",
				DueAmount: nil,
				DueDate:   "2023-12-25",
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

		setTestToken(client, "test-access-token")
		response, err := client.GetAccountDueBalances(context.Background(), "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.DueBalances.DueAmount != nil {
			t.Errorf("expected DueAmount nil, got %v", response.DueBalances.DueAmount)
		}
	})

	t.Run("success case: due balance with since parameter", func(t *testing.T) {
		t.Parallel()

		dueAmount := 50000.00
		expectedResponse := AccountDueBalances{
			DueBalances: AccountDueBalance{
				ID:        123,
				AccountID: 456,
				Date:      "2023-12-01",
				DueAmount: &dueAmount,
				DueDate:   "2023-12-25",
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedSince := "2023-01-01"
			actualSince := r.URL.Query().Get("since")
			if actualSince != expectedSince {
				t.Errorf("expected since parameter %s, got %s", expectedSince, actualSince)
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

		setTestToken(client, "test-access-token")
		response, err := client.GetAccountDueBalances(context.Background(), "account_key_123",
			WithSinceForDueBalances("2023-01-01"),
		)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.DueBalances.DueDate != "2023-12-25" {
			t.Errorf("expected DueDate 2023-12-25, got %s", response.DueBalances.DueDate)
		}
	})

	t.Run("success case: due balance with page parameter", func(t *testing.T) {
		t.Parallel()

		dueAmount := 50000.00
		expectedResponse := AccountDueBalances{
			DueBalances: AccountDueBalance{
				ID:        123,
				AccountID: 456,
				Date:      "2023-12-01",
				DueAmount: &dueAmount,
				DueDate:   "2023-12-25",
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedPage := "2"
			actualPage := r.URL.Query().Get("page")
			if actualPage != expectedPage {
				t.Errorf("expected page parameter %s, got %s", expectedPage, actualPage)
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

		setTestToken(client, "test-access-token")
		response, err := client.GetAccountDueBalances(context.Background(), "account_key_123",
			WithPageForDueBalances(2),
		)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
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
		_, err = client.GetAccountDueBalances(context.Background(), "account_key_123")
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

		setTestToken(client, "test-token")
		_, err = client.GetAccountDueBalances(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when start_date is specified without end_date", func(t *testing.T) {
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
		_, err = client.GetAccountDueBalances(context.Background(), "account_key_123",
			WithStartDateForDueBalances("2023-01-01"),
		)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when end_date is specified without start_date", func(t *testing.T) {
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
		_, err = client.GetAccountDueBalances(context.Background(), "account_key_123",
			WithEndDateForDueBalances("2023-12-31"),
		)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error case: returns error when since date format is invalid", func(t *testing.T) {
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
		_, err = client.GetAccountDueBalances(context.Background(), "account_key_123",
			WithSinceForDueBalances("2023/01/01"),
		)
		if err == nil {
			t.Error("expected error, got nil")
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

		setTestToken(client, "invalid-token")
		_, err = client.GetAccountDueBalances(context.Background(), accountID)
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

		setTestToken(client, "test-token")
		// nolint:staticcheck // passing nil context for testing purposes
		_, err = client.GetAccountDueBalances(nil, accountID)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
