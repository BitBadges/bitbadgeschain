package gamm

import (
	"fmt"
	"strings"

	sdkerrors "cosmossdk.io/errors"

	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
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
	// ErrorCodePoolNotFound represents pool not found
	ErrorCodePoolNotFound ErrorCode = 2
	// ErrorCodeSwapFailed represents swap operation failed
	ErrorCodeSwapFailed ErrorCode = 3
	// ErrorCodeQueryFailed represents query operation failed
	ErrorCodeQueryFailed ErrorCode = 4
	// ErrorCodeInternalError represents internal error
	ErrorCodeInternalError ErrorCode = 5
	// ErrorCodeUnauthorized represents unauthorized operation
	ErrorCodeUnauthorized ErrorCode = 6
	// ErrorCodeJoinPoolFailed represents join pool operation failed
	ErrorCodeJoinPoolFailed ErrorCode = 7
	// ErrorCodeExitPoolFailed represents exit pool operation failed
	ErrorCodeExitPoolFailed ErrorCode = 8
	// ErrorCodeIBCTransferFailed represents IBC transfer operation failed
	ErrorCodeIBCTransferFailed ErrorCode = 9
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

	// Check keeper errors
	if sdkerrors.IsOf(err, gammtypes.ErrPoolNotFound) {
		return ErrorCodePoolNotFound, true
	}
	if sdkerrors.IsOf(err, gammtypes.ErrLimitMaxAmount) {
		return ErrorCodeSwapFailed, true
	}
	if sdkerrors.IsOf(err, gammtypes.ErrLimitMinAmount) {
		return ErrorCodeSwapFailed, true
	}
	if sdkerrors.IsOf(err, gammtypes.ErrNotPositiveRequireAmount) {
		return ErrorCodeInvalidInput, true
	}
	if sdkerrors.IsOf(err, gammtypes.ErrEmptyRoutes) {
		return ErrorCodeInvalidInput, true
	}
	if sdkerrors.IsOf(err, gammtypes.ErrDenomNotFoundInPool) {
		return ErrorCodeSwapFailed, true
	}
	if sdkerrors.IsOf(err, gammtypes.ErrInvalidMathApprox) {
		return ErrorCodeSwapFailed, true
	}
	if sdkerrors.IsOf(err, gammtypes.ErrTooManyTokensOut) {
		return ErrorCodeExitPoolFailed, true
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

func ErrPoolNotFound(poolId uint64) *PrecompileError {
	return NewPrecompileError(ErrorCodePoolNotFound, "pool not found", fmt.Sprintf("poolId: %d", poolId))
}

func ErrSwapFailed(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeSwapFailed, "swap failed", details)
}

func ErrQueryFailed(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeQueryFailed, "query failed", details)
}

func ErrInternalError(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeInternalError, "internal error", details)
}

func ErrUnauthorized(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeUnauthorized, "unauthorized operation", details)
}

func ErrJoinPoolFailed(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeJoinPoolFailed, "join pool failed", details)
}

func ErrExitPoolFailed(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeExitPoolFailed, "exit pool failed", details)
}

func ErrIBCTransferFailed(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeIBCTransferFailed, "IBC transfer failed", details)
}
