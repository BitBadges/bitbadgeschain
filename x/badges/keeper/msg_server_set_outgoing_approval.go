package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetOutgoingApproval(goCtx context.Context, msg *types.MsgSetOutgoingApproval) (*types.MsgSetOutgoingApprovalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	collectionId, err := k.resolveCollectionIdWithAutoPrev(ctx, msg.CollectionId)
	if err != nil {
		return nil, err
	}

	// Get current user balance to build the new approval list
	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	userBalance, _, err := k.GetBalanceOrApplyDefault(ctx, collection, msg.Creator)
	if err != nil {
		return nil, err
	}
	if userBalance.UserPermissions == nil {
		userBalance.UserPermissions = &types.UserPermissions{}
	}

	// Build new outgoing approvals list - replace existing or add new
	newOutgoingApprovals := make([]*types.UserOutgoingApproval, 0, len(userBalance.OutgoingApprovals)+1)
	foundExisting := false

	for _, approval := range userBalance.OutgoingApprovals {
		if approval.ApprovalId == msg.Approval.ApprovalId {
			newOutgoingApprovals = append(newOutgoingApprovals, msg.Approval)
			foundExisting = true
		} else {
			newOutgoingApprovals = append(newOutgoingApprovals, approval)
		}
	}

	if !foundExisting {
		newOutgoingApprovals = append(newOutgoingApprovals, msg.Approval)
	}

	// Create UpdateUserApprovals message with the new approval list
	updateMsg := &types.MsgUpdateUserApprovals{
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       newOutgoingApprovals,
		// Set other fields to false to keep existing values
		UpdateIncomingApprovals:                         false,
		UpdateAutoApproveSelfInitiatedIncomingTransfers: false,
		UpdateAutoApproveSelfInitiatedOutgoingTransfers: false,
		UpdateAutoApproveAllIncomingTransfers:           false,
		UpdateUserPermissions:                           false,
	}

	// Execute the update using the helper function
	err = k.executeUpdateUserApprovals(ctx, msg.Creator, collectionId, updateMsg)
	if err != nil {
		return nil, err
	}

	// Emit events
	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "set_outgoing_approval"),
		sdk.NewAttribute("msg", msgStr),
		sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
	)

	return &types.MsgSetOutgoingApprovalResponse{}, nil
}
