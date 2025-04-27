package keeper

import (
	"context"
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateUserApprovals(goCtx context.Context, msg *types.MsgUpdateUserApprovals) (*types.MsgUpdateUserApprovalsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := k.GetCreator(ctx, msg.Creator, msg.CreatorOverride)
	if err != nil {
		return nil, err
	}
	msg.Creator = creator

	err = msg.CheckAndCleanMsg(ctx, true)
	if err != nil {
		return nil, err
	}

	collection, found := k.GetCollectionFromStore(ctx, msg.CollectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	isArchived := types.GetIsArchived(ctx, collection)
	if isArchived {
		return nil, ErrCollectionIsArchived
	}

	if !IsStandardBalances(collection) {
		return nil, ErrWrongBalancesType
	}

	userBalance, appliedDefault := k.GetBalanceOrApplyDefault(ctx, collection, msg.Creator)
	if userBalance.UserPermissions == nil {
		userBalance.UserPermissions = &types.UserPermissions{}
	}

	if msg.UpdateOutgoingApprovals {
		if err := k.ValidateUserOutgoingApprovalsUpdate(ctx, collection, userBalance.OutgoingApprovals, msg.OutgoingApprovals, userBalance.UserPermissions.CanUpdateOutgoingApprovals, msg.Creator); err != nil {
			return nil, err
		}
		userBalance.OutgoingApprovals = msg.OutgoingApprovals

		// If we didn't apply the default, we need to increment the versions
		// Else, we did apply the default and we should ensure the version is kept at 0 - no need to double increment
		if !appliedDefault {
			newOutgoingApprovalsWithVersion := []*types.UserOutgoingApproval{}
			for _, approval := range msg.OutgoingApprovals {
				newVersion := k.IncrementApprovalVersion(ctx, collection.CollectionId, "outgoing", msg.Creator, approval.ApprovalId)
				approval.Version = newVersion
				newOutgoingApprovalsWithVersion = append(newOutgoingApprovalsWithVersion, approval)
			}
			userBalance.OutgoingApprovals = newOutgoingApprovalsWithVersion
		} else {
			// We did apply the default, so we need to ensure the version is kept at 0 and no need to increment again
			for _, approval := range msg.OutgoingApprovals {
				approval.Version = sdkmath.NewUint(0)
			}
			userBalance.OutgoingApprovals = msg.OutgoingApprovals
		}
	}

	if msg.UpdateIncomingApprovals {
		if err := k.ValidateUserIncomingApprovalsUpdate(ctx, collection, userBalance.IncomingApprovals, msg.IncomingApprovals, userBalance.UserPermissions.CanUpdateIncomingApprovals, msg.Creator); err != nil {
			return nil, err
		}
		userBalance.IncomingApprovals = msg.IncomingApprovals

		// If we didn't apply the default, we need to increment the versions
		// Else, we did apply the default and we should ensure the version is kept at 0 - no need to double increment
		if !appliedDefault {
			newIncomingApprovalsWithVersion := []*types.UserIncomingApproval{}
			for _, approval := range msg.IncomingApprovals {
				newVersion := k.IncrementApprovalVersion(ctx, collection.CollectionId, "incoming", msg.Creator, approval.ApprovalId)
				approval.Version = newVersion
				newIncomingApprovalsWithVersion = append(newIncomingApprovalsWithVersion, approval)
			}
			userBalance.IncomingApprovals = newIncomingApprovalsWithVersion
		} else {
			// We did apply the default, so we need to ensure the version is kept at 0 and no need to increment again
			for _, approval := range msg.IncomingApprovals {
				approval.Version = sdkmath.NewUint(0)
			}
			userBalance.IncomingApprovals = msg.IncomingApprovals
		}
	}

	if msg.UpdateAutoApproveSelfInitiatedIncomingTransfers && userBalance.AutoApproveSelfInitiatedIncomingTransfers != msg.AutoApproveSelfInitiatedIncomingTransfers {
		//Check permission is valid for current time
		err = k.CheckIfActionPermissionPermits(ctx, userBalance.UserPermissions.CanUpdateAutoApproveSelfInitiatedIncomingTransfers, "can update auto approve self initiated incoming transfers")
		if err != nil {
			return nil, err
		}
		userBalance.AutoApproveSelfInitiatedIncomingTransfers = msg.AutoApproveSelfInitiatedIncomingTransfers
	}

	if msg.UpdateAutoApproveSelfInitiatedOutgoingTransfers && userBalance.AutoApproveSelfInitiatedOutgoingTransfers != msg.AutoApproveSelfInitiatedOutgoingTransfers {
		//Check permission is valid for current time
		err = k.CheckIfActionPermissionPermits(ctx, userBalance.UserPermissions.CanUpdateAutoApproveSelfInitiatedOutgoingTransfers, "can update auto approve self initiated outgoing transfers")
		if err != nil {
			return nil, err
		}
		userBalance.AutoApproveSelfInitiatedOutgoingTransfers = msg.AutoApproveSelfInitiatedOutgoingTransfers
	}

	if msg.UpdateAutoApproveAllIncomingTransfers && userBalance.AutoApproveAllIncomingTransfers != msg.AutoApproveAllIncomingTransfers {
		//Check permission is valid for current time
		err = k.CheckIfActionPermissionPermits(ctx, userBalance.UserPermissions.CanUpdateAutoApproveAllIncomingTransfers, "can update auto approve all incoming transfers")
		if err != nil {
			return nil, err
		}
		userBalance.AutoApproveAllIncomingTransfers = msg.AutoApproveAllIncomingTransfers
	}

	if msg.UpdateUserPermissions {
		err := k.ValidateUserPermissionsUpdate(ctx, userBalance.UserPermissions, msg.UserPermissions)
		if err != nil {
			return nil, err
		}

		userBalance.UserPermissions = msg.UserPermissions
	}

	err = k.SetBalanceForAddress(ctx, collection, msg.Creator, userBalance)
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
			sdk.NewAttribute("msg_type", "update_user_approvals"),
			sdk.NewAttribute("msg", string(msgBytes)),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("indexer",
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "update_user_approvals"),
			sdk.NewAttribute("msg", string(msgBytes)),
		),
	)

	return &types.MsgUpdateUserApprovalsResponse{}, nil
}
