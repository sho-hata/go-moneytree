package moneytree

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Category represents a category returned by the Moneytree LINK API.
type Category struct {
	// ID is the unique identifier for the category.
	// Deprecated: Use EntityKey instead, as ID varies by environment (staging/production).
	// ID may differ between staging and production environments.
	ID int64 `json:"id"`
	// EntityKey is the unique identifier for the category.
	// Use this instead of ID, as EntityKey is consistent across environments (staging/production).
	// If IsSystem is true, EntityKey is a unique ID.
	// If IsSystem is false (individual categories), the combination of ID and EntityKey differs for each guest user.
	// Also, the combination of ID and EntityKey may differ between staging and production environments.
	// EntityKey is only assigned to categories defined by Moneytree.
	// For categories created by users, EntityKey will be null.
	EntityKey *string `json:"entity_key,omitempty"`
	// CategoryType is the type of category.
	// Possible values: "expense" (支出), "income" (収入).
	// For special categories like transfer, repayment, investment (and their subcategories), this will be null.
	CategoryType *string `json:"category_type,omitempty"`
	// Name is the name of the category.
	// For common categories (IsSystem == true), names are available in English and Japanese.
	Name string `json:"name"`
	// ParentID is the ID of the parent category.
	// It will be null if a parent category does not exist.
	ParentID *int64 `json:"parent_id,omitempty"`
	// IsSystem indicates whether this is a system category defined by Moneytree.
	// If true, this category is a common category that all guests can use.
	// If false, this category was added by this guest (it cannot be seen by other guests).
	IsSystem bool `json:"is_system"`
	// CreatedAt is the time when the category was registered in Moneytree.
	// Format: ISO 8601 date-time.
	CreatedAt string `json:"created_at"`
	// UpdatedAt is the last updated time.
	// This is updated by Moneytree or user changes.
	// Format: ISO 8601 date-time.
	UpdatedAt string `json:"updated_at"`
}

// Categories represents the response from the categories endpoint.
type Categories struct {
	// Categories is a list of categories available to the guest user at login.
	Categories []Category `json:"categories"`
}

// GetCategoriesOption configures options for the GetCategories API call.
type GetCategoriesOption func(*getCategoriesOptions)

type getCategoriesOptions struct {
	Page   *int
	Locale *string
}

// WithPageForCategories specifies the page number for pagination.
// Page numbers start from 1. The default value is 1.
// Valid range is 1 to 100000.
func WithPageForCategories(page int) GetCategoriesOption {
	return func(opts *getCategoriesOptions) {
		opts.Page = &page
	}
}

// WithLocale specifies the display language for category names.
// Possible values: "en" (English), "ja" (Japanese).
func WithLocale(locale string) GetCategoriesOption {
	return func(opts *getCategoriesOptions) {
		opts.Locale = &locale
	}
}

// GetCategories retrieves the list of categories available to the guest user at login.
// This endpoint requires the transactions_read OAuth scope.
//
// This API returns all categories that can be used by the guest user.
// Categories defined by Moneytree have an EntityKey and IsSystem is true.
// Categories created by users or through the app/web service have EntityKey as null and IsSystem is false.
//
// Note: The ID field may differ between staging and production environments.
// Use EntityKey instead of ID to identify categories consistently across environments.
//
// Example:
//
//	client := moneytree.NewClient("jp-api-staging")
//	response, err := client.GetCategories(ctx, accessToken)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, category := range response.Categories {
//		if category.IsSystem {
//			fmt.Printf("System category: %s (EntityKey: %s)\n", category.Name, *category.EntityKey)
//		} else {
//			fmt.Printf("User category: %s\n", category.Name)
//		}
//	}
//
// Example with pagination and locale:
//
//	response, err := client.GetCategories(ctx, accessToken,
//		moneytree.WithPageForCategories(1),
//		moneytree.WithLocale("ja"),
//	)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-categories
func (c *Client) GetCategories(ctx context.Context, accessToken string, opts ...GetCategoriesOption) (*Categories, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	options := &getCategoriesOptions{}
	for _, opt := range opts {
		opt(options)
	}

	urlPath := "link/categories.json"
	queryParams := url.Values{}
	if options.Page != nil {
		queryParams.Set("page", fmt.Sprintf("%d", *options.Page))
	}
	if options.Locale != nil {
		if *options.Locale != "en" && *options.Locale != "ja" {
			return nil, fmt.Errorf("locale must be either 'en' or 'ja', got %s", *options.Locale)
		}
		queryParams.Set("locale", *options.Locale)
	}
	if len(queryParams) > 0 {
		urlPath = fmt.Sprintf("%s?%s", urlPath, queryParams.Encode())
	}

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res Categories
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// CreateCategoryRequest represents a request to create a new category.
type CreateCategoryRequest struct {
	// Name is the name of the category.
	Name string `json:"name"`
	// ParentID is the ID of the parent category.
	ParentID int64 `json:"parent_id"`
}

// CreateCategory creates a new category.
// This endpoint requires the transactions_write OAuth scope.
//
// This API creates a new category for the guest user.
// The created category will have IsSystem set to false, meaning it is a user-created category
// that cannot be seen by other guests.
//
// Example:
//
//	request := &moneytree.CreateCategoryRequest{
//		Name:     "新しいカテゴリー",
//		ParentID: 0,
//	}
//	category, err := client.CreateCategory(ctx, accessToken, request)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Created category: ID=%d, Name=%s\n", category.ID, category.Name)
//
// Reference: https://docs.link.getmoneytree.com/reference/post-link-categories
func (c *Client) CreateCategory(ctx context.Context, accessToken string, req *CreateCategoryRequest) (*Category, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	urlPath := "link/categories.json"

	httpReq, err := c.NewRequest(ctx, http.MethodPost, urlPath, req, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res Category
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetCategory retrieves a specific category by its ID.
// This endpoint requires the transactions_read OAuth scope.
//
// This API returns the category information for the specified category ID.
//
// Example:
//
//	category, err := client.GetCategory(ctx, accessToken, 1048)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Category: ID=%d, Name=%s, IsSystem=%v\n", category.ID, category.Name, category.IsSystem)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-link-category
func (c *Client) GetCategory(ctx context.Context, accessToken string, categoryID int64) (*Category, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	urlPath := fmt.Sprintf("link/categories/%d.json", categoryID)

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil, WithBearerToken(accessToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res Category
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
