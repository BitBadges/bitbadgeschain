package keeper

import (
	"context"
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NewCollection(goCtx context.Context, msg *types.MsgNewCollection) (*types.MsgNewCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	NextCollectionId := k.GetNextCollectionId(ctx)
	k.IncrementNextCollectionId(ctx)

	collection := &types.BadgeCollection{
		CollectionId:       				NextCollectionId,
		CollectionMetadataTimeline: msg.CollectionMetadataTimeline,
		OffChainBalancesMetadataTimeline:   msg.OffChainBalancesMetadataTimeline,
		BadgeMetadataTimeline:      msg.BadgeMetadataTimeline,
		ManagerTimeline:            []*types.ManagerTimeline{
			{
				Manager: msg.Creator,
				TimelineTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(0),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
			},
		},
		Permissions:        msg.Permissions,
		CollectionApprovedTransfersTimeline:  msg.CollectionApprovedTransfersTimeline,
		CustomDataTimeline:         msg.CustomDataTimeline,
		ContractAddressTimeline:    msg.ContractAddressTimeline,
		StandardsTimeline:          msg.StandardsTimeline,
		NextBadgeId:        sdkmath.NewUint(1),
		BalancesType:       msg.BalancesType,
		IsArchivedTimeline: []*types.IsArchivedTimeline{
			{
				IsArchived: false,
				TimelineTimes:      []*types.UintRange{
					{
						Start: sdkmath.NewUint(0),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
			},
		},
		UnmintedSupplys:    []*types.Balance{},
		TotalSupplys:         []*types.Balance{},
		InheritedBalancesTimeline: msg.InheritedBalancesTimeline,
		DefaultUserApprovedOutgoingTransfersTimeline: msg.DefaultApprovedOutgoingTransfersTimeline,
		DefaultUserApprovedIncomingTransfersTimeline: msg.DefaultApprovedIncomingTransfersTimeline,
	}

	for _, addressMapping := range msg.AddressMappings {
		if err := k.CreateAddressMapping(ctx, addressMapping); err != nil {
			return nil, err
		}
	}

	if len(msg.BadgesToCreate) > 0 {
		err := *new(error)
		collection, err = k.CreateBadges(ctx, collection, msg.BadgesToCreate, msg.Transfers)
		if err != nil {
			return nil, err
		}
	}

	if msg.InheritedBalancesTimeline != nil && len(msg.InheritedBalancesTimeline) > 0 {
		if err := k.ValidateInheritedBalancesUpdate(ctx, collection, collection.InheritedBalancesTimeline, collection.InheritedBalancesTimeline, collection.Permissions.CanUpdateInheritedBalances); err != nil {
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
