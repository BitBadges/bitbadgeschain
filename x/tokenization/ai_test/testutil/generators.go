package testutil

import (
	"math"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// GenerateUintRange generates a UintRange with the given start and end
func GenerateUintRange(start, end uint64) *types.UintRange {
	return &types.UintRange{
		Start: sdkmath.NewUint(start),
		End:   sdkmath.NewUint(end),
	}
}

// GenerateMaxUintRange generates a UintRange from start to MaxUint64
func GenerateMaxUintRange(start uint64) *types.UintRange {
	return &types.UintRange{
		Start: sdkmath.NewUint(start),
		End:   sdkmath.NewUint(math.MaxUint64),
	}
}

// GenerateBalance generates a balance with the given parameters
func GenerateBalance(amount, tokenIdStart, tokenIdEnd, ownershipTimeStart, ownershipTimeEnd uint64) *types.Balance {
	return &types.Balance{
		Amount: sdkmath.NewUint(amount),
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(tokenIdStart), End: sdkmath.NewUint(tokenIdEnd)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(ownershipTimeStart), End: sdkmath.NewUint(ownershipTimeEnd)},
		},
	}
}

// GenerateSimpleBalance generates a simple balance for a single token ID
func GenerateSimpleBalance(amount, tokenId uint64) *types.Balance {
	return GenerateBalance(amount, tokenId, tokenId, 1, math.MaxUint64)
}

// GenerateCollectionApproval generates a collection approval with default settings
func GenerateCollectionApproval(approvalId string, fromListId, toListId string) *types.CollectionApproval {
	return &types.CollectionApproval{
		ApprovalId:        approvalId,
		FromListId:        fromListId,
		ToListId:          toListId,
		InitiatedByListId: "All", // Default to "All" for initiated by list
		TransferTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		ApprovalCriteria: &types.ApprovalCriteria{},
		Version:          sdkmath.NewUint(0),
	}
}

// GenerateUserOutgoingApproval generates a user outgoing approval
func GenerateUserOutgoingApproval(approvalId, toListId string) *types.UserOutgoingApproval {
	return &types.UserOutgoingApproval{
		ApprovalId:        approvalId,
		ToListId:          toListId,
		InitiatedByListId: "All", // Default to "All" for initiated by list
		TransferTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		ApprovalCriteria: &types.OutgoingApprovalCriteria{},
		Version:          sdkmath.NewUint(0),
	}
}

// GenerateUserIncomingApproval generates a user incoming approval
func GenerateUserIncomingApproval(approvalId, fromListId string) *types.UserIncomingApproval {
	return &types.UserIncomingApproval{
		ApprovalId:        approvalId,
		FromListId:        fromListId,
		InitiatedByListId: "All", // Default to "All" for initiated by list
		TransferTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		ApprovalCriteria: &types.IncomingApprovalCriteria{},
		Version:          sdkmath.NewUint(0),
	}
}

// GenerateAddressList generates an address list
func GenerateAddressList(listId string, addresses []string, whitelist bool) *types.AddressList {
	return &types.AddressList{
		ListId:     listId,
		Addresses:  addresses,
		Whitelist:  whitelist,
		Uri:        "",
		CustomData: "",
	}
}

// GenerateMerkleChallenge generates a Merkle challenge
func GenerateMerkleChallenge(challengeTrackerId, root string, expectedProofLength uint64) *types.MerkleChallenge {
	return &types.MerkleChallenge{
		Root:                    root,
		ExpectedProofLength:     sdkmath.NewUint(expectedProofLength),
		UseCreatorAddressAsLeaf: false,
		LeafSigner:              "",
		MaxUsesPerLeaf:          sdkmath.NewUint(1),
		ChallengeTrackerId:      challengeTrackerId,
		Uri:                     "",
		CustomData:              "",
	}
}

// GenerateETHSignatureChallenge generates an ETH signature challenge
func GenerateETHSignatureChallenge(challengeTrackerId, signer string) *types.ETHSignatureChallenge {
	return &types.ETHSignatureChallenge{
		Signer:             signer,
		ChallengeTrackerId: challengeTrackerId,
		Uri:                "",
		CustomData:         "",
	}
}

// GenerateTransfer generates a transfer message
func GenerateTransfer(from string, toAddresses []string, balances []*types.Balance) *types.Transfer {
	return &types.Transfer{
		From:                                    from,
		ToAddresses:                             toAddresses,
		Balances:                                balances,
		MerkleProofs:                            []*types.MerkleProof{},
		EthSignatureProofs:                      []*types.ETHSignatureProof{},
		Memo:                                    "",
		PrioritizedApprovals:                    []*types.ApprovalIdentifierDetails{},
		OnlyCheckPrioritizedCollectionApprovals: false,
		OnlyCheckPrioritizedIncomingApprovals:   false,
		OnlyCheckPrioritizedOutgoingApprovals:   false,
	}
}

// GenerateActionPermission generates an action permission
func GenerateActionPermission(permittedTimes, forbiddenTimes []*types.UintRange) *types.ActionPermission {
	return &types.ActionPermission{
		PermanentlyPermittedTimes: permittedTimes,
		PermanentlyForbiddenTimes: forbiddenTimes,
	}
}

// GenerateTokenIdsActionPermission generates a token IDs action permission
func GenerateTokenIdsActionPermission(tokenIds []*types.UintRange, permittedTimes, forbiddenTimes []*types.UintRange) *types.TokenIdsActionPermission {
	return &types.TokenIdsActionPermission{
		TokenIds:                  tokenIds,
		PermanentlyPermittedTimes: permittedTimes,
		PermanentlyForbiddenTimes: forbiddenTimes,
	}
}

// GenerateManagerTimeline generates a manager timeline entry
func GenerateManagerTimeline(timelineTimes []*types.UintRange, manager string) *types.ManagerTimeline {
	return &types.ManagerTimeline{
		TimelineTimes: timelineTimes,
		Manager:       manager,
	}
}

// GenerateCollectionMetadata generates collection metadata
func GenerateCollectionMetadata(uri, customData string) *types.CollectionMetadata {
	return &types.CollectionMetadata{
		Uri:        uri,
		CustomData: customData,
	}
}

// GenerateTokenMetadata generates token metadata
func GenerateTokenMetadata(tokenIds []*types.UintRange, uri, customData string) *types.TokenMetadata {
	return &types.TokenMetadata{
		TokenIds:   tokenIds,
		Uri:        uri,
		CustomData: customData,
	}
}

// GeneratePrioritizedApproval generates a prioritized approval identifier for collection-level approvals
// Approvals with merkle challenges, ETH signatures, or other non-auto-scannable criteria must be prioritized
func GeneratePrioritizedApproval(approvalId string) *types.ApprovalIdentifierDetails {
	return &types.ApprovalIdentifierDetails{
		ApprovalId:      approvalId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		Version:         sdkmath.NewUint(0),
	}
}
