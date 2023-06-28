package keeper

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ExpandCollectionApprovedTransfers(approvedTransfers []*types.CollectionApprovedTransfer) []*types.CollectionApprovedTransfer {
	newCurrApprovedTransfers := []*types.CollectionApprovedTransfer{}
	for _, approvedTransfer := range approvedTransfers {
		for _, allowedCombination := range approvedTransfer.AllowedCombinations {
			badgeIds := approvedTransfer.BadgeIds
			if allowedCombination.InvertBadgeIds {
				badgeIds = types.InvertIdRanges(badgeIds, sdk.NewUint(math.MaxUint64))
			}

			times := approvedTransfer.TransferTimes
			if allowedCombination.InvertTransferTimes {
				times = types.InvertIdRanges(times, sdk.NewUint(math.MaxUint64))
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
				ToMappingId: toMappingId,
				FromMappingId: fromMappingId,
				InitiatedByMappingId: initiatedByMappingId,
				TransferTimes: times,
				BadgeIds: badgeIds,
				AllowedCombinations: []*types.IsCollectionTransferAllowed{
					{
						IsAllowed: allowedCombination.IsAllowed,
					},
				},
				Approvals: approvedTransfer.Approvals,
				PerAddressApprovals: approvedTransfer.PerAddressApprovals,
				MaxNumTransfers: approvedTransfer.MaxNumTransfers,
				PerAddressMaxNumTransfers: approvedTransfer.PerAddressMaxNumTransfers,
				RequireFromEqualsInitiatedBy: approvedTransfer.RequireFromEqualsInitiatedBy,
				RequireFromDoesNotEqualInitiatedBy: approvedTransfer.RequireFromDoesNotEqualInitiatedBy,
				RequireToEqualsInitiatedBy: approvedTransfer.RequireToEqualsInitiatedBy,
				RequireToDoesNotEqualInitiatedBy: approvedTransfer.RequireToDoesNotEqualInitiatedBy,
				OverridesFromApprovedOutgoingTransfers: approvedTransfer.OverridesFromApprovedOutgoingTransfers,
				OverridesToApprovedIncomingTransfers: approvedTransfer.OverridesToApprovedIncomingTransfers,
				CustomData: approvedTransfer.CustomData,
				Uri: approvedTransfer.Uri,
				TrackerId: approvedTransfer.TrackerId,
				Challenges: approvedTransfer.Challenges,
			})
		}
	}

	return newCurrApprovedTransfers
}

func ExpandUserApprovedIncomingTransfers(currApprovedTransfers []*types.UserApprovedIncomingTransfer, userAddress string) []*types.UserApprovedIncomingTransfer {
	newCurrApprovedTransfers := []*types.UserApprovedIncomingTransfer{}
	for _, approvedTransfer := range currApprovedTransfers {
		for _, allowedCombination := range approvedTransfer.AllowedCombinations {
			badgeIds := approvedTransfer.BadgeIds
			if allowedCombination.InvertBadgeIds {
				badgeIds = types.InvertIdRanges(badgeIds, sdk.NewUint(math.MaxUint64))
			}

			times := approvedTransfer.TransferTimes
			if allowedCombination.InvertTransferTimes {
				times = types.InvertIdRanges(times, sdk.NewUint(math.MaxUint64))
			}

			fromMappingId := approvedTransfer.FromMappingId
			if allowedCombination.InvertFrom {
				fromMappingId = "!" + fromMappingId
			}

			initiatedByMappingId := approvedTransfer.InitiatedByMappingId
			if allowedCombination.InvertInitiatedBy {
				initiatedByMappingId = "!" + initiatedByMappingId
			}

			newCurrApprovedTransfers = append(newCurrApprovedTransfers, &types.UserApprovedIncomingTransfer{
				FromMappingId: fromMappingId,
				InitiatedByMappingId: initiatedByMappingId,
				TransferTimes: times,
				BadgeIds: badgeIds,
				AllowedCombinations: []*types.IsUserIncomingTransferAllowed{
					{
						IsAllowed: allowedCombination.IsAllowed,
					},
				},
				Approvals: approvedTransfer.Approvals,
				PerAddressApprovals: approvedTransfer.PerAddressApprovals,
				MaxNumTransfers: approvedTransfer.MaxNumTransfers,
				PerAddressMaxNumTransfers: approvedTransfer.PerAddressMaxNumTransfers,
				Challenges: approvedTransfer.Challenges,
				RequireFromEqualsInitiatedBy: approvedTransfer.RequireFromEqualsInitiatedBy,
				RequireFromDoesNotEqualInitiatedBy: approvedTransfer.RequireFromDoesNotEqualInitiatedBy,
				CustomData: approvedTransfer.CustomData,
				Uri: approvedTransfer.Uri,
				TrackerId: approvedTransfer.TrackerId,
			})
		}
	}


	newCurrApprovedTransfers = append(newCurrApprovedTransfers, &types.UserApprovedIncomingTransfer{
		FromMappingId: "All", //everyone
		InitiatedByMappingId: userAddress,
		TransferTimes: []*types.IdRange{
			{
				Start: sdk.NewUint(0),
				End: sdk.NewUint(uint64(math.MaxUint64)),
			},
		},
		BadgeIds: []*types.IdRange{
			{
				Start: sdk.NewUint(1),
				End: sdk.NewUint(math.MaxUint64),
			},
		},
		AllowedCombinations: []*types.IsUserIncomingTransferAllowed{
			{
				IsAllowed: true,
			},
		},
	})

	return newCurrApprovedTransfers
}

func ExpandUserApprovedOutgoingTransfers(currApprovedTransfers []*types.UserApprovedOutgoingTransfer, address string) []*types.UserApprovedOutgoingTransfer {
	newCurrApprovedTransfers := []*types.UserApprovedOutgoingTransfer{}
	for _, approvedTransfer := range currApprovedTransfers {
		for _, allowedCombination := range approvedTransfer.AllowedCombinations {
			badgeIds := approvedTransfer.BadgeIds
			if allowedCombination.InvertBadgeIds {
				badgeIds = types.InvertIdRanges(badgeIds, sdk.NewUint(math.MaxUint64))
			}

			times := approvedTransfer.TransferTimes
			if allowedCombination.InvertTransferTimes {
				times = types.InvertIdRanges(times, sdk.NewUint(uint64(math.MaxUint64)))
			}

			toMappingId := approvedTransfer.ToMappingId
			if allowedCombination.InvertTo {
				toMappingId = "!" + toMappingId
			}

			initiatedByMappingId := approvedTransfer.InitiatedByMappingId
			if allowedCombination.InvertInitiatedBy {
				initiatedByMappingId = "!" + initiatedByMappingId
			}

			newCurrApprovedTransfers = append(newCurrApprovedTransfers, &types.UserApprovedOutgoingTransfer{
				ToMappingId: toMappingId,
				InitiatedByMappingId: initiatedByMappingId,
				TransferTimes: times,
				BadgeIds: badgeIds,
				AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
					{
						IsAllowed: allowedCombination.IsAllowed,
					},
				},
				Approvals: approvedTransfer.Approvals,
				PerAddressApprovals: approvedTransfer.PerAddressApprovals,
				MaxNumTransfers: approvedTransfer.MaxNumTransfers,
				PerAddressMaxNumTransfers: approvedTransfer.PerAddressMaxNumTransfers,
				Challenges: approvedTransfer.Challenges,
				RequireToEqualsInitiatedBy: approvedTransfer.RequireToEqualsInitiatedBy,
				RequireToDoesNotEqualInitiatedBy: approvedTransfer.RequireToDoesNotEqualInitiatedBy,
				CustomData: approvedTransfer.CustomData,
				Uri: approvedTransfer.Uri,
				TrackerId: approvedTransfer.TrackerId,
			})
		}
	}


	newCurrApprovedTransfers = append(newCurrApprovedTransfers, &types.UserApprovedOutgoingTransfer{
		ToMappingId: "All",
		InitiatedByMappingId: address,
		TransferTimes: []*types.IdRange{
			{
				Start: sdk.NewUint(0),
				End: sdk.NewUint(uint64(math.MaxUint64)),
			},
		},
		BadgeIds: []*types.IdRange{
			{
				Start: sdk.NewUint(1),
				End: sdk.NewUint(math.MaxUint64),
			},
		},
		AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
			{
				IsAllowed: true,
			},
		},
	})

	return newCurrApprovedTransfers
}