package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v13"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastOldSetOutgoingApprovalToNewType(oldMsg *oldtypes.MsgSetOutgoingApproval) (*types.MsgSetOutgoingApproval, error) {
	// Convert to JSON
	jsonBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return nil, err
	}

	// Unmarshal into new type
	var newMsg types.MsgSetOutgoingApproval
	if err := json.Unmarshal(jsonBytes, &newMsg); err != nil {
		return nil, err
	}

	return &newMsg, nil
}

func (k msgServer) SetOutgoingApprovalV13(goCtx context.Context, msg *oldtypes.MsgSetOutgoingApproval) (*types.MsgSetOutgoingApprovalResponse, error) {
	newMsg, err := CastOldSetOutgoingApprovalToNewType(msg)
	if err != nil {
		return nil, err
	}
	return k.SetOutgoingApproval(goCtx, newMsg)
}

func (k msgServer) SetOutgoingApprovalV14(goCtx context.Context, msg *types.MsgSetOutgoingApproval) (*types.MsgSetOutgoingApprovalResponse, error) {
	return k.SetOutgoingApproval(goCtx, msg)
}

func (k msgServer) SetOutgoingApproval(goCtx context.Context, msg *types.MsgSetOutgoingApproval) (*types.MsgSetOutgoingApprovalResponse, error) {
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
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "set_outgoing_approval"),
			sdk.NewAttribute("msg", string(msgBytes)),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("indexer",
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "set_outgoing_approval"),
			sdk.NewAttribute("msg", string(msgBytes)),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
		),
	)

	return &types.MsgSetOutgoingApprovalResponse{}, nil
}
