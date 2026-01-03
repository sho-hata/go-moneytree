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

func TestGetInstitutions(t *testing.T) {
	t.Parallel()

	t.Run("success case: institutions list is retrieved correctly", func(t *testing.T) {
		t.Parallel()

		displayName1 := stringPtr("Test Bank 1")
		displayName2 := stringPtr("Test Bank 2")
		displayName3 := stringPtr("Test Bank 3")
		displayNameReading1 := stringPtr("テストバンク1")
		displayNameReading2 := stringPtr("テストバンク2")
		displayNameReading3 := stringPtr("テストバンク3")
		status1 := stringPtr("active")
		status2 := stringPtr("inactive")
		statusReason1 := stringPtr("maintenance")
		statusReason2 := stringPtr("legacy")
		loginURL1 := stringPtr("https://example.com/login")
		guidanceURL1 := stringPtr("https://example.com/guidance")

		expectedResponse := Institutions{
			Institutions: []Institution{
				{
					EntityKey:                "test_bank_1",
					InstitutionType:          "bank",
					DisplayName:              displayName1,
					DisplayNameReading:       displayNameReading1,
					Status:                   status1,
					StatusReason:             nil,
					LoginURL:                 loginURL1,
					GuidanceURL:              guidanceURL1,
					BillingGroup:             nil,
					Tags:                     []string{"bank", "individual"},
					DefaultAuthorizationType: 0,
				},
				{
					EntityKey:                "test_bank_2",
					InstitutionType:          "bank",
					DisplayName:              displayName2,
					DisplayNameReading:       displayNameReading2,
					Status:                   status2,
					StatusReason:             statusReason1,
					LoginURL:                 nil,
					GuidanceURL:              nil,
					BillingGroup:             stringPtr("2"),
					Tags:                     []string{"bank"},
					DefaultAuthorizationType: 1,
				},
				{
					EntityKey:                "test_bank_3",
					InstitutionType:          "bank",
					DisplayName:              displayName3,
					DisplayNameReading:       displayNameReading3,
					Status:                   status2,
					StatusReason:             statusReason2,
					LoginURL:                 nil,
					GuidanceURL:              nil,
					BillingGroup:             nil,
					Tags:                     []string{"bank", "legacy"},
					DefaultAuthorizationType: 0,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/institutions.json" {
				t.Errorf("expected path /link/institutions.json, got %s", r.URL.Path)
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
		response, err := client.GetInstitutions(context.Background())
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Institutions) != 3 {
			t.Fatalf("expected 3 institutions, got %d", len(response.Institutions))
		}

		inst1 := response.Institutions[0]
		if inst1.EntityKey != expectedResponse.Institutions[0].EntityKey {
			t.Errorf("expected EntityKey %s, got %s", expectedResponse.Institutions[0].EntityKey, inst1.EntityKey)
		}
		if inst1.InstitutionType != expectedResponse.Institutions[0].InstitutionType {
			t.Errorf("expected InstitutionType %s, got %s", expectedResponse.Institutions[0].InstitutionType, inst1.InstitutionType)
		}
		if inst1.DisplayName == nil || *inst1.DisplayName != *expectedResponse.Institutions[0].DisplayName {
			t.Errorf("expected DisplayName %s, got %v", *expectedResponse.Institutions[0].DisplayName, inst1.DisplayName)
		}
		if inst1.Status == nil || *inst1.Status != *expectedResponse.Institutions[0].Status {
			t.Errorf("expected Status %s, got %v", *expectedResponse.Institutions[0].Status, inst1.Status)
		}
		if inst1.StatusReason != nil {
			t.Errorf("expected StatusReason nil, got %v", inst1.StatusReason)
		}
		if len(inst1.Tags) != len(expectedResponse.Institutions[0].Tags) {
			t.Errorf("expected Tags length %d, got %d", len(expectedResponse.Institutions[0].Tags), len(inst1.Tags))
		}

		inst2 := response.Institutions[1]
		if inst2.StatusReason == nil {
			t.Error("expected StatusReason, got nil")
		} else if *inst2.StatusReason != *expectedResponse.Institutions[1].StatusReason {
			t.Errorf("expected StatusReason %s, got %s", *expectedResponse.Institutions[1].StatusReason, *inst2.StatusReason)
		}
	})

	t.Run("success case: institutions list with since parameter", func(t *testing.T) {
		t.Parallel()

		sinceTime := "2023-01-01"
		displayName := stringPtr("Test Bank 1")
		status := stringPtr("active")

		expectedResponse := Institutions{
			Institutions: []Institution{
				{
					EntityKey:                "test_bank_1",
					InstitutionType:          "bank",
					DisplayName:              displayName,
					DisplayNameReading:       displayName,
					Status:                   status,
					StatusReason:             nil,
					Tags:                     []string{"bank"},
					DefaultAuthorizationType: 0,
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected method %s, got %s", http.MethodGet, r.Method)
			}
			if r.URL.Path != "/link/institutions.json" {
				t.Errorf("expected path /link/institutions.json, got %s", r.URL.Path)
			}
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
		response, err := client.GetInstitutions(context.Background(), WithSince(sinceTime))
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Institutions) != 1 {
			t.Fatalf("expected 1 institution, got %d", len(response.Institutions))
		}
	})

	t.Run("success case: empty institutions list", func(t *testing.T) {
		t.Parallel()

		expectedResponse := Institutions{
			Institutions: []Institution{},
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
		response, err := client.GetInstitutions(context.Background())
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}

		if response == nil {
			t.Fatal("expected response, got nil")
		}
		if len(response.Institutions) != 0 {
			t.Fatalf("expected 0 institutions, got %d", len(response.Institutions))
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
		_, err = client.GetInstitutions(context.Background())
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
		_, err = client.GetInstitutions(context.Background())
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
		_, err = client.GetInstitutions(nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "context must be non-nil") {
			t.Errorf("expected error about context, got %v", err)
		}
	})
}

func TestWithSince_InvalidDateFormat(t *testing.T) {
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
				_, err = client.GetInstitutions(context.Background(),
					WithSince(invalidDate),
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

				opt := WithSince(validDate)
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
				_, err = client.GetInstitutions(context.Background(),
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
