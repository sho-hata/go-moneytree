package moneytree

import (
	"context"
	"fmt"
	"net/http"
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

	httpReq, err := c.NewRequest(http.MethodGet, "link/profile.json", nil, WithBearerToken(accessToken))
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

	httpReq, err := c.NewRequest(http.MethodPost, "link/profile/revoke.json", nil, WithBearerToken(accessToken))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if _, err = c.Do(ctx, httpReq, nil); err != nil {
		return err
	}
	return nil
}
