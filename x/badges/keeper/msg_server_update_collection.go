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
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		ValidTokenIds:               msg.ValidTokenIds,
		UpdateValidTokenIds:         msg.UpdateValidTokenIds,
		UpdateCollectionPermissions: msg.UpdateCollectionPermissions,
		CollectionPermissions:       msg.CollectionPermissions,
		UpdateManager:               msg.UpdateManager,
		Manager:                     msg.Manager,
		UpdateCollectionMetadata:    msg.UpdateCollectionMetadata,
		CollectionMetadata:          msg.CollectionMetadata,
		UpdateTokenMetadata:         msg.UpdateTokenMetadata,
		TokenMetadata:               msg.TokenMetadata,
		UpdateCustomData:            msg.UpdateCustomData,
		CustomData:                  msg.CustomData,
		UpdateCollectionApprovals:   msg.UpdateCollectionApprovals,
		CollectionApprovals:         msg.CollectionApprovals,
		UpdateStandards:             msg.UpdateStandards,
		Standards:                   msg.Standards,
		UpdateIsArchived:            msg.UpdateIsArchived,
		IsArchived:                  msg.IsArchived,
		MintEscrowCoinsToTransfer:   msg.MintEscrowCoinsToTransfer,
		CosmosCoinWrapperPathsToAdd: msg.CosmosCoinWrapperPathsToAdd,
		Invariants:                  msg.Invariants,
		AliasPathsToAdd:             msg.AliasPathsToAdd,
	}
	res, err := k.UniversalUpdateCollection(ctx, &newMsg)
	if err != nil {
		return nil, err
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "update_collection"),
		sdk.NewAttribute("msg", string(msgBytes)),
	)

	return &types.MsgUpdateCollectionResponse{
		CollectionId: res.CollectionId,
	}, nil
}
