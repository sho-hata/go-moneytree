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

func TestGetCorporateAccounts(t *testing.T) {
	t.Parallel()

	t.Run("success case: accounts list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		id1 := int64(123)
		id2 := int64(456)
		balance1 := float64Ptr(100000.50)
		balance2 := float64Ptr(-5000.00)
		lastAggregatedAt := "2023-01-01T00:00:00Z"
		createdAt := "2023-01-01T00:00:00Z"
		updatedAt := "2023-01-01T00:00:00Z"

		expectedResponse := CorporateAccounts{
			Accounts: []CorporateAccount{
				{
					ID:                      id1,
					AccountKey:              "account_key_1",
					AccountGroup:            789,
					AccountSubtype:          "savings",
					AccountType:             "bank",
					Currency:                "JPY",
					InstitutionEntityKey:    "test_bank_1",
					InstitutionID:           1,
					InstitutionAccountName:  "普通預金",
					InstitutionAccountNumber: stringPtr("1234567"),
					Nickname:                "普通預金",
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
					AccountKey:              "account_key_2",
					AccountGroup:            789,
					AccountSubtype:          "credit_card",
					AccountType:             "credit_card",
					Currency:                "JPY",
					InstitutionEntityKey:    "test_bank_2",
					InstitutionID:           2,
					InstitutionAccountName:  "クレジットカード",
					InstitutionAccountNumber: stringPtr("****1234"),
					Nickname:                "クレジットカード",
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
			if r.URL.Path != "/link/corporate/accounts.json" {
				t.Errorf("expected path /link/corporate/accounts.json, got %s", r.URL.Path)
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

		response, err := client.GetCorporateAccounts(context.Background(), "test-access-token")
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

		accountKey := "account_key_1"
		lastAggregatedAt := "2023-01-01T00:00:00Z"
		createdAt := "2023-01-01T00:00:00Z"
		updatedAt := "2023-01-01T00:00:00Z"

		expectedResponse := CorporateAccounts{
			Accounts: []CorporateAccount{
				{
					ID:                      123,
					AccountKey:              accountKey,
					AccountGroup:            789,
					AccountSubtype:          "savings",
					AccountType:             "bank",
					Currency:                "JPY",
					InstitutionEntityKey:    "test_bank_1",
					InstitutionID:           1,
					InstitutionAccountName:  "普通預金",
					InstitutionAccountNumber: nil,
					Nickname:                "普通預金",
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

		response, err := client.GetCorporateAccounts(context.Background(), "test-access-token")
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
		if account.CurrentBalance != nil {
			t.Errorf("expected CurrentBalance nil, got %v", account.CurrentBalance)
		}
		if account.LastAggregatedSuccess != nil {
			t.Errorf("expected LastAggregatedSuccess nil, got %v", account.LastAggregatedSuccess)
		}
	})

	t.Run("success case: empty accounts list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := CorporateAccounts{
			Accounts: []CorporateAccount{},
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

		response, err := client.GetCorporateAccounts(context.Background(), "test-access-token")
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

		balance := float64Ptr(100000.50)
		lastAggregatedAt := "2023-01-01T00:00:00Z"
		createdAt := "2023-01-01T00:00:00Z"
		updatedAt := "2023-01-01T00:00:00Z"

		expectedResponse := CorporateAccounts{
			Accounts: []CorporateAccount{
				{
					ID:                      123,
					AccountKey:              "account_key_1",
					AccountGroup:            789,
					AccountSubtype:          "savings",
					AccountType:             "bank",
					Currency:                "JPY",
					InstitutionEntityKey:    "test_bank_1",
					InstitutionID:           1,
					InstitutionAccountName:  "普通預金",
					InstitutionAccountNumber: stringPtr("1234567"),
					Nickname:                "普通預金",
					BranchName:              stringPtr("本店"),
					BranchCode:              stringPtr("001"),
					AggregationState:        "success",
					AggregationStatus:       "success",
					LastAggregatedAt:        lastAggregatedAt,
					LastAggregatedSuccess:   stringPtr(lastAggregatedAt),
					CurrentBalance:          balance,
					CurrentBalanceInBase:    balance,
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
			if r.URL.Path != "/link/corporate/accounts.json" {
				t.Errorf("expected path /link/corporate/accounts.json, got %s", r.URL.Path)
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

		response, err := client.GetCorporateAccounts(context.Background(), "test-access-token", WithPageForCorporateAccounts(2))
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


	t.Run("success case: accounts list with account_attributes", func(t *testing.T) {
		t.Parallel()

		balance := float64Ptr(100000.50)
		lastAggregatedAt := "2023-01-01T00:00:00Z"
		createdAt := "2023-01-01T00:00:00Z"
		updatedAt := "2023-01-01T00:00:00Z"
		accountHolderNameKatakanaRaw := "カ）テストカイシャ"
		accountHolderNameKatakanaZengin := "カ）テストカイシャ"

		expectedResponse := CorporateAccounts{
			Accounts: []CorporateAccount{
				{
					ID:                      123,
					AccountKey:              "account_key_1",
					AccountGroup:            789,
					AccountSubtype:          "savings",
					AccountType:             "bank",
					Currency:                "JPY",
					InstitutionEntityKey:    "test_bank_1",
					InstitutionID:           1,
					InstitutionAccountName:  "普通預金",
					InstitutionAccountNumber: stringPtr("1234567"),
					Nickname:                "普通預金",
					BranchName:              stringPtr("本店"),
					BranchCode:              stringPtr("001"),
					AggregationState:        "success",
					AggregationStatus:       "success",
					LastAggregatedAt:        lastAggregatedAt,
					LastAggregatedSuccess:   stringPtr(lastAggregatedAt),
					CurrentBalance:          balance,
					CurrentBalanceInBase:    balance,
					CurrentBalanceDataSource: stringPtr("institution"),
					CreatedAt:               createdAt,
					UpdatedAt:               updatedAt,
					AccountAttributes: &CorporateAccountAttributes{
						AccountHolderNameKatakanaRaw:     &accountHolderNameKatakanaRaw,
						AccountHolderNameKatakanaZengin:  &accountHolderNameKatakanaZengin,
					},
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

		response, err := client.GetCorporateAccounts(context.Background(), "test-access-token")
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
		if account.AccountAttributes == nil {
			t.Fatal("expected AccountAttributes, got nil")
		}
		if account.AccountAttributes.AccountHolderNameKatakanaRaw == nil || *account.AccountAttributes.AccountHolderNameKatakanaRaw != accountHolderNameKatakanaRaw {
			t.Errorf("expected AccountHolderNameKatakanaRaw %s, got %v", accountHolderNameKatakanaRaw, account.AccountAttributes.AccountHolderNameKatakanaRaw)
		}
		if account.AccountAttributes.AccountHolderNameKatakanaZengin == nil || *account.AccountAttributes.AccountHolderNameKatakanaZengin != accountHolderNameKatakanaZengin {
			t.Errorf("expected AccountHolderNameKatakanaZengin %s, got %v", accountHolderNameKatakanaZengin, account.AccountAttributes.AccountHolderNameKatakanaZengin)
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

		_, err = client.GetCorporateAccounts(context.Background(), "")
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

		_, err = client.GetCorporateAccounts(context.Background(), "invalid-token")
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
		_, err = client.GetCorporateAccounts(nil, "test-token")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

