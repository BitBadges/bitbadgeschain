package precompile

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
)

// EmitSendEvent emits an event for a send operation via precompile
func EmitSendEvent(
	ctx sdk.Context,
	from common.Address,
	toAddress string,
	amount sdk.Coins,
) {
	fromStr := sdk.AccAddress(from.Bytes()).String()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"precompile_send",
			sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
			sdk.NewAttribute("from", fromStr),
			sdk.NewAttribute("to_address", toAddress),
			sdk.NewAttribute("amount", amount.String()),
		),
	)
}

// LogPrecompileUsage logs precompile method usage for monitoring
func LogPrecompileUsage(ctx sdk.Context, methodName string, success bool, gasUsed uint64, err error) {
	// Emit event for monitoring
	event := sdk.NewEvent(
		"precompile_usage",
		sdk.NewAttribute(sdk.AttributeKeyModule, "evm_precompile"),
		sdk.NewAttribute("precompile", "sendmanager"),
		sdk.NewAttribute("method", methodName),
		sdk.NewAttribute("success", fmt.Sprintf("%t", success)),
	)

	if err != nil {
		event = event.AppendAttributes(
			sdk.NewAttribute("error", err.Error()),
		)
	}

	if gasUsed > 0 {
		event = event.AppendAttributes(
			sdk.NewAttribute("gas_used", fmt.Sprintf("%d", gasUsed)),
		)
	}

	ctx.EventManager().EmitEvent(event)
}
