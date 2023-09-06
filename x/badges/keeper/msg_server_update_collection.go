package keeper

import (
	"context"
	"fmt"
	"math"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateCollection(goCtx context.Context, msg *types.MsgUpdateCollection) (*types.MsgUpdateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection := &types.BadgeCollection{}
	if msg.CollectionId.Equal(sdkmath.NewUint(0)) {
		nextCollectionId := k.GetNextCollectionId(ctx)
		k.IncrementNextCollectionId(ctx)

		collection = &types.BadgeCollection{
			CollectionId:                     nextCollectionId,
			InheritedCollectionId: 					 	msg.InheritedCollectionId,
			CollectionPermissions:            &types.CollectionPermissions{},
			BalancesType:                     msg.BalancesType,
			DefaultUserApprovedOutgoingTransfersTimeline: msg.DefaultApprovedOutgoingTransfersTimeline,
			DefaultUserApprovedIncomingTransfersTimeline: msg.DefaultApprovedIncomingTransfersTimeline,
			DefaultUserPermissions: 				msg.DefaultUserPermissions,
			CreatedBy: 											msg.Creator,
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
		}

		if IsInheritedBalances(collection) && (collection.InheritedCollectionId.IsNil() || collection.InheritedCollectionId.IsZero() ) {
			return nil, sdkerrors.Wrapf(ErrWrongBalancesType, "inherited balances are being set but collection %s does not have inherited balances", collection.CollectionId)
		}
		
	} else {
		found := false
		collection, found = k.GetCollectionFromStore(ctx, msg.CollectionId)
		if !found {
			return nil, ErrCollectionNotExists
		}
	}

	//Check must be manager
	err := k.UniversalValidate(ctx, collection, UniversalValidationParams{
		Creator:       msg.Creator,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

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
	//Other cases: 
	//previouslyArchived && not stillArchived - we have just unarchived the collection
	//not previouslyArchived && stillArchived - we have just archived the collection (all TXs moving forward will fail, but we allow this one)
	//not previouslyArchived && not stillArchived - unarchived before and now so we allow


	if msg.UpdateCollectionApprovedTransfersTimeline {
		if err := k.ValidateCollectionApprovedTransfersUpdate(ctx, collection, collection.CollectionApprovedTransfersTimeline, msg.CollectionApprovedTransfersTimeline, collection.CollectionPermissions.CanUpdateCollectionApprovedTransfers, msg.Creator); err != nil {
			return nil, err
		}
		collection.CollectionApprovedTransfersTimeline = msg.CollectionApprovedTransfersTimeline
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

	if msg.UpdateContractAddressTimeline {
		if err := k.ValidateContractAddressUpdate(ctx, collection.ContractAddressTimeline, msg.ContractAddressTimeline, collection.CollectionPermissions.CanUpdateContractAddress); err != nil {
			return nil, err
		}
		collection.ContractAddressTimeline = msg.ContractAddressTimeline
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
	
	collection, err = k.CreateBadges(ctx, collection, msg.BadgesToCreate)
	if err != nil {
		return nil, err
	}

	if msg.UpdateCollectionPermissions {
		err = k.ValidatePermissionsUpdate(ctx, collection.CollectionPermissions, msg.CollectionPermissions, msg.Creator)
		if err != nil {
			return nil, err
		}

		//iterate through the non-nil values
		if msg.CollectionPermissions.CanDeleteCollection != nil {
			collection.CollectionPermissions.CanDeleteCollection = msg.CollectionPermissions.CanDeleteCollection
		}

		if msg.CollectionPermissions.CanArchiveCollection != nil {
			collection.CollectionPermissions.CanArchiveCollection = msg.CollectionPermissions.CanArchiveCollection
		}

		if msg.CollectionPermissions.CanUpdateContractAddress != nil {
			collection.CollectionPermissions.CanUpdateContractAddress = msg.CollectionPermissions.CanUpdateContractAddress
		}

		if msg.CollectionPermissions.CanUpdateOffChainBalancesMetadata != nil {
			collection.CollectionPermissions.CanUpdateOffChainBalancesMetadata = msg.CollectionPermissions.CanUpdateOffChainBalancesMetadata
		}

		if msg.CollectionPermissions.CanUpdateCustomData != nil {
			collection.CollectionPermissions.CanUpdateCustomData = msg.CollectionPermissions.CanUpdateCustomData
		}

		if msg.CollectionPermissions.CanUpdateStandards != nil {
			collection.CollectionPermissions.CanUpdateStandards = msg.CollectionPermissions.CanUpdateStandards
		}

		if msg.CollectionPermissions.CanUpdateManager != nil {
			collection.CollectionPermissions.CanUpdateManager = msg.CollectionPermissions.CanUpdateManager
		}

		if msg.CollectionPermissions.CanUpdateCollectionMetadata != nil {
			collection.CollectionPermissions.CanUpdateCollectionMetadata = msg.CollectionPermissions.CanUpdateCollectionMetadata
		}

		if msg.CollectionPermissions.CanCreateMoreBadges != nil {
			collection.CollectionPermissions.CanCreateMoreBadges = msg.CollectionPermissions.CanCreateMoreBadges
		}

		if msg.CollectionPermissions.CanUpdateBadgeMetadata != nil {
			collection.CollectionPermissions.CanUpdateBadgeMetadata = msg.CollectionPermissions.CanUpdateBadgeMetadata
		}

		if msg.CollectionPermissions.CanUpdateCollectionApprovedTransfers != nil {
			collection.CollectionPermissions.CanUpdateCollectionApprovedTransfers = msg.CollectionPermissions.CanUpdateCollectionApprovedTransfers
		}

		collection.CollectionPermissions = msg.CollectionPermissions
	}

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
		),
	)

	return &types.MsgUpdateCollectionResponse{
		CollectionId: collection.CollectionId,
	}, nil
}
