package keeper

import (
	"math"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// Little hack to make AllowedCombinations a 1-element array so we know if disallowed/allowed for ArbitraryValue
func ExpandCollectionApprovals(approvals []*types.CollectionApproval) []*types.CollectionApproval {
	newCurrApprovals := []*types.CollectionApproval{}
	//TODO: delete this; relic of old interface design
	for _, approval := range approvals {
		newCurrApprovals = append(newCurrApprovals, approval)
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
		}}, currApprovals...)

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
		}}, currApprovals...)

	return currApprovals
}
