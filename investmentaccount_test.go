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
					ID:                       id1,
					AccountKey:               "investment_account_key_1",
					AccountGroup:             789,
					AccountSubtype:           "brokerage",
					AccountType:              "stock",
					Currency:                 "JPY",
					InstitutionEntityKey:     "test_brokerage_1",
					InstitutionID:            1,
					InstitutionAccountName:   "証券口座",
					InstitutionAccountNumber: stringPtr("1234567"),
					Nickname:                 "証券口座",
					BranchName:               stringPtr("本店"),
					BranchCode:               stringPtr("001"),
					AggregationState:         "success",
					AggregationStatus:        "success",
					LastAggregatedAt:         lastAggregatedAt,
					LastAggregatedSuccess:    stringPtr(lastAggregatedAt),
					CurrentBalance:           balance1,
					CurrentBalanceInBase:     balance1,
					CurrentBalanceDataSource: stringPtr("institution"),
					CreatedAt:                createdAt,
					UpdatedAt:                updatedAt,
				},
				{
					ID:                       id2,
					AccountKey:               "investment_account_key_2",
					AccountGroup:             789,
					AccountSubtype:           "defined_contribution_pension",
					AccountType:              "pension",
					Currency:                 "JPY",
					InstitutionEntityKey:     "test_pension_1",
					InstitutionID:            2,
					InstitutionAccountName:   "確定拠出年金",
					InstitutionAccountNumber: stringPtr("9876543"),
					Nickname:                 "確定拠出年金",
					BranchName:               nil,
					BranchCode:               nil,
					AggregationState:         "success",
					AggregationStatus:        "success",
					LastAggregatedAt:         lastAggregatedAt,
					LastAggregatedSuccess:    stringPtr(lastAggregatedAt),
					CurrentBalance:           balance2,
					CurrentBalanceInBase:     balance2,
					CurrentBalanceDataSource: stringPtr("institution"),
					CreatedAt:                createdAt,
					UpdatedAt:                updatedAt,
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

		setTestToken(client, "test-access-token")
		response, err := client.GetInvestmentAccounts(context.Background())
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
					ID:                       123,
					AccountKey:               accountKey,
					AccountGroup:             789,
					AccountSubtype:           "brokerage",
					AccountType:              "stock",
					Currency:                 "JPY",
					InstitutionEntityKey:     "test_brokerage_1",
					InstitutionID:            1,
					InstitutionAccountName:   "証券口座",
					InstitutionAccountNumber: nil,
					Nickname:                 "証券口座",
					BranchName:               nil,
					BranchCode:               nil,
					AggregationState:         "success",
					AggregationStatus:        "success",
					LastAggregatedAt:         lastAggregatedAt,
					LastAggregatedSuccess:    nil,
					CurrentBalance:           nil,
					CurrentBalanceInBase:     nil,
					CurrentBalanceDataSource: nil,
					CreatedAt:                createdAt,
					UpdatedAt:                updatedAt,
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
		response, err := client.GetInvestmentAccounts(context.Background())
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

		setTestToken(client, "test-access-token")
		response, err := client.GetInvestmentAccounts(context.Background(),
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

		setTestToken(client, "test-access-token")
		response, err := client.GetInvestmentAccounts(context.Background())
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

		// Token is not set, so refreshToken should fail
		_, err = client.GetInvestmentAccounts(context.Background())
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

		setTestToken(client, "invalid-token")
		_, err = client.GetInvestmentAccounts(context.Background())
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

		setTestToken(client, "test-token")
		// nolint:staticcheck // passing nil context for testing purposes
		_, err = client.GetInvestmentAccounts(nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestGetInvestmentPositions(t *testing.T) {
	t.Parallel()

	t.Run("success case: positions list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		id1 := int64(123)
		id2 := int64(456)
		marketValue1 := 1000000.50
		marketValue2 := 500000.00
		acquisitionValue1 := 950000.00
		profit1 := 50000.50
		quantity1 := 100.0
		quantity2 := 50.0
		tickerCode1 := "7203"
		nameClean1 := "トヨタ自動車"
		nameClean2 := "日本株式インデックス"
		taxType1 := []string{"ippan"}
		taxSubType1 := "ippan"
		createdAt := "2023-01-01T00:00:00Z"
		updatedAt := "2023-01-01T00:00:00Z"

		expectedResponse := InvestmentPositions{
			Positions: []InvestmentPosition{
				{
					ID:               id1,
					Date:             "2023-01-01",
					AssetClass:       "stock",
					AssetSubclass:    stringPtr("common_stock"),
					TickerCode:       &tickerCode1,
					NameRaw:          stringPtr("トヨタ自動車株式会社"),
					NameClean:        &nameClean1,
					Currency:         "JPY",
					TaxType:          taxType1,
					TaxSubType:       &taxSubType1,
					MarketValue:      marketValue1,
					Value:            marketValue1,
					AcquisitionValue: &acquisitionValue1,
					CostBasis:        &acquisitionValue1,
					Profit:           &profit1,
					Quantity:         &quantity1,
					CreatedAt:        createdAt,
					UpdatedAt:        updatedAt,
				},
				{
					ID:               id2,
					Date:             "2023-01-01",
					AssetClass:       "investment_trust",
					AssetSubclass:    nil,
					TickerCode:       nil,
					NameRaw:          stringPtr("日本株式インデックスファンド"),
					NameClean:        &nameClean2,
					Currency:         "JPY",
					TaxType:          []string{"NISA"},
					TaxSubType:       stringPtr("tsumitate"),
					MarketValue:      marketValue2,
					Value:            marketValue2,
					AcquisitionValue: nil,
					CostBasis:        nil,
					Profit:           nil,
					Quantity:         &quantity2,
					CreatedAt:        createdAt,
					UpdatedAt:        updatedAt,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/investments/accounts/account_key_123/positions.json" {
				t.Errorf("expected path /link/investments/accounts/account_key_123/positions.json, got %s", r.URL.Path)
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
		response, err := client.GetInvestmentPositions(context.Background(), "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Positions) != 2 {
			t.Fatalf("expected 2 positions, got %d", len(response.Positions))
		}

		position1 := response.Positions[0]
		if position1.ID != expectedResponse.Positions[0].ID {
			t.Errorf("expected ID %d, got %d", expectedResponse.Positions[0].ID, position1.ID)
		}
		if position1.AssetClass != expectedResponse.Positions[0].AssetClass {
			t.Errorf("expected AssetClass %s, got %s", expectedResponse.Positions[0].AssetClass, position1.AssetClass)
		}
		if position1.MarketValue != expectedResponse.Positions[0].MarketValue {
			t.Errorf("expected MarketValue %f, got %f", expectedResponse.Positions[0].MarketValue, position1.MarketValue)
		}
		if position1.NameClean == nil || *position1.NameClean != *expectedResponse.Positions[0].NameClean {
			t.Errorf("expected NameClean %s, got %v", *expectedResponse.Positions[0].NameClean, position1.NameClean)
		}
		if position1.TickerCode == nil || *position1.TickerCode != *expectedResponse.Positions[0].TickerCode {
			t.Errorf("expected TickerCode %s, got %v", *expectedResponse.Positions[0].TickerCode, position1.TickerCode)
		}
		if position1.Quantity == nil || *position1.Quantity != *expectedResponse.Positions[0].Quantity {
			t.Errorf("expected Quantity %f, got %v", *expectedResponse.Positions[0].Quantity, position1.Quantity)
		}

		position2 := response.Positions[1]
		if position2.AssetClass != expectedResponse.Positions[1].AssetClass {
			t.Errorf("expected AssetClass %s, got %s", expectedResponse.Positions[1].AssetClass, position2.AssetClass)
		}
		if len(position2.TaxType) != len(expectedResponse.Positions[1].TaxType) {
			t.Errorf("expected TaxType length %d, got %d", len(expectedResponse.Positions[1].TaxType), len(position2.TaxType))
		}
	})

	t.Run("success case: positions list with null optional fields", func(t *testing.T) {
		t.Parallel()

		expectedResponse := InvestmentPositions{
			Positions: []InvestmentPosition{
				{
					ID:               123,
					Date:             "2023-01-01",
					AssetClass:       "cash",
					AssetSubclass:    nil,
					TickerCode:       nil,
					NameRaw:          nil,
					NameClean:        nil,
					Currency:         "JPY",
					TaxType:          nil,
					TaxSubType:       nil,
					MarketValue:      100000.00,
					Value:            100000.00,
					AcquisitionValue: nil,
					CostBasis:        nil,
					Profit:           nil,
					Quantity:         nil,
					CreatedAt:        "2023-01-01T00:00:00Z",
					UpdatedAt:        "2023-01-01T00:00:00Z",
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
		response, err := client.GetInvestmentPositions(context.Background(), "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Positions) != 1 {
			t.Fatalf("expected 1 position, got %d", len(response.Positions))
		}
		if response.Positions[0].TickerCode != nil {
			t.Errorf("expected TickerCode nil, got %v", response.Positions[0].TickerCode)
		}
		if response.Positions[0].Quantity != nil {
			t.Errorf("expected Quantity nil, got %v", response.Positions[0].Quantity)
		}
	})

	t.Run("success case: positions list with pagination", func(t *testing.T) {
		t.Parallel()

		expectedResponse := InvestmentPositions{
			Positions: []InvestmentPosition{
				{
					ID:          123,
					Date:        "2023-01-01",
					AssetClass:  "stock",
					Currency:    "JPY",
					MarketValue: 1000000.00,
					Value:       1000000.00,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T00:00:00Z",
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

		setTestToken(client, "test-access-token")
		response, err := client.GetInvestmentPositions(context.Background(), "account_key_123",
			WithPageForInvestmentPositions(2),
		)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Positions) != 1 {
			t.Fatalf("expected 1 position, got %d", len(response.Positions))
		}
	})

	t.Run("success case: empty positions list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := InvestmentPositions{
			Positions: []InvestmentPosition{},
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
		response, err := client.GetInvestmentPositions(context.Background(), "account_key_123")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Positions) != 0 {
			t.Fatalf("expected 0 positions, got %d", len(response.Positions))
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
		_, err = client.GetInvestmentPositions(context.Background(), "account_key_123")
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
		_, err = client.GetInvestmentPositions(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "account ID is required") {
			t.Errorf("expected error about account ID, got %v", err)
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

		setTestToken(client, "invalid-token")
		_, err = client.GetInvestmentPositions(context.Background(), "account_key_123")
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

		setTestToken(client, "test-token")
		// nolint:staticcheck // passing nil context for testing purposes
		_, err = client.GetInvestmentPositions(nil, "account_key_123")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestGetInvestmentAccountTransactions(t *testing.T) {
	t.Parallel()

	t.Run("success case: transactions list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		descriptionGuest := "投資取引"
		descriptionPretty := "投資取引（補正済み）"
		descriptionRaw := "投資取引（生データ）"
		categoryEntityKey := "category_key_123"

		expectedResponse := InvestmentAccountTransactions{
			Transactions: []InvestmentAccountTransaction{
				{
					ID:                1048,
					Amount:            -100000.00,
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
			if r.URL.Path != "/link/investments/accounts/account_key_123/transactions.json" {
				t.Errorf("expected path /link/investments/accounts/account_key_123/transactions.json, got %s", r.URL.Path)
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
		response, err := client.GetInvestmentAccountTransactions(context.Background(), "account_key_123")
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
		if transaction.Amount != -100000.00 {
			t.Errorf("expected Amount -100000.00, got %f", transaction.Amount)
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

		expectedResponse := InvestmentAccountTransactions{
			Transactions: []InvestmentAccountTransaction{},
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
		response, err := client.GetInvestmentAccountTransactions(context.Background(), "account_key_123")
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

		expectedResponse := InvestmentAccountTransactions{
			Transactions: []InvestmentAccountTransaction{
				{
					ID:         1048,
					Amount:     -100000.00,
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
		response, err := client.GetInvestmentAccountTransactions(context.Background(), "account_key_123",
			WithPageForInvestmentTransactions(2),
			WithPerPageForInvestmentTransactions(100),
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

		expectedResponse := InvestmentAccountTransactions{
			Transactions: []InvestmentAccountTransaction{
				{
					ID:         1048,
					Amount:     -100000.00,
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
		response, err := client.GetInvestmentAccountTransactions(context.Background(), "account_key_123",
			WithSortKeyForInvestmentTransactions("date"),
			WithSortByForInvestmentTransactions("desc"),
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

		expectedResponse := InvestmentAccountTransactions{
			Transactions: []InvestmentAccountTransaction{
				{
					ID:         1048,
					Amount:     -100000.00,
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
		response, err := client.GetInvestmentAccountTransactions(context.Background(), "account_key_123",
			WithSinceForInvestmentTransactions("2023-01-01"),
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
		_, err = client.GetInvestmentAccountTransactions(context.Background(), "account_key_123")
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
		_, err = client.GetInvestmentAccountTransactions(context.Background(), "")
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
		_, err = client.GetInvestmentAccountTransactions(context.Background(), "account_key_123",
			WithSortByForInvestmentTransactions("invalid"),
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
		_, err = client.GetInvestmentAccountTransactions(context.Background(), "account_key_123",
			WithSinceForInvestmentTransactions("2023/01/01"),
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
		_, err = client.GetInvestmentAccountTransactions(context.Background(), accountID)
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
		_, err = client.GetInvestmentAccountTransactions(nil, accountID)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}
