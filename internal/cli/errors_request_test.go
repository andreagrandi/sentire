package cli

import (
	"testing"

	"github.com/andreagrandi/sentire/internal/client"
)

func TestWrapErrorTimeout(t *testing.T) {
	err := wrapError(&client.RequestError{Message: "request timed out", Timeout: true})

	cliErr, ok := err.(*CLIError)
	if !ok {
		t.Fatalf("Expected *CLIError, got %T", err)
	}
	if cliErr.Code != CodeTimeout {
		t.Errorf("Expected code %q, got %q", CodeTimeout, cliErr.Code)
	}
	if cliErr.ExitCode != ExitAPI {
		t.Errorf("Expected exit code %d, got %d", ExitAPI, cliErr.ExitCode)
	}
}

func TestWrapErrorCanceled(t *testing.T) {
	err := wrapError(&client.RequestError{Message: "request canceled", Canceled: true})

	cliErr, ok := err.(*CLIError)
	if !ok {
		t.Fatalf("Expected *CLIError, got %T", err)
	}
	if cliErr.Code != CodeCanceled {
		t.Errorf("Expected code %q, got %q", CodeCanceled, cliErr.Code)
	}
	if cliErr.ExitCode != ExitAPI {
		t.Errorf("Expected exit code %d, got %d", ExitAPI, cliErr.ExitCode)
	}
}

func TestWrapErrorGenericRequestError(t *testing.T) {
	// A transport failure that is neither a timeout nor a cancellation keeps
	// the general exit code, matching behavior before request errors were typed.
	reqErr := &client.RequestError{Message: "http request failed: connection refused"}

	err := wrapError(reqErr)
	if _, ok := err.(*CLIError); ok {
		t.Error("Expected a generic request error not to be wrapped as a CLIError")
	}
	if got := exitCodeFromError(err); got != ExitGeneral {
		t.Errorf("Expected exit code %d, got %d", ExitGeneral, got)
	}
}
