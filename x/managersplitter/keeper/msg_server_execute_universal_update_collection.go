package keeper

import (
	"context"
	"slices"

	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// isAddressApproved checks if an address is in the approved addresses list
func isAddressApproved(address string, approvedAddresses []string) bool {
	return slices.Contains(approvedAddresses, address)
}

// checkPermission checks if the executor has permission to execute a specific action
func (k Keeper) checkPermission(ctx sdk.Context, executor string, managerSplitter *types.ManagerSplitter, permissionName string) error {
	// Admin always has full permissions
	if executor == managerSplitter.Admin {
		return nil
	}

	// Get the permission criteria for this permission
	var criteria *types.PermissionCriteria
	switch permissionName {
	case "canDeleteCollection":
		if managerSplitter.Permissions != nil && managerSplitter.Permissions.CanDeleteCollection != nil {
			criteria = managerSplitter.Permissions.CanDeleteCollection
		}
	case "canArchiveCollection":
		if managerSplitter.Permissions != nil && managerSplitter.Permissions.CanArchiveCollection != nil {
			criteria = managerSplitter.Permissions.CanArchiveCollection
		}
	case "canUpdateStandards":
		if managerSplitter.Permissions != nil && managerSplitter.Permissions.CanUpdateStandards != nil {
			criteria = managerSplitter.Permissions.CanUpdateStandards
		}
	case "canUpdateCustomData":
		if managerSplitter.Permissions != nil && managerSplitter.Permissions.CanUpdateCustomData != nil {
			criteria = managerSplitter.Permissions.CanUpdateCustomData
		}
	case "canUpdateManager":
		if managerSplitter.Permissions != nil && managerSplitter.Permissions.CanUpdateManager != nil {
			criteria = managerSplitter.Permissions.CanUpdateManager
		}
	case "canUpdateCollectionMetadata":
		if managerSplitter.Permissions != nil && managerSplitter.Permissions.CanUpdateCollectionMetadata != nil {
			criteria = managerSplitter.Permissions.CanUpdateCollectionMetadata
		}
	case "canUpdateValidTokenIds":
		if managerSplitter.Permissions != nil && managerSplitter.Permissions.CanUpdateValidTokenIds != nil {
			criteria = managerSplitter.Permissions.CanUpdateValidTokenIds
		}
	case "canUpdateTokenMetadata":
		if managerSplitter.Permissions != nil && managerSplitter.Permissions.CanUpdateTokenMetadata != nil {
			criteria = managerSplitter.Permissions.CanUpdateTokenMetadata
		}
	case "canUpdateCollectionApprovals":
		if managerSplitter.Permissions != nil && managerSplitter.Permissions.CanUpdateCollectionApprovals != nil {
			criteria = managerSplitter.Permissions.CanUpdateCollectionApprovals
		}
	}

	// If no criteria set, deny by default (except admin)
	if criteria == nil {
		return sdkerrors.Wrap(types.ErrPermissionDenied, "no permission criteria set for "+permissionName)
	}

	// Check if executor is in approved addresses
	if !isAddressApproved(executor, criteria.ApprovedAddresses) {
		return sdkerrors.Wrap(types.ErrPermissionDenied, "executor not approved for "+permissionName)
	}

	return nil
}

// checkAllPermissions checks all permissions that would be used by the UniversalUpdateCollection message
func (k Keeper) checkAllPermissions(ctx sdk.Context, executor string, managerSplitter *types.ManagerSplitter, msg *badgestypes.MsgUniversalUpdateCollection) error {
	// Check permissions based on which fields are being updated
	if msg.UpdateValidTokenIds {
		if err := k.checkPermission(ctx, executor, managerSplitter, "canUpdateValidTokenIds"); err != nil {
			return err
		}
	}

	if msg.UpdateCollectionPermissions {
		// Updating permissions requires checking all permission types
		// For simplicity, we'll check canUpdateCollectionApprovals as a proxy
		// In a more sophisticated implementation, we might want to check each permission individually
		if err := k.checkPermission(ctx, executor, managerSplitter, "canUpdateCollectionApprovals"); err != nil {
			return err
		}
	}

	if msg.UpdateManagerTimeline {
		if err := k.checkPermission(ctx, executor, managerSplitter, "canUpdateManager"); err != nil {
			return err
		}
	}

	if msg.UpdateCollectionMetadataTimeline {
		if err := k.checkPermission(ctx, executor, managerSplitter, "canUpdateCollectionMetadata"); err != nil {
			return err
		}
	}

	if msg.UpdateTokenMetadataTimeline {
		if err := k.checkPermission(ctx, executor, managerSplitter, "canUpdateTokenMetadata"); err != nil {
			return err
		}
	}

	if msg.UpdateCustomDataTimeline {
		if err := k.checkPermission(ctx, executor, managerSplitter, "canUpdateCustomData"); err != nil {
			return err
		}
	}

	if msg.UpdateCollectionApprovals {
		if err := k.checkPermission(ctx, executor, managerSplitter, "canUpdateCollectionApprovals"); err != nil {
			return err
		}
	}

	if msg.UpdateStandardsTimeline {
		if err := k.checkPermission(ctx, executor, managerSplitter, "canUpdateStandards"); err != nil {
			return err
		}
	}

	if msg.UpdateIsArchivedTimeline {
		if err := k.checkPermission(ctx, executor, managerSplitter, "canArchiveCollection"); err != nil {
			return err
		}
	}

	// Note: canDeleteCollection is checked separately if needed
	// For now, we don't have a delete flag in UniversalUpdateCollection

	return nil
}

func (k msgServer) ExecuteUniversalUpdateCollection(goCtx context.Context, msg *types.MsgExecuteUniversalUpdateCollection) (*types.MsgExecuteUniversalUpdateCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate executor address
	_, err := sdk.AccAddressFromBech32(msg.Executor)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAddress, "invalid executor address")
	}

	// Get manager splitter
	managerSplitter, found := k.GetManagerSplitterFromStore(ctx, msg.ManagerSplitterAddress)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrManagerSplitterNotFound, msg.ManagerSplitterAddress)
	}

	// Check all permissions before executing
	if err := k.checkAllPermissions(ctx, msg.Executor, managerSplitter, msg.UniversalUpdateCollectionMsg); err != nil {
		return nil, err
	}

	// Create a new message with the manager splitter address as the creator
	// This ensures the badges module sees the manager splitter as the manager
	badgesMsg := *msg.UniversalUpdateCollectionMsg
	badgesMsg.Creator = managerSplitter.Address

	// Call the badges module's UniversalUpdateCollection through msgServer
	// We need to create a msgServer instance to call the method
	badgesMsgServer := badgeskeeper.NewMsgServerImpl(k.badgesKeeper)
	response, err := badgesMsgServer.UniversalUpdateCollection(goCtx, &badgesMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to execute UniversalUpdateCollection")
	}

	return &types.MsgExecuteUniversalUpdateCollectionResponse{
		CollectionId: response.CollectionId,
	}, nil
}
