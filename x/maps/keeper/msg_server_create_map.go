package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/maps/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	badgetypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func (k msgServer) CreateMap(goCtx context.Context, msg *types.MsgCreateMap) (*types.MsgCreateMapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	//Map IDs cannot have a "-" in them (reserved for future use cases)
	for _, char := range msg.MapId {
		if char == '-' {
			return nil, sdkerrors.Wrap(ErrInvalidMapId, "Map ID cannot contain '-'")
		}
	}

	//Maps w/ ID as a valid address are reserved for that address
	_, err := sdk.AccAddressFromBech32(msg.MapId)
	if err == nil {
		if msg.MapId != msg.Creator {
			return nil, sdkerrors.Wrap(ErrInvalidMapId, "Map ID cannot be a valid address that is not your own")
		}
	}

	//Numeric map IDs are reserved for specific badge collections
	collId, err := sdkmath.ParseUint(msg.MapId)
	if err == nil {
		currCollectionRes, err := k.badgesKeeper.GetCollection(ctx, &badgetypes.QueryGetCollectionRequest{
			CollectionId: collId,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(ErrInvalidMapId, "Map ID must be a valid collection ID. Could not find collection in store")
		}

		//Check if user is manager of collection in x/badges
		currManager := badgetypes.GetCurrentManager(ctx, currCollectionRes.Collection)
		if currManager != msg.Creator {
			return nil, sdkerrors.Wrap(ErrInvalidMapId, "Numeric map IDs are reserved for specific badge collections. To create a map for this collection, you must be the manager of the collection")
		}
	}

	mapToAdd := types.Map{
		Creator: msg.Creator,
		MapId:   msg.MapId,
		UpdateCriteria: msg.UpdateCriteria,
		ValueOptions: msg.ValueOptions,
		DefaultValue: msg.DefaultValue,
		ManagerTimeline: msg.ManagerTimeline,
		MetadataTimeline: msg.MetadataTimeline,
		Permissions: msg.Permissions,
		InheritManagerTimelineFrom: msg.InheritManagerTimelineFrom,
	}
	if msg.Permissions == nil {
		mapToAdd.Permissions = &types.MapPermissions{}
	}
	if msg.ValueOptions == nil {
		mapToAdd.ValueOptions = &types.ValueOptions{}
	}
	if msg.UpdateCriteria == nil {
		mapToAdd.UpdateCriteria = &types.MapUpdateCriteria{}
	}


	//Check if protocol already exists
	if k.StoreHasMapID(ctx, msg.MapId) {
		return nil, sdkerrors.Wrap(ErrMapExists, "Protocol already exists")
	}

	//Add protocol to store
	err = k.SetMapInStore(ctx, &mapToAdd)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to add protocol to store")
	}

	return &types.MsgCreateMapResponse{}, nil
}
