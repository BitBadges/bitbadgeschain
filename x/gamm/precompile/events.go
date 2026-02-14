package gamm

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"

	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

// EmitJoinPoolEvent emits an event for a join pool operation via precompile
func EmitJoinPoolEvent(
	ctx sdk.Context,
	poolId uint64,
	sender common.Address,
	shareOutAmount sdkmath.Int,
	tokenIn sdk.Coins,
) {
	senderStr := sdk.AccAddress(sender.Bytes()).String()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"precompile_join_pool",
			sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
			sdk.NewAttribute("pool_id", fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute("sender", senderStr),
			sdk.NewAttribute("share_out_amount", shareOutAmount.String()),
			sdk.NewAttribute("token_in", tokenIn.String()),
		),
	)
}

// EmitExitPoolEvent emits an event for an exit pool operation via precompile
func EmitExitPoolEvent(
	ctx sdk.Context,
	poolId uint64,
	sender common.Address,
	tokenOut sdk.Coins,
) {
	senderStr := sdk.AccAddress(sender.Bytes()).String()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"precompile_exit_pool",
			sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
			sdk.NewAttribute("pool_id", fmt.Sprintf("%d", poolId)),
			sdk.NewAttribute("sender", senderStr),
			sdk.NewAttribute("token_out", tokenOut.String()),
		),
	)
}

// EmitSwapEvent emits an event for a swap operation via precompile
func EmitSwapEvent(
	ctx sdk.Context,
	sender common.Address,
	routes []poolmanagertypes.SwapAmountInRoute,
	tokenIn sdk.Coin,
	tokenOutAmount sdkmath.Int,
) {
	senderStr := sdk.AccAddress(sender.Bytes()).String()

	// Format routes as string
	routesStr := formatRoutes(routes)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"precompile_swap_exact_amount_in",
			sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
			sdk.NewAttribute("sender", senderStr),
			sdk.NewAttribute("routes", routesStr),
			sdk.NewAttribute("token_in", tokenIn.String()),
			sdk.NewAttribute("token_out_amount", tokenOutAmount.String()),
		),
	)
}

// EmitIBCTransferEvent emits an event for an IBC transfer operation via precompile
func EmitIBCTransferEvent(
	ctx sdk.Context,
	sender common.Address,
	sourceChannel string,
	receiver string,
	tokenOutAmount sdkmath.Int,
) {
	senderStr := sdk.AccAddress(sender.Bytes()).String()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"precompile_swap_exact_amount_in_with_ibc_transfer",
			sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
			sdk.NewAttribute("sender", senderStr),
			sdk.NewAttribute("source_channel", sourceChannel),
			sdk.NewAttribute("receiver", receiver),
			sdk.NewAttribute("token_out_amount", tokenOutAmount.String()),
		),
	)
}

// formatRoutes formats swap routes into a string representation
func formatRoutes(routes []poolmanagertypes.SwapAmountInRoute) string {
	if len(routes) == 0 {
		return "[]"
	}
	result := "["
	for i, route := range routes {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d:%s", route.PoolId, route.TokenOutDenom)
	}
	result += "]"
	return result
}

