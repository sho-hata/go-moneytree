package moneytree

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type RetrieveTokenRequest struct {
	GrantType    *string `json:"grant_type,omitempty"`
	RedirectURI  *string `json:"redirect_uri,omitempty"`
	Code         *string `json:"code,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	CodeVerifier *string `json:"code_verifier,omitempty"`
	Scope        *string `json:"scope,omitempty"`
}

type retrieveTokenRequest struct {
	RetrieveTokenRequest
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

type OauthToken struct {
	AccessToken    *string `json:"access_token,omitempty"`
	TokenType      *string `json:"token_type,omitempty"`
	CreatedAt      *int    `json:"created_at,omitempty"`
	ExpiresIn      *int    `json:"expires_in,omitempty"`
	RefreshToken   *string `json:"refresh_token,omitempty"`
	Scope          *string `json:"scope,omitempty"`
	ResourceServer *string `json:"resource_server,omitempty"`
}

// RevokeTokenRequest represents a request to revoke an access token or refresh token.
type RevokeTokenRequest struct {
	// Token is the access token or refresh token to revoke.
	Token string
}

// RetrieveToken retrieves an access token or refresh token.
//
// Reference: https://docs.link.getmoneytree.com/reference/post-oauth-token
func (c *Client) RetrieveToken(ctx context.Context, req *RetrieveTokenRequest) (*OauthToken, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	body := retrieveTokenRequest{
		RetrieveTokenRequest: *req,
		ClientID:             c.config.ClientID,
		ClientSecret:         c.config.ClientSecret,
	}
	httpReq, err := c.NewRequest(http.MethodPost, "oauth/token", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var t OauthToken
	if _, err = c.Do(ctx, httpReq, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// RevokeToken revokes an access token or refresh token.
// According to the API documentation, this endpoint returns 200 OK even if the token
// does not exist or has already been revoked.
//
// Reference: https://docs.link.getmoneytree.com/reference/post-oauth-revoke
func (c *Client) RevokeToken(ctx context.Context, req *RevokeTokenRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if req.Token == "" {
		return fmt.Errorf("token is required")
	}

	form := url.Values{}
	form.Set("token", req.Token)
	form.Set("client_id", c.config.ClientID)
	form.Set("client_secret", c.config.ClientSecret)

	body := strings.NewReader(form.Encode())
	httpReq, err := c.NewFormRequest("oauth/revoke", body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if _, err = c.Do(ctx, httpReq, nil); err != nil {
		return err
	}
	return nil
}
