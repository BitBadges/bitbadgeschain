package precompile

import (
	"fmt"
	"strings"

	sdkerrors "cosmossdk.io/errors"

	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	sendmanagertypes "github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

// PrecompileError represents a structured error from the precompile
type PrecompileError struct {
	Code    ErrorCode
	Message string
	Details string
}

// ErrorCode represents different error categories
type ErrorCode uint32

const (
	// ErrorCodeInvalidInput represents invalid input parameters
	ErrorCodeInvalidInput ErrorCode = 1
	// ErrorCodeSendFailed represents send operation failed
	ErrorCodeSendFailed ErrorCode = 2
	// ErrorCodeInsufficientBalance represents insufficient balance
	ErrorCodeInsufficientBalance ErrorCode = 3
	// ErrorCodeInternalError represents internal error
	ErrorCodeInternalError ErrorCode = 4
	// ErrorCodeUnauthorized represents unauthorized operation
	ErrorCodeUnauthorized ErrorCode = 5
)

// Error implements the error interface
func (e *PrecompileError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("precompile error [code=%d]: %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("precompile error [code=%d]: %s", e.Code, e.Message)
}

// NewPrecompileError creates a new PrecompileError
func NewPrecompileError(code ErrorCode, message string, details string) *PrecompileError {
	return &PrecompileError{
		Code:    code,
		Message: message,
		Details: sanitizeErrorDetails(details),
	}
}

// sanitizeErrorDetails removes sensitive information from error details
func sanitizeErrorDetails(details string) string {
	if details == "" {
		return ""
	}

	// Remove potential sensitive paths and internal details
	sanitized := details

	// List of sensitive patterns to remove or redact
	sensitivePatterns := []struct {
		pattern     string
		replacement string
	}{
		// File paths (Unix-style)
		{"/home/", "[path]/"},
		{"/root/", "[path]/"},
		{"/usr/", "[path]/"},
		{"/var/", "[path]/"},
		{"/etc/", "[path]/"},
		{"/tmp/", "[path]/"},
		// File paths (Windows-style)
		{"C:\\", "[path]\\"},
		{"D:\\", "[path]\\"},
		// Go-specific internals
		{".go:", "[file]:"},
		// Stack trace indicators
		{"goroutine ", "[goroutine] "},
		{"panic:", "[panic]:"},
		{"runtime.", "[runtime]."},
		// Module paths
		{"github.com/bitbadges/", "[module]/"},
		{"github.com/cosmos/", "[module]/"},
		// IP addresses (simple pattern)
		{"127.0.0.1", "[localhost]"},
		{"0.0.0.0", "[anyaddr]"},
	}

	// Use strings.ReplaceAll for better performance
	for _, sp := range sensitivePatterns {
		sanitized = strings.ReplaceAll(sanitized, sp.pattern, sp.replacement)
	}

	// Truncate very long error messages that might contain stack traces
	const maxLength = 500
	if len(sanitized) > maxLength {
		sanitized = sanitized[:maxLength] + "... [truncated]"
	}

	return sanitized
}

// MapCosmosErrorToPrecompileError maps Cosmos SDK errors to appropriate precompile error codes
// Returns the mapped error code and a boolean indicating if a mapping was found
func MapCosmosErrorToPrecompileError(err error) (ErrorCode, bool) {
	if err == nil {
		return 0, false
	}

	// Check cosmos-sdk error types
	if sdkerrors.IsOf(err, errortypes.ErrInsufficientFunds) {
		return ErrorCodeInsufficientBalance, true
	}
	if sdkerrors.IsOf(err, errortypes.ErrInvalidCoins) {
		return ErrorCodeInvalidInput, true
	}
	if sdkerrors.IsOf(err, errortypes.ErrInvalidAddress) {
		return ErrorCodeInvalidInput, true
	}

	// Check sendmanager errors
	if sdkerrors.IsOf(err, sendmanagertypes.ErrInvalidRequest) {
		return ErrorCodeInvalidInput, true
	}

	return 0, false
}

// WrapError wraps a standard error into a PrecompileError with appropriate code
// It first attempts to map Cosmos SDK errors to specific error codes
// If no mapping is found, it uses the provided default code
func WrapError(err error, defaultCode ErrorCode, message string) *PrecompileError {
	details := ""
	if err != nil {
		details = err.Error()
	}

	// Try to map Cosmos SDK error to precompile error code
	if mappedCode, found := MapCosmosErrorToPrecompileError(err); found {
		return NewPrecompileError(mappedCode, message, details)
	}

	// Fall back to provided default code
	return NewPrecompileError(defaultCode, message, details)
}

// Common error constructors
func ErrInvalidInput(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeInvalidInput, "invalid input parameters", details)
}

func ErrSendFailed(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeSendFailed, "send failed", details)
}

func ErrInsufficientBalance(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeInsufficientBalance, "insufficient balance", details)
}

func ErrInternalError(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeInternalError, "internal error", details)
}

func ErrUnauthorized(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeUnauthorized, "unauthorized operation", details)
}

