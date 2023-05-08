package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NewCollection(goCtx context.Context, msg *types.MsgNewCollection) (*types.MsgNewCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	NextCollectionId := k.GetNextCollectionId(ctx)
	k.IncrementNextCollectionId(ctx)

	collection := types.BadgeCollection{
		CollectionId:             NextCollectionId,
		CollectionUri:            msg.CollectionUri,
		BalancesUri: 							msg.BalancesUri,
		BadgeUris:                msg.BadgeUris,
		Manager:                  msg.Creator,
		Permissions:              msg.Permissions,
		AllowedTransfers:         msg.AllowedTransfers,
		ManagerApprovedTransfers: msg.ManagerApprovedTransfers,
		Bytes:                    msg.Bytes,
		Standard:                 msg.Standard,
		NextBadgeId:              sdk.NewUint(1),
		NextClaimId: 						  sdk.NewUint(1),
	}

	if len(msg.BadgeSupplys) != 0 {
		err := *new(error)
		collection, err = k.CreateBadges(ctx, collection, msg.BadgeSupplys, msg.Transfers, msg.Claims, msg.Creator)
		if err != nil {
			return nil, err
		}
	}

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("collectionId", fmt.Sprint(NextCollectionId)),
		),
	)

	return &types.MsgNewCollectionResponse{
		CollectionId: NextCollectionId,
	}, nil
}
