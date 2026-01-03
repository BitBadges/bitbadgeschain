package keeper

import (
	"context"
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PurgeApprovals handles the purging of specific user-level approvals (incoming/outgoing)
func (k msgServer) PurgeApprovals(goCtx context.Context, msg *types.MsgPurgeApprovals) (*types.MsgPurgeApprovalsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	collectionId, err := k.resolveCollectionIdWithAutoPrev(ctx, msg.CollectionId)
	if err != nil {
		return nil, err
	}

	collection, found := k.GetCollectionFromStore(ctx, collectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	targetAddress := msg.ApproverAddress
	if targetAddress == "" {
		targetAddress = msg.Creator
	}

	if targetAddress != msg.Creator {
		if !msg.PurgeCounterpartyApprovals {
			return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "creator %s cannot purge approvals for address %s without PurgeCounterpartyApprovals flag", msg.Creator, targetAddress)
		}
	}

	numPurged := sdkmath.NewUint(0)
	currentTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))

	// Get current user balance to filter out approvals to purge
	userBalance, _ := k.GetBalanceOrApplyDefault(ctx, collection, targetAddress)
	if userBalance.UserPermissions == nil {
		userBalance.UserPermissions = &types.UserPermissions{}
	}

	// Find user approvals to purge
	incomingApprovalsToPurge := []*types.ApprovalIdentifierDetails{}
	outgoingApprovalsToPurge := []*types.ApprovalIdentifierDetails{}

	for _, approval := range msg.ApprovalsToPurge {
		if approval.ApprovalLevel == "incoming" && approval.ApproverAddress == targetAddress {
			incomingApprovalsToPurge = append(incomingApprovalsToPurge, approval)
		} else if approval.ApprovalLevel == "outgoing" && approval.ApproverAddress == targetAddress {
			outgoingApprovalsToPurge = append(outgoingApprovalsToPurge, approval)
		}
	}

	// Filter incoming approvals
	if len(incomingApprovalsToPurge) > 0 {
		newIncomingApprovals := k.filterApprovalsToPurge(ctx, userBalance.IncomingApprovals, incomingApprovalsToPurge, msg, targetAddress, currentTime, &numPurged)
		userBalance.IncomingApprovals = newIncomingApprovals.([]*types.UserIncomingApproval)
	}

	// Filter outgoing approvals
	if len(outgoingApprovalsToPurge) > 0 {
		newOutgoingApprovals := k.filterApprovalsToPurge(ctx, userBalance.OutgoingApprovals, outgoingApprovalsToPurge, msg, targetAddress, currentTime, &numPurged)
		userBalance.OutgoingApprovals = newOutgoingApprovals.([]*types.UserOutgoingApproval)
	}

	// Create UpdateUserApprovals message with filtered approvals
	updateMsg := &types.MsgUpdateUserApprovals{
		UpdateIncomingApprovals: len(incomingApprovalsToPurge) > 0,
		IncomingApprovals:       userBalance.IncomingApprovals,
		UpdateOutgoingApprovals: len(outgoingApprovalsToPurge) > 0,
		OutgoingApprovals:       userBalance.OutgoingApprovals,
		// Set other fields to false to keep existing values
		UpdateAutoApproveSelfInitiatedIncomingTransfers: false,
		UpdateAutoApproveSelfInitiatedOutgoingTransfers: false,
		UpdateAutoApproveAllIncomingTransfers:           false,
		UpdateUserPermissions:                           false,
	}

	// Execute the update using the helper function
	err = k.executeUpdateUserApprovals(ctx, targetAddress, collectionId, updateMsg)
	if err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "purge_approvals"),
		sdk.NewAttribute("msg", msgStr),
		sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
	)

	return &types.MsgPurgeApprovalsResponse{NumPurged: numPurged}, nil
}

// filterApprovalsToPurge filters out approvals that should be purged
func (k msgServer) filterApprovalsToPurge(ctx sdk.Context, approvals interface{}, approvalsToPurge []*types.ApprovalIdentifierDetails, msg *types.MsgPurgeApprovals, targetAddress string, currentTime sdkmath.Uint, numPurged *sdkmath.Uint) interface{} {
	// Create a map for quick lookup with version
	approvalsToPurgeMap := make(map[string]sdkmath.Uint)
	for _, approval := range approvalsToPurge {
		key := fmt.Sprintf("%s:%s", approval.ApprovalId, targetAddress)
		approvalsToPurgeMap[key] = approval.Version
	}

	switch a := approvals.(type) {
	case []*types.UserIncomingApproval:
		newApprovals := []*types.UserIncomingApproval{}
		for _, approval := range a {
			key := fmt.Sprintf("%s:%s", approval.ApprovalId, targetAddress)
			if expectedVersion, exists := approvalsToPurgeMap[key]; exists {
				// Check if version matches
				if approval.Version.Equal(expectedVersion) {
					if k.canPurgeUserApproval(ctx, approval, msg, targetAddress, currentTime) {
						*numPurged = numPurged.Add(sdkmath.NewUint(1))
						continue // skip adding to newApprovals
					}
				}
			}
			newApprovals = append(newApprovals, approval)
		}
		return newApprovals
	case []*types.UserOutgoingApproval:
		newApprovals := []*types.UserOutgoingApproval{}
		for _, approval := range a {
			key := fmt.Sprintf("%s:%s", approval.ApprovalId, targetAddress)
			if expectedVersion, exists := approvalsToPurgeMap[key]; exists {
				// Check if version matches
				if approval.Version.Equal(expectedVersion) {
					if k.canPurgeUserApproval(ctx, approval, msg, targetAddress, currentTime) {
						*numPurged = numPurged.Add(sdkmath.NewUint(1))
						continue // skip adding to newApprovals
					}
				}
			}
			newApprovals = append(newApprovals, approval)
		}
		return newApprovals
	default:
		return approvals
	}
}

// canPurgeUserApproval checks if a user approval can be purged
func (k msgServer) canPurgeUserApproval(ctx sdk.Context, approval interface{}, msg *types.MsgPurgeApprovals, targetAddress string, currentTime sdkmath.Uint) bool {
	// Check if approval has future transfer times
	hasFutureTransferTimes := false

	switch a := approval.(type) {
	case *types.UserIncomingApproval:
		for _, transferTime := range a.TransferTimes {
			if transferTime.End.GT(currentTime) {
				hasFutureTransferTimes = true
				break
			}
		}
	case *types.UserOutgoingApproval:
		for _, transferTime := range a.TransferTimes {
			if transferTime.End.GT(currentTime) {
				hasFutureTransferTimes = true
				break
			}
		}
	}

	// If approval has future transfer times, it cannot be purged
	if hasFutureTransferTimes {
		return false
	}

	// Self-purge with expired approvals is always allowed
	if targetAddress == msg.Creator {
		return true
	}

	// Check auto-deletion options for permission
	switch a := approval.(type) {
	case *types.UserIncomingApproval:
		// For incoming approvals, check if counterparty purge is allowed
		if msg.PurgeCounterpartyApprovals && a.ApprovalCriteria != nil && a.ApprovalCriteria.AutoDeletionOptions != nil && a.ApprovalCriteria.AutoDeletionOptions.AllowCounterpartyPurge {
			// Check if creator is the only initiator
			if k.IsWhitelistWithSingleAddress(ctx, a.InitiatedByListId, msg.Creator) {
				return true
			}
		}
		// Also check if others can purge expired approvals
		if a.ApprovalCriteria != nil && a.ApprovalCriteria.AutoDeletionOptions != nil {
			return a.ApprovalCriteria.AutoDeletionOptions.AllowPurgeIfExpired
		}
		return false
	case *types.UserOutgoingApproval:
		// For outgoing approvals, check if counterparty purge is allowed
		if msg.PurgeCounterpartyApprovals && a.ApprovalCriteria != nil && a.ApprovalCriteria.AutoDeletionOptions != nil && a.ApprovalCriteria.AutoDeletionOptions.AllowCounterpartyPurge {
			// Check if creator is the only initiator
			if k.IsWhitelistWithSingleAddress(ctx, a.InitiatedByListId, msg.Creator) {
				return true
			}
		}
		// Also check if others can purge expired approvals
		if a.ApprovalCriteria != nil && a.ApprovalCriteria.AutoDeletionOptions != nil {
			return a.ApprovalCriteria.AutoDeletionOptions.AllowPurgeIfExpired
		}
		return false
	}

	return false
}

// IsWhitelistWithSingleAddress checks if the list is a whitelist with exactly one address matching the given address.
func (k msgServer) IsWhitelistWithSingleAddress(ctx sdk.Context, listId string, address string) bool {
	keeper := k.Keeper // msgServer embeds Keeper
	list, err := keeper.GetAddressListById(ctx, listId)
	if err != nil {
		return false
	}
	if !list.Whitelist {
		return false
	}
	if len(list.Addresses) != 1 {
		return false
	}
	return list.Addresses[0] == address
}
