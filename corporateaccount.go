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
