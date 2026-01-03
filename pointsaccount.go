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
func (c *Client) GetPointAccounts(ctx context.Context, opts ...GetPointAccountsOption) (*PointAccounts, error) {
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

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res PointAccounts
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// PointAccountTransaction represents a transaction record for a point account returned by the Moneytree LINK API.
// The specification is the same as personal account transactions.
// This type is an alias for PersonalAccountTransaction for clarity and consistency.
type PointAccountTransaction = PersonalAccountTransaction

// PointAccountTransactions represents the response from the point account transactions endpoint.
type PointAccountTransactions struct {
	// Transactions is a list of transaction records for the account.
	Transactions []PointAccountTransaction `json:"transactions"`
}

// GetPointAccountTransactionsOption configures options for the GetPointAccountTransactions API call.
type GetPointAccountTransactionsOption func(*getTransactionsOptions)

// WithPageForPointAccountTransactions specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPageForPointAccountTransactions(page int) GetPointAccountTransactionsOption {
	return func(opts *getTransactionsOptions) {
		opts.Page = &page
	}
}

// WithPerPageForPointAccountTransactions specifies the number of items per page.
// The default value is 500. Valid range is 1 to 500.
func WithPerPageForPointAccountTransactions(perPage int) GetPointAccountTransactionsOption {
	return func(opts *getTransactionsOptions) {
		opts.PerPage = &perPage
	}
}

// WithSortKeyForPointAccountTransactions specifies the sort key for transaction details.
// If not provided, the database's id key is used by default.
// Using sort_key may affect response time, so it is recommended to use it only when necessary.
// If "date" is specified as the sort key, the database sorts by the transaction date
// (which is the actual transaction date, not the date Moneytree obtained it).
// The default value is "id".
func WithSortKeyForPointAccountTransactions(sortKey string) GetPointAccountTransactionsOption {
	return func(opts *getTransactionsOptions) {
		opts.SortKey = &sortKey
	}
}

// WithSortByForPointAccountTransactions specifies the sort order.
// Possible values: "asc" (ascending, default), "desc" (descending).
// The default value is "asc".
func WithSortByForPointAccountTransactions(sortBy string) GetPointAccountTransactionsOption {
	return func(opts *getTransactionsOptions) {
		opts.SortBy = &sortBy
	}
}

// WithSinceForPointAccountTransactions specifies a date to retrieve only records updated after this time (updated_at).
// This is useful for incremental updates to avoid fetching all transactions every time.
// Date format: "2006-01-02" (YYYY-MM-DD).
func WithSinceForPointAccountTransactions(since string) GetPointAccountTransactionsOption {
	return func(opts *getTransactionsOptions) {
		opts.Since = &since
	}
}

// GetPointAccountTransactions retrieves the transaction records for a specific point account.
// This endpoint requires the points_read OAuth scope.
//
// This API returns transaction records for point accounts.
// The specification is the same as personal account transactions.
// Only the API path and required scope (points_read) differ.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetPointAccountTransactions(ctx, accessToken, 1048)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, transaction := range response.Transactions {
//		fmt.Printf("Date: %s, Amount: %v, Description: %s\n", transaction.Date, transaction.Amount, *transaction.DescriptionPretty)
//	}
//
// Example with pagination and sorting:
//
//	response, err := client.GetPointAccountTransactions(ctx, accessToken, 1048,
//		moneytree.WithPageForPointAccountTransactions(1),
//		moneytree.WithPerPageForPointAccountTransactions(100),
//		moneytree.WithSortKeyForPointAccountTransactions("date"),
//		moneytree.WithSortByForPointAccountTransactions("desc"),
//	)
//
// Example with since parameter:
//
//	response, err := client.GetPointAccountTransactions(ctx, accessToken, 1048,
//		moneytree.WithSinceForPointAccountTransactions("2023-01-01"),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-points-accounts-transactions
func (c *Client) GetPointAccountTransactions(ctx context.Context, accountID int64, opts ...GetPointAccountTransactionsOption) (*PointAccountTransactions, error) {
	options := &getTransactionsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if options.Since != nil {
		if err := validateDateFormat(*options.Since); err != nil {
			return nil, err
		}
	}

	if options.SortBy != nil {
		if *options.SortBy != "asc" && *options.SortBy != "desc" {
			return nil, fmt.Errorf("sort_by must be 'asc' or 'desc', got: %s", *options.SortBy)
		}
	}

	urlPath := fmt.Sprintf("link/points/accounts/%d/transactions.json", accountID)
	queryParams := url.Values{}
	applyPaginationParams(queryParams, &options.paginationOptions)
	if options.SortKey != nil {
		queryParams.Set("sort_key", *options.SortKey)
	}
	if options.SortBy != nil {
		queryParams.Set("sort_by", *options.SortBy)
	}
	if options.Since != nil {
		queryParams.Set("since", *options.Since)
	}
	if len(queryParams) > 0 {
		urlPath = fmt.Sprintf("%s?%s", urlPath, queryParams.Encode())
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res PointAccountTransactions
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// PointExpiration represents a point expiration record returned by the Moneytree LINK API.
// This record contains information about points that are expiring for a point account.
type PointExpiration struct {
	// ID is the point expiration record ID.
	ID int64 `json:"id"`
	// AccountID is the point account ID.
	AccountID int64 `json:"account_id"`
	// ExpirationAmount is the point balance reaching expiration.
	ExpirationAmount float64 `json:"expiration_amount"`
	// ExpirationDate is the expiration date.
	// Format: "2006-01-02" (YYYY-MM-DD).
	ExpirationDate string `json:"expiration_date"`
	// Date is the date when points reaching expiration were confirmed on the financial institution's website.
	// Format: ISO 8601 date-time.
	Date string `json:"date"`
}

// PointExpirations represents the response from the point expirations endpoint.
type PointExpirations struct {
	// PointExpirations is a list of point expiration records for the account.
	PointExpirations []PointExpiration `json:"point_expirations"`
}

// GetPointExpirationsOption configures options for the GetPointExpirations API call.
type GetPointExpirationsOption func(*getPointExpirationsOptions)

type getPointExpirationsOptions struct {
	paginationOptions
	Since *string
}

// WithPageForPointExpirations specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPageForPointExpirations(page int) GetPointExpirationsOption {
	return func(opts *getPointExpirationsOptions) {
		opts.Page = &page
	}
}

// WithPerPageForPointExpirations specifies the number of items per page.
// The default value is 500. Valid range is 1 to 500.
func WithPerPageForPointExpirations(perPage int) GetPointExpirationsOption {
	return func(opts *getPointExpirationsOptions) {
		opts.PerPage = &perPage
	}
}

// WithSinceForPointExpirations specifies a date to retrieve only records updated after this time (updated_at).
// This is useful for incremental updates to avoid fetching all expirations every time.
// Date format: "2006-01-02" (YYYY-MM-DD).
func WithSinceForPointExpirations(since string) GetPointExpirationsOption {
	return func(opts *getPointExpirationsOptions) {
		opts.Since = &since
	}
}

// GetPointExpirations retrieves the point expiration details for a specific point account.
// This endpoint requires the points_read OAuth scope.
//
// This API returns point expiration records for the specified account.
// Each record contains information about points that are expiring, including the expiration date
// and the amount of points that will expire.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetPointExpirations(ctx, 1048)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, expiration := range response.PointExpirations {
//		fmt.Printf("Expiration Date: %s, Amount: %v\n", expiration.ExpirationDate, expiration.ExpirationAmount)
//	}
//
// Example with pagination:
//
//	response, err := client.GetPointExpirations(ctx, 1048,
//		moneytree.WithPageForPointExpirations(1),
//		moneytree.WithPerPageForPointExpirations(100),
//	)
//
// Example with since parameter:
//
//	response, err := client.GetPointExpirations(ctx, 1048,
//		moneytree.WithSinceForPointExpirations("2023-01-01"),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-points-accounts-expirations
func (c *Client) GetPointExpirations(ctx context.Context, accountID int64, opts ...GetPointExpirationsOption) (*PointExpirations, error) {
	options := &getPointExpirationsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if options.Since != nil {
		if err := validateDateFormat(*options.Since); err != nil {
			return nil, err
		}
	}

	urlPath := fmt.Sprintf("link/points/accounts/%d/expirations.json", accountID)
	queryParams := url.Values{}
	applyPaginationParams(queryParams, &options.paginationOptions)
	if options.Since != nil {
		queryParams.Set("since", *options.Since)
	}
	if len(queryParams) > 0 {
		urlPath = fmt.Sprintf("%s?%s", urlPath, queryParams.Encode())
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res PointExpirations
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
