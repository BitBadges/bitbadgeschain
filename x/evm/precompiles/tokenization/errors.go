package tokenization

import (
	"fmt"

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
	// Remove potential sensitive paths or internal details
	// This is a basic implementation - can be enhanced
	if details == "" {
		return ""
	}
	// For now, we keep details but could filter out sensitive info
	return details
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

