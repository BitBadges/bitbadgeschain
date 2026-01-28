package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(2)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(3)),
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
	// Restrict TokensToCreate to only token ID 1 to match the approval
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	// Keep the mint-test approval and add the default approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{
		{
			ApprovalId:        "default",
			ApprovalCriteria:  &types.ApprovalCriteria{OverridesFromOutgoingApprovals: true},
			TransferTimes:     GetFullUintRanges(),
			TokenIds:          GetFullUintRanges(), // Cover all token IDs to allow collection creation
			OwnershipTimes:    GetFullUintRanges(),
			FromListId:        "Mint",
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
		},
	}, collectionsToCreate[0].CollectionApprovals...)

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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, haltTokenCollectionId),
			},
		},
	})
	suite.Require().Nil(err, "Error minting halt token to charlie")

	// Verify charlie owns the halt token and get the actual token ID
	haltTokenCollection, found := suite.app.TokenizationKeeper.GetCollectionFromStore(suite.ctx, haltTokenCollectionId)
	suite.Require().True(found, "Halt token collection should exist")
	charlieBalance, _, _ := suite.app.TokenizationKeeper.GetBalanceOrApplyDefault(suite.ctx, haltTokenCollection, charlie)
	suite.Require().True(len(charlieBalance.Balances) > 0, "Charlie should own the halt token")
	actualTokenId := charlieBalance.Balances[0].TokenIds[0].Start // Get the actual token ID that was minted

	// Create main collection (will be collection ID 2) with approval that requires halt token owner to own halt token
	collectionsToCreate2 := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate2[0].CollectionApprovals[1].ApprovalCriteria.MustOwnTokens = []*types.MustOwnTokens{
		{
			CollectionId: haltTokenCollectionId, // Halt token collection
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			TokenIds:            []*types.UintRange{{Start: actualTokenId, End: actualTokenId}}, // Use the actual token ID that was minted
			OwnershipTimes:      GetFullUintRanges(),
			OwnershipCheckParty: charlie, // Use bb1 address directly - halt token owner
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate2)
	suite.Require().Nil(err)

	mainCollectionId := sdkmath.NewUint(2) // Second collection gets ID 2

	// Get the test approval version
	testVersion, found := suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(mainCollectionId, "collection", "", "test"))
	suite.Require().True(found, "Test approval version should be found")
	mainPrioritizedApprovals := []*types.ApprovalIdentifierDetails{
		{
			ApprovalId:      "test",
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			Version:         testVersion,
		},
	}

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
				PrioritizedApprovals: mainPrioritizedApprovals,
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
						TokenIds:       []*types.UintRange{{Start: actualTokenId, End: actualTokenId}}, // Use the actual token ID
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, haltTokenCollectionId),
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
				PrioritizedApprovals: mainPrioritizedApprovals,
			},
		},
	})
	suite.Require().NotNil(err, "Transfer should fail when halt token owner doesn't own the halt token")
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
	checkers := suite.app.TokenizationKeeper.GetApprovalCriteriaCheckers(approval1)
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
	checkers = suite.app.TokenizationKeeper.GetApprovalCriteriaCheckers(approval2)
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
	checkers = suite.app.TokenizationKeeper.GetApprovalCriteriaCheckers(approval3)
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
	checkers = suite.app.TokenizationKeeper.GetApprovalCriteriaCheckers(approval4)
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
	checkers = suite.app.TokenizationKeeper.GetApprovalCriteriaCheckers(approval5)
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
	checkers = suite.app.TokenizationKeeper.GetApprovalCriteriaCheckers(approval6)
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
