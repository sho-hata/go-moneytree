package moneytree

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// InvestmentAccount represents an investment account returned by the Moneytree LINK API.
// Investment accounts include securities accounts, pension accounts, etc.
type InvestmentAccount struct {
	// ID is the unique ID of the account assigned by Moneytree.
	ID int64 `json:"id"`
	// AccountKey is the unique identifier for the account.
	// Use this instead of ID, as ID varies by environment (staging/production).
	AccountKey string `json:"account_key"`
	// AccountGroup is the unique ID for the financial service registration group.
	// This value corresponds to the account_group value in account group information API.
	// Account groups are displayed as one user on the financial institution's website.
	AccountGroup int64 `json:"account_group"`
	// AccountSubtype describes the specific type of account.
	// Possible values: "brokerage", "brokerage_cash", "pension_cash", "defined_contribution_pension",
	// "term_life", "whole_life", etc.
	AccountSubtype string `json:"account_subtype"`
	// AccountType describes the type of account.
	// Deprecated: Use AccountSubtype instead, as it provides more detailed information.
	AccountType string `json:"account_type"`
	// Currency is the currency code of the account based on ISO4217 (e.g., "JPY", "USD").
	Currency string `json:"currency"`
	// InstitutionEntityKey is the key that identifies the financial service.
	// The name that can be displayed to customers can be obtained via the Financial Institution List API.
	// Use this instead of InstitutionID, as the institution_entity_key may change due to financial institution mergers.
	InstitutionEntityKey string `json:"institution_entity_key"`
	// InstitutionID is the ID of the financial institution.
	// Deprecated: Use InstitutionEntityKey instead.
	InstitutionID int64 `json:"institution_id"`
	// InstitutionAccountName is the account name of the financial institution.
	InstitutionAccountName string `json:"institution_account_name"`
	// InstitutionAccountNumber is something like the "account number" of the account.
	// Specifications (number of digits, alphanumeric characters, etc.) vary by financial institution
	// and may change without notice.
	InstitutionAccountNumber *string `json:"institution_account_number"`
	// InstitutionAccountType is the account type of the financial institution.
	// Deprecated: This field always returns null, so please use AccountSubtype instead.
	InstitutionAccountType *string `json:"institution_account_type"`
	// Nickname is the name given to the account by the user within Moneytree.
	// If the user has not set a nickname individually, the same value as InstitutionAccountName will be returned.
	Nickname string `json:"nickname"`
	// BranchName is the branch name of the financial institution's account.
	// Examples: "本店" (Head Office), "原宿支店" (Harajuku Branch), "309".
	// The data returned varies by financial institution. This value is provided for display purposes only.
	BranchName *string `json:"branch_name"`
	// BranchCode is the branch number. Returns null if not applicable.
	BranchCode *string `json:"branch_code"`
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
	LastAggregatedSuccess *string `json:"last_aggregated_success"`
	// CurrentBalance is the account balance as of LastAggregatedSuccess.
	CurrentBalance *float64 `json:"current_balance"`
	// CurrentBalanceInBase is the account balance converted to JPY.
	// If the financial service provides the foreign currency equivalent, that equivalent will be stored and returned.
	// If not supported, it will be calculated using the exchange rate used by Moneytree.
	CurrentBalanceInBase *float64 `json:"current_balance_in_base"`
	// CurrentBalanceDataSource is a label indicating the source of the data.
	// Possible values: "guest", "institution".
	CurrentBalanceDataSource *string `json:"current_balance_data_source"`
	// CreatedAt is the time registered with Moneytree.
	// Format: ISO 8601 date-time.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the last updated time (updated by Moneytree or user changes, etc.).
	// Format: ISO 8601 date-time.
	UpdatedAt string `json:"updated_at"`
}

// InvestmentAccounts represents the response from the investment accounts endpoint.
type InvestmentAccounts struct {
	// Accounts is a list of investment accounts.
	Accounts []InvestmentAccount `json:"accounts"`
}

// GetInvestmentAccountsOption configures options for the GetInvestmentAccounts API call.
type GetInvestmentAccountsOption func(*getInvestmentAccountsOptions)

type getInvestmentAccountsOptions struct {
	Page *int
}

// WithPageForInvestmentAccounts specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPageForInvestmentAccounts(page int) GetInvestmentAccountsOption {
	return func(opts *getInvestmentAccountsOptions) {
		opts.Page = &page
	}
}

// GetInvestmentAccounts retrieves the list of all investment accounts.
// This endpoint requires the investment_accounts_read OAuth scope.
//
// This API returns all investment or investment-like accounts including securities accounts,
// pension accounts, etc. registered by the guest user.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetInvestmentAccounts(ctx, accessToken)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, account := range response.Accounts {
//		fmt.Printf("Account: %s, Subtype: %s, Balance: %v\n", account.AccountKey, account.AccountSubtype, account.CurrentBalance)
//	}
//
// Example with pagination:
//
//	response, err := client.GetInvestmentAccounts(ctx, accessToken,
//		moneytree.WithPageForInvestmentAccounts(1),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-investments-accounts
func (c *Client) GetInvestmentAccounts(ctx context.Context, accessToken string, opts ...GetInvestmentAccountsOption) (*InvestmentAccounts, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	options := &getInvestmentAccountsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	urlPath := "link/investments/accounts.json"
	queryParams := url.Values{}
	if options.Page != nil {
		queryParams.Set("page", fmt.Sprintf("%d", *options.Page))
	}
	if len(queryParams) > 0 {
		urlPath = fmt.Sprintf("%s?%s", urlPath, queryParams.Encode())
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res InvestmentAccounts
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// InvestmentPosition represents a position record for an investment account returned by the Moneytree LINK API.
// Unlike transaction details, position details represent what assets the customer currently holds at a point in time.
// Positions change over time as market values fluctuate, so this API returns the most recently updated position details
// that Moneytree has confirmed, rather than historical records.
type InvestmentPosition struct {
	// ID is the unique identifier for the position record.
	// A new identifier is returned from the API each time the position is updated.
	ID int64 `json:"id"`
	// Date is the date associated with the position.
	// Format: "2006-01-02" (YYYY-MM-DD).
	Date string `json:"date"`
	// AssetClass is the classification of investment assets.
	// Possible values: "stock", "investment_trust", "bond", "cash", "commodity", "alternative", "other", "unknown".
	// Values not present in this enumeration may be added without warning.
	AssetClass string `json:"asset_class"`
	// AssetSubclass is a further classification of investment assets.
	AssetSubclass *string `json:"asset_subclass"`
	// TickerCode is the alphanumeric asset code for stocks being invested in.
	// Applicable only to stock investment accounts.
	TickerCode *string `json:"ticker_code"`
	// Ticker is the alphanumeric asset code for stocks being invested in.
	// Deprecated: Use TickerCode instead.
	Ticker *string `json:"ticker"`
	// NameRaw is the unformatted name assigned to the position.
	NameRaw *string `json:"name_raw"`
	// NameClean is the formatted name assigned to the position.
	NameClean *string `json:"name_clean"`
	// Currency is the currency code (ISO 4217) for calculating the value of the position.
	Currency string `json:"currency"`
	// TaxType is the classification of taxes applied to the assets of the position.
	// Possible values: "ippan", "tokutei", "NISA", "dc pension", "stock option", "unknown".
	TaxType []string `json:"tax_type"`
	// TaxSubType is detailed information on TaxType.
	// Possible values: "ippan", "tsumitate", "junior", "growth_investment", "tsumitate_investment".
	// It will be null if detailed information is not provided by the financial institution.
	TaxSubType *string `json:"tax_sub_type"`
	// MarketValue is the total currency value of the position evaluated at the updated_at timestamp.
	MarketValue float64 `json:"market_value"`
	// Value is the total currency value of the position evaluated at the updated_at timestamp.
	// Deprecated: Use MarketValue instead.
	Value float64 `json:"value"`
	// AcquisitionValue is the total currency value of the position when it was created.
	AcquisitionValue *float64 `json:"acquisition_value"`
	// CostBasis is the total currency value of the position when it was created.
	// Deprecated: Use AcquisitionValue instead.
	CostBasis *float64 `json:"cost_basis"`
	// Profit is the unrealized gain/loss (plus or minus) of the position,
	// considering the difference between acquisition cost and market value.
	Profit *float64 `json:"profit"`
	// Quantity is the quantity of assets included in the position.
	Quantity *float64 `json:"quantity"`
	// CreatedAt is the time registered with Moneytree.
	// Format: ISO 8601 date-time.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the last updated time (updated by Moneytree or user changes, etc.).
	// Format: ISO 8601 date-time.
	UpdatedAt string `json:"updated_at"`
}

// InvestmentPositions represents the response from the investment positions endpoint.
type InvestmentPositions struct {
	// Positions is a list of position records for the account.
	Positions []InvestmentPosition `json:"positions"`
}

// GetInvestmentPositionsOption configures options for the GetInvestmentPositions API call.
type GetInvestmentPositionsOption func(*getInvestmentPositionsOptions)

type getInvestmentPositionsOptions struct {
	Page *int
}

// WithPageForInvestmentPositions specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPageForInvestmentPositions(page int) GetInvestmentPositionsOption {
	return func(opts *getInvestmentPositionsOptions) {
		opts.Page = &page
	}
}

// GetInvestmentPositions retrieves the position records for a specific investment account.
// This endpoint requires the investment_transactions_read OAuth scope.
//
// Unlike transaction details, position details represent what assets the customer currently holds at a point in time.
// Positions change over time as market values fluctuate, so this API returns the most recently updated position details
// that Moneytree has confirmed, rather than historical records.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetInvestmentPositions(ctx, accessToken, "account_key_123")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, position := range response.Positions {
//		fmt.Printf("Asset: %s, Market Value: %v, Quantity: %v\n",
//			*position.NameClean, position.MarketValue, *position.Quantity)
//	}
//
// Example with pagination:
//
//	response, err := client.GetInvestmentPositions(ctx, accessToken, "account_key_123",
//		moneytree.WithPageForInvestmentPositions(1),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-investments-accounts-positions
func (c *Client) GetInvestmentPositions(ctx context.Context, accessToken string, accountID string, opts ...GetInvestmentPositionsOption) (*InvestmentPositions, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	options := &getInvestmentPositionsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	urlPath := fmt.Sprintf("link/investments/accounts/%s/positions.json", url.PathEscape(accountID))
	queryParams := url.Values{}
	if options.Page != nil {
		queryParams.Set("page", fmt.Sprintf("%d", *options.Page))
	}
	if len(queryParams) > 0 {
		urlPath = fmt.Sprintf("%s?%s", urlPath, queryParams.Encode())
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res InvestmentPositions
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

