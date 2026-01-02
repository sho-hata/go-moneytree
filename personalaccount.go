package moneytree

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// paginationOptions represents common pagination options used across multiple API endpoints.
type paginationOptions struct {
	Page    *int
	PerPage *int
}

// applyPaginationParams applies pagination parameters to the query parameters.
func applyPaginationParams(queryParams url.Values, opts *paginationOptions) {
	if opts.Page != nil {
		queryParams.Set("page", fmt.Sprintf("%d", *opts.Page))
	}
	if opts.PerPage != nil {
		queryParams.Set("per_page", fmt.Sprintf("%d", *opts.PerPage))
	}
}

// PersonalAccount represents an individual account returned by the Moneytree LINK API.
// Individual accounts include bank accounts, credit cards, digital money, etc.
type PersonalAccount struct {
	// ID is the unique ID of the account.
	// Deprecated: Use AccountKey instead, as ID varies by environment (staging/production).
	ID *int64 `json:"id,omitempty"`
	// AccountKey is the unique identifier for the account.
	// Use this instead of ID, as ID varies by environment (staging/production).
	AccountKey string `json:"account_key"`
	// AccountGroup is the unique ID for the financial service registration group.
	// This value corresponds to the account_group value in account group information API.
	AccountGroup int64 `json:"account_group"`
	// InstitutionEntityKey is the key that identifies the financial service.
	// The name that can be displayed to customers can be obtained via the Financial Institution List API.
	InstitutionEntityKey string `json:"institution_entity_key"`
	// AccountType describes the type of account.
	// Possible values: "bank" (bank account), "credit_card" (credit card),
	// "stored_value" (electronic money), "point" (point card), "stock" (securities).
	AccountType string `json:"account_type"`
	// Name is the display name of the account.
	Name *string `json:"name,omitempty"`
	// Balance is the current balance of the account.
	// This value is null if the balance cannot be retrieved.
	Balance *float64 `json:"balance,omitempty"`
	// Currency is the currency code of the account (e.g., "JPY", "USD").
	Currency *string `json:"currency,omitempty"`
	// LastAggregatedAt is the last time data was acquired for this account.
	LastAggregatedAt *time.Time `json:"last_aggregated_at,omitempty"`
}

// PersonalAccounts represents the response from the individual accounts endpoint.
type PersonalAccounts struct {
	// Accounts is a list of individual accounts.
	Accounts []PersonalAccount `json:"accounts"`
}

// GetPersonalAccountsOption configures options for the GetPersonalAccounts API call.
type GetPersonalAccountsOption func(*getPersonalAccountsOptions)

type getPersonalAccountsOptions struct {
	paginationOptions
}

// WithPage specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPage(page int) GetPersonalAccountsOption {
	return func(opts *getPersonalAccountsOptions) {
		opts.Page = &page
	}
}

// WithPerPage specifies the number of items per page.
// The default value is 500. Valid range is 1 to 500.
func WithPerPage(perPage int) GetPersonalAccountsOption {
	return func(opts *getPersonalAccountsOptions) {
		opts.PerPage = &perPage
	}
}

// GetPersonalAccounts retrieves the list of all individual accounts.
// This endpoint requires the accounts_read OAuth scope.
//
// This API returns all individual accounts including bank accounts, credit cards,
// digital money, etc. registered by the guest user.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetPersonalAccounts(ctx, accessToken)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, account := range response.Accounts {
//		fmt.Printf("Account: %s, Type: %s, Balance: %v\n", account.AccountKey, account.AccountType, account.Balance)
//	}
//
// Example with pagination:
//
//	response, err := client.GetPersonalAccounts(ctx, accessToken,
//		moneytree.WithPage(1),
//		moneytree.WithPerPage(100),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-accounts
func (c *Client) GetPersonalAccounts(ctx context.Context, accessToken string, opts ...GetPersonalAccountsOption) (*PersonalAccounts, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	options := &getPersonalAccountsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	urlPath := "link/accounts.json"
	queryParams := url.Values{}
	applyPaginationParams(queryParams, &options.paginationOptions)
	if len(queryParams) > 0 {
		urlPath = fmt.Sprintf("%s?%s", urlPath, queryParams.Encode())
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res PersonalAccounts
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// PersonalAccountBalance represents a balance record for a personal account returned by the Moneytree LINK API.
type PersonalAccountBalance struct {
	// ID is the balance record ID.
	ID int64 `json:"id"`
	// AccountID is the account ID.
	AccountID int64 `json:"account_id"`
	// Date is the date when the balance was confirmed on the financial institution's website.
	Date time.Time `json:"date"`
	// Balance is the account balance.
	Balance float64 `json:"balance"`
	// BalanceInBase is the account balance converted to JPY.
	// If the financial service provides the converted amount for foreign currency,
	// that amount is stored and returned in this field. If not supported,
	// it is calculated using the exchange rate used by Moneytree.
	BalanceInBase float64 `json:"balance_in_base"`
}

// PersonalAccountBalances represents the response from the personal account balances endpoint.
type PersonalAccountBalances struct {
	// AccountBalances is a list of balance records for the account.
	AccountBalances []PersonalAccountBalance `json:"account_balances"`
}

// GetPersonalAccountBalancesOption configures options for the GetPersonalAccountBalances API call.
type GetPersonalAccountBalancesOption func(*getPersonalAccountBalancesOptions)

type getPersonalAccountBalancesOptions struct {
	paginationOptions
	Since *time.Time
}

// WithPageForBalances specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPageForBalances(page int) GetPersonalAccountBalancesOption {
	return func(opts *getPersonalAccountBalancesOptions) {
		opts.Page = &page
	}
}

// WithPerPageForBalances specifies the number of items per page.
// This sets the number of results to return per page when paginating the result set.
func WithPerPageForBalances(perPage int) GetPersonalAccountBalancesOption {
	return func(opts *getPersonalAccountBalancesOptions) {
		opts.PerPage = &perPage
	}
}

// WithSinceForBalances specifies a date to retrieve only records updated after this time (updated_at).
// This parameter takes precedence over start_date and end_date parameters.
// This is useful for incremental updates to avoid fetching all balances every time.
func WithSinceForBalances(t time.Time) GetPersonalAccountBalancesOption {
	return func(opts *getPersonalAccountBalancesOptions) {
		opts.Since = &t
	}
}

// GetPersonalAccountBalances retrieves the balance history for a specific personal account.
// This endpoint requires the accounts_read OAuth scope.
//
// This API returns balance records for the specified account. The balance history
// can be used to track changes in account balance over time.
//
// Note: This API can also retrieve balances for investment accounts, but you need
// the investment_accounts_read scope to get the account ID from the investment accounts list API.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetPersonalAccountBalances(ctx, accessToken, "account_key_123")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, balance := range response.AccountBalances {
//		fmt.Printf("Date: %s, Balance: %v, BalanceInBase: %v\n", balance.Date.Format("2006-01-02"), balance.Balance, balance.BalanceInBase)
//	}
//
// Example with since parameter:
//
//	sinceTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
//	response, err := client.GetPersonalAccountBalances(ctx, accessToken, "account_key_123",
//		moneytree.WithSinceForBalances(sinceTime),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-account-balances
func (c *Client) GetPersonalAccountBalances(ctx context.Context, accessToken string, accountID string, opts ...GetPersonalAccountBalancesOption) (*PersonalAccountBalances, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	options := &getPersonalAccountBalancesOptions{}
	for _, opt := range opts {
		opt(options)
	}

	urlPath := fmt.Sprintf("link/accounts/%s/balances.json", url.PathEscape(accountID))
	queryParams := url.Values{}
	applyPaginationParams(queryParams, &options.paginationOptions)
	if options.Since != nil {
		queryParams.Set("since", options.Since.Format("2006-01-02"))
	}
	if len(queryParams) > 0 {
		urlPath = fmt.Sprintf("%s?%s", urlPath, queryParams.Encode())
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res PersonalAccountBalances
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
