package keeper

import (
	"encoding/json"
	"fmt"

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

// MarshalMessageForEvent marshals a message to JSON for event emission.
// Note: We do NOT truncate or limit message size, as this would break indexers that need to parse
// complete JSON. Cosmos SDK has its own transaction/event size limits at the protocol level that
// prevent DoS attacks. The data in messages is public blockchain data anyway, so there's no
// security concern with including it in events.
func MarshalMessageForEvent(msg interface{}) (string, error) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal message: %w", err)
	}

	return string(msgBytes), nil
}










