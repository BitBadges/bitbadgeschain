package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteCollection(goCtx context.Context, msg *types.MsgDeleteCollection) (*types.MsgDeleteCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	collection, found := k.GetCollectionFromStore(ctx, msg.CollectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	err := k.UniversalValidate(ctx, collection, UniversalValidationParams{
		Creator:       msg.Creator,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

	// Check deleted permission is valid for current time
	err = k.CheckIfActionPermissionPermits(ctx, collection.CollectionPermissions.CanDeleteCollection, "can delete collection")
	if err != nil {
		return nil, err
	}

	// Escrowed supply stays redeemable only through this collection, so it must be
	// fully unwound before the collection can go away.
	for _, path := range collection.CosmosCoinWrapperPaths {
		escrow, _, err := k.GetBalanceOrApplyDefault(ctx, collection, path.Address)
		if err != nil {
			return nil, err
		}
		if hasNonZeroBalance(escrow) {
			return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "wrapper path %s still escrows tokens for outstanding %s coins", path.Denom, WrappedDenomPrefix+collection.CollectionId.String()+":"+path.Denom)
		}
	}
	if collection.Invariants != nil && collection.Invariants.CosmosCoinBackedPath != nil {
		stats, _ := k.GetCollectionStatsFromStore(ctx, collection.CollectionId)
		for _, bal := range stats.Balances {
			if bal.Amount.GT(sdkmath.ZeroUint()) {
				return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "backed tokens are still in circulation; they must be backed before the collection can be deleted")
			}
		}
	}

	// Purge all collection-related state before deleting the collection
	// This includes balances, approval trackers, challenge trackers, approval versions, and ETH signature trackers
	if err := k.PurgeCollectionState(ctx, collection.CollectionId); err != nil {
		return nil, err
	}

	k.DeleteCollectionFromStore(ctx, collection.CollectionId)

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "tokenization"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "delete_collection"),
		sdk.NewAttribute("msg", msgStr),
	)
	return &types.MsgDeleteCollectionResponse{}, nil
}
