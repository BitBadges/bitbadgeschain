package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// createBackedCollection creates a collection with cosmosCoinBackedPath and returns the collection.
// It creates the collection without any allowBackedMinting approvals so that tests can add them later.
func (suite *TestSuite) createBackedCollection() *types.TokenCollection {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	// Remove mint approvals since they're not allowed when cosmosCoinBackedPath is set
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{
					Amount: sdkmath.NewUint(1),
					Denom:  "ibc/1234567890ABCDEF",
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating backed collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "error getting backed collection")
	suite.Require().NotNil(collection.Invariants)
	suite.Require().NotNil(collection.Invariants.CosmosCoinBackedPath)

	return collection
}

// TestBackedMintingGuardrails_FromIsBackingAddr_Deposit tests that an approval with
// allowBackedMinting=true and backing address on fromListId passes validation (deposit/withdrawal from backing).
func (suite *TestSuite) TestBackedMintingGuardrails_FromIsBackingAddr_Deposit() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collection := suite.createBackedCollection()
	backingAddr := collection.Invariants.CosmosCoinBackedPath.Address

	// Get current approvals and add a new one with backing address as fromListId
	approvals := collection.CollectionApprovals
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        "backed-deposit",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        backingAddr, // Exactly the backing address
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true,
			MustPrioritize:     true,
		},
	})

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Nil(err, "Should pass: backing address on fromListId with mustPrioritize=true")
}

// TestBackedMintingGuardrails_ToIsBackingAddr_Withdrawal tests that an approval with
// allowBackedMinting=true and backing address on toListId passes validation (withdrawal direction).
func (suite *TestSuite) TestBackedMintingGuardrails_ToIsBackingAddr_Withdrawal() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collection := suite.createBackedCollection()
	backingAddr := collection.Invariants.CosmosCoinBackedPath.Address

	approvals := collection.CollectionApprovals
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        "backed-withdrawal",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          backingAddr, // Exactly the backing address
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true,
			MustPrioritize:     true,
		},
	})

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Nil(err, "Should pass: backing address on toListId with mustPrioritize=true")
}

// TestBackedMintingGuardrails_BothSidesMatch_Rejected tests that an approval with
// allowBackedMinting=true and backing address on BOTH fromListId and toListId is rejected.
func (suite *TestSuite) TestBackedMintingGuardrails_BothSidesMatch_Rejected() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collection := suite.createBackedCollection()
	backingAddr := collection.Invariants.CosmosCoinBackedPath.Address

	approvals := collection.CollectionApprovals
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        "backed-both",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        backingAddr,
		ToListId:          backingAddr,
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true,
			MustPrioritize:     true,
		},
	})

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Error(err, "Should reject: backing address on both sides")
	suite.Require().Contains(err.Error(), "both fromListId and toListId match the backing address")
}

// TestBackedMintingGuardrails_NeitherSideMatches_Rejected tests that an approval with
// allowBackedMinting=true and backing address on NEITHER fromListId nor toListId is rejected.
func (suite *TestSuite) TestBackedMintingGuardrails_NeitherSideMatches_Rejected() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collection := suite.createBackedCollection()

	approvals := collection.CollectionApprovals
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        "backed-neither",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true,
			MustPrioritize:     true,
		},
	})

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Error(err, "Should reject: neither side is backing address")
	suite.Require().Contains(err.Error(), "neither fromListId nor toListId is exactly the backing address")
}

// TestBackedMintingGuardrails_FromListIdIsAllWithoutMint_ToIsBackingAddr tests that an approval with
// allowBackedMinting=true, fromListId="AllWithoutMint" (not the backing address), and
// toListId=backingAddr passes because exactly one side matches the backing address.
func (suite *TestSuite) TestBackedMintingGuardrails_FromListIdIsAllWithoutMint_ToIsBackingAddr() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collection := suite.createBackedCollection()
	backingAddr := collection.Invariants.CosmosCoinBackedPath.Address

	approvals := collection.CollectionApprovals
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        "backed-broad-from",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint", // Broad list, not the backing address
		ToListId:          backingAddr,       // Exactly the backing address
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true,
			MustPrioritize:     true,
		},
	})

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	// "AllWithoutMint" is not exactly the backing address and toListId IS the backing address,
	// so only toListId matches. This should pass since exactly one side (toListId) matches.
	suite.Require().Nil(err, "Should pass: toListId is exactly backing address, fromListId='AllWithoutMint' is not a match")
}

// TestBackedMintingGuardrails_FromListIdIsAll_NeitherExact_Rejected tests that an approval with
// allowBackedMinting=true and fromListId="All" and toListId="All" is rejected.
func (suite *TestSuite) TestBackedMintingGuardrails_FromListIdIsAll_NeitherExact_Rejected() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	_ = suite.createBackedCollection()

	approvals := []*types.CollectionApproval{
		{
			ApprovalId:        "backed-all-both",
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			FromListId:        "All",
			ToListId:          "All",
			InitiatedByListId: "AllWithoutMint",
			ApprovalCriteria: &types.ApprovalCriteria{
				AllowBackedMinting: true,
				MustPrioritize:     true,
			},
		},
	}

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Error(err, "Should reject: neither fromListId=All nor toListId=All is exactly the backing address")
	suite.Require().Contains(err.Error(), "neither fromListId nor toListId is exactly the backing address")
}

// TestBackedMintingGuardrails_MustPrioritizeFalse_Rejected tests that an approval with
// allowBackedMinting=true but mustPrioritize=false is rejected.
func (suite *TestSuite) TestBackedMintingGuardrails_MustPrioritizeFalse_Rejected() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collection := suite.createBackedCollection()
	backingAddr := collection.Invariants.CosmosCoinBackedPath.Address

	approvals := collection.CollectionApprovals
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        "backed-no-prioritize",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        backingAddr,
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true,
			MustPrioritize:     false, // Must be true
		},
	})

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Error(err, "Should reject: mustPrioritize is false")
	suite.Require().Contains(err.Error(), "mustPrioritize is not set to true")
}

// TestBackedMintingGuardrails_AllowBackedMintingFalse_NoValidation tests that an approval with
// allowBackedMinting=false does not trigger any of the new validations.
func (suite *TestSuite) TestBackedMintingGuardrails_AllowBackedMintingFalse_NoValidation() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	_ = suite.createBackedCollection()

	approvals := []*types.CollectionApproval{
		{
			ApprovalId:        "normal-transfer",
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			FromListId:        "AllWithoutMint",
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			ApprovalCriteria: &types.ApprovalCriteria{
				AllowBackedMinting: false, // No backed minting
			},
		},
	}

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Nil(err, "Should pass: allowBackedMinting=false should not trigger guardrail validation")
}

// TestBackedMintingGuardrails_NoBackedPath_Rejected tests that allowBackedMinting=true
// on a collection without cosmosCoinBackedPath is rejected.
func (suite *TestSuite) TestBackedMintingGuardrails_NoBackedPath_Rejected() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection WITHOUT cosmosCoinBackedPath
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	// Try to add an approval with allowBackedMinting=true
	approvals := []*types.CollectionApproval{
		{
			ApprovalId:        "backed-no-path",
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			FromListId:        "AllWithoutMint",
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			ApprovalCriteria: &types.ApprovalCriteria{
				AllowBackedMinting: true,
				MustPrioritize:     true,
			},
		},
	}

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Error(err, "Should reject: no cosmosCoinBackedPath in invariants")
	suite.Require().Contains(err.Error(), "no cosmosCoinBackedPath in invariants")
}
