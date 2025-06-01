package keeper

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
)

func AppendSelfInitiatedIncomingApproval(currApprovals []*types.UserIncomingApproval, userAddress string) []*types.UserIncomingApproval {
	currApprovals = append([]*types.UserIncomingApproval{
		{
			FromListId:        "AllWithMint", //everyone
			InitiatedByListId: userAddress,
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
			ApprovalId: "self-initiated-incoming",
			Version:    sdkmath.NewUint(0),
		}}, currApprovals...)

	return currApprovals
}

// By default, we approve all transfers if from === initiatedBy
func AppendSelfInitiatedOutgoingApproval(currApprovals []*types.UserOutgoingApproval, userAddress string) []*types.UserOutgoingApproval {
	//prepend it
	currApprovals = append([]*types.UserOutgoingApproval{
		{
			ToListId:          "AllWithMint", //everyone
			InitiatedByListId: userAddress,
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
			ApprovalId: "self-initiated-outgoing",
			Version:    sdkmath.NewUint(0),
		}}, currApprovals...)

	return currApprovals
}
