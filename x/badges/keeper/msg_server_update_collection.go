package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateCollection(goCtx context.Context, msg *types.MsgUpdateCollection) (*types.MsgUpdateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	newMsg := types.MsgUniversalUpdateCollection{
		Creator:                                msg.Creator,
		CollectionId:                           msg.CollectionId,
		BadgeIdsToAdd:                          msg.BadgeIdsToAdd,
		UpdateCollectionPermissions:            msg.UpdateCollectionPermissions,
		CollectionPermissions:                  msg.CollectionPermissions,
		UpdateManagerTimeline:                  msg.UpdateManagerTimeline,
		ManagerTimeline:                        msg.ManagerTimeline,
		UpdateCollectionMetadataTimeline:       msg.UpdateCollectionMetadataTimeline,
		CollectionMetadataTimeline:             msg.CollectionMetadataTimeline,
		UpdateBadgeMetadataTimeline:            msg.UpdateBadgeMetadataTimeline,
		BadgeMetadataTimeline:                  msg.BadgeMetadataTimeline,
		UpdateOffChainBalancesMetadataTimeline: msg.UpdateOffChainBalancesMetadataTimeline,
		OffChainBalancesMetadataTimeline:       msg.OffChainBalancesMetadataTimeline,
		UpdateCustomDataTimeline:               msg.UpdateCustomDataTimeline,
		CustomDataTimeline:                     msg.CustomDataTimeline,
		UpdateCollectionApprovals:              msg.UpdateCollectionApprovals,
		CollectionApprovals:                    msg.CollectionApprovals,
		UpdateStandardsTimeline:                msg.UpdateStandardsTimeline,
		StandardsTimeline:                      msg.StandardsTimeline,
		UpdateIsArchivedTimeline:               msg.UpdateIsArchivedTimeline,
		IsArchivedTimeline:                     msg.IsArchivedTimeline,
	}
	res, err := k.UniversalUpdateCollection(ctx, &newMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateCollectionResponse{
		CollectionId: res.CollectionId,
	}, nil
}
