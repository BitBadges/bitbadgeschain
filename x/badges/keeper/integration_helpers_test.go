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
		actual, err = types.SubtractBalance(actual, balance, false)
		suite.Require().Nil(err, "Underflow error comparing balances: %s")
	}

	suite.Require().Equal(len(actual), 0, "Balances not equal %s %s", expected, actual)
}

func GetFullUintRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(math.MaxUint64),
		},
	}
}

func GetBottomHalfUintRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(math.MaxUint32),
		},
	}
}

func GetTopHalfUintRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(math.MaxUint32 + 1),
			End:   sdkmath.NewUint(math.MaxUint64),
		},
	}
}

func GetOneUintRange() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(1),
		},
	}
}

func GetTwoUintRanges() []*types.UintRange {
	return []*types.UintRange{
		{
			Start: sdkmath.NewUint(2),
			End:   sdkmath.NewUint(2),
		},
	}
}

func GetCollectionsToCreate() []*types.MsgNewCollection {
	collectionsToCreate := []*types.MsgNewCollection{
		{
			Creator:      bob,
			BalancesType: sdkmath.NewUint(1),
			CollectionApprovedTransfers:  []*types.CollectionApprovedTransfer{
				{
					ToMappingId:          "AllWithoutMint",
					FromMappingId:        "AllWithoutMint",
					InitiatedByMappingId: "AllWithoutMint",
					TransferTimes:        GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
					BadgeIds:             GetFullUintRanges(),
					AllowedCombinations: []*types.IsCollectionTransferAllowed{
						{
							IsApproved: true,
						},
					},
					ApprovalId: "test",
					ApprovalTrackerId:                 "test",
					ApprovalDetails: &types.ApprovalDetails{
							MaxNumTransfers: &types.MaxNumTransfers{
								OverallMaxNumTransfers: sdkmath.NewUint(1000),
							},
							ApprovalAmounts: &types.ApprovalAmounts{
								PerFromAddressApprovalAmount: sdkmath.NewUint(1),
							},
						
					},
				}},
			DefaultApprovedIncomingTransfers:[]*types.UserApprovedIncomingTransfer{
						{
							FromMappingId:        "AllWithoutMint",
							InitiatedByMappingId: "AllWithoutMint",
							TransferTimes:        GetFullUintRanges(),
							OwnershipTimes: 			GetFullUintRanges(),
							BadgeIds:             GetFullUintRanges(),
							AllowedCombinations: []*types.IsUserIncomingTransferAllowed{
								{
									IsApproved: true,
								},
							},
							ApprovalId: "test",
							ApprovalTrackerId:                 "test",
							ApprovalDetails: &types.IncomingApprovalDetails{
								
									
									MaxNumTransfers: &types.MaxNumTransfers{
										OverallMaxNumTransfers: sdkmath.NewUint(1000),
									},
									ApprovalAmounts: &types.ApprovalAmounts{
										PerFromAddressApprovalAmount: sdkmath.NewUint(1),
									},
								
							},
						},
					
			},
			DefaultApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
						{
							ToMappingId:          "AllWithoutMint",
							InitiatedByMappingId: "AllWithoutMint",
							TransferTimes:        GetFullUintRanges(),
							OwnershipTimes: GetFullUintRanges(),
							BadgeIds:             GetFullUintRanges(),
							AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
								{
									IsApproved: true,
								},
							},
							ApprovalId: "test",
							ApprovalTrackerId:                 		"test",
							ApprovalDetails: &types.OutgoingApprovalDetails{
								MaxNumTransfers: &types.MaxNumTransfers{
									OverallMaxNumTransfers: sdkmath.NewUint(1000),
								},
								ApprovalAmounts: &types.ApprovalAmounts{
									PerFromAddressApprovalAmount: sdkmath.NewUint(1),
								},
							},
						},
					
			},
			BadgesToCreate: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			Permissions: &types.CollectionPermissions{
				CanArchiveCollection:                           []*types.TimedUpdatePermission{},
				CanUpdateContractAddress:             []*types.TimedUpdatePermission{},
				CanUpdateOffChainBalancesMetadata:    []*types.TimedUpdatePermission{},
				CanUpdateStandards:                   []*types.TimedUpdatePermission{},
				CanUpdateCustomData:                  []*types.TimedUpdatePermission{},
				CanDeleteCollection:                  []*types.ActionPermission{},
				CanUpdateManager:                     []*types.TimedUpdatePermission{},
				CanUpdateCollectionMetadata:          []*types.TimedUpdatePermission{},
				CanUpdateBadgeMetadata:               []*types.TimedUpdateWithBadgeIdsPermission{},
				CanUpdateCollectionApprovedTransfers: []*types.CollectionApprovedTransferPermission{},
				CanCreateMoreBadges: []*types.BalancesActionPermission{
					{
						DefaultValues: &types.BalancesActionDefaultValues{
							// BadgeIds:       GetFullUintRanges(),
							// OwnershipTimes:     GetFullUintRanges(),
							PermittedTimes: GetFullUintRanges(),
							ForbiddenTimes: []*types.UintRange{},
						},
						Combinations: []*types.BalancesActionCombination{{}},
					},
				},
			},
		},
	}

	return collectionsToCreate
}

func GetTransferableCollectionToCreateAllMintedToCreator(creator string) []*types.MsgNewCollection {
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfers[0].ApprovalDetails.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(uint64(math.MaxUint64))
	collectionsToCreate[0].CollectionApprovedTransfers = append([]*types.CollectionApprovedTransfer{{

		ToMappingId:                            "AllWithoutMint",
		FromMappingId:                          "Mint",
		InitiatedByMappingId:                   "AllWithoutMint",
		
		TransferTimes:                          GetFullUintRanges(),
		BadgeIds:                               GetFullUintRanges(),
		OwnershipTimes: GetFullUintRanges(),
		AllowedCombinations: []*types.IsCollectionTransferAllowed{
			{
				IsApproved: true,
			},
		},
		ApprovalId: "mint-test",
		ApprovalTrackerId:                 "mint-test",
		ApprovalDetails: &types.ApprovalDetails{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(1000),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				},
				OverridesFromApprovedOutgoingTransfers: true,
				OverridesToApprovedIncomingTransfers:   true,
		},
	},
	}, collectionsToCreate[0].CollectionApprovedTransfers...,
	)

	collectionsToCreate[0].Transfers = []*types.Transfer{
		{
			From:        "Mint",
			ToAddresses: []string{creator},
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
		},
	}

	return collectionsToCreate
}
