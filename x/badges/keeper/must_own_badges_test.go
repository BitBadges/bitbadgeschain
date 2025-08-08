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
	suite.Require().Nil(err, "Error transferring badge: %s")
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
	suite.Require().Nil(err, "Error transferring badge: %s")
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
	suite.Require().Error(err, "Error transferring badge: %s")
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
	suite.Require().Nil(err, "Error transferring badge: %s")
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
	suite.Require().Error(err, "Error transferring badge: %s")
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
	suite.Require().Error(err, "Error transferring badge: %s")
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
	suite.Require().Error(err, "Error transferring badge: %s")

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
	suite.Require().Nil(err, "Error transferring badge: %s")
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
	suite.Require().Error(err, "Error transferring badge: %s")

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
	suite.Require().Error(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestMustOwnBadgesOwnershipCheckParty() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	// Test 1: Check ownership for initiator (default behavior)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:            GetFullUintRanges(),
			OwnershipTimes:      GetFullUintRanges(),
			OwnershipCheckParty: "initiator", // Explicitly set to initiator
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	// This should succeed because bob (initiator) owns the badges
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
	suite.Require().Nil(err, "Error transferring badge with initiator ownership check: %s")

	// Test 2: Check ownership for sender
	collectionsToCreate2 := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate2[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(2),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:            GetFullUintRanges(),
			OwnershipTimes:      GetFullUintRanges(),
			OwnershipCheckParty: "sender", // Check sender ownership
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate2)
	suite.Require().Nil(err)

	// This should succeed because bob (sender) owns the badges
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(2),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(2)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge with sender ownership check: %s")

	// Test 3: Check ownership for recipient (should fail because alice doesn't own badges)
	collectionsToCreate3 := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate3[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(3),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:            GetFullUintRanges(),
			OwnershipTimes:      GetFullUintRanges(),
			OwnershipCheckParty: "recipient", // Check recipient ownership
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate3)
	suite.Require().Nil(err)

	// This should fail because alice (recipient) doesn't own the badges
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(3),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(3)),
			},
		},
	})
	suite.Require().NotNil(err, "Transfer should fail when recipient doesn't own required badges")
}
