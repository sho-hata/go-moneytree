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

func TestGetInvestmentAccounts(t *testing.T) {
	t.Parallel()

	t.Run("success case: accounts list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		id1 := int64(123)
		id2 := int64(456)
		balance1 := float64Ptr(1000000.50)
		balance2 := float64Ptr(500000.00)
		lastAggregatedAt := "2023-01-01T00:00:00Z"
		createdAt := "2023-01-01T00:00:00Z"
		updatedAt := "2023-01-01T00:00:00Z"

		expectedResponse := InvestmentAccounts{
			Accounts: []InvestmentAccount{
				{
					ID:                      id1,
					AccountKey:              "investment_account_key_1",
					AccountGroup:            789,
					AccountSubtype:          "brokerage",
					AccountType:             "stock",
					Currency:                "JPY",
					InstitutionEntityKey:    "test_brokerage_1",
					InstitutionID:           1,
					InstitutionAccountName:  "証券口座",
					InstitutionAccountNumber: stringPtr("1234567"),
					Nickname:                "証券口座",
					BranchName:              stringPtr("本店"),
					BranchCode:              stringPtr("001"),
					AggregationState:        "success",
					AggregationStatus:       "success",
					LastAggregatedAt:        lastAggregatedAt,
					LastAggregatedSuccess:   stringPtr(lastAggregatedAt),
					CurrentBalance:          balance1,
					CurrentBalanceInBase:    balance1,
					CurrentBalanceDataSource: stringPtr("institution"),
					CreatedAt:               createdAt,
					UpdatedAt:               updatedAt,
				},
				{
					ID:                      id2,
					AccountKey:              "investment_account_key_2",
					AccountGroup:            789,
					AccountSubtype:          "defined_contribution_pension",
					AccountType:             "pension",
					Currency:                "JPY",
					InstitutionEntityKey:    "test_pension_1",
					InstitutionID:           2,
					InstitutionAccountName:  "確定拠出年金",
					InstitutionAccountNumber: stringPtr("9876543"),
					Nickname:                "確定拠出年金",
					BranchName:              nil,
					BranchCode:              nil,
					AggregationState:        "success",
					AggregationStatus:       "success",
					LastAggregatedAt:        lastAggregatedAt,
					LastAggregatedSuccess:   stringPtr(lastAggregatedAt),
					CurrentBalance:          balance2,
					CurrentBalanceInBase:    balance2,
					CurrentBalanceDataSource: stringPtr("institution"),
					CreatedAt:               createdAt,
					UpdatedAt:               updatedAt,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/investments/accounts.json" {
				t.Errorf("expected path /link/investments/accounts.json, got %s", r.URL.Path)
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

		response, err := client.GetInvestmentAccounts(context.Background(), "test-access-token")
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
		if account1.ID != expectedResponse.Accounts[0].ID {
			t.Errorf("expected ID %d, got %d", expectedResponse.Accounts[0].ID, account1.ID)
		}
		if account1.AccountKey != expectedResponse.Accounts[0].AccountKey {
			t.Errorf("expected AccountKey %s, got %s", expectedResponse.Accounts[0].AccountKey, account1.AccountKey)
		}
		if account1.AccountSubtype != expectedResponse.Accounts[0].AccountSubtype {
			t.Errorf("expected AccountSubtype %s, got %s", expectedResponse.Accounts[0].AccountSubtype, account1.AccountSubtype)
		}
		if account1.AccountType != expectedResponse.Accounts[0].AccountType {
			t.Errorf("expected AccountType %s, got %s", expectedResponse.Accounts[0].AccountType, account1.AccountType)
		}
		if account1.Currency != expectedResponse.Accounts[0].Currency {
			t.Errorf("expected Currency %s, got %s", expectedResponse.Accounts[0].Currency, account1.Currency)
		}
		if account1.Nickname != expectedResponse.Accounts[0].Nickname {
			t.Errorf("expected Nickname %s, got %s", expectedResponse.Accounts[0].Nickname, account1.Nickname)
		}
		if account1.CurrentBalance == nil || *account1.CurrentBalance != *expectedResponse.Accounts[0].CurrentBalance {
			t.Errorf("expected CurrentBalance %v, got %v", *expectedResponse.Accounts[0].CurrentBalance, account1.CurrentBalance)
		}
		if account1.AccountGroup != expectedResponse.Accounts[0].AccountGroup {
			t.Errorf("expected AccountGroup %d, got %d", expectedResponse.Accounts[0].AccountGroup, account1.AccountGroup)
		}
		if account1.InstitutionEntityKey != expectedResponse.Accounts[0].InstitutionEntityKey {
			t.Errorf("expected InstitutionEntityKey %s, got %s", expectedResponse.Accounts[0].InstitutionEntityKey, account1.InstitutionEntityKey)
		}

		account2 := response.Accounts[1]
		if account2.AccountSubtype != expectedResponse.Accounts[1].AccountSubtype {
			t.Errorf("expected AccountSubtype %s, got %s", expectedResponse.Accounts[1].AccountSubtype, account2.AccountSubtype)
		}
		if account2.AccountType != expectedResponse.Accounts[1].AccountType {
			t.Errorf("expected AccountType %s, got %s", expectedResponse.Accounts[1].AccountType, account2.AccountType)
		}
		if account2.CurrentBalance == nil || *account2.CurrentBalance != *expectedResponse.Accounts[1].CurrentBalance {
			t.Errorf("expected CurrentBalance %v, got %v", *expectedResponse.Accounts[1].CurrentBalance, account2.CurrentBalance)
		}
	})

	t.Run("success case: accounts list with null balance", func(t *testing.T) {
		t.Parallel()

		accountKey := "investment_account_key_1"
		lastAggregatedAt := "2023-01-01T00:00:00Z"
		createdAt := "2023-01-01T00:00:00Z"
		updatedAt := "2023-01-01T00:00:00Z"

		expectedResponse := InvestmentAccounts{
			Accounts: []InvestmentAccount{
				{
					ID:                      123,
					AccountKey:              accountKey,
					AccountGroup:            789,
					AccountSubtype:          "brokerage",
					AccountType:             "stock",
					Currency:                "JPY",
					InstitutionEntityKey:    "test_brokerage_1",
					InstitutionID:           1,
					InstitutionAccountName:  "証券口座",
					InstitutionAccountNumber: nil,
					Nickname:                "証券口座",
					BranchName:              nil,
					BranchCode:              nil,
					AggregationState:        "success",
					AggregationStatus:       "success",
					LastAggregatedAt:        lastAggregatedAt,
					LastAggregatedSuccess:   nil,
					CurrentBalance:          nil,
					CurrentBalanceInBase:    nil,
					CurrentBalanceDataSource: nil,
					CreatedAt:               createdAt,
					UpdatedAt:               updatedAt,
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

		response, err := client.GetInvestmentAccounts(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Accounts) != 1 {
			t.Fatalf("expected 1 account, got %d", len(response.Accounts))
		}
		if response.Accounts[0].CurrentBalance != nil {
			t.Errorf("expected CurrentBalance nil, got %v", response.Accounts[0].CurrentBalance)
		}
	})

	t.Run("success case: accounts list with pagination", func(t *testing.T) {
		t.Parallel()

		expectedResponse := InvestmentAccounts{
			Accounts: []InvestmentAccount{
				{
					ID:                     123,
					AccountKey:             "investment_account_key_1",
					AccountGroup:           789,
					AccountSubtype:         "brokerage",
					AccountType:            "stock",
					Currency:               "JPY",
					InstitutionEntityKey:   "test_brokerage_1",
					InstitutionID:          1,
					InstitutionAccountName: "証券口座",
					Nickname:               "証券口座",
					AggregationState:       "success",
					AggregationStatus:      "success",
					LastAggregatedAt:       "2023-01-01T00:00:00Z",
					CreatedAt:              "2023-01-01T00:00:00Z",
					UpdatedAt:              "2023-01-01T00:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("page") != "2" {
				t.Errorf("expected page=2, got %s", r.URL.Query().Get("page"))
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

		response, err := client.GetInvestmentAccounts(context.Background(), "test-access-token",
			WithPageForInvestmentAccounts(2),
		)
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

	t.Run("success case: empty accounts list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := InvestmentAccounts{
			Accounts: []InvestmentAccount{},
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

		response, err := client.GetInvestmentAccounts(context.Background(), "test-access-token")
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

		_, err = client.GetInvestmentAccounts(context.Background(), "")
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
			_, _ = w.Write([]byte(`{"error": "invalid_token", "error_description": "The access token provided is invalid."}`))
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

		_, err = client.GetInvestmentAccounts(context.Background(), "invalid-token")
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
		_, err = client.GetInvestmentAccounts(nil, "test-token")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

