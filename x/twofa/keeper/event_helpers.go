package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EmitMessageAndIndexerEvents emits both a message event and an indexer event with the same attributes
// This is commonly used for indexing purposes where both events are needed with identical data
func EmitMessageAndIndexerEvents(ctx sdk.Context, attributes ...sdk.Attribute) {
	// Emit the standard message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage, attributes...),
	)

	// Emit the indexer event with the same attributes
	ctx.EventManager().EmitEvent(
		sdk.NewEvent("indexer", attributes...),
	)
}

