package moneytree

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// oauthTokenPath is the path for the OAuth token endpoint.
	oauthTokenPath = "oauth/token"
	// oauthRevokePath is the path for the OAuth revoke endpoint.
	oauthRevokePath = "oauth/revoke"
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

// Valid checks if the token is valid (not expired).
// It returns true if the token has an access token and is not expired.
// The token is considered expired if CreatedAt + ExpiresIn is before the current time.
// A buffer time of 1 minute is used to account for clock skew and network delays.
func (t *OauthToken) Valid() bool {
	if t == nil {
		return false
	}
	if t.AccessToken == nil {
		return false
	}
	if t.CreatedAt == nil || t.ExpiresIn == nil {
		return false
	}
	// Calculate expiration time: CreatedAt (Unix timestamp) + ExpiresIn (seconds)
	expiresAt := time.Unix(int64(*t.CreatedAt), 0).Add(time.Duration(*t.ExpiresIn) * time.Second)
	// Use a 1-minute buffer to account for clock skew and network delays
	bufferTime := 1 * time.Minute
	return time.Now().Add(bufferTime).Before(expiresAt)
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
	httpReq, err := c.NewRequest(ctx, http.MethodPost, oauthTokenPath, body)
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
	httpReq, err := c.NewFormRequest(ctx, oauthRevokePath, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if _, err = c.Do(ctx, httpReq, nil); err != nil {
		return err
	}
	return nil
}

// SetToken sets the OAuth token for the client.
// This method allows you to set a token that was obtained externally.
//
// Example:
//
//	token, err := client.RetrieveToken(ctx, &moneytree.RetrieveTokenRequest{...})
//	if err != nil {
//		log.Fatal(err)
//	}
//	client.SetToken(token)
func (c *Client) SetToken(token *OauthToken) {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()
	c.token = token
	c.getTokenErr = nil
}

// sleepWithContext sleeps for the specified duration, but returns early if the context is canceled.
func sleepWithContext(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

// refreshToken refreshes the token if necessary.
// This method implements a goroutine-safe token refresh mechanism.
// It checks if the current token is valid, and if not, attempts to refresh it
// using the refresh_token grant type with RetrieveToken.
// If another goroutine is already refreshing the token, it waits for that to complete.
func (c *Client) refreshToken(ctx context.Context) error {
	maxAttempts := 5
	for i := 0; i < maxAttempts; i++ {
		// Check if token is valid without locking (read-only check)
		c.tokenMutex.Lock()
		tokenValid := c.token.Valid()
		getTokenErr := c.getTokenErr
		c.tokenMutex.Unlock()

		if tokenValid {
			// Token is valid, use it
			return nil
		}
		if getTokenErr != nil {
			// Another goroutine encountered an error
			return getTokenErr
		}

		// Try to acquire the lock for token refresh
		if c.tokenMutex.TryLock() {
			// We got the lock, proceed with token refresh
			defer c.tokenMutex.Unlock()

			// Double-check after acquiring the lock
			if c.token.Valid() {
				return nil
			}
			if c.getTokenErr != nil {
				return c.getTokenErr
			}

			// Refresh the token using refresh_token grant type
			if c.token == nil {
				c.getTokenErr = fmt.Errorf("token is not set: call SetToken() with a token obtained from RetrieveToken()")
				return c.getTokenErr
			}
			if c.token.RefreshToken == nil {
				c.token = nil
				c.getTokenErr = fmt.Errorf("no refresh token available: the current token does not have a refresh token")
				return c.getTokenErr
			}

			grantType := "refresh_token"
			token, err := c.RetrieveToken(ctx, &RetrieveTokenRequest{
				GrantType:    &grantType,
				RefreshToken: c.token.RefreshToken,
			})

			if err != nil {
				c.token = nil
				c.getTokenErr = fmt.Errorf("refresh token error: %w", err)
				return c.getTokenErr
			}
			c.token = token
			c.getTokenErr = nil
			return nil
		}

		// Another goroutine is refreshing the token, wait a bit and retry
		waitTime := time.Duration(100+rand.Intn(100)) * time.Millisecond
		if err := sleepWithContext(ctx, waitTime); err != nil {
			return err
		}
	}
	return fmt.Errorf("max attempts exceeded while waiting for token refresh")
}
