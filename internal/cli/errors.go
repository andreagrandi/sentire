package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sentire/internal/cli/formatter"
	"sentire/internal/client"
	"sentire/internal/config"
)

// Exit codes for different error categories
const (
	ExitSuccess       = 0
	ExitGeneral       = 1
	ExitAuth          = 2
	ExitAPI           = 3
	ExitInvalidInput  = 4
	ExitInvalidFormat = 4
)

// Error codes for structured error output
const (
	CodeAuthMissing   = "auth_missing"
	CodeAPIError      = "api_error"
	CodeInvalidInput  = "invalid_input"
	CodeInvalidFormat = "invalid_format"
)

// CLIError represents a structured error with a machine-readable code
type CLIError struct {
	Message  string `json:"error"`
	Code     string `json:"code"`
	ExitCode int    `json:"-"`
}

func (e *CLIError) Error() string {
	return e.Message
}

// NewAuthError creates an authentication error
func NewAuthError(message string) *CLIError {
	return &CLIError{
		Message:  message,
		Code:     CodeAuthMissing,
		ExitCode: ExitAuth,
	}
}

// NewAPIError creates an API error
func NewAPIError(message string) *CLIError {
	return &CLIError{
		Message:  message,
		Code:     CodeAPIError,
		ExitCode: ExitAPI,
	}
}

// NewInvalidInputError creates an invalid input error
func NewInvalidInputError(message string) *CLIError {
	return &CLIError{
		Message:  message,
		Code:     CodeInvalidInput,
		ExitCode: ExitInvalidInput,
	}
}

// NewInvalidFormatError creates an invalid format error
func NewInvalidFormatError(message string) *CLIError {
	return &CLIError{
		Message:  message,
		Code:     CodeInvalidFormat,
		ExitCode: ExitInvalidFormat,
	}
}

// wrapError converts known error types into CLIError
func wrapError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*CLIError); ok {
		return err
	}
	var authErr *config.AuthError
	if errors.As(err, &authErr) {
		return NewAuthError(authErr.Message)
	}
	var apiErr *client.APIError
	if errors.As(err, &apiErr) {
		return NewAPIError(apiErr.Message)
	}
	var fmtErr *formatter.FormatError
	if errors.As(err, &fmtErr) {
		return NewInvalidFormatError(fmtErr.Message)
	}
	return err
}

// writeErrorOutput writes the error to stderr in the appropriate format
func writeErrorOutput(w io.Writer, err error, format string) {
	wrapped := wrapError(err)
	if cliErr, ok := wrapped.(*CLIError); ok && format == "json" {
		json.NewEncoder(w).Encode(cliErr)
	} else {
		fmt.Fprintf(w, "Error: %v\n", err)
	}
}

// exitCodeFromError returns the appropriate exit code for an error
func exitCodeFromError(err error) int {
	wrapped := wrapError(err)
	if cliErr, ok := wrapped.(*CLIError); ok {
		return cliErr.ExitCode
	}
	return ExitGeneral
}
