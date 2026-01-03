package moneytree

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Institution represents a financial institution returned by the Moneytree LINK API.
type Institution struct {
	// Deprecated: Use EntityKey instead, as ID varies by environment (staging/production).
	// ID is the unique ID of the financial institution.
	ID *int64 `json:"id,omitempty"`
	// EntityKey is the unique ID of the financial institution.
	// Use this instead of ID, as ID varies by environment (staging/production).
	EntityKey string `json:"entity_key"`
	// InstitutionType describes the type of financial institution.
	// Possible values: "bank" (individual bank), "credit_card" (individual credit card),
	// "stored_value" (electronic money), "point" (point card), "corporate" (corporate bank account, corporate card),
	// "stock" (securities). More types may be added in the future.
	InstitutionType string `json:"institution_type"`
	// DisplayName is the localized display name of the financial institution.
	DisplayName *string `json:"display_name"`
	// DisplayNameReading is the phonetic reading (kana) of the financial institution's name if Japanese is specified.
	// Otherwise, it's the same as DisplayName.
	DisplayNameReading *string `json:"display_name_reading"`
	// Status indicates whether data can be acquired.
	// Possible values: "active" (data can be acquired), "inactive" (data cannot be acquired, refer to StatusReason for details).
	// Returns null if unknown.
	Status *string `json:"status"`
	// StatusReason provides the reason for an inactive status.
	// Possible values: "maintenance" (temporarily inactive due to maintenance),
	// "unavailable" (inactive for a long time), "unsupported" (not supported by Moneytree, but future support is planned/considered),
	// "wont_support" (not supported for other reasons), "legacy" (previously supported but no longer exists, e.g., merged banks),
	// "test" (for testing purposes, not for general customers).
	// Returns null if unknown.
	StatusReason *string `json:"status_reason"`
	// Deprecated: This field was previously used for handling electronic certificates but is no longer in use.
	// CertificateRequired indicates the certificate requirement level.
	// Possible values: 0, 1, 2.
	CertificateRequired *int `json:"certificate_required"`
	// LoginURL is the URL for logging in on the financial institution's website.
	LoginURL *string `json:"login_url"`
	// GuidanceURL is a URL providing guidance on how to register the financial institution
	// (can be provided by Moneytree or the institution).
	GuidanceURL *string `json:"guidance_url"`
	// BillingGroup is for financial services with update restrictions within a certain period.
	// A number (1 or more) is returned; otherwise, null.
	// Contact a sales representative for update intervals.
	// Possible values: null, "2", "3", "4".
	BillingGroup *string `json:"billing_group"`
	// Tags are tags associated with the financial institution.
	Tags []string `json:"tags"`
	// DefaultAuthorizationType describes how Moneytree acquires data.
	// Possible values: 0 (web scraping), 1 (API scraping).
	DefaultAuthorizationType int `json:"default_authorization_type"`
}

// Institutions represents the response from the institutions list endpoint.
type Institutions struct {
	// Institutions is a list of financial institutions.
	Institutions []Institution `json:"institutions"`
}

// GetInstitutionsOption configures options for the GetInstitutions API call.
type GetInstitutionsOption func(*getInstitutionsOptions)

type getInstitutionsOptions struct {
	Since *string
}

// WithSince specifies a date to retrieve only institutions updated after this time.
// This is useful for incremental updates to avoid fetching all institutions every time.
// Date format: "2006-01-02" (YYYY-MM-DD).
func WithSince(since string) GetInstitutionsOption {
	return func(opts *getInstitutionsOptions) {
		opts.Since = &since
	}
}

// GetInstitutions retrieves the list of financial institutions.
// This endpoint does not require any OAuth scope.
//
// This API returns all financial institutions on Moneytree, including not only available ones
// but also institutions in various states (e.g., under maintenance, merged banks).
// Always check the status and status_reason attributes when consuming this API.
//
// Note: This is a system-level API, not a per-customer API. The access token used
// for this API is different from customer-specific APIs. Please refer to the authentication
// documentation for details.
//
// To reduce response size, use WithSince option to fetch only incremental updates
// rather than fetching all institutions every time.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetInstitutions(ctx, systemAccessToken)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, inst := range response.Institutions {
//		if inst.Status == "available" {
//			fmt.Printf("Available: %s (%s)\n", inst.Name, inst.EntityKey)
//		}
//	}
//
// Example with since parameter:
//
//	response, err := client.GetInstitutions(ctx, systemAccessToken, moneytree.WithSince("2023-01-01"))
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Reference: https://docs.link.getmoneytree.com/reference/get-institutions
func (c *Client) GetInstitutions(ctx context.Context, opts ...GetInstitutionsOption) (*Institutions, error) {
	options := &getInstitutionsOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if options.Since != nil {
		if err := validateDateFormat(*options.Since); err != nil {
			return nil, err
		}
	}

	urlPath := "link/institutions.json"
	if options.Since != nil {
		urlPath = fmt.Sprintf("%s?since=%s", urlPath, url.QueryEscape(*options.Since))
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res Institutions
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
