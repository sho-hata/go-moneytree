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

func TestGetPointAccounts(t *testing.T) {
	t.Parallel()

	t.Run("success case: point accounts list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		id1 := int64(123)
		id2 := int64(456)
		balance1 := float64Ptr(10000.50)
		balance2 := float64Ptr(5000.00)
		lastAggregatedAt := "2023-01-01T00:00:00Z"
		lastAggregatedSuccess := "2023-01-01T00:00:00Z"
		createdAt := "2023-01-01T00:00:00Z"
		updatedAt := "2023-01-01T00:00:00Z"

		expectedResponse := PointAccounts{
			PointAccounts: []PointAccount{
				{
					ID:                     id1,
					AccountGroup:           789,
					AccountType:            "point",
					Currency:               "JPY",
					InstitutionEntityKey:   "test_point_1",
					InstitutionAccountName: "ポイントカード",
					Nickname:               "ポイントカード",
					CurrentBalance:         balance1,
					AggregationState:       "success",
					AggregationStatus:      "success",
					LastAggregatedAt:       lastAggregatedAt,
					LastAggregatedSuccess:  stringPtr(lastAggregatedSuccess),
					CreatedAt:              createdAt,
					UpdatedAt:              updatedAt,
				},
				{
					ID:                     id2,
					AccountGroup:           789,
					AccountType:            "point",
					Currency:               "JPY",
					InstitutionEntityKey:   "test_point_2",
					InstitutionAccountName: "クレジットカードポイント",
					Nickname:               "クレジットカードポイント",
					CurrentBalance:         balance2,
					AggregationState:       "success",
					AggregationStatus:      "success",
					LastAggregatedAt:       lastAggregatedAt,
					LastAggregatedSuccess:  stringPtr(lastAggregatedSuccess),
					CreatedAt:              createdAt,
					UpdatedAt:              updatedAt,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/points/accounts.json" {
				t.Errorf("expected path /link/points/accounts.json, got %s", r.URL.Path)
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

		response, err := client.GetPointAccounts(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.PointAccounts) != 2 {
			t.Fatalf("expected 2 point accounts, got %d", len(response.PointAccounts))
		}

		account1 := response.PointAccounts[0]
		if account1.ID != expectedResponse.PointAccounts[0].ID {
			t.Errorf("expected ID %d, got %d", expectedResponse.PointAccounts[0].ID, account1.ID)
		}
		if account1.AccountGroup != expectedResponse.PointAccounts[0].AccountGroup {
			t.Errorf("expected AccountGroup %d, got %d", expectedResponse.PointAccounts[0].AccountGroup, account1.AccountGroup)
		}
		if account1.AccountType != expectedResponse.PointAccounts[0].AccountType {
			t.Errorf("expected AccountType %s, got %s", expectedResponse.PointAccounts[0].AccountType, account1.AccountType)
		}
		if account1.Currency != expectedResponse.PointAccounts[0].Currency {
			t.Errorf("expected Currency %s, got %s", expectedResponse.PointAccounts[0].Currency, account1.Currency)
		}
		if account1.InstitutionEntityKey != expectedResponse.PointAccounts[0].InstitutionEntityKey {
			t.Errorf("expected InstitutionEntityKey %s, got %s", expectedResponse.PointAccounts[0].InstitutionEntityKey, account1.InstitutionEntityKey)
		}
		if account1.InstitutionAccountName != expectedResponse.PointAccounts[0].InstitutionAccountName {
			t.Errorf("expected InstitutionAccountName %s, got %s", expectedResponse.PointAccounts[0].InstitutionAccountName, account1.InstitutionAccountName)
		}
		if account1.Nickname != expectedResponse.PointAccounts[0].Nickname {
			t.Errorf("expected Nickname %s, got %s", expectedResponse.PointAccounts[0].Nickname, account1.Nickname)
		}
		if account1.CurrentBalance == nil || *account1.CurrentBalance != *expectedResponse.PointAccounts[0].CurrentBalance {
			t.Errorf("expected CurrentBalance %v, got %v", *expectedResponse.PointAccounts[0].CurrentBalance, account1.CurrentBalance)
		}
		if account1.AggregationState != expectedResponse.PointAccounts[0].AggregationState {
			t.Errorf("expected AggregationState %s, got %s", expectedResponse.PointAccounts[0].AggregationState, account1.AggregationState)
		}
		if account1.AggregationStatus != expectedResponse.PointAccounts[0].AggregationStatus {
			t.Errorf("expected AggregationStatus %s, got %s", expectedResponse.PointAccounts[0].AggregationStatus, account1.AggregationStatus)
		}
		if account1.LastAggregatedAt != expectedResponse.PointAccounts[0].LastAggregatedAt {
			t.Errorf("expected LastAggregatedAt %s, got %s", expectedResponse.PointAccounts[0].LastAggregatedAt, account1.LastAggregatedAt)
		}
		if account1.LastAggregatedSuccess == nil || *account1.LastAggregatedSuccess != *expectedResponse.PointAccounts[0].LastAggregatedSuccess {
			t.Errorf("expected LastAggregatedSuccess %v, got %v", *expectedResponse.PointAccounts[0].LastAggregatedSuccess, account1.LastAggregatedSuccess)
		}
		if account1.CreatedAt != expectedResponse.PointAccounts[0].CreatedAt {
			t.Errorf("expected CreatedAt %s, got %s", expectedResponse.PointAccounts[0].CreatedAt, account1.CreatedAt)
		}
		if account1.UpdatedAt != expectedResponse.PointAccounts[0].UpdatedAt {
			t.Errorf("expected UpdatedAt %s, got %s", expectedResponse.PointAccounts[0].UpdatedAt, account1.UpdatedAt)
		}

		account2 := response.PointAccounts[1]
		if account2.ID != expectedResponse.PointAccounts[1].ID {
			t.Errorf("expected ID %d, got %d", expectedResponse.PointAccounts[1].ID, account2.ID)
		}
		if account2.AccountType != expectedResponse.PointAccounts[1].AccountType {
			t.Errorf("expected AccountType %s, got %s", expectedResponse.PointAccounts[1].AccountType, account2.AccountType)
		}
		if account2.CurrentBalance == nil || *account2.CurrentBalance != *expectedResponse.PointAccounts[1].CurrentBalance {
			t.Errorf("expected CurrentBalance %v, got %v", *expectedResponse.PointAccounts[1].CurrentBalance, account2.CurrentBalance)
		}
	})

	t.Run("success case: point accounts list with null balance", func(t *testing.T) {
		t.Parallel()

		lastAggregatedAt := "2023-01-01T00:00:00Z"
		createdAt := "2023-01-01T00:00:00Z"
		updatedAt := "2023-01-01T00:00:00Z"

		expectedResponse := PointAccounts{
			PointAccounts: []PointAccount{
				{
					ID:                     123,
					AccountGroup:           789,
					AccountType:            "point",
					Currency:               "JPY",
					InstitutionEntityKey:   "test_point_1",
					InstitutionAccountName: "ポイントカード",
					Nickname:               "ポイントカード",
					CurrentBalance:         nil,
					AggregationState:       "error",
					AggregationStatus:      "error.temporary",
					LastAggregatedAt:       lastAggregatedAt,
					LastAggregatedSuccess:  nil,
					CreatedAt:              createdAt,
					UpdatedAt:              updatedAt,
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

		response, err := client.GetPointAccounts(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.PointAccounts) != 1 {
			t.Fatalf("expected 1 point account, got %d", len(response.PointAccounts))
		}
		if response.PointAccounts[0].CurrentBalance != nil {
			t.Errorf("expected CurrentBalance nil, got %v", response.PointAccounts[0].CurrentBalance)
		}
		if response.PointAccounts[0].LastAggregatedSuccess != nil {
			t.Errorf("expected LastAggregatedSuccess nil, got %v", response.PointAccounts[0].LastAggregatedSuccess)
		}
	})

	t.Run("success case: point accounts list with pagination", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PointAccounts{
			PointAccounts: []PointAccount{
				{
					ID:                     123,
					AccountGroup:           789,
					AccountType:            "point",
					Currency:               "JPY",
					InstitutionEntityKey:   "test_point_1",
					InstitutionAccountName: "ポイントカード",
					Nickname:               "ポイントカード",
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
			if r.URL.Query().Get("per_page") != "100" {
				t.Errorf("expected per_page=100, got %s", r.URL.Query().Get("per_page"))
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

		response, err := client.GetPointAccounts(context.Background(), "test-access-token",
			WithPageForPointAccounts(2),
			WithPerPageForPointAccounts(100),
		)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.PointAccounts) != 1 {
			t.Fatalf("expected 1 point account, got %d", len(response.PointAccounts))
		}
	})

	t.Run("error case: access token is empty", func(t *testing.T) {
		t.Parallel()

		client := &Client{
			httpClient: http.DefaultClient,
			config: &Config{
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "jp-api-staging.getmoneytree.com",
				},
			},
		}

		_, err := client.GetPointAccounts(context.Background(), "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error message to contain 'access token is required', got %v", err)
		}
	})

	t.Run("error case: HTTP error response", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse := map[string]interface{}{
				"error":             "invalid_token",
				"error_description": "The access token provided is invalid.",
			}
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				t.Errorf("failed to encode error response: %v", err)
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

		_, err = client.GetPointAccounts(context.Background(), "invalid-token")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
		}
		if apiErr.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, apiErr.StatusCode)
		}
	})

	t.Run("error case: invalid JSON response", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("invalid json")); err != nil {
				t.Errorf("failed to write response: %v", err)
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

		_, err = client.GetPointAccounts(context.Background(), "test-access-token")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("success case: empty point accounts list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PointAccounts{
			PointAccounts: []PointAccount{},
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

		response, err := client.GetPointAccounts(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.PointAccounts) != 0 {
			t.Fatalf("expected 0 point accounts, got %d", len(response.PointAccounts))
		}
	})
}

func TestGetPointAccountTransactions(t *testing.T) {
	t.Parallel()

	t.Run("success case: transactions list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		descriptionGuest := "ポイント獲得"
		descriptionPretty := "ポイント獲得"
		descriptionRaw := "ポイント獲得"
		categoryEntityKey := "category_key_123"

		expectedResponse := PointAccountTransactions{
			Transactions: []PointAccountTransaction{
				{
					ID:                1048,
					Amount:            1000.00,
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
			if r.URL.Path != "/link/points/accounts/123/transactions.json" {
				t.Errorf("expected path /link/points/accounts/123/transactions.json, got %s", r.URL.Path)
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

		response, err := client.GetPointAccountTransactions(context.Background(), "test-access-token", 123)
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
		if transaction.Amount != 1000.00 {
			t.Errorf("expected Amount 1000.00, got %f", transaction.Amount)
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

		expectedResponse := PointAccountTransactions{
			Transactions: []PointAccountTransaction{},
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

		response, err := client.GetPointAccountTransactions(context.Background(), "test-access-token", 123)
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

		expectedResponse := PointAccountTransactions{
			Transactions: []PointAccountTransaction{
				{
					ID:         1048,
					Amount:     1000.00,
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
			if r.URL.Path != "/link/points/accounts/123/transactions.json" {
				t.Errorf("expected path /link/points/accounts/123/transactions.json, got %s", r.URL.Path)
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

		response, err := client.GetPointAccountTransactions(context.Background(), "test-access-token", 123,
			WithPageForPointAccountTransactions(2),
			WithPerPageForPointAccountTransactions(100),
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

		expectedResponse := PointAccountTransactions{
			Transactions: []PointAccountTransaction{
				{
					ID:         1048,
					Amount:     1000.00,
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

		response, err := client.GetPointAccountTransactions(context.Background(), "test-access-token", 123,
			WithSortKeyForPointAccountTransactions("date"),
			WithSortByForPointAccountTransactions("desc"),
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

		expectedResponse := PointAccountTransactions{
			Transactions: []PointAccountTransaction{
				{
					ID:         1048,
					Amount:     1000.00,
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

		response, err := client.GetPointAccountTransactions(context.Background(), "test-access-token", 123,
			WithSinceForPointAccountTransactions("2023-01-01"),
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

		_, err = client.GetPointAccountTransactions(context.Background(), "", 123)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error about access token, got %v", err)
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

		_, err = client.GetPointAccountTransactions(context.Background(), "test-token", 123,
			WithSortByForPointAccountTransactions("invalid"),
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

		_, err = client.GetPointAccountTransactions(context.Background(), "test-token", 123,
			WithSinceForPointAccountTransactions("2023/01/01"),
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

		accountID := int64(123)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte(`{"error": "invalid_token", "error_description": "The access token is invalid or expired."}`)); err != nil {
				t.Errorf("failed to write response: %v", err)
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

		_, err = client.GetPointAccountTransactions(context.Background(), "invalid-token", accountID)
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

		accountID := int64(123)

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
		_, err = client.GetPointAccountTransactions(nil, "test-token", accountID)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestGetPointExpirations(t *testing.T) {
	t.Parallel()

	t.Run("success case: point expirations list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PointExpirations{
			PointExpirations: []PointExpiration{
				{
					ID:               1,
					AccountID:        123,
					ExpirationAmount: 1000.00,
					ExpirationDate:   "2024-12-31",
					Date:             "2023-12-01T10:00:00Z",
				},
				{
					ID:               2,
					AccountID:        123,
					ExpirationAmount: 500.00,
					ExpirationDate:   "2025-01-15",
					Date:             "2023-12-01T10:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/points/accounts/123/expirations.json" {
				t.Errorf("expected path /link/points/accounts/123/expirations.json, got %s", r.URL.Path)
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

		response, err := client.GetPointExpirations(context.Background(), "test-access-token", 123)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.PointExpirations) != 2 {
			t.Fatalf("expected 2 point expirations, got %d", len(response.PointExpirations))
		}

		expiration1 := response.PointExpirations[0]
		if expiration1.ID != expectedResponse.PointExpirations[0].ID {
			t.Errorf("expected ID %d, got %d", expectedResponse.PointExpirations[0].ID, expiration1.ID)
		}
		if expiration1.AccountID != expectedResponse.PointExpirations[0].AccountID {
			t.Errorf("expected AccountID %d, got %d", expectedResponse.PointExpirations[0].AccountID, expiration1.AccountID)
		}
		if expiration1.ExpirationAmount != expectedResponse.PointExpirations[0].ExpirationAmount {
			t.Errorf("expected ExpirationAmount %f, got %f", expectedResponse.PointExpirations[0].ExpirationAmount, expiration1.ExpirationAmount)
		}
		if expiration1.ExpirationDate != expectedResponse.PointExpirations[0].ExpirationDate {
			t.Errorf("expected ExpirationDate %s, got %s", expectedResponse.PointExpirations[0].ExpirationDate, expiration1.ExpirationDate)
		}
		if expiration1.Date != expectedResponse.PointExpirations[0].Date {
			t.Errorf("expected Date %s, got %s", expectedResponse.PointExpirations[0].Date, expiration1.Date)
		}

		expiration2 := response.PointExpirations[1]
		if expiration2.ID != expectedResponse.PointExpirations[1].ID {
			t.Errorf("expected ID %d, got %d", expectedResponse.PointExpirations[1].ID, expiration2.ID)
		}
		if expiration2.ExpirationAmount != expectedResponse.PointExpirations[1].ExpirationAmount {
			t.Errorf("expected ExpirationAmount %f, got %f", expectedResponse.PointExpirations[1].ExpirationAmount, expiration2.ExpirationAmount)
		}
	})

	t.Run("success case: empty point expirations list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PointExpirations{
			PointExpirations: []PointExpiration{},
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

		response, err := client.GetPointExpirations(context.Background(), "test-access-token", 123)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.PointExpirations) != 0 {
			t.Fatalf("expected 0 point expirations, got %d", len(response.PointExpirations))
		}
	})

	t.Run("success case: point expirations list with pagination parameters", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PointExpirations{
			PointExpirations: []PointExpiration{
				{
					ID:               1,
					AccountID:        123,
					ExpirationAmount: 1000.00,
					ExpirationDate:   "2024-12-31",
					Date:             "2023-12-01T10:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/points/accounts/123/expirations.json" {
				t.Errorf("expected path /link/points/accounts/123/expirations.json, got %s", r.URL.Path)
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

		response, err := client.GetPointExpirations(context.Background(), "test-access-token", 123,
			WithPageForPointExpirations(2),
			WithPerPageForPointExpirations(100),
		)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.PointExpirations) != 1 {
			t.Fatalf("expected 1 point expiration, got %d", len(response.PointExpirations))
		}
	})

	t.Run("success case: point expirations list with since parameter", func(t *testing.T) {
		t.Parallel()

		expectedResponse := PointExpirations{
			PointExpirations: []PointExpiration{
				{
					ID:               1,
					AccountID:        123,
					ExpirationAmount: 1000.00,
					ExpirationDate:   "2024-12-31",
					Date:             "2023-12-01T10:00:00Z",
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

		response, err := client.GetPointExpirations(context.Background(), "test-access-token", 123,
			WithSinceForPointExpirations("2023-01-01"),
		)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.PointExpirations) != 1 {
			t.Fatalf("expected 1 point expiration, got %d", len(response.PointExpirations))
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

		_, err = client.GetPointExpirations(context.Background(), "", 123)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error about access token, got %v", err)
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

		_, err = client.GetPointExpirations(context.Background(), "test-token", 123,
			WithSinceForPointExpirations("2023/01/01"),
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

		accountID := int64(123)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte(`{"error": "invalid_token", "error_description": "The access token is invalid or expired."}`)); err != nil {
				t.Errorf("failed to write response: %v", err)
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

		_, err = client.GetPointExpirations(context.Background(), "invalid-token", accountID)
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

		accountID := int64(123)

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
		_, err = client.GetPointExpirations(nil, "test-token", accountID)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}
