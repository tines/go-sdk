package tines

import (
	"fmt"
	"strings"
)

const (
	errEmptyApiKey         = "API Token must not be empty"
	errEmptyTenant         = "Tines Tenant must not be empty"
	errMalformedTenant     = "Tines Tenant must be in the format https://example.tines.com/"
	errInternalServerError = "internal server error"
	errDoRequestError      = "error while attempting to make the HTTP request"
	errUnmarshalError      = "error unmarshalling the JSON response"
	errReadBodyError       = "error reading the HTTP response body bytes"
	errParseError          = "error parsing the input"
)

type ErrorType string

const (
	ErrorTypeRequest        ErrorType = "request"
	ErrorTypeAuthentication ErrorType = "authentication"
	ErrorTypeAuthorization  ErrorType = "authorization"
	ErrorTypeNotFound       ErrorType = "not_found"
	ErrorTypeRateLimit      ErrorType = "rate_limit"
	ErrorTypeServer         ErrorType = "server"
)

type Error struct {
	Type       ErrorType      `json:"type,omitempty"`
	StatusCode int            `json:"status_code,omitempty"`
	Errors     []ErrorMessage `json:"errors,omitempty"`
}

type ErrorMessage struct {
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}

// Implements the standard Go error interface for compatibility with default errors.
func (e Error) Error() string {
	var errString string
	errMessages := []string{}
	errCount := 0

	for _, err := range e.Errors {
		if err.Message != "" {
			msg := fmt.Sprintf("%s: %s", err.Message, err.Details)
			errMessages = append(errMessages, msg)
			errCount += 1
		}
	}
	errString = fmt.Sprintf("%d error(s) occurred: %s", errCount, strings.Join(errMessages, ", "))
	return errString
}

// Check to see if an instantiated error object has more than zero `ErrorMessages` that have been
// appended to it.
func (e Error) HasErrors() bool {
	return e.Errors != nil
}
