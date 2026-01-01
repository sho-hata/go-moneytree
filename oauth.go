package moneytree

import (
	"context"
	"fmt"
	"net/http"
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

func (c *Client) RetrieveToken(ctx context.Context, req *RetrieveTokenRequest) (*OauthToken, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	body := retrieveTokenRequest{
		RetrieveTokenRequest: *req,
		ClientID:             c.config.ClientID,
		ClientSecret:         c.config.ClientSecret,
	}
	httpReq, err := c.NewRequest(ctx, http.MethodPost, "oauth/token", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var t OauthToken
	if _, err = c.Do(ctx, httpReq, &t); err != nil {
		return nil, err
	}
	return &t, nil
}
