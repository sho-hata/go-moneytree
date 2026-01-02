package moneytree

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Profile represents the user profile information returned by the Moneytree LINK API.
type Profile struct {
	// LocaleIdentifier is the language and region identifier desired by the customer.
	// The format is [language]_[region], where _[region] is optional.
	// Language is a 2-digit ISO639 code representing the language.
	// Examples: "ja_JP" (Japanese, residing in Japan), "en_US" (English, residing in America), "ja_US" (Japanese, residing in America).
	LocaleIdentifier *string `json:"locale_identifier,omitempty"`
	// Email is the customer's current email address.
	// Note: Since it may be changed by the customer, it should not be used as a unique identifier.
	// Also, if a guest user changes their email address, it may take several minutes for the new data to be reflected.
	Email *string `json:"email,omitempty"`
	// MoneytreeID is a value that is unique within the system and cannot identify the customer.
	// Used as a unique identifier within the system.
	MoneytreeID *string `json:"moneytree_id,omitempty"`
}

// GetProfile retrieves the user profile information.
// This endpoint requires the guest_read OAuth scope.
func (c *Client) GetProfile(ctx context.Context, accessToken string) (*Profile, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, "link/profile.json", nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var profile Profile
	if _, err = c.Do(ctx, httpReq, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// RevokeProfile revokes the guest account connection.
// This endpoint requires the guest_read OAuth scope.
func (c *Client) RevokeProfile(ctx context.Context, accessToken string) error {
	if accessToken == "" {
		return fmt.Errorf("access token is required")
	}

	httpReq, err := c.NewRequest(ctx, http.MethodPost, "link/profile/revoke.json", nil, WithBearerToken(accessToken))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if _, err = c.Do(ctx, httpReq, nil); err != nil {
		return err
	}
	return nil
}

// AccountGroup represents an account group status returned by the Moneytree LINK API.
// An account group is a collection of accounts that were registered together
// through a single financial service registration.
type AccountGroup struct {
	// AggregationState represents the current data acquisition status.
	// Possible values: "success", "running", "error".
	AggregationState string `json:"aggregation_state"`
	// AggregationStatus represents the current data acquisition status in more detail than AggregationState.
	// For a list of possible values and their meanings, refer to the aggregation_status list guide.
	AggregationStatus string `json:"aggregation_status"`
	// LastAggregatedAt is the last time data was acquired.
	LastAggregatedAt time.Time `json:"last_aggregated_at"`
	// LastAggregatedSuccess is the last time data was successfully acquired.
	// This value is null if data has never been successfully acquired.
	LastAggregatedSuccess *time.Time `json:"last_aggregated_success"`
	// ID is the account group ID.
	// Deprecated: Use AccountGroup instead.
	ID *int64 `json:"id,omitempty"`
	// AccountGroup is the unique ID for the financial service registration group.
	// This value corresponds to the account_group value in each account information API.
	AccountGroup int64 `json:"account_group"`
	// InstitutionEntityKey is the key that identifies the financial service.
	// The name that can be displayed to customers can be obtained via the Financial Institution List API.
	InstitutionEntityKey string `json:"institution_entity_key"`
}

// AccountGroups represents the response from the account groups status endpoint.
type AccountGroups struct {
	// AccountGroups is a list of account groups registered by the guest user.
	AccountGroups []AccountGroup `json:"account_groups"`
}

// GetAccountGroups retrieves the status of all account groups for the guest user.
// This endpoint requires the accounts_read OAuth scope.
//
// Account groups represent collections of accounts that were registered together
// through a single financial service registration. For example, a single bank registration
// may provide access to checking accounts, savings accounts, and card loans.
//
// This API can be used to check the processing status and completion of synchronization requests.
// If last_aggregated_at is null, it indicates that the financial service registration
// has not been completed (either still in progress or failed during initial registration).
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetAccountGroups(ctx, accessToken)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, ag := range response.AccountGroups {
//		fmt.Printf("Account Group: %d, Status: %s\n", ag.AccountGroup, ag.AggregationStatus)
//	}
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-profile-account-groups
func (c *Client) GetAccountGroups(ctx context.Context, accessToken string) (*AccountGroups, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, "link/profile/account_groups.json", nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res AccountGroups
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// RefreshProfile requests a refresh of financial institution data.
// This endpoint requires the request_refresh OAuth scope.
//
// This API requests Moneytree to update financial institution data. When the request is accepted,
// Moneytree starts a job to access registered financial services using stored authentication credentials
// or access/refresh tokens to retrieve updated information.
//
// Note: This API is limited to 4 requests per guest per day (resets at 00:00 JST).
// Even if 202 is returned, some financial services may have update restrictions.
// Refer to the Financial Institution List API for details on restricted financial services
// and their update interval conditions.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	err := client.RefreshProfile(ctx, accessToken)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Reference: https://docs.link.getmoneytree.com/reference/post-link-profile-refresh
func (c *Client) RefreshProfile(ctx context.Context, accessToken string) error {
	if accessToken == "" {
		return fmt.Errorf("access token is required")
	}

	httpReq, err := c.NewRequest(ctx, http.MethodPost, "link/profile/refresh.json", nil, WithBearerToken(accessToken))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if _, err = c.Do(ctx, httpReq, nil); err != nil {
		return err
	}
	return nil
}

// RefreshAccountGroup requests a refresh of data for a specific account group.
// This endpoint requires the request_refresh OAuth scope.
//
// This API requests Moneytree to update data for the specified account group. When the request is accepted,
// Moneytree starts a job to access the registered financial service using stored authentication credentials
// or access/refresh tokens to retrieve updated information.
//
// This API is useful for business use cases where you want to synchronize specific account groups
// at specific times, rather than synchronizing all account groups at once.
//
// Note: Even if 202 is returned, some financial services may have update restrictions.
// Refer to the Financial Institution List API for details on restricted financial services
// and their update interval conditions.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	err := client.RefreshAccountGroup(ctx, accessToken, 12345)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Reference: https://docs.link.getmoneytree.com/reference/post-link-account-group-refresh
func (c *Client) RefreshAccountGroup(ctx context.Context, accessToken string, accountGroup int64) error {
	if accessToken == "" {
		return fmt.Errorf("access token is required")
	}

	urlPath := fmt.Sprintf("link/account_groups/%d/refresh.json", accountGroup)
	httpReq, err := c.NewRequest(ctx, http.MethodPost, urlPath, nil, WithBearerToken(accessToken))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if _, err = c.Do(ctx, httpReq, nil); err != nil {
		return err
	}
	return nil
}
