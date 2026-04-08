package keeper

import (
	"encoding/json"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ComputeCollectionApprovalChanges computes the approval changes between existing and new approval lists
// for collection-level approvals. Returns the list of changes and a summary review item string.
func ComputeCollectionApprovalChanges(
	existingApprovals []*types.CollectionApproval,
	newApprovals []*types.CollectionApproval,
) []*types.ApprovalChange {
	changes := []*types.ApprovalChange{}

	existingMap := make(map[string]*types.CollectionApproval)
	for _, a := range existingApprovals {
		existingMap[a.ApprovalId] = a
	}

	newMap := make(map[string]bool)
	for _, newApproval := range newApprovals {
		newMap[newApproval.ApprovalId] = true
		existing, exists := existingMap[newApproval.ApprovalId]
		if !exists {
			changes = append(changes, &types.ApprovalChange{
				ApprovalId:    newApproval.ApprovalId,
				ApprovalLevel: "collection",
				Action:        "created",
				Version:       newApproval.Version.String(),
			})
		} else if !collectionApprovalEqual(existing, newApproval) {
			changes = append(changes, &types.ApprovalChange{
				ApprovalId:    newApproval.ApprovalId,
				ApprovalLevel: "collection",
				Action:        "edited",
				Version:       newApproval.Version.String(),
			})
		}
		// If equal, no change - skip
	}

	// Check for deleted approvals
	for _, existing := range existingApprovals {
		if !newMap[existing.ApprovalId] {
			changes = append(changes, &types.ApprovalChange{
				ApprovalId:    existing.ApprovalId,
				ApprovalLevel: "collection",
				Action:        "deleted",
				Version:       existing.Version.String(),
			})
		}
	}

	return changes
}

// ComputeOutgoingApprovalChanges computes approval changes for user outgoing approvals.
func ComputeOutgoingApprovalChanges(
	existingApprovals []*types.UserOutgoingApproval,
	newApprovals []*types.UserOutgoingApproval,
) []*types.ApprovalChange {
	changes := []*types.ApprovalChange{}

	existingMap := make(map[string]*types.UserOutgoingApproval)
	for _, a := range existingApprovals {
		existingMap[a.ApprovalId] = a
	}

	newMap := make(map[string]bool)
	for _, newApproval := range newApprovals {
		newMap[newApproval.ApprovalId] = true
		existing, exists := existingMap[newApproval.ApprovalId]
		if !exists {
			changes = append(changes, &types.ApprovalChange{
				ApprovalId:    newApproval.ApprovalId,
				ApprovalLevel: "outgoing",
				Action:        "created",
				Version:       newApproval.Version.String(),
			})
		} else if !userOutgoingApprovalEqual(existing, newApproval) {
			changes = append(changes, &types.ApprovalChange{
				ApprovalId:    newApproval.ApprovalId,
				ApprovalLevel: "outgoing",
				Action:        "edited",
				Version:       newApproval.Version.String(),
			})
		}
	}

	for _, existing := range existingApprovals {
		if !newMap[existing.ApprovalId] {
			changes = append(changes, &types.ApprovalChange{
				ApprovalId:    existing.ApprovalId,
				ApprovalLevel: "outgoing",
				Action:        "deleted",
				Version:       existing.Version.String(),
			})
		}
	}

	return changes
}

// ComputeIncomingApprovalChanges computes approval changes for user incoming approvals.
func ComputeIncomingApprovalChanges(
	existingApprovals []*types.UserIncomingApproval,
	newApprovals []*types.UserIncomingApproval,
) []*types.ApprovalChange {
	changes := []*types.ApprovalChange{}

	existingMap := make(map[string]*types.UserIncomingApproval)
	for _, a := range existingApprovals {
		existingMap[a.ApprovalId] = a
	}

	newMap := make(map[string]bool)
	for _, newApproval := range newApprovals {
		newMap[newApproval.ApprovalId] = true
		existing, exists := existingMap[newApproval.ApprovalId]
		if !exists {
			changes = append(changes, &types.ApprovalChange{
				ApprovalId:    newApproval.ApprovalId,
				ApprovalLevel: "incoming",
				Action:        "created",
				Version:       newApproval.Version.String(),
			})
		} else if !userIncomingApprovalEqual(existing, newApproval) {
			changes = append(changes, &types.ApprovalChange{
				ApprovalId:    newApproval.ApprovalId,
				ApprovalLevel: "incoming",
				Action:        "edited",
				Version:       newApproval.Version.String(),
			})
		}
	}

	for _, existing := range existingApprovals {
		if !newMap[existing.ApprovalId] {
			changes = append(changes, &types.ApprovalChange{
				ApprovalId:    existing.ApprovalId,
				ApprovalLevel: "incoming",
				Action:        "deleted",
				Version:       existing.Version.String(),
			})
		}
	}

	return changes
}

// EmitApprovalChangeEvents emits one event per approval change.
func EmitApprovalChangeEvents(ctx sdk.Context, collectionId sdkmath.Uint, approverAddress string, changes []*types.ApprovalChange) {
	for _, change := range changes {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent("approvalChange",
				sdk.NewAttribute("collectionId", collectionId.String()),
				sdk.NewAttribute("approvalId", change.ApprovalId),
				sdk.NewAttribute("approvalLevel", change.ApprovalLevel),
				sdk.NewAttribute("approverAddress", approverAddress),
				sdk.NewAttribute("action", change.Action),
				sdk.NewAttribute("version", change.Version),
			),
		)
	}
}

// ApprovalChangeSummary returns a human-readable summary of approval changes.
func ApprovalChangeSummary(changes []*types.ApprovalChange) string {
	created, edited, deleted := 0, 0, 0
	for _, c := range changes {
		switch c.Action {
		case "created":
			created++
		case "edited":
			edited++
		case "deleted":
			deleted++
		}
	}
	return fmt.Sprintf("%d approvals created, %d edited, %d deleted", created, edited, deleted)
}

// CollectionReviewItems generates advisory review strings for collection updates.
func CollectionReviewItems(collection *types.TokenCollection, approvalChanges []*types.ApprovalChange) []string {
	items := []string{}

	if len(approvalChanges) > 0 {
		items = append(items, ApprovalChangeSummary(approvalChanges))
	}

	// Check if valid token IDs are empty
	if len(collection.ValidTokenIds) == 0 {
		items = append(items, "No valid token IDs set")
	}

	// Check if no transfer approvals are set
	if len(collection.CollectionApprovals) == 0 {
		items = append(items, "No transfer approvals set — tokens will be non-transferable")
	}

	return items
}

// TransferReviewItems generates advisory review strings for token transfers.
func TransferReviewItems(approvalsUsed []ApprovalsUsed) []string {
	items := []string{}

	for _, au := range approvalsUsed {
		items = append(items, fmt.Sprintf("Transfer used approval '%s' (level: %s, version: %s)", au.ApprovalId, au.ApprovalLevel, au.Version))
	}

	return items
}

// ConvertApprovalsUsedToProto converts the Go-only ApprovalsUsed to proto types.
func ConvertApprovalsUsedToProto(approvalsUsed []ApprovalsUsed) []*types.ApprovalUsed {
	result := make([]*types.ApprovalUsed, len(approvalsUsed))
	for i, au := range approvalsUsed {
		result[i] = &types.ApprovalUsed{
			ApprovalId:      au.ApprovalId,
			ApprovalLevel:   au.ApprovalLevel,
			ApproverAddress: au.ApproverAddress,
			Version:         au.Version,
		}
	}
	return result
}

// ConvertCoinTransfersToProto converts the Go-only CoinTransfers to proto types.
func ConvertCoinTransfersToProto(coinTransfers []CoinTransfers) []*types.CoinTransferProto {
	result := make([]*types.CoinTransferProto, len(coinTransfers))
	for i, ct := range coinTransfers {
		result[i] = &types.CoinTransferProto{
			From:          ct.From,
			To:            ct.To,
			Amount:        ct.Amount,
			Denom:         ct.Denom,
			IsProtocolFee: ct.IsProtocolFee,
		}
	}
	return result
}

// EmitApprovalChangeEventsWithJSON emits events with the full approval JSON attached.
func EmitApprovalChangeEventsWithJSON(ctx sdk.Context, collectionId sdkmath.Uint, approverAddress string, changes []*types.ApprovalChange, approvalJSON map[string]string) {
	for _, change := range changes {
		attrs := []sdk.Attribute{
			sdk.NewAttribute("collectionId", collectionId.String()),
			sdk.NewAttribute("approvalId", change.ApprovalId),
			sdk.NewAttribute("approvalLevel", change.ApprovalLevel),
			sdk.NewAttribute("approverAddress", approverAddress),
			sdk.NewAttribute("action", change.Action),
			sdk.NewAttribute("version", change.Version),
		}
		if jsonStr, ok := approvalJSON[change.ApprovalId]; ok {
			attrs = append(attrs, sdk.NewAttribute("approval", jsonStr))
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent("approvalChange", attrs...),
		)
	}
}

// BuildApprovalJSONMap builds a map of approvalId -> JSON string for collection approvals.
func BuildCollectionApprovalJSONMap(approvals []*types.CollectionApproval) map[string]string {
	result := make(map[string]string)
	for _, a := range approvals {
		jsonBytes, err := json.Marshal(a)
		if err == nil {
			result[a.ApprovalId] = string(jsonBytes)
		}
	}
	return result
}

// BuildOutgoingApprovalJSONMap builds a map of approvalId -> JSON string for outgoing approvals.
func BuildOutgoingApprovalJSONMap(approvals []*types.UserOutgoingApproval) map[string]string {
	result := make(map[string]string)
	for _, a := range approvals {
		jsonBytes, err := json.Marshal(a)
		if err == nil {
			result[a.ApprovalId] = string(jsonBytes)
		}
	}
	return result
}

// BuildIncomingApprovalJSONMap builds a map of approvalId -> JSON string for incoming approvals.
func BuildIncomingApprovalJSONMap(approvals []*types.UserIncomingApproval) map[string]string {
	result := make(map[string]string)
	for _, a := range approvals {
		jsonBytes, err := json.Marshal(a)
		if err == nil {
			result[a.ApprovalId] = string(jsonBytes)
		}
	}
	return result
}
