package keeper

import (
	"math"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// Little hack to make AllowedCombinations a 1-element array so we know if disallowed/allowed for ArbitraryValue
func ExpandCollectionApprovals(approvals []*types.CollectionApproval) []*types.CollectionApproval {
	newCurrApprovals := []*types.CollectionApproval{}
	for _, approval := range approvals {
		badgeIds := types.GetUintRangesWithOptions(approval.BadgeIds, approval.BadgeIdsOptions, true)
		ownershipTimes := types.GetUintRangesWithOptions(approval.OwnershipTimes, approval.OwnershipTimesOptions, true)
		times := types.GetUintRangesWithOptions(approval.TransferTimes, approval.TransferTimesOptions, true)
		toMappingId := types.GetMappingIdWithOptions(approval.ToMappingId, approval.ToMappingOptions, true)
		fromMappingId := types.GetMappingIdWithOptions(approval.FromMappingId, approval.FromMappingOptions, true)
		initiatedByMappingId := types.GetMappingIdWithOptions(approval.InitiatedByMappingId, approval.InitiatedByMappingOptions, true)
		
		newCurrApprovals = append(newCurrApprovals, &types.CollectionApproval{
			ToMappingId:          toMappingId,
			FromMappingId:        fromMappingId,
			InitiatedByMappingId: initiatedByMappingId,
			TransferTimes:        times,
			BadgeIds:             badgeIds,
			OwnershipTimes: 		 	ownershipTimes,
			IsApproved: approval.IsApproved,
			Uri: approval.Uri,
			CustomData: approval.CustomData,
			ApprovalCriteria: approval.ApprovalCriteria,
			ApprovalId: approval.ApprovalId,
			ApprovalTrackerId: approval.ApprovalTrackerId,
			ChallengeTrackerId: approval.ChallengeTrackerId,

			//Leave all options nil bc we applied them already
		})
	}

	return newCurrApprovals
}

// By default, we approve all transfers if to === initiatedBy
func AppendDefaultForIncoming(currApprovals []*types.UserIncomingApproval, userAddress string) []*types.UserIncomingApproval {
	currApprovals = append([]*types.UserIncomingApproval{
		{
		FromMappingId:        "AllWithMint", //everyone
		InitiatedByMappingId: userAddress,
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(uint64(math.MaxUint64)),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(uint64(math.MaxUint64)),
			},
		},
		BadgeIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(math.MaxUint64),
			},
		},
		IsApproved: true,
	}	}, currApprovals...)

	return currApprovals
}

// By default, we approve all transfers if from === initiatedBy
func AppendDefaultForOutgoing(currApprovals []*types.UserOutgoingApproval, userAddress string) []*types.UserOutgoingApproval {
	//prepend it 
	currApprovals = append([]*types.UserOutgoingApproval{
		{
		ToMappingId:          "AllWithMint", //everyone
		InitiatedByMappingId: userAddress,
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(uint64(math.MaxUint64)),
			},
		},
		OwnershipTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(uint64(math.MaxUint64)),
			},
		},
		BadgeIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(math.MaxUint64),
			},
		},
		IsApproved: true,
	}}, currApprovals...)


	return currApprovals
}
