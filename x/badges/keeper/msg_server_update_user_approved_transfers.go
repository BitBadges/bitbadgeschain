package keeper

import (
	"context"
	"encoding/json"

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

	userBalance := k.GetBalanceOrApplyDefault(ctx, collection, msg.Creator)
	if userBalance.UserPermissions == nil {
		userBalance.UserPermissions = &types.UserPermissions{}
	}

	if msg.UpdateOutgoingApprovals {
		if err := k.ValidateUserOutgoingApprovalsUpdate(ctx, collection, userBalance.OutgoingApprovals, msg.OutgoingApprovals, userBalance.UserPermissions.CanUpdateOutgoingApprovals, msg.Creator); err != nil {
			return nil, err
		}
		userBalance.OutgoingApprovals = msg.OutgoingApprovals
	}

	if msg.UpdateIncomingApprovals {
		if err := k.ValidateUserIncomingApprovalsUpdate(ctx, collection, userBalance.IncomingApprovals, msg.IncomingApprovals, userBalance.UserPermissions.CanUpdateIncomingApprovals, msg.Creator); err != nil {
			return nil, err
		}
		userBalance.IncomingApprovals = msg.IncomingApprovals
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
			sdk.NewAttribute("msg_type", "update_user_approved_transfers"),
			sdk.NewAttribute("msg", string(msgBytes)),
		),
	)

	return &types.MsgUpdateUserApprovalsResponse{}, nil
}
