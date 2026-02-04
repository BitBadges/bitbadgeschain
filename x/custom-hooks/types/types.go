package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
)

// SetDeterministicError sets a deterministic error message in the transient store.
// This should be called before returning an error to capture a deterministic error string.
// The error message should be deterministic (no traces, logs, or non-deterministic values).
// Panics if the error message contains patterns indicating stack traces or non-deterministic content.
// Uses transient store so the value is automatically cleared at the end of the transaction.
// SetDeterministicError stores an error message in the transient store for deterministic error handling.
// This function validates that the error message doesn't contain non-deterministic content like stack traces.
//
// Limitations: This uses pattern matching to detect non-deterministic content, which may:
//   - False-positive on legitimate content containing these patterns (e.g., ".go" in a message about Go files)
//   - Miss some non-deterministic patterns not covered by the checks
//
// For production use, consider:
//   - Using structured error types with predefined error codes
//   - Creating an allowlist of safe error message formats
//   - Implementing more sophisticated validation based on error source
func SetDeterministicError(ctx sdk.Context, errorMsg string) {
	// Validate that the error message doesn't contain stack traces or non-deterministic patterns
	
	// Check for file paths (indicated by ".go" or file extensions)
	// Note: This may false-positive on legitimate messages about Go files
	if strings.Contains(errorMsg, ".go") {
		panic(fmt.Sprintf("SetDeterministicError: error message contains '.go' (likely a stack trace), which is non-deterministic. Error message: %s", errorMsg))
	}

	// Check for goroutine IDs (e.g., "goroutine 123")
	if strings.Contains(errorMsg, "goroutine") {
		panic(fmt.Sprintf("SetDeterministicError: error message contains 'goroutine' (likely from stack trace), which is non-deterministic. Error message: %s", errorMsg))
	}

	// Check for runtime package references (runtime.xxx)
	if strings.Contains(errorMsg, "runtime.") {
		panic(fmt.Sprintf("SetDeterministicError: error message contains 'runtime.' (likely from stack trace), which is non-deterministic. Error message: %s", errorMsg))
	}

	// Check for "panic" keyword (often in stack traces)
	if strings.Contains(errorMsg, "panic(") || strings.Contains(errorMsg, "panic:") {
		panic(fmt.Sprintf("SetDeterministicError: error message contains 'panic' (likely from stack trace), which is non-deterministic. Error message: %s", errorMsg))
	}

	// Check for full package paths (github.com/... or similar)
	// Using regex for more precise matching of package paths
	if matched, _ := regexp.MatchString(`(github\.com|golang\.org|go\.pkg\.dev|gopkg\.in)/`, errorMsg); matched {
		panic(fmt.Sprintf("SetDeterministicError: error message contains package path (likely from stack trace), which is non-deterministic. Error message: %s", errorMsg))
	}
	
	// Check for common stack trace patterns: file paths with line numbers (e.g., "file.go:123")
	if matched, _ := regexp.MatchString(`\w+\.go:\d+`, errorMsg); matched {
		panic(fmt.Sprintf("SetDeterministicError: error message contains file path with line number (likely from stack trace), which is non-deterministic. Error message: %s", errorMsg))
	}

	// Store in transient store - this persists across function calls even when context is passed by value
	// Transient stores are automatically cleared at the end of each transaction
	store := ctx.TransientStore(TransientStoreKey)
	store.Set(DeterministicErrorKey, []byte(errorMsg))
}

// GetDeterministicError retrieves the deterministic error message from the transient store, if any.
// Returns the error message and true if found, empty string and false otherwise.
func GetDeterministicError(ctx sdk.Context) (string, bool) {
	store := ctx.TransientStore(TransientStoreKey)
	value := store.Get(DeterministicErrorKey)
	if len(value) == 0 {
		return "", false
	}
	return string(value), true
}

// ClearDeterministicError clears any deterministic error message from the transient store.
// This should be called before starting a new operation to avoid using stale error messages.
func ClearDeterministicError(ctx sdk.Context) {
	store := ctx.TransientStore(TransientStoreKey)
	store.Delete(DeterministicErrorKey)
}

// WrapErr sets a deterministic error in transient store and returns a wrapped error.
// The detMsg is used both for caching and as the format string for the error.
// Usage: return customhookstypes.WrapErr(&ctx, errType, detMsg, args...)
func WrapErr(ctx *sdk.Context, errType error, detMsg string, args ...interface{}) error {
	// Use detMsg as the format string - if it has format verbs, format it first for caching
	var cachedMsg string
	if len(args) > 0 {
		cachedMsg = fmt.Sprintf(detMsg, args...)
	} else {
		cachedMsg = detMsg
	}
	SetDeterministicError(*ctx, cachedMsg)
	// Use Wrap when no args (to avoid linter warning about non-constant format string)
	if len(args) == 0 {
		return errorsmod.Wrap(errType, detMsg)
	}
	return errorsmod.Wrapf(errType, detMsg, args...)
}

// WrapErrSimple sets a deterministic error in transient store and returns a wrapped error.
// This is for cases where detMsg is already a complete string (not a format string).
// Usage: return customhookstypes.WrapErrSimple(&ctx, errType, detMsg)
func WrapErrSimple(ctx *sdk.Context, errType error, detMsg string) error {
	SetDeterministicError(*ctx, detMsg)
	return errorsmod.Wrap(errType, detMsg)
}

// Err sets a deterministic error in transient store and returns a new error.
// The detMsg is used both for caching and as the format string for the error.
// Usage: return customhookstypes.Err(&ctx, detMsg, args...)
func Err(ctx *sdk.Context, detMsg string, args ...interface{}) error {
	// Use detMsg as the format string - if it has format verbs, format it first for caching
	var cachedMsg string
	if len(args) > 0 {
		cachedMsg = fmt.Sprintf(detMsg, args...)
	} else {
		cachedMsg = detMsg
	}
	SetDeterministicError(*ctx, cachedMsg)
	return errorsmod.Wrap(ErrDeterministicError, cachedMsg)
}

// HookData represents the data that can be executed via IBC hooks
type HookData struct {
	SwapAndAction *SwapAndAction `json:"swap_and_action,omitempty"`
}

// SwapAndAction represents the Skip-style format for swap and action
type SwapAndAction struct {
	UserSwap                  *UserSwap       `json:"user_swap,omitempty"`
	MinAsset                  *MinAsset       `json:"min_asset,omitempty"`
	TimeoutTimestamp          *uint64         `json:"timeout_timestamp,omitempty"`
	PostSwapAction            *PostSwapAction `json:"post_swap_action,omitempty"`
	DestinationRecoverAddress string          `json:"destination_recover_address,omitempty"`
	Affiliates                []Affiliate     `json:"affiliates,omitempty"`
}

// UserSwap contains the swap operation
type UserSwap struct {
	SwapExactAssetIn *SwapExactAssetIn `json:"swap_exact_asset_in,omitempty"`
}

// SwapExactAssetIn contains the swap operations
type SwapExactAssetIn struct {
	SwapVenueName string      `json:"swap_venue_name,omitempty"`
	Operations    []Operation `json:"operations"`
}

// Operation represents a single swap operation
type Operation struct {
	Pool     string `json:"pool"` // Pool ID as string
	DenomIn  string `json:"denom_in"`
	DenomOut string `json:"denom_out"`
}

// MinAsset represents the minimum asset output
type MinAsset struct {
	Native *NativeAsset `json:"native,omitempty"`
}

// NativeAsset represents a native asset
type NativeAsset struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

// PostSwapAction contains actions to execute after swap
type PostSwapAction struct {
	IBCTransfer *IBCTransferInfo `json:"ibc_transfer,omitempty"`
	Transfer    *TransferInfo    `json:"transfer,omitempty"`
}

// TransferInfo contains local transfer information
type TransferInfo struct {
	ToAddress string `json:"to_address"`
}

// IBCTransferInfo contains IBC transfer information
type IBCTransferInfo struct {
	IBCInfo *IBCInfo `json:"ibc_info,omitempty"`
}

// IBCInfo contains IBC transfer details
type IBCInfo struct {
	SourceChannel  string `json:"source_channel"`
	Receiver       string `json:"receiver"`
	Memo           string `json:"memo,omitempty"`
	RecoverAddress string `json:"recover_address,omitempty"`
}

// Affiliate represents an affiliate fee recipient
type Affiliate struct {
	BasisPointsFee string `json:"basis_points_fee"`
	Address        string `json:"address"`
}

const (
	// MaxMemoSize is the maximum size (in bytes) for IBC memo data
	// Security: LOW-009 - Prevents DoS attacks via extremely large memo payloads
	// 64KB is a reasonable limit that allows most legitimate memos while preventing abuse
	MaxMemoSize = 64 * 1024 // 64KB
)

// ParseHookDataFromMemo parses hook data from IBC memo
// Security: LOW-009 - Validates memo size before unmarshaling to prevent DoS attacks
func ParseHookDataFromMemo(memo string) (*HookData, error) {
	if len(memo) == 0 {
		return nil, nil
	}

	// Security: Validate memo size before unmarshaling to prevent DoS
	if len(memo) > MaxMemoSize {
		return nil, errorsmod.Wrapf(ErrMemoSizeExceedsMaximum, "size: %d bytes, max: %d bytes", len(memo), MaxMemoSize)
	}

	var memoObj map[string]interface{}
	if err := json.Unmarshal([]byte(memo), &memoObj); err != nil {
		return nil, err
	}

	// Check for swap_and_action key
	if swapAndActionRaw, ok := memoObj["swap_and_action"]; ok {
		swapAndActionBytes, err := json.Marshal(swapAndActionRaw)
		if err != nil {
			return nil, err
		}

		var swapAndAction SwapAndAction
		if err := json.Unmarshal(swapAndActionBytes, &swapAndAction); err != nil {
			return nil, err
		}

		return &HookData{
			SwapAndAction: &swapAndAction,
		}, nil
	}

	return nil, nil
}

// NewCustomErrorAcknowledgement creates a custom error acknowledgement with a deterministic error string
// IMPORTANT: The error string must be deterministic (no traces, logs, or non-deterministic values)
// This is used instead of channeltypes.NewErrorAcknowledgement to provide more friendly error messages
func NewCustomErrorAcknowledgement(errorMsg string) ibcexported.Acknowledgement {
	return channeltypes.Acknowledgement{
		Response: &channeltypes.Acknowledgement_Error{
			Error: fmt.Sprintf("swap-and-action-hooks: %s", errorMsg),
		},
	}
}

// NewSuccessAcknowledgement creates a success acknowledgement
// This is used internally to indicate successful execution
func NewSuccessAcknowledgement() ibcexported.Acknowledgement {
	return channeltypes.NewResultAcknowledgement([]byte("success"))
}

// IsSuccessAcknowledgement checks if an acknowledgement indicates success
func IsSuccessAcknowledgement(ack ibcexported.Acknowledgement) bool {
	return ack.Success()
}

// GetAckError extracts the error message from an acknowledgement.
// Returns the error message and true if the acknowledgement is an error, empty string and false otherwise.
func GetAckError(ack ibcexported.Acknowledgement) (string, bool) {
	channelAck, ok := ack.(channeltypes.Acknowledgement)
	if !ok {
		return "", false
	}
	if errResp, ok := channelAck.Response.(*channeltypes.Acknowledgement_Error); ok {
		return errResp.Error, true
	}
	return "", false
}
