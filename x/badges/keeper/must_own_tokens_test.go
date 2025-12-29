package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestMustOwnTokens() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnTokensMustSatisfyForAllAssets() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnTokensMustSatisfyForAllAssets2() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(2),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
		{
			CollectionId: sdkmath.NewUint(2),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(2),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnTokensMustOwnOne() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
		{
			CollectionId: sdkmath.NewUint(2),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	// Create second collection (will be collection ID 2) so bob owns tokens in it
	collectionsToCreate2 := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err = CreateCollections(suite, wctx, collectionsToCreate2)
	suite.Require().Nil(err, "Error creating second collection")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnTokensMustOwnOne2() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(2),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(2),
			},
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnTokensDoesntOwnBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnTokensMustOwnZero() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(0),
				End:   sdkmath.NewUint(0),
			},
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnTokensMustOwnGreaterThan() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(2),
				End:   sdkmath.NewUint(100),
			},
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestMustOwnTokensOwnershipCheckParty() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	// Test 1: Check ownership for initiator (default behavior)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:            GetFullUintRanges(),
			OwnershipTimes:      GetFullUintRanges(),
			OwnershipCheckParty: "initiator", // Explicitly set to initiator
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	// This should succeed because bob (initiator) owns the tokens
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token with initiator ownership check: %s")

	// Test 2: Check ownership for sender
	collectionsToCreate2 := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate2[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(2),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:            GetFullUintRanges(),
			OwnershipTimes:      GetFullUintRanges(),
			OwnershipCheckParty: "sender", // Check sender ownership
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate2)
	suite.Require().Nil(err)

	// This should succeed because bob (sender) owns the tokens
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(2),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(2)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token with sender ownership check: %s")

	// Test 3: Check ownership for recipient (should fail because alice doesn't own badges)
	collectionsToCreate3 := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate3[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(3),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:            GetFullUintRanges(),
			OwnershipTimes:      GetFullUintRanges(),
			OwnershipCheckParty: "recipient", // Check recipient ownership
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate3)
	suite.Require().Nil(err)

	// This should fail because alice (recipient) doesn't own the tokens
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(3),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(3)),
			},
		},
	})
	suite.Require().NotNil(err, "Transfer should fail when recipient doesn't own required tokens")
}

func (suite *TestSuite) TestMustOwnTokensBb1AddressSupport() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	// Test: Check ownership for arbitrary bb1 address (halt token scenario)
	// Create a halt token collection first (will be collection ID 1)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(charlie)
	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		{
			ApprovalId: "default",
			ApprovalCriteria: &types.ApprovalCriteria{},
			TransferTimes:     GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			OwnershipTimes:    GetFullUintRanges(),
			FromListId:        "Mint",
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	haltTokenCollectionId := sdkmath.NewUint(1) // First collection gets ID 1

	// Mint halt token to charlie
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      charlie,
		CollectionId: haltTokenCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, haltTokenCollectionId),
			},
		},
	})
	suite.Require().Nil(err, "Error minting halt token to charlie")

	// Create main collection (will be collection ID 2) with approval that requires halt token owner to own halt token
	collectionsToCreate2 := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate2[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: haltTokenCollectionId, // Halt token collection
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:            []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes:      GetFullUintRanges(),
			OwnershipCheckParty: charlie, // Use bb1 address directly - halt token owner
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate2)
	suite.Require().Nil(err)

	mainCollectionId := sdkmath.NewUint(2) // Second collection gets ID 2

	// Test 1: Transfer should succeed because charlie (halt token owner) owns the halt token
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, mainCollectionId),
			},
		},
	})
	suite.Require().Nil(err, "Transfer should succeed when halt token owner owns the halt token")

	// Test 2: Transfer halt token away from charlie (simulating halt)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      charlie,
		CollectionId: haltTokenCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        charlie,
				ToAddresses: []string{alice}, // Transfer halt token to alice
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, haltTokenCollectionId),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring halt token")

	// Test 3: Transfer should now fail because charlie (halt token owner) no longer owns the halt token
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, mainCollectionId),
			},
		},
	})
	suite.Require().NotNil(err, "Transfer should fail when halt token owner doesn't own the halt token")
	suite.Require().Contains(err.Error(), "token ownership requirement", "Error should mention token ownership requirement")
}

func (suite *TestSuite) Test2FAVaultEndToEnd() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	// This test demonstrates a complete 2FA vault setup with:
	// 1. Time-sensitive signatures (via time tokens with ownership times)
	// 2. Quick halting (via halt tokens with bb1 address support)

	// Step 1: Create time token collection (collection ID 1)
	timeTokenCollections := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	timeTokenCollections[0].CollectionApprovals = []*types.CollectionApproval{
		{
			ApprovalId: "default",
			ApprovalCriteria: &types.ApprovalCriteria{},
			TransferTimes:     GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			OwnershipTimes:    GetFullUintRanges(),
			FromListId:        "Mint",
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
		},
	}

	err = CreateCollections(suite, wctx, timeTokenCollections)
	suite.Require().Nil(err)
	timeTokenCollectionId := sdkmath.NewUint(1)

	// Step 2: Create halt token collection (collection ID 2)
	haltTokenCollections := GetTransferableCollectionToCreateAllMintedToCreator(charlie)
	haltTokenCollections[0].CollectionApprovals = []*types.CollectionApproval{
		{
			ApprovalId: "default",
			ApprovalCriteria: &types.ApprovalCriteria{},
			TransferTimes:     GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			OwnershipTimes:    GetFullUintRanges(),
			FromListId:        "Mint",
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
		},
	}

	err = CreateCollections(suite, wctx, haltTokenCollections)
	suite.Require().Nil(err)
	haltTokenCollectionId := sdkmath.NewUint(2)

	// Step 3: Mint halt token to charlie (halt token owner)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      charlie,
		CollectionId: haltTokenCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, haltTokenCollectionId),
			},
		},
	})
	suite.Require().Nil(err, "Error minting halt token to charlie")

	// Step 4: Create main collection (collection ID 3) with 2FA vault approval
	mainCollections := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	mainCollections[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		// Time token check (expires after TTL)
		{
			CollectionId: timeTokenCollectionId,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:            []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes:      GetFullUintRanges(),
			OverrideWithCurrentTime: true, // Always check current time
			OwnershipCheckParty: "initiator", // Check initiator owns time token
		},
		// Halt token check (transfer token = halt all approvals)
		{
			CollectionId: haltTokenCollectionId,
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:            []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
			OwnershipTimes:      GetFullUintRanges(),
			OwnershipCheckParty: charlie, // Use bb1 address directly - halt token owner
		},
	}

	err = CreateCollections(suite, wctx, mainCollections)
	suite.Require().Nil(err)
	mainCollectionId := sdkmath.NewUint(3)

	// Step 5: Mint time token to bob with ownership time = [now, now + 5 minutes]
	currentTime := suite.ctx.BlockTime().UnixMilli()
	ttl := int64(5 * 60 * 1000) // 5 minutes in milliseconds
	timeTokenOwnershipTimes := []*types.UintRange{
		{
			Start: sdkmath.NewUint(uint64(currentTime)),
			End:   sdkmath.NewUint(uint64(currentTime + ttl)),
		},
	}

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: timeTokenCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: timeTokenOwnershipTimes,
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, timeTokenCollectionId),
			},
		},
	})
	suite.Require().Nil(err, "Error minting time token to bob")

	// Step 6: Transfer should succeed - bob owns time token and charlie owns halt token
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, mainCollectionId),
			},
		},
	})
	suite.Require().Nil(err, "Transfer should succeed when both time token and halt token requirements are met")

	// Step 7: Simulate time token expiration by advancing block time
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(6 * 60 * 1000 * 1000000)) // Add 6 minutes (in nanoseconds)

	// Step 8: Transfer should fail - time token has expired
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, mainCollectionId),
			},
		},
	})
	suite.Require().NotNil(err, "Transfer should fail when time token has expired")
	suite.Require().Contains(err.Error(), "token ownership requirement", "Error should mention token ownership requirement")

	// Step 9: Mint new time token and reset block time
	currentTime = suite.ctx.BlockTime().UnixMilli()
	timeTokenOwnershipTimes = []*types.UintRange{
		{
			Start: sdkmath.NewUint(uint64(currentTime)),
			End:   sdkmath.NewUint(uint64(currentTime + ttl)),
		},
	}

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: timeTokenCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: timeTokenOwnershipTimes,
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, timeTokenCollectionId),
			},
		},
	})
	suite.Require().Nil(err, "Error minting new time token to bob")

	// Step 10: Transfer should succeed again with new time token
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, mainCollectionId),
			},
		},
	})
	suite.Require().Nil(err, "Transfer should succeed again with new time token")

	// Step 11: Emergency halt - transfer halt token away from charlie
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      charlie,
		CollectionId: haltTokenCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        charlie,
				ToAddresses: []string{alice}, // Transfer halt token to alice
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, haltTokenCollectionId),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring halt token")

	// Step 12: Transfer should now fail - halt token has been transferred (emergency halt)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, mainCollectionId),
			},
		},
	})
	suite.Require().NotNil(err, "Transfer should fail after emergency halt")
	suite.Require().Contains(err.Error(), "token ownership requirement", "Error should mention token ownership requirement")
}

// TestCheckMustOwnTokensBreakLogic tests the break/continue logic for MustSatisfyForAllAssets
func (suite *TestSuite) TestCheckMustOwnTokensBreakLogic() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create first collection (will be collection ID 1)
	collectionsToCreate1 := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(suite, wctx, collectionsToCreate1)
	suite.Require().Nil(err, "Error creating first collection")

	// Create second collection (will be collection ID 2)
	collectionsToCreate2 := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err = CreateCollections(suite, wctx, collectionsToCreate2)
	suite.Require().Nil(err, "Error creating second collection")

	// Test 1: MustSatisfyForAllAssets = false, first requirement satisfied - should continue to check second
	// Collection 1 exists and bob owns tokens
	mustOwnTokens1 := []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1), // bob owns this - should pass
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: false,
		},
		{
			CollectionId: sdkmath.NewUint(2), // bob owns this - should pass
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: false,
		},
	}
	approval1 := &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			MustOwnTokens: mustOwnTokens1,
		},
	}
	checkers := suite.app.BadgesKeeper.GetApprovalCriteriaCheckers(approval1)
	var detErrMsg string
	for _, checker := range checkers {
		detErrMsg, err = checker.Check(suite.ctx, approval1, nil, alice, bob, bob, "", "", nil, nil, "", false)
		if err != nil {
			break
		}
	}
	suite.Require().Nil(err, "Should succeed - both requirements satisfied, should continue through all")
	suite.Require().Equal("", detErrMsg, "Should have no error message")

	// Test 2: MustSatisfyForAllAssets = false, first requirement satisfied, second failed - should continue and fail
	mustOwnTokens2 := []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1), // bob owns this - should pass
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: false,
		},
		{
			CollectionId: sdkmath.NewUint(999), // doesn't exist - should fail
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: false,
		},
	}
	approval2 := &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			MustOwnTokens: mustOwnTokens2,
		},
	}
	checkers = suite.app.BadgesKeeper.GetApprovalCriteriaCheckers(approval2)
	detErrMsg = ""
	err = nil
	for _, checker := range checkers {
		detErrMsg, err = checker.Check(suite.ctx, approval2, nil, alice, bob, bob, "", "", nil, nil, "", false)
		if err != nil {
			break
		}
	}
	suite.Require().NotNil(err, "Should fail - second requirement failed, should continue and return error")
	suite.Require().Contains(detErrMsg, "token ownership requirement idx 1 failed", "Should have error message for second requirement")

	// Test 3: MustSatisfyForAllAssets = false, all requirements failed - should continue through all and return error
	mustOwnTokens3 := []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(999), // doesn't exist
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: false,
		},
		{
			CollectionId: sdkmath.NewUint(998), // doesn't exist
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: false,
		},
	}
	approval3 := &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			MustOwnTokens: mustOwnTokens3,
		},
	}
	checkers = suite.app.BadgesKeeper.GetApprovalCriteriaCheckers(approval3)
	detErrMsg = ""
	err = nil
	for _, checker := range checkers {
		detErrMsg, err = checker.Check(suite.ctx, approval3, nil, alice, bob, bob, "", "", nil, nil, "", false)
		if err != nil {
			break
		}
	}
	suite.Require().NotNil(err, "Should fail - all requirements failed")
	suite.Require().Contains(detErrMsg, "token ownership requirement idx 0 failed", "Should have error message for first requirement")

	// Test 4: MustSatisfyForAllAssets = true, first requirement failed - should break early
	mustOwnTokens4 := []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(999), // doesn't exist - should fail and break
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
		{
			CollectionId: sdkmath.NewUint(1), // bob owns this - should not be checked if first fails
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
	}
	approval4 := &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			MustOwnTokens: mustOwnTokens4,
		},
	}
	checkers = suite.app.BadgesKeeper.GetApprovalCriteriaCheckers(approval4)
	detErrMsg = ""
	err = nil
	for _, checker := range checkers {
		detErrMsg, err = checker.Check(suite.ctx, approval4, nil, alice, bob, bob, "", "", nil, nil, "", false)
		if err != nil {
			break
		}
	}
	suite.Require().NotNil(err, "Should fail - first requirement failed, should break early")
	suite.Require().Contains(detErrMsg, "token ownership requirement idx 0 failed", "Should have error message for first requirement")

	// Test 5: MustSatisfyForAllAssets = true, first satisfied, second failed - should break after second
	mustOwnTokens5 := []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1), // bob owns this - should pass
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
		{
			CollectionId: sdkmath.NewUint(999), // doesn't exist - should fail and break
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
	}
	approval5 := &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			MustOwnTokens: mustOwnTokens5,
		},
	}
	checkers = suite.app.BadgesKeeper.GetApprovalCriteriaCheckers(approval5)
	detErrMsg = ""
	err = nil
	for _, checker := range checkers {
		detErrMsg, err = checker.Check(suite.ctx, approval5, nil, alice, bob, bob, "", "", nil, nil, "", false)
		if err != nil {
			break
		}
	}
	suite.Require().NotNil(err, "Should fail - second requirement failed, should break after second")
	suite.Require().Contains(detErrMsg, "token ownership requirement idx 1 failed", "Should have error message for second requirement")

	// Test 6: MustSatisfyForAllAssets = false, first requirement fails amount check, second satisfied - should fail because all requirements must pass
	mustOwnTokens6 := []*types.MustOwnTokens{
		{
			CollectionId: sdkmath.NewUint(1), // bob owns 1 token
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(2), // requires 2, but bob only has 1 - should fail
				End:   sdkmath.NewUint(2),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: false,
		},
		{
			CollectionId: sdkmath.NewUint(2), // bob owns this - should pass
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: false,
		},
	}
	approval6 := &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			MustOwnTokens: mustOwnTokens6,
		},
	}
	checkers = suite.app.BadgesKeeper.GetApprovalCriteriaCheckers(approval6)
	detErrMsg = ""
	err = nil
	for _, checker := range checkers {
		detErrMsg, err = checker.Check(suite.ctx, approval6, nil, alice, bob, bob, "", "", nil, nil, "", false)
		if err != nil {
			break
		}
	}
	suite.Require().NotNil(err, "Should fail - first requirement failed, all requirements must pass")
	suite.Require().Contains(detErrMsg, "token ownership requirement idx 0 failed", "Should have error message for first requirement")
}
