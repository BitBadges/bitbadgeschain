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
		UpdateCollectionPermissions: true,
		CollectionPermissions: msg.CollectionPermissions,
		UpdateManagerTimeline: true,
		ManagerTimeline: msg.ManagerTimeline,
		UpdateCollectionMetadataTimeline: true,
		CollectionMetadataTimeline: msg.CollectionMetadataTimeline,
		UpdateBadgeMetadataTimeline: true,
		BadgeMetadataTimeline: msg.BadgeMetadataTimeline,
		UpdateOffChainBalancesMetadataTimeline: true,
		OffChainBalancesMetadataTimeline: msg.OffChainBalancesMetadataTimeline,
		UpdateCustomDataTimeline: true,
		CustomDataTimeline: msg.CustomDataTimeline,
		UpdateCollectionApprovals: true,
		CollectionApprovals: msg.CollectionApprovals,
		UpdateStandardsTimeline: true,
		StandardsTimeline: msg.StandardsTimeline,
		UpdateIsArchivedTimeline: true,
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
