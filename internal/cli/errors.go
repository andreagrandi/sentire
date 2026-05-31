package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andreagrandi/sentire/internal/cli/formatter"
	"github.com/andreagrandi/sentire/internal/client"
	"github.com/andreagrandi/sentire/internal/config"
	"github.com/andreagrandi/sentire/internal/redact"
	"io"
	"os"
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
	CodeTimeout       = "timeout"
	CodeCanceled      = "canceled"
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

// NewTimeoutError creates a request timeout error
func NewTimeoutError(message string) *CLIError {
	return &CLIError{
		Message:  message,
		Code:     CodeTimeout,
		ExitCode: ExitAPI,
	}
}

// NewCanceledError creates a request cancellation error
func NewCanceledError(message string) *CLIError {
	return &CLIError{
		Message:  message,
		Code:     CodeCanceled,
		ExitCode: ExitAPI,
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
	var reqErr *client.RequestError
	if errors.As(err, &reqErr) {
		if reqErr.Timeout {
			return NewTimeoutError(reqErr.Message)
		}
		if reqErr.Canceled {
			return NewCanceledError(reqErr.Message)
		}
		// Other transport failures keep the general exit code, matching
		// the behavior before request errors were typed.
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

// writeErrorOutput writes the error to stderr in the appropriate format.
// All output is run through token redaction as a defense-in-depth measure so
// that even an error message constructed outside the client never leaks the
// Sentry API token.
func writeErrorOutput(w io.Writer, err error, format string) {
	wrapped := wrapError(err)
	secret := loadSecretForRedaction()
	if cliErr, ok := wrapped.(*CLIError); ok && (format == "json" || format == "ndjson") {
		cliErr.Message = redact.Secret(cliErr.Message, secret)
		json.NewEncoder(w).Encode(cliErr)
	} else {
		fmt.Fprintf(w, "Error: %s\n", redact.Secret(err.Error(), secret))
	}
}

// loadSecretForRedaction returns the configured Sentry API token, falling back
// to the config file when the env var is unset. Returns "" if neither is set,
// in which case redaction is a no-op.
func loadSecretForRedaction() string {
	if t := os.Getenv("SENTRY_API_TOKEN"); t != "" {
		return t
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		return ""
	}
	return cfg.SentryAPIToken
}

// exitCodeFromError returns the appropriate exit code for an error
func exitCodeFromError(err error) int {
	wrapped := wrapError(err)
	if cliErr, ok := wrapped.(*CLIError); ok {
		return cliErr.ExitCode
	}
	return ExitGeneral
}
