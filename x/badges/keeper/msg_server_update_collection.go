package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateCollection(goCtx context.Context, msg *types.MsgUpdateCollection) (*types.MsgUpdateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	newMsg := types.MsgUniversalUpdateCollection{
		Creator:                                msg.Creator,
		CollectionId:                           msg.CollectionId,
		ValidBadgeIds:                          msg.ValidBadgeIds,
		UpdateValidBadgeIds:                    msg.UpdateValidBadgeIds,
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
		MintEscrowCoinsToTransfer:              msg.MintEscrowCoinsToTransfer,
		CosmosCoinWrapperPathsToAdd:            msg.CosmosCoinWrapperPathsToAdd,
	}
	res, err := k.UniversalUpdateCollection(ctx, &newMsg)
	if err != nil {
		return nil, err
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "update_collection"),
			sdk.NewAttribute("msg", string(msgBytes)),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("indexer",
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "update_collection"),
			sdk.NewAttribute("msg", string(msgBytes)),
		),
	)

	return &types.MsgUpdateCollectionResponse{
		CollectionId: res.CollectionId,
	}, nil
}
