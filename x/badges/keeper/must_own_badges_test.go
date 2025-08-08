package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestMustOwnBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnBadgesMustSatisfyForAllAssets() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnBadgesMustSatisfyForAllAssets2() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(2),
			},
			BadgeIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
		{
			CollectionId: sdkmath.NewUint(2),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(2),
			},
			BadgeIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnBadgesMustOwnOne() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
		{
			CollectionId: sdkmath.NewUint(2),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnBadgesMustOwnOne2() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(2),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(2),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnBadgesDoesntOwnBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnBadgesMustOwnZero() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(0),
				End:   sdkmath.NewUint(0),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnBadgesMustOwnGreaterThan() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(2),
				End:   sdkmath.NewUint(100),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")
}
