package keeper

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"

	"github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// ExecuteTransferTokens executes a TransferTokensAction hook.
// It mints/transfers BitBadges tokens on behalf of the IBC intermediate sender.
func (k Keeper) ExecuteTransferTokens(ctx sdk.Context, sender sdk.AccAddress, action *types.TransferTokensAction, tokenIn sdk.Coin, originalSender string) ibcexported.Acknowledgement {
	if action == nil {
		return types.NewSuccessAcknowledgement()
	}

	// 1. Early validation: collection_id must be parseable
	collectionId, err := strconv.ParseUint(action.CollectionId, 10, 64)
	if err != nil {
		return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid transfer_tokens collection_id: %s", action.CollectionId))
	}

	// 2. Transfers must be non-empty
	if len(action.Transfers) == 0 || string(action.Transfers) == "[]" || string(action.Transfers) == "null" {
		return types.NewCustomErrorAcknowledgement("transfer_tokens transfers cannot be empty")
	}

	// 3. If fail_on_error is false, recover_address must be provided and valid
	if !action.FailOnError {
		if action.RecoverAddress == "" {
			return types.NewCustomErrorAcknowledgement("recover_address is required when fail_on_error is false")
		}
		_, err := sdk.AccAddressFromBech32(action.RecoverAddress)
		if err != nil {
			return types.NewCustomErrorAcknowledgement(fmt.Sprintf("invalid recover_address: %s", action.RecoverAddress))
		}
	}

	// 4. Get collection and set auto-approval flags for intermediate sender
	collection, found := k.tokenizationKeeper.GetCollectionFromStore(ctx, sdkmath.NewUint(collectionId))
	if !found {
		return types.NewCustomErrorAcknowledgement(fmt.Sprintf("transfer_tokens collection not found: %s", action.CollectionId))
	}

	// Set all auto-approval flags for the intermediate sender on this collection.
	// Uses main ctx (not cacheCtx) intentionally — these flags are idempotent and
	// intermediate addresses are NEVER controllable by users, so persisting them
	// even on transfer failure is safe and avoids redundant writes on retry.
	err = k.tokenizationKeeper.SetAllAutoApprovalFlagsForIntermediateAddress(ctx, collection, sender.String())
	if err != nil {
		return types.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to set auto-approve for intermediate address on collection %s", action.CollectionId))
	}

	// 5. Unmarshal snake_case JSON transfers directly into protobuf types
	protoTransfers, err := types.UnmarshalTransfersFromJSON(action.Transfers)
	if err != nil {
		return types.NewCustomErrorAcknowledgement(fmt.Sprintf("failed to parse transfer_tokens transfers: %s", err.Error()))
	}
	if len(protoTransfers) == 0 {
		return types.NewCustomErrorAcknowledgement("transfer_tokens transfers cannot be empty")
	}

	// 6. Construct MsgTransferTokens
	msg := &tokenizationtypes.MsgTransferTokens{
		Creator:      sender.String(),
		CollectionId: sdkmath.NewUint(collectionId),
		Transfers:    protoTransfers,
	}

	// 7. Execute in a cached context for atomicity
	cacheCtx, writeCache := ctx.CacheContext()
	types.ClearDeterministicError(cacheCtx)

	_, err = k.tokenizationMsgServer.TransferTokens(cacheCtx, msg)
	if err != nil {
		// Transfer failed
		transferErrMsg := fmt.Sprintf("transfer_tokens failed: collection=%s, sender=%s (derived from %s)", action.CollectionId, sender.String(), originalSender)

		// Only use deterministic error from transient store — never use err.Error() directly
		// as it may contain stack traces or non-deterministic content that would halt consensus
		if detErrMsg, found := types.GetDeterministicError(cacheCtx); found {
			transferErrMsg = fmt.Sprintf("%s: %s", transferErrMsg, detErrMsg)
		}

		if action.FailOnError {
			// Return error ack — IBC will refund tokens to source chain
			return types.NewCustomErrorAcknowledgement(transferErrMsg)
		}

		// fail_on_error is false — send IBC tokens to recover_address
		extraAttrs := []sdk.Attribute{
			sdk.NewAttribute("collection_id", action.CollectionId),
		}
		return k.ExecuteFallbackToRecoverAddress(ctx, sender, action.RecoverAddress, tokenIn, originalSender, transferErrMsg, "transfer_tokens_fallback", extraAttrs)
	}

	// 8. Transfer succeeded — commit cache
	//
	// tokenIn disposition: the IBC coins (tokenIn) are intentionally retained at
	// the intermediate sender address as payment for this hook operation. They are
	// NOT forwarded or burned here. On the failure path, ExecuteFallbackToRecoverAddress
	// handles recovery by sending them on to RecoverAddress instead. For non-payment
	// use cases where no fee is required, callers should send a trivial (1uatom) amount.
	writeCache()

	// Emit success event. token_in_denom and token_in_amount are emitted as
	// separate attributes (in addition to the combined token_in string) so that
	// indexers can filter/aggregate by denom or amount independently.
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"transfer_tokens_success",
		sdk.NewAttribute("module", "custom-hooks"),
		sdk.NewAttribute("sender", sender.String()),
		sdk.NewAttribute("original_sender", originalSender),
		sdk.NewAttribute("collection_id", action.CollectionId),
		sdk.NewAttribute("num_transfers", strconv.Itoa(len(protoTransfers))),
		sdk.NewAttribute("token_in", tokenIn.String()),
		sdk.NewAttribute("token_in_denom", tokenIn.Denom),
		sdk.NewAttribute("token_in_amount", tokenIn.Amount.String()),
	))

	return types.NewSuccessAcknowledgement()
}
