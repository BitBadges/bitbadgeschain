package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UnwrapIBCDenom(goCtx context.Context, msg *types.MsgUnwrapIBCDenom) (*types.MsgUnwrapIBCDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Validate that the denom starts with "ibc/"
	if !strings.HasPrefix(msg.Amount.Denom, "ibc/") {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "denom must start with 'ibc/', got: %s", msg.Amount.Denom)
	}

	// Calculate the balances to mint (path.Balances * amount)
	amountToMint := sdkmath.NewUintFromString(msg.Amount.Amount.String())
	if amountToMint.IsZero() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "amount cannot be zero")
	}

	// Get the collection
	collection, exists := k.GetCollectionFromStore(ctx, msg.CollectionId)
	if !exists {
		return nil, ErrCollectionNotExists
	}

	// Extract the SHA256 hash from the IBC denom
	ibcHash := strings.TrimPrefix(msg.Amount.Denom, "ibc/")
	if len(ibcHash) != 64 { // SHA256 hex string length
		return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "invalid IBC denom format, expected ibc/SHA256, got: %s", msg.Amount.Denom)
	}

	creatorAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAddress, "invalid creator address: %s", err)
	}

	// Iterate over all unwrap paths for the collection
	for _, path := range collection.IbcUnwrapPaths {
		// Reconstruct the core denomination from the path details
		// Format: "portId/channelId/badges:{sourceCollectionId}:{denom}:{address}"
		// where address is only included if withAddress is true

		// Handle denom with potential {id} placeholder replacement
		denom := path.Denom
		if path.AllowOverrideWithAnyValidToken && msg.OverrideTokenId != "" {
			// Replace {id} placeholder with the override token ID
			denom = strings.ReplaceAll(denom, "{id}", msg.OverrideTokenId)
		}

		coreDenom := fmt.Sprintf("%s/%s/badges:%s:%s", path.PortId, path.ChannelId, path.SourceCollectionId.String(), denom)

		// Add suffixes if denomSuffixDetails specifies them
		if path.DenomSuffixDetails != nil {
			if path.DenomSuffixDetails.WithAddress {
				coreDenom += fmt.Sprintf(":%s", creatorAddr.String())
			}

			// Add destination collection ID if specified
			if path.DenomSuffixDetails.DestinationCollectionId != "" {
				coreDenom += fmt.Sprintf(":%s", path.DenomSuffixDetails.DestinationCollectionId)
			}

			// Add destination chain ID if specified
			if path.DenomSuffixDetails.DestinationChainId != "" {
				coreDenom += fmt.Sprintf(":%s", path.DenomSuffixDetails.DestinationChainId)
			}
		}

		// Calculate SHA256 hash of the core denom
		hash := sha256.Sum256([]byte(coreDenom))
		expectedHash := hex.EncodeToString(hash[:])

		// Check if this matches the IBC denom hash
		if expectedHash == ibcHash {
			// Validate destination collection ID and chain ID if specified
			if path.DenomSuffixDetails != nil {
				// Validate destination collection ID
				if path.DenomSuffixDetails.DestinationCollectionId != "" {
					if path.DenomSuffixDetails.DestinationCollectionId != msg.CollectionId.String() {
						return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "destination collection ID mismatch: expected %s, got %s", path.DenomSuffixDetails.DestinationCollectionId, msg.CollectionId.String())
					}
				}

				// Validate destination chain ID
				if path.DenomSuffixDetails.DestinationChainId != "" {
					// Get current chain ID from context
					currentChainId := ctx.ChainID()
					if path.DenomSuffixDetails.DestinationChainId != currentChainId {
						return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "destination chain ID mismatch: expected %s, got %s", path.DenomSuffixDetails.DestinationChainId, currentChainId)
					}
				}
			}

			// Create the balances to mint
			var balancesToMint []*types.Balance
			for _, balance := range path.Balances {
				// Multiply each balance amount by the unwrap amount
				newAmount := balance.Amount.Mul(amountToMint)
				balancesToMint = append(balancesToMint, &types.Balance{
					Amount:         newAmount,
					OwnershipTimes: balance.OwnershipTimes,
					BadgeIds:       balance.BadgeIds,
				})
			}

			// Get current user balance or create default
			balanceKey := ConstructBalanceKey(msg.Creator, msg.CollectionId)
			userBalance, _ := k.GetBalanceOrApplyDefault(ctx, collection, msg.Creator)

			// Add the new balances to the user's existing balances
			userBalance.Balances, err = types.AddBalances(ctx, userBalance.Balances, balancesToMint)
			if err != nil {
				return nil, sdkerrors.Wrapf(err, "failed to add balances to user balance")
			}

			// Store the updated balance
			err = k.SetUserBalanceInStore(ctx, balanceKey, userBalance)
			if err != nil {
				return nil, sdkerrors.Wrapf(err, "failed to update user balance")
			}

			// Burn the IBC tokens from the creator
			err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, sdk.NewCoins(*msg.Amount))
			if err != nil {
				return nil, sdkerrors.Wrapf(err, "failed to burn IBC tokens")
			}

			err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(*msg.Amount))
			if err != nil {
				return nil, sdkerrors.Wrapf(err, "failed to burn IBC tokens")
			}

			// Emit events
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
					sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
					sdk.NewAttribute("msg_type", "unwrap_ibc_denom"),
					sdk.NewAttribute("collection_id", msg.CollectionId.String()),
					sdk.NewAttribute("ibc_denom", msg.Amount.Denom),
					sdk.NewAttribute("amount", msg.Amount.Amount.String()),
					sdk.NewAttribute("channel_id", path.ChannelId),
					sdk.NewAttribute("port_id", path.PortId),
					sdk.NewAttribute("source_collection_id", path.SourceCollectionId.String()),
					sdk.NewAttribute("override_token_id", msg.OverrideTokenId),
					sdk.NewAttribute("allow_override", fmt.Sprintf("%t", path.AllowOverrideWithAnyValidToken)),
				),
			)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent("indexer",
					sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
					sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
					sdk.NewAttribute("msg_type", "unwrap_ibc_denom"),
					sdk.NewAttribute("collection_id", msg.CollectionId.String()),
					sdk.NewAttribute("ibc_denom", msg.Amount.Denom),
					sdk.NewAttribute("amount", msg.Amount.Amount.String()),
					sdk.NewAttribute("channel_id", path.ChannelId),
					sdk.NewAttribute("port_id", path.PortId),
					sdk.NewAttribute("source_collection_id", path.SourceCollectionId.String()),
					sdk.NewAttribute("override_token_id", msg.OverrideTokenId),
					sdk.NewAttribute("allow_override", fmt.Sprintf("%t", path.AllowOverrideWithAnyValidToken)),
				),
			)

			return &types.MsgUnwrapIBCDenomResponse{}, nil
		}
	}

	// No matching path found
	return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "no matching IBC unwrap path found for denom: %s", msg.Amount.Denom)
}
