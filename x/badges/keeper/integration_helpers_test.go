package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
)

func AssertUintsEqual(suite *TestSuite, expected sdkmath.Uint, actual sdkmath.Uint) {
	suite.Require().Equal(expected.Equal(actual), true, "Uints not equal %s %s", expected.String(), actual.String())
}

func AssertUintRangesEqual(suite *TestSuite, expected []*types.UintRange, actual []*types.UintRange) {
	remainingOne, _ := types.RemoveUintRangeFromUintRange(actual, expected)
	remainingTwo, _ := types.RemoveUintRangeFromUintRange(expected, actual)
	suite.Require().Equal(len(remainingOne), 0, "UintRanges not equal %s %s", expected, actual)
	suite.Require().Equal(len(remainingTwo), 0, "UintRanges not equal %s %s", expected, actual)
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

func GetFullUintRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:  sdkmath.NewUint(math.MaxUint64),
		},
	}
}

func GetBottomHalfUintRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:  sdkmath.NewUint(math.MaxUint32),
		},
	}
}

func GetTopHalfUintRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(math.MaxUint32 + 1),
			End:  sdkmath.NewUint(math.MaxUint64),
		},
	}
}

func GetOneUintRange() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:  sdkmath.NewUint(1),
		},
	}
}

func GetTwoUintRanges() []*types.UintRange {
	return []*types.UintRange{
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
						TimelineTimes: GetFullUintRanges(),
						ApprovedTransfers: []*types.CollectionApprovedTransfer{
						{
							ToMappingId: "All",
							FromMappingId: "All",
							InitiatedByMappingId: "All",
							TransferTimes: GetFullUintRanges(),
							BadgeIds: GetFullUintRanges(),
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
											OwnershipTimes: GetFullUintRanges(),
											BadgeIds: GetFullUintRanges(),
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
								TransferTimes: GetFullUintRanges(),
								BadgeIds: GetFullUintRanges(),
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
												OwnershipTimes: GetFullUintRanges(),
												BadgeIds: GetFullUintRanges(),
											},
										},
										NumTransfers: sdkmath.NewUint(1000),
									},
								},
							},
						},
						TimelineTimes: GetFullUintRanges(),
					},
				},
				DefaultApprovedOutgoingTransfersTimeline: []*types.UserApprovedOutgoingTransferTimeline{
					{
						ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
							{
								ToMappingId: "All",
								InitiatedByMappingId: "All",
								TransferTimes: GetFullUintRanges(),
								BadgeIds: GetFullUintRanges(),
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
												OwnershipTimes: GetFullUintRanges(),
												BadgeIds: GetFullUintRanges(),
											},
										},
										NumTransfers: sdkmath.NewUint(1000),
									},
								},
							},
						},
						TimelineTimes: GetFullUintRanges(),
					},
				},
				BadgesToCreate: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
								BadgeIds: GetFullUintRanges(),
								PermittedTimes: GetFullUintRanges(),
								ForbiddenTimes: []*types.UintRange{},
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
				TransferTimes: GetFullUintRanges(),
				BadgeIds: GetFullUintRanges(),
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
								OwnershipTimes: GetFullUintRanges(),
								BadgeIds: GetFullUintRanges(),
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
					BadgeIds: GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
		},
	}

	return collectionsToCreate
}