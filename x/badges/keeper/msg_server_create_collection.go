package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// NewCollectionId represents the ID used to indicate a new collection creation
	NewCollectionId = 0
)

func (k msgServer) CreateCollection(goCtx context.Context, msg *types.MsgCreateCollection) (*types.MsgCreateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message before processing
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	newMsg := types.MsgUniversalUpdateCollection{
		Creator:      msg.Creator,
		CollectionId: sdkmath.NewUint(NewCollectionId), //We use 0 to indicate a new collection

		//Exclusive to collection creations
		DefaultBalances: msg.DefaultBalances,

		//Applicable to creations and updates
		ValidTokenIds:                    msg.ValidTokenIds,
		UpdateCollectionPermissions:      true,
		CollectionPermissions:            msg.CollectionPermissions,
		UpdateManagerTimeline:            true,
		ManagerTimeline:                  msg.ManagerTimeline,
		UpdateCollectionMetadataTimeline: true,
		CollectionMetadataTimeline:       msg.CollectionMetadataTimeline,
		UpdateTokenMetadataTimeline:      true,
		TokenMetadataTimeline:            msg.TokenMetadataTimeline,
		UpdateCustomDataTimeline:         true,
		CustomDataTimeline:               msg.CustomDataTimeline,
		UpdateCollectionApprovals:        true,
		CollectionApprovals:              msg.CollectionApprovals,
		UpdateStandardsTimeline:          true,
		StandardsTimeline:                msg.StandardsTimeline,
		UpdateIsArchivedTimeline:         true,
		IsArchivedTimeline:               msg.IsArchivedTimeline,

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

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "create_collection"),
		sdk.NewAttribute("msg", string(msgBytes)),
	)

	return &types.MsgCreateCollectionResponse{
		CollectionId: res.CollectionId,
	}, nil
}
