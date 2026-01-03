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

func TestGetCategories(t *testing.T) {
	t.Parallel()

	t.Run("success case: categories list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		entityKey1 := stringPtr("food")
		entityKey2 := stringPtr("transportation")
		categoryType1 := stringPtr("expense")
		categoryType2 := stringPtr("expense")
		parentID1 := int64Ptr(0)

		expectedResponse := Categories{
			Categories: []Category{
				{
					ID:          1,
					EntityKey:   entityKey1,
					CategoryType: categoryType1,
					Name:        "食費",
					ParentID:    parentID1,
					IsSystem:    true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T00:00:00Z",
				},
				{
					ID:          2,
					EntityKey:   entityKey2,
					CategoryType: categoryType2,
					Name:        "交通費",
					ParentID:    nil,
					IsSystem:    true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T00:00:00Z",
				},
				{
					ID:          3,
					EntityKey:   nil,
					CategoryType: nil,
					Name:        "ユーザー作成カテゴリー",
					ParentID:    nil,
					IsSystem:    false,
					CreatedAt:   "2023-01-02T00:00:00Z",
					UpdatedAt:   "2023-01-02T00:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/categories.json" {
				t.Errorf("expected path /link/categories.json, got %s", r.URL.Path)
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

		response, err := client.GetCategories(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Categories) != 3 {
			t.Fatalf("expected 3 categories, got %d", len(response.Categories))
		}

		category1 := response.Categories[0]
		if category1.ID != 1 {
			t.Errorf("expected ID 1, got %d", category1.ID)
		}
		if category1.EntityKey == nil || *category1.EntityKey != *entityKey1 {
			t.Errorf("expected EntityKey %s, got %v", *entityKey1, category1.EntityKey)
		}
		if category1.CategoryType == nil || *category1.CategoryType != "expense" {
			t.Errorf("expected CategoryType 'expense', got %v", category1.CategoryType)
		}
		if !category1.IsSystem {
			t.Errorf("expected IsSystem true, got %v", category1.IsSystem)
		}
		if category1.Name != "食費" {
			t.Errorf("expected Name '食費', got %s", category1.Name)
		}
		if category1.CreatedAt == "" {
			t.Error("expected CreatedAt, got empty")
		}
		if category1.UpdatedAt == "" {
			t.Error("expected UpdatedAt, got empty")
		}

		category2 := response.Categories[1]
		if category2.EntityKey == nil || *category2.EntityKey != *entityKey2 {
			t.Errorf("expected EntityKey %s, got %v", *entityKey2, category2.EntityKey)
		}
		if !category2.IsSystem {
			t.Errorf("expected IsSystem true, got %v", category2.IsSystem)
		}

		category3 := response.Categories[2]
		if category3.EntityKey != nil {
			t.Errorf("expected EntityKey nil for user-created category, got %v", category3.EntityKey)
		}
		if category3.IsSystem {
			t.Errorf("expected IsSystem false for user-created category, got %v", category3.IsSystem)
		}
	})

	t.Run("success case: empty categories list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := Categories{
			Categories: []Category{},
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

		response, err := client.GetCategories(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Categories) != 0 {
			t.Fatalf("expected 0 categories, got %d", len(response.Categories))
		}
	})

	t.Run("success case: categories with null entity_key", func(t *testing.T) {
		t.Parallel()

		entityKey1 := stringPtr("food")
		categoryType1 := stringPtr("expense")

		expectedResponse := Categories{
			Categories: []Category{
				{
					ID:          1,
					EntityKey:   entityKey1,
					CategoryType: categoryType1,
					Name:        "食費",
					ParentID:    nil,
					IsSystem:    true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T00:00:00Z",
				},
				{
					ID:          2,
					EntityKey:   nil,
					CategoryType: nil,
					Name:        "ユーザー作成カテゴリー",
					ParentID:    nil,
					IsSystem:    false,
					CreatedAt:   "2023-01-02T00:00:00Z",
					UpdatedAt:   "2023-01-02T00:00:00Z",
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

		response, err := client.GetCategories(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Categories) != 2 {
			t.Fatalf("expected 2 categories, got %d", len(response.Categories))
		}

		if response.Categories[1].EntityKey != nil {
			t.Errorf("expected EntityKey nil for user-created category, got %v", response.Categories[1].EntityKey)
		}
		if response.Categories[1].IsSystem {
			t.Errorf("expected IsSystem false for user-created category, got %v", response.Categories[1].IsSystem)
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

		_, err = client.GetCategories(context.Background(), "")
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
			_, _ = w.Write([]byte(`{"error": "invalid_token", "error_description": "The access token is invalid."}`))
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

		_, err = client.GetCategories(context.Background(), "test-token")
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
		_, err = client.GetCategories(nil, "test-token") //nolint:staticcheck
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})

	t.Run("success case: categories with pagination", func(t *testing.T) {
		t.Parallel()

		entityKey1 := stringPtr("food")
		categoryType1 := stringPtr("expense")

		expectedResponse := Categories{
			Categories: []Category{
				{
					ID:          1,
					EntityKey:   entityKey1,
					CategoryType: categoryType1,
					Name:        "食費",
					ParentID:    nil,
					IsSystem:    true,
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

		response, err := client.GetCategories(context.Background(), "test-access-token", WithPageForCategories(2))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Categories) != 1 {
			t.Fatalf("expected 1 category, got %d", len(response.Categories))
		}
	})

	t.Run("success case: categories with locale", func(t *testing.T) {
		t.Parallel()

		entityKey1 := stringPtr("food")
		categoryType1 := stringPtr("expense")

		expectedResponse := Categories{
			Categories: []Category{
				{
					ID:          1,
					EntityKey:   entityKey1,
					CategoryType: categoryType1,
					Name:        "Food",
					ParentID:    nil,
					IsSystem:    true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T00:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("locale") != "en" {
				t.Errorf("expected locale=en, got %s", r.URL.Query().Get("locale"))
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

		response, err := client.GetCategories(context.Background(), "test-access-token", WithLocale("en"))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Categories) != 1 {
			t.Fatalf("expected 1 category, got %d", len(response.Categories))
		}
		if response.Categories[0].Name != "Food" {
			t.Errorf("expected Name 'Food', got %s", response.Categories[0].Name)
		}
	})

	t.Run("success case: categories with pagination and locale", func(t *testing.T) {
		t.Parallel()

		entityKey1 := stringPtr("food")
		categoryType1 := stringPtr("expense")

		expectedResponse := Categories{
			Categories: []Category{
				{
					ID:          1,
					EntityKey:   entityKey1,
					CategoryType: categoryType1,
					Name:        "Food",
					ParentID:    nil,
					IsSystem:    true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T00:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("page") != "1" {
				t.Errorf("expected page=1, got %s", r.URL.Query().Get("page"))
			}
			if r.URL.Query().Get("locale") != "ja" {
				t.Errorf("expected locale=ja, got %s", r.URL.Query().Get("locale"))
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

		response, err := client.GetCategories(context.Background(), "test-access-token",
			WithPageForCategories(1),
			WithLocale("ja"),
		)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
	})

	t.Run("error case: returns error when locale is invalid", func(t *testing.T) {
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

		_, err = client.GetCategories(context.Background(), "test-token", WithLocale("fr"))
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "locale must be either 'en' or 'ja'") {
			t.Errorf("expected error about locale, got %v", err)
		}
	})

	t.Run("error case: invalid JSON response", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("invalid json"))
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

		_, err = client.GetCategories(context.Background(), "test-token")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestCreateCategory(t *testing.T) {
	t.Parallel()

	t.Run("success case: category is created correctly", func(t *testing.T) {
		t.Parallel()

		expectedResponse := Category{
			ID:          123,
			EntityKey:   nil,
			CategoryType: nil,
			Name:        "新しいカテゴリー",
			ParentID:    int64Ptr(0),
			IsSystem:    false,
			CreatedAt:   "2023-01-01T00:00:00Z",
			UpdatedAt:   "2023-01-01T00:00:00Z",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("expected method %s, got %s", http.MethodPost, r.Method)
			}
			if r.URL.Path != "/link/categories.json" {
				t.Errorf("expected path /link/categories.json, got %s", r.URL.Path)
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				t.Errorf("expected Authorization header with Bearer prefix, got %s", authHeader)
			}
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			var req CreateCategoryRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}
			if req.Name != "新しいカテゴリー" {
				t.Errorf("expected Name '新しいカテゴリー', got %s", req.Name)
			}
			if req.ParentID != 0 {
				t.Errorf("expected ParentID 0, got %d", req.ParentID)
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

		request := &CreateCategoryRequest{
			Name:     "新しいカテゴリー",
			ParentID: 0,
		}

		response, err := client.CreateCategory(context.Background(), "test-access-token", request)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.ID != 123 {
			t.Errorf("expected ID 123, got %d", response.ID)
		}
		if response.Name != "新しいカテゴリー" {
			t.Errorf("expected Name '新しいカテゴリー', got %s", response.Name)
		}
		if response.ParentID == nil || *response.ParentID != 0 {
			t.Errorf("expected ParentID 0, got %v", response.ParentID)
		}
		if response.IsSystem {
			t.Errorf("expected IsSystem false, got %v", response.IsSystem)
		}
		if response.CreatedAt == "" {
			t.Error("expected CreatedAt, got empty")
		}
		if response.UpdatedAt == "" {
			t.Error("expected UpdatedAt, got empty")
		}
	})

	t.Run("success case: category is created with parent_id", func(t *testing.T) {
		t.Parallel()

		parentID := int64(10)
		expectedResponse := Category{
			ID:          456,
			EntityKey:   nil,
			CategoryType: nil,
			Name:        "サブカテゴリー",
			ParentID:    &parentID,
			IsSystem:    false,
			CreatedAt:   "2023-01-01T00:00:00Z",
			UpdatedAt:   "2023-01-01T00:00:00Z",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req CreateCategoryRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}
			if req.ParentID != parentID {
				t.Errorf("expected ParentID %d, got %d", parentID, req.ParentID)
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

		request := &CreateCategoryRequest{
			Name:     "サブカテゴリー",
			ParentID: parentID,
		}

		response, err := client.CreateCategory(context.Background(), "test-access-token", request)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.ParentID == nil || *response.ParentID != parentID {
			t.Errorf("expected ParentID %d, got %v", parentID, response.ParentID)
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

		request := &CreateCategoryRequest{
			Name:     "テストカテゴリー",
			ParentID: 0,
		}

		_, err = client.CreateCategory(context.Background(), "", request)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error about access token, got %v", err)
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

		_, err = client.CreateCategory(context.Background(), "test-token", nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "request cannot be nil") {
			t.Errorf("expected error about request, got %v", err)
		}
	})

	t.Run("error case: returns error when name is empty", func(t *testing.T) {
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

		request := &CreateCategoryRequest{
			Name:     "",
			ParentID: 0,
		}

		_, err = client.CreateCategory(context.Background(), "test-token", request)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "name is required") {
			t.Errorf("expected error about name, got %v", err)
		}
	})

	t.Run("error case: returns error when API returns an error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "invalid_request", "error_description": "Parent category does not exist."}`))
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

		request := &CreateCategoryRequest{
			Name:     "テストカテゴリー",
			ParentID: 99999,
		}

		_, err = client.CreateCategory(context.Background(), "test-token", request)
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

		request := &CreateCategoryRequest{
			Name:     "テストカテゴリー",
			ParentID: 0,
		}

		// nolint:staticcheck // passing nil context for testing purposes
		_, err = client.CreateCategory(nil, "test-token", request) //nolint:staticcheck
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestGetCategory(t *testing.T) {
	t.Parallel()

	t.Run("success case: category is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		categoryID := int64(1048)
		entityKey := stringPtr("food")
		categoryType := stringPtr("expense")

		expectedResponse := Category{
			ID:          categoryID,
			EntityKey:   entityKey,
			CategoryType: categoryType,
			Name:        "食費",
			ParentID:    nil,
			IsSystem:    true,
			CreatedAt:   "2023-01-01T00:00:00Z",
			UpdatedAt:   "2023-01-01T00:00:00Z",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			expectedPath := fmt.Sprintf("/link/categories/%d.json", categoryID)
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

		response, err := client.GetCategory(context.Background(), "test-access-token", categoryID)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.ID != categoryID {
			t.Errorf("expected ID %d, got %d", categoryID, response.ID)
		}
		if response.Name != "食費" {
			t.Errorf("expected Name '食費', got %s", response.Name)
		}
		if response.EntityKey == nil || *response.EntityKey != *entityKey {
			t.Errorf("expected EntityKey %s, got %v", *entityKey, response.EntityKey)
		}
		if response.CategoryType == nil || *response.CategoryType != "expense" {
			t.Errorf("expected CategoryType 'expense', got %v", response.CategoryType)
		}
		if !response.IsSystem {
			t.Errorf("expected IsSystem true, got %v", response.IsSystem)
		}
		if response.CreatedAt == "" {
			t.Error("expected CreatedAt, got empty")
		}
		if response.UpdatedAt == "" {
			t.Error("expected UpdatedAt, got empty")
		}
	})

	t.Run("success case: category with null entity_key", func(t *testing.T) {
		t.Parallel()

		categoryID := int64(123)

		expectedResponse := Category{
			ID:          categoryID,
			EntityKey:   nil,
			CategoryType: nil,
			Name:        "ユーザー作成カテゴリー",
			ParentID:    nil,
			IsSystem:    false,
			CreatedAt:   "2023-01-02T00:00:00Z",
			UpdatedAt:   "2023-01-02T00:00:00Z",
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

		response, err := client.GetCategory(context.Background(), "test-access-token", categoryID)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.EntityKey != nil {
			t.Errorf("expected EntityKey nil for user-created category, got %v", response.EntityKey)
		}
		if response.IsSystem {
			t.Errorf("expected IsSystem false for user-created category, got %v", response.IsSystem)
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

		_, err = client.GetCategory(context.Background(), "", 1048)
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
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "not_found", "error_description": "Category not found."}`))
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

		_, err = client.GetCategory(context.Background(), "test-token", 99999)
		if err == nil {
			t.Error("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Errorf("expected APIError, got %T", err)
		}
		if apiErr.StatusCode != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, apiErr.StatusCode)
		}
		if !strings.Contains(err.Error(), "not_found") {
			t.Errorf("expected error about not_found, got %v", err)
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
		_, err = client.GetCategory(nil, "test-token", 1048) //nolint:staticcheck
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})

	t.Run("error case: invalid JSON response", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("invalid json"))
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

		_, err = client.GetCategory(context.Background(), "test-token", 1048)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUpdateCategory(t *testing.T) {
	t.Parallel()

	t.Run("success case: category is updated correctly", func(t *testing.T) {
		t.Parallel()

		categoryID := int64(123)
		expectedResponse := Category{
			ID:          categoryID,
			EntityKey:   nil,
			CategoryType: nil,
			Name:        "更新されたカテゴリー名",
			ParentID:    int64Ptr(0),
			IsSystem:    false,
			CreatedAt:   "2023-01-01T00:00:00Z",
			UpdatedAt:   "2023-01-02T00:00:00Z",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("expected method %s, got %s", http.MethodPut, r.Method)
			}
			expectedPath := fmt.Sprintf("/link/categories/%d.json", categoryID)
			if r.URL.Path != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
			}
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				t.Errorf("expected Authorization header with Bearer prefix, got %s", authHeader)
			}
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("expected Content-Type application/json, got %s", contentType)
			}

			var req UpdateCategoryRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}
			if req.Name != "更新されたカテゴリー名" {
				t.Errorf("expected Name '更新されたカテゴリー名', got %s", req.Name)
			}
			if req.ParentID != 0 {
				t.Errorf("expected ParentID 0, got %d", req.ParentID)
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

		request := &UpdateCategoryRequest{
			Name:     "更新されたカテゴリー名",
			ParentID: 0,
		}

		response, err := client.UpdateCategory(context.Background(), "test-access-token", categoryID, request)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.ID != categoryID {
			t.Errorf("expected ID %d, got %d", categoryID, response.ID)
		}
		if response.Name != "更新されたカテゴリー名" {
			t.Errorf("expected Name '更新されたカテゴリー名', got %s", response.Name)
		}
		if response.ParentID == nil || *response.ParentID != 0 {
			t.Errorf("expected ParentID 0, got %v", response.ParentID)
		}
		if response.UpdatedAt == "" {
			t.Error("expected UpdatedAt, got empty")
		}
	})

	t.Run("success case: category is updated with parent_id", func(t *testing.T) {
		t.Parallel()

		categoryID := int64(456)
		parentID := int64(10)
		expectedResponse := Category{
			ID:          categoryID,
			EntityKey:   nil,
			CategoryType: nil,
			Name:        "サブカテゴリー",
			ParentID:    &parentID,
			IsSystem:    false,
			CreatedAt:   "2023-01-01T00:00:00Z",
			UpdatedAt:   "2023-01-02T00:00:00Z",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req UpdateCategoryRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request: %v", err)
			}
			if req.ParentID != parentID {
				t.Errorf("expected ParentID %d, got %d", parentID, req.ParentID)
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

		request := &UpdateCategoryRequest{
			Name:     "サブカテゴリー",
			ParentID: parentID,
		}

		response, err := client.UpdateCategory(context.Background(), "test-access-token", categoryID, request)
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if response.ParentID == nil || *response.ParentID != parentID {
			t.Errorf("expected ParentID %d, got %v", parentID, response.ParentID)
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

		request := &UpdateCategoryRequest{
			Name:     "テストカテゴリー",
			ParentID: 0,
		}

		_, err = client.UpdateCategory(context.Background(), "", 123, request)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error about access token, got %v", err)
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

		_, err = client.UpdateCategory(context.Background(), "test-token", 123, nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "request cannot be nil") {
			t.Errorf("expected error about request, got %v", err)
		}
	})

	t.Run("error case: returns error when name is empty", func(t *testing.T) {
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

		request := &UpdateCategoryRequest{
			Name:     "",
			ParentID: 0,
		}

		_, err = client.UpdateCategory(context.Background(), "test-token", 123, request)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "name is required") {
			t.Errorf("expected error about name, got %v", err)
		}
	})

	t.Run("error case: returns error when API returns an error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error": "invalid_request", "error_description": "Category not found or cannot be updated."}`))
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

		request := &UpdateCategoryRequest{
			Name:     "テストカテゴリー",
			ParentID: 0,
		}

		_, err = client.UpdateCategory(context.Background(), "test-token", 99999, request)
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

		request := &UpdateCategoryRequest{
			Name:     "テストカテゴリー",
			ParentID: 0,
		}

		// nolint:staticcheck // passing nil context for testing purposes
		_, err = client.UpdateCategory(nil, "test-token", 123, request) //nolint:staticcheck
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestDeleteCategory(t *testing.T) {
	t.Parallel()

	t.Run("success case: category is deleted correctly", func(t *testing.T) {
		t.Parallel()

		categoryID := int64(123)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("expected method %s, got %s", http.MethodDelete, r.Method)
			}
			expectedPath := fmt.Sprintf("/link/categories/%d.json", categoryID)
			if r.URL.Path != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
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

		err = client.DeleteCategory(context.Background(), "test-access-token", categoryID)
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

		err = client.DeleteCategory(context.Background(), "", 123)
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
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "not_found", "error_description": "Category not found or cannot be deleted."}`))
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

		err = client.DeleteCategory(context.Background(), "test-token", 99999)
		if err == nil {
			t.Error("expected error, got nil")
		}

		var apiErr *APIError
		if !errors.As(err, &apiErr) {
			t.Errorf("expected APIError, got %T", err)
		}
		if apiErr.StatusCode != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, apiErr.StatusCode)
		}
		if !strings.Contains(err.Error(), "not_found") {
			t.Errorf("expected error about not_found, got %v", err)
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
		err = client.DeleteCategory(nil, "test-token", 123) //nolint:staticcheck
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestGetSystemCategories(t *testing.T) {
	t.Parallel()

	t.Run("success case: system categories list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		entityKey1 := stringPtr("food")
		entityKey2 := stringPtr("transportation")
		categoryType1 := stringPtr("expense")
		categoryType2 := stringPtr("expense")
		parentID1 := int64Ptr(0)

		expectedResponse := Categories{
			Categories: []Category{
				{
					ID:          1,
					EntityKey:   entityKey1,
					CategoryType: categoryType1,
					Name:        "食費",
					ParentID:    parentID1,
					IsSystem:    true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T00:00:00Z",
				},
				{
					ID:          2,
					EntityKey:   entityKey2,
					CategoryType: categoryType2,
					Name:        "交通費",
					ParentID:    nil,
					IsSystem:    true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T00:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/categories/system.json" {
				t.Errorf("expected path /link/categories/system.json, got %s", r.URL.Path)
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

		response, err := client.GetSystemCategories(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Categories) != 2 {
			t.Fatalf("expected 2 categories, got %d", len(response.Categories))
		}

		category1 := response.Categories[0]
		if category1.ID != 1 {
			t.Errorf("expected ID 1, got %d", category1.ID)
		}
		if category1.EntityKey == nil || *category1.EntityKey != *entityKey1 {
			t.Errorf("expected EntityKey %s, got %v", *entityKey1, category1.EntityKey)
		}
		if category1.CategoryType == nil || *category1.CategoryType != "expense" {
			t.Errorf("expected CategoryType 'expense', got %v", category1.CategoryType)
		}
		if !category1.IsSystem {
			t.Errorf("expected IsSystem true, got %v", category1.IsSystem)
		}
		if category1.Name != "食費" {
			t.Errorf("expected Name '食費', got %s", category1.Name)
		}

		category2 := response.Categories[1]
		if category2.EntityKey == nil || *category2.EntityKey != *entityKey2 {
			t.Errorf("expected EntityKey %s, got %v", *entityKey2, category2.EntityKey)
		}
		if !category2.IsSystem {
			t.Errorf("expected IsSystem true, got %v", category2.IsSystem)
		}
	})

	t.Run("success case: empty system categories list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := Categories{
			Categories: []Category{},
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

		response, err := client.GetSystemCategories(context.Background(), "test-access-token")
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Categories) != 0 {
			t.Fatalf("expected 0 categories, got %d", len(response.Categories))
		}
	})

	t.Run("success case: system categories with pagination", func(t *testing.T) {
		t.Parallel()

		entityKey1 := stringPtr("food")
		categoryType1 := stringPtr("expense")

		expectedResponse := Categories{
			Categories: []Category{
				{
					ID:          1,
					EntityKey:   entityKey1,
					CategoryType: categoryType1,
					Name:        "食費",
					ParentID:    nil,
					IsSystem:    true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T00:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/link/categories/system.json" {
				t.Errorf("expected path /link/categories/system.json, got %s", r.URL.Path)
			}
			if r.URL.RawQuery != "page=2" {
				t.Errorf("expected query page=2, got %s", r.URL.RawQuery)
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

		response, err := client.GetSystemCategories(context.Background(), "test-access-token", WithPageForCategories(2))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Categories) != 1 {
			t.Fatalf("expected 1 category, got %d", len(response.Categories))
		}
	})

	t.Run("success case: system categories with locale", func(t *testing.T) {
		t.Parallel()

		entityKey1 := stringPtr("food")
		categoryType1 := stringPtr("expense")

		expectedResponse := Categories{
			Categories: []Category{
				{
					ID:          1,
					EntityKey:   entityKey1,
					CategoryType: categoryType1,
					Name:        "Food",
					ParentID:    nil,
					IsSystem:    true,
					CreatedAt:   "2023-01-01T00:00:00Z",
					UpdatedAt:   "2023-01-01T00:00:00Z",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/link/categories/system.json" {
				t.Errorf("expected path /link/categories/system.json, got %s", r.URL.Path)
			}
			if r.URL.RawQuery != "locale=en" {
				t.Errorf("expected query locale=en, got %s", r.URL.RawQuery)
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

		response, err := client.GetSystemCategories(context.Background(), "test-access-token", WithLocale("en"))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Categories) != 1 {
			t.Fatalf("expected 1 category, got %d", len(response.Categories))
		}
		if response.Categories[0].Name != "Food" {
			t.Errorf("expected Name 'Food', got %s", response.Categories[0].Name)
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

		_, err = client.GetSystemCategories(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "access token is required") {
			t.Errorf("expected error about access token, got %v", err)
		}
	})

	t.Run("error case: returns error when locale is invalid", func(t *testing.T) {
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

		_, err = client.GetSystemCategories(context.Background(), "test-token", WithLocale("fr"))
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "locale must be either 'en' or 'ja'") {
			t.Errorf("expected error about locale, got %v", err)
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

		_, err = client.GetSystemCategories(context.Background(), "invalid-token")
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

	t.Run("error case: returns error when invalid JSON response", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("invalid json"))
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

		_, err = client.GetSystemCategories(context.Background(), "test-token")
		if err == nil {
			t.Error("expected error, got nil")
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
		_, err = client.GetSystemCategories(nil, "test-token") //nolint:staticcheck
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

