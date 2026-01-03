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
)

func TestGetPersonalAccounts(t *testing.T) {
	t.Parallel()

	t.Run("success case: accounts list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		id1 := int64(123)
		id2 := int64(456)
		name1 := stringPtr("普通預金")
		name2 := stringPtr("クレジットカード")
		balance1 := float64Ptr(100000.50)
		balance2 := float64Ptr(-5000.00)
		currency := stringPtr("JPY")
		lastAggregatedAt := "2023-01-01"

		expectedResponse := PersonalAccounts{
			Accounts: []PersonalAccount{
				{
					ID:                   &id1,
					AccountKey:           "account_key_1",
					AccountGroup:         789,
					InstitutionEntityKey: "test_bank_1",
					AccountType:          "bank",
					Name:                 name1,
					Balance:              balance1,
					Currency:             currency,
					LastAggregatedAt:     stringPtr(lastAggregatedAt),
				},
				{
					ID:                   &id2,
					AccountKey:           "account_key_2",
					AccountGroup:         789,
					InstitutionEntityKey: "test_bank_2",
					AccountType:          "credit_card",
					Name:                 name2,
					Balance:              balance2,
					Currency:             currency,
					LastAggregatedAt:     stringPtr(lastAggregatedAt),
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/accounts.json" {
				t.Errorf("expected path /link/accounts.json, got %s", r.URL.Path)
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
		response, err := client.GetPersonalAccounts(context.Background())
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Accounts) != 2 {
			t.Fatalf("expected 2 accounts, got %d", len(response.Accounts))
		}

		account1 := response.Accounts[0]
		if account1.AccountKey != expectedResponse.Accounts[0].AccountKey {
			t.Errorf("expected AccountKey %s, got %s", expectedResponse.Accounts[0].AccountKey, account1.AccountKey)
		}
		if account1.AccountType != expectedResponse.Accounts[0].AccountType {
			t.Errorf("expected AccountType %s, got %s", expectedResponse.Accounts[0].AccountType, account1.AccountType)
		}
		if account1.Name == nil || *account1.Name != *expectedResponse.Accounts[0].Name {
			t.Errorf("expected Name %s, got %v", *expectedResponse.Accounts[0].Name, account1.Name)
		}
		if account1.Balance == nil || *account1.Balance != *expectedResponse.Accounts[0].Balance {
			t.Errorf("expected Balance %v, got %v", *expectedResponse.Accounts[0].Balance, account1.Balance)
		}
		if account1.AccountGroup != expectedResponse.Accounts[0].AccountGroup {
			t.Errorf("expected AccountGroup %d, got %d", expectedResponse.Accounts[0].AccountGroup, account1.AccountGroup)
		}
		if account1.InstitutionEntityKey != expectedResponse.Accounts[0].InstitutionEntityKey {
			t.Errorf("expected InstitutionEntityKey %s, got %s", expectedResponse.Accounts[0].InstitutionEntityKey, account1.InstitutionEntityKey)
		}

		account2 := response.Accounts[1]
		if account2.AccountType != expectedResponse.Accounts[1].AccountType {
			t.Errorf("expected AccountType %s, got %s", expectedResponse.Accounts[1].AccountType, account2.AccountType)
		}
		if account2.Balance == nil || *account2.Balance != *expectedResponse.Accounts[1].Balance {
			t.Errorf("expected Balance %v, got %v", *expectedResponse.Accounts[1].Balance, account2.Balance)
		}
	})

	t.Run("success case: accounts list with null balance", func(t *testing.T) {
		t.Parallel()

		accountKey := "account_key_1"
		name := stringPtr("普通預金")

		expectedResponse := PersonalAccounts{
			Accounts: []PersonalAccount{
				{
					AccountKey:           accountKey,
					AccountGroup:         789,
					InstitutionEntityKey: "test_bank_1",
					AccountType:          "bank",
					Name:                 name,
					Balance:              nil,
					Currency:             stringPtr("JPY"),
					LastAggregatedAt:     nil,
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

		setTestToken(client, "test-access-token")
		response, err := client.GetPersonalAccounts(context.Background())
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Accounts) != 1 {
			t.Fatalf("expected 1 account, got %d", len(response.Accounts))
		}

		account := response.Accounts[0]
		if account.Balance != nil {
			t.Errorf("expected Balance nil, got %v", account.Balance)
		}
		if account.LastAggregatedAt != nil {
			t.Errorf("expected LastAggregatedAt nil, got %v", account.LastAggregatedAt)
		}
	})

	t.Run("success case: empty accounts list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PersonalAccounts{
			Accounts: []PersonalAccount{},
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
		response, err := client.GetPersonalAccounts(context.Background())
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Accounts) != 0 {
			t.Fatalf("expected 0 accounts, got %d", len(response.Accounts))
		}
	})

	t.Run("success case: accounts list with page parameter", func(t *testing.T) {
		t.Parallel()

		name := stringPtr("普通預金")
		balance := float64Ptr(100000.50)
		currency := stringPtr("JPY")

		expectedResponse := PersonalAccounts{
			Accounts: []PersonalAccount{
				{
					AccountKey:           "account_key_1",
					AccountGroup:         789,
					InstitutionEntityKey: "test_bank_1",
					AccountType:          "bank",
					Name:                 name,
					Balance:              balance,
					Currency:             currency,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/accounts.json" {
				t.Errorf("expected path /link/accounts.json, got %s", r.URL.Path)
			}
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
		response, err := client.GetPersonalAccounts(context.Background(), WithPage(2))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Accounts) != 1 {
			t.Fatalf("expected 1 account, got %d", len(response.Accounts))
		}
	})

	t.Run("success case: accounts list with per_page parameter", func(t *testing.T) {
		t.Parallel()

		name := stringPtr("普通預金")
		balance := float64Ptr(100000.50)
		currency := stringPtr("JPY")

		expectedResponse := PersonalAccounts{
			Accounts: []PersonalAccount{
				{
					AccountKey:           "account_key_1",
					AccountGroup:         789,
					InstitutionEntityKey: "test_bank_1",
					AccountType:          "bank",
					Name:                 name,
					Balance:              balance,
					Currency:             currency,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedPerPage := "100"
			actualPerPage := r.URL.Query().Get("per_page")
			if actualPerPage != expectedPerPage {
				t.Errorf("expected per_page parameter %s, got %s", expectedPerPage, actualPerPage)
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
		response, err := client.GetPersonalAccounts(context.Background(), WithPerPage(100))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Accounts) != 1 {
			t.Fatalf("expected 1 account, got %d", len(response.Accounts))
		}
	})

	t.Run("success case: accounts list with both page and per_page parameters", func(t *testing.T) {
		t.Parallel()

		name := stringPtr("普通預金")
		balance := float64Ptr(100000.50)
		currency := stringPtr("JPY")

		expectedResponse := PersonalAccounts{
			Accounts: []PersonalAccount{
				{
					AccountKey:           "account_key_1",
					AccountGroup:         789,
					InstitutionEntityKey: "test_bank_1",
					AccountType:          "bank",
					Name:                 name,
					Balance:              balance,
					Currency:             currency,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedPage := "3"
			actualPage := r.URL.Query().Get("page")
			if actualPage != expectedPage {
				t.Errorf("expected page parameter %s, got %s", expectedPage, actualPage)
			}

			expectedPerPage := "50"
			actualPerPage := r.URL.Query().Get("per_page")
			if actualPerPage != expectedPerPage {
				t.Errorf("expected per_page parameter %s, got %s", expectedPerPage, actualPerPage)
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
		response, err := client.GetPersonalAccounts(context.Background(), WithPage(3), WithPerPage(50))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Accounts) != 1 {
			t.Fatalf("expected 1 account, got %d", len(response.Accounts))
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
		_, err = client.GetPersonalAccounts(context.Background())
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

		setTestToken(client, "invalid-token")
		_, err = client.GetPersonalAccounts(context.Background())
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
		setTestToken(client, "test-token")
		_, err = client.GetPersonalAccounts(nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestWithSinceForBalances_InvalidDateFormat(t *testing.T) {
	t.Parallel()

	t.Run("error case: returns error when date format is invalid", func(t *testing.T) {
		t.Parallel()

		invalidDates := []string{
			"2023/01/01",
			"2023-1-1",
			"01-01-2023",
			"2023-01-01T00:00:00Z",
			"invalid",
			"",
		}

		for _, invalidDate := range invalidDates {
			invalidDate := invalidDate
			t.Run(fmt.Sprintf("invalid date: %s", invalidDate), func(t *testing.T) {
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
		_, err = client.GetPersonalAccountBalances(context.Background(), "account_key_123",
					WithSinceForBalances(invalidDate),
				)
				if err == nil {
					t.Errorf("expected error for invalid date format: %s", invalidDate)
				}
				if !strings.Contains(err.Error(), "date must be in format YYYY-MM-DD") {
					t.Errorf("expected error about date format, got: %v", err)
				}
			})
		}
	})

	t.Run("success case: accepts valid date format", func(t *testing.T) {
		t.Parallel()

		validDates := []string{
			"2023-01-01",
			"2020-11-08",
			"2000-12-31",
		}

		for _, validDate := range validDates {
			validDate := validDate
			t.Run(fmt.Sprintf("valid date: %s", validDate), func(t *testing.T) {
				t.Parallel()

				opt := WithSinceForBalances(validDate)
				if opt == nil {
					t.Error("expected non-nil option function")
				}

				// オプション関数が正常に適用されることを確認（エラーが発生しない）
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

				// オプション関数を適用してもエラーが発生しないことを確認
				// （実際のAPI呼び出しは失敗するが、日付フォーマットエラーではない）
				setTestToken(client, "test-token")
		_, err = client.GetPersonalAccountBalances(context.Background(), "account_key_123",
					opt,
				)
				// 日付フォーマットエラーではないことを確認
				if err != nil && strings.Contains(err.Error(), "date must be in format YYYY-MM-DD") {
					t.Errorf("unexpected date format error for valid date: %s, error: %v", validDate, err)
				}
			})
		}
	})
}

func TestGetPersonalAccountBalances(t *testing.T) {
	t.Parallel()

	t.Run("success case: balances list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		accountID := "account_key_123"
		id1 := int64(1)
		id2 := int64(2)
		accountIDValue := int64(123)
		balance1 := 100000.50
		balance2 := 105000.75
		balanceInBase1 := 100000.50
		balanceInBase2 := 105000.75
		date1 := "2023-01-01"
		date2 := "2023-01-02"

		expectedResponse := PersonalAccountBalances{
			AccountBalances: []PersonalAccountBalance{
				{
					ID:            id1,
					AccountID:     accountIDValue,
					Date:          date1,
					Balance:       balance1,
					BalanceInBase: balanceInBase1,
				},
				{
					ID:            id2,
					AccountID:     accountIDValue,
					Date:          date2,
					Balance:       balance2,
					BalanceInBase: balanceInBase2,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			expectedPath := fmt.Sprintf("/link/accounts/%s/balances.json", accountID)
			if r.URL.Path != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
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
		response, err := client.GetPersonalAccountBalances(context.Background(), accountID)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.AccountBalances) != 2 {
			t.Fatalf("expected 2 balances, got %d", len(response.AccountBalances))
		}

		bal1 := response.AccountBalances[0]
		if bal1.ID != expectedResponse.AccountBalances[0].ID {
			t.Errorf("expected ID %d, got %d", expectedResponse.AccountBalances[0].ID, bal1.ID)
		}
		if bal1.AccountID != expectedResponse.AccountBalances[0].AccountID {
			t.Errorf("expected AccountID %d, got %d", expectedResponse.AccountBalances[0].AccountID, bal1.AccountID)
		}
		if bal1.Balance != expectedResponse.AccountBalances[0].Balance {
			t.Errorf("expected Balance %v, got %v", expectedResponse.AccountBalances[0].Balance, bal1.Balance)
		}
		if bal1.BalanceInBase != expectedResponse.AccountBalances[0].BalanceInBase {
			t.Errorf("expected BalanceInBase %v, got %v", expectedResponse.AccountBalances[0].BalanceInBase, bal1.BalanceInBase)
		}
		if bal1.Date != expectedResponse.AccountBalances[0].Date {
			t.Errorf("expected Date %s, got %s", expectedResponse.AccountBalances[0].Date, bal1.Date)
		}

		bal2 := response.AccountBalances[1]
		if bal2.Balance != expectedResponse.AccountBalances[1].Balance {
			t.Errorf("expected Balance %v, got %v", expectedResponse.AccountBalances[1].Balance, bal2.Balance)
		}
		if bal2.BalanceInBase != expectedResponse.AccountBalances[1].BalanceInBase {
			t.Errorf("expected BalanceInBase %v, got %v", expectedResponse.AccountBalances[1].BalanceInBase, bal2.BalanceInBase)
		}
		if bal2.Date != expectedResponse.AccountBalances[1].Date {
			t.Errorf("expected Date %s, got %s", expectedResponse.AccountBalances[1].Date, bal2.Date)
		}
	})

	t.Run("success case: balances list with since parameter", func(t *testing.T) {
		t.Parallel()

		accountID := "account_key_123"
		sinceTime := "2023-01-01"
		id := int64(1)
		accountIDValue := int64(123)
		balance := 100000.50
		balanceInBase := 100000.50
		date := "2023-01-02"

		expectedResponse := PersonalAccountBalances{
			AccountBalances: []PersonalAccountBalance{
				{
					ID:            id,
					AccountID:     accountIDValue,
					Date:          date,
					Balance:       balance,
					BalanceInBase: balanceInBase,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualSince := r.URL.Query().Get("since")
			if actualSince != sinceTime {
				t.Errorf("expected since parameter %s, got %s", sinceTime, actualSince)
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
		response, err := client.GetPersonalAccountBalances(context.Background(), accountID, WithSinceForBalances(sinceTime))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.AccountBalances) != 1 {
			t.Fatalf("expected 1 balance, got %d", len(response.AccountBalances))
		}
	})

	t.Run("success case: balances list with page and per_page parameters", func(t *testing.T) {
		t.Parallel()

		accountID := "account_key_123"
		id := int64(1)
		accountIDValue := int64(123)
		balance := 100000.50
		balanceInBase := 100000.50
		date := "2023-01-02"

		expectedResponse := PersonalAccountBalances{
			AccountBalances: []PersonalAccountBalance{
				{
					ID:            id,
					AccountID:     accountIDValue,
					Date:          date,
					Balance:       balance,
					BalanceInBase: balanceInBase,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedPage := "2"
			actualPage := r.URL.Query().Get("page")
			if actualPage != expectedPage {
				t.Errorf("expected page parameter %s, got %s", expectedPage, actualPage)
			}

			expectedPerPage := "50"
			actualPerPage := r.URL.Query().Get("per_page")
			if actualPerPage != expectedPerPage {
				t.Errorf("expected per_page parameter %s, got %s", expectedPerPage, actualPerPage)
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
		response, err := client.GetPersonalAccountBalances(context.Background(), accountID,
			WithPageForBalances(2),
			WithPerPageForBalances(50),
		)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.AccountBalances) != 1 {
			t.Fatalf("expected 1 balance, got %d", len(response.AccountBalances))
		}
	})

	t.Run("success case: empty balances list", func(t *testing.T) {
		t.Parallel()

		accountID := "account_key_123"

		expectedResponse := PersonalAccountBalances{
			AccountBalances: []PersonalAccountBalance{},
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
		response, err := client.GetPersonalAccountBalances(context.Background(), accountID)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.AccountBalances) != 0 {
			t.Fatalf("expected 0 balances, got %d", len(response.AccountBalances))
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
		_, err = client.GetPersonalAccountBalances(context.Background(), "account_key_123")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "token refresh function is not set") && !strings.Contains(err.Error(), "refresh token") {
			t.Errorf("expected error about token refresh, got %v", err)
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
		_, err = client.GetPersonalAccountBalances(context.Background(), "")
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

		setTestToken(client, "invalid-token")
		_, err = client.GetPersonalAccountBalances(context.Background(), accountID)
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

		setTestToken(client, "test-token")
		// nolint:staticcheck // passing nil context for testing purposes
		_, err = client.GetPersonalAccountBalances(nil, accountID)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestGetTermDeposits(t *testing.T) {
	t.Parallel()

	t.Run("success case: term deposits list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		purchaseDate := "2023-01-01"
		maturityDate := "2025-01-01"
		nameRaw := "定期預金"
		nameClean := "定期預金（補正済み）"
		termLengthYear := 2
		termLengthMonth := 0
		termLengthDay := 0

		expectedResponse := TermDeposits{
			TermDeposits: []TermDeposit{
				{
					ID:              1048,
					AccountID:       123,
					Date:            "2023-12-01",
					PurchaseDate:    &purchaseDate,
					MaturityDate:    &maturityDate,
					NameRaw:         &nameRaw,
					NameClean:       &nameClean,
					Value:           1050000.00,
					CostBasis:       1000000.00,
					InterestRate:    0.25,
					Currency:        "JPY",
					TermLengthYear:  &termLengthYear,
					TermLengthMonth: &termLengthMonth,
					TermLengthDay:   &termLengthDay,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/accounts/account_key_123/term_deposits.json" {
				t.Errorf("expected path /link/accounts/account_key_123/term_deposits.json, got %s", r.URL.Path)
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
		response, err := client.GetTermDeposits(context.Background(), "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.TermDeposits) != 1 {
			t.Fatalf("expected 1 term deposit, got %d", len(response.TermDeposits))
		}

		deposit := response.TermDeposits[0]
		if deposit.ID != 1048 {
			t.Errorf("expected ID 1048, got %d", deposit.ID)
		}
		if deposit.AccountID != 123 {
			t.Errorf("expected AccountID 123, got %d", deposit.AccountID)
		}
		if deposit.Date != "2023-12-01" {
			t.Errorf("expected Date 2023-12-01, got %s", deposit.Date)
		}
		if deposit.Value != 1050000.00 {
			t.Errorf("expected Value 1050000.00, got %f", deposit.Value)
		}
		if deposit.CostBasis != 1000000.00 {
			t.Errorf("expected CostBasis 1000000.00, got %f", deposit.CostBasis)
		}
		if deposit.InterestRate != 0.25 {
			t.Errorf("expected InterestRate 0.25, got %f", deposit.InterestRate)
		}
		if deposit.Currency != "JPY" {
			t.Errorf("expected Currency JPY, got %s", deposit.Currency)
		}
	})

	t.Run("success case: empty term deposits list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := TermDeposits{
			TermDeposits: []TermDeposit{},
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
		response, err := client.GetTermDeposits(context.Background(), "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.TermDeposits) != 0 {
			t.Fatalf("expected 0 term deposits, got %d", len(response.TermDeposits))
		}
	})

	t.Run("success case: term deposits list with page parameter", func(t *testing.T) {
		t.Parallel()

		purchaseDate := "2023-01-01"
		maturityDate := "2025-01-01"

		expectedResponse := TermDeposits{
			TermDeposits: []TermDeposit{
				{
					ID:           1048,
					AccountID:    123,
					Date:         "2023-12-01",
					PurchaseDate: &purchaseDate,
					MaturityDate: &maturityDate,
					Value:        1050000.00,
					CostBasis:    1000000.00,
					InterestRate: 0.25,
					Currency:     "JPY",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/accounts/account_key_123/term_deposits.json" {
				t.Errorf("expected path /link/accounts/account_key_123/term_deposits.json, got %s", r.URL.Path)
			}
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
		response, err := client.GetTermDeposits(context.Background(), "account_key_123", WithPageForTermDeposits(2))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.TermDeposits) != 1 {
			t.Fatalf("expected 1 term deposit, got %d", len(response.TermDeposits))
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
		_, err = client.GetTermDeposits(context.Background(), "account_key_123")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "token refresh function is not set") && !strings.Contains(err.Error(), "refresh token") {
			t.Errorf("expected error about token refresh, got %v", err)
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
		_, err = client.GetTermDeposits(context.Background(), "")
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

		setTestToken(client, "invalid-token")
		_, err = client.GetTermDeposits(context.Background(), accountID)
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

		setTestToken(client, "test-token")
		// nolint:staticcheck // passing nil context for testing purposes
		_, err = client.GetTermDeposits(nil, accountID)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestGetPersonalAccountTransactions(t *testing.T) {
	t.Parallel()

	t.Run("success case: transactions list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		descriptionGuest := "テスト取引"
		descriptionPretty := "テスト取引（補正済み）"
		descriptionRaw := "テスト取引（生データ）"
		categoryEntityKey := "category_key_123"

		expectedResponse := PersonalAccountTransactions{
			Transactions: []PersonalAccountTransaction{
				{
					ID:                1048,
					Amount:            -5000.00,
					Date:              "2023-12-01T10:00:00Z",
					DescriptionGuest:  &descriptionGuest,
					DescriptionPretty: &descriptionPretty,
					DescriptionRaw:    &descriptionRaw,
					AccountID:         123,
					CategoryID:        456,
					Attributes:        PersonalAccountTransactionAttributes{},
					CategoryEntityKey: &categoryEntityKey,
					CreatedAt:         "2023-12-01T09:00:00Z",
					UpdatedAt:         "2023-12-01T09:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/accounts/account_key_123/transactions.json" {
				t.Errorf("expected path /link/accounts/account_key_123/transactions.json, got %s", r.URL.Path)
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
		response, err := client.GetPersonalAccountTransactions(context.Background(), "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Transactions) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(response.Transactions))
		}

		transaction := response.Transactions[0]
		if transaction.ID != 1048 {
			t.Errorf("expected ID 1048, got %d", transaction.ID)
		}
		if transaction.Amount != -5000.00 {
			t.Errorf("expected Amount -5000.00, got %f", transaction.Amount)
		}
		if transaction.Date != "2023-12-01T10:00:00Z" {
			t.Errorf("expected Date 2023-12-01T10:00:00Z, got %s", transaction.Date)
		}
		if transaction.AccountID != 123 {
			t.Errorf("expected AccountID 123, got %d", transaction.AccountID)
		}
		if transaction.CategoryID != 456 {
			t.Errorf("expected CategoryID 456, got %d", transaction.CategoryID)
		}
	})

	t.Run("success case: empty transactions list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PersonalAccountTransactions{
			Transactions: []PersonalAccountTransaction{},
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
		response, err := client.GetPersonalAccountTransactions(context.Background(), "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Transactions) != 0 {
			t.Fatalf("expected 0 transactions, got %d", len(response.Transactions))
		}
	})

	t.Run("success case: transactions list with pagination parameters", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PersonalAccountTransactions{
			Transactions: []PersonalAccountTransaction{
				{
					ID:         1048,
					Amount:     -5000.00,
					Date:       "2023-12-01T10:00:00Z",
					AccountID:  123,
					CategoryID: 456,
					Attributes: PersonalAccountTransactionAttributes{},
					CreatedAt:  "2023-12-01T09:00:00Z",
					UpdatedAt:  "2023-12-01T09:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/accounts/account_key_123/transactions.json" {
				t.Errorf("expected path /link/accounts/account_key_123/transactions.json, got %s", r.URL.Path)
			}
			expectedPage := "2"
			actualPage := r.URL.Query().Get("page")
			if actualPage != expectedPage {
				t.Errorf("expected page parameter %s, got %s", expectedPage, actualPage)
			}
			expectedPerPage := "100"
			actualPerPage := r.URL.Query().Get("per_page")
			if actualPerPage != expectedPerPage {
				t.Errorf("expected per_page parameter %s, got %s", expectedPerPage, actualPerPage)
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
		response, err := client.GetPersonalAccountTransactions(context.Background(), "account_key_123",
			WithPageForTransactions(2),
			WithPerPageForTransactions(100),
		)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Transactions) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(response.Transactions))
		}
	})

	t.Run("success case: transactions list with sort parameters", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PersonalAccountTransactions{
			Transactions: []PersonalAccountTransaction{
				{
					ID:         1048,
					Amount:     -5000.00,
					Date:       "2023-12-01T10:00:00Z",
					AccountID:  123,
					CategoryID: 456,
					Attributes: PersonalAccountTransactionAttributes{},
					CreatedAt:  "2023-12-01T09:00:00Z",
					UpdatedAt:  "2023-12-01T09:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			expectedSortKey := "date"
			actualSortKey := r.URL.Query().Get("sort_key")
			if actualSortKey != expectedSortKey {
				t.Errorf("expected sort_key parameter %s, got %s", expectedSortKey, actualSortKey)
			}
			expectedSortBy := "desc"
			actualSortBy := r.URL.Query().Get("sort_by")
			if actualSortBy != expectedSortBy {
				t.Errorf("expected sort_by parameter %s, got %s", expectedSortBy, actualSortBy)
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
		response, err := client.GetPersonalAccountTransactions(context.Background(), "account_key_123",
			WithSortKeyForTransactions("date"),
			WithSortByForTransactions("desc"),
		)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Transactions) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(response.Transactions))
		}
	})

	t.Run("success case: transactions list with since parameter", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PersonalAccountTransactions{
			Transactions: []PersonalAccountTransaction{
				{
					ID:         1048,
					Amount:     -5000.00,
					Date:       "2023-12-01T10:00:00Z",
					AccountID:  123,
					CategoryID: 456,
					Attributes: PersonalAccountTransactionAttributes{},
					CreatedAt:  "2023-12-01T09:00:00Z",
					UpdatedAt:  "2023-12-01T09:00:00Z",
				},
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
		response, err := client.GetPersonalAccountTransactions(context.Background(), "account_key_123",
			WithSinceForTransactions("2023-01-01"),
		)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Transactions) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(response.Transactions))
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
		_, err = client.GetPersonalAccountTransactions(context.Background(), "account_key_123")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "token refresh function is not set") && !strings.Contains(err.Error(), "refresh token") {
			t.Errorf("expected error about token refresh, got %v", err)
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
		_, err = client.GetPersonalAccountTransactions(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "account ID is required") {
			t.Errorf("expected error about account ID, got %v", err)
		}
	})

	t.Run("error case: returns error when sort_by is invalid", func(t *testing.T) {
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
		_, err = client.GetPersonalAccountTransactions(context.Background(), "account_key_123",
			WithSortByForTransactions("invalid"),
		)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "sort_by must be 'asc' or 'desc'") {
			t.Errorf("expected error about sort_by, got %v", err)
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
		_, err = client.GetPersonalAccountTransactions(context.Background(), "account_key_123",
			WithSinceForTransactions("2023/01/01"),
		)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "date must be in format YYYY-MM-DD") {
			t.Errorf("expected error about date format, got %v", err)
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
		_, err = client.GetPersonalAccountTransactions(context.Background(), accountID)
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

		setTestToken(client, "test-token")
		// nolint:staticcheck // passing nil context for testing purposes
		_, err = client.GetPersonalAccountTransactions(nil, accountID)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func int64Ptr(i int64) *int64 {
	return &i
}

func TestUpdatePersonalAccountTransaction(t *testing.T) {
	t.Parallel()

	t.Run("success case: transaction is updated correctly", func(t *testing.T) {
		t.Parallel()

		descriptionGuest := "新しいメモ"
		descriptionPretty := "マネーツリーによる補正"
		descriptionRaw := "生データ"
		categoryEntityKey := "category_key_123"

		expectedResponse := PersonalAccountTransaction{
			ID:                1337,
			Amount:            -5000.00,
			Date:              "2023-12-01T10:00:00Z",
			DescriptionGuest:  &descriptionGuest,
			DescriptionPretty: &descriptionPretty,
			DescriptionRaw:    &descriptionRaw,
			AccountID:         1048,
			CategoryID:        123,
			Attributes:        PersonalAccountTransactionAttributes{},
			CategoryEntityKey: &categoryEntityKey,
			CreatedAt:         "2023-12-01T09:00:00Z",
			UpdatedAt:         "2023-12-01T09:00:00Z",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("expected method %s, got %s", http.MethodPut, r.Method)
			}
			if r.URL.Path != "/link/accounts/account_key_123/transactions/1337.json" {
				t.Errorf("expected path /link/accounts/account_key_123/transactions/1337.json, got %s", r.URL.Path)
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				t.Errorf("expected Authorization header with Bearer prefix, got %s", authHeader)
			}
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			var req UpdatePersonalAccountTransactionRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}
			if req.DescriptionGuest == nil || *req.DescriptionGuest != "新しいメモ" {
				t.Errorf("expected DescriptionGuest '新しいメモ', got %v", req.DescriptionGuest)
			}
			if req.CategoryID == nil || *req.CategoryID != 123 {
				t.Errorf("expected CategoryID 123, got %v", req.CategoryID)
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

		request := &UpdatePersonalAccountTransactionRequest{
			DescriptionGuest: &descriptionGuest,
			CategoryID:       int64Ptr(123),
		}

		setTestToken(client, "test-access-token")
		response, err := client.UpdatePersonalAccountTransaction(context.Background(), "account_key_123", 1337, request)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.ID != 1337 {
			t.Errorf("expected ID 1337, got %d", response.ID)
		}
		if response.Amount != -5000.00 {
			t.Errorf("expected Amount -5000.00, got %f", response.Amount)
		}
		if response.AccountID != 1048 {
			t.Errorf("expected AccountID 1048, got %d", response.AccountID)
		}
		if response.CategoryID != 123 {
			t.Errorf("expected CategoryID 123, got %d", response.CategoryID)
		}
		if response.DescriptionGuest == nil || *response.DescriptionGuest != descriptionGuest {
			t.Errorf("expected DescriptionGuest %s, got %v", descriptionGuest, response.DescriptionGuest)
		}
	})

	t.Run("success case: update only description_guest", func(t *testing.T) {
		t.Parallel()

		descriptionGuest := "メモのみ更新"

		expectedResponse := PersonalAccountTransaction{
			ID:               1337,
			Amount:           -5000.00,
			Date:             "2023-12-01T10:00:00Z",
			DescriptionGuest: &descriptionGuest,
			AccountID:        1048,
			CategoryID:       456,
			Attributes:       PersonalAccountTransactionAttributes{},
			CreatedAt:        "2023-12-01T09:00:00Z",
			UpdatedAt:        "2023-12-01T09:00:00Z",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req UpdatePersonalAccountTransactionRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}
			if req.DescriptionGuest == nil || *req.DescriptionGuest != descriptionGuest {
				t.Errorf("expected DescriptionGuest %s, got %v", descriptionGuest, req.DescriptionGuest)
			}
			if req.CategoryID != nil {
				t.Errorf("expected CategoryID nil, got %v", req.CategoryID)
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

		request := &UpdatePersonalAccountTransactionRequest{
			DescriptionGuest: &descriptionGuest,
		}

		setTestToken(client, "test-access-token")
		response, err := client.UpdatePersonalAccountTransaction(context.Background(), "account_key_123", 1337, request)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.DescriptionGuest == nil || *response.DescriptionGuest != descriptionGuest {
			t.Errorf("expected DescriptionGuest %s, got %v", descriptionGuest, response.DescriptionGuest)
		}
	})

	t.Run("success case: update only category_id", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PersonalAccountTransaction{
			ID:         1337,
			Amount:     -5000.00,
			Date:       "2023-12-01T10:00:00Z",
			AccountID:  1048,
			CategoryID: 789,
			Attributes: PersonalAccountTransactionAttributes{},
			CreatedAt:  "2023-12-01T09:00:00Z",
			UpdatedAt:  "2023-12-01T09:00:00Z",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req UpdatePersonalAccountTransactionRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}
			if req.DescriptionGuest != nil {
				t.Errorf("expected DescriptionGuest nil, got %v", req.DescriptionGuest)
			}
			if req.CategoryID == nil || *req.CategoryID != 789 {
				t.Errorf("expected CategoryID 789, got %v", req.CategoryID)
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

		request := &UpdatePersonalAccountTransactionRequest{
			CategoryID: int64Ptr(789),
		}

		setTestToken(client, "test-access-token")
		response, err := client.UpdatePersonalAccountTransaction(context.Background(), "account_key_123", 1337, request)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.CategoryID != 789 {
			t.Errorf("expected CategoryID 789, got %d", response.CategoryID)
		}
	})

	t.Run("success case: update date and amount for manually entered account", func(t *testing.T) {
		t.Parallel()

		date := "2023-12-01T10:00:00Z"
		amount := -5000.00
		descriptionGuest := "手入力取引"

		expectedResponse := PersonalAccountTransaction{
			ID:               1337,
			Amount:           amount,
			Date:             date,
			DescriptionGuest: &descriptionGuest,
			AccountID:        1048,
			CategoryID:       456,
			Attributes:       PersonalAccountTransactionAttributes{},
			CreatedAt:        "2023-12-01T09:00:00Z",
			UpdatedAt:        "2023-12-01T09:00:00Z",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req UpdatePersonalAccountTransactionRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}
			if req.Date == nil || *req.Date != date {
				t.Errorf("expected Date %s, got %v", date, req.Date)
			}
			if req.Amount == nil || *req.Amount != amount {
				t.Errorf("expected Amount %f, got %v", amount, req.Amount)
			}
			if req.DescriptionGuest == nil || *req.DescriptionGuest != descriptionGuest {
				t.Errorf("expected DescriptionGuest %s, got %v", descriptionGuest, req.DescriptionGuest)
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

		request := &UpdatePersonalAccountTransactionRequest{
			Date:             &date,
			Amount:           &amount,
			DescriptionGuest: &descriptionGuest,
		}

		setTestToken(client, "test-access-token")
		response, err := client.UpdatePersonalAccountTransaction(context.Background(), "account_key_123", 1337, request)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.Date != date {
			t.Errorf("expected Date %s, got %s", date, response.Date)
		}
		if response.Amount != amount {
			t.Errorf("expected Amount %f, got %f", amount, response.Amount)
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

		request := &UpdatePersonalAccountTransactionRequest{
			DescriptionGuest: stringPtr("test"),
		}

		// Token is not set, so refreshToken should fail
		_, err = client.UpdatePersonalAccountTransaction(context.Background(), "account_key_123", 1337, request)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "token refresh function is not set") && !strings.Contains(err.Error(), "refresh token") {
			t.Errorf("expected error about token refresh, got %v", err)
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

		request := &UpdatePersonalAccountTransactionRequest{
			DescriptionGuest: stringPtr("test"),
		}

		setTestToken(client, "test-token")
		_, err = client.UpdatePersonalAccountTransaction(context.Background(), "", 1337, request)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "account ID is required") {
			t.Errorf("expected error about account ID, got %v", err)
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
		_, err = client.UpdatePersonalAccountTransaction(context.Background(), "account_key_123", 1337, nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "request cannot be nil") {
			t.Errorf("expected error about request, got %v", err)
		}
	})

	t.Run("error case: returns error when description_guest exceeds 255 characters", func(t *testing.T) {
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

		longDescription := strings.Repeat("a", 256)
		request := &UpdatePersonalAccountTransactionRequest{
			DescriptionGuest: &longDescription,
		}

		setTestToken(client, "test-token")
		_, err = client.UpdatePersonalAccountTransaction(context.Background(), "account_key_123", 1337, request)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "description_guest must be 255 characters or less") {
			t.Errorf("expected error about description_guest length, got %v", err)
		}
	})

	t.Run("error case: returns error when API returns an error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "invalid_request", "error_description": "Category ID does not exist."}`))
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

		request := &UpdatePersonalAccountTransactionRequest{
			CategoryID: int64Ptr(99999),
		}

		setTestToken(client, "test-token")
		_, err = client.UpdatePersonalAccountTransaction(context.Background(), "account_key_123", 1337, request)
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
		if !strings.Contains(err.Error(), "invalid_request") {
			t.Errorf("expected error about invalid_request, got %v", err)
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

		request := &UpdatePersonalAccountTransactionRequest{
			DescriptionGuest: stringPtr("test"),
		}

		setTestToken(client, "test-token")
		// nolint:staticcheck // passing nil context for testing purposes
		_, err = client.UpdatePersonalAccountTransaction(nil, "account_key_123", 1337, request)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}
