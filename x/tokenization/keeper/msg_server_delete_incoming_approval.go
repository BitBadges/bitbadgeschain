package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteIncomingApproval(goCtx context.Context, msg *types.MsgDeleteIncomingApproval) (*types.MsgDeleteIncomingApprovalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	collectionId, err := k.resolveCollectionIdWithAutoPrev(ctx, msg.CollectionId)
	if err != nil {
		return nil, err
	}

	// Get current user balance to filter out the approval to delete
	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, errorsmod.Wrapf(ErrCollectionNotExists, "collection ID %s not found", collectionId.String())
	}

	userBalance, _, err := k.GetBalanceOrApplyDefault(ctx, collection, msg.Creator)
	if err != nil {
		return nil, err
	}
	if userBalance.UserPermissions == nil {
		userBalance.UserPermissions = &types.UserPermissions{}
	}

	// Find and remove the approval
	foundApproval := false
	newIncomingApprovals := []*types.UserIncomingApproval{}
	for _, approval := range userBalance.IncomingApprovals {
		if approval.ApprovalId == msg.ApprovalId {
			foundApproval = true
			// Skip this approval (delete it)
		} else {
			newIncomingApprovals = append(newIncomingApprovals, approval)
		}
	}

	if !foundApproval {
		return nil, errorsmod.Wrapf(ErrApprovalNotFound, "approval ID: %s", msg.ApprovalId)
	}

	// Create UpdateUserApprovals message with filtered incoming approvals
	updateMsg := &types.MsgUpdateUserApprovals{
		UpdateIncomingApprovals: true,
		IncomingApprovals:       newIncomingApprovals,
		// Set other fields to false to keep existing values
		UpdateOutgoingApprovals:                         false,
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
		sdk.NewAttribute("msg_type", "delete_incoming_approval"),
		sdk.NewAttribute("msg", msgStr),
		sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
	)

	return &types.MsgDeleteIncomingApprovalResponse{}, nil
}
