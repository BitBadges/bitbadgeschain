package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteIncomingApproval(goCtx context.Context, msg *types.MsgDeleteIncomingApproval) (*types.MsgDeleteIncomingApprovalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	collectionId := msg.CollectionId
	if collectionId.Equal(sdkmath.NewUint(0)) {
		nextCollectionId := k.GetNextCollectionId(ctx)
		collectionId = nextCollectionId.Sub(sdkmath.NewUint(1))
	}

	// Get current user balance to filter out the approval to delete
	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	userBalance, _ := k.GetBalanceOrApplyDefault(ctx, collection, msg.Creator)
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
		return nil, fmt.Errorf("approval with ID %s not found", msg.ApprovalId)
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
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "delete_incoming_approval"),
			sdk.NewAttribute("msg", string(msgBytes)),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("indexer",
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "delete_incoming_approval"),
			sdk.NewAttribute("msg", string(msgBytes)),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
		),
	)

	return &types.MsgDeleteIncomingApprovalResponse{}, nil
}
