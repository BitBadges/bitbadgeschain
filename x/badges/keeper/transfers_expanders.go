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
			badgeIds := types.GetUintRangesWithOptions(approvedTransfer.BadgeIds, allowedCombination.BadgeIdsOptions, true)
			ownedTimes := types.GetUintRangesWithOptions(approvedTransfer.OwnedTimes, allowedCombination.OwnedTimesOptions, true)
			times := types.GetUintRangesWithOptions(approvedTransfer.TransferTimes, allowedCombination.TransferTimesOptions, true)
			toMappingId := types.GetMappingIdWithOptions(approvedTransfer.ToMappingId, allowedCombination.ToMappingOptions, true)
			fromMappingId := types.GetMappingIdWithOptions(approvedTransfer.FromMappingId, allowedCombination.FromMappingOptions, true)
			initiatedByMappingId := types.GetMappingIdWithOptions(approvedTransfer.InitiatedByMappingId, allowedCombination.InitiatedByMappingOptions, true)
			
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
				ApprovalDetails: approvedTransfer.ApprovalDetails,
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
	})

	return currApprovedTransfers
}
