package moneytree

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// CorporateAccount represents a corporate account returned by the Moneytree LINK API.
// Corporate accounts include bank accounts, credit cards, etc. (excluding point accounts).
type CorporateAccount struct {
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
	// Possible values: "savings", "checking", "chochiku", "term_deposit", "term_deposit_builder",
	// "term_deposit_shikumi", "zaikei", "card_loan", "debit_card", "tax_payment_reserve_deposit",
	// "credit_card", "loan_installment", "asset_management", "home_loan", "stored_value",
	// "brokerage", "brokerage_cash", "pension_cash", "defined_contribution_pension",
	// "term_life", "whole_life".
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
	// Example: "積立貯金" (Savings deposit).
	InstitutionAccountName string `json:"institution_account_name"`
	// InstitutionAccountNumber is something like the "account number" of the account.
	// Specifications (number of digits, alphanumeric characters, etc.) vary by financial institution
	// and may change without notice. For credit cards, the number may be partially or fully masked.
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
	// AccountAttributes contains optional attributes for the account.
	// When the access token has both accounts_read and account_holder_read scopes,
	// this field may contain account holder information.
	// Note: account_holder_read scope and related fields are currently only available in Staging environment.
	AccountAttributes *CorporateAccountAttributes `json:"account_attributes,omitempty"`
}

// CorporateAccountAttributes represents optional attributes for a corporate account.
// This object may be empty depending on the account and OAuth scopes.
type CorporateAccountAttributes struct {
	// AccountHolderNameKatakanaRaw is the account holder name in katakana (raw format).
	// This field is only available when account_holder_read scope is included.
	// Note: Currently only available in Staging environment.
	AccountHolderNameKatakanaRaw *string `json:"account_holder_name_katakana_raw,omitempty"`
	// AccountHolderNameKatakanaZengin is the account holder name in katakana (Zengin format).
	// This field is only available when account_holder_read scope is included.
	// Note: Currently only available in Staging environment.
	AccountHolderNameKatakanaZengin *string `json:"account_holder_name_katakana_zengin,omitempty"`
}

// CorporateAccounts represents the response from the corporate accounts endpoint.
type CorporateAccounts struct {
	// Accounts is a list of corporate accounts.
	Accounts []CorporateAccount `json:"accounts"`
}

// GetCorporateAccountsOption configures options for the GetCorporateAccounts API call.
type GetCorporateAccountsOption func(*getCorporateAccountsOptions)

type getCorporateAccountsOptions struct {
	Page *int
}

// WithPageForCorporateAccounts specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPageForCorporateAccounts(page int) GetCorporateAccountsOption {
	return func(opts *getCorporateAccountsOptions) {
		opts.Page = &page
	}
}

// GetCorporateAccounts retrieves the list of all corporate accounts (excluding point accounts).
// This endpoint requires the accounts_read OAuth scope.
//
// This API returns all corporate accounts including bank accounts, credit cards,
// digital money, etc. registered by the guest user.
//
// When the access token has both accounts_read and account_holder_read scopes,
// the response includes account holder information in account_attributes.
// Note: account_holder_read scope and related fields are currently only available in Staging environment.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetCorporateAccounts(ctx, accessToken)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, account := range response.Accounts {
//		fmt.Printf("Account: %s, Subtype: %s, Balance: %v\n", account.AccountKey, account.AccountSubtype, account.CurrentBalance)
//	}
//
// Example with pagination:
//
//	response, err := client.GetCorporateAccounts(ctx, accessToken,
//		moneytree.WithPageForCorporateAccounts(1),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-corporate-accounts
func (c *Client) GetCorporateAccounts(ctx context.Context, accessToken string, opts ...GetCorporateAccountsOption) (*CorporateAccounts, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	options := &getCorporateAccountsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	urlPath := "link/corporate/accounts.json"
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

	var res CorporateAccounts
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// CorporateAccountBalance represents a balance record for a corporate account returned by the Moneytree LINK API.
type CorporateAccountBalance struct {
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
}

// CorporateAccountBalances represents the response from the corporate account balances endpoint.
type CorporateAccountBalances struct {
	// AccountBalances is a list of balance records for the account.
	AccountBalances []CorporateAccountBalance `json:"account_balances"`
}

// GetCorporateAccountBalancesOption configures options for the GetCorporateAccountBalances API call.
type GetCorporateAccountBalancesOption func(*getCorporateAccountBalancesOptions)

type getCorporateAccountBalancesOptions struct {
	paginationOptions
	SortKey *string
	SortBy  *string
	Since   *string
}

// WithPageForCorporateBalances specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPageForCorporateBalances(page int) GetCorporateAccountBalancesOption {
	return func(opts *getCorporateAccountBalancesOptions) {
		opts.Page = &page
	}
}

// WithPerPageForCorporateBalances specifies the number of items per page.
// The default value is 500. Valid range is 1 to 500.
func WithPerPageForCorporateBalances(perPage int) GetCorporateAccountBalancesOption {
	return func(opts *getCorporateAccountBalancesOptions) {
		opts.PerPage = &perPage
	}
}

// WithSortKeyForCorporateBalances specifies the sort key for balance records.
// If not provided, the database's id key is used by default.
// Using sort_key may affect response time, so it is recommended to use it only when necessary.
// If "date" is specified as the sort key, the database sorts by the balance date
// (which is the actual balance date, not the date Moneytree obtained it).
// The default value is "id".
func WithSortKeyForCorporateBalances(sortKey string) GetCorporateAccountBalancesOption {
	return func(opts *getCorporateAccountBalancesOptions) {
		opts.SortKey = &sortKey
	}
}

// WithSortByForCorporateBalances specifies the sort order.
// Possible values: "asc" (ascending, default), "desc" (descending).
// The default value is "asc".
func WithSortByForCorporateBalances(sortBy string) GetCorporateAccountBalancesOption {
	return func(opts *getCorporateAccountBalancesOptions) {
		opts.SortBy = &sortBy
	}
}

// WithSinceForCorporateBalances specifies a date to retrieve only records updated after this time (updated_at).
// This parameter takes precedence over start_date and end_date parameters.
// This is useful for incremental updates to avoid fetching all balances every time.
// Date format: "2006-01-02" (YYYY-MM-DD).
func WithSinceForCorporateBalances(since string) GetCorporateAccountBalancesOption {
	return func(opts *getCorporateAccountBalancesOptions) {
		opts.Since = &since
	}
}

// GetCorporateAccountBalances retrieves the balance history for a specific corporate account.
// This endpoint requires the accounts_read OAuth scope.
//
// This API returns balance records for the specified account. The balance history
// can be used to track changes in account balance over time.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetCorporateAccountBalances(ctx, accessToken, "account_key_123")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, balance := range response.AccountBalances {
//		fmt.Printf("Date: %s, Balance: %v, BalanceInBase: %v\n", balance.Date, balance.Balance, balance.BalanceInBase)
//	}
//
// Example with pagination and sorting:
//
//	response, err := client.GetCorporateAccountBalances(ctx, accessToken, "account_key_123",
//		moneytree.WithPageForCorporateBalances(1),
//		moneytree.WithPerPageForCorporateBalances(100),
//		moneytree.WithSortKeyForCorporateBalances("date"),
//		moneytree.WithSortByForCorporateBalances("desc"),
//	)
//
// Example with since parameter:
//
//	response, err := client.GetCorporateAccountBalances(ctx, accessToken, "account_key_123",
//		moneytree.WithSinceForCorporateBalances("2023-01-01"),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-corporate-account-balances
func (c *Client) GetCorporateAccountBalances(ctx context.Context, accessToken string, accountID string, opts ...GetCorporateAccountBalancesOption) (*CorporateAccountBalances, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	options := &getCorporateAccountBalancesOptions{}
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

	urlPath := fmt.Sprintf("link/corporate/accounts/%s/balances.json", url.PathEscape(accountID))
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

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res CorporateAccountBalances
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// CorporateAccountTransactionAttributes represents optional attributes for a corporate account transaction.
// This object may be empty depending on the transaction.
// The properties returned depend on the account's subtype.
type CorporateAccountTransactionAttributes struct {
	// FXBaseCurrency is the currency used for foreign currency credit card transactions.
	FXBaseCurrency *string `json:"fx_base_currency,omitempty"`
	// FXBaseAmount is the amount in local currency for foreign currency credit card transactions.
	FXBaseAmount *float64 `json:"fx_base_amount,omitempty"`
	// AuthorizationCode is the authorization code for debit card transactions.
	// Usually a 6-digit number.
	AuthorizationCode *string `json:"authorization_code,omitempty"`
	// Balance is the balance after this transaction for bank accounts.
	Balance *float64 `json:"balance,omitempty"`
	// IsRevolving indicates if this statement is a revolving payment.
	IsRevolving *bool `json:"is_revolving,omitempty"`
	// IsBonus indicates if this statement is a bonus payment.
	IsBonus *bool `json:"is_bonus,omitempty"`
	// IsCashAdvance indicates if this statement was a cash advance
	// (e.g., withdrawing cash using a credit card).
	IsCashAdvance *bool `json:"is_cash_advance,omitempty"`
	// InstallmentCount is the number of installment payments for this statement.
	// May be "1" if installment payments were not used.
	InstallmentCount *int `json:"installment_count,omitempty"`
	// ExpenseType is the type of transaction.
	// Deprecated: This field is deprecated.
	// Possible values: 0 = Unknown (assumed private use), 1 = Private use, 2 = Business.
	ExpenseType *int `json:"expense_type,omitempty"`
	// PredictedExpenseType is the predicted type of transaction.
	// Deprecated: This field is deprecated.
	// Possible values: 0 = Unknown (assumed private use), 1 = Private use, 2 = Business.
	PredictedExpenseType *int `json:"predicted_expense_type,omitempty"`
	// DataSource indicates the data source.
	// Deprecated: This field is deprecated.
	DataSource *string `json:"data_source,omitempty"`
	// TransactionType is the transaction type.
	// Deprecated: This field always returns null.
	TransactionType *string `json:"transaction_type,omitempty"`
}

// CorporateAccountTransaction represents a transaction record for a corporate account returned by the Moneytree LINK API.
type CorporateAccountTransaction struct {
	// ID is the transaction ID (unique across the entire system).
	// For example, if the same financial institution account is registered twice
	// with the same authentication information, different IDs will be assigned to each entity.
	ID int64 `json:"id"`
	// Amount is the transaction amount.
	Amount float64 `json:"amount"`
	// Date is the transaction date.
	// Format: ISO 8601 date-time.
	Date string `json:"date"`
	// DescriptionGuest is the content of the transaction entered by the customer.
	DescriptionGuest *string `json:"description_guest"`
	// DescriptionPretty is the content of the transaction corrected by Moneytree.
	DescriptionPretty *string `json:"description_pretty"`
	// DescriptionRaw is the unedited transaction content (raw data).
	// Regarding the details (summary field), there are digit restrictions depending on the bank API specifications.
	// For details, please check the publicly available API specifications of each bank.
	DescriptionRaw *string `json:"description_raw"`
	// AccountID is the account ID.
	AccountID int64 `json:"account_id"`
	// CategoryID is the category ID of the transaction detail.
	CategoryID int64 `json:"category_id"`
	// Attributes contains optional attributes for the transaction.
	// This object may be empty depending on the transaction.
	// The properties returned depend on the account's subtype.
	Attributes CorporateAccountTransactionAttributes `json:"attributes"`
	// CategoryEntityKey is the entity key of the specified category in the transaction details.
	// If it is a user-defined category, this value is null. Otherwise, it has a value.
	CategoryEntityKey *string `json:"category_entity_key"`
	// CreatedAt is the time registered with Moneytree.
	// Format: ISO 8601 date-time.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the last updated time (updated by Moneytree or user changes, etc.).
	// Format: ISO 8601 date-time.
	UpdatedAt string `json:"updated_at"`
}

// CorporateAccountTransactions represents the response from the corporate account transactions endpoint.
type CorporateAccountTransactions struct {
	// Transactions is a list of transaction records for the account.
	Transactions []CorporateAccountTransaction `json:"transactions"`
}

// GetCorporateAccountTransactionsOption configures options for the GetCorporateAccountTransactions API call.
type GetCorporateAccountTransactionsOption func(*getCorporateTransactionsOptions)

type getCorporateTransactionsOptions struct {
	paginationOptions
	SortKey *string
	SortBy  *string
	Since   *string
}

// WithPageForCorporateTransactions specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPageForCorporateTransactions(page int) GetCorporateAccountTransactionsOption {
	return func(opts *getCorporateTransactionsOptions) {
		opts.Page = &page
	}
}

// WithPerPageForCorporateTransactions specifies the number of items per page.
// The default value is 500. Valid range is 1 to 500.
func WithPerPageForCorporateTransactions(perPage int) GetCorporateAccountTransactionsOption {
	return func(opts *getCorporateTransactionsOptions) {
		opts.PerPage = &perPage
	}
}

// WithSortKeyForCorporateTransactions specifies the sort key for transaction details.
// If not provided, the database's id key is used by default.
// Using sort_key may affect response time, so it is recommended to use it only when necessary.
// If "date" is specified as the sort key, the database sorts by the transaction date
// (which is the actual transaction date, not the date Moneytree obtained it).
// The default value is "id".
func WithSortKeyForCorporateTransactions(sortKey string) GetCorporateAccountTransactionsOption {
	return func(opts *getCorporateTransactionsOptions) {
		opts.SortKey = &sortKey
	}
}

// WithSortByForCorporateTransactions specifies the sort order.
// Possible values: "asc" (ascending, default), "desc" (descending).
// The default value is "asc".
func WithSortByForCorporateTransactions(sortBy string) GetCorporateAccountTransactionsOption {
	return func(opts *getCorporateTransactionsOptions) {
		opts.SortBy = &sortBy
	}
}

// WithSinceForCorporateTransactions specifies a date to retrieve only records updated after this time (updated_at).
// This is useful for incremental updates to avoid fetching all transactions every time.
// Date format: "2006-01-02" (YYYY-MM-DD).
func WithSinceForCorporateTransactions(since string) GetCorporateAccountTransactionsOption {
	return func(opts *getCorporateTransactionsOptions) {
		opts.Since = &since
	}
}

// GetCorporateAccountTransactions retrieves the transaction records for a specific corporate account.
// This endpoint requires the transactions_read OAuth scope.
//
// This API returns transaction records for the specified account.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetCorporateAccountTransactions(ctx, accessToken, "account_key_123")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, transaction := range response.Transactions {
//		fmt.Printf("Date: %s, Amount: %v, Description: %s\n", transaction.Date, transaction.Amount, *transaction.DescriptionPretty)
//	}
//
// Example with pagination and sorting:
//
//	response, err := client.GetCorporateAccountTransactions(ctx, accessToken, "account_key_123",
//		moneytree.WithPageForCorporateTransactions(1),
//		moneytree.WithPerPageForCorporateTransactions(100),
//		moneytree.WithSortKeyForCorporateTransactions("date"),
//		moneytree.WithSortByForCorporateTransactions("desc"),
//	)
//
// Example with since parameter:
//
//	response, err := client.GetCorporateAccountTransactions(ctx, accessToken, "account_key_123",
//		moneytree.WithSinceForCorporateTransactions("2023-01-01"),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-corporate-accounts-transactions
func (c *Client) GetCorporateAccountTransactions(ctx context.Context, accessToken string, accountID string, opts ...GetCorporateAccountTransactionsOption) (*CorporateAccountTransactions, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	options := &getCorporateTransactionsOptions{}
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

	urlPath := fmt.Sprintf("link/corporate/accounts/%s/transactions.json", url.PathEscape(accountID))
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

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res CorporateAccountTransactions
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateCorporateAccountTransactionRequest represents a request to update a corporate account transaction.
type UpdateCorporateAccountTransactionRequest struct {
	// DescriptionGuest is a description/memo for transaction details, up to 255 characters.
	// If null is set, previous data will be deleted.
	// Do not set this parameter if you are not changing the value.
	DescriptionGuest *string `json:"description_guest,omitempty"`
	// CategoryID is the category of the transaction details.
	// If the corresponding ID (common category or this guest user's category) does not exist, 400 will be returned.
	// Do not set this parameter if you are not changing the value.
	CategoryID *int64 `json:"category_id,omitempty"`
}

// UpdateCorporateAccountTransaction updates a corporate account transaction.
// This endpoint requires the transactions_write OAuth scope.
//
// This API allows guest users to add memos (transaction content) to transaction details
// and edit category data automatically registered by Moneytree.
//
// Example:
//
//	descriptionGuest := "新しいメモ"
//	categoryID := int64(123)
//	request := &moneytree.UpdateCorporateAccountTransactionRequest{
//		DescriptionGuest: &descriptionGuest,
//		CategoryID:       &categoryID,
//	}
//	transaction, err := client.UpdateCorporateAccountTransaction(ctx, accessToken, "account_key_123", 1337, request)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Updated transaction: ID=%d, Description=%s\n", transaction.ID, *transaction.DescriptionGuest)
//
// Example with only description:
//
//	descriptionGuest := "取引メモ"
//	request := &moneytree.UpdateCorporateAccountTransactionRequest{
//		DescriptionGuest: &descriptionGuest,
//	}
//	transaction, err := client.UpdateCorporateAccountTransaction(ctx, accessToken, "account_key_123", 1337, request)
//
// Example to delete description:
//
//	request := &moneytree.UpdateCorporateAccountTransactionRequest{
//		DescriptionGuest: nil,
//	}
//	transaction, err := client.UpdateCorporateAccountTransaction(ctx, accessToken, "account_key_123", 1337, request)
//
// Reference: https://docs.link.getmoneytree.com/reference/put-link-corporate-account-transaction
func (c *Client) UpdateCorporateAccountTransaction(ctx context.Context, accessToken string, accountID string, transactionID int64, req *UpdateCorporateAccountTransactionRequest) (*CorporateAccountTransaction, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	if req.DescriptionGuest != nil && len(*req.DescriptionGuest) > 255 {
		return nil, fmt.Errorf("description_guest must be 255 characters or less, got %d characters", len(*req.DescriptionGuest))
	}

	urlPath := fmt.Sprintf("link/corporate/accounts/%s/transactions/%d.json", url.PathEscape(accountID), transactionID)

	httpReq, err := c.NewRequest(ctx, http.MethodPut, urlPath, req, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res CorporateAccountTransaction
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
