package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
)

func GetFullIdRanges() []*types.IdRange {
	return []*types.IdRange{
		{
			Start: sdkmath.NewUint(1),
			End:  sdkmath.NewUint(math.MaxUint64),
		},
	}
}

func GetBottomHalfIdRanges() []*types.IdRange {
	return []*types.IdRange{
		{
			Start: sdkmath.NewUint(1),
			End:  sdkmath.NewUint(math.MaxUint32),
		},
	}
}

func GetTopHalfIdRanges() []*types.IdRange {
	return []*types.IdRange{
		{
			Start: sdkmath.NewUint(math.MaxUint32 + 1),
			End:  sdkmath.NewUint(math.MaxUint64),
		},
	}
}

func GetOneIdRange() []*types.IdRange {
	return []*types.IdRange{
		{
			Start: sdkmath.NewUint(1),
			End:  sdkmath.NewUint(1),
		},
	}
}

func GetTwoIdRanges() []*types.IdRange {
	return []*types.IdRange{
		{
			Start: sdkmath.NewUint(2),
			End:  sdkmath.NewUint(2),
		},
	}
}

func GetCollectionsToCreate() []CollectionsToCreate {
	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				Creator: bob,
				BalancesType: sdkmath.NewUint(0),
				ApprovedTransfersTimeline: []*types.CollectionApprovedTransferTimeline{
					{
						Times: GetFullIdRanges(),
						ApprovedTransfers: []*types.CollectionApprovedTransfer{
						{
							ToMappingId: "All",
							FromMappingId: "All",
							InitiatedByMappingId: "All",
							TransferTimes: GetFullIdRanges(),
							BadgeIds: GetFullIdRanges(),
							AllowedCombinations: []*types.IsCollectionTransferAllowed{
								{
									IsAllowed: true,
								},
							},
							Challenges: []*types.Challenge{},
							TrackerId: "test",
							IncrementIdsBy: sdkmath.NewUint(0),
							IncrementTimesBy: sdkmath.NewUint(0),
							PerAddressApprovals: &types.PerAddressApprovals{
								ApprovalsPerFromAddress: &types.ApprovalsTracker{
									Amounts: []*types.Balance{
										{
											Amount: sdkmath.NewUint(1),
											Times: GetFullIdRanges(),
											BadgeIds: GetFullIdRanges(),
										},
									},
									NumTransfers: sdkmath.NewUint(1000),
								},
							},
						}},
					},
				},
				DefaultApprovedIncomingTransfersTimeline: []*types.UserApprovedIncomingTransferTimeline{
					{
						ApprovedIncomingTransfers: []*types.UserApprovedIncomingTransfer{
							{
								FromMappingId: "All",
								InitiatedByMappingId: "All",
								TransferTimes: GetFullIdRanges(),
								BadgeIds: GetFullIdRanges(),
								AllowedCombinations: []*types.IsUserIncomingTransferAllowed{
									{
										IsAllowed: true,
									},
								},
								Challenges: []*types.Challenge{},
								TrackerId: "test",
								IncrementIdsBy: sdkmath.NewUint(0),
								IncrementTimesBy: sdkmath.NewUint(0),
								PerAddressApprovals: &types.PerAddressApprovals{
									ApprovalsPerFromAddress: &types.ApprovalsTracker{
										Amounts: []*types.Balance{
											{
												Amount: sdkmath.NewUint(1),
												Times: GetFullIdRanges(),
												BadgeIds: GetFullIdRanges(),
											},
										},
										NumTransfers: sdkmath.NewUint(1000),
									},
								},
							},
						},
						Times: GetFullIdRanges(),
					},
				},
				DefaultApprovedOutgoingTransfersTimeline: []*types.UserApprovedOutgoingTransferTimeline{
					{
						ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
							{
								ToMappingId: "All",
								InitiatedByMappingId: "All",
								TransferTimes: GetFullIdRanges(),
								BadgeIds: GetFullIdRanges(),
								AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
									{
										IsAllowed: true,
									},
								},
								Challenges: []*types.Challenge{},
								TrackerId: "test",
								IncrementIdsBy: sdkmath.NewUint(0),
								IncrementTimesBy: sdkmath.NewUint(0),
								PerAddressApprovals: &types.PerAddressApprovals{
									ApprovalsPerFromAddress: &types.ApprovalsTracker{
										Amounts: []*types.Balance{
											{
												Amount: sdkmath.NewUint(1),
												Times: GetFullIdRanges(),
												BadgeIds: GetFullIdRanges(),
											},
										},
										NumTransfers: sdkmath.NewUint(1000),
									},
								},
							},
						},
						Times: GetFullIdRanges(),
					},
				},
				BadgesToCreate: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetFullIdRanges(),
						Times: GetFullIdRanges(),
					},
				},
				Permissions: &types.CollectionPermissions{
					CanCreateMoreBadges: []*types.ActionWithBadgeIdsAndTimesPermission{
						{
							DefaultValues: &types.ActionWithBadgeIdsAndTimesDefaultValues{
								BadgeIds: GetFullIdRanges(),
								PermittedTimes: GetFullIdRanges(),
								ForbiddenTimes: []*types.IdRange{},
							},
							Combinations: []*types.ActionWithBadgeIdsAndTimesCombination{{}},
						},
					},
				},
	
			},
			Amount:  sdkmath.NewUint(1),
		
		},
	}

	return collectionsToCreate
}


