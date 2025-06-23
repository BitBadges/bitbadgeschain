package keeper

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Legacy function that is all-inclusive (creates and updates)
func (k msgServer) UniversalUpdateCollection(goCtx context.Context, msg *types.MsgUniversalUpdateCollection) (*types.MsgUniversalUpdateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.CheckAndCleanMsg(ctx, true)
	if err != nil {
		return nil, err
	}

	collection := &types.BadgeCollection{}
	if msg.CollectionId.Equal(sdkmath.NewUint(0)) {
		//Creation case
		nextCollectionId := k.GetNextCollectionId(ctx)
		k.IncrementNextCollectionId(ctx)

		// From cosmos SDK x/group moduleAdd commentMore actions
		// Generate account address of collection
		var accountAddr sdk.AccAddress
		// loop here in the rare case where a ADR-028-derived address creates a
		// collision with an existing address.
		for {
			derivationKey := make([]byte, 8)
			binary.BigEndian.PutUint64(derivationKey, nextCollectionId.Uint64())

			ac, err := authtypes.NewModuleCredential(types.ModuleName, AccountGenerationPrefix, derivationKey)
			if err != nil {
				return nil, err
			}
			//generate the address from the credential
			accountAddr = sdk.AccAddress(ac.Address())

			break
		}

		collection = &types.BadgeCollection{
			CollectionId:          nextCollectionId,
			CollectionPermissions: &types.CollectionPermissions{},
			BalancesType:          msg.BalancesType,
			DefaultBalances:       msg.DefaultBalances,
			CreatedBy:             msg.Creator,
			ManagerTimeline: []*types.ManagerTimeline{
				{
					Manager: msg.Creator,
					TimelineTimes: []*types.UintRange{
						{
							Start: sdkmath.NewUint(1),
							End:   sdkmath.NewUint(math.MaxUint64),
						},
					},
				},
			},
			MintEscrowAddress: accountAddr.String(),
		}
	} else {
		//Update case
		found := false
		collection, found = k.GetCollectionFromStore(ctx, msg.CollectionId)
		if !found {
			return nil, ErrCollectionNotExists
		}
	}

	//Check must be manager
	err = k.UniversalValidate(ctx, collection, UniversalValidationParams{
		Creator:       msg.Creator,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

	//Other cases:
	//previouslyArchived && not stillArchived - we have just unarchived the collection
	//not previouslyArchived && stillArchived - we have just archived the collection (all TXs moving forward will fail, but we allow this one)
	//not previouslyArchived && not stillArchived - unarchived before and now so we allow
	previouslyArchived := types.GetIsArchived(ctx, collection)
	if msg.UpdateIsArchivedTimeline {
		if err := k.ValidateIsArchivedUpdate(ctx, collection.IsArchivedTimeline, msg.IsArchivedTimeline, collection.CollectionPermissions.CanArchiveCollection); err != nil {
			return nil, err
		}
		collection.IsArchivedTimeline = msg.IsArchivedTimeline
	}
	stillArchived := types.GetIsArchived(ctx, collection)

	if previouslyArchived && stillArchived {
		return nil, ErrCollectionIsArchived
	}

	if msg.UpdateCollectionApprovals {
		if err := k.ValidateCollectionApprovalsUpdate(ctx, collection, collection.CollectionApprovals, msg.CollectionApprovals, collection.CollectionPermissions.CanUpdateCollectionApprovals); err != nil {
			return nil, err
		}

		// Create a map of existing approvals for quick lookup
		existingApprovals := make(map[string]*types.CollectionApproval)
		for _, approval := range collection.CollectionApprovals {
			existingApprovals[approval.ApprovalId] = approval
		}

		// Only increment versions for approvals that have changed
		newApprovalsWithVersion := []*types.CollectionApproval{}
		for _, newApproval := range msg.CollectionApprovals {
			existingApproval, exists := existingApprovals[newApproval.ApprovalId]
			// Only increment version if approval is new or changed
			if !exists || existingApproval.String() != newApproval.String() {
				newVersion := k.IncrementApprovalVersion(ctx, collection.CollectionId, "collection", "", newApproval.ApprovalId)
				newApproval.Version = newVersion
			} else {
				// Keep existing version if approval hasn't changed
				newApproval.Version = existingApproval.Version
			}
			newApprovalsWithVersion = append(newApprovalsWithVersion, newApproval)
		}
		collection.CollectionApprovals = newApprovalsWithVersion
	}

	if msg.UpdateCollectionMetadataTimeline {
		if err := k.ValidateCollectionMetadataUpdate(ctx, collection.CollectionMetadataTimeline, msg.CollectionMetadataTimeline, collection.CollectionPermissions.CanUpdateCollectionMetadata); err != nil {
			return nil, err
		}
		collection.CollectionMetadataTimeline = msg.CollectionMetadataTimeline
	}

	if msg.UpdateOffChainBalancesMetadataTimeline {
		if err := k.ValidateOffChainBalancesMetadataUpdate(ctx, collection, collection.OffChainBalancesMetadataTimeline, msg.OffChainBalancesMetadataTimeline, collection.CollectionPermissions.CanUpdateOffChainBalancesMetadata); err != nil {
			return nil, err
		}
		collection.OffChainBalancesMetadataTimeline = msg.OffChainBalancesMetadataTimeline
	}

	if msg.UpdateBadgeMetadataTimeline {
		if err := k.ValidateBadgeMetadataUpdate(ctx, collection.BadgeMetadataTimeline, msg.BadgeMetadataTimeline, collection.CollectionPermissions.CanUpdateBadgeMetadata); err != nil {
			return nil, err
		}
		collection.BadgeMetadataTimeline = msg.BadgeMetadataTimeline
	}

	if msg.UpdateManagerTimeline {
		if err := k.ValidateManagerUpdate(ctx, collection.ManagerTimeline, msg.ManagerTimeline, collection.CollectionPermissions.CanUpdateManager); err != nil {
			return nil, err
		}
		collection.ManagerTimeline = msg.ManagerTimeline
	}

	if msg.UpdateStandardsTimeline {
		if err := k.ValidateStandardsUpdate(ctx, collection.StandardsTimeline, msg.StandardsTimeline, collection.CollectionPermissions.CanUpdateStandards); err != nil {
			return nil, err
		}
		collection.StandardsTimeline = msg.StandardsTimeline
	}

	if msg.UpdateCustomDataTimeline {
		if err := k.ValidateCustomDataUpdate(ctx, collection.CustomDataTimeline, msg.CustomDataTimeline, collection.CollectionPermissions.CanUpdateCustomData); err != nil {
			return nil, err
		}
		collection.CustomDataTimeline = msg.CustomDataTimeline
	}

	if msg.UpdateValidBadgeIds {
		collection, err = k.CreateBadges(ctx, collection, msg.ValidBadgeIds)
		if err != nil {
			return nil, err
		}
	}

	if msg.UpdateCollectionPermissions {
		err = k.ValidatePermissionsUpdate(ctx, collection.CollectionPermissions, msg.CollectionPermissions)
		if err != nil {
			return nil, err
		}

		collection.CollectionPermissions = msg.CollectionPermissions
	}

	if len(msg.MintEscrowCoinsToTransfer) > 0 {
		from := sdk.MustAccAddressFromBech32(msg.Creator)
		to := sdk.MustAccAddressFromBech32(collection.MintEscrowAddress)

		for _, coin := range msg.MintEscrowCoinsToTransfer {
			err = k.bankKeeper.SendCoins(ctx, from, to, sdk.NewCoins(*coin))
			if err != nil {
				return nil, err
			}
		}
	}

	if len(msg.CosmosCoinWrapperPathsToAdd) > 0 {
		pathsToAdd := make([]*types.CosmosCoinWrapperPath, len(msg.CosmosCoinWrapperPathsToAdd))
		for i, path := range msg.CosmosCoinWrapperPathsToAdd {
			var accountAddr sdk.AccAddress
			for {
				fullPath := path.Denom
				fullPathBytes := []byte(fullPath)

				ac, err := authtypes.NewModuleCredential(types.ModuleName, WrapperPathGenerationPrefix, fullPathBytes)
				if err != nil {
					return nil, err
				}
				//generate the address from the credential
				accountAddr = sdk.AccAddress(ac.Address())

				break
			}

			pathsToAdd[i] = &types.CosmosCoinWrapperPath{
				Address:        accountAddr.String(),
				Denom:          path.Denom,
				OwnershipTimes: path.OwnershipTimes,
				BadgeIds:       path.BadgeIds,
			}
		}

		collection.CosmosCoinWrapperPaths = append(collection.CosmosCoinWrapperPaths, pathsToAdd...)
	}

	// Ensure no duplicate denom paths
	denomPaths := make(map[string]bool)
	for _, path := range collection.CosmosCoinWrapperPaths {
		if _, ok := denomPaths[path.Denom]; ok {
			return nil, fmt.Errorf("duplicate ibc wrapper path denom: %s", path.Denom)
		}
		denomPaths[path.Denom] = true
	}

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
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
			sdk.NewAttribute("msg_type", "universal_update_collection"),
			sdk.NewAttribute("msg", string(msgBytes)),
			sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("indexer",
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "universal_update_collection"),
			sdk.NewAttribute("msg", string(msgBytes)),
			sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
		),
	)

	return &types.MsgUniversalUpdateCollectionResponse{
		CollectionId: collection.CollectionId,
	}, nil
}
