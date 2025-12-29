package moneytree

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error returned by the Moneytree LINK API.
type APIError struct {
	StatusCode int `json:"-"`
	// ErrorType is the value of the error field set by moneytree.
	// It is empty when an unexpected error occurs during response decoding.
	ErrorType string `json:"error,omitempty"`
	// ErrorDescription is the value of the error_description field set by moneytree.
	// However, if an unexpected error occurs during response decoding, it contains a message set by the library.
	ErrorDescription string `json:"error_description,omitempty"`
	RawMessage       string `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.ErrorDescription != "" {
		if e.ErrorType != "" {
			return fmt.Sprintf("%d: %s - %s", e.StatusCode, e.ErrorType, e.ErrorDescription)
		}
		return fmt.Sprintf("%d: %s", e.StatusCode, e.ErrorDescription)
	}
	if e.ErrorType != "" {
		return fmt.Sprintf("%d: %s", e.StatusCode, e.ErrorType)
	}
	return fmt.Sprintf("%d", e.StatusCode)
}

// checks the response, and in case of error, maps it to the error structure.
func checkResponseError(r *http.Response) error {
	if r == nil {
		return errors.New("response cannot be nil")
	}

	if !isErrorStatusCode(r.StatusCode) {
		return nil
	}

	apiErr := &APIError{
		StatusCode: r.StatusCode,
	}

	if r.Body != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return &APIError{
				StatusCode:       r.StatusCode,
				ErrorDescription: fmt.Sprintf("unable to read response from moneytree: %s", err.Error()),
			}
		}

		apiErr.RawMessage = string(body)

		if err := json.Unmarshal(body, apiErr); err != nil {
			return &APIError{
				StatusCode:       r.StatusCode,
				ErrorDescription: fmt.Sprintf("unable to decode response from moneytree: %s", err.Error()),
				RawMessage:       string(body),
			}
		}
	}
	return apiErr
}

func isErrorStatusCode(statusCode int) bool {
	return statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError
}
