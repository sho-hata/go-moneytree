package moneytree

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// SubmitAccount2FAKeyValues represents the key-values for 2FA submission.
// This object must contain exactly one key-value pair: either "otp" or "captcha", but not both.
type SubmitAccount2FAKeyValues struct {
	// OTP is the one-time password value.
	// This field should be set when the account's authentication status is "suspended.missing-answer.auth.otp".
	// Maximum length is 255 characters.
	OTP *string `json:"otp,omitempty"`
	// Captcha is the CAPTCHA answer value.
	// This field should be set when the account's authentication status is "suspended.missing-answer.auth.captcha".
	// Maximum length is 255 characters.
	Captcha *string `json:"captcha,omitempty"`
}

// SubmitAccount2FARequest represents a request to submit 2FA information for an account.
type SubmitAccount2FARequest struct {
	// KeyValues contains the 2FA information.
	// This object must contain exactly one key-value pair: either "otp" or "captcha", but not both.
	KeyValues SubmitAccount2FAKeyValues `json:"key_values"`
}

// SubmitAccount2FA submits 2FA (two-factor authentication) information for an account that requires additional authentication.
// This endpoint requires the accounts_read OAuth scope.
//
// This API is used when an account's authentication status requires OTP (one-time password) or CAPTCHA input.
// The request must contain exactly one key-value pair in key_values: either "otp" or "captcha", but not both.
//
// Account status requirements:
//   - For OTP submission: The account's authentication status must be "suspended.missing-answer.auth.otp"
//   - For CAPTCHA submission: The account's authentication status must be "suspended.missing-answer.auth.captcha"
//
// Upon successful submission, the authentication workflow resumes and the account status changes to "running.auth".
//
// Example with OTP:
//
//	otp := "123456"
//	request := &moneytree.SubmitAccount2FARequest{
//		KeyValues: moneytree.SubmitAccount2FAKeyValues{
//			OTP: &otp,
//		},
//	}
//	err := client.SubmitAccount2FA(ctx, accessToken, "account_key_123", request)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Example with CAPTCHA:
//
//	captcha := "captcha_answer"
//	request := &moneytree.SubmitAccount2FARequest{
//		KeyValues: moneytree.SubmitAccount2FAKeyValues{
//			Captcha: &captcha,
//		},
//	}
//	err := client.SubmitAccount2FA(ctx, accessToken, "account_key_123", request)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Reference: https://docs.link.getmoneytree.com/reference/put-account-2fa
func (c *Client) SubmitAccount2FA(ctx context.Context, accountID string, req *SubmitAccount2FARequest) error {
	if accountID == "" {
		return fmt.Errorf("account ID is required")
	}
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate that exactly one of OTP or Captcha is set
	hasOTP := req.KeyValues.OTP != nil
	hasCaptcha := req.KeyValues.Captcha != nil

	if !hasOTP && !hasCaptcha {
		return fmt.Errorf("key_values must contain either 'otp' or 'captcha', but both are missing")
	}
	if hasOTP && hasCaptcha {
		return fmt.Errorf("key_values must contain either 'otp' or 'captcha', but not both")
	}

	// Validate maximum length
	if hasOTP && len(*req.KeyValues.OTP) > 255 {
		return fmt.Errorf("otp must be 255 characters or less, got %d characters", len(*req.KeyValues.OTP))
	}
	if hasCaptcha && len(*req.KeyValues.Captcha) > 255 {
		return fmt.Errorf("captcha must be 255 characters or less, got %d characters", len(*req.KeyValues.Captcha))
	}

	urlPath := fmt.Sprintf("link/accounts/%s/2fa.json", url.PathEscape(accountID))

	httpReq, err := c.NewRequest(ctx, http.MethodPut, urlPath, req)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if _, err := c.Do(ctx, httpReq, nil); err != nil {
		return err
	}
	return nil
}

// CaptchaImage represents the CAPTCHA image information returned by the Moneytree LINK API.
type CaptchaImage struct {
	// URL is the URL pointing to the CAPTCHA image that needs to be solved.
	// Users should download the image from this URL, display it to end users,
	// collect text input, and submit it via the SubmitAccount2FA endpoint.
	URL string `json:"url"`
}

// GetAccountCaptcha retrieves the CAPTCHA image URL for an account that requires CAPTCHA authentication.
// This endpoint requires the accounts_read OAuth scope.
//
// This API is used when an account's authentication status is "suspended.missing-answer.auth.captcha".
// The response contains a URL pointing to the CAPTCHA image. Users should:
// 1. Download the image from the provided URL
// 2. Display the image to end users
// 3. Collect text input from users
// 4. Submit the text via the SubmitAccount2FA endpoint
//
// Account status requirement:
//   - The account's authentication status must be "suspended.missing-answer.auth.captcha"
//   - If the account is not in the correct status, 400 Bad Request will be returned
//
// Example:
//
//	captchaImage, err := client.GetAccountCaptcha(ctx, accessToken, "account_key_123")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("CAPTCHA image URL: %s\n", captchaImage.URL)
//
// Reference: https://docs.link.getmoneytree.com/reference/get-account-captcha
func (c *Client) GetAccountCaptcha(ctx context.Context, accountID string) (*CaptchaImage, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	urlPath := fmt.Sprintf("link/accounts/%s/captcha.json", url.PathEscape(accountID))

	httpReq, err := c.NewRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var res CaptchaImage
	if _, err = c.Do(ctx, httpReq, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
