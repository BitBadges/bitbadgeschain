package tokenization

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// EmitTransferEvent emits an event for a token transfer via precompile
func EmitTransferEvent(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	from common.Address,
	toAddresses []common.Address,
	amount sdkmath.Uint,
	tokenIds []*tokenizationtypes.UintRange,
	ownershipTimes []*tokenizationtypes.UintRange,
) {
	// Convert addresses to strings for event
	fromStr := sdk.AccAddress(from.Bytes()).String()
	toStrs := make([]string, len(toAddresses))
	for i, addr := range toAddresses {
		toStrs[i] = sdk.AccAddress(addr.Bytes()).String()
	}

	// Build token IDs string
	tokenIdsStr := formatUintRanges(tokenIds)
	ownershipTimesStr := formatUintRanges(ownershipTimes)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"precompile_transfer_tokens",
			sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
			sdk.NewAttribute("collection_id", collectionId.String()),
			sdk.NewAttribute("from", fromStr),
			sdk.NewAttribute("to_addresses", fmt.Sprintf("%v", toStrs)),
			sdk.NewAttribute("amount", amount.String()),
			sdk.NewAttribute("token_ids", tokenIdsStr),
			sdk.NewAttribute("ownership_times", ownershipTimesStr),
		),
	)
}

// EmitIncomingApprovalEvent emits an event for setting an incoming approval via precompile
func EmitIncomingApprovalEvent(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	from common.Address,
	approvalId string,
) {
	fromStr := sdk.AccAddress(from.Bytes()).String()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"precompile_set_incoming_approval",
			sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
			sdk.NewAttribute("collection_id", collectionId.String()),
			sdk.NewAttribute("from", fromStr),
			sdk.NewAttribute("approval_id", approvalId),
		),
	)
}

// EmitOutgoingApprovalEvent emits an event for setting an outgoing approval via precompile
func EmitOutgoingApprovalEvent(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	from common.Address,
	approvalId string,
) {
	fromStr := sdk.AccAddress(from.Bytes()).String()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"precompile_set_outgoing_approval",
			sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
			sdk.NewAttribute("collection_id", collectionId.String()),
			sdk.NewAttribute("from", fromStr),
			sdk.NewAttribute("approval_id", approvalId),
		),
	)
}

// EmitGetBalanceAmountEventSingle emits an event for a getBalanceAmount query via precompile (single ID/time)
func EmitGetBalanceAmountEventSingle(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	userAddress string,
	tokenId sdkmath.Uint,
	ownershipTime sdkmath.Uint,
	totalAmount sdkmath.Uint,
) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"precompile_get_balance_amount",
			sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
			sdk.NewAttribute("collection_id", collectionId.String()),
			sdk.NewAttribute("user_address", userAddress),
			sdk.NewAttribute("token_id", tokenId.String()),
			sdk.NewAttribute("ownership_time", ownershipTime.String()),
			sdk.NewAttribute("total_amount", totalAmount.String()),
		),
	)
}

// EmitGetTotalSupplyEventSingle emits an event for a getTotalSupply query via precompile (single ID/time)
func EmitGetTotalSupplyEventSingle(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	tokenId sdkmath.Uint,
	ownershipTime sdkmath.Uint,
	totalAmount sdkmath.Uint,
) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"precompile_get_total_supply",
			sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
			sdk.NewAttribute("collection_id", collectionId.String()),
			sdk.NewAttribute("token_id", tokenId.String()),
			sdk.NewAttribute("ownership_time", ownershipTime.String()),
			sdk.NewAttribute("total_amount", totalAmount.String()),
		),
	)
}

// formatUintRanges formats a slice of UintRange into a string representation
func formatUintRanges(ranges []*tokenizationtypes.UintRange) string {
	if len(ranges) == 0 {
		return "[]"
	}
	result := "["
	for i, r := range ranges {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%s-%s", r.Start.String(), r.End.String())
	}
	result += "]"
	return result
}

