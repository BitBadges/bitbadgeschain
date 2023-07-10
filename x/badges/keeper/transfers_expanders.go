package keeper

import (
	"math"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// Little hack to make AllowedCombinations a 1-element array so we know if disallowed/allowed for ArbitraryValue
func ExpandCollectionApprovedTransfers(approvedTransfers []*types.CollectionApprovedTransfer) []*types.CollectionApprovedTransfer {
	newCurrApprovedTransfers := []*types.CollectionApprovedTransfer{}
	for _, approvedTransfer := range approvedTransfers {
		for _, allowedCombination := range approvedTransfer.AllowedCombinations {
			badgeIds := approvedTransfer.BadgeIds
			if allowedCombination.InvertBadgeIds {
				badgeIds = types.InvertUintRanges(badgeIds, sdkmath.NewUint(math.MaxUint64))
			}

			ownedTimes := approvedTransfer.OwnedTimes
			if allowedCombination.InvertOwnedTimes {
				ownedTimes = types.InvertUintRanges(ownedTimes, sdkmath.NewUint(math.MaxUint64))
			}

			times := approvedTransfer.TransferTimes
			if allowedCombination.InvertTransferTimes {
				times = types.InvertUintRanges(times, sdkmath.NewUint(math.MaxUint64))
			}

			toMappingId := approvedTransfer.ToMappingId
			if allowedCombination.InvertTo {
				toMappingId = "!" + toMappingId
			}

			fromMappingId := approvedTransfer.FromMappingId
			if allowedCombination.InvertFrom {
				fromMappingId = "!" + fromMappingId
			}

			initiatedByMappingId := approvedTransfer.InitiatedByMappingId
			if allowedCombination.InvertInitiatedBy {
				initiatedByMappingId = "!" + initiatedByMappingId
			}

			newCurrApprovedTransfers = append(newCurrApprovedTransfers, &types.CollectionApprovedTransfer{
				ToMappingId:          toMappingId,
				FromMappingId:        fromMappingId,
				InitiatedByMappingId: initiatedByMappingId,
				TransferTimes:        times,
				BadgeIds:             badgeIds,
				OwnedTimes: 		 	ownedTimes,
				AllowedCombinations: []*types.IsCollectionTransferAllowed{
					{
						IsAllowed: allowedCombination.IsAllowed,
					},
				},
				PredeterminedBalances: 				  			approvedTransfer.PredeterminedBalances,
				ApprovalAmounts: 								  			approvedTransfer.ApprovalAmounts,
				MaxNumTransfers: 								  			approvedTransfer.MaxNumTransfers,
				RequireFromEqualsInitiatedBy:           approvedTransfer.RequireFromEqualsInitiatedBy,
				RequireFromDoesNotEqualInitiatedBy:     approvedTransfer.RequireFromDoesNotEqualInitiatedBy,
				RequireToEqualsInitiatedBy:             approvedTransfer.RequireToEqualsInitiatedBy,
				RequireToDoesNotEqualInitiatedBy:       approvedTransfer.RequireToDoesNotEqualInitiatedBy,
				OverridesFromApprovedOutgoingTransfers: approvedTransfer.OverridesFromApprovedOutgoingTransfers,
				OverridesToApprovedIncomingTransfers:   approvedTransfer.OverridesToApprovedIncomingTransfers,
				CustomData:                             approvedTransfer.CustomData,
				Uri:                                    approvedTransfer.Uri,
				ApprovalId:                              approvedTransfer.ApprovalId,
				Challenges:                             approvedTransfer.Challenges,
			})
		}
	}

	return newCurrApprovedTransfers
}

// By default, we approve all transfers if to === initiatedBy
func AppendDefaultForIncoming(currApprovedTransfers []*types.UserApprovedIncomingTransfer, userAddress string) []*types.UserApprovedIncomingTransfer {
	currApprovedTransfers = append(currApprovedTransfers, &types.UserApprovedIncomingTransfer{
		FromMappingId:        "All", //everyone
		InitiatedByMappingId: userAddress,
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(uint64(math.MaxUint64)),
			},
		},
		OwnedTimes: []*types.UintRange{
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
		AllowedCombinations: []*types.IsUserIncomingTransferAllowed{
			{
				IsAllowed: true,
			},
		},
		ApprovalAmounts: &types.ApprovalAmounts{},
		MaxNumTransfers: &types.MaxNumTransfers{},
	})

	return currApprovedTransfers
}

// By default, we approve all transfers if from === initiatedBy
func AppendDefaultForOutgoing(currApprovedTransfers []*types.UserApprovedOutgoingTransfer, userAddress string) []*types.UserApprovedOutgoingTransfer {
	currApprovedTransfers = append(currApprovedTransfers, &types.UserApprovedOutgoingTransfer{
		ToMappingId:          "All", //everyone
		InitiatedByMappingId: userAddress,
		TransferTimes: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(uint64(math.MaxUint64)),
			},
		},
		OwnedTimes: []*types.UintRange{
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
		AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
			{
				IsAllowed: true,
			},
		},
		ApprovalAmounts: &types.ApprovalAmounts{},
		MaxNumTransfers: &types.MaxNumTransfers{},
	})

	return currApprovedTransfers
}
