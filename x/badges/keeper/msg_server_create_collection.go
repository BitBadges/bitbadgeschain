package keeper

import (
	"context"

	"bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateCollection(goCtx context.Context, msg *types.MsgCreateCollection) (*types.MsgCreateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	newMsg := types.MsgUniversalUpdateCollection{
		Creator:      msg.Creator,
		CollectionId: sdkmath.NewUint(0), //We use 0 to indicate a new collection

		//Exclusive to collection creations
		BalancesType:    msg.BalancesType,
		DefaultBalances: msg.DefaultBalances,

		//Applicable to creations and updates
		BadgeIdsToAdd:                          msg.BadgeIdsToAdd,
		UpdateCollectionPermissions:            true,
		CollectionPermissions:                  msg.CollectionPermissions,
		UpdateManagerTimeline:                  true,
		ManagerTimeline:                        msg.ManagerTimeline,
		UpdateCollectionMetadataTimeline:       true,
		CollectionMetadataTimeline:             msg.CollectionMetadataTimeline,
		UpdateBadgeMetadataTimeline:            true,
		BadgeMetadataTimeline:                  msg.BadgeMetadataTimeline,
		UpdateOffChainBalancesMetadataTimeline: true,
		OffChainBalancesMetadataTimeline:       msg.OffChainBalancesMetadataTimeline,
		UpdateCustomDataTimeline:               true,
		CustomDataTimeline:                     msg.CustomDataTimeline,
		UpdateCollectionApprovals:              true,
		CollectionApprovals:                    msg.CollectionApprovals,
		UpdateStandardsTimeline:                true,
		StandardsTimeline:                      msg.StandardsTimeline,
		UpdateIsArchivedTimeline:               true,
		IsArchivedTimeline:                     msg.IsArchivedTimeline,
	}

	res, err := k.UniversalUpdateCollection(ctx, &newMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateCollectionResponse{
		CollectionId: res.CollectionId,
	}, nil
}
