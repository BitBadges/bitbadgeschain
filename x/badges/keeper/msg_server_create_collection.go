package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v13"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastOldCreateCollectionToNewType(oldMsg *oldtypes.MsgCreateCollection) (*types.MsgCreateCollection, error) {
	// Convert to JSON
	jsonBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return nil, err
	}

	// Unmarshal into new type
	var newCollection types.MsgCreateCollection
	if err := json.Unmarshal(jsonBytes, &newCollection); err != nil {
		return nil, err
	}

	return &newCollection, nil
}

func (k msgServer) CreateCollectionV13(goCtx context.Context, msg *oldtypes.MsgCreateCollection) (*types.MsgCreateCollectionResponse, error) {
	newMsg, err := CastOldCreateCollectionToNewType(msg)
	if err != nil {
		return nil, err
	}
	return k.CreateCollection(goCtx, newMsg)
}

func (k msgServer) CreateCollectionV14(goCtx context.Context, msg *types.MsgCreateCollection) (*types.MsgCreateCollectionResponse, error) {
	return k.CreateCollection(goCtx, msg)
}

func (k msgServer) CreateCollection(goCtx context.Context, msg *types.MsgCreateCollection) (*types.MsgCreateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	newMsg := types.MsgUniversalUpdateCollection{
		Creator:      msg.Creator,
		CollectionId: sdkmath.NewUint(0), //We use 0 to indicate a new collection

		//Exclusive to collection creations
		BalancesType:    msg.BalancesType,
		DefaultBalances: msg.DefaultBalances,

		//Applicable to creations and updates
		ValidBadgeIds:                          msg.ValidBadgeIds,
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

		MintEscrowCoinsToTransfer:   msg.MintEscrowCoinsToTransfer,
		CosmosCoinWrapperPathsToAdd: msg.CosmosCoinWrapperPathsToAdd,
		Invariants:                  msg.Invariants,
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
			sdk.NewAttribute("msg_type", "create_collection"),
			sdk.NewAttribute("msg", string(msgBytes)),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("indexer",
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "create_collection"),
			sdk.NewAttribute("msg", string(msgBytes)),
		),
	)

	return &types.MsgCreateCollectionResponse{
		CollectionId: res.CollectionId,
	}, nil
}
