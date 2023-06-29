package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferBadge(goCtx context.Context, msg *types.MsgTransferBadge) (*types.MsgTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                 msg.Creator,
		CollectionId:            msg.CollectionId,
	})
	if err != nil {
		return nil, err
	}

	if collection.BalancesType != sdkmath.NewUint(0) {
		return nil, ErrOffChainBalances
	}

	
	if err := k.Keeper.HandleTransfers(ctx, collection, msg.Transfers, "Manager", msg.OnlyDeductApprovals); err != nil {
		return nil, err
	}

	//The "mint" balances are stored via the collection's unminted supplys
	for _, transfer := range msg.Transfers {
		if transfer.From == "Mint" {
			if err := k.SetCollectionInStore(ctx, collection); err != nil {
				return nil, err
			}
			break;
		}
	}

	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute("collection_id", fmt.Sprint(msg.CollectionId)),
		),
	)

	return &types.MsgTransferBadgeResponse{}, nil
}
