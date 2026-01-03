package moneytree

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// AccountBalanceDetail represents a balance detail record for an account returned by the Moneytree LINK API.
// This API returns detailed balance information for all accounts that have available balances.
type AccountBalanceDetail struct {
	// ID is the balance record ID.
	ID int64 `json:"id"`
	// AccountID is the account ID.
	AccountID int64 `json:"account_id"`
	// Date is the date when the balance was confirmed on the financial institution's website.
	// Format: "2006-01-02" (YYYY-MM-DD).
	Date string `json:"date"`
	// Balance is the account balance.
	Balance float64 `json:"balance"`
	// BalanceInBase is the account balance converted to JPY.
	// If the financial service provides the converted amount for foreign currency,
	// that amount is stored and returned in this field. If not supported,
	// it is calculated using the exchange rate used by Moneytree.
	BalanceInBase float64 `json:"balance_in_base"`
	// BalanceType indicates the type of balance.
	// Possible values:
	//   0 = Total credit card amount. For non-debt accounts, refers to ordinary balance.
	//   1 = Undetermined amount for credit cards, etc.
	//   2 = Confirmed amount for credit cards, etc.
	//   3 = Long-term debt for credit cards, etc. (revolving, bonus payments, installment payments, etc.) amount
	BalanceType *int `json:"balance_type"`
}

// AccountBalanceDetails represents the response from the account balance details endpoint.
type AccountBalanceDetails struct {
	// AccountBalances is a list of balance detail records for the account.
	AccountBalances []AccountBalanceDetail `json:"account_balances"`
}

// GetAccountBalanceDetails retrieves the detailed balance information for a specific account.
// This endpoint requires the accounts_read OAuth scope.
//
// This API returns detailed balance information for all accounts that have available balances.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetAccountBalanceDetails(ctx, accessToken, "account_key_123")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, detail := range response.AccountBalances {
//		fmt.Printf("Date: %s, Balance: %v, BalanceInBase: %v\n", detail.Date, detail.Balance, detail.BalanceInBase)
//	}
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-account-balance-details-1
func (c *Client) GetAccountBalanceDetails(ctx context.Context, accountID string) (*AccountBalanceDetails, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	urlPath := fmt.Sprintf("link/accounts/%s/balances/details.json", url.PathEscape(accountID))

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res AccountBalanceDetails
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// AccountDueBalance represents a due balance record for an account returned by the Moneytree LINK API.
// This API returns payment due dates and amounts for accounts (mainly personal and corporate credit cards).
type AccountDueBalance struct {
	// ID is the balance record ID.
	ID int64 `json:"id"`
	// AccountID is the account ID.
	AccountID int64 `json:"account_id"`
	// Date is the date when the balance was confirmed on the financial institution's website.
	// Format: "2006-01-02" (YYYY-MM-DD).
	Date string `json:"date"`
	// DueAmount is the payment amount.
	DueAmount *float64 `json:"due_amount"`
	// DueDate is the payment date.
	DueDate string `json:"due_date"`
}

// AccountDueBalances represents the response from the account due balances endpoint.
type AccountDueBalances struct {
	// DueBalances is the due balance record for the account.
	DueBalances AccountDueBalance `json:"due_balances"`
}

// GetAccountDueBalancesOption configures options for the GetAccountDueBalances API call.
type GetAccountDueBalancesOption func(*getAccountDueBalancesOptions)

type getAccountDueBalancesOptions struct {
	Page      *int
	Since     *string
	StartDate *string
	EndDate   *string
}

// WithPageForDueBalances specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPageForDueBalances(page int) GetAccountDueBalancesOption {
	return func(opts *getAccountDueBalancesOptions) {
		opts.Page = &page
	}
}

// WithSinceForDueBalances specifies a date to retrieve only records updated after this time (updated_at).
// This parameter takes precedence over start_date and end_date parameters.
// This is useful for incremental updates to avoid fetching all due balances every time.
// Date format: "2006-01-02" (YYYY-MM-DD).
func WithSinceForDueBalances(since string) GetAccountDueBalancesOption {
	return func(opts *getAccountDueBalancesOptions) {
		opts.Since = &since
	}
}

// WithStartDateForDueBalances specifies the start date for transaction details.
// Deprecated: Use WithSinceForDueBalances instead.
// If specified, end_date is also required.
// Date format: "2006-01-02" (YYYY-MM-DD).
func WithStartDateForDueBalances(startDate string) GetAccountDueBalancesOption {
	return func(opts *getAccountDueBalancesOptions) {
		opts.StartDate = &startDate
	}
}

// WithEndDateForDueBalances specifies the end date for transaction details.
// Deprecated: Use WithSinceForDueBalances instead.
// If specified, start_date is also required.
// Date format: "2006-01-02" (YYYY-MM-DD).
func WithEndDateForDueBalances(endDate string) GetAccountDueBalancesOption {
	return func(opts *getAccountDueBalancesOptions) {
		opts.EndDate = &endDate
	}
}

// GetAccountDueBalances retrieves the payment due dates and amounts for a specific account.
// This endpoint requires the accounts_read OAuth scope.
//
// This API returns payment due dates and amounts for accounts (mainly personal and corporate credit cards).
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetAccountDueBalances(ctx, accessToken, "account_key_123")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Due Date: %s, Due Amount: %v\n", response.DueBalances.DueDate, *response.DueBalances.DueAmount)
//
// Example with since parameter:
//
//	response, err := client.GetAccountDueBalances(ctx, accessToken, "account_key_123",
//		moneytree.WithSinceForDueBalances("2023-01-01"),
//	)
//
// Example with pagination:
//
//	response, err := client.GetAccountDueBalances(ctx, accessToken, "account_key_123",
//		moneytree.WithPageForDueBalances(1),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-account-due-balances-1
func (c *Client) GetAccountDueBalances(ctx context.Context, accountID string, opts ...GetAccountDueBalancesOption) (*AccountDueBalances, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	options := &getAccountDueBalancesOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if options.Since != nil {
		if err := validateDateFormat(*options.Since); err != nil {
			return nil, err
		}
	}

	if options.StartDate != nil {
		if err := validateDateFormat(*options.StartDate); err != nil {
			return nil, err
		}
		if options.EndDate == nil {
			return nil, fmt.Errorf("end_date is required when start_date is specified")
		}
	}

	if options.EndDate != nil {
		if err := validateDateFormat(*options.EndDate); err != nil {
			return nil, err
		}
		if options.StartDate == nil {
			return nil, fmt.Errorf("start_date is required when end_date is specified")
		}
	}

	urlPath := fmt.Sprintf("link/accounts/%s/due_balances.json", url.PathEscape(accountID))
	queryParams := url.Values{}
	if options.Page != nil {
		queryParams.Set("page", fmt.Sprintf("%d", *options.Page))
	}
	if options.Since != nil {
		queryParams.Set("since", *options.Since)
	}
	if options.StartDate != nil {
		queryParams.Set("start_date", *options.StartDate)
	}
	if options.EndDate != nil {
		queryParams.Set("end_date", *options.EndDate)
	}
	if len(queryParams) > 0 {
		urlPath = fmt.Sprintf("%s?%s", urlPath, queryParams.Encode())
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res AccountDueBalances
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
