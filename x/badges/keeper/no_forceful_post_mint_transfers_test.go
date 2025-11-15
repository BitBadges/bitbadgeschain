package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestIsOnlyMint() {
	// Test isOnlyMint logic by testing through validation
	// We'll test it indirectly by creating collections and checking behavior

	// Test case 1: Address list with only Mint should allow overrides
	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate1 := []*types.MsgNewCollection{
		{
			Creator: bob,
			CollectionApprovals: []*types.CollectionApproval{
				{
					ApprovalId:        "approval1",
					FromListId:        types.MintAddress, // Only Mint
					ToListId:          "All",
					InitiatedByListId: "All",
					TokenIds:          GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TransferTimes:     GetFullUintRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						OverridesFromOutgoingApprovals: true,
					},
				},
			},
			Invariants: &types.InvariantsAddObject{
				NoForcefulPostMintTransfers: true,
			},
		},
	}
	err := CreateCollections(suite, wctx, collectionsToCreate1)
	require.Nil(suite.T(), err, "Address list with only Mint should allow overrides")

	// Test case 2: Address list with Mint and another address should NOT allow overrides
	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
		ListId:    "mintandalice",
		Addresses: []string{types.MintAddress, alice},
		Whitelist: true,
		CreatedBy: bob,
	})
	require.Nil(suite.T(), err, "Should be able to create address list")

	collectionsToCreate2 := []*types.MsgNewCollection{
		{
			Creator: bob,
			CollectionApprovals: []*types.CollectionApproval{
				{
					ApprovalId:        "approval2",
					FromListId:        "mintandalice", // Includes Mint but not only Mint
					ToListId:          "All",
					InitiatedByListId: "All",
					TokenIds:          GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TransferTimes:     GetFullUintRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						OverridesFromOutgoingApprovals: true,
					},
				},
			},
			Invariants: &types.InvariantsAddObject{
				NoForcefulPostMintTransfers: true,
			},
		},
	}
	err = CreateCollections(suite, wctx, collectionsToCreate2)
	require.NotNil(suite.T(), err, "Address list with Mint and another address should NOT allow overrides")
}

func (suite *TestSuite) TestNoForcefulPostMintTransfers_WithOnlyMint() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with noForcefulPostMintTransfers enabled
	collectionsToCreate := []*types.MsgNewCollection{
		{
			Creator: bob,
			CollectionApprovals: []*types.CollectionApproval{
				{
					ApprovalId:        "approval1",
					FromListId:        types.MintAddress, // Only Mint
					ToListId:          "All",
					InitiatedByListId: "All",
					TokenIds:          GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TransferTimes:     GetFullUintRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						OverridesFromOutgoingApprovals: true, // This should be allowed since FromListId is only Mint
					},
				},
			},
			Invariants: &types.InvariantsAddObject{
				NoForcefulPostMintTransfers: true,
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Collection with only Mint FromListId and overrides should be allowed")
}

func (suite *TestSuite) TestNoForcefulPostMintTransfers_WithNonMint() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with noForcefulPostMintTransfers enabled
	collectionsToCreate := []*types.MsgNewCollection{
		{
			Creator: bob,
			CollectionApprovals: []*types.CollectionApproval{
				{
					ApprovalId:        "approval1",
					FromListId:        alice, // Not only Mint
					ToListId:          "All",
					InitiatedByListId: "All",
					TokenIds:          GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TransferTimes:     GetFullUintRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						OverridesFromOutgoingApprovals: true, // This should be disallowed
					},
				},
			},
			Invariants: &types.InvariantsAddObject{
				NoForcefulPostMintTransfers: true,
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NotNil(err, "Collection with non-Mint FromListId and overrides should be disallowed")
	suite.Require().Contains(err.Error(), "overridesFromOutgoingApprovals")
	suite.Require().Contains(err.Error(), "noForcefulPostMintTransfers")
}

func (suite *TestSuite) TestNoForcefulPostMintTransfers_WithOverridesToIncoming() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with noForcefulPostMintTransfers enabled
	collectionsToCreate := []*types.MsgNewCollection{
		{
			Creator: bob,
			CollectionApprovals: []*types.CollectionApproval{
				{
					ApprovalId:        "approval1",
					FromListId:        alice, // Not only Mint
					ToListId:          "All",
					InitiatedByListId: "All",
					TokenIds:          GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TransferTimes:     GetFullUintRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						OverridesToIncomingApprovals: true, // This should be disallowed
					},
				},
			},
			Invariants: &types.InvariantsAddObject{
				NoForcefulPostMintTransfers: true,
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NotNil(err, "Collection with non-Mint FromListId and overridesToIncomingApprovals should be disallowed")
	suite.Require().Contains(err.Error(), "overridesToIncomingApprovals")
	suite.Require().Contains(err.Error(), "noForcefulPostMintTransfers")
}

func (suite *TestSuite) TestNoForcefulPostMintTransfers_WithMintAndOtherAddresses() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create an address list that includes Mint but also other addresses
	err := suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
		ListId:    "mintandalice",
		Addresses: []string{types.MintAddress, alice},
		Whitelist: true,
		CreatedBy: bob,
	})
	suite.Require().Nil(err, "Should be able to create address list")

	// Create a collection with noForcefulPostMintTransfers enabled
	collectionsToCreate := []*types.MsgNewCollection{
		{
			Creator: bob,
			CollectionApprovals: []*types.CollectionApproval{
				{
					ApprovalId:        "approval1",
					FromListId:        "mintandalice", // Includes Mint but not only Mint
					ToListId:          "All",
					InitiatedByListId: "All",
					TokenIds:          GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TransferTimes:     GetFullUintRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						OverridesFromOutgoingApprovals: true, // This should be disallowed since it's not ONLY Mint
					},
				},
			},
			Invariants: &types.InvariantsAddObject{
				NoForcefulPostMintTransfers: true,
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NotNil(err, "Collection with FromListId that includes Mint but also other addresses should be disallowed")
	suite.Require().Contains(err.Error(), "overridesFromOutgoingApprovals")
}

func (suite *TestSuite) TestNoForcefulPostMintTransfers_Disabled() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with noForcefulPostMintTransfers disabled
	collectionsToCreate := []*types.MsgNewCollection{
		{
			Creator: bob,
			CollectionApprovals: []*types.CollectionApproval{
				{
					ApprovalId:        "approval1",
					FromListId:        alice,
					ToListId:          "All",
					InitiatedByListId: "All",
					TokenIds:          GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TransferTimes:     GetFullUintRanges(),
					ApprovalCriteria: &types.ApprovalCriteria{
						OverridesFromOutgoingApprovals: true, // This should be allowed when invariant is disabled
					},
				},
			},
			Invariants: &types.InvariantsAddObject{
				NoForcefulPostMintTransfers: false, // Disabled
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Collection with overrides should be allowed when noForcefulPostMintTransfers is disabled")
}
