package keeper

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ProcessOwnershipQuery handles an incoming ICQ ownership query and returns a response.
// This is the core method that verifies token ownership for cross-chain queries.
func (k Keeper) ProcessOwnershipQuery(
	ctx sdk.Context,
	query *types.OwnershipQueryPacket,
) *types.OwnershipQueryResponsePacket {
	// Validate the query
	if err := query.ValidateBasic(); err != nil {
		return types.NewErrorOwnershipQueryResponsePacket(query.QueryId, err)
	}

	// Parse collection ID
	collectionId, err := sdkmath.ParseUint(query.CollectionId)
	if err != nil {
		return types.NewErrorOwnershipQueryResponsePacket(query.QueryId,
			sdkerrors.Wrapf(types.ErrInvalidRequest, "invalid collection_id: %s", err.Error()))
	}

	// Verify collection exists
	_, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return types.NewErrorOwnershipQueryResponsePacket(query.QueryId,
			sdkerrors.Wrapf(ErrCollectionNotExists, "collection %s does not exist", collectionId.String()))
	}

	// Normalize address (handle both bech32 and EVM hex formats)
	address := normalizeAddressForICQ(query.Address)

	// Construct balance key and get user balance
	balanceKey := ConstructBalanceKey(address, collectionId)
	balance, found := k.GetUserBalanceFromStore(ctx, balanceKey)
	if !found || len(balance.Balances) == 0 {
		// User has no balance for this collection
		return types.NewOwnershipQueryResponsePacket(
			query.QueryId,
			false,
			sdkmath.ZeroUint(),
			nil,
			uint64(ctx.BlockHeight()),
			"",
		)
	}

	// Parse token ID and ownership time
	tokenId, err := sdkmath.ParseUint(query.TokenId)
	if err != nil {
		return types.NewErrorOwnershipQueryResponsePacket(query.QueryId,
			sdkerrors.Wrapf(types.ErrInvalidRequest, "invalid token_id: %s", err.Error()))
	}
	ownershipTime, err := sdkmath.ParseUint(query.OwnershipTime)
	if err != nil {
		return types.NewErrorOwnershipQueryResponsePacket(query.QueryId,
			sdkerrors.Wrapf(types.ErrInvalidRequest, "invalid ownership_time: %s", err.Error()))
	}

	// Get the balance amount for this specific token ID and ownership time
	totalAmount := getBalanceForIdAndTime(balance.Balances, tokenId, ownershipTime)

	// User owns tokens if total_amount > 0
	ownsTokens := totalAmount.GT(sdkmath.ZeroUint())

	return types.NewOwnershipQueryResponsePacket(
		query.QueryId,
		ownsTokens,
		totalAmount,
		nil, // balance_proof - not implemented yet
		uint64(ctx.BlockHeight()),
		"",
	)
}

// ProcessBulkOwnershipQuery handles a bulk ICQ ownership query.
// It processes multiple queries in a single IBC packet for efficiency.
func (k Keeper) ProcessBulkOwnershipQuery(
	ctx sdk.Context,
	bulk *types.BulkOwnershipQueryPacket,
) *types.BulkOwnershipQueryResponsePacket {
	// Validate the bulk query
	if err := bulk.ValidateBasic(); err != nil {
		// Return a single error response for the whole bulk
		return &types.BulkOwnershipQueryResponsePacket{
			QueryId: bulk.QueryId,
			Responses: []*types.OwnershipQueryResponsePacket{
				types.NewErrorOwnershipQueryResponsePacket(bulk.QueryId, err),
			},
		}
	}

	// Process each individual query
	responses := make([]*types.OwnershipQueryResponsePacket, len(bulk.Queries))
	for i, query := range bulk.Queries {
		responses[i] = k.ProcessOwnershipQuery(ctx, query)
	}

	return types.NewBulkOwnershipQueryResponsePacket(bulk.QueryId, responses)
}

// normalizeAddressForICQ normalizes an address for ICQ queries.
// It handles both bech32 Cosmos addresses and EVM hex addresses.
func normalizeAddressForICQ(address string) string {
	// If it starts with 0x, convert from EVM hex to bech32
	if strings.HasPrefix(address, "0x") || strings.HasPrefix(address, "0X") {
		// Remove 0x prefix and convert to lowercase
		hexAddr := strings.ToLower(strings.TrimPrefix(strings.TrimPrefix(address, "0X"), "0x"))

		// Pad to 40 characters if needed
		for len(hexAddr) < 40 {
			hexAddr = "0" + hexAddr
		}

		// Convert hex to bytes
		addrBytes := make([]byte, 20)
		for i := 0; i < 20; i++ {
			b, err := parseHexByte(hexAddr[i*2 : i*2+2])
			if err != nil {
				// If conversion fails, return original address
				return address
			}
			addrBytes[i] = b
		}

		// Convert to bech32 with bb prefix
		bech32Addr, err := sdk.Bech32ifyAddressBytes("bb", addrBytes)
		if err != nil {
			return address
		}
		return bech32Addr
	}

	// Already a bech32 address, return as-is
	return address
}

// parseHexByte parses a two-character hex string into a byte
func parseHexByte(s string) (byte, error) {
	if len(s) != 2 {
		return 0, sdkerrors.Wrap(types.ErrInvalidRequest, "invalid hex byte length")
	}

	var result byte
	for i := 0; i < 2; i++ {
		result <<= 4
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			result |= c - '0'
		case c >= 'a' && c <= 'f':
			result |= c - 'a' + 10
		case c >= 'A' && c <= 'F':
			result |= c - 'A' + 10
		default:
			return 0, sdkerrors.Wrap(types.ErrInvalidRequest, "invalid hex character")
		}
	}
	return result, nil
}

// getBalanceForIdAndTime returns the balance amount for a specific token ID and ownership time.
// This mirrors the SDK's getBalanceForIdAndTime function - it sums all matching balance amounts
// since the same (id, time) could appear in multiple overlapping Balance objects.
func getBalanceForIdAndTime(balances []*types.Balance, tokenId sdkmath.Uint, ownershipTime sdkmath.Uint) sdkmath.Uint {
	amount := sdkmath.ZeroUint()

	for _, balance := range balances {
		if balance == nil {
			continue
		}

		// Check if tokenId is in this balance's token ID ranges
		foundTokenId := false
		for _, tokenRange := range balance.TokenIds {
			if tokenRange != nil && tokenId.GTE(tokenRange.Start) && tokenId.LTE(tokenRange.End) {
				foundTokenId = true
				break
			}
		}

		if !foundTokenId {
			continue
		}

		// Check if ownershipTime is in this balance's ownership time ranges
		foundTime := false
		for _, timeRange := range balance.OwnershipTimes {
			if timeRange != nil && ownershipTime.GTE(timeRange.Start) && ownershipTime.LTE(timeRange.End) {
				foundTime = true
				break
			}
		}

		if foundTime {
			amount = amount.Add(balance.Amount)
		}
	}

	return amount
}

// EmitICQRequestEvent emits an event for an incoming ICQ request
func (k Keeper) EmitICQRequestEvent(ctx sdk.Context, query *types.OwnershipQueryPacket) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeICQRequest,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyQueryId, query.QueryId),
			sdk.NewAttribute(types.AttributeKeyAddress, query.Address),
			sdk.NewAttribute(types.AttributeKeyCollectionId, query.CollectionId),
		),
	)
}

// EmitICQResponseEvent emits an event for an ICQ response
func (k Keeper) EmitICQResponseEvent(ctx sdk.Context, response *types.OwnershipQueryResponsePacket) {
	attrs := []sdk.Attribute{
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyQueryId, response.QueryId),
		sdk.NewAttribute(types.AttributeKeyOwnsTokens, boolToString(response.OwnsTokens)),
		sdk.NewAttribute(types.AttributeKeyTotalAmount, response.TotalAmount.String()),
	}

	if response.Error != "" {
		attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyError, response.Error))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeICQResponse, attrs...),
	)
}

// boolToString converts a boolean to a string
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// ProcessFullBalanceQuery handles an incoming ICQ full balance query and returns a response.
// This returns the complete UserBalanceStore for a user including balances, approvals, and permissions.
func (k Keeper) ProcessFullBalanceQuery(
	ctx sdk.Context,
	query *types.FullBalanceQueryPacket,
) *types.FullBalanceQueryResponsePacket {
	// Validate the query
	if err := query.ValidateBasic(); err != nil {
		return types.NewErrorFullBalanceQueryResponsePacket(query.QueryId, err)
	}

	// Parse collection ID
	collectionId, err := sdkmath.ParseUint(query.CollectionId)
	if err != nil {
		return types.NewErrorFullBalanceQueryResponsePacket(query.QueryId,
			sdkerrors.Wrapf(types.ErrInvalidRequest, "invalid collection_id: %s", err.Error()))
	}

	// Verify collection exists
	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return types.NewErrorFullBalanceQueryResponsePacket(query.QueryId,
			sdkerrors.Wrapf(ErrCollectionNotExists, "collection %s does not exist", collectionId.String()))
	}

	// Normalize address (handle both bech32 and EVM hex formats)
	address := normalizeAddressForICQ(query.Address)

	// Get user balance store with defaults applied
	userBalanceStore, _, err := k.GetBalanceOrApplyDefault(ctx, collection, address)
	if err != nil {
		return types.NewErrorFullBalanceQueryResponsePacket(query.QueryId,
			sdkerrors.Wrapf(types.ErrInvalidRequest, "failed to get balance: %s", err.Error()))
	}

	// Serialize the balance store
	balanceStoreBytes, err := types.ModuleCdc.Marshal(userBalanceStore)
	if err != nil {
		return types.NewErrorFullBalanceQueryResponsePacket(query.QueryId,
			sdkerrors.Wrapf(types.ErrInvalidRequest, "failed to serialize balance store: %s", err.Error()))
	}

	return types.NewFullBalanceQueryResponsePacket(
		query.QueryId,
		balanceStoreBytes,
		uint64(ctx.BlockHeight()),
		"",
	)
}

// EmitFullBalanceQueryRequestEvent emits an event for an incoming full balance query request
func (k Keeper) EmitFullBalanceQueryRequestEvent(ctx sdk.Context, query *types.FullBalanceQueryPacket) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeICQRequest,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyQueryId, query.QueryId),
			sdk.NewAttribute(types.AttributeKeyAddress, query.Address),
			sdk.NewAttribute(types.AttributeKeyCollectionId, query.CollectionId),
			sdk.NewAttribute("query_type", "full_balance"),
		),
	)
}

// EmitFullBalanceQueryResponseEvent emits an event for a full balance query response
func (k Keeper) EmitFullBalanceQueryResponseEvent(ctx sdk.Context, response *types.FullBalanceQueryResponsePacket) {
	attrs := []sdk.Attribute{
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyQueryId, response.QueryId),
		sdk.NewAttribute("query_type", "full_balance"),
	}

	if response.Error != "" {
		attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyError, response.Error))
	} else {
		attrs = append(attrs, sdk.NewAttribute("balance_store_size", string(rune(len(response.BalanceStore)))))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeICQResponse, attrs...),
	)
}
