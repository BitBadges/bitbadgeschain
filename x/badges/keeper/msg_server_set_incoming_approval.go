package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetIncomingApproval(goCtx context.Context, msg *types.MsgSetIncomingApproval) (*types.MsgSetIncomingApprovalResponse, error) {
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

	// Get current user balance to build the new approval list
	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	userBalance, _ := k.GetBalanceOrApplyDefault(ctx, collection, msg.Creator)
	if userBalance.UserPermissions == nil {
		userBalance.UserPermissions = &types.UserPermissions{}
	}

	// Build new incoming approvals list - replace existing or add new
	newIncomingApprovals := make([]*types.UserIncomingApproval, 0, len(userBalance.IncomingApprovals)+1)
	foundExisting := false

	for _, approval := range userBalance.IncomingApprovals {
		if approval.ApprovalId == msg.Approval.ApprovalId {
			newIncomingApprovals = append(newIncomingApprovals, msg.Approval)
			foundExisting = true
		} else {
			newIncomingApprovals = append(newIncomingApprovals, approval)
		}
	}

	if !foundExisting {
		newIncomingApprovals = append(newIncomingApprovals, msg.Approval)
	}

	// Create UpdateUserApprovals message with the new approval list
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
			sdk.NewAttribute("msg_type", "set_incoming_approval"),
			sdk.NewAttribute("msg", string(msgBytes)),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("indexer",
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "set_incoming_approval"),
			sdk.NewAttribute("msg", string(msgBytes)),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
		),
	)

	return &types.MsgSetIncomingApprovalResponse{}, nil
}
