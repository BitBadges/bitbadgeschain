package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateCollection(goCtx context.Context, msg *types.MsgCreateCollection) (*types.MsgCreateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	newMsg := types.MsgUniversalUpdateCollection{
		Creator: msg.Creator,
		CollectionId: sdk.NewUint(0), //We use 0 to indicate a new collection

		//Exclusive to collection creations
		BalancesType: msg.BalancesType,
		DefaultOutgoingApprovals: msg.DefaultOutgoingApprovals,
		DefaultIncomingApprovals: msg.DefaultIncomingApprovals,
		DefaultAutoApproveSelfInitiatedOutgoingTransfers: msg.DefaultAutoApproveSelfInitiatedOutgoingTransfers,
		DefaultAutoApproveSelfInitiatedIncomingTransfers: msg.DefaultAutoApproveSelfInitiatedIncomingTransfers,
		DefaultUserPermissions: msg.DefaultUserPermissions,

		//Applicable to creations and updates
		BadgesToCreate: msg.BadgesToCreate,
		UpdateCollectionPermissions: msg.UpdateCollectionPermissions,
		CollectionPermissions: msg.CollectionPermissions,
		UpdateManagerTimeline: msg.UpdateManagerTimeline,
		ManagerTimeline: msg.ManagerTimeline,
		UpdateCollectionMetadataTimeline: msg.UpdateCollectionMetadataTimeline,
		CollectionMetadataTimeline: msg.CollectionMetadataTimeline,
		UpdateBadgeMetadataTimeline: msg.UpdateBadgeMetadataTimeline,
		BadgeMetadataTimeline: msg.BadgeMetadataTimeline,
		UpdateOffChainBalancesMetadataTimeline: msg.UpdateOffChainBalancesMetadataTimeline,
		OffChainBalancesMetadataTimeline: msg.OffChainBalancesMetadataTimeline,
		UpdateCustomDataTimeline: msg.UpdateCustomDataTimeline,
		CustomDataTimeline: msg.CustomDataTimeline,
		UpdateCollectionApprovals: msg.UpdateCollectionApprovals,
		CollectionApprovals: msg.CollectionApprovals,
		UpdateStandardsTimeline: msg.UpdateStandardsTimeline,
		StandardsTimeline: msg.StandardsTimeline,
		UpdateIsArchivedTimeline: msg.UpdateIsArchivedTimeline,
		IsArchivedTimeline: msg.IsArchivedTimeline,
	}
	res, err := k.UniversalUpdateCollection(ctx, &newMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateCollectionResponse{
		CollectionId: res.CollectionId,
	}, nil
}
