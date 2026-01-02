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

		response, err := client.GetPersonalAccounts(context.Background(), "test-access-token")
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

		response, err := client.GetPersonalAccounts(context.Background(), "test-access-token")
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

		response, err := client.GetPersonalAccounts(context.Background(), "test-access-token")
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

		response, err := client.GetPersonalAccounts(context.Background(), "test-access-token", WithPage(2))
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

		response, err := client.GetPersonalAccounts(context.Background(), "test-access-token", WithPerPage(100))
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

		response, err := client.GetPersonalAccounts(context.Background(), "test-access-token", WithPage(3), WithPerPage(50))
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

		_, err = client.GetPersonalAccounts(context.Background(), "")
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

		_, err = client.GetPersonalAccounts(context.Background(), "invalid-token")
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
		_, err = client.GetPersonalAccounts(nil, "test-token")
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

				_, err = client.GetPersonalAccountBalances(context.Background(), "test-token", "account_key_123",
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
					config: &Config{
						BaseURL: baseURL,
					},
				}

				// オプション関数を適用してもエラーが発生しないことを確認
				// （実際のAPI呼び出しは失敗するが、日付フォーマットエラーではない）
				_, err = client.GetPersonalAccountBalances(context.Background(), "test-token", "account_key_123",
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

		response, err := client.GetPersonalAccountBalances(context.Background(), "test-access-token", accountID)
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

		response, err := client.GetPersonalAccountBalances(context.Background(), "test-access-token", accountID, WithSinceForBalances(sinceTime))
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

		response, err := client.GetPersonalAccountBalances(context.Background(), "test-access-token", accountID,
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

		response, err := client.GetPersonalAccountBalances(context.Background(), "test-access-token", accountID)
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

		_, err = client.GetPersonalAccountBalances(context.Background(), "", "account_key_123")
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

		_, err = client.GetPersonalAccountBalances(context.Background(), "test-token", "")
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

		_, err = client.GetPersonalAccountBalances(context.Background(), "invalid-token", accountID)
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
		_, err = client.GetPersonalAccountBalances(nil, "test-token", accountID)
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
					ID:             1048,
					AccountID:      123,
					Date:           "2023-12-01",
					PurchaseDate:   &purchaseDate,
					MaturityDate:   &maturityDate,
					NameRaw:        &nameRaw,
					NameClean:      &nameClean,
					Value:          1050000.00,
					CostBasis:      1000000.00,
					InterestRate:   0.25,
					Currency:       "JPY",
					TermLengthYear: &termLengthYear,
					TermLengthMonth: &termLengthMonth,
					TermLengthDay:  &termLengthDay,
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

		response, err := client.GetTermDeposits(context.Background(), "test-access-token", "account_key_123")
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

		response, err := client.GetTermDeposits(context.Background(), "test-access-token", "account_key_123")
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

		response, err := client.GetTermDeposits(context.Background(), "test-access-token", "account_key_123", WithPageForTermDeposits(2))
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

		_, err = client.GetTermDeposits(context.Background(), "", "account_key_123")
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

		_, err = client.GetTermDeposits(context.Background(), "test-token", "")
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

		_, err = client.GetTermDeposits(context.Background(), "invalid-token", accountID)
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
		_, err = client.GetTermDeposits(nil, "test-token", accountID)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}
