package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
)

func AssertUintsEqual(suite *TestSuite, expected sdkmath.Uint, actual sdkmath.Uint) {
	suite.Require().Equal(expected.Equal(actual), true, "Uints not equal %s %s", expected.String(), actual.String())
}

func AssertIdRangesEqual(suite *TestSuite, expected []*types.IdRange, actual []*types.IdRange) {
	remainingOne, _ := types.RemoveIdRangeFromIdRange(actual, expected)
	remainingTwo, _ := types.RemoveIdRangeFromIdRange(expected, actual)
	suite.Require().Equal(len(remainingOne), 0, "IdRanges not equal %s %s", expected, actual)
	suite.Require().Equal(len(remainingTwo), 0, "IdRanges not equal %s %s", expected, actual)
}

func AssertBalancesEqual(suite *TestSuite, expected []*types.Balance, actual []*types.Balance) {
	// suite.Require().Equal(len(expected), len(actual), "Balances length not equal %d %d", len(expected), len(actual))

	err := *new(error)
	for _, balance := range expected {
		actual, err = types.SubtractBalance(actual, balance)
		suite.Require().Nil(err, "Underflow error comparing balances: %s")
	}

	suite.Require().Equal(len(actual), 0, "Balances not equal %s %s", expected, actual)
}

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
				CollectionApprovedTransfersTimeline: []*types.CollectionApprovedTransferTimeline{
					{
						TimelineTimes: GetFullIdRanges(),
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
							IncrementBadgeIdsBy: sdkmath.NewUint(0),
							IncrementOwnershipTimesBy: sdkmath.NewUint(0),
							PerAddressApprovals: &types.PerAddressApprovals{
								ApprovalsPerFromAddress: &types.ApprovalsTracker{
									Amounts: []*types.Balance{
										{
											Amount: sdkmath.NewUint(1),
											OwnershipTimes: GetFullIdRanges(),
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
								IncrementBadgeIdsBy: sdkmath.NewUint(0),
								IncrementOwnershipTimesBy: sdkmath.NewUint(0),
								PerAddressApprovals: &types.PerAddressApprovals{
									ApprovalsPerFromAddress: &types.ApprovalsTracker{
										Amounts: []*types.Balance{
											{
												Amount: sdkmath.NewUint(1),
												OwnershipTimes: GetFullIdRanges(),
												BadgeIds: GetFullIdRanges(),
											},
										},
										NumTransfers: sdkmath.NewUint(1000),
									},
								},
							},
						},
						TimelineTimes: GetFullIdRanges(),
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
								IncrementBadgeIdsBy: sdkmath.NewUint(0),
								IncrementOwnershipTimesBy: sdkmath.NewUint(0),
								PerAddressApprovals: &types.PerAddressApprovals{
									ApprovalsPerFromAddress: &types.ApprovalsTracker{
										Amounts: []*types.Balance{
											{
												Amount: sdkmath.NewUint(1),
												OwnershipTimes: GetFullIdRanges(),
												BadgeIds: GetFullIdRanges(),
											},
										},
										NumTransfers: sdkmath.NewUint(1000),
									},
								},
							},
						},
						TimelineTimes: GetFullIdRanges(),
					},
				},
				BadgesToCreate: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetFullIdRanges(),
						OwnershipTimes: GetFullIdRanges(),
					},
				},
				Permissions: &types.CollectionPermissions{
					CanArchive: []*types.TimedUpdatePermission{},
					CanUpdateContractAddress: []*types.TimedUpdatePermission{},
					CanUpdateOffChainBalancesMetadata: []*types.TimedUpdatePermission{},
					CanUpdateStandards: []*types.TimedUpdatePermission{},
					CanUpdateCustomData: []*types.TimedUpdatePermission{},
					CanDeleteCollection: []*types.ActionPermission{},
					CanUpdateManager: []*types.TimedUpdatePermission{},
					CanUpdateCollectionMetadata: []*types.TimedUpdatePermission{},
					CanUpdateBadgeMetadata: []*types.TimedUpdateWithBadgeIdsPermission{},
					CanUpdateInheritedBalances: []*types.TimedUpdateWithBadgeIdsPermission{},
					CanUpdateCollectionApprovedTransfers: []*types.CollectionApprovedTransferPermission{},
					CanCreateMoreBadges: []*types.BalancesActionPermission{
						{
							DefaultValues: &types.BalancesActionDefaultValues{
								BadgeIds: GetFullIdRanges(),
								PermittedTimes: GetFullIdRanges(),
								ForbiddenTimes: []*types.IdRange{},
							},
							Combinations: []*types.BalancesActionCombination{{

							}},
						},
					},
					
				},
	
			},
			Amount:  sdkmath.NewUint(1),
		
		},
	}

	return collectionsToCreate
}


func GetTransferableCollectionToCreateAllMintedToCreator(creator string) []CollectionsToCreate {
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].PerAddressApprovals.ApprovalsPerFromAddress.Amounts[0].Amount = sdkmath.NewUint(uint64(math.MaxUint64))
	collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline[0].ApprovedTransfers = append([]*types.CollectionApprovedTransfer{{
		
				ToMappingId: "All",
				FromMappingId: "Mint",
				InitiatedByMappingId: "All",
				OverridesFromApprovedOutgoingTransfers: true,
				OverridesToApprovedIncomingTransfers: true,
				TransferTimes: GetFullIdRanges(),
				BadgeIds: GetFullIdRanges(),
				AllowedCombinations: []*types.IsCollectionTransferAllowed{
					{
						IsAllowed: true,
					},
				},
				Challenges: []*types.Challenge{},
				TrackerId: "test",
				IncrementBadgeIdsBy: sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				PerAddressApprovals: &types.PerAddressApprovals{
					ApprovalsPerFromAddress: &types.ApprovalsTracker{
						Amounts: []*types.Balance{
							{
								Amount: sdkmath.NewUint(1000),
								OwnershipTimes: GetFullIdRanges(),
								BadgeIds: GetFullIdRanges(),
							},
						},
						NumTransfers: sdkmath.NewUint(1000),
					},
				},
			},
		}, collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline[0].ApprovedTransfers...
	)

	collectionsToCreate[0].Collection.Transfers = []*types.Transfer{
		{
			From: "Mint",
			ToAddresses: []string{creator},
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: GetFullIdRanges(),
					OwnershipTimes: GetFullIdRanges(),
				},
			},
		},
	}

	return collectionsToCreate
}