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
func (c *Client) GetAccountBalanceDetails(ctx context.Context, accessToken string, accountID string) (*AccountBalanceDetails, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	urlPath := fmt.Sprintf("link/accounts/%s/balances/details.json", url.PathEscape(accountID))

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res AccountBalanceDetails
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
