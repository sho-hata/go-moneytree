package moneytree

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// PointAccount represents a point account returned by the Moneytree LINK API.
// Point accounts include various point systems such as credit card points, shopping points, etc.
type PointAccount struct {
	// ID is the unique ID of the point account assigned by Moneytree.
	ID int64 `json:"id"`
	// AccountGroup is the unique ID for the financial service registration group.
	// This value corresponds to the account_group value in account group information API.
	// Account groups are displayed as one user on the financial institution's website.
	AccountGroup int64 `json:"account_group"`
	// AccountType describes the type of account.
	// For point accounts, this value is always "point".
	AccountType string `json:"account_type"`
	// Currency is the currency code based on ISO4217.
	// For point accounts, this is usually the currency representation of the point system.
	Currency string `json:"currency"`
	// InstitutionEntityKey is the key that identifies the financial service.
	// The name that can be displayed to customers can be obtained via the Financial Institution List API.
	// Use this instead of InstitutionID, as the institution_entity_key may change due to financial institution mergers.
	// It is recommended to store a new value each time the API is read if a new value exists.
	InstitutionEntityKey string `json:"institution_entity_key"`
	// InstitutionAccountName is the account name of the financial institution.
	// Example: "積立貯金" (Savings deposit).
	InstitutionAccountName string `json:"institution_account_name"`
	// Nickname is the name given to the account by the user within Moneytree.
	// If the user has not set a nickname individually, the same value as InstitutionAccountName will be returned.
	Nickname string `json:"nickname"`
	// CurrentBalance is the current point balance.
	// This value is null if the balance cannot be retrieved.
	CurrentBalance *float64 `json:"current_balance"`
	// AggregationState is the status of the latest data acquisition.
	// Possible values: "success", "running", "error".
	AggregationState string `json:"aggregation_state"`
	// AggregationStatus is the status of the latest data acquisition.
	// This provides more detailed content than AggregationState.
	// Possible values: "success", "running.auth", "running.data", "running.intelligence",
	// "suspended.missing-answer.auth.security", "suspended.missing-answer.auth.otp",
	// "suspended.missing-answer.auth.captcha", "suspended.missing-answer.auth.puzzle",
	// "inactive", "auth.creds.security.invalid", "auth.creds.otp.invalid",
	// "auth.creds.captcha.invalid", "auth.creds.puzzle.invalid",
	// "auth.creds.certificate.required", "guest.intervention.required",
	// "auth.creds.invalid", "auth.creds.locked.temporary", "auth.creds.locked.permanent",
	// "error.permanent", "error.temporary", "error.session", "error.network",
	// "error.service.unavailable", "error.unsupported", "unknown".
	AggregationStatus string `json:"aggregation_status"`
	// LastAggregatedAt is the last time data was acquired for this account.
	// Format: ISO 8601 date-time.
	LastAggregatedAt string `json:"last_aggregated_at"`
	// LastAggregatedSuccess is the last time data was successfully acquired for this account.
	// Format: ISO 8601 date-time.
	// This value is null if data has never been successfully acquired.
	LastAggregatedSuccess *string `json:"last_aggregated_success"`
	// CreatedAt is the time registered with Moneytree.
	// Format: ISO 8601 date-time.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the last updated time (updated by Moneytree or user changes, etc.).
	// Format: ISO 8601 date-time.
	UpdatedAt string `json:"updated_at"`
}

// PointAccounts represents the response from the point accounts endpoint.
type PointAccounts struct {
	// PointAccounts is a list of point accounts.
	PointAccounts []PointAccount `json:"point_accounts"`
}

// GetPointAccountsOption configures options for the GetPointAccounts API call.
type GetPointAccountsOption func(*getPointAccountsOptions)

type getPointAccountsOptions struct {
	paginationOptions
}

// WithPageForPointAccounts specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPageForPointAccounts(page int) GetPointAccountsOption {
	return func(opts *getPointAccountsOptions) {
		opts.Page = &page
	}
}

// WithPerPageForPointAccounts specifies the number of items per page.
// The default value is 500. Valid range is 1 to 500.
func WithPerPageForPointAccounts(perPage int) GetPointAccountsOption {
	return func(opts *getPointAccountsOptions) {
		opts.PerPage = &perPage
	}
}

// GetPointAccounts retrieves the list of all point accounts.
// This endpoint requires the points_read OAuth scope.
//
// This API returns all point accounts registered by the guest user.
// Point accounts include various point systems such as credit card points, shopping points, etc.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetPointAccounts(ctx, accessToken)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, account := range response.PointAccounts {
//		fmt.Printf("Account: ID=%d, Type=%s, Balance=%v\n", account.ID, account.AccountType, account.CurrentBalance)
//	}
//
// Example with pagination:
//
//	response, err := client.GetPointAccounts(ctx, accessToken,
//		moneytree.WithPageForPointAccounts(1),
//		moneytree.WithPerPageForPointAccounts(100),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-points-accounts
func (c *Client) GetPointAccounts(ctx context.Context, accessToken string, opts ...GetPointAccountsOption) (*PointAccounts, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	options := &getPointAccountsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	urlPath := "link/points/accounts.json"
	queryParams := url.Values{}
	applyPaginationParams(queryParams, &options.paginationOptions)
	if len(queryParams) > 0 {
		urlPath = fmt.Sprintf("%s?%s", urlPath, queryParams.Encode())
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res PointAccounts
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
