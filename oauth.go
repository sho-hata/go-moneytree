package moneytree

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GetAccessTokenRequest represents a request to get an access token.
type GetAccessTokenRequest struct {
	// Code is the authorization code received from the authorization endpoint.
	Code string
	// RedirectURI is the redirect URI that was used in the authorization request.
	RedirectURI string
}

// TokenResponse represents the response from the OAuth token endpoint.
type TokenResponse struct {
	// AccessToken is the access token that can be used to make API requests.
	AccessToken string `json:"access_token"`
	// TokenType is the type of token, typically "Bearer".
	TokenType string `json:"token_type"`
	// ExpiresIn is the number of seconds until the access token expires.
	ExpiresIn int `json:"expires_in"`
	// RefreshToken is the refresh token that can be used to obtain a new access token.
	RefreshToken string `json:"refresh_token"`
	// Scope is the scope of the access token.
	Scope string `json:"scope"`
}

// GetAccessToken exchanges an authorization code for an access token and refresh token.
//
// This method implements the OAuth 2.0 authorization code flow. After a user authorizes
// the application, the authorization server redirects to the redirect URI with an
// authorization code. This code is then exchanged for an access token using this method.
//
// Example:
//
//	client, _ := moneytree.NewClient(&moneytree.Config{
//		BaseURL:      "https://myaccount-staging.getmoneytree.com",
//		ClientID:     "your-client-id",
//		ClientSecret: "your-client-secret",
//	})
//
//	req := &moneytree.GetAccessTokenRequest{
//		Code:        "authorization-code-from-redirect",
//		RedirectURI: "https://your-app.com/callback",
//	}
//
//	token, err := client.GetAccessToken(context.Background(), req)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println("Access Token:", token.AccessToken)
func (c *Client) GetAccessToken(ctx context.Context, req *GetAccessTokenRequest) (*TokenResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if req.Code == "" {
		return nil, fmt.Errorf("code is required")
	}
	if req.RedirectURI == "" {
		return nil, fmt.Errorf("redirect_uri is required")
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", req.Code)
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)
	data.Set("redirect_uri", req.RedirectURI)

	tokenURL := fmt.Sprintf("%s/oauth/token", c.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	if err := checkResponseError(httpResp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &tokenResp, nil
}
