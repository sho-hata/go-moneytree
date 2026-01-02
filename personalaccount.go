package moneytree

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

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
	Page    *int
	PerPage *int
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
	if options.Page != nil {
		queryParams.Set("page", fmt.Sprintf("%d", *options.Page))
	}
	if options.PerPage != nil {
		queryParams.Set("per_page", fmt.Sprintf("%d", *options.PerPage))
	}
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
