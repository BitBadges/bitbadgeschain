package tokenization

import (
	"fmt"
	"strings"

	sdkerrors "cosmossdk.io/errors"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
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
	// ErrorCodeCollectionNotFound represents collection not found
	ErrorCodeCollectionNotFound ErrorCode = 2
	// ErrorCodeBalanceNotFound represents balance not found
	ErrorCodeBalanceNotFound ErrorCode = 3
	// ErrorCodeTransferFailed represents transfer operation failed
	ErrorCodeTransferFailed ErrorCode = 4
	// ErrorCodeApprovalFailed represents approval operation failed
	ErrorCodeApprovalFailed ErrorCode = 5
	// ErrorCodeQueryFailed represents query operation failed
	ErrorCodeQueryFailed ErrorCode = 6
	// ErrorCodeInternalError represents internal error
	ErrorCodeInternalError ErrorCode = 7
	// ErrorCodeUnauthorized represents unauthorized operation
	ErrorCodeUnauthorized ErrorCode = 8
	// ErrorCodeCollectionArchived represents collection is archived (read-only)
	ErrorCodeCollectionArchived ErrorCode = 9
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

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

// findSubstring returns the index of the first occurrence of substr in s, or -1 if not found
func findSubstring(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(s) < len(substr) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// replaceFirst replaces the first occurrence of old with new in s
func replaceFirst(s, old, new string) string {
	idx := findSubstring(s, old)
	if idx == -1 {
		return s
	}
	return s[:idx] + new + s[idx+len(old):]
}

// MapCosmosErrorToPrecompileError maps Cosmos SDK errors to appropriate precompile error codes
// Returns the mapped error code and a boolean indicating if a mapping was found
func MapCosmosErrorToPrecompileError(err error) (ErrorCode, bool) {
	if err == nil {
		return 0, false
	}

	// Check keeper errors
	if sdkerrors.IsOf(err, tokenizationkeeper.ErrCollectionNotExists) {
		return ErrorCodeCollectionNotFound, true
	}
	if sdkerrors.IsOf(err, tokenizationkeeper.ErrUserBalanceNotExists) {
		return ErrorCodeBalanceNotFound, true
	}
	if sdkerrors.IsOf(err, tokenizationkeeper.ErrInadequateApprovals) {
		return ErrorCodeUnauthorized, true
	}
	if sdkerrors.IsOf(err, tokenizationkeeper.ErrCollectionIsArchived) {
		return ErrorCodeCollectionArchived, true
	}
	if sdkerrors.IsOf(err, tokenizationkeeper.ErrDisallowedTransfer) {
		return ErrorCodeTransferFailed, true
	}
	if sdkerrors.IsOf(err, tokenizationkeeper.ErrAccountNotFound) {
		return ErrorCodeQueryFailed, true
	}
	if sdkerrors.IsOf(err, tokenizationkeeper.ErrAddressListNotFound) {
		return ErrorCodeQueryFailed, true
	}
	if sdkerrors.IsOf(err, tokenizationkeeper.ErrApprovalNotFound) {
		return ErrorCodeQueryFailed, true
	}

	// Check types errors
	if sdkerrors.IsOf(err, tokenizationtypes.ErrUnauthorized) {
		return ErrorCodeUnauthorized, true
	}
	if sdkerrors.IsOf(err, tokenizationtypes.ErrNotFound) {
		return ErrorCodeQueryFailed, true
	}
	if sdkerrors.IsOf(err, tokenizationtypes.ErrInvalidCollectionID) {
		return ErrorCodeCollectionNotFound, true
	}
	if sdkerrors.IsOf(err, tokenizationtypes.ErrInvalidRequest) {
		return ErrorCodeInvalidInput, true
	}
	if sdkerrors.IsOf(err, tokenizationtypes.ErrUnderflow) {
		return ErrorCodeTransferFailed, true
	}
	if sdkerrors.IsOf(err, tokenizationtypes.ErrOverflow) {
		return ErrorCodeTransferFailed, true
	}
	if sdkerrors.IsOf(err, tokenizationtypes.ErrInvalidAddress) {
		return ErrorCodeInvalidInput, true
	}
	if sdkerrors.IsOf(err, tokenizationtypes.ErrInvalidCollectionID) {
		return ErrorCodeInvalidInput, true
	}
	if sdkerrors.IsOf(err, tokenizationtypes.ErrAmountEqualsZero) {
		return ErrorCodeInvalidInput, true
	}

	return 0, false
}

// FieldPathBuilder helps build field paths for nested structures
// Example: Field("defaultBalances").Field("balances").Index(0).Field("amount") -> "defaultBalances.balances[0].amount"
type FieldPathBuilder struct {
	path []string
}

// NewFieldPathBuilder creates a new FieldPathBuilder
func NewFieldPathBuilder() *FieldPathBuilder {
	return &FieldPathBuilder{path: make([]string, 0)}
}

// Field adds a field name to the path
func (b *FieldPathBuilder) Field(name string) *FieldPathBuilder {
	b.path = append(b.path, name)
	return b
}

// Index adds an array index to the path
func (b *FieldPathBuilder) Index(i int) *FieldPathBuilder {
	b.path = append(b.path, fmt.Sprintf("[%d]", i))
	return b
}

// Key adds a map key to the path
func (b *FieldPathBuilder) Key(key string) *FieldPathBuilder {
	b.path = append(b.path, fmt.Sprintf("[%q]", key))
	return b
}

// String returns the complete field path as a string
func (b *FieldPathBuilder) String() string {
	if len(b.path) == 0 {
		return ""
	}
	result := b.path[0]
	for i := 1; i < len(b.path); i++ {
		// If the current segment starts with '[', it's an index/key, so don't add a dot
		if len(b.path[i]) > 0 && b.path[i][0] == '[' {
			result += b.path[i]
		} else {
			result += "." + b.path[i]
		}
	}
	return result
}

// Clone creates a copy of the FieldPathBuilder
func (b *FieldPathBuilder) Clone() *FieldPathBuilder {
	newPath := make([]string, len(b.path))
	copy(newPath, b.path)
	return &FieldPathBuilder{path: newPath}
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

// WrapErrorWithContext wraps an error with additional context about where it occurred
func WrapErrorWithContext(err error, defaultCode ErrorCode, message string, context string) *PrecompileError {
	details := ""
	if err != nil {
		details = err.Error()
	}
	if context != "" {
		details = fmt.Sprintf("[%s] %s", context, details)
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

func ErrCollectionNotFound(collectionId string) *PrecompileError {
	return NewPrecompileError(ErrorCodeCollectionNotFound, "collection not found", fmt.Sprintf("collectionId: %s", collectionId))
}

func ErrBalanceNotFound(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeBalanceNotFound, "balance not found", details)
}

func ErrTransferFailed(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeTransferFailed, "transfer failed", details)
}

func ErrApprovalFailed(details string) *PrecompileError {
	return NewPrecompileError(ErrorCodeApprovalFailed, "approval operation failed", details)
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

func ErrCollectionArchived(collectionId string) *PrecompileError {
	return NewPrecompileError(ErrorCodeCollectionArchived, "collection is archived (read-only)", fmt.Sprintf("collectionId: %s", collectionId))
}

// Enhanced error helpers with field path support

// ErrInvalidField creates an error with field path
func ErrInvalidField(fieldPath string, reason string, expected string) *PrecompileError {
	details := fmt.Sprintf("field '%s': %s", fieldPath, reason)
	if expected != "" {
		details += fmt.Sprintf(" (expected: %s)", expected)
	}
	return NewPrecompileError(ErrorCodeInvalidInput, "invalid input parameters", details)
}

// ErrInvalidFieldValue creates an error for invalid field values
func ErrInvalidFieldValue(fieldPath string, value interface{}, reason string) *PrecompileError {
	details := fmt.Sprintf("field '%s': invalid value '%v' - %s", fieldPath, value, reason)
	return NewPrecompileError(ErrorCodeInvalidInput, "invalid input parameters", details)
}

// ErrMissingRequiredField creates an error for missing required fields
func ErrMissingRequiredField(fieldPath string) *PrecompileError {
	details := fmt.Sprintf("field '%s': required field is missing or empty", fieldPath)
	return NewPrecompileError(ErrorCodeInvalidInput, "invalid input parameters", details)
}

// ErrFieldTypeMismatch creates an error for type mismatches
func ErrFieldTypeMismatch(fieldPath string, got interface{}, expected string) *PrecompileError {
	gotType := "unknown"
	if got != nil {
		gotType = fmt.Sprintf("%T", got)
	}
	details := fmt.Sprintf("field '%s': type mismatch - got %s, expected %s", fieldPath, gotType, expected)
	return NewPrecompileError(ErrorCodeInvalidInput, "invalid input parameters", details)
}
